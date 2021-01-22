package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := getListenPort()
	fmt.Printf("Starting HTTP server on %d", port)

	http.HandleFunc("/", stanRequestHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func getListenPort() int {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		return 8080
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	return port
}

// helper test functions for filtering the request payloads
// for the required criteria of the response
func hasDRM(p payload) bool {
	return p.DRM
}

func hasEpisodes(p payload) bool {
	return p.EpisodeCount >= 1
}

// fairly large request handler, should be broken up if it were any
// larger (in terms of logical steps)
func stanRequestHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Stan's custom error payload
	stanRequest, err := parseRequest(body)
	if err != nil {
		log.Printf("Error parsing JSON: %s", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(errorResponse))
		return
	}

	// filter the request down to just those that meet the criteria test functions
	filtered := filterRequest(stanRequest, []func(payload) bool{hasDRM, hasEpisodes})

	// transform the request payloads into the response format
	respW := transformPayload(filtered)

	// write and return
	jsonResponse, _ := json.Marshal(respW)
	rw.Write(jsonResponse)
}

func parseRequest(document []byte) (request, error) {
	var req request
	err := json.Unmarshal(document, &req)
	if err != nil {
		return req, err
	}

	// Upon reading the spec again, I don't think an empty payload
	// is an error, but I'll leave this in for now.
	if len(req.Payloads) == 0 {
		return req, fmt.Errorf("request contained no payloads")
	}

	return req, err
}

func filterRequest(req request, tests []func(payload) bool) (newReq request) {
	// I'm not bothering to add all the other response fields, sorry.
	// I'm sure there's a better functional way of doing this

	for _, p := range req.Payloads {
		passed := true

		for _, test := range tests {
			if !test(p) {
				passed = false
			}
		}

		if passed {
			newReq.Payloads = append(newReq.Payloads, p)
		}
	}

	return
}

// Currently a static selection of fields to use in the output.
// We assume that all entries in the request have the required fields.
func transformPayload(req request) (respW responseWrapper) {
	for _, payload := range req.Payloads {
		newResponse := response{
			Image: payload.Image.ShowImage,
			Slug:  payload.Slug,
			Title: payload.Title,
		}

		respW.Response = append(respW.Response, newResponse)
	}

	return
}
