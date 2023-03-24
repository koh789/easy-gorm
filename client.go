package egorm

import (
	"errors"
	"reflect"

	"gorm.io/gorm"
)

type CRUD[T any, ID comparable, S ~[]T] interface {
	FindAll() (S, error)

	FindByID(pk ID) (*T, error)

	FindByIDs(pks []ID) (S, error)

	Save(r *T) (*T, error)

	SaveAll(rs S) (S, error)
}

type CRUDClient[T any, ID comparable, S ~[]T] struct {
	db *gorm.DB
}

func NewCRUDClient[T any, ID comparable, S ~[]T](db *gorm.DB) *CRUDClient[T, ID, S] {
	return &CRUDClient[T, ID, S]{db: db}
}

func (c *CRUDClient[T, ID, S]) DB() *gorm.DB {
	return c.db
}

func (c *CRUDClient[T, ID, S]) FindAll() (S, error) {
	var results S
	err := c.db.Find(&results).Error
	return results, err
}

func (c *CRUDClient[T, ID, S]) FindByID(id ID) (*T, error) {
	if ContainZeroValues(id) {
		return nil, nil
	}
	var result T
	err := c.db.First(&result, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, err
}

func (c *CRUDClient[T, ID, S]) FindByIDs(ids []ID) (S, error) {
	var notEmptyIDs []ID
	for _, id := range ids {
		if !ContainZeroValues(id) {
			notEmptyIDs = append(notEmptyIDs, id)
		}
	}
	if len(notEmptyIDs) == 0 {
		return nil, nil
	}
	t := reflect.TypeOf(notEmptyIDs[0])
	var results S
	var err error
	db := c.db
	// Separate processing for embedded PK and primitive types
	if t.Kind() == reflect.Struct {
		for _, id := range notEmptyIDs {
			db = db.Or(id)
		}
		err = db.Find(&results).Error
	} else {
		err = db.Find(&results, notEmptyIDs).Error
	}
	return results, err
}

func (c *CRUDClient[T, ID, S]) Save(entity *T) (*T, error) {
	if entity == nil {
		return nil, nil
	}
	entities := []T{*entity}
	err := c.db.Save(entities).Error
	if err != nil {
		return nil, err
	}
	return &entities[0], nil
}

func (c *CRUDClient[T, ID, S]) SaveAll(entities S) (S, error) {
	if len(entities) == 0 {
		return nil, nil
	}
	err := c.db.Save(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
