package model

type ExecType uint8

const (
	ExecTypeShell ExecType = iota + 1
	ExecTypeHttp
	ExecTypeFile
)

type NotifyStrategy uint8

const (
	NotifyAfterSuccess NotifyStrategy = iota + 1 // 成功后通知
	NotifyAfterFailed                            // 失败后通知
	NotifyAlways                                 // 总是通知
)

type NotifyType uint8

const (
	NotifyTypeEmail NotifyType = iota + 1
)

type NotifyStatus uint8

const (
	NotifyStatusEnabled NotifyStatus = iota + 1
	NotifyStatusDisabled
)
