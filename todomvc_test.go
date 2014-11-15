package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	setupDb("/tmp/TodoDabaseTest", true)
}

func PerformRequestJson(r http.Handler, method string, url string, p interface{}) (*httptest.ResponseRecorder, error) {
	var err error
	var stream []byte
	var req *http.Request
	if p != nil {
		stream, err = json.Marshal(p)
		byteBuffer := bytes.NewBuffer(stream)
		req, err = http.NewRequest(method, url, byteBuffer)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp, err
}

func TestAdd(t *testing.T) {
	r := gin.New()
	configApi(r)
	w, _ := PerformRequestJson(r, "POST", "/api/todos", Todo{Title: "Task1", Completed: false})
	assert.Equal(t, http.StatusOK, w.Code)
	w, _ = PerformRequestJson(r, "POST", "/api/todos", Todo{Title: "Task2", Completed: true})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCount(t *testing.T) {
	r := gin.New()
	configApi(r)
	w, _ := PerformRequestJson(r, "GET", "/api/todos", Todo{})
	todos := []Todo{}
	json.Unmarshal(w.Body.Bytes(), &todos)
	assert.Equal(t, 2, len(todos))
}
