package main

import (
	"os"
	"sort"

	"github.com/dariubs/percent"
	log "github.com/sirupsen/logrus"

	"paintbot-client/basebot"
	"paintbot-client/models"
	"paintbot-client/utilities/maputility"
)

type PossibleAction struct {
	Action models.Action
	Points int
}

func main() {
	basebot.Start("Simple Go Bot", models.Training, desiredGameSettings, calculateMove)
}

//var moves = []models.Action{models.Left, models.Down, models.Right, models.Up} //, models.Stay}
var lastDir = models.Stay
var noOfRecursions = 5

// Implement your paintbot here
func calculateMove(updateEvent models.MapUpdateEvent) models.Action {

	utility := maputility.MapUtility{Map: updateEvent.Map, CurrentPlayerID: *updateEvent.ReceivingPlayerID}
	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	if me.StunnedForGameTicks > 0 {
		return models.Stay
	}

	if me.CarryingPowerUp && shouldExplode(coordinates, utility, 4) {
		return models.Explode
	}

	possibleActions := calculateBestDirection(coordinates, utility, noOfRecursions, lastDir)

	// log.Info(fmt.Printf("x: %d y: %d \n %s %d \n %s %d \n %s %d \n %s %d \n",
	// 	coordinates.X, coordinates.Y,
	// 	possibleActions[0].Action, possibleActions[0].Points,
	// 	possibleActions[1].Action, possibleActions[1].Points,
	// 	possibleActions[2].Action, possibleActions[2].Points,
	// 	possibleActions[3].Action, possibleActions[3].Points))

	// log.Info(possibleActions[0].Action)
	lastDir = possibleActions[0].Action
	return possibleActions[0].Action
}

func calculateBestDirection(coordinats models.Coordinates, utility maputility.MapUtility, recursion int, last models.Action) []PossibleAction {

	moves := []PossibleAction{
		{Action: models.Left},
		{Action: models.Down},
		{Action: models.Right},
		{Action: models.Up},
	}

	for i, action := range moves {
		newCoordinates := utility.TranslateCoordinateByAction(action.Action, coordinats)

		points := 0

		if action.Action == getReverseMove(last) {
			points += -10
		}

		switch tile := utility.GetTileAt(newCoordinates); tile {
		case "OBSTACLE":

			points += -100
		case "POWERUP":
			if utility.GetMyCharacterInfo().CarryingPowerUp {
				points += 10
			} else {
				points += 100
			}
		case "PLAYER":
			points += -50
		case "OPEN":
			if !utility.CanIMoveInDirection(getReverseMove(action.Action)) {
				points += 10
			} else {
				points += 5
			}

			points += calculatePlayerOnPosition(newCoordinates, utility)
		}

		//points = points * recursion
		if recursion > 1 && utility.CanIMoveInDirection(action.Action) {
			possibleActions := calculateBestDirection(newCoordinates, utility, recursion-1, action.Action)
			points += possibleActions[0].Points
		}
		moves[i].Points = points
	}

	// Sort best moves first
	sort.SliceStable(moves, func(i, j int) bool {
		return moves[i].Points > moves[j].Points
	})

	return moves
}

func calculatePlayerOnPosition(coordinates models.Coordinates, utility maputility.MapUtility) int {
	switch playerName := getPlayerOnPosition(coordinates, utility); playerName {
	case "":
		return 0
	case utility.GetMyCharacterInfo().ID:
		return -5
	default:
		return 10
	}
}

func getPlayerOnPosition(coordinates models.Coordinates, utility maputility.MapUtility) string {
	position := utility.ConvertCoordinatesToPosition(coordinates)
	for _, character := range utility.Map.CharacterInfos {

		for _, pos := range character.ColouredPosition {
			if pos == position {
				return character.ID
			}
		}
	}
	return ""
}

func getReverseMove(action models.Action) models.Action {
	switch action {
	case models.Action("RIGHT"):
		return models.Action("LEFT")
	case models.Action("LEFT"):
		return models.Action("RIGHT")
	case models.Action("UP"):
		return models.Action("DOWN")
	case models.Action("DOWN"):
		return models.Action("UP")
	}

	return action
}

func shouldExplode(myCoordinate models.Coordinates, utility maputility.MapUtility, explosionRange int) bool {

	column := explosionRange*2 + 1
	noOfTiles := column * column

	startCoordinate := models.Coordinates{X: myCoordinate.X - explosionRange, Y: myCoordinate.Y - explosionRange}

	powerups := 0
	open := 0
	otherPlayers := 0
	myTile := 0

	for y := 0; y <= explosionRange*2; y++ {
		for x := 0; x <= explosionRange*2; x++ {

			coordinate := models.Coordinates{X: startCoordinate.X + x, Y: startCoordinate.Y + y}

			if coordinate.X == myCoordinate.X && coordinate.Y == myCoordinate.Y {
				continue
			}

			switch tile := utility.GetTileAt(coordinate); tile {
			case "OBSTACLE":
				continue
			case "POWERUP":
				powerups += 1
			case "PLAYER":
				// sure, explode and take him down!
				return true
			case "OPEN":
				if getPlayerOnPosition(coordinate, utility) == utility.GetMyCharacterInfo().ID {
					myTile += 1
				} else if getPlayerOnPosition(coordinate, utility) == "" {
					open += 1
				} else {
					otherPlayers += 1
				}
			}

		}
	}

	if percent.PercentOf(open+otherPlayers, noOfTiles-1) > 30 && powerups >= 2 {
		return true
	}

	if percent.PercentOf(otherPlayers, noOfTiles-1) > 25 {
		return true
	}

	if percent.PercentOf(otherPlayers+open, noOfTiles-1) > 70 {
		return true
	}

	return false
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		ForceQuote:             true,
		FullTimestamp:          true,
		TimestampFormat:        "15:04:05.999",
		DisableLevelTruncation: true,
		PadLevelText:           true,
		QuoteEmptyFields:       true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

// desired game settings can be changed to nil to get default settings
var desiredGameSettings = &models.GameSettings{
	MaxNOOFPlayers:                 5,
	TimeInMSPerTick:                250,
	ObstaclesEnabled:               true,
	PowerUpsEnabled:                true,
	AddPowerUpLikelihood:           38,
	RemovePowerUpLikelihood:        5,
	TrainingGame:                   true,
	PointsPerTileOwned:             1,
	PointsPerCausedStun:            5,
	NOOFTicksInvulnerableAfterStun: 3,
	NOOFTicksStunned:               10,
	StartObstacles:                 40,
	StartPowerUps:                  41,
	GameDurationInSeconds:          15,
	ExplosionRange:                 4,
	PointsPerTick:                  false,
}
