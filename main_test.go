package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_ParseJson(t *testing.T) {
	document, err := ioutil.ReadFile("test_fixtures/request_payload.json")
	if err != nil {
		t.Fatal("could not read test fixture")
	}

	_, err = parseRequest(document)
	if err != nil {
		t.Fatal("expected valid document to parse but it did not")
	}

	document, err = ioutil.ReadFile("test_fixtures/empty_request.json")
	if err != nil {
		t.Fatal("could not read test fixture")
	}
	_, err = parseRequest(document)
	if err == nil {
		t.Fatal("expected empty document to fail but it did not")
	}
}

func Test_FilterRequest(t *testing.T) {
	req := request{
		Payloads: []payload{
			{
				Title: "foo",
				DRM:   true,
			},
			{
				Title: "bar",
				DRM:   false,
			},
		},
	}

	oneMatch := filterRequest(req, []func(payload) bool{
		func(p payload) bool {
			return p.DRM
		},
	})

	if want, got := 1, len(oneMatch.Payloads); want != got {
		t.Fatalf("expected %d payloads but got %d", want, got)
	}

	noMatches := filterRequest(req, []func(payload) bool{
		func(p payload) bool {
			return p.Title == "a non-existent title"
		},
	})

	if want, got := 0, len(noMatches.Payloads); want != got {
		t.Fatalf("expected %d payloads but got %d", want, got)
	}
}

func deepCompareResponses(expected, got responseWrapper) error {
	if len(expected.Response) != len(got.Response) {
		return fmt.Errorf("expected length %d, got length %d", len(expected.Response), len(got.Response))
	}

	for i := 0; i < len(expected.Response); i++ {
		if expected.Response[i] != got.Response[i] {
			return fmt.Errorf("expected: %s\ngot: %s", expected.Response[i], got.Response[i])
		}
	}

	return nil
}

func Test_TransformPayload(t *testing.T) {
	req := request{
		Payloads: []payload{
			{
				Title: "a tale of foo",
				Slug:  "foo",
				Image: struct {
					ShowImage string
				}{
					ShowImage: "http://image.url",
				},
			},
		},
	}

	expectedResponseWrapper := responseWrapper{
		Response: []response{
			{
				Image: "http://image.url",
				Slug:  "foo",
				Title: "a tale of foo",
			},
		},
	}

	respW := transformPayload(req)

	if err := deepCompareResponses(expectedResponseWrapper, respW); err != nil {
		t.Fatalf("responses did not match: %s", err.Error())
	}

	badResponseWrapper := responseWrapper{}

	respW = transformPayload(req)

	if err := deepCompareResponses(badResponseWrapper, respW); err == nil {
		t.Fatalf("responses matched but were not the same")
	}
}

func Test_EndToEnd(t *testing.T) {
	// grab the corresponding response document for the request fixture
	responseFixtureJSON, err := ioutil.ReadFile("test_fixtures/response_payload.json")
	if err != nil {
		t.Fatal("could not read test fixture")
	}
	var responseFixture responseWrapper
	json.Unmarshal(responseFixtureJSON, &responseFixture)

	// grab the request fixture and parse it
	requestFixture, err := ioutil.ReadFile("test_fixtures/request_payload.json")
	if err != nil {
		t.Fatal("could not read test fixture")
	}
	req, _ := parseRequest(requestFixture)

	// filter for the code challenge's required test predicates and transform
	// to the required output format
	filtered := filterRequest(req, []func(payload) bool{hasDRM, hasEpisodes})
	response := transformPayload(filtered)

	// the generated response and the test fixture should match exactly
	if err := deepCompareResponses(responseFixture, response); err != nil {
		t.Fatalf("generated response did not match response test fixture:\n %s", err.Error())
	}
}
