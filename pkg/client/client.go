package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rune/internal/auth"
)

type Client struct {
	BaseURL string
	Token   string
}

func New(baseURL string, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
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
	req.Header.Set("Authorization", "Bearer "+c.Token)
	_, err := http.DefaultClient.Do(req)
	return err
}

func (c *Client) Put(key, value string) error {

	body := map[string]string{
		"key":   key,
		"value": value,
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.BaseURL+"/secret/put", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

	req, err := http.NewRequest("GET", c.BaseURL+"/secret/get?key="+key, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := http.DefaultClient.Do(req)
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

func (c *Client) List() ([]string, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/secret/list", nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(bodyBytes))
	}

	var keys []string
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
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

func (c *Client) CreateToken(name string) (string, auth.TokenRecord, error) {
	body := map[string]string{
		"name": name,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", c.BaseURL+"/token/create", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", auth.TokenRecord{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", auth.TokenRecord{}, fmt.Errorf(string(b))
	}
	var result struct {
		Token  string           `json:"token"`
		Record auth.TokenRecord `json:"record"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Token, result.Record, nil
}

func (c *Client) ListTokens() ([]auth.TokenRecord, error) {

	req, _ := http.NewRequest("GET", c.BaseURL+"/token/list", nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(b))
	}
	var tokens []auth.TokenRecord
	json.NewDecoder(resp.Body).Decode(&tokens)
	return tokens, nil

}

func (c *Client) RevokeToken(id string) error {

	body := map[string]string{
		"id": id,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", c.BaseURL+"/token/revoke", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(string(b))
	}
	return nil

}
