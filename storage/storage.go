package storage

type Storage interface {
	Put(data []byte) (string, error)
	Delete(id string) error
	Get(id string) ([]byte, error)
}
