package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)
	http.HandleFunc("/gen", generateHandler)

	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main.html")
}

type Response struct {
type RespSuccess struct {
	Created int64
	Data    []struct {
		Url string
	}
}

	Error struct {
		Code    string
		Message string
		Param   string
		Type    string
	}
}

type DataResponse struct {
	StatusCode int
	Data       string
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle parameter
	content := strings.TrimSpace(r.FormValue("content"))

	url := "https://api.openai.com/v1/images/generations"
	method := "POST"
	openaiKey := os.Getenv("OPENAI_KEY") //Your API key
	param := fmt.Sprintf(`{
		"prompt": "%s",
		"n": 1,
		"size": "256x256",
		"response_format": "url"
	  }`, content)

	payload := strings.NewReader(param)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+openaiKey)
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
	fmt.Println(string(body))

	if res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusUnauthorized {
		var respFail RespFail
		json.Unmarshal([]byte(body), &respFail)
		response, err := json.Marshal(DataResponse{1, respFail.Error.Message})
		if err != nil {
			panic(err)
		}
		// Send response
		fmt.Fprintf(w, "%s", response)
		return
	}

	var respSuccess RespSuccess
	json.Unmarshal([]byte(body), &respSuccess)

	response, err := json.Marshal(DataResponse{0, respSuccess.Data[0].Url})
	if err != nil {
		panic(err)
	}
	// Send response
	fmt.Fprintf(w, "%s", response)
	return
}
