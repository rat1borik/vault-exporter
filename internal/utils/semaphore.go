package utils

type sema struct {
	C chan struct{}
}

func NewSemaphore(len int) *sema {
	return &sema{
		C: make(chan struct{}, len),
	}
}

func (s *sema) Acquire() {
	s.C <- struct{}{}
}

func (s *sema) Release() {
	<-s.C
}
