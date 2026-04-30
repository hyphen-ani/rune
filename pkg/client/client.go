package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) Unseal(passphrase string) error {
	body := map[string]string{
		"passphrase": passphrase,
	}

	data, _ := json.Marshal(body)
	resp, err := http.Post(c.BaseURL+"/unseal", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: %s", string(bodyBytes))
	}

	return nil

}

func (c *Client) Seal() error {
	req, _ := http.NewRequest("POST", c.BaseURL+"/seal", nil)
	_, err := http.DefaultClient.Do(req)
	return err
}

func (c *Client) Put(key, value string) error {

	body := map[string]string{
		"key":   key,
		"value": value,
	}

	data, _ := json.Marshal(body)
	resp, err := http.Post(c.BaseURL+"/secret/put", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(string(bodyBytes))
	}
	return nil

}

func (c *Client) Get(key string) (string, error) {

	resp, err := http.Get(c.BaseURL + "/secret/get?key=" + key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(string(body))
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["value"], nil

}

func (c *Client) Status() (string, error) {

	resp, err := http.Get(c.BaseURL + "/status")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(string(body))
	}
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil

}
