package ifacetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gmon/watch/config"
	p "gmon/watch/process"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

/* ifacetest package contains testing functions for the REST api of gmon.
The functions are used by testMonitoring.go test application with the proper CLI flags.
The testing functions follow the uri of the API
*/

type testCase struct {
	name     string
	test     testParams
	expected expected
}

// XXX FIXME make a switch function instead of that fancy dynamic function thingy.
// just refacta that
type testParams struct {
	uri    string
	params string
	method string
}

type expected struct {
	err     error
	res     p.WatchedProcess
	resList []p.WatchedProcess
}

type responser func(testCase) (string, bool)

// RestTest contains the main HTTP wrapper
type RestTest struct {
	response responser
}

// NewRESTTest creates a test struct using the Config values
func NewRESTTest(c *config.Config) RestTest {
	return RestTest{response: getResponser(c)}
}

// Run runs the tests declared in testcases.go
func (t *RestTest) Run() {
	for c, tCase := range testCases {
		fmt.Printf("Test %d : ", c)
		if resp, ok := t.response(tCase); ok == true { // call test passed
			assert(resp, tCase)
		}
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

func assert(got, exp expected, name string) bool {

	errorMess := "Failing Test: " + name + ":"
	// compare it via the compareRespExp
	if compareResponses(got.resp, exp.resp) == false {
		fmt.Println(errorMess + " respose differs: " + fmt.Sprintln(got) + " : " + fmt.Sprintln(exp))
	}
	fmt.Println(errorMess + " error differs: " + fmt.Sprintln(got) + " : " + fmt.Sprintln(exp))
	return false

}

func getResponser(c *config.Config) responser {

	port := strconv.Itoa(c.Port)
	url := "http://localhost:" + port + "/"
	mime := "application/json"

	return responser(
		func(t testCase) (string, bool) {

			var response *http.Response
			var err error
			uri := t.test.uri

			b := t.test.params
			body := bytes.NewBuffer([]byte(t.test.params))
			switch t.test.method {
			case "post":
				response, err := http.Post(url+uri, mime, body)
			case "put":
				request, _ := http.NewRequest("PUT", url+uri, body)
				response, err = http.DefaultClient.Do(request) // TODO: is it correct?
			case "get":
				response, err := http.Get(url + uri + "/" + t.test.params)
			default:
				panic("Not handled method in getResponse!")
			}
			if err != nil {
				fmt.Println(err)
				return "", false
			}
			return getBody(response)
		})
}

//
func getBody(r *http.Response) (string, bool) {
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return string(r.StatusCode), false
	}
	return string(contents), true
}
