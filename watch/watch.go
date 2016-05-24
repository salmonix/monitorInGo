package watch

import (
	"fmt"
	"gmon/ps"
	"gmon/watch/process"
	"os"
	"strconv"
	"strings"
)

// WatchingContainer contains a read channel to receive instrucions
// and a map of id of processes to scan.
// TODO: add a watch for the system itself
type WatchingContainer struct {
	processes map[int]*process.WatchedProcess
	treshold  float64
	// system    System // TODO: add system
}

// Dummy is an empty watched process with -1 pid and no values.
var Dummy, _ = process.NewWatchedProcess(-1, 0)

// NewContainer return a new *WatchingContainer
func NewContainer(tr float64) *WatchingContainer {
	processes := make(map[int]*process.WatchedProcess)
	Watch := &WatchingContainer{processes, float64(tr)}
	Watch.Add(os.Getpid(), 0) // register self
	return Watch
}

// Add registers a process in the WatchingContainer
func (w *WatchingContainer) Add(p, ppid int) (process.WatchedProcess, error) {
	if _, ok := w.processes[p]; ok == false {
		np, err := process.NewWatchedProcess(p, ppid)
		if err != nil {
			fmt.Printf("New process came back with error: %s", err)
			return *np, err
		}
		w.processes[p] = np
	}
	np, _ := w.processes[p]
	return *np, nil
}

// Delete removes a process from the watchlist
func (w *WatchingContainer) Delete(p int) {
	if _, ok := w.processes[p]; ok == true {
		delete(w.processes, p)
	}
}

// Get returns the ([]struct, ok) for a watched process and true if the process exists
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
		return ret, true
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

	psTable, err := ps.GetProcessTable(0)
	if err != nil {
		return err
	}

	for _, r := range strings.Split(string(psTable), "\n") {
		psRow := strings.Fields(r)

		rowPid, _ := strconv.ParseInt(psRow[3], 10, 64) // not the best error handling ever
		if proc, ok := w.processes[int(rowPid)]; ok == true {

			newStatus, err := process.NewWatchedProcess(proc.Pid, proc.Ppid)
			if err != nil {
				return err
			}

			w.processes[proc.Pid] = proc.Update(newStatus, w.treshold)
		}
	}
	return nil
}
