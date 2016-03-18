package watch

import (
	"gmon/watch/process"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// WatchingContainer contains a read channel to receive instrucions
// and a map of id of processes to scan.
type WatchingContainer struct {
	processes map[int]*process.WatchedProcess
	treshold  float64
}

// Dummy is an empty watched process with -1 pid and no values.
var Dummy = process.NewWatchedProcess(-1, 0)

// NewContainer return a new *WatchingContainer
func NewContainer(tr float64) *WatchingContainer {
	processes := make(map[int]*process.WatchedProcess)
	Watch := &WatchingContainer{processes, float64(tr)}
	Watch.Add(os.Getpid())
	return Watch
}

// Add registers a process in the WatchingContainer
func (w *WatchingContainer) Add(p int) process.WatchedProcess {
	if _, ok := w.processes[p]; ok == false {
		np := process.NewWatchedProcess(p, 0)
		w.processes[p] = np
	}
	np, _ := w.processes[p]
	return *np
}

// Delete removes a process from the watchlist
func (w *WatchingContainer) Delete(p int) {
	if _, ok := w.processes[p]; ok == true {
		delete(w.processes, p)
	}
}

// Get returns the []struct for a watched process and true if the process exists
// a dummy value and false if not. If pid is <0 all the processes are returned.
func (w *WatchingContainer) Get(p int) ([]*process.WatchedProcess, bool) {

	if p < 0 {
		ret := make([]*process.WatchedProcess, len(w.processes))
		p := 0
		for _, v := range w.processes {
			ret[p] = v
			p++
		}
		return ret, false
	}
	ret := make([]*process.WatchedProcess, 1)
	if pr, ok := w.processes[p]; ok == true {
		ret[0] = pr
		return ret, true
	}
	ret[0] = Dummy
	return ret, false
}

// Refresh re-reads the ps table and refreshes the process data
func (w *WatchingContainer) Refresh() error {

	fields := 9
	//                     0  1    2   3    4    5   6    7    8
	commandString := "-eo cmd,pcpu,pid,ppid,size,rss,vsize,uid,utime"
	psTable, err := exec.Command("ps", commandString).Output() // check fields when it changes
	if err != nil {
		return err
	}

	// what about the - in the ps output?
	for _, r := range strings.Split(string(psTable), "\n") {
		psRow := strings.Fields(r)

		rowPid, _ := strconv.ParseInt(psRow[3], 10, 64) // not the best error handling ever
		if proc, ok := w.processes[int(rowPid)]; ok == true {

			asFloat := make([]float64, fields-2)
			for c := 2; c < fields; c++ {
				psVal, _ := strconv.ParseFloat(psRow[2], 64) // XXX can it be other than number?
				asFloat[c] = psVal
			}

			newStatus := process.NewWatchedProcess(proc.Pid, proc.Ppid)
			newStatus.Cmd = psRow[0]
			cpu, _ := strconv.ParseFloat(psRow[1], 32)
			newStatus.CPU = cpu
			newStatus.Mem = asFloat[4]
			newStatus.Rss = asFloat[5]
			newStatus.Vss = asFloat[6]
			newStatus.UID = int(asFloat[7])
			newStatus.Utime = int(asFloat[8])

			w.processes[proc.Pid] = proc.Update(newStatus, w.treshold)
		}
	}
	return nil
}
