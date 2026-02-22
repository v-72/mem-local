package store

type Store interface {
	Init() error
	Set(string, string) bool
	Get(string) (string, bool)
}

type MemLocalStore struct {
	Data map[string]string
}

type GetResult struct {
	value string
	ok    bool
}

func (s *MemLocalStore) Init() error {
	s.Data = make(map[string]string)
	return nil
}

func (s *MemLocalStore) getData(key string, dataChan chan GetResult) {
	if val, ok := s.Data[key]; ok {
		dataChan <- GetResult{value: val, ok: ok}
	} else {
		dataChan <- GetResult{value: "", ok: false}
	}
}

func (s *MemLocalStore) Set(k string, v string) bool {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}
	s.Data[k] = v
	return true
}

func (s *MemLocalStore) Get(k string) (string, bool) {
	dataChan := make(chan GetResult)
	go s.getData(k, dataChan)
	result := <-dataChan
	return result.value, result.ok
}
