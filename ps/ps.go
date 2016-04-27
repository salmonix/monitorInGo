package ps

import (
	"os/exec"
	"strconv"
)

var psString = "pid,cmd,ppid,vsize,rss,pcpu,size,uid,utime,state"

// PsMask is scanf mask to convert ps []byte into WatchedProcess
var PsMask = "%d %s %d %d %d %d %d %d %s"

// GetProcessTable calls ps with or without the pid parameter
// returning the output as []byte.
func GetProcessTable(p int) ([]byte, error) {

	if p > 0 {
		return exec.Command("ps", "h", "-q", strconv.Itoa(p), "-eo", psString).Output()
	}

	return exec.Command("ps", "h", "-eo", psString).Output()
}
