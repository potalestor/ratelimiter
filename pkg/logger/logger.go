package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const fmtAppInfo = `
	%s starting...

`

// Initialize logger using cfg.
func Initialize(logfile string) {
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open/create log file '%s'", logfile)
		}

		multi := io.MultiWriter(file, os.Stdout)
		log.SetOutput(multi)
	}

	log.SetFlags(0)

	log.Printf(fmtAppInfo, strings.ToUpper(filepath.Base(os.Args[0])))
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
