package internal

type StorageRepository interface {
	GetAll() (map[int]Product, error)
	Save(map[int]Product) error
}
