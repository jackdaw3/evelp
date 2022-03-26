package cache

type Cacheable interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}
