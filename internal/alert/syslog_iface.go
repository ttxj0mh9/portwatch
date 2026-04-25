package alert

// syslogIface is the subset of *syslog.Writer consumed by SyslogHandler.
// Keeping it internal lets tests inject a fake without importing log/syslog
// in test files.
type syslogIface interface {
	Notice(string) error
	Warning(string) error
	Close() error
}

// Ensure *syslog.Writer satisfies syslogIface at compile time.
// (Checked via the real struct stored in SyslogHandler.writer.)

// SyslogHandler is redeclared here with the interface field so that both
// the production constructor (using *syslog.Writer) and the test helper
// (using fakeSyslogWriter) work without reflection.
//
// NOTE: the struct definition lives in syslog_handler.go; this file only
// documents the interface contract.
