package main

import (
	"os"
	"yaddd/cmd/yaddd/cmd"

	"github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"
)

const (
	versionCmd = "version"
	runCmd     = "run"
)

// compile passing -ldflags "-X yaddd.version <version> \
// -X yaddd.buildTime <YYYY-MM-DD HH:MM:SS MSK>".
var version, buildTime string //nolint:gochecknoglobals

func logVersion(log *logrus.Logger) {
	log.WithField("version", version).
		WithField("buildTime", buildTime).
		Info("YaDDD service started")
}

func run(args []string) {
	name := args[0]

	var err error

	var command string

	if len(args) > 1 {
		runFunc := func(cmdFunc func([]string) error) (err error) {
			logVersion(logrus.New())

			return cmdFunc(args)
		}

		command = args[1]

		switch command {
		case versionCmd:
			printVersion(os.Stdout)
		case runCmd:
			err = runFunc(cmd.Run)
		default:
			printHelp(os.Stderr, name)
		}
	} else {
		printHelp(os.Stderr, name)
	}

	if err != nil {
		logrus.WithField("cmd", command).
			WithError(err).
			Fatal("Execution problem")
	}
}

func main() {
	run(os.Args)
}
