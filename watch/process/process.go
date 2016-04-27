package process

import (
	"fmt"
	"gmon/ps"
	"math"
)

// WatchedProcess contains the process metrics and the PID of the process
type WatchedProcess struct {
	Cmd                 string
	CPU                 float64
	Pid                 int
	Ppid                int
	Mem                 float64
	VirtualSizeMb       float64
	VirtualSizePercent  float64
	ResidentSizeMb      float64
	ResidentSizePercent float64 // TODO: how shall I implement this?
	UID                 int
	Utime               int
	State               string // enum: R is running, S is sleeping, D is sleeping in an uninterruptible wait, Z is zombie, T is traced or stopped
	Checked             int
	Children            []int
}

// NewWatchedProcess returns an initiated struct  -reading the PS data from the PS table
// If PID == 0 the data refers to the host.
func NewWatchedProcess(pid int, ppid int) (*WatchedProcess, error) {

	var newProcess WatchedProcess

	if pid < 0 {
		return &newProcess, nil
	}

	psRaw, err := ps.GetProcessTable(pid)
	if err != nil {
		return &newProcess, err
	}

	fmt.Scanf(ps.PsMask, psRaw, &newProcess)
	return &newProcess, nil

}

// Update returns the updated pointer using the int as treshold for changes.
func (old *WatchedProcess) Update(new *WatchedProcess, tr float64) *WatchedProcess {
	if overTreshold(old.Mem, new.Mem, tr) {
		return new
	}
	if overTreshold(old.ResidentSizeMb, new.ResidentSizeMb, tr) {
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
