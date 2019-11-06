package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/gorilla/mux"
)

func handlePanic(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			trace := formatStackTrace(string(debug.Stack()))

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("content-type", "application/json")
			fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", r, trace)

		}
	}()

	panic("Ohh..!")
}

func formatStackTrace(trace string) string {

	slice := strings.Split(trace, "\n")

	trace = ""
	for _, v := range slice {
		if strings.Contains(v, ".go:") {
			re1 := regexp.MustCompile(`(.*?).go:`)
			str1 := []byte(re1.FindString(v))

			re2 := regexp.MustCompile(`.go:([0-9]+)`)
			str2 := []byte(strings.TrimSpace(re2.FindString(v)))

			filename := strings.TrimSpace(string(str1[:len(str1)-1]))
			line := string(str2[4:])

			fmt.Sprintf("%v,%v", filename, line)

			trace = trace + "<a href=\"/debug?filename=" + url.QueryEscape(filename) + "&line=" + line + "\">" + v + "</a>" + "\n"
		} else {
			trace = trace + v + "\n"
		}
	}

	return trace
}

func handleDebug(w http.ResponseWriter, r *http.Request) {

	qparam := r.URL.Query()
	filename, _ := url.QueryUnescape(qparam.Get("filename"))
	//line, _ := strconv.Atoi(qparam.Get("line"))

	var details string

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for n := 1; sc.Scan(); n++ {
		details = details + sc.Text() + "\n"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = quick.Highlight(w, details, "go", "html", "github")
}

func handlers() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/panic", handlePanic).Methods("GET")
	router.HandleFunc("/debug", handleDebug).Methods("GET")

	return router
}

func main() {
	router := handlers()
	log.Fatal(http.ListenAndServe(":8000", router))
}
