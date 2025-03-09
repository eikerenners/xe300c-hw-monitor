package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "http://192.168.8.1/rpc"
	//url := "http://127.0.0.1/rpc"
	jsonData := []byte(`{"jsonrpc":"2.0","id":1,"method":"call","params":["","system","get_status"]}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("glinet", "1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Request was successful: %d\n", resp.StatusCode)
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}

	// add a check for the "code" field inside the returned json object. If it's equal to -32000 pring "access denied" otherwise print the entire object. AI!

	fmt.Println("Response:", result)
}
