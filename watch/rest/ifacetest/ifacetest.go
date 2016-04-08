package ifacetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gmon/watch/config"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

/* ifacetest package contains testing functions for the REST api of gmon.
The functions are used by testMonitoring.go test application with the proper CLI flags.
The testing functions follow the uri of the API
*/

type restTest struct {
	post func(string, string) cReturn
	put  func(string, string) cReturn
	get  func(string, string) cReturn
}

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

// Run runs the tests declared in testcases.go
func (t *restTest) Run() {
	for c, v := range testCases {
		fmt.Printf("Test %d : ", c)
		fmt.Println(v)
	}

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

func assert(got, exp cReturn, name string) bool {

	errorMess := "Failing Test: " + name + ":"
	// compare it via the compareRespExp
	if compareResponses(got.resp, exp.resp) == false {
		fmt.Println(errorMess + " respose differs: " + fmt.Sprintln(got) + " : " + fmt.Sprintln(exp))
	}
	fmt.Println(errorMess + " error differs: " + fmt.Sprintln(got) + " : " + fmt.Sprintln(exp))
	return false

}

// NewRESTTest returns a struct with the REST calls prepared
func NewRESTTest(c config.Config) *restTest {

	port := strconv.Itoa(c.Port)
	url := "http://localhost:" + port + "/"
	mime := "application/json"

	makeCReturn := func(r *http.Response, e error) cReturn {
		gRet := cReturn{err: e, resp: r}
		return gRet
	}

	test := restTest{
		post: func(uri string, body string) cReturn {
			toPost := bytes.NewBuffer([]byte(body))
			return makeCReturn(http.Post(url+uri, mime, toPost))
		},
		put: func(uri string, body string) cReturn {
			toPut := bytes.NewBuffer([]byte(body))
			request, _ := http.NewRequest("PUT", url+uri, toPut)
			return makeCReturn(http.DefaultClient.Do(request))
		},
		get: func(uri string, body string) cReturn {
			// toGet := bytes.NewBuffer([]byte(body))
			return makeCReturn(http.Get(url + uri + "/" + body))
		},
	}
	return &test
}

func getBody(r *http.Response) (string, bool) {
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return string(r.StatusCode), false
	}
	return string(contents), true
}
