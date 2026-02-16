package apiutil

import (
	"encoding/json"

	"github.com/bernardinorafael/gogem/pkg/uid"
)

type Expandable[T any] struct {
	id       string
	data     *T
	expanded bool
}

func NewExpandableField[T any](id string, data *T) *Expandable[T] {
	if !uid.IsValid(id) {
		// indicates a programming error, not user input error
		panic("invalid resource id")
	}
	if data == nil {
		return &Expandable[T]{
			id:       id,
			data:     nil,
			expanded: false,
		}
	}
	return &Expandable[T]{
		id:       id,
		data:     data,
		expanded: true,
	}
}

func (e *Expandable[T]) MarshalJSON() ([]byte, error) {
	if e.expanded {
		return json.Marshal(e.data)
	}
	return json.Marshal(e.id)
}
