#!/bin/bash

# Настройка безопасного выхода: убиваем воркеры при любом прерывании скрипта
trap 'echo "Очистка..."; kill $W1 $W2 $W3 $W4 2>/dev/null; rm -f test_1gb.txt; exit' INT TERM EXIT

# 1. Компилируем утилиту
echo "Компиляция mygrep..."
go build -o mygrep main.go

# 2. Генерируем тестовый файл (~1 ГБ = ~22 млн строк)
# Использование seq и awk работает в десятки раз быстрее, чем цикл for в bash
echo "Генерация тестового файла (~1GB)... Это займет несколько секунд."
rm -f test_1gb.txt
seq 1 22000000 | awk '{ 
    if ($1 % 10000 == 0) 
        print "Line " $1 ": CRITICAL ERROR: Connection lost" 
    else 
        print "Line " $1 ": System is running normally" 
}' > test_1gb.txt

echo "Файл сгенерирован. Размер:"
du -h test_1gb.txt

# 3. Поднимаем 4 воркера в фоне
echo "Запуск воркеров..."
./mygrep --mode=worker --address=127.0.0.1:9001 > worker1.log 2>&1 &
W1=$!
./mygrep --mode=worker --address=127.0.0.1:9002 > worker2.log 2>&1 &
W2=$!
./mygrep --mode=worker --address=127.0.0.1:9003 > worker3.log 2>&1 &
W3=$!
./mygrep --mode=worker --address=127.0.0.1:9004 > worker4.log 2>&1 &
W4=$!

sleep 2 # Ждем инициализации сети

# 4. Бенчмарк системного grep
echo "======================================"
echo "Тест 1: Системный GNU grep"
time grep -c "CRITICAL ERROR" test_1gb.txt
echo "======================================"

# 5. Бенчмарк распределённого mygrep
echo "Тест 2: Распределённый mygrep (4 воркера, Идеальные условия)"
time ./mygrep --mode=master --workers=127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003,127.0.0.1:9004 -c "CRITICAL ERROR" test_1gb.txt
echo "======================================"

# 6. Chaos Engineering: Тест отказоустойчивости (Убиваем воркер в процессе)
echo "Тест 3: Распределённый mygrep (Отказоустойчивость)"
echo "Убиваем Воркер 4 ($W4) прямо сейчас..."
kill $W4

echo "Запускаем поиск на оставшихся 3 воркерах (Кворум = 3, должно работать)..."
time ./mygrep --mode=master --workers=127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003,127.0.0.1:9004 -c "CRITICAL ERROR" test_1gb.txt
echo "======================================"

# Очистка произойдет автоматически благодаря trap в начале скрипта
echo "Тесты завершены успешно!"