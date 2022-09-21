package illustrator

import (
	"fmt"

	"github.com/go-playground/validator"
)

// Max. resolution constraints
// - Width max. 50 characters
// - Height max. 100 characters
const (
	MaxCanvasWidth  int = 50
	MaxCanvasHeight int = 100
)

func DrawingModelValidation(sl validator.StructLevel) {
	if drawing, ok := sl.Current().Interface().(DrawingModel); ok {
		var setFields int

		// Validate filler - at least one set
		if drawing.Fill == nil {
			setFields++
		}
		if drawing.Outline == nil {
			setFields++
		}
		if setFields == 2 {
			sl.ReportError(drawing, "Fill/Outline", "Fill/Outline", "at least one field must be set", "")
		}

		// Validate coordinates
		if len(drawing.Coordinates) != 2 {
			sl.ReportError(drawing, "Coordinates", "Coordinates", "only two entries allowed", "")
		}
		if drawing.Coordinates[0] > MaxCanvasHeight {
			tag := fmt.Sprintf("invalid 'i' coordinate value, must be less or equal than %d", MaxCanvasHeight)
			sl.ReportError(drawing, "Coordinates", "Coordinates", tag, "")
		}
		if drawing.Coordinates[1] > MaxCanvasWidth {
			tag := fmt.Sprintf("invalid 'j' coordinate value, must be less or equal than %d", MaxCanvasWidth)
			sl.ReportError(drawing, "Coordinates", "Coordinates", tag, "")
		}

		// Validate drawing dimensions
		if drawing.Width > MaxCanvasWidth {
			tag := fmt.Sprintf("drawing width max. value %d exceeded", MaxCanvasWidth)
			sl.ReportError(drawing, "Width", "Width", tag, "")
		}
		if drawing.Height > MaxCanvasHeight {
			tag := fmt.Sprintf("drawing height max. value %d exceeded", MaxCanvasHeight)
			sl.ReportError(drawing, "Height", "Height", tag, "")
		}
	}
}

func CanvasModelValidation(sl validator.StructLevel) {
	if canvas, ok := sl.Current().Interface().(CanvasModel); ok {
		// Validate canvas dimensions
		if canvas.Width > MaxCanvasWidth {
			tag := fmt.Sprintf("canvas width max. value %d exceeded", MaxCanvasWidth)
			sl.ReportError(canvas, "Width", "Width", tag, "")
		}
		if canvas.Height > MaxCanvasHeight {
			tag := fmt.Sprintf("canvas height max. value %d exceeded", MaxCanvasHeight)
			sl.ReportError(canvas, "Height", "Height", tag, "")
		}
	}
}

func RegisterValidation(v *validator.Validate) {
	v.RegisterStructValidation(DrawingModelValidation, DrawingModel{})
	v.RegisterStructValidation(CanvasModelValidation, CanvasModel{})
}
