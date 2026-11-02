// Package client es un cliente minimo de la API REST de Umami (self-hosted y cloud).
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	endpoint string
	http     *http.Client
	token    string // JWT (self-hosted)
	apiKey   string // x-umami-api-key (Umami Cloud)
}

type Config struct {
	Endpoint string
	Username string
	Password string
	APIKey   string
	Insecure bool
}

type NotFoundError struct{ Path string }

func (e *NotFoundError) Error() string { return "no encontrado: " + e.Path }

// New crea el cliente; con usuario/contrasena hace login y guarda el JWT.
func New(ctx context.Context, cfg Config) (*Client, error) {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	if cfg.Insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // opt-in explicito
	}
	c := &Client{
		endpoint: strings.TrimRight(cfg.Endpoint, "/"),
		http:     &http.Client{Timeout: 30 * time.Second, Transport: tr},
		apiKey:   cfg.APIKey,
	}
	if c.apiKey != "" {
		return c, nil
	}
	var out struct {
		Token string `json:"token"`
	}
	body := map[string]string{"username": cfg.Username, "password": cfg.Password}
	if err := c.do(ctx, http.MethodPost, "/api/auth/login", body, &out); err != nil {
		return nil, fmt.Errorf("login en Umami: %w", err)
	}
	if out.Token == "" {
		return nil, fmt.Errorf("login en Umami: respuesta sin token")
	}
	c.token = out.Token
	return c, nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var rdr io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.apiKey != "" {
		req.Header.Set("x-umami-api-key", c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Path: path}
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: %s: %s", method, path, resp.Status, strings.TrimSpace(string(data)))
	}
	if out != nil && len(data) > 0 {
		return json.Unmarshal(data, out)
	}
	return nil
}
