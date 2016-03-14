package glog

import (
	"gmon/glog"
	"testing"
)

func Test_ProperObejct(t *testing.T) {
	l := glog.GetLogger()
	l.Info("FSFSF")
	// Output: gmon: 2016/03/12 00:02:50 glog_test.go:10: FSFSF
}
