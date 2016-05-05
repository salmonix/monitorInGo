package process

import (
	"fmt"
	"gmon/conversions"
	"io/ioutil"
	"math"
	"strings"
)

// WatchedProcess contains the process metrics and the PID of the process
type WatchedProcess struct {
	Cmd            string
	CPU            float64
	Pid            int
	Ppid           int
	Mem            float64
	VirtualSize    float64
	VirtualSizeMb  float32
	ResidentSize   float64
	ResidentSizeMb float32
	UID            int
	Utime          int
	Checked        int
	Children       []int
}

// GetPPID returns the parent process ID for the process identified with its PID
// or 0, err if is does not exists
// XXX: we do not check for cases when no ppid exists
func GetPPID(p int) (int, error) {
	stat, err := readStatm(p)
	if err != nil {
		return -1, err
	}
	ppid := stat[3]
	return int(ppid), nil
}

// NewWatchedProcess returns an initiated struct
func NewWatchedProcess(pid int, ppid int) *WatchedProcess {
	if ppid == 0 {
		ppid, _ = GetPPID(pid)
	}

	return &WatchedProcess{Pid: pid, Ppid: ppid}
}

// Update returns the updated pointer using the int as treshold for changes.
func (old *WatchedProcess) Update(new *WatchedProcess, tr float64) *WatchedProcess {
	if overTreshold(old.Mem, new.Mem, tr) {
		return new
	}
	if overTreshold(old.ResidentSize, new.ResidentSize, tr) {
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

func readStatm(p int) ([]int, error) {
	statPath := fmt.Sprintf("/proc/%d/statm", p)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return nil, err
	}
	stat, err := conversions.StringsToInt(strings.Fields(string(dataBytes)))
	return stat, err
}
