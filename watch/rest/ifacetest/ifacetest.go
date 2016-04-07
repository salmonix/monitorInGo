package ifacetest

import (
	"bytes"
	"fmt"
	"gmon/watch/config"
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

/* RunTests runs the test functions
func RunTests(t *restTest) {
	for v,_ in range(testCases) {



	}

}
*/

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
func NewRESTTest(c *config.Config) *restTest {

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
