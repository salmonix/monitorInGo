package ifacetest

import "os"
import p "gmon/watch/process"

var testCases = []testCase{
	testCase{
		name: "Config", test: testParams{uri: "config", params: "", method: "POST"},
		expected: responseContainer{err: nil, res: p.WatchedProcess{Pid: os.Getpid()}},
	},
}
