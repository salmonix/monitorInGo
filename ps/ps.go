package ps

import (
	"fmt"
	"os/exec"
	"strconv"
)

var psString = "pid,cmd,ppid,vsize,rss,pcpu,size,uid,etime,state"

// GetProcessTable calls ps with or without the pid parameter
// returning the output as []byte.
func GetProcessTable(p int) ([]byte, error) {

	if p > 0 {
		return exec.Command("ps", "h", "-q", strconv.Itoa(p), "-eo", psString).Output()
	}

	return exec.Command("ps", "h", "-eo", psString).Output()
}

// GetPsPID returns the PID for the process in the psRow
func GetPsPID(psRow []byte) (int, error) {
	var pid int
	if _, err := fmt.Sscan(string(psRow), &pid); err != nil {
		return 0, err
	}
	return pid, nil
}
