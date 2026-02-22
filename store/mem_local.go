package store

type Store interface {
	Init() error
	Set(string, string) bool
	Get(string) (string, bool)
}

type MemLocalStore struct {
	Data map[string]string 
}

func (s *MemLocalStore) Init() error {
	s.Data = make(map[string]string)
	return nil
}


func (s *MemLocalStore) Set(k string, v string) bool {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}
	s.Data[k] = v
	return true
}

func (s *MemLocalStore) Get(k string) (string, bool) {
	if s.Data == nil {
		return "", false
	}
	val, ok:= s.Data[k]
	if !ok{
		return val, false
	} 
	return val, true
}
 