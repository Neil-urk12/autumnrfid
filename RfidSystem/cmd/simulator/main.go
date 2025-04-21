package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}

	for {
		// Simulate RFID scan
		rfidData := map[string]string{
			"rfid": "ACLC-2023-001",
		}

		jsonData, err := json.Marshal(rfidData)
		if err != nil {
			fmt.Printf("Error marshaling data: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/card-scan", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}

		fmt.Printf("Response status: %s\n", resp.Status)
		resp.Body.Close()

		// Wait for 5 seconds before next simulation
		time.Sleep(5 * time.Second)
	}
}
