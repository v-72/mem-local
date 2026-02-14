package store

type MemLocalStore struct {
	// Declaration: map[keyType]valueType
	Data map[string]string 
}

func (s *MemLocalStore) Set(k string, v string) bool {
	// Initialize the map if it's nil
	if s.Data == nil {
		s.Data = make(map[string]string)
	}
	s.Data[k] = v
	return true
}

func (s *MemLocalStore) Get(k string) (string, bool) {
	// Initialize the map if it's nil
	if s.Data == nil {
		return "", false
	}
	val, ok:= s.Data[k]
	if !ok{
		return val, false
	} 
	return val, true
}
 