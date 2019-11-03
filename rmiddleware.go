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

			trace := string(debug.Stack())

			trace = formatStackTrace(trace)

			fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", r, trace)

			// //details := fmt.Sprintf("path : %q line : %q", str1[:len(str1)-1], str2[4:])

		}
	}()

	panic("Ohh..!")
}

func formatStackTrace(trace string) string {

	slice := strings.Split(trace, "\n")

	trace = ""
	for _, v := range slice {
		if strings.Contains(v, ".go:") {
			re := regexp.MustCompile(`(.*?).go:`)
			str1 := []byte(re.FindString(v))

			re = regexp.MustCompile(`.go:([0-9]+) `)
			str2 := []byte(strings.TrimSpace(re.FindString(v)))

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
		log.Fatal(err)
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

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/panic", handlePanic).Methods("GET")
	router.HandleFunc("/debug", handleDebug).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
