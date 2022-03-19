// Package dexlogconfig is a policy package to configure xlog as I like it.
package dexlogconfig

import "github.com/hlandau/xlog"
import "gopkg.in/hlandau/easyconfig.v1/cflag"
import "os"
import _ "github.com/hlandau/buildinfo"
import "gopkg.in/hlandau/svcutils.v1/systemd"

var (
	flagGroup          = cflag.NewGroup(nil, "xlog")
	logSeverityFlag    = cflag.String(flagGroup, "severity", "NOTICE", "Log severity (any syslog severity name or number (0-7) or 'trace' (8) (most verbose))")
	logFileFlag        = cflag.String(flagGroup, "file", "", "Log to filename")
	fileSeverityFlag   = cflag.String(flagGroup, "fileseverity", "TRACE", "File logging severity limit")
	logStderrFlag      = cflag.Bool(flagGroup, "stderr", true, "Log to stderr?")
	stderrSeverityFlag = cflag.String(flagGroup, "stderrseverity", "TRACE", "stderr logging severity limit")
)

func openStderr() {
	if logStderrFlag.Value() {
		if sev, ok := xlog.ParseSeverity(stderrSeverityFlag.Value()); ok {
			xlog.StderrSink.SetSeverity(sev)
		}

		if systemd.IsRunningUnder() {
			xlog.StderrSink.Systemd = true
		}

		return
	}

	xlog.RootSink.Remove(xlog.StderrSink)
}

func openFile() {
	fn := logFileFlag.Value()
	if fn == "" {
		return
	}

	f, err := os.Create(fn)
	if err != nil {
		return
	}

	sink := xlog.NewWriterSink(f)
	if sev, ok := xlog.ParseSeverity(fileSeverityFlag.Value()); ok {
		sink.SetSeverity(sev)
	}

	xlog.RootSink.Add(sink)
}

func setSeverity() {
	sevs := logSeverityFlag.Value()
	sev, ok := xlog.ParseSeverity(sevs)
	if !ok {
		return
	}

	xlog.Root.SetSeverity(sev)

	xlog.VisitSites(func(s xlog.Site) error {
		s.SetSeverity(sev)
		return nil
	})
}

// Parse registered configurables and setup logging.
func Init() {
	setSeverity()
	openStderr()
	openSyslog()
	openJournal()
	openEventLog()
	openFile()
}
