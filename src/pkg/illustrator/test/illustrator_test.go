package illustrator_test

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/sketch-home-task/src/pkg/illustrator"
	"github.com/stretchr/testify/assert"
)

type testEntry struct {
	name           string
	canvas         illustrator.CanvasModel
	expectedCanvas string
	validEntry     bool
}

var (
	validCanvas string = `@-XXX-OOOO--OOO
--XXX-O$$O--O--
@-----O$$O--OOO
@---XXOOOO-----
----X***X------
@@--XXXXX----$$
@@-----------$$`
)

func TestIllustrator(t *testing.T) {
	asteriskRune := '*'
	atRune := '@'
	bigORune := 'O'
	dollerRune := '$'
	bigXRune := 'X'
	invalidRune := rune(138)

	testTable := []testEntry{
		// Valid test cases
		{
			name: "Test illustrator with valid canvas",
			canvas: illustrator.CanvasModel{
				Width:  15,
				Height: 7,
				Drawings: []illustrator.DrawingModel{
					{
						Coordinates: []int{0, 0},
						Width:       0,
						Height:      0,
						Fill:        &asteriskRune,
						Outline:     &atRune,
					},
					{
						Coordinates: []int{0, 0},
						Width:       1,
						Height:      1,
						Fill:        &asteriskRune,
						Outline:     &atRune,
					},
					{
						Coordinates: []int{2, 0},
						Width:       1,
						Height:      2,
						Fill:        &asteriskRune,
						Outline:     &atRune,
					},
					{
						Coordinates: []int{5, 0},
						Width:       2,
						Height:      2,
						Fill:        &asteriskRune,
						Outline:     &atRune,
					},
					{
						Coordinates: []int{0, 2},
						Width:       3,
						Height:      2,
						Fill:        &asteriskRune,
						Outline:     &bigXRune,
					},
					{
						Coordinates: []int{3, 4},
						Width:       5,
						Height:      3,
						Fill:        &asteriskRune,
						Outline:     &bigXRune,
					},
					{
						Coordinates: []int{0, 6},
						Width:       4,
						Height:      4,
						Fill:        &dollerRune,
						Outline:     &bigORune,
					},
					{
						Coordinates: []int{0, 12},
						Width:       4,
						Height:      3,
						Fill:        nil,
						Outline:     &bigORune,
					},
					{
						Coordinates: []int{4, 12},
						Width:       4,
						Height:      4,
						Fill:        &dollerRune,
						Outline:     nil,
					},
				},
			},
			expectedCanvas: validCanvas,
			validEntry:     true,
		},
		// Invalid test cases
		{
			name: "Test illustrator with invalid drawing dimensions",
			canvas: illustrator.CanvasModel{
				Width:  15,
				Height: 7,
				Drawings: []illustrator.DrawingModel{
					{
						Coordinates: []int{0, 0},
						Width:       6,
						Height:      illustrator.CanvasMaxHeight + 1,
						Fill:        &dollerRune,
						Outline:     nil,
					},
				},
			},
			expectedCanvas: "",
			validEntry:     false,
		},
		{
			name: "Test illustrator with invalid canvas dimensions",
			canvas: illustrator.CanvasModel{
				Width:  illustrator.CanvasMaxWidth + 1,
				Height: 7,
				Drawings: []illustrator.DrawingModel{
					{
						Coordinates: []int{0, 0},
						Width:       4,
						Height:      5,
						Fill:        &dollerRune,
						Outline:     nil,
					},
				},
			},
			expectedCanvas: "",
			validEntry:     false,
		},
		{
			name: "Test illustrator with invalid filler",
			canvas: illustrator.CanvasModel{
				Width:  20,
				Height: 7,
				Drawings: []illustrator.DrawingModel{
					{
						Coordinates: []int{0, 0},
						Width:       4,
						Height:      5,
						Fill:        &invalidRune,
						Outline:     &invalidRune,
					},
				},
			},
			expectedCanvas: "",
			validEntry:     false,
		},
		{
			name: "Test illustrator with invalid number of drawing coordinates",
			canvas: illustrator.CanvasModel{
				Width:  20,
				Height: 7,
				Drawings: []illustrator.DrawingModel{
					{
						Coordinates: []int{0, 0, 1},
						Width:       4,
						Height:      5,
						Fill:        &dollerRune,
						Outline:     nil,
					},
				},
			},
			expectedCanvas: "",
			validEntry:     false,
		},
		{
			name: "Test illustrator with invalid range of drawing coordinates",
			canvas: illustrator.CanvasModel{
				Width:  20,
				Height: 7,
				Drawings: []illustrator.DrawingModel{
					{
						Coordinates: []int{0, illustrator.CanvasMaxWidth + 1},
						Width:       4,
						Height:      5,
						Fill:        &dollerRune,
						Outline:     nil,
					},
				},
			},
			expectedCanvas: "",
			validEntry:     false,
		},
	}

	validator := validator.New()
	illustrator.RegisterValidation(validator)

	a := assert.New(t)
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actualCanvas, err := tt.canvas.GetString('-', "\n", validator)
			if tt.validEntry {
				a.NoError(err)
			} else {
				a.Error(err)
			}
			a.Equal(tt.expectedCanvas, actualCanvas)
		})
	}
}
