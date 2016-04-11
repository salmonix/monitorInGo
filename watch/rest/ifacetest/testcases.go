package ifacetest

var testCases = []testCase{
	testCase{
		name: "Config", test: testParams{uri: "config", params: "", method: "POST"}, expected: getExpReturn("nil", 200, ""), // watchedProcess
	},
}
