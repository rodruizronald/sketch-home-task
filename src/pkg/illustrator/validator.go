package illustrator

import (
	"fmt"

	"github.com/go-playground/validator"
)

const (
	// Canvas max. resolution constraints
	// - Width max. 50 characters
	// - Height max. 100 characters
	CanvasMaxWidth  int = 50
	CanvasMaxHeight int = 100
	// Canvas max. name length
	CanvasMaxNameSize int = 15
	// Drawing ASCII characters lower/higher limit
	DrawingCharLowerLimit  rune = 32
	DrawingCharHigherLimit rune = 126
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

		// Validate drawing character range
		if drawing.Fill != nil &&
			(*drawing.Fill < DrawingCharLowerLimit || *drawing.Fill > DrawingCharHigherLimit) {
			tag := fmt.Sprintf("invalid character, must be between %d - %d", DrawingCharLowerLimit, DrawingCharHigherLimit)
			sl.ReportError(drawing, "Fill", "Fill", tag, "")
		}
		if drawing.Outline != nil &&
			(*drawing.Outline < DrawingCharLowerLimit || *drawing.Outline > DrawingCharHigherLimit) {
			tag := fmt.Sprintf("invalid character, must be between %d - %d", DrawingCharLowerLimit, DrawingCharHigherLimit)
			sl.ReportError(drawing, "Outline", "Outline", tag, "")
		}

		// Validate number of coordinates
		if len(drawing.Coordinates) != 2 {
			sl.ReportError(drawing, "Coordinates", "Coordinates", "only two entries allowed", "")
		}

		// Validate coordinates value - must be within canvas size range
		if drawing.Coordinates[0] > CanvasMaxHeight {
			tag := fmt.Sprintf("invalid 'i' coordinate value, must be less or equal than %d", CanvasMaxHeight)
			sl.ReportError(drawing, "Coordinates", "Coordinates", tag, "")
		}
		if drawing.Coordinates[1] > CanvasMaxWidth {
			tag := fmt.Sprintf("invalid 'j' coordinate value, must be less or equal than %d", CanvasMaxWidth)
			sl.ReportError(drawing, "Coordinates", "Coordinates", tag, "")
		}

		// Validate drawing dimensions - must be within canvas size range
		if drawing.Width > CanvasMaxWidth {
			tag := fmt.Sprintf("drawing width max. value %d exceeded", CanvasMaxWidth)
			sl.ReportError(drawing, "Width", "Width", tag, "")
		}
		if drawing.Height > CanvasMaxHeight {
			tag := fmt.Sprintf("drawing height max. value %d exceeded", CanvasMaxHeight)
			sl.ReportError(drawing, "Height", "Height", tag, "")
		}
	}
}

func CanvasModelValidation(sl validator.StructLevel) {
	if canvas, ok := sl.Current().Interface().(CanvasModel); ok {
		// Validate canvas name max. length
		if len(canvas.Name) == 0 && len(canvas.Name) > CanvasMaxNameSize {
			tag := fmt.Sprintf("canvas name max. size %d exceeded", CanvasMaxNameSize)
			sl.ReportError(canvas, "Name", "Name", tag, "")
		}

		// Validate canvas dimensions
		if canvas.Width > CanvasMaxWidth {
			tag := fmt.Sprintf("canvas width max. value %d exceeded", CanvasMaxWidth)
			sl.ReportError(canvas, "Width", "Width", tag, "")
		}
		if canvas.Height > CanvasMaxHeight {
			tag := fmt.Sprintf("canvas height max. value %d exceeded", CanvasMaxHeight)
			sl.ReportError(canvas, "Height", "Height", tag, "")
		}
	}
}

func RegisterValidation(v *validator.Validate) {
	v.RegisterStructValidation(DrawingModelValidation, DrawingModel{})
	v.RegisterStructValidation(CanvasModelValidation, CanvasModel{})
}
