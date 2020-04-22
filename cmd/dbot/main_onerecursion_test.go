package main

import (
	"paintbot-client/models"
	"paintbot-client/utilities/maputility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllMovesIsEmptyAndOpen(t *testing.T) {

	utility := maputility.MapUtility{
		Map: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{{
				Position:        4,
				CarryingPowerUp: false,
				ID:              "myId",
			}},
		},
		CurrentPlayerID: "myId",
	}

	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	moves := calculateBestDirection(coordinates, utility, 1, models.Stay)
	assert.Equal(t, 10, moves[0].Points)
	assert.Equal(t, 10, moves[1].Points)
	assert.Equal(t, 10, moves[2].Points)
	assert.Equal(t, 10, moves[3].Points)
}

func TestMovesIsObstacleAndPowerup(t *testing.T) {

	utility := maputility.MapUtility{
		Map: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{{
				Position:        4,
				CarryingPowerUp: false,
				ID:              "myId",
			}},
			ObstacleUpPositions: []int{1},
			PowerUpPositions:    []int{3},
		},
		CurrentPlayerID: "myId",
	}

	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	moves := calculateBestDirection(coordinates, utility, 1, models.Stay)
	assert.Equal(t, 100, moves[0].Points)
	assert.Equal(t, 10, moves[1].Points)
	assert.Equal(t, 10, moves[2].Points)
	assert.Equal(t, -100, moves[3].Points)
	assert.Equal(t, models.Action("UP"), moves[3].Action)
	assert.Equal(t, models.Action("LEFT"), moves[0].Action)
}

func TestMovesIsOtherCharacters(t *testing.T) {

	utility := maputility.MapUtility{
		Map: models.Map{
			Width:  3,
			Height: 3,
			CharacterInfos: []models.CharacterInfo{
				{
					Position:         4,
					CarryingPowerUp:  false,
					ID:               "myId",
					ColouredPosition: []int{5},
				},
				{
					Position:         1,
					CarryingPowerUp:  false,
					ID:               "otherId",
					ColouredPosition: []int{3},
				},
			},
		},
		CurrentPlayerID: "myId",
	}

	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	moves := calculateBestDirection(coordinates, utility, 1, models.Stay)
	assert.Equal(t, 20, moves[0].Points)
	assert.Equal(t, 10, moves[1].Points)
	assert.Equal(t, 5, moves[2].Points)
	assert.Equal(t, -50, moves[3].Points)
	assert.Equal(t, models.Action("LEFT"), moves[0].Action)
	assert.Equal(t, models.Action("DOWN"), moves[1].Action)
	assert.Equal(t, models.Action("RIGHT"), moves[2].Action)
	assert.Equal(t, models.Action("UP"), moves[3].Action)
}
