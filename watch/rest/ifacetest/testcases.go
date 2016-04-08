package ifacetest

var testCases = []testCase{
	testCase{
		name: "Config", params: testParams{uri: "config", params: ""}, expected: getExpReturn("nil", 200, ""),
	},
}
