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

// Entity
// If you want T to have a behavior,
// define it in Entity and implement it in a struct that is mapped to a table.
type Entity[T any, ID comparable, S ~[]T] interface {
	*T
}

type CRUDClient[T any, ID comparable, S ~[]T, E Entity[T, ID, S]] struct {
	reader *gorm.DB
	writer *gorm.DB
}

func NewCRUDClient[T any, ID comparable, S ~[]T, E Entity[T, ID, S]](reader, writer *gorm.DB) *CRUDClient[T, ID, S, E] {
	return &CRUDClient[T, ID, S, E]{reader: reader, writer: writer}
}

func (c *CRUDClient[T, ID, S, E]) Reader() *gorm.DB {
	return c.reader
}

func (c *CRUDClient[T, ID, S, E]) Writer() *gorm.DB {
	return c.writer
}

func (c *CRUDClient[T, ID, S, E]) FindAll() (S, error) {
	var results S
	err := c.reader.Find(&results).Error
	return results, err
}

func (c *CRUDClient[T, ID, S, E]) FindByID(id ID) (*T, error) {
	if ContainZeroValues(id) {
		return nil, nil
	}
	var result T
	err := c.reader.First(&result, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, err
}

func (c *CRUDClient[T, ID, S, E]) FindByIDs(ids []ID) (S, error) {
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
	reader := c.reader
	// Separate processing for embedded PK and primitive types
	if t.Kind() == reflect.Struct {
		for _, id := range notEmptyIDs {
			reader = reader.Or(id)
		}
		err = reader.Find(&results).Error
	} else {
		err = reader.Find(&results, notEmptyIDs).Error
	}
	return results, err
}

func (c *CRUDClient[T, ID, S, E]) Save(entity *T) (*T, error) {
	if entity == nil {
		return nil, nil
	}
	entities := []T{*entity}
	err := c.writer.Save(entities).Error
	if err != nil {
		return nil, err
	}
	return &entities[0], nil
}

func (c *CRUDClient[T, ID, S, E]) SaveAll(entities S) (S, error) {
	if len(entities) == 0 {
		return nil, nil
	}
	err := c.writer.Save(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
