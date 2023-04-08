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

	fs := http.FileServer(http.Dir("view/static"))
	http.Handle("/view/static/", http.StripPrefix("/view/static/", fs))

	http.HandleFunc("/chat", handlerChat)
	http.HandleFunc("/", handler)

	http.HandleFunc("/gen", generateHandler)
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(":8080", nil)
}

type RespSuccess struct {
	Created int64
	Data    []struct {
		Url string
	}
}

type RespFail struct {
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

//Main page
func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "view/main.html")
}

//Generate image API
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
}

//Handle websocket
func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	ServeWs(wsServer, w, r)
}

func handlerChat(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "view/chat.html")

}
