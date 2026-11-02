package client

import (
	"context"
	"net/http"
)

type Website struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Domain  string  `json:"domain"`
	ShareID *string `json:"shareId"`
}

type WebsiteInput struct {
	Name    string  `json:"name"`
	Domain  string  `json:"domain"`
	ShareID *string `json:"shareId"`
}

func (c *Client) CreateWebsite(ctx context.Context, in WebsiteInput) (*Website, error) {
	var w Website
	if err := c.do(ctx, http.MethodPost, "/api/websites", in, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *Client) GetWebsite(ctx context.Context, id string) (*Website, error) {
	var w Website
	if err := c.do(ctx, http.MethodGet, "/api/websites/"+id, nil, &w); err != nil {
		return nil, err
	}
	// Umami devuelve 200 con cuerpo vacio/id vacio para sitios borrados (deletedAt)
	if w.ID == "" {
		return nil, &NotFoundError{Path: "/api/websites/" + id}
	}
	return &w, nil
}

func (c *Client) UpdateWebsite(ctx context.Context, id string, in WebsiteInput) (*Website, error) {
	var w Website
	if err := c.do(ctx, http.MethodPost, "/api/websites/"+id, in, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *Client) DeleteWebsite(ctx context.Context, id string) error {
	err := c.do(ctx, http.MethodDelete, "/api/websites/"+id, nil, nil)
	if _, nf := err.(*NotFoundError); nf {
		return nil
	}
	return err
}
