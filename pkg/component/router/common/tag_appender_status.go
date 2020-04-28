package common

type TagAppenderStatus int

const (
	NotExist TagAppenderStatus = iota
	UpToDate
	Outdated
	Unknown
)

func IsTagAppenderExist(status TagAppenderStatus) bool {
	return status == UpToDate || status == Outdated
}

func IsTagAppenderNeedUpdate(status TagAppenderStatus) bool {
	return status == NotExist || status == Outdated
}
