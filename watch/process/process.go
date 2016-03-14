package process

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

// WatchedProcess contains the process metrics and the PID of the process
type WatchedProcess struct {
	Mem  float64
	Vss  float64
	Rss  float64
	Pid  int
	Ppid int
}

// GetPPID returns the parent process ID for the process identified with its PID
// or 0, err if is does not exists
// XXX: we do not check for cases when no ppid exists
func GetPPID(p int) (int, error) {
	statPath := fmt.Sprintf("/proc/%d/statm", p)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return p, err
	}

	stat := strings.Fields(string(dataBytes))
	ppid, _ := strconv.ParseInt(stat[3], 10, 64)
	return int(ppid), nil
}

// NewWatchedProcess returns an initiated struct
func NewWatchedProcess(pid int, ppid int) *WatchedProcess {
	if ppid == 0 {
		ppid, _ = GetPPID(pid)
	}

	return &WatchedProcess{Pid: pid, Ppid: ppid}
}

// Update returns the updated pointer using the th int as treshold for changes.
func (old *WatchedProcess) Update(new *WatchedProcess, tr float64) *WatchedProcess {
	if overTreshold(old.Mem, new.Mem, tr) {
		return new
	}
	if overTreshold(old.Rss, new.Rss, tr) {
		return new
	}
	return old
}

func overTreshold(x, y, tr float64) bool {
	diff := x / math.Abs(x-y)
	if diff > tr {
		return true
	}
	return false
}
