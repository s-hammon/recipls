package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerReadiness(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			"GET 200",
			"GET",
			http.StatusOK,
			`{"status":"ok"}`,
		},
		{
			"POST 405",
			"POST",
			http.StatusMethodNotAllowed,
			`{"status":"error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/v1/healthz", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlerReadiness)

			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
