package controllers

import (
	"bytes"
	"fmt"
	"golang-proxy-server/common/model"
	"golang-proxy-server/config"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// TCPExternalAPI - TCP External API config
// ClientID: client id
// IP: external ip, if ip is 'mock', return mock response
// Method: support POST/GET(default)
// Data: When method is GET, Data is a query string(e.q. 'key=value'), When method is POST, Data is a json(e.q. '{"key1":"value1", "key2":"value2"}')
type TCPExternalAPI struct {
	ClientID string
	IP       string
	Method   string
	Data     interface{}
}

// Request - a TCPExternalAPI receiver for send request to ExternalRequestQueue channel
func (t *TCPExternalAPI) Request() (str string, err error) {
	config.Info.Printf("Client %v request to '%v' data: %v", t.ClientID, t.IP, t.Data)
	if len(model.ExternalRequestQueue) < config.DeploySet.External.ExternalRequestQueue {
		model.ExternalRequestQueue <- t
		if stringData, ok := t.Data.(string); ok {
			return "MockExternal request data: (" + stringData + ") ing...", nil
		}
		return "MockExternal requesting...", nil
	}
	config.Warning.Printf("External api request queue(%v) is full, please wait", config.DeploySet.External.ExternalRequestQueue)
	return "", fmt.Errorf("External api request queue(%v) is full, please wait", config.DeploySet.External.ExternalRequestQueue)
}

// Worker - send http request from ExternalRequestQueue channel
func Worker(conn *net.TCPConn, clientID, clientIP string, request *TCPExternalAPI) {
	defer model.ExternalAPIRequestFinish(nil)
	defer func() {
		// recover panic
		if r := recover(); r != nil {
			config.Error.Printf("External api panic: %v", r)
			return
		}
	}()

	// if request.IP == "mock" use mock response
	if request.IP == "mock" {
		conn.Write([]byte("Mock External response data:" + request.Data.(string) + "\n"))
		return
	}

	queryString := ""
	var requestbody io.Reader
	if request.Method == "POST" {
		body, ok := request.Data.(string)
		if ok {
			requestbody = bytes.NewBuffer([]byte(body))
		} else {
			requestbody = nil
		}
		fmt.Println("requestbody:", requestbody)
	} else {
		// default use GET
		request.Method = "GET"
		var ok bool
		queryString, ok = request.Data.(string)
		if ok {
			queryString = "?" + queryString
		}
		requestbody = nil
	}

	client := &http.Client{
		Timeout: time.Duration(config.DeploySet.External.ExternalRequestTimeout) * time.Second,
	}
	req, err := http.NewRequest(request.Method, request.IP+queryString, requestbody)
	if request.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		config.Error.Printf("External request api (%v) error: %v", request.IP+queryString, err)
		conn.Write([]byte("External request api (" + request.IP + queryString + ") error: " + err.Error() + "\n"))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		config.Error.Printf("External request api (%v) error: %v", request.IP+queryString, err)
		conn.Write([]byte("External request api (" + request.IP + queryString + ") error: " + err.Error() + "\n"))
		return
	}

	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		config.Error.Printf("External api (%v) parse response body error: %v", request.IP+queryString, err)
		conn.Write([]byte("External api (" + request.IP + queryString + ") parse response body error: " + err.Error() + "\n"))
	}

	conn.Write([]byte("External response data:" + string(result) + "\n"))
	return
}
