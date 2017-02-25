package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {

}

func requestToURL(url string) {
	jsonInit := "["
	for i := 0; i < 1999; i++ {
		jsonInit += "{\"name\":\"test\",\"year\":2016,\"author\":\"Myself\"},"
	}
	jsonInit += "{\"name\":\"test\",\"year\":2016,\"author\":\"Myself\"}]"
	jsonStr := []byte(jsonInit)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func main() {

	for i := 0; i < 50; i++ {
		fmt.Printf("Starting clients for server:%d\n", i+1)
		wg.Add(1)
		go func(serverID int) {
			for j := 0; j < 500; j++ {

				requestToURL(fmt.Sprintf("http://iti%d:8%03d/insertmany", serverID, serverID))

				time.Sleep(10 * time.Millisecond)

			}
			wg.Done()
		}(i + 1)
	}

	wg.Wait()
	fmt.Printf("Done")
}
