package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type glError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type glErrorObject struct {
	Name    int     `json:"name"`
	Jsonrpc string  `json:"jsonrpc"`
	Error   glError `json:"error"`
}

func main() {
	// add code to open a logfile called status.logs, continuously get the updated status every minute (getStatus) and write information like the temperature and battery level (chargePercent) to the logfile. When exited, close the logfile. AI!

	status, err := getStatus()

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

		// Okay, no successful call and no error messages. Let's parse the response into an object.
		status, err := ParseGetStatusMsg(string(body))
		if err != nil {
			fmt.Println("HTTP request successful but error parsing message: ", err)
			return nil, err
		}

		// For now, we only print. This will be later used somewhere.
		fmt.Println("Battery: ", status.Result.System.MCU.ChargePercent)
		fmt.Println("Temperature: ", status.Result.System.MCU.Temperature)
		return status, nil

	} else {
		fmt.Println("getStatus Request failed with Code: ", statusCode, "and Msg: ", body)
		return nil, err
	}
}

func mcuHttpQuery(url string, commands []string) ([]byte, int, error) {

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

func parseError(body []byte) error {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if error, exists := result["error"]; exists && error.(map[string]interface{})["code"].(float64) < float64(0) {
		return fmt.Errorf(fmt.Sprintf("%v", result["error"]))
	}
	return nil
}
