package printer

import (
	"fmt"
	"strings"
	"time"

	"github.com/speeddem0n/shinovTZ/internal/task"
)

type Printer struct {
	Header    string
	TasksChan <-chan task.Task
}

func (p Printer) Print() {
	ticker := time.NewTicker(3 * time.Second) // Тикер на 3 секунды
	defer ticker.Stop()

	printerStrBuilder := strings.Builder{}
	printerStrBuilder.WriteString(fmt.Sprintf("%s\n", p.Header))
	for {
		select {
		case task, ok := <-p.TasksChan:
			if !ok { // Если канал закрыт, выводим финальные ошибочные задачи
				fmt.Printf("%s=========================\n", printerStrBuilder.String())
				return
			}

			printerStrBuilder.WriteString(fmt.Sprintf("%s\n", task.String()))
		case <-ticker.C: // Периодически выводим промежуточные результаты
			fmt.Printf("%s=========================\n", printerStrBuilder.String())
		}
	}
}
