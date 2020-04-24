package main

import (
	"errors"
	"math/rand"
	"os"
	"sort"
	"time"

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
	basebot.Start("Thomat-Kickup", models.Tournament, desiredGameSettings, calculateMove)
}

var lastDir = models.Stay
var noOfRecursions = 7
var explosionRange = 4

// Implement your paintbot here
func calculateMove(updateEvent models.MapUpdateEvent) models.Action {

	utility := maputility.MapUtility{Map: updateEvent.Map, CurrentPlayerID: *updateEvent.ReceivingPlayerID}
	me := utility.GetMyCharacterInfo()
	coordinates := utility.ConvertPositionToCoordinates(me.Position)

	if me.StunnedForGameTicks > 0 {
		return models.Stay
	}

	if me.CarryingPowerUp && shouldExplode(coordinates, utility, explosionRange) {
		return models.Explode
	}

	possibleActions := calculateBestDirection(coordinates, utility, noOfRecursions, lastDir)

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

	moves = shuffle(moves)

	for i, action := range moves {
		newCoordinates := utility.TranslateCoordinateByAction(action.Action, coordinats)

		moves[i].Points = calculateBestDirectionSync(action, newCoordinates, utility, recursion, last)
	}

	// Sort best moves first
	sort.SliceStable(moves, func(i, j int) bool {
		return moves[i].Points > moves[j].Points
	})

	return moves
}

func calculateBestDirectionSync(action PossibleAction, coordinats models.Coordinates, utility maputility.MapUtility, recursion int, last models.Action) int {

	points := 0

	if action.Action == getReverseMove(last) {
		points += -10
	}

	switch tile := utility.GetTileAt(coordinats); tile {
	case "OBSTACLE":
		points += -150
	case "POWERUP":
		if utility.GetMyCharacterInfo().CarryingPowerUp {
			points += 10
		} else {
			points += 500
		}
	case "PLAYER":
		player, err := getPlayerOnPosition(coordinats, utility)
		if err == nil {
			if player.ID == utility.GetMyCharacterInfo().ID {
				points += -10
			} else {
				if player.CarryingPowerUp {
					points += -100
				} else {
					points += -5
				}
			}
		}

	case "OPEN":
		points += 10

		points += calculatePlayerColorOnPosition(coordinats, utility)
	}

	if recursion > 1 && utility.CanIMoveInDirection(action.Action) {
		possibleActions := calculateBestDirection(coordinats, utility, recursion-1, action.Action)
		points += possibleActions[0].Points
	}

	return points
}

func calculatePlayerColorOnPosition(coordinates models.Coordinates, utility maputility.MapUtility) int {
	switch playerName := getPlayerColorIDOnPosition(coordinates, utility); playerName {
	case "":
		return 0
	case utility.GetMyCharacterInfo().ID:
		return -10
	default:
		return 10
	}
}

func getPlayerColorIDOnPosition(coordinates models.Coordinates, utility maputility.MapUtility) string {
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

func getPlayerOnPosition(coordinates models.Coordinates, utility maputility.MapUtility) (models.CharacterInfo, error) {

	position := utility.ConvertCoordinatesToPosition(coordinates)
	for _, character := range utility.Map.CharacterInfos {

		if position == character.Position {
			return character, nil
		}
	}

	return models.CharacterInfo{}, errors.New("no player found")
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
				powerups++
			case "PLAYER":
				// sure, explode and take him down!
				return true
			case "OPEN":
				if getPlayerColorIDOnPosition(coordinate, utility) == utility.GetMyCharacterInfo().ID {
					myTile++
				} else if getPlayerColorIDOnPosition(coordinate, utility) == "" {
					open++
				} else {
					otherPlayers++
				}
			}

		}
	}

	if percent.PercentOf(open+otherPlayers, noOfTiles-1) > 30 && powerups >= 1 {
		return true
	}

	if percent.PercentOf(otherPlayers, noOfTiles-1) > 45 {
		return true
	}

	if percent.PercentOf(otherPlayers+open, noOfTiles-1) > 70 {
		return true
	}

	return false
}

func shuffle(l []PossibleAction) []PossibleAction {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for i := range l {
		n := r.Intn(len(l) - 1)
		l[i], l[n] = l[n], l[i]
	}
	return l
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
	MaxNOOFPlayers:                 3,
	TimeInMSPerTick:                250,
	ObstaclesEnabled:               true,
	PowerUpsEnabled:                true,
	AddPowerUpLikelihood:           15,
	RemovePowerUpLikelihood:        5,
	TrainingGame:                   true,
	PointsPerTileOwned:             1,
	PointsPerCausedStun:            5,
	NOOFTicksInvulnerableAfterStun: 3,
	NOOFTicksStunned:               10,
	StartObstacles:                 50,
	StartPowerUps:                  10,
	GameDurationInSeconds:          15,
	ExplosionRange:                 4,
	PointsPerTick:                  false,
}
