package trace

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleTraceOn(t *testing.T) {
	requestURL := fmt.Sprintf("http://199.199.199.199:9999/trace/on")
	var reader io.Reader
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, requestURL, reader)

	HandleTraceOn(w, r)
	body := make([]byte, 50)
	defer w.Result().Body.Close()
	count, _ := w.Result().Body.Read(body)
	if count != 17 {
		t.Fatalf("incorrect body length: %v", count)
	}
	if string(body[:17]) != "Tracing activated" {
		t.Errorf("incorrect response received: %v", string(body))
	}
}

func TestHandleTraceOff(t *testing.T) {
	requestURL := fmt.Sprintf("http://199.199.199.199:9999/trace/off")
	var reader io.Reader
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, requestURL, reader)

	HandleTraceOff(w, r)
	body := make([]byte, 50)
	defer w.Result().Body.Close()
	count, _ := w.Result().Body.Read(body)
	if count != 19 {
		t.Fatalf("incorrect body length: %v", count)
	}
	if string(body[:19]) != "Tracing deactivated" {
		t.Errorf("incorrect response received: %v", string(body))
	}
}
