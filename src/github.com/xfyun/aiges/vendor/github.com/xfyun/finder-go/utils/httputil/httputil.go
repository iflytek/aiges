package httputil

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func DoGet(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("status_code:", resp.StatusCode)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
	}

	return result, nil
}

func DoPost(client *http.Client, contentType string, url string, data []byte) ([]byte, error) {
	resp, err := client.Post(url, contentType, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
