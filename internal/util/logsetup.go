package util

import (
	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
)

func SetupLogs() {
	logrus.SetLevel(logrus.DebugLevel)
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	logrus.AddHook(filenameHook)
}
