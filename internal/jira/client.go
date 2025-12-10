package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client é o cliente para interagir com a API do Jira
type Client struct {
	config     *JiraConfig
	httpClient *http.Client
	baseURL    string
}

// NewClient cria um novo cliente Jira
func NewClient(config *JiraConfig) (*Client, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("URL do Jira não configurada")
	}
	if config.Email == "" {
		return nil, fmt.Errorf("email não configurado")
	}
	if config.APIToken == "" {
		return nil, fmt.Errorf("API token não configurado")
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: config.URL,
	}, nil
}

// createAuthHeader cria o header de autenticação
func (c *Client) createAuthHeader() string {
	auth := fmt.Sprintf("%s:%s", c.config.Email, c.config.APIToken)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// doRequest executa uma requisição HTTP
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/api/3/%s", c.baseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("erro ao serializar body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("erro na API do Jira (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Issue representa uma issue do Jira
type Issue struct {
	Key    string `json:"key"`
	ID     string `json:"id"`
	Self   string `json:"self"`
	Fields struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		IssueType   struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"issuetype"`
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
		EpicLink string `json:"customfield_10014,omitempty"` // Epic Link field
		Status   struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

// CreateIssueRequest representa a requisição para criar uma issue
type CreateIssueRequest struct {
	Fields struct {
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
		Summary     string `json:"summary"`
		Description struct {
			Type    string `json:"type"`
			Version int    `json:"version"`
			Content []struct {
				Type    string `json:"type"`
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"content"`
		} `json:"description"`
		IssueType struct {
			ID string `json:"id"`
		} `json:"issuetype"`
		EpicLink string `json:"customfield_10014,omitempty"` // Epic Link
	} `json:"fields"`
}

// CreateEpic cria um Epic no Jira
func (c *Client) CreateEpic(summary, description string) (*Issue, error) {
	// Primeiro, precisamos obter o ID do tipo "Epic"
	issueTypes, err := c.GetIssueTypes()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter tipos de issue: %w", err)
	}

	epicTypeID := ""
	for _, it := range issueTypes {
		if it.Name == "Epic" {
			epicTypeID = it.ID
			break
		}
	}

	if epicTypeID == "" {
		return nil, fmt.Errorf("tipo 'Epic' não encontrado no projeto")
	}

	req := CreateIssueRequest{}
	req.Fields.Project.Key = c.config.Project
	req.Fields.Summary = summary
	req.Fields.Description.Type = "doc"
	req.Fields.Description.Version = 1
	req.Fields.Description.Content = []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}{
		{
			Type: "paragraph",
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{
				{
					Type: "text",
					Text: description,
				},
			},
		},
	}
	req.Fields.IssueType.ID = epicTypeID

	respBody, err := c.doRequest("POST", "issue", req)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return &issue, nil
}

// CreateIssue cria uma issue (card) no Jira
func (c *Client) CreateIssue(summary, description, issueTypeName, epicKey string) (*Issue, error) {
	// Obter tipos de issue
	issueTypes, err := c.GetIssueTypes()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter tipos de issue: %w", err)
	}

	issueTypeID := ""
	for _, it := range issueTypes {
		if it.Name == issueTypeName {
			issueTypeID = it.ID
			break
		}
	}

	if issueTypeID == "" {
		// Tentar usar "Task" como padrão
		for _, it := range issueTypes {
			if it.Name == "Task" {
				issueTypeID = it.ID
				break
			}
		}
		if issueTypeID == "" {
			return nil, fmt.Errorf("tipo de issue '%s' não encontrado", issueTypeName)
		}
	}

	req := CreateIssueRequest{}
	req.Fields.Project.Key = c.config.Project
	req.Fields.Summary = summary
	req.Fields.Description.Type = "doc"
	req.Fields.Description.Version = 1
	req.Fields.Description.Content = []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}{
		{
			Type: "paragraph",
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{
				{
					Type: "text",
					Text: description,
				},
			},
		},
	}
	req.Fields.IssueType.ID = issueTypeID

	if epicKey != "" {
		req.Fields.EpicLink = epicKey
	}

	respBody, err := c.doRequest("POST", "issue", req)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return &issue, nil
}

// IssueType representa um tipo de issue
type IssueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
}

// GetIssueTypes obtém os tipos de issue disponíveis no projeto
func (c *Client) GetIssueTypes() ([]IssueType, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("project/%s", c.config.Project), nil)
	if err != nil {
		return nil, err
	}

	var project struct {
		IssueTypes []IssueType `json:"issueTypes"`
	}
	if err := json.Unmarshal(respBody, &project); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return project.IssueTypes, nil
}

// GetEpic obtém um Epic pelo key
func (c *Client) GetEpic(epicKey string) (*Issue, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("issue/%s", epicKey), nil)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return &issue, nil
}

