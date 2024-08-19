package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	stomp "github.com/go-stomp/stomp/v3"
)

func NewHttpSender(name, url, body, requestType string, headers, params map[string]string, timeout int) (*http.Response, error) {
	_timeout := 5
	if timeout > 0 && timeout < 1000 {
		_timeout = timeout
	}

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    time.Duration(_timeout),
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	if requestType != "GET" && requestType != "POST" && requestType != "PUT" && requestType != "DELETE" {
		return nil, fmt.Errorf("request type %s not supported", requestType)
	}
	_body := []byte(body)
	_params := "?"
	for key, value := range params {
		_params += key + "=" + value + "&"
	}
	if len(_params) > 1 {
		url = url + _params
	}
	req, _ := http.NewRequest(requestType, url, bytes.NewBuffer(_body))
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err == nil {
		return resp, nil
	}
	return nil, err
}

func NewStompSender(destination, server, username, password, data string) error {
	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(username, password),
		stomp.ConnOpt.Host("/"),
	}
	conn, err := stomp.Dial("tcp", server, options...)
	if err != nil {
		println("cannot connect to server", err.Error())
		return err
	}
	err = conn.Send(destination, "text/plain", []byte(data), nil)
	if err != nil {
		println("failed to send to server", err)
		return err
	}
	return nil
}
