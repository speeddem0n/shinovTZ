package printer

import (
	"fmt"
	"strings"
	"time"

	"github.com/speeddem0n/shinovTZ/internal/task"
)

type Printer struct {
	Header    string // Заголовок строки
	TasksChan <-chan task.Task
}

// Принтер результатов
func (p Printer) Print() {
	ticker := time.NewTicker(3 * time.Second) // Тикер на 3 секунды
	defer ticker.Stop()

	printerStrBuilder := strings.Builder{}
	printerStrBuilder.WriteString(fmt.Sprintf("%s\n", p.Header))
	for {
		select {
		case task, ok := <-p.TasksChan:
			if !ok { // Если канал закрыт, выводим финальные результаты
				fmt.Printf("%s=========================\n", printerStrBuilder.String())
				return
			}

			printerStrBuilder.WriteString(fmt.Sprintf("%s\n", task.String()))
		case <-ticker.C: // Выводим результаты раз в 3 секунды
			fmt.Printf("%s=========================\n", printerStrBuilder.String())
		}
	}
}
