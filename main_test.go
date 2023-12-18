package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIncrementOrderUID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"BasicIncrement", "123", "124"},
		{"ZeroValue", "0", "1"},
		{"LargeNumber", "9999", "10000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := incrementOrderUID(tc.input)
			if got != tc.expected {
				t.Errorf("incrementOrderUID(%v) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestFileReadAndWrite(t *testing.T) {
	// Создание временного файла
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // очистка

	// Подготовка и запись исходных данных
	initialData := []byte(`{"order_uid": "100"}`)
	if err := ioutil.WriteFile(tmpfile.Name(), initialData, 0666); err != nil {
		t.Fatal(err)
	}

	processFile(tmpfile.Name())

	// Чтение и проверка результата
	resultData, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	expectedData := `{"order_uid":"101","track_number":"","entry":"","delivery":{"name":"","phone":"","zip":"","city":"","address":"","region":"","email":""},"payment":{"transaction":"","request_id":"","currency":"","provider":"","amount":0,"payment_dt":0,"bank":"","delivery_cost":0,"goods_total":0,"custom_fee":0},"items":null,"locale":"","internal_signature":"","customer_id":"","delivery_service":"","shardkey":"","sm_id":0,"date_created":"0001-01-01T00:00:00Z","oof_shard":""}`
	if string(resultData) != expectedData {
		t.Errorf("Expected file content %v, got %v", expectedData, string(resultData))
	}
}
