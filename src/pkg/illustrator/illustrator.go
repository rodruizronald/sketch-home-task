package illustrator

import "github.com/go-playground/validator"

// Max. resolution constraints
// - Width max. 50 characters
// - Height max. 100 characters
const (
	MaxCanvasWidth  int = 50
	MaxCanvasHeight int = 100
)

type Canvas struct {
	Name     string    `json:"name"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Drawings []Drawing `json:"drawings" validate:"dive"`
}

type Drawing struct {
	Coordinates []int `json:"coordinates"`
	Width       int   `json:"width"`
	Height      int   `json:"height"`
	Fill        *rune `json:"fill" validate:"omitempty,gte=32,lte=126"`
	Outline     *rune `json:"outline" validate:"omitempty,gte=32,lte=126"`
}

func (c *Canvas) GetStringCanvas(emptyFiller rune, v *validator.Validate) (canvas string, err error) {
	if v != nil {
		if err = v.Struct(c); err != nil {
			return
		}
	}

	canvasRunes := make([][]rune, c.Height)
	for i := range canvasRunes {
		canvasRunes[i] = make([]rune, c.Width)
		for j := range canvasRunes[i] {
			canvasRunes[i][j] = emptyFiller
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

				canvasRunes[i][j] = char
			}
		}
	}

	for i := range canvasRunes {
		if i > 0 {
			canvas += "\n"
		}
		canvas += string(canvasRunes[i])
	}

	return
}
