package generator

import (
	"context"
	"time"

	"github.com/speeddem0n/shinovTZ/internal/task"
)

type Generator struct {
	TasksPoolSize int // Размер буфера для канала задач
}

func (g Generator) Generate(ctx context.Context) <-chan task.Task {
	tasksCh := make(chan task.Task, g.TasksPoolSize)

	go func() {
		id := 1
		for {
			select {
			case <-ctx.Done():
				close(tasksCh)
				return

			default:
				// Время создания задачи
				creationTime := time.Now().Format(time.RFC3339)

				// Условие для ошибочных задач (Nanosconds заменено на Seconds т.к. Windwos не обладает достаточной точностью вычисления наносекунд, и всегда оставляет на конце два ноля)
				if time.Now().Second()%2 > 0 {
					creationTime = "Some error occurred"
				}

				tasksCh <- task.Task{
					ID:        id,
					CreatedAt: creationTime,
				}

				id++
				time.Sleep(time.Millisecond * 150) // Симуляция работы
			}
		}
	}()

	return tasksCh
}
