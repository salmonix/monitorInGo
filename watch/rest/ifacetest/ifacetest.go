package ifacetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gmon/watch/config"
	p "gmon/watch/process"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

/* ifacetest package contains testing functions for the REST api of gmon.
The functions are used by testMonitoring.go test application with the proper CLI flags.
The testing functions follow the uri of the API
*/

type testCase struct {
	name     string
	test     testParams
	expected responseContainer
}

type testParams struct {
	uri    string
	params string
	method string
}

type responseContainer struct {
	err     error
	res     p.WatchedProcess
	resList []p.WatchedProcess
}

type responser func(testCase) (responseContainer, bool)

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
			assert(resp, tCase.expected, tCase.name)
		}
	}
}

func assert(got, exp responseContainer, name string) bool {

	errorMess := "Failing Test: " + name + ":"
	// compare it via the compareRespExp
	if compareContainers(got, exp) == false {
		fmt.Println(errorMess + " respose differs: " + fmt.Sprintln(got) + " : " + fmt.Sprintln(exp))
	}
	fmt.Println(errorMess + " error differs: " + fmt.Sprintln(got) + " : " + fmt.Sprintln(exp))
	return false

}

func compareContainers(got, exp responseContainer) bool {

	dummy := responseContainer{}
	if got.err == exp.err {
		if exp.err != nil {
			if reflect.DeepEqual(exp.res, dummy.res) { // reflect DeepEqual
				return compareStructs(got.res, exp.res)
			}
			if reflect.DeepEqual(exp.resList, dummy.resList) == false {
				for c, e := range exp.resList {
					g := got.resList[c]
					if compareStructs(e, g) == false {
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

func compareStructs(got, exp p.WatchedProcess) bool {
	gotValues := reflect.ValueOf(got)
	expValues := reflect.ValueOf(exp)
	identical := true

	for i := 0; i < gotValues.NumField(); i++ {
		e := expValues.Field(i)
		g := gotValues.Field(i)
		if e != g {
			identical = false
		}
	}
	return identical
}

func getResponser(c *config.Config) responser {

	port := strconv.Itoa(c.Port)
	url := "http://localhost:" + port + "/"
	mime := "application/json"

	return responser(
		func(t testCase) (responseContainer, bool) {

			var response *http.Response
			var err error
			uri := t.test.uri

			body := bytes.NewBuffer([]byte(t.test.params))

			switch t.test.method {
			case "post":
				response, err = http.Post(url+uri, mime, body)
			case "put":
				request, _ := http.NewRequest("PUT", url+uri, body)
				response, err = http.DefaultClient.Do(request) // TODO: is it correct?
			case "get":
				response, err = http.Get(url + uri + "/" + t.test.params)
			default:
				panic("Not handled method in getResponse!")
			}
			if err != nil {
				fmt.Println(err)
				return responseContainer{}, false
			}
			return getBody(response)
		})
}

func getBody(r *http.Response) (responseContainer, bool) {
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return responseContainer{err: err}, true
	}

	var resList []p.WatchedProcess
	if err := json.Unmarshal(contents, &resList); err == nil {
		return responseContainer{err: nil, resList: resList}, true
	}

	var res p.WatchedProcess
	if err := json.Unmarshal(contents, &res); err == nil {
		return responseContainer{err: nil, res: res}, true
	}
	fmt.Printf("%s", err)
	return responseContainer{}, false
}
