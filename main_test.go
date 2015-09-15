package main

import (
	`fmt`
	`io/ioutil`
	`net/http`
	`strings`
	`testing`
)

const (
	LocalHost = "http://localhost:8088/cache/"
)

type CachePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (cp CachePair) toJson() string {
	return fmt.Sprintf("{\"key\": \"%s\", \"value\": \"%s\"}", cp.Key, cp.Value)
}

func TestChallenge1Post(t *testing.T) {
	// Start from a clean slate.
	deleteAllAndAssert(t)

	// Populate a single key / value pair via a POST call.
	cp := &CachePair{Key: "falfa", Value: "fbeta"}

	// Status code 201 returned with a link to new cache item's key.
	postAndAssert(t, cp.toJson())

}

func deleteAllAndAssert(t *testing.T) {
	expectedStatus := http.StatusNoContent
	req, err := http.NewRequest("DELETE", LocalHost, nil)
	if err != nil {
		t.Errorf("Delete call failed: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
	}

	if resp.StatusCode != expectedStatus {
		t.Errorf("Response code of %s doesn't match expected %s.", resp.StatusCode, expectedStatus)
	}
}

func postAndAssert(t *testing.T, json string) {
	expectedStatus := http.StatusCreated
	resp, err := http.Post(LocalHost, "application/json", strings.NewReader(json))
	if err != nil {
		t.Errorf("Initial connection failed: %s", err)
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
	}

	if resp.StatusCode != expectedStatus {
		t.Errorf("Response code of %s doesn't match expected %s.", resp.StatusCode, expectedStatus)
	}
}
