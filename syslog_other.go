// +build windows

package dexlogconfig

import "github.com/hlandau/xlog"
import "gopkg.in/hlandau/easyconfig.v1/cflag"
import "gopkg.in/hlandau/svcutils.v1/exepath"
import "github.com/btcsuite/winsvc/eventlog"
import "fmt"

var (
	eventLogFlag         = cflag.Bool(flagGroup, "eventlog", false, "Log to event log?")
	eventLogSeverityFlag = cflag.String(flagGroup, "eventlogseverity", "DEBUG", "Event log severity limit")
)

func openEventLog() {
	var err error

	if !eventLogFlag.Value() {
		return
	}

	pn := exepath.ProgramName
	esink.Log, err = eventlog.Open(pn)
	if err != nil {
		return
	}

	if sev, ok := xlog.ParseSeverity(eventLogSeverityFlag.Value()); ok {
		esink.MinSeverity = sev
	}

	xlog.RootSink.Add(&esink)
}

type eventSink struct {
	Log         *eventlog.Log
	MinSeverity xlog.Severity
}

func (s *eventSink) ReceiveLocally(sev xlog.Severity, format string, params ...interface{}) {
	s.ReceiveFromChild(sev, format, params...)
}

func (s *eventSink) ReceiveFromChild(sev xlog.Severity, format string, params ...interface{}) {
	if sev > s.MinSeverity {
		return
	}

	var eid uint32
	eid = 1

	msg := fmt.Sprintf(format, params...)

	if sev <= xlog.SevError {
		s.Log.Error(eid, msg)
	} else if sev <= xlog.SevWarn {
		s.Log.Warning(eid, msg)
	} else {
		s.Log.Info(eid, msg)
	}

	// ignore errors
}

var esink = eventSink{
	MinSeverity: xlog.SevDebug,
}

func openSyslog()  {}
func openJournal() {}
