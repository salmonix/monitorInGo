package process

import (
	"fmt"
	conv "gmon/conversions"
	"gmon/glog"
	"gmon/ps"
	"math"
	"time"
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

var l = glog.GetLogger("watch")

// NewWatchedProcess returns an initiated struct  -reading the PS data from the PS table
// If PID == 0 the data refers to the host.
func NewWatchedProcess(pid int, ppid int) (*WatchedProcess, error) {

	var newProcess WatchedProcess
	l.Debug("NewWatchedProcess is requested for pid ", pid)

	if pid < 0 {
		l.Debug("Pid is smaller than zero")
		return &newProcess, nil
	}

	psRaw, err := ps.GetProcessTable(pid)
	if err != nil {
		panic(err)
		// return &newProcess, err
	}

	// TODO: FIXME: Still wrong parsing
	// var psString = "pid,cmd,ppid,vsize,rss,pcpu,size,uid,utime,state"
	// var PsMask = "%d %s %d %d %d %d %d %d %s %s"
	var vsize, rss, pcpu, size, uid int
	var cmd, utime, state string
	c, err := fmt.Sscan(string(psRaw), ps.PsMask, &pid, &cmd, &ppid, &vsize, &rss, &pcpu, &size, &uid, &utime, &state)
	l.Debug("Scanned", c, "element of the 10")
	if err != nil {
		panic(err)
		//return &newProcess, err
	}

	// TODO: I need the system mem size from 0
	newProcess.Cmd = cmd
	newProcess.CPU = float64(pcpu)
	newProcess.Pid = pid
	newProcess.Ppid = ppid
	//newProcess.Mem =  // TODO calculate
	newProcess.VirtualSizeMb = conv.BytesToMb(float64(vsize))
	// newProcess.VirtualSizePercent =
	newProcess.ResidentSizeMb = conv.BytesToMb(float64(rss))
	// newProcess.ResidentSizePercent // TODO calculate
	newProcess.UID = uid

	if uT, err := time.Parse("00:00:00", utime); err == nil {
		newProcess.Utime = int(uT.Unix())
	} else {
		newProcess.Utime = -1
	}

	newProcess.State = state
	// Checked = int // TODO: not implemented
	// Children =             []int // TODO: not implemented
	l.Warning("Returning from new process")
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
