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

	psTable, err := exec.Command("ps", "-aux").Output()
	if err != nil {
		return err
	}

	// ps aux output is USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND
	// this is a conversion, and may be an interface if different types are used
	psPositions := []int{1, 3, 4, 5} // NOTE: pos 10 will be the cmdline
	for _, r := range strings.Split(string(psTable), "\n") {
		psRow := strings.Fields(r)
		prDat := make([]float64, 5)

		if proc, ok := w.processes[int(prDat[0])]; ok == true {

			for c, psDat := range psPositions {
				psVal, _ := strconv.ParseFloat(psRow[psDat], 64) // XXX can it be other than number?
				prDat[c] = psVal
			}

			newStatus := process.NewWatchedProcess(proc.Pid, proc.Ppid)
			newStatus.Mem = prDat[1]
			newStatus.Vss = prDat[2]
			newStatus.Rss = prDat[3]
			w.processes[proc.Pid] = proc.Update(newStatus, w.treshold)
		}
	}
	return nil
}
