package ps

import (
	"os/exec"
	"strconv"
)

// GetProcessTable calls ps with or without the pid parameter
// returning the output as []byte.
func GetProcessTable(p int) ([]byte, error) {

	if p > 0 {
		// FIXME  cd this is still incorrect
		return exec.Command("ps", "h", "-q", strconv.Itoa(p), "-eo", "pid,cmd,ppid,vsize,rss,pcpu,size,uid,utime,state").Output()
	}

	return exec.Command("ps", "h", "-eo", "pid,cmd,ppid,vsize,rss,pcpu,size,uid,utime,state").Output()
}
