package watch

import (
	"bytes"
	"gmon/glog"
	"gmon/ps"
	"gmon/watch/process"
	"os"
	"strconv"
	"strings"
)

var l = glog.GetLogger("watch")

// WatchingContainer contains a read channel to receive instrucions
// and a map of id : processes to scan.
// TODO: add a watch for the system itself
type WatchingContainer struct {
	Processes map[int]*process.WatchedProcess
	treshold  uint32
	// system    System // TODO: add system parameters
}

// Dummy is a not valid watched process with -1 pid and no values.
var Dummy = process.NewWatchedProcess(-1, 0)

// NewContainer return a new *WatchingContainer
func NewContainer(tr uint32) *WatchingContainer {
	processes := make(map[int]*process.WatchedProcess)
	watch := &WatchingContainer{processes, tr}
	watch.Add(os.Getpid(), 0)
	return watch
}

// Add registers a process in the WatchingContainer returning the new process
func (w *WatchingContainer) Add(p, ppid int) *process.WatchedProcess {
	if _, ok := w.Processes[p]; ok == false {
		np := process.NewWatchedProcess(p, ppid)
		w.Processes[p] = np
	}
	np, _ := w.Processes[p]
	return np
}

// Delete removes a process from the watchlist
func (w *WatchingContainer) Delete(p int) {
	if _, ok := w.Processes[p]; ok == true {
		delete(w.Processes, p)
	}
}

// Get returns the ([]struct, ok) for a watched process and true if the process exists
// a dummy value and false if not. If pid is <0 all the processes are returned.
// XXX that might not be a good idea
func (w *WatchingContainer) Get(p int) ([]*process.WatchedProcess, bool) {

	if p < 0 {
		ret := make([]*process.WatchedProcess, len(w.Processes))
		p := 0
		for _, v := range w.Processes {
			ret[p] = v
			p++
		}
		return ret, true
	}
	ret := make([]*process.WatchedProcess, 1)
	if pr, ok := w.Processes[p]; ok == true {
		ret[0] = pr
		return ret, true
	}
	ret[0] = Dummy
	return ret, false
}

// Refresh re-reads the ps table and refreshes the process data
func (w *WatchingContainer) Refresh() error {

	var pids []string // read only those processes that are
	for k, _ := range w.Processes {
		pids = append(pids, strconv.Itoa(k))
	}

	psTable, err := ps.GetProcessTable(strings.Join(pids, ","))

	if err != nil {
		l.Warning(err, pids)
		return err
	}

	newProcesses := make(map[int]*process.WatchedProcess)

	for _, psRow := range bytes.Split(bytes.TrimSpace(psTable), []byte("\n")) {
		newState := new(process.WatchedProcess)
		newState.CastPsRow2Process(psRow)

		if err != nil { // ps is gone...?
			l.Warning("Houston, we have a problem")
			return err
		}

		if proc, ok := w.Processes[newState.Pid]; ok == true {
			newProcesses[proc.Pid] = proc.Update(newState, w.treshold)
			delete(w.Processes, proc.Pid)
		}
	}
	// remove all remaining processes are gone, so we delete them
	for k, proc := range w.Processes {
		if proc.State == "" {
			delete(w.Processes, k)
		}
	}
	w.Processes = newProcesses
	return nil
}
