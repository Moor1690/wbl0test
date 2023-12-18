package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
)

const clusterID = "test-cluster"
const clientID = "producer-client"
const subject = "test-subject"

func main() {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://host.docker.internal:4222"))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	for {
		// Чтение файла
		data, err := ioutil.ReadFile("./ord1.json")
		if err != nil {
			log.Fatalf("Ошибка при чтении файла: %v", err)
		}

		// Декодирование JSON
		var order Order
		if err := json.Unmarshal(data, &order); err != nil {
			log.Fatalf("Ошибка декодирования JSON: %v", err)
		}

		// Изменение order_uid
		order.OrderUID = incrementOrderUID(order.OrderUID)

		// Кодирование обратно в JSON
		modifiedData, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("Ошибка кодирования JSON: %v", err)
		}

		// Запись обратно в файл
		if err := ioutil.WriteFile("./ord1.json", modifiedData, 0644); err != nil {
			log.Fatalf("Ошибка при записи в файл : %v", err)
		}

		// Отправка сообщения
		if err := sc.Publish(subject, modifiedData); err != nil {
			log.Fatalf("Ошибка при отправке сообщения: %v", err)
		}
		fmt.Println("Message sent!")

		// Пауза в 5 секунд
		time.Sleep(5 * time.Second)
	}
}

func incrementOrderUID(uid string) string {
	// Преобразование строки в число
	num, err := strconv.Atoi(uid)
	if err != nil {
		log.Fatalf("Ошибка при преобразовании order_uid в число: %v", err)
	}

	// Увеличение числа на 1
	num++

	// Преобразование числа обратно в строку
	return strconv.Itoa(num)
}
