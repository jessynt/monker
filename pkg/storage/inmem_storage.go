package storage

type InmemStorage struct {
	queue *Queue
}

func NewInmemStorage() Storage {
	return &InmemStorage{
		queue: NewQueue(),
	}
}

func (s *InmemStorage) Put(data []byte) error {
	s.queue.Push(data)
	return nil
}

func (s *InmemStorage) Get() ([]byte, error) {
	value := s.queue.Pop()
	if value == nil {
		return nil, ErrStorageIsEmpty
	}

	return value, nil
}

func (s *InmemStorage) Close() error {
	return nil
}
