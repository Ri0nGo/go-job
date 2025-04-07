package model

type ExecType uint8

const (
	ExecTypeShell ExecType = iota + 1
	ExecTypeHttp
	ExecTypeFile
)
