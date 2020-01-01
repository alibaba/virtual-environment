package shared

type VirtualEnvChangedError struct {
}

func (v VirtualEnvChangedError) Error() string {
	return "detected virtual environment instance change"
}

func IsVirtualEnvChanged(err error) bool {
	_, ok := err.(VirtualEnvChangedError)
	return ok
}
