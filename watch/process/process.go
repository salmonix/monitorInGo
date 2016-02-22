package process

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

// Pid type
type Pid int

// WatchedProcess contains the process metrics and the PID of the process
type WatchedProcess struct {
	Mem  float64
	Vss  float64
	Rss  float64
	Pid  Pid
	Ppid Pid
}

// GetPPID returns the parent process ID for the process identified with its PID
// or 0, err if is does not exists
// XXX: we do not check for cases when no ppid exists
func GetPPID(p Pid) (Pid, error) {
	statPath := fmt.Sprintf("/proc/%d/statm", p)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return Pid(p), err
	}

	stat := strings.Fields(string(dataBytes))
	ppid, _ := strconv.ParseInt(stat[3], 10, 64)
	return Pid(ppid), nil
}

// NewWatchedProcess returns an initiated struct
func NewWatchedProcess(pid Pid, ppid Pid) *WatchedProcess {
	if ppid == 0 {
		ppid, _ = GetPPID(pid)
	}

	return &WatchedProcess{Pid: pid, Ppid: ppid}
}

// Update updates the process data. Currently it is a dummy, but future logic can be added here.
func (old *WatchedProcess) Update(new *WatchedProcess, tr float64) *WatchedProcess {
	return new
}

func overTreshold(x, y, tr float64) bool {
	diff := x / math.Abs(x-y)
	if diff > tr {
		return true
	}
	return false
}
