package shutdowner

type shutdownOption struct {
	finalFunc finalFunc
	priority  int
}

type finalFunc func() error
