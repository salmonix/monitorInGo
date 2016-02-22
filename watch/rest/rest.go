package rest

import (
	"fmt"
	"gmon/watch"
	"gmon/watch/process"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRouter implements the REST interface to add, remove a query processes.
// POST: /process/:name?pid=int
//    add a new process to the process list. The process must be identified by PID
//    at later calls. If the process is already exist, nothing happens
// GET: /process/:pid
//    Get the metrics of the process
// DELETE: /process/:pid
//    Delete the process.
// GET: /process
//    Return a [process] list with the process data.
// TODO: cat /proc/sys/kernel/pid_max returns the max of PID. we should add it to ParseInt
func GetRouter(w *watch.WatchingContainer) *gin.Engine {

	router := gin.Default()
	router.GET("/process/:id", func(c *gin.Context) {
		pid := parseInt(c.Param("id"), c)
		proc, _ := w.Get(process.Pid(pid))
		c.JSON(http.StatusOK, proc)
	})

	router.POST("/process/:id", func(c *gin.Context) {
		pid := parseInt(c.Param("id"), c)
		c.JSON(http.StatusOK, w.Add(process.Pid(pid)))
	})

	return router
}

func parseInt(s string, c *gin.Context) int {
	if i, err := strconv.ParseInt(c.Param("name"), 10, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
		return int(i)
	}
	c.String(http.StatusBadRequest, "unable to serve data for %s", s)
	return 0
}
