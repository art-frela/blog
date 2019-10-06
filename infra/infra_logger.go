package infra

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func (bs *BlogServer) setLogger(version string) {
	log := logrus.WithFields(logrus.Fields{
		"program":     bs.config.GetString("app.name"),
		"hostname":    bs.config.GetString("httpd.host"),
		"version":     version,
		"environment": bs.config.GetString("env"),
	})

	logrus.SetReportCaller(true)
	if bs.config.GetString("env") == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		log.Logger.Formatter = &logrus.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
	}
	llevel, err := logrus.ParseLevel(bs.config.GetString("log.level"))
	if err != nil {
		log.Errorf("set log level error, %v", err)
		llevel = logrus.InfoLevel

	}
	logrus.SetLevel(llevel) // set loglevel

	// set output
	if bs.config.GetString("log.file") != "" || bs.config.GetString("log.file") != "stdout" {
		f, err := os.Create(bs.config.GetString("log.file"))
		if err != nil {
			log.Errorf("open log file [%s] error, %v", err)
		} else {
			logrus.SetOutput(f)
		}

	}
	bs.log = log
}
