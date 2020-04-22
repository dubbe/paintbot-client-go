package main

import (
	"paintbot-client/models"
	"paintbot-client/utilities/maputility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllEmptyShouldExplode(t *testing.T) {

	utility := maputility.MapUtility{
		Map: models.Map{
			Width:  6,
			Height: 6,
			CharacterInfos: []models.CharacterInfo{{
				Position:        14,
				CarryingPowerUp: true,
				ID:              "myId",
			}},
		},
		CurrentPlayerID: "myId",
	}

	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	shouldExplode := shouldExplode(coordinates, utility, 2)
	assert.Equal(t, true, shouldExplode)

}
