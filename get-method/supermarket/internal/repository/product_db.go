package repository

import (
	"database/sql"
	"errors"
	"supermarket/internal"

	"github.com/go-sql-driver/mysql"
)

func NewProductDB(db *sql.DB) *ProductDB {
	return &ProductDB{db: db}
}

type ProductDB struct {
	db *sql.DB
}

// Queries
const (
	GetProductById = "SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE id = ?"
	CreateProduct  = "INSERT INTO products (name, quantity, code_value, is_published, expiration, price) VALUES (?, ?, ?, ?, ?, ?)"
	UpdateProduct  = "UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ? WHERE id = ?"
	DeleteProduct  = "DELETE FROM products WHERE id = ?"
)

func (pdb *ProductDB) GetById(id int) (internal.Product, error) {
	row := pdb.db.QueryRow(GetProductById, id)
	if err := row.Err(); err != nil {
		return internal.Product{}, err
	}

	var product internal.Product
	if err := row.Scan(&product.Id, &product.Name, &product.Quantity, &product.Code, &product.IsPublished, &product.Expiration, &product.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Product{}, internal.NewProductNotFoundError()
		}
		return internal.Product{}, err
	}

	return product, nil
}

func (pdb *ProductDB) Create(product *internal.Product) error {
	result, err := pdb.db.Exec(
		CreateProduct,
		product.Name,
		product.Quantity,
		product.Code,
		product.IsPublished,
		product.Expiration,
		product.Price,
	)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			switch mysqlError.Number {
			case 1062:
				return internal.NewProductAlreadyExistsError()
			}
			return err
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.Id = int(id)

	return nil
}

func (pdb *ProductDB) Update(product *internal.Product) error {
	row, err := pdb.db.Exec(
		UpdateProduct,
		product.Name,
		product.Quantity,
		product.Code,
		product.IsPublished,
		product.Expiration,
		product.Price,
		product.Id,
	)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			switch mysqlError.Number {
			case 1062:
				return internal.NewProductAlreadyExistsError()
			}
			return err
		}
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return internal.NewProductNotFoundError()
	}

	return nil
}

func (pdb *ProductDB) Delete(id int) error {
	row, err := pdb.db.Exec(DeleteProduct, id)
	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return internal.NewProductNotFoundError()
	}

	return nil
}
