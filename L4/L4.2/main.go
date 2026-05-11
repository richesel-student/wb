package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// --- 1. Конфигурация и Структуры ---

type Options struct {
	Mode       string
	Address    string
	Workers    string
	CountOnly  bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	ShowLine   bool
}

type GrepArgs struct {
	TaskID     int
	Pattern    string
	ChunkData  []byte // Raw bytes для минимизации аллокаций памяти (защита от OOM)
	StartLine  int
	CountOnly  bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
}

type GrepReply struct {
	TaskID       int
	MatchedLines []string
	MatchCount   int
}

type TaskContext struct {
	Args    GrepArgs
	Retries int
}

const (
	MaxRetries   = 3
	ChunkSize    = 1024 * 1024 // 1 MB чанки
	RetryBackoff = 2 * time.Second
)

// --- 2. Воркер (Рабочий Узел) ---

type WorkerService struct {
	mu         sync.Mutex
	regexCache map[string]*regexp.Regexp
}

func NewWorkerService() *WorkerService {
	return &WorkerService{
		regexCache: make(map[string]*regexp.Regexp),
	}
}

// getRegex безопасно (с Double-Checked Locking) возвращает скомпилированное выражение
func (w *WorkerService) getRegex(pattern string, ignoreCase bool, fixed bool) (*regexp.Regexp, error) {
	key := pattern
	if ignoreCase {
		key = "i:" + key
	}
	if fixed {
		key = "f:" + key
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if re, ok := w.regexCache[key]; ok {
		return re, nil
	}

	expr := pattern
	if fixed {
		expr = regexp.QuoteMeta(pattern)
	}
	if ignoreCase {
		expr = "(?i)" + expr
	}

	re, err := regexp.Compile(expr)
	if err == nil {
		w.regexCache[key] = re
	}
	return re, err
}

func (w *WorkerService) Process(args *GrepArgs, reply *GrepReply) error {
	reply.TaskID = args.TaskID

	var re *regexp.Regexp
	var err error
	var patternBytes []byte

	// Для точного поиска с учетом регистра используем bytes.Contains (Zero-allocation)
	useBytesContains := args.Fixed && !args.IgnoreCase

	if useBytesContains {
		patternBytes = []byte(args.Pattern)
	} else {
		re, err = w.getRegex(args.Pattern, args.IgnoreCase, args.Fixed)
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(bytes.NewReader(args.ChunkData))
	lineIndex := 0

	for scanner.Scan() {
		lineBytes := scanner.Bytes() // Читаем байты напрямую без копирования в string

		var ok bool
		if useBytesContains {
			ok = bytes.Contains(lineBytes, patternBytes)
		} else {
			ok = re.Match(lineBytes)
		}

		if args.Invert {
			ok = !ok
		}

		if ok {
			reply.MatchCount++
			if !args.CountOnly {
				// Аллокация новой строки происходит ТОЛЬКО при успешном совпадении
				out := string(lineBytes)
				if args.StartLine > 0 {
					out = fmt.Sprintf("%d:%s", args.StartLine+lineIndex, out)
				}
				reply.MatchedLines = append(reply.MatchedLines, out)
			}
		}
		lineIndex++
	}

	return scanner.Err()
}

// --- 3. Координатор (Мастер) ---

func runMaster(ctx context.Context, cancel context.CancelFunc, opts Options, pattern, filename string) {
	workerAddrs := strings.Split(opts.Workers, ",")
	totalWorkers := len(workerAddrs)
	if totalWorkers == 0 || workerAddrs[0] == "" {
		log.Fatal("Для режима master необходимо указать --workers")
	}

	quorum := int32(totalWorkers/2 + 1)
	var activeWorkers int32 = 0 // Атомарный счетчик живых узлов

	tasks := make(chan TaskContext, 100)
	results := make(chan *GrepReply, 100)
	var wg sync.WaitGroup

	// 1. Запуск пула соединений с воркерами
	for _, addr := range workerAddrs {
		go workerConnector(ctx, cancel, addr, tasks, results, &wg, &activeWorkers, quorum)
	}

	// 2. Механизм прогрева (Initial Quorum Check)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			currentActive := atomic.LoadInt32(&activeWorkers)
			if currentActive < quorum {
				log.Printf("ФАТАЛЬНАЯ ОШИБКА: Не удалось собрать стартовый кворум за 5 секунд (активно %d, требуется %d)", currentActive, quorum)
				cancel()
			}
		}
	}()

	var readerWg sync.WaitGroup // WaitGroup для отслеживания завершения чтения файла
	readerWg.Add(1)

	// 3. Чтение файла и генерация задач
	go func() {
		defer readerWg.Done() // Сигнализируем, что весь файл прочитан и разбит на чанки

		var reader io.Reader
		if filename != "" {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatalf("Ошибка открытия файла: %v", err)
			}
			defer file.Close()
			reader = file
		} else {
			reader = os.Stdin
		}

		bufReader := bufio.NewReader(reader)
		var chunkBuffer bytes.Buffer
		taskID := 0
		startLine := 1
		linesInChunk := 0

		for {
			select {
			case <-ctx.Done():
				return // Немедленный выход, если система падает
			default:
				line, err := bufReader.ReadBytes('\n')
				if len(line) > 0 {
					chunkBuffer.Write(line)
					linesInChunk++
				}

				if chunkBuffer.Len() >= ChunkSize || err != nil {
					if chunkBuffer.Len() > 0 {
						data := make([]byte, chunkBuffer.Len())
						copy(data, chunkBuffer.Bytes())

						wg.Add(1) // Добавляем задачу для воркеров

						// Безопасная отправка задачи
						select {
						case <-ctx.Done():
							wg.Done() // Откатываем счетчик, так как задача не ушла
							return
						case tasks <- TaskContext{
							Retries: 0,
							Args: GrepArgs{
								TaskID:     taskID,
								Pattern:    pattern,
								ChunkData:  data,
								StartLine:  startLine,
								CountOnly:  opts.CountOnly,
								IgnoreCase: opts.IgnoreCase,
								Invert:     opts.Invert,
								Fixed:      opts.Fixed,
							},
						}:
						}

						taskID++
						startLine += linesInChunk
						chunkBuffer.Reset()
						linesInChunk = 0
					}
				}

				if err == io.EOF {
					return
				} else if err != nil {
					log.Printf("Ошибка чтения файла: %v", err)
					cancel()
					return
				}
			}
		}
	}()

	// 4. Ждем выполнения всех задач
	go func() {
		readerWg.Wait() // Ждем, пока горутина-читатель закончит плодить задачи
		wg.Wait()       // Ждем, пока воркеры выполнят все эти задачи
		close(results)  // Закрываем канал (сигнал успешного завершения)
	}()

	// 5. Сбор результатов
	var finalReplies []*GrepReply
	var aborted bool

