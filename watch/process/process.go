package process

import (
	"gmon/ps"
	"math"
	"strconv"
	"strings"
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

	psRow := strings.Fields(string(psRaw))
	fields := len(psRow)
	asFloat := make([]float64, fields-2)
	for c := 2; c < fields; c++ {
		psVal, _ := strconv.ParseFloat(psRow[2], 64) // XXX can it be other than number?
		asFloat[c] = psVal
	}

	newProcess.Pid = pid
	newProcess.Ppid = ppid
	newProcess.Cmd = psRow[0]
	cpu, _ := strconv.ParseFloat(psRow[1], 32)
	newProcess.CPU = cpu
	newProcess.Mem = asFloat[4]
	newProcess.ResidentSizeMb = asFloat[5]
	newProcess.VirtualSizeMb = asFloat[6]
	newProcess.UID = int(asFloat[7])
	newProcess.Utime = int(asFloat[8])
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
