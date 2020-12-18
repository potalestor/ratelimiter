package operator

// Semaphore is on the buffered channel.
type Semaphore struct {
	bufchan chan struct{}
}

// NewSemaphore returns new instance.
func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		bufchan: make(chan struct{}, size),
	}
}

// Give reads from channel and frees up space on the channel.
func (s *Semaphore) Give() {
	<-s.bufchan
}

// Take writes struct to channel.
func (s *Semaphore) Take() error {
	s.bufchan <- struct{}{}

	return nil
}

// Close channel.
func (s *Semaphore) Close() error {
	close(s.bufchan)

	return nil
}
