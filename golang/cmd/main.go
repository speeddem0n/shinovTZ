package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/speeddem0n/shinovTZ/internal/generator"
	"github.com/speeddem0n/shinovTZ/internal/printer"
	"github.com/speeddem0n/shinovTZ/internal/task"
	"github.com/speeddem0n/shinovTZ/internal/worker"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// Мы даем тестовое задание чтобы:
// * уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших ребят;
// * увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
// * снизить число коротких собеседований, когда мы отказываем сразу же.

// Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

// Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

const (
	workerPoolSize = 5
)

func main() {
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	// Запуск генератора задач
	generator := generator.Generator{
		TasksPoolSize: 10,
	}
	tasksChan := generator.Generate(ctx)

	// Основные каналы и группы ожидания
	successChan := make(chan task.Task, generator.TasksPoolSize/2)
	errorChan := make(chan task.Task, generator.TasksPoolSize/2)
	wg := &sync.WaitGroup{}

	go func() {
		for i := 0; i < workerPoolSize; i++ {
			wg.Add(1)

			go worker.Worker{
				TasksChan:        tasksChan,
				SuccessTasksChan: successChan,
				ErrorTasksChan:   errorChan,
			}.Work(wg)
		}
	}()

	// Запуск принтера успешных задач
	successfulTasksPrinter := printer.Printer{
		Header:    "==== Successful Tasks ===",
		TasksChan: successChan,
	}
	go successfulTasksPrinter.Print()

	// Запуск принтера ошибочных задач
	erroredTasksPrinter := printer.Printer{
		Header:    "===== Errored Tasks =====",
		TasksChan: errorChan,
	}
	go erroredTasksPrinter.Print()

	// Работать ровно 10 секунд
	time.Sleep(10 * time.Second)
	ctxCancelFn()

	// Ожидаем завершения всех воркеров
	wg.Wait()
	close(successChan) // Закрываем канал успешных задач после завершения воркеров
	close(errorChan)   // Закрываем канал ошибочных задач после завершения воркеров

	fmt.Println("All tasks processed and results printed.")
}
