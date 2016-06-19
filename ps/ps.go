package ps

import (
	"os/exec"
)

// ps : the current approach is simply calling system ps and parses the return value.
// That could be implemented as a go library to save the call cost, but it seems that
// it will not going to cause problem to call the C ps in every x seconds on a mem fs.

// PsString are the parameters we read from the ps command
var PsString = "pid,cmd,ppid,vsize,rss,pcpu,size,uid,etime,state"

// GetProcessTable calls ps with or without the pid parameter
// returning the output as []byte.
func GetProcessTable(p string) ([]byte, error) {
	return exec.Command("ps", "h", "-q", p, "-eo", PsString).Output()
}
