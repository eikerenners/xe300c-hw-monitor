package main

import (
	"encoding/json"
	"fmt"
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
