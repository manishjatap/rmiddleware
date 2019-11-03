package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicApiSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Fatal(err)
	}

	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePanic)
	handler.ServeHTTP(rw, req)
	assert.Equal(t, rw.Code, http.StatusOK, "API should return 200 OK response")

}

func TestMain(m *testing.M) {
	//fmt.Println(main())
	os.Exit(m.Run())
}
