package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func fakeUmami(t *testing.T) *httptest.Server {
	sites := map[string]Website{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/auth/login" {
			json.NewEncoder(w).Encode(map[string]string{"token": "jwt"})
			return
		}
		if r.Header.Get("Authorization") != "Bearer jwt" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/websites/")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/websites":
			var in WebsiteInput
			json.NewDecoder(r.Body).Decode(&in)
			sites["abc"] = Website{ID: "abc", Name: in.Name, Domain: in.Domain}
			json.NewEncoder(w).Encode(sites["abc"])
		case r.Method == http.MethodGet:
			s, ok := sites[id]
			if !ok {
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode(s)
		case r.Method == http.MethodPost:
			var in WebsiteInput
			json.NewDecoder(r.Body).Decode(&in)
			sites[id] = Website{ID: id, Name: in.Name, Domain: in.Domain}
			json.NewEncoder(w).Encode(sites[id])
		case r.Method == http.MethodDelete:
			delete(sites, id)
			w.Write([]byte(`{"ok":true}`))
		}
	}))
}

func TestWebsiteLifecycle(t *testing.T) {
	srv := fakeUmami(t)
	defer srv.Close()
	ctx := context.Background()
	c, err := New(ctx, Config{Endpoint: srv.URL, Username: "u", Password: "p"})
	if err != nil {
		t.Fatal(err)
	}
	w, err := c.CreateWebsite(ctx, WebsiteInput{Name: "n", Domain: "d.com"})
	if err != nil || w.ID != "abc" {
		t.Fatalf("create: %v %+v", err, w)
	}
	if _, err = c.UpdateWebsite(ctx, "abc", WebsiteInput{Name: "n2", Domain: "d.com"}); err != nil {
		t.Fatal(err)
	}
	if g, _ := c.GetWebsite(ctx, "abc"); g.Name != "n2" {
		t.Fatalf("update not applied: %+v", g)
	}
	if err = c.DeleteWebsite(ctx, "abc"); err != nil {
		t.Fatal(err)
	}
	var nf *NotFoundError
	if _, err = c.GetWebsite(ctx, "abc"); !errors.As(err, &nf) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
	if err = c.DeleteWebsite(ctx, "abc"); err != nil {
		t.Fatalf("delete must be idempotent: %v", err)
	}
}
