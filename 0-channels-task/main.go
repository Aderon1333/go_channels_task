// Новая ветка "for_code_review"
package main

import (
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"
)

// Тип "пара", чтобы удобно хранить значения для четных чисел
type Pair struct {
	value   int
	comment string
}

// Горутина, заполняющая канал входных значений
func inputStreamCreator(inputChannel chan<- *int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		loc_int := rand.Intn(2000) // 2000 - чтобы числа не были огромными
		inputChannel <- &loc_int
	}
}

// Горутина, разделяющая входной канал
func channelSplitter(inputChannel <-chan *int, outputChannel chan<- interface{}) {
	for loc_i := range inputChannel {
		if *loc_i%2 == 0 {
			newPair := &Pair{value: 0, comment: ""}
			if *loc_i > 1000 {
				newPair.value = *loc_i
				newPair.comment = "big"
			} else {
				newPair.value = *loc_i
				newPair.comment = "small"
			}
			outputChannel <- newPair
		} else if *loc_i%3 == 0 {
			outputChannel <- loc_i
		}
	}
}

// Горутина, которая выводит канал с результатами
func resultCreator(outputChannel <-chan interface{}) {
	for channelValue := range outputChannel {
		switch channelValue := channelValue.(type) {
		case *int:
			fmt.Printf("Type int: %d\n", *channelValue)
		case *Pair:
			fmt.Printf("Type pair{value, comment}: %d, %s\n", channelValue.value, channelValue.comment)
		default:
			panic("Unexpected type!")
		}
	}
}

func main() {
	// Входной канал
	inputChannel := make(chan *int)

	// Выходной канал
	outputChannel := make(chan interface{})

	// Канал для обработки сигналов от ОС
	sigCh := make(chan os.Signal, 1)

	// Канал для завершения процесса
	done := make(chan bool, 1)

	// Флаг того, что процесс в завершении
	var stopping bool

	go inputStreamCreator(inputChannel)

	go channelSplitter(inputChannel, outputChannel)

	go resultCreator(outputChannel)

	// Не уверен, что правильно понял gracefull shutdown, и возможно усложнил, но вроде работает
	go func() {
		for {
			sig := <-sigCh
			if sig == syscall.SIGINT {
				if stopping {
					os.Exit(1)
				} else {
					stopping = true
					go func() {
						done <- true
					}()
				}
			}
		}
	}()

	<-done
}
