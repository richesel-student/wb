package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// текущая выполняемая команда (нужно для обработки Ctrl+C)
var currentCmd *exec.Cmd

func main() {

	// чтение команд из стандартного ввода
	reader := bufio.NewReader(os.Stdin)

	// канал для получения сигнала Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// обработка Ctrl+C
	go func() {
		for range sig {
			if currentCmd != nil && currentCmd.Process != nil {
				currentCmd.Process.Signal(syscall.SIGINT)
			}
			fmt.Println()
		}
	}()

	// основной цикл shell
	for {
		fmt.Print("minishell> ")

		line, err := reader.ReadString('\n')

		// завершение shell при Ctrl+D
		if err == io.EOF {
			fmt.Println("exit")
			return
		}

		line = strings.TrimSpace(line)

		// пропускаем пустую строку
		if line == "" {
			continue
		}

		// подстановка переменных окружения
		line = os.ExpandEnv(line)

		runConditional(line)
	}
}

// обработка условных операторов && и ||
func runConditional(line string) {

	if strings.Contains(line, "&&") {

		parts := strings.Split(line, "&&")

		for _, p := range parts {
			if !execute(strings.TrimSpace(p)) {
				return
			}
		}

		return
	}

	if strings.Contains(line, "||") {

		parts := strings.Split(line, "||")

		for _, p := range parts {
			if execute(strings.TrimSpace(p)) {
				return
			}
		}

		return
	}

	execute(line)
}

// выполнение одной команды
func execute(line string) bool {

	// если есть конвейер — запускаем pipeline
	if strings.Contains(line, "|") {
		return executePipeline(line)
	}

	// разбиваем строку на аргументы
	args := strings.Fields(line)

	if len(args) == 0 {
		return true
	}

	// проверяем встроенные команды
	switch args[0] {

	case "cd":
		return builtinCd(args)

	case "pwd":
		return builtinPwd()

	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
		return true

	case "kill":
		return builtinKill(args)

	case "ps":
		return runExternal("ps", []string{"aux"}, "", "", false)

	case "exit":
		os.Exit(0)
	}

	// запуск внешней команды
	return runCommand(args)
}

// встроенная команда cd
func builtinCd(args []string) bool {

	if len(args) < 2 {
		fmt.Println("cd: отсутствует путь")
		return false
	}

	err := os.Chdir(args[1])

	if err != nil {
		fmt.Println("cd:", err)
		return false
	}

	return true
}

// встроенная команда pwd
func builtinPwd() bool {

	dir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(dir)

	return true
}

// встроенная команда kill
func builtinKill(args []string) bool {

	if len(args) < 2 {
		fmt.Println("kill: требуется PID")
		return false
	}

	pid, err := strconv.Atoi(args[1])

	if err != nil {
		fmt.Println("неверный PID")
		return false
	}

	proc, err := os.FindProcess(pid)

	if err != nil {
		fmt.Println(err)
		return false
	}

	// отправляем сигнал завершения процессу
	proc.Signal(syscall.SIGKILL)

	return true
}

// разбор редиректов ввода и вывода (<, >, >>)
func parseRedirects(args []string) ([]string, string, string, bool) {

	var input string
	var output string
	appendMode := false

	var cleaned []string

	for i := 0; i < len(args); i++ {

		switch args[i] {

		case ">":
			output = args[i+1]
			i++

		case ">>":
			output = args[i+1]
			appendMode = true
			i++

		case "<":
			input = args[i+1]
			i++

		default:
			cleaned = append(cleaned, args[i])
		}
	}

	return cleaned, input, output, appendMode
}

// запуск команды с учетом редиректов
func runCommand(args []string) bool {

	args, input, output, appendMode := parseRedirects(args)

	return runExternal(args[0], args[1:], input, output, appendMode)
}

// запуск внешней программы
func runExternal(cmd string, args []string, input string, output string, appendMode bool) bool {

	command := exec.Command(cmd, args...)

	// настройка входного потока
	if input != "" {

		file, err := os.Open(input)
		if err != nil {
			fmt.Println(err)
			return false
		}

		command.Stdin = file

	} else {

		command.Stdin = os.Stdin
	}

	// настройка выходного потока
	if output != "" {

		var file *os.File
		var err error

		if appendMode {
			file, err = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			file, err = os.Create(output)
		}

		if err != nil {
			fmt.Println(err)
			return false
		}

		command.Stdout = file

	} else {

		command.Stdout = os.Stdout
	}

	command.Stderr = os.Stderr

	// сохраняем текущую команду для обработки Ctrl+C
	currentCmd = command

	err := command.Run()

	currentCmd = nil

	if err != nil {
		return false
	}

	return true
}

// выполнение конвейера команд (|)
func executePipeline(line string) bool {

	parts := strings.Split(line, "|")

	var cmds []*exec.Cmd

	for _, p := range parts {

		args := strings.Fields(strings.TrimSpace(p))

		if len(args) == 0 {
			return false
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr

		cmds = append(cmds, cmd)
	}

	// соединяем вывод одной команды со входом следующей
	for i := 0; i < len(cmds)-1; i++ {

		pipe, err := cmds[i].StdoutPipe()

		if err != nil {
			fmt.Println(err)
			return false
		}

		cmds[i+1].Stdin = pipe
	}

	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout

	// запускаем все команды
	for _, cmd := range cmds {

		currentCmd = cmd

		err := cmd.Start()

		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	// ожидаем завершения
	for _, cmd := range cmds {
		cmd.Wait()
	}

	currentCmd = nil

	return true
}