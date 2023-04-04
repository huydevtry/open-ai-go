package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/gen", generateHandler)

	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main.html")
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle parameter
	content := r.FormValue("content")

	url := "https://api.openai.com/v1/images/generations"
	method := "POST"

	payload := strings.NewReader(`{
	"prompt": "` + content + `",
	"n": 1,
	"size": "1024x1024",
	"response_format": "url"
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send response
	fmt.Fprintf(w, "%s", body)
}
