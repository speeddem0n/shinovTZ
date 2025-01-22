package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/speeddem0n/shinovTZ/internal/task"
)

type Worker struct {
	TasksChan        <-chan task.Task
	SuccessTasksChan chan<- task.Task
	ErrorTasksChan   chan<- task.Task
}

// Воркер
func (w Worker) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	// Считываем задачи из канала taskChan
	for task := range w.TasksChan {
		// Проверяем задачу на наличие ошибки
		parsedTime, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil {
			task.Result = []byte(fmt.Sprintf("failed to parse date: %s", err.Error()))
			task.FinishedAt = time.Now().Format(time.RFC3339)
			w.ErrorTasksChan <- task
			continue
		}

		// Проверяем, не истекло ли время выполнения задачи
		if parsedTime.Before(time.Now().Add(-20 * time.Second)) {
			task.Result = []byte("task timeout error")
			task.FinishedAt = time.Now().Format(time.RFC3339)
			w.ErrorTasksChan <- task
			continue
		}

		// Возвращаем успушную задачу
		task.Result = []byte("task has been succeeded")
		task.FinishedAt = time.Now().Format(time.RFC3339)
		w.SuccessTasksChan <- task
	}
}
