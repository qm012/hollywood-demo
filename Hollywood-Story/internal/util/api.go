package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/qm012/dun"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"time"
)

type APIInternalPromptResp struct {
	Content string       `json:"content"`
	Usage   openai.Usage `json:"usage"`
}

func APIInternalPrompt(ID string, variables map[string]string) (*APIInternalPromptResp, error) {
	if len(ID) == 0 {
		return nil, errors.New("ID cant be empty")
	}
	if len(variables) == 0 {
		variables = map[string]string{}
	}
	params := struct {
		Variables map[string]string `json:"variables"`
	}{
		Variables: variables,
	}
	url := "http://127.0.0.1:8092/internal/prompts/" + ID

	paramsDate, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Post(url, "application/json", bytes.NewReader(paramsDate))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var responseBody = &dun.Response{}
	if err = json.Unmarshal(body, responseBody); err != nil {
		return nil, err
	}
	if responseBody.Code != http.StatusOK {
		return nil, errors.New(responseBody.Error())
	}
	promptRespByte, err := json.Marshal(responseBody.Data)
	if err != nil {
		return nil, err
	}

	var promptResp = new(APIInternalPromptResp)
	if err = json.Unmarshal(promptRespByte, promptResp); err != nil {
		return nil, err
	}
	return promptResp, nil
}
