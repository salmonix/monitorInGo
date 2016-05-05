package watch

import (
	"fmt"
	"gmon/watch/process"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// WatchingContainer contains a read channel to receive instrucions
// and a map of id of processes to scan.
// TODO: add a watch for the system itself
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
	Watch.Add(os.Getpid(), 0)
	return Watch
}

// Add registers a process in the WatchingContainer
// TODO : on init scan the process immediately
func (w *WatchingContainer) Add(p, ppid int) process.WatchedProcess {
	if _, ok := w.processes[p]; ok == false {
		np := process.NewWatchedProcess(p, ppid)
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

	fmt.Printf("Requested: %d", p)
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
		fmt.Printf("Requested: %d", p)
		ret[0] = pr
		return ret, true
	}
	ret[0] = Dummy
	return ret, false
}

// Refresh re-reads the ps table and refreshes the process data
func (w *WatchingContainer) Refresh() error {

	psTable, err := getProcessTable(0)
	if err != nil {
		return err
	}

	for _, r := range strings.Split(string(psTable), "\n") {
		psRow := strings.Fields(r)

		rowPid, _ := strconv.ParseInt(psRow[3], 10, 64) // not the best error handling ever
		if proc, ok := w.processes[int(rowPid)]; ok == true {

			newStatus := convertPS2Process(proc.Pid, proc.Ppid, psRow)

			w.processes[proc.Pid] = proc.Update(newStatus, w.treshold)
		}
	}
	return nil
}

// passing pid and ppid forces that these are the values that must be known
// at this point
func convertPS2Process(pid, ppid int, psRow []string) *process.WatchedProcess {
	fields := 9
	asFloat := make([]float64, fields-2)
	for c := 2; c < fields; c++ {
		psVal, _ := strconv.ParseFloat(psRow[2], 64) // XXX can it be other than number?
		asFloat[c] = psVal
	}

	newStatus := process.NewWatchedProcess(pid, ppid)
	newStatus.Cmd = psRow[0]
	cpu, _ := strconv.ParseFloat(psRow[1], 32)
	newStatus.CPU = cpu
	newStatus.Mem = asFloat[4]
	newStatus.ResidentSize = asFloat[5]
	newStatus.VirtualSize = asFloat[6]
	newStatus.UID = int(asFloat[7])
	newStatus.Utime = int(asFloat[8])
	return newStatus
}

func getProcessTable(p int) ([]byte, error) {

	var commandString string
	if p != 0 {
		commandString = "-p " + strconv.Itoa(p)
	} else {
		commandString = ""
	}
	//                     0  1    2   3    4    5   6    7    8
	commandString += "-eo cmd,pcpu,pid,ppid,size,rss,vsize,uid,utime"
	psTable, err := exec.Command("ps", commandString).Output() // check fields when it changes
	if err != nil {
		return psTable, err
	}
	return psTable, nil
}
