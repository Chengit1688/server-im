package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url string) (response []byte, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetWithBasicAuth(url string, username string, password string, timeoutSecond int) (content []byte, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		err = errors.New(fmt.Sprintf("get with basic auth error, new request error, error: %v", err))
		return nil, err
	}

	req.Close = true
	req.SetBasicAuth(username, password)

	var resp *http.Response
	client := http.Client{Timeout: time.Duration(timeoutSecond) * time.Second}
	if resp, err = client.Do(req); err != nil {
		err = errors.New(fmt.Sprintf("get with basic auth error, client do error, error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("get with basic auth error, read all error, error: %v", err))
		return nil, err
	}
	return result, nil
}

func Post(url string, data interface{}, timeoutSecond int) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, json marshal error, error: %v", err))
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, new request error, error: %v", err))
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: time.Duration(timeoutSecond) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, client do error, error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, read all error, error: %v", err))
		return nil, err
	}
	return result, nil
}

func PostWithToken(url string, data interface{}, token string, timeoutSecond int) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, json marshal error, error: %v", err))
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, new request error, error: %v", err))
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")
	req.Header.Add("token", token)

	client := &http.Client{Timeout: time.Duration(timeoutSecond) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, client do error, error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("post error, read all error, error: %v", err))
		return nil, err
	}
	return result, nil
}

func PostWithBasicAuth(url string, data interface{}, username string, password string, timeoutSecond int) ([]byte, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		err = errors.New(fmt.Sprintf("post with basic auth error, json marshal error, error: %v", err))
		return nil, err
	}

	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr)); err != nil {
		err = errors.New(fmt.Sprintf("post with basic auth error, new request error, error: %v", err))
		return nil, err
	}

	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")
	req.SetBasicAuth(username, password)

	var resp *http.Response
	client := http.Client{Timeout: time.Duration(timeoutSecond) * time.Second}
	if resp, err = client.Do(req); err != nil {
		err = errors.New(fmt.Sprintf("post with basic auth error, client do error, error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("post with basic auth error, read all error, error: %v", err))
		return nil, err
	}
	return result, nil
}

func DeleteWithBasicAuth(url string, username string, password string, timeoutSecond int) (content []byte, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodDelete, url, nil); err != nil {
		err = errors.New(fmt.Sprintf("delete with basic auth error, new request error, error: %v", err))
		return nil, err
	}

	req.Close = true
	req.SetBasicAuth(username, password)

	var resp *http.Response
	client := http.Client{Timeout: time.Duration(timeoutSecond) * time.Second}
	if resp, err = client.Do(req); err != nil {
		err = errors.New(fmt.Sprintf("delete with basic auth error, client do error, error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("delete with basic auth error, read all error, error: %v", err))
		return nil, err
	}
	return result, nil
}
