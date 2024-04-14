package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func makeRequest(url string, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		<-ch

		// Отключение проверки сертификата
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Response: %v\n", resp.Status)
	}
}

func main() {
	url := "https://localhost:9000/user_banner?tag_id=1&feature_id=2&use_last_revision=false"
	numRequests := 1000
	duration := 1 * time.Second // Длительность теста 1 секунда

	var wg sync.WaitGroup
	ch := make(chan bool, numRequests)

	// Инициализация горутин
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(url, ch, &wg)
	}

	// Запуск горутин с интервалом времени, чтобы добиться 1000 RPS
	start := time.Now()
	for time.Since(start) < duration {
		ch <- true
		time.Sleep(time.Second / time.Duration(numRequests))
	}

	close(ch) // Закрываем канал, чтобы завершить горутины
	wg.Wait() // Ожидаем завершения всех горутин
}
