package illustrator

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type CanvasStorage interface {
	Close()
	FindByName(ctx context.Context, name string) (canvas *CanvasModel, err error)
	Create(ctx context.Context, canvas *CanvasModel) (err error)
	Update(ctx context.Context, canvas *CanvasModel) (err error)
	Delete(ctx context.Context, name string) (err error)
}

// TODO: REMOVE VALIDATION FROM TAGS
// validate name length lte = 15

type CanvasModel struct {
	Name     string       `json:"name"`
	Width    int          `json:"width"`
	Height   int          `json:"height"`
	Drawings DrawingSlice `json:"drawings" validate:"dive"`
}

type DrawingModel struct {
	Coordinates []int `json:"coordinates"`
	Width       int   `json:"width"`
	Height      int   `json:"height"`
	Fill        *rune `json:"fill" validate:"omitempty,gte=32,lte=126"`
	Outline     *rune `json:"outline" validate:"omitempty,gte=32,lte=126"`
}

type DrawingSlice []DrawingModel

// DrawingSlice Scanner/Valuer database/sql interface for serialization in databases

func (o *DrawingSlice) Scan(value interface{}) (err error) {
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("value type '%T' is not supported", value)
	}
	return json.Unmarshal(data, o)
}

func (o DrawingSlice) Value() (ret driver.Value, err error) {
	return json.Marshal(o)
}
