package process

import (
	"fmt"
	conv "gmon/conversions"
	"gmon/glog"
	"gmon/ps"
	"strconv"
	"time"
)

// WatchedProcess contains the process metrics and the PID of the process
type WatchedProcess struct {
	Cmd            string
	CPU            float64
	Pid            int
	Ppid           int
	Mem            uint32
	VirtualSizeMb  uint32
	ResidentSizeMb uint32
	UID            int
	RegisteredAt   int32  // these both are unix ts.
	StartedAt      int32  // might be 0 if process does not provide this info
	State          string // enum: R is running, S is sleeping, D is sleeping in an uninterruptible wait, Z is zombie, T is traced or stopped
	Checked        int
	Children       []int
	Tags           string // a # separated list of strings for flat custom metadata
}

var l = glog.GetLogger("watch")

// NewWatchedProcess returns an initiated struct - reading the PS data from the PS table
// If PID == 0 the data refers to the host.
func NewWatchedProcess(pid int, ppid int) *WatchedProcess {

	newProcess := new(WatchedProcess)
	if pid < 0 {
		return newProcess
	}
	psRaw, err := ps.GetProcessTable(strconv.Itoa(pid)) // XXX Test for not existing PID

	if err != nil {
		l.Warning(err)
		return newProcess
	}

	newProcess.CastPsRow2Process(psRaw)
	newProcess.RegisteredAt = int32(time.Now().Unix())
	return newProcess
}

// Update returns the updated pointer using the int as treshold for changes.
func (old *WatchedProcess) Update(new *WatchedProcess, tr uint32) *WatchedProcess {
	if overTreshold(old.Mem, new.Mem, tr) {
		return new
	}
	if overTreshold(old.ResidentSizeMb, new.ResidentSizeMb, tr) {
		return new
	}
	return old
}

func overTreshold(x, y, tr uint32) bool {
	diff := (x - y)

	if diff < 0 {
		diff = -diff
	}

	if diff == 0 {
		return false
	}

	diffPerc := x / diff
	if diffPerc > tr {
		return true
	}
	return false
}

// CastPsRow2Process casts the row from the ps output into the process struct.
func (p *WatchedProcess) CastPsRow2Process(psRaw []byte) {

	var uid int
	var vsize, size, rss uint64
	c, _ := fmt.Sscan(string(psRaw), &p.Pid, &p.Cmd, &p.Ppid, &vsize, &rss, &p.CPU, &size, &uid, &p.State)

	if c != 9 { // this is for a not-yet understood failure possibility
		l.Warning("Not enough parameters got from psRaw:", c)
		l.Debug("Scanned", c, "element of the 10")
		l.Debug(string(psRaw))
		p.Pid = -1
		return
	}
	// TODO: I need the system mem size from 0

	p.Mem = conv.KBytesToMb(size)
	p.VirtualSizeMb = conv.KBytesToMb(vsize)
	// p.VirtualSizePercent =
	p.ResidentSizeMb = conv.KBytesToMb(rss)
	// p.ResidentSizePercent
	p.UID = uid

	// Checked = int
	// Children =  []int
}
