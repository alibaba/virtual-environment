package event

type TagAppenderInitFailed struct {
}

func (t TagAppenderInitFailed) Error() string {
	return "failed to create tag appender, caused by unknown configure unmarshal error"
}

func IsTagAppenderInitFailed(err error) bool {
	_, ok := err.(TagAppenderInitFailed)
	return ok
}
