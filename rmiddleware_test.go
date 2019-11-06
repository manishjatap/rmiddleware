package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestPanicApiSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePanic)
	handler.ServeHTTP(res, req)
	assert.Equal(t, res.Code, http.StatusInternalServerError, "API should return [500 Internal Server Error]")
}

func TestDebugApiSuccess(t *testing.T) {
	filename := createTempFile()

	req, err := http.NewRequest("GET", "/debug?filename="+url.QueryEscape(filename)+"&line=1", nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDebug)
	handler.ServeHTTP(res, req)
	assert.Equal(t, res.Code, http.StatusOK, "API should return [200 Ok]")

	deleteTempFile(filename)

}

func TestDebugApiOpenFileError(t *testing.T) {

	req, err := http.NewRequest("GET", "/debug?filename="+url.QueryEscape("/fake/path/temp.txt")+"&line=1", nil)

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDebug)
	handler.ServeHTTP(res, req)
	assert.Equal(t, res.Code, http.StatusInternalServerError, "Expected : Error")
}

func TestHandler(t *testing.T) {
	router := handlers()
	assert.NotNil(t, router, "Expected : No nil handler")
}

func createTempFile() string {
	if path, err := os.Getwd(); err == nil {
		path = path + "/deleteMe.txt"
		if file, err := os.Create(path); err == nil {
			defer file.Close()
			if _, err := file.WriteString("Line1\nLine2"); err == nil {
				return path
			}
			panic(err)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	return ""
}

func deleteTempFile(filename string) {
	if err := os.Remove(filename); err != nil {
		log.Fatal(err)
	}
}
