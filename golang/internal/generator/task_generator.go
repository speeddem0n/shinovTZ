package generator

import (
	"context"
	"time"

	"github.com/speeddem0n/shinovTZ/internal/task"
)

type Generator struct {
	TasksPoolSize int
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

				// Условие для ошибочных задач (Если ОС Windows то точнось вычесления Nanosecond оставляем желать лучшего, и ошибочные таски не появляются)
				if time.Now().Nanosecond()%2 > 0 {
					creationTime = "Some error occurred"
				}

				// Передаем задачу в канал
				tasksCh <- task.Task{
					ID:        id,
					CreatedAt: creationTime,
				}

				id++
				time.Sleep(time.Millisecond * 250)
			}
		}
	}()

	return tasksCh
}
