package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gmon/glog"
	"gmon/watch"
	c "gmon/watch/config"
	"gmon/watch/process"
	"io/ioutil"
	"net/http"
	"strconv"
)

var l = glog.GetLogger("watch")

// GetRouter implements the REST interface to add, remove a query processes.
// For details see the Wiki ( https://bitbucket.org/monitoringo/monitoringo/wiki/Monitoring )
// TODO: missing DELETE
func GetRouter(w *watch.WatchingContainer, conf *c.Config) *gin.Engine {

	router := gin.Default()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("watch/templates/index.html")

	router.GET("/monitoring", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Monitoring WebUI",
		})
	})

	router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, conf)
	})

	router.GET("/processes", func(c *gin.Context) {
		if proc, err := w.Get(-1); err == true {
			c.JSON(http.StatusOK, proc)
		} else {
			c.JSON(http.StatusBadRequest, "Not found") // TODO: find a better error code
		}
	})

	router.GET("/processes/*id", func(c *gin.Context) {
		fmt.Printf("GET %s", c.Param("id"))
		pid := parseInt(c.Param("id"), c)
		proc, ok := w.Get(pid)
		if ok == true {
			c.JSON(http.StatusOK, proc)
			return
		}
		c.JSON(http.StatusExpectationFailed, "Unable to return process list")
		return

	})

	router.POST("/processes", func(c *gin.Context) {
		var process process.WatchedProcess
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &process)

		proc := w.Add(process.Pid, process.Ppid)
		c.JSON(http.StatusOK, proc)

	})

	return router
}

func parseInt(s string, c *gin.Context) int {
	if s == "/" {
		return -1
	}
	if i, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
		return int(i)
	}
	c.String(http.StatusBadRequest, "unable to serve data for %s", s)
	return 0
}
