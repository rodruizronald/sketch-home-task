package illustrator

import (
	"github.com/go-playground/validator"
)

func (c *CanvasModel) GetString(emptyFiller rune, newLine string, validator *validator.Validate) (str string, err error) {
	if validator != nil {
		if err = validator.Struct(c); err != nil {
			return
		}
	}

	runes := make([][]rune, c.Height)
	for i := range runes {
		runes[i] = make([]rune, c.Width)
		for j := range runes[i] {
			runes[i][j] = emptyFiller
		}
	}

	char := emptyFiller
	for _, drawing := range c.Drawings {
		fillChar := emptyFiller
		if drawing.Fill != nil {
			fillChar = *drawing.Fill
		}
		outlineChar := emptyFiller
		if drawing.Outline != nil {
			outlineChar = *drawing.Outline
		}

		iStartPoint := drawing.Coordinates[0]
		jStartPoint := drawing.Coordinates[1]

		iEndPoint := iStartPoint + drawing.Height - 1
		jEndPoint := jStartPoint + drawing.Width - 1

		for i := iStartPoint; i <= iEndPoint; i++ {
			for j := jStartPoint; j <= jEndPoint; j++ {
				char = fillChar
				if i == iStartPoint ||
					j == jStartPoint ||
					i == iEndPoint ||
					j == jEndPoint {
					char = outlineChar
				}

				// Skip sections out of range
				if i >= c.Height || j >= c.Width {
					continue
				}

				runes[i][j] = char
			}
		}
	}

	for i := range runes {
		if i > 0 {
			str += newLine
		}
		str += string(runes[i])
	}

	return
}
