package api

type ExecType uint8

const (
	ExecTypeFile ExecType = iota + 1
	ExecTypeShell
	ExecTypeHttp
)
