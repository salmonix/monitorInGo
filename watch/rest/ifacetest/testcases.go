package ifacetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type testCase struct {
	name     string
	params   testParams
	expected cReturn
}

type testParams struct {
	uri    string
	params string
}

type cReturn struct {
	err  error
	resp *http.Response
}

func getTests() []testCase {
	r := []testCase{
		testCase{name: "Config", params: testParams{uri: "config", params: ""},
			expected: getExpReturn("nil", 200, ""), // put here a JSON
		},
	}
	return r
}

func getExpReturn(err string, statusCode int, body string) cReturn {

	h := map[string][]string{
		"Content-Type": []string{"application/json", "charset=utf-8"},
	}

	// make the string 'body' an IO readCloser: one example of ugly type lifting due to shallowness
	r := http.Response{StatusCode: statusCode, Body: ioutil.NopCloser(bytes.NewReader([]byte(body))), Header: h}

	if err == "nil" {
		return cReturn{err: nil, resp: &r}
	} else {
		return cReturn{err: fmt.Errorf("%s", err), resp: &r}
	}
}

func compareResponses(got, exp *http.Response) bool {
	gotMap, gotErr := getBodyJSON(got.Body)
	expMap, expErr := getBodyJSON(exp.Body)
	if gotErr != expErr {
		return false
	}

	// THIS IS NOT DEEP COMPARE
	for k, expV := range expMap {
		gotV, ok := gotMap[k]
		if ok == true {
			if gotV != expV {
				return false
			}
		}
	}
	return true
}

func getBodyJSON(r io.ReadCloser) (map[string]interface{}, error) {
	var gotStruct map[string]interface{}

	gotBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return gotStruct, err
	}

	json.Unmarshal(gotBytes, &gotStruct)
	return gotStruct, nil
}
