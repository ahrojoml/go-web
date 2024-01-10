package api

type InvalidProductError struct {
	Field string
}

func (e InvalidProductError) Error() string {
	return "invalid product"
}

func NewInvalidProductError(field string) InvalidProductError {
	return InvalidProductError{Field: field}
}
