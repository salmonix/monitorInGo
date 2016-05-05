package rest

import (
	"fmt"
	"gmon/watch"
	c "gmon/watch/config"
	"gmon/watch/process"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRouter implements the REST interface to add, remove a query processes.
// For details see the Wiki ( https://bitbucket.org/monitoringo/monitoringo/wiki/Monitoring )
func GetRouter(w *watch.WatchingContainer, conf *c.Config) *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("watch/templates/index.html")

	router.GET("/monitoring", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Monitoring WebUI",
		})
	})

	router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, conf)
	})

	router.GET("/processes/*id", func(c *gin.Context) {
		fmt.Printf("GET %s", c.Param("id"))
		pid := parseInt(c.Param("id"), c)
		proc, _ := w.Get(pid)
		c.JSON(http.StatusOK, proc)
	})

	router.POST("/processes", func(c *gin.Context) {
		var process process.WatchedProcess
		c.BindJSON(&process)
		c.JSON(http.StatusOK, w.Add(process.Pid, process.Ppid))
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
