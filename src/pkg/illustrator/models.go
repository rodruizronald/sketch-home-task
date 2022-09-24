package illustrator

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type CanvasStorage interface {
	Close() (err error)
	FindByName(ctx context.Context, name string) (canvas *CanvasModel, err error)
	Create(ctx context.Context, canvas *CanvasModel) (res sql.Result, err error)
	Update(ctx context.Context, canvas *CanvasModel) (res sql.Result, err error)
	Delete(ctx context.Context, name string) (res sql.Result, err error)
}

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
	Fill        *rune `json:"fill"`
	Outline     *rune `json:"outline"`
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