loop:
	for {
		select {
		case <-ctx.Done():
			aborted = true // Контекст был отменен из-за потери кворума
			break loop
		case reply, ok := <-results:
			if !ok {
				break loop // Канал закрыт штатно, все результаты получены
			}
			finalReplies = append(finalReplies, reply)
		}
	}

	// 6. Проверка статуса завершения
	if aborted {
		log.Fatal("Выполнение прервано из-за потери сети (кворум нарушен) или фатальной ошибки.")
	}

	// 7. Упорядочиваем результаты
	sort.Slice(finalReplies, func(i, j int) bool {
		return finalReplies[i].TaskID < finalReplies[j].TaskID
	})

	totalCount := 0
	for _, reply := range finalReplies {
		if opts.CountOnly {
			totalCount += reply.MatchCount
		} else {
			for _, line := range reply.MatchedLines {
				fmt.Print(line)
			}
		}
	}

	if opts.CountOnly {
		fmt.Println(totalCount)
	}
}

// workerConnector управляет подключением, мониторит кворум и безопасно возвращает задачи
func workerConnector(
	ctx context.Context, cancel context.CancelFunc, addr string,
	tasks chan TaskContext, results chan<- *GrepReply, wg *sync.WaitGroup,
	activeWorkers *int32, quorum int32,
) {
	var client *rpc.Client
	var err error
	isConnected := false

	for {
		// Логика подключения
		if client == nil {
			client, err = rpc.Dial("tcp", addr)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case <-time.After(RetryBackoff):
					continue
				}
			}

			// Успешное подключение
			atomic.AddInt32(activeWorkers, 1)
			isConnected = true
			
		}

		select {
		case <-ctx.Done():
			if client != nil {
				client.Close()
			}
			return
		case taskContext, ok := <-tasks:
			if !ok {
				return
			}

			reply := new(GrepReply)
			err = client.Call("WorkerService.Process", &taskContext.Args, reply)

			if err != nil {
				// Разрыв соединения
				if isConnected {
					currentActive := atomic.AddInt32(activeWorkers, -1)
					isConnected = false
					if currentActive < quorum {
						log.Printf("КРИТИЧЕСКАЯ ОШИБКА: Потерян кворум (активно %d, требуется %d)", currentActive, quorum)
						cancel()
					}
				}

				client.Close()
				client = nil

				taskContext.Retries++
				if taskContext.Retries > MaxRetries {
					log.Printf("ЗАДАЧА %d ПРОВАЛЕНА окончательно.", taskContext.Args.TaskID)
					wg.Done()
					cancel()
					continue
				}

				// Возврат задачи
				go func(tc TaskContext) {
					select {
					case <-ctx.Done():
						return
					case tasks <- tc:
					}
				}(taskContext)

				continue
			}

			// Успех
			results <- reply
			wg.Done()
		}
	}
}

// --- 4. Точка входа ---

func parseFlags() Options {
	opts := Options{}
	flag.StringVar(&opts.Mode, "mode", "master", "Режим: master или worker")
	flag.StringVar(&opts.Address, "address", ":8080", "Адрес для worker")
	flag.StringVar(&opts.Workers, "workers", "", "Адреса воркеров через запятую (только для master)")

	flag.BoolVar(&opts.CountOnly, "c", false, "count only")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "ignore case")
	flag.BoolVar(&opts.Invert, "v", false, "invert match")
	flag.BoolVar(&opts.Fixed, "F", false, "fixed string")
	flag.BoolVar(&opts.ShowLine, "n", false, "show line numbers")
	flag.Parse()

	return opts
}

func main() {
	opts := parseFlags()

	if opts.Mode == "worker" {
		worker := NewWorkerService()
		err := rpc.Register(worker)
		if err != nil {
			log.Fatalf("Ошибка регистрации RPC: %v", err)
		}

		listener, err := net.Listen("tcp", opts.Address)
		if err != nil {
			log.Fatalf("Ошибка запуска на %s: %v", opts.Address, err)
		}
		
		rpc.Accept(listener)
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Использование: mygrep [flags] pattern [file]")
		os.Exit(1)
	}

	pattern := args[0]
	filename := ""
	if len(args) >= 2 {
		filename = args[1]
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runMaster(ctx, cancel, opts, pattern, filename)
}