package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile("status.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	for {
		status, err := getStatus()
		if err != nil {
			logger.Printf("Error getting status: %v", err)
		} else {
			logger.Printf("Temperature: %.2f, Battery Level: %d%%, Load: %.2f",
				status.Result.System.MCU.Temperature,
				status.Result.System.MCU.ChargePercent,
				status.Result.System.LoadAverage[0])
			fmt.Printf("Temperature: %.2f, Battery Level: %d%%, Load: %.2f\n",
				status.Result.System.MCU.Temperature,
				status.Result.System.MCU.ChargePercent,
				status.Result.System.LoadAverage[0])
		}
		time.Sleep(1 * time.Minute)
	}

}

func getStatus() (*glStatusResponse, error) {
	//url := "http://192.168.8.1/rpc"
	url := "http://127.0.0.1/rpc"

	commands := []string{"system", "get_status"}

	body, statusCode, err := mcuHttpQuery(url, commands)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	if statusCode >= 200 && statusCode <= 299 {
		// Here we can now parse the returned value...

		// First, let's check if there's an error present. The RPC call always returns 200.
		// Even if for example authorization failed, it will return 200 with error messages.
		err = parseError(body)
		if err != nil {
			fmt.Println("HTTP request successful but error in return value: ", err)
			return nil, err
		}

		// Call was successful and no error messages in the returned body. Let's parse the response into an object.
		status, err := ParseGetStatusMsg(string(body))
		if err != nil {
			fmt.Println("HTTP request successful but error parsing message: ", err)
			return nil, err
		}
		return status, nil

	} else {
		fmt.Println("getStatus Request failed with Code: ", statusCode, "and Msg: ", body)
		return nil, err
	}
}

func mcuHttpQuery(url string, commands []string) ([]byte, int, error) {

	// first, construct the http request. It's fairly simple, they always follow this pattern:
	// '{"jsonrpc":"2.0","id":1,"method":"call","params":["<token>","<param 1>","<param 2>",...,"<param N>"]}'
	// Since this is sent from the host system (localhost), a security token is not needed.
	req, err := makeQuery(url, commands)
	if err != nil {
		return nil, -1, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	return body, resp.StatusCode, nil
}

func makeQuery(url string, commands []string) (*http.Request, error) {
	params := string(``)
	start := string(`{"jsonrpc":"2.0","id":1,"method":"call","params":[""`)
	end := string(`]}`)

	for _, cmd := range commands {
		params = params + string(fmt.Sprintf(`,"%s"`, cmd))
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(start+params+end)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("glinet", "1")
	req.Header.Set("Content-Type", "application/json")

	return req, err
}
