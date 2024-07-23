package main

import (
	"flag"
	"fmt"
	"log/slog"
)

// handleFlags is a function that handles the command line flags. The only flag
// currently processed is -l (logging level) which sets the log level to
// error (e), warning (w), information (i) or debug (d). The function returns
// the type slog.Level and an error. If logging level is incorrectly an error
// is returned.
// In the future, this function may support more flags such as WebMode and input
// and output specifiers.
func handleFlags() (slog.Level, bool, error) {
	logLevelPtr := flag.String("l", "w", "used to set the logging level, may be 'e' error, 'w' warning, 'i' information or 'd' debug")
	printCfgPtr := flag.Bool("printConfig", false, "used to print to screen and then quit, no output file will be generated")
	flag.Parse()

	if *printCfgPtr {
		/* When return here the log level is always set to info*/
		return slog.LevelWarn, true, nil
	}

	switch *logLevelPtr {
	case "e":
		return slog.LevelError, false, nil
	case "w":
		return slog.LevelWarn, false, nil
	case "i":
		return slog.LevelInfo, false, nil
	case "d":
		return slog.LevelDebug, false, nil
	default:
		return slog.LevelDebug, false, fmt.Errorf("%s not supported for -l", *logLevelPtr)
	}
}
