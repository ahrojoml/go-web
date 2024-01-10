package product

type InvalidProductError struct {
	Field string
}

func (e InvalidProductError) Error() string {
	return "invalid product"
}

func NewInvalidProductError(field string) InvalidProductError {
	return InvalidProductError{Field: field}
}

type ProductNotFoundError struct{}

func (e ProductNotFoundError) Error() string {
	return "product not found"
}

func NewProductNotFoundError() ProductNotFoundError {
	return ProductNotFoundError{}
}
