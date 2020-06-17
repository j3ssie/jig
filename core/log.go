package core

import (
	"fmt"
	"io/ioutil"
	"os"
	// "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = logrus.New()

// InitLog init log
func InitLog(options *Options) {
	logger.SetLevel(logrus.InfoLevel)
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &prefixed.TextFormatter{
			ForceColors:     true,
			ForceFormatting: true,
		},
	}

	if options.Debug == true {
		logger.SetLevel(logrus.DebugLevel)
	} else if options.Verbose == true {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetOutput(ioutil.Discard)
	}
}

// GoodF print good message
func GoodF(format string, args ...interface{}) {
	good := color.HiGreenString("[+]")
	fmt.Printf("%s %s\n", good, fmt.Sprintf(format, args...))
}

// BannerF print info message
func BannerF(format string, data string) {
	banner := color.BlueString("[*] %v", format)
	logger.Info(fmt.Sprintf("%v%v", banner, color.HiGreenString(data)))
}

// InforF print info message
func InforF(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

// WarningF print good message
func WarningF(format string, args ...interface{}) {
	good := color.YellowString("[!]")
	fmt.Printf("%s %s\n", good, fmt.Sprintf(format, args...))
}

// DebugF print debug message
func DebugF(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

// ErrorF print good message
func ErrorF(format string, args ...interface{}) {
	good := color.RedString("[-]")
	fmt.Printf("%s %s\n", good, fmt.Sprintf(format, args...))
}
