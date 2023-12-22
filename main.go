package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nats-io/stan.go"
)

const clusterID = "test-cluster"
const clientID = "producer-client"
const subject = "test-subject"

func getEnvWithDefault(key string, defaultVal int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}

	return intValue
}

func connectToNATS(port int) (stan.Conn, error) {
	clusterID := "yourClusterID" // Замените на ваш clusterID
	clientID := "yourClientID"   // Замените на ваш clientID

	connectionString := fmt.Sprintf("nats://host.docker.internal:%d", port)
	return stan.Connect(clusterID, clientID, stan.NatsURL(connectionString))
}

func errorHandler(natsError error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, fmt.Sprintf("Error connecting to NATS Streaming: %v", natsError), http.StatusInternalServerError)
	}
}

func main() {
	port := getEnvWithDefault("PORT", 4222)

	sc, err := connectToNATS(port)
	if err != nil {
		log.Printf("Failed to connect to NATS Streaming: %v", err)

		http.HandleFunc("/", errorHandler(err))
		log.Printf("Starting server on http://localhost:%d", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	} else {
		defer sc.Close()
		// Тут ваша логика при успешном подключении к NATS
	}

	// for {
	// 	// Чтение файла
	// 	data, err := ioutil.ReadFile("./ord1.json")
	// 	if err != nil {
	// 		log.Fatalf("Ошибка при чтении файла: %v", err)
	// 	}

	// 	// Декодирование JSON
	// 	var order Order
	// 	if err := json.Unmarshal(data, &order); err != nil {
	// 		log.Fatalf("Ошибка декодирования JSON: %v", err)
	// 	}

	// 	// Изменение order_uid
	// 	order.OrderUID = incrementOrderUID(order.OrderUID)

	// 	// Кодирование обратно в JSON
	// 	modifiedData, err := json.Marshal(order)
	// 	if err != nil {
	// 		log.Fatalf("Ошибка кодирования JSON: %v", err)
	// 	}

	// 	// Запись обратно в файл
	// 	if err := ioutil.WriteFile("./ord1.json", modifiedData, 0644); err != nil {
	// 		log.Fatalf("Ошибка при записи в файл : %v", err)
	// 	}

	// 	// Отправка сообщения
	// 	if err := sc.Publish(subject, modifiedData); err != nil {
	// 		log.Fatalf("Ошибка при отправке сообщения : %v", err)
	// 	}
	// 	fmt.Println(" Message sent! ")

	// 	// Пауза в 5 секунд
	// 	time.Sleep(5 * time.Second)
	// }
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

func processFile(filePath string) {
	// Чтение файла
	data, err := ioutil.ReadFile(filePath)
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
	if err := ioutil.WriteFile(filePath, modifiedData, 0644); err != nil {
		log.Fatalf("Ошибка при записи в файл: %v", err)
	}
}
