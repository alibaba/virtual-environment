package shared

type VirtualEnvChangeDetected struct {
}

func (v VirtualEnvChangeDetected) Error() string {
	return "detected virtual environment instance change"
}

func IsVirtualEnvChanged(err error) bool {
	_, ok := err.(VirtualEnvChangeDetected)
	return ok
}
