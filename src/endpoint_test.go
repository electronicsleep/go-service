package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: %v : %v", status, http.StatusOK)
	}

	expected := "{\"status\":\"up\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: %v: %v", rr.Body.String(), expected)
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statusHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: %v : %v", status, http.StatusOK)
	}

	expected := "go-service: ok"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: %v: %v", rr.Body.String(), expected)
	}
}
