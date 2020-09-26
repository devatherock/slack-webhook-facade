package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Attachment struct {
	Text  string
	Title string `json:",omitempty"`
}

type SlackRequest struct {
	Text        string       `json:",omitempty"`
	Channel     string       `json:",omitempty"`
	Attachments []Attachment `json:",omitempty"`
}

func init() {
	level, ok := os.LookupEnv("PARAMETER_LOG_LEVEL")

	// LOG_LEVEL not set, default to info
	if !ok {
		level = "info"
	}

	// parse string, this is built-in feature of logrus
	logLevel, error := log.ParseLevel(level)
	if error != nil {
		logLevel = log.InfoLevel
	}

	// set global log level
	log.SetLevel(logLevel)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/zulip/{authorization}", ZulipHandler)
	http.Handle("/", router)

	http.HandleFunc("/api/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("UP"))
	})

	// Read from PORT environment variable available on heroku
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	log.Println("Http server listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

// Handles requests to Zulip
func ZulipHandler(writer http.ResponseWriter, request *http.Request) {
	// Read request
	requestBody, error := ioutil.ReadAll(request.Body)
	if error != nil {
		writeErrorResponse(writer, error, 400)
		return
	}

	// Parse request
	slackRequest := SlackRequest{}
	error = json.Unmarshal(requestBody, &slackRequest)
	if error != nil {
		writeErrorResponse(writer, error, 400)
		return
	}

	zulipPayload := url.Values{}
	zulipPayload.Set("type", "stream")
	zulipPayload.Set("to", slackRequest.Channel)
	zulipPayload.Set("topic", slackRequest.Attachments[0].Title)

	text := slackRequest.Attachments[0].Text
	if text == "" {
		text = slackRequest.Text
	}
	zulipPayload.Set("content", text)

	pathVariables := mux.Vars(request)
	zulipRequest, _ := http.NewRequest("POST", request.URL.Query().Get("server")+"/api/v1/messages",
		strings.NewReader(zulipPayload.Encode())) // URL-encoded payload
	zulipRequest.Header.Add("Authorization", "Basic "+pathVariables["authorization"])
	zulipRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	zulipResponse, error := client.Do(zulipRequest)
	if error != nil {
		writeErrorResponse(writer, error, 500)
		return
	}
	defer zulipResponse.Body.Close()
	if zulipResponse.StatusCode < 400 {
		log.Println("Message posted to Zulip with http status", zulipResponse.StatusCode)
	} else {
		log.Println("Posting message to Zulip failed with http status", zulipResponse.StatusCode)
		zulipResponseBody, error := ioutil.ReadAll(zulipResponse.Body)
		if error != nil {
			writeErrorResponse(writer, error, zulipResponse.StatusCode)
		} else {
			log.Error("Zulip Response: ", string(zulipResponseBody))
			writer.WriteHeader(zulipResponse.StatusCode)
			writeResponse(writer, false)
		}
		return
	}

	writeResponse(writer, true)
}

// Writes an error HTTP status code
func writeErrorResponse(writer http.ResponseWriter, err error, status int) {
	log.Error("error: ", err)
	writer.WriteHeader(status)
	writeResponse(writer, false)
}

// Writes the response body
func writeResponse(writer http.ResponseWriter, status bool) {
	writer.Header().Set("Content-Type", "application/json")
	converterResponse := map[string]interface{}{
		"success": status,
	}

	responseBody, error := json.Marshal(&converterResponse)
	if error != nil {
		log.Error("error: ", error)
		return
	} else {
		writer.Write(responseBody)
	}
}
