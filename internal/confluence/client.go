package confluence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client é o cliente para interagir com a API do Confluence
type Client struct {
	config     *ConfluenceConfig
	httpClient *http.Client
	baseURL    string
}

// NewClient cria um novo cliente Confluence
func NewClient(config *ConfluenceConfig) (*Client, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("URL do Confluence não configurada")
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
	url := fmt.Sprintf("%s/wiki/rest/api/%s", c.baseURL, endpoint)

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
		return nil, fmt.Errorf("erro na API do Confluence (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Page representa uma página do Confluence
type Page struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Space struct {
		Key string `json:"key"`
	} `json:"space"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
}

// CreatePageRequest representa a requisição para criar uma página
type CreatePageRequest struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Space struct {
		Key string `json:"key"`
	} `json:"space"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
}

// CreatePage cria uma página no Confluence
func (c *Client) CreatePage(title, content string, parentID string) (*Page, error) {
	req := CreatePageRequest{
		Type:  "page",
		Title: title,
	}
	req.Space.Key = c.config.Space
	req.Body.Storage.Value = content
	req.Body.Storage.Representation = "storage"

	// Se houver parent, adicionar
	if parentID != "" {
		reqBody := map[string]interface{}{
			"type":  "page",
			"title": title,
			"space": map[string]string{
				"key": c.config.Space,
			},
			"body": map[string]interface{}{
				"storage": map[string]interface{}{
					"value":          content,
					"representation": "storage",
				},
			},
			"ancestors": []map[string]string{
				{"id": parentID},
			},
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("erro ao serializar body: %w", err)
		}

		url := fmt.Sprintf("%s/wiki/rest/api/content", c.baseURL)
		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("erro ao criar requisição: %w", err)
		}

		httpReq.Header.Set("Authorization", c.createAuthHeader())
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("erro ao executar requisição: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler resposta: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("erro na API do Confluence (status %d): %s", resp.StatusCode, string(respBody))
		}

		var page Page
		if err := json.Unmarshal(respBody, &page); err != nil {
			return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
		}

		return &page, nil
	}

	respBody, err := c.doRequest("POST", "content", req)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(respBody, &page); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return &page, nil
}

// UpdatePage atualiza uma página existente
func (c *Client) UpdatePage(pageID, title, content string) (*Page, error) {
	// Primeiro, obter a página atual para pegar a versão
	getResp, err := c.doRequest("GET", fmt.Sprintf("content/%s?expand=body.storage,version", pageID), nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter página: %w", err)
	}

	var currentPage Page
	if err := json.Unmarshal(getResp, &currentPage); err != nil {
		return nil, fmt.Errorf("erro ao parsear página atual: %w", err)
	}

	// Criar requisição de atualização
	reqBody := map[string]interface{}{
		"id":    pageID,
		"type":  "page",
		"title": title,
		"body": map[string]interface{}{
			"storage": map[string]interface{}{
				"value":          content,
				"representation": "storage",
			},
		},
		"version": map[string]interface{}{
			"number": currentPage.Version.Number + 1,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar body: %w", err)
	}

	url := fmt.Sprintf("%s/wiki/rest/api/content/%s", c.baseURL, pageID)
	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	httpReq.Header.Set("Authorization", c.createAuthHeader())
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("erro na API do Confluence (status %d): %s", resp.StatusCode, string(respBody))
	}

	var page Page
	if err := json.Unmarshal(respBody, &page); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return &page, nil
}

// SearchPages busca páginas no espaço
func (c *Client) SearchPages(query string) ([]Page, error) {
	url := fmt.Sprintf("%s/wiki/rest/api/content/search?cql=space=%s+and+text~\"%s\"", c.baseURL, c.config.Space, query)
	
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	httpReq.Header.Set("Authorization", c.createAuthHeader())
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("erro na API do Confluence (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Results []Page `json:"results"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return result.Results, nil
}

