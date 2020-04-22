package main

import (
	"paintbot-client/models"
	"paintbot-client/utilities/maputility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test2AllMovesIsEmptyAndOpen(t *testing.T) {

	utility := maputility.MapUtility{
		Map: models.Map{
			Width:  6,
			Height: 6,
			CharacterInfos: []models.CharacterInfo{{
				Position:        14,
				CarryingPowerUp: false,
				ID:              "myId",
			}},
		},
		CurrentPlayerID: "myId",
	}

	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	moves := calculateBestDirection(coordinates, utility, 2, models.Stay)
	assert.Equal(t, 20, moves[0].Points)
	assert.Equal(t, 20, moves[1].Points)
	assert.Equal(t, 20, moves[2].Points)
	assert.Equal(t, 20, moves[3].Points)
}

func Test2AllMovesIsMine(t *testing.T) {

	utility := maputility.MapUtility{
		Map: models.Map{
			Width:  6,
			Height: 6,
			CharacterInfos: []models.CharacterInfo{{
				Position:         14,
				CarryingPowerUp:  false,
				ID:               "myId",
				ColouredPosition: []int{0, 1, 2, 3, 4, 6, 7, 8, 9, 10, 12, 13, 15, 16, 18, 19, 20, 21, 22, 24, 25, 26, 27, 28},
			}},
		},
		CurrentPlayerID: "myId",
	}

	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	moves := calculateBestDirection(coordinates, utility, 2, models.Left)
	assert.Equal(t, 10, moves[0].Points)
	assert.Equal(t, 10, moves[1].Points)
	assert.Equal(t, 10, moves[2].Points)
	assert.Equal(t, 0, moves[3].Points)
	assert.Equal(t, models.Action("RIGHT"), moves[3].Action)
}
