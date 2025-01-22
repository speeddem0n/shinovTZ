package task

import "fmt"

type Task struct {
	ID         int    // ID задачи
	CreatedAt  string // Время создания задачи
	FinishedAt string // Время выполнения задачи
	Result     []byte // Результат
}

func (t Task) String() string {
	return fmt.Sprintf("ID: %d, CreatedAt: %s, FinishedAt: %s, Result: %s", t.ID, t.CreatedAt, t.FinishedAt, t.Result)
}
