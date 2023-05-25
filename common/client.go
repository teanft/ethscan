package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/teanft/ethscan/config"
	"github.com/teanft/ethscan/util"
	"io"
	"net/http"
)

type RequestData struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int64    `json:"id"`
}

func sendHTTPRequest(url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(util.NewErr("failed to close body", err))
			return
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("Request failed: %s", respBody))
	}

	return respBody, nil
}

func getJSONRPCRequestBody(jsonrpc, method string, params []string, id int64) ([]byte, error) {
	requestData := RequestData{
		Jsonrpc: jsonrpc,
		Method:  method,
		Params:  params,
		Id:      id,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	return requestBody, nil
}

func CallRPC(c *gin.Context, method string, params []string) (map[string]interface{}, error) {

	requestBody, err := getJSONRPCRequestBody("2.0", method, params, 83)
	if err != nil {
		Fail(c, nil, fmt.Sprintf("Failed to construct JSON-RPC request: %s", err.Error()))
		return nil, util.NewErr("failed to construct JSON-RPC request", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	responseBytes, err := sendHTTPRequest(config.Cfg.Client.URL, headers, bytes.NewBuffer(requestBody))
	if err != nil {
		Fail(c, nil, fmt.Sprintf("Failed to send HTTP request: %s", err.Error()))
		return nil, util.NewErr("failed to send HTTP request", err)
	}

	var responseData map[string]interface{}

	if err = json.Unmarshal(responseBytes, &responseData); err != nil {
		Fail(c, nil, fmt.Sprintf("Failed to decode response body: %s", err.Error()))
		return nil, util.NewErr("failed to decode response body", err)
	}
	return responseData, nil
}
