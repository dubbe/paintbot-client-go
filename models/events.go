package models

type GameMessage struct {
	Type string `json:"type"`

	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type GameSettings struct {
	MaxNOOFPlayers                 int  `json:"maxNoofPlayers"`
	TimeInMSPerTick                int  `json:"timeInMsPerTick"`
	ObstaclesEnabled               bool `json:"obstaclesEnabled"`
	PowerUpsEnabled                bool `json:"powerUpsEnabled"`
	AddPowerUpLikelihood           int  `json:"addPowerUpLikelihood"`
	RemovePowerUpLikelihood        int  `json:"removePowerUpLikelihood"`
	TrainingGame                   bool `json:"trainingGame"`
	PointsPerTileOwned             int  `json:"pointsPerTileOwned"`
	PointsPerCausedStun            int  `json:"pointsPerCausedStun"`
	NOOFTicksInvulnerableAfterStun int  `json:"noOfTicksInvulnerableAfterStun"`
	NOOFTicksStunned               int  `json:"noOfTicksStunned"`
	StartObstacles                 int  `json:"startObstacles"`
	StartPowerUps                  int  `json:"startPowerUps"`
	GameDurationInSeconds          int  `json:"gameDurationInSeconds"`
	ExplosionRange                 int  `json:"explosionRange"`
	PointsPerTick                  bool `json:"pointsPerTick"`
}

type RegisterPlayerEvent struct {
	Type              string        `json:"type"`
	PlayerName        string        `json:"playerName"`
	GameSettings      *GameSettings `json:"gameSettings"`
	ReceivingPlayerID *string       `json:"receivingPlayerId"`
	Timestamp         int           `json:"timestamp"`
}

type PlayerRegisteredEvent struct {
	Type              string       `json:"type"`
	GameID            string       `json:"gameId"`
	PlayerName        string       `json:"name"`
	GameSettings      GameSettings `json:"gameSettings"`
	GameMode          string       `json:"gameMode"`
	ReceivingPlayerID *string      `json:"receivingPlayerId"`
	Timestamp         int          `json:"timestamp"`
}

type StartGameEvent struct {
	Type              string  `json:"type"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type ClientInfoMSG struct {
	Type                   string  `json:"type"`
	Language               string  `json:"language"`
	LanguageVersion        string  `json:"languageVersion"`
	OperatingSystem        string  `json:"operatingSystem"`
	OperatingSystemVersion string  `json:"operatingSystemVersion"`
	ClientVersion          string  `json:"clientVersion"`
	ReceivingPlayerID      *string `json:"receivingPlayerId"`
	Timestamp              int     `json:"timestamp"`
}

type GameLinkEvent struct {
	Type              string  `json:"type"`
	GameID            string  `json:"gameId"`
	URL               string  `json:"url"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type GameStartingEvent struct {
	Type              string       `json:"type"`
	GameID            string       `json:"gameId"`
	NOOFPlayers       int          `json:"noOfPlayers"`
	Width             int          `json:"width"`
	Height            int          `json:"height"`
	GameSettings      GameSettings `json:"gameSettings"`
	ReceivingPlayerID *string      `json:"receivingPlayerId"`
	Timestamp         int          `json:"timestamp"`
}

type CharacterInfo struct {
	Name                string `json:"name"`
	Points              int    `json:"points"`
	Position            int    `json:"position"`
	ColouredPosition    []int  `json:"colouredPositions"`
	StunnedForGameTicks int    `json:"stunnedForGameTicks"`
	ID                  string `json:"id"`
	CarryingPowerUp     bool   `json:"carryingPowerUp"`
}

type Map struct {
	Width               int             `json:"width"`
	Height              int             `json:"height"`
	WorldTick           int             `json:"worldTick"`
	CharacterInfos      []CharacterInfo `json:"characterInfos"`
	PowerUpPositions    []int           `json:"powerUpPositions"`
	ObstacleUpPositions []int           `json:"obstaclePositions"`
	CollisionInfos      []int           `json:"collisionInfos"`
	ExplosionInfos      []int           `json:"explosionInfos"`
}

type MapUpdateEvent struct {
	Type              string  `json:"type"`
	GameID            string  `json:"gameId"`
	GameTick          int     `json:"gameTick"`
	Map               Map     `json:"map"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type RegisterMoveEvent struct {
	Type              string  `json:"type"`
	GameID            string  `json:"gameId"`
	GameTick          int     `json:"gameTick"`
	Action            string  `json:"direction"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type PlayerRank struct {
	PlayerName string `json:"playerName"`
	PlayerId   string `json:"playerId"`
	Rank       int    `json:"rank"`
	Points     int    `json:"points"`
	Alive      bool   `json:"alive"`
}

type GameResultEvent struct {
	Type              string       `json:"type"`
	GameID            string       `json:"gameId"`
	PlayerRanks       []PlayerRank `json:"playerRanks"`
	ReceivingPlayerID *string      `json:"receivingPlayerId"`
	Timestamp         int          `json:"timestamp"`
}

type GameEndedEvent struct {
	Type              string  `json:"type"`
	PlayerWinnerID    string  `json:"playerWinnerId"`
	PlayerWinnerName  string  `json:"playerWinnerName"`
	GameID            string  `json:"gameId"`
	GameTick          int     `json:"gameTick"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type PlayerPoint struct {
	Name     string `json:"name"`
	PlayerID string `json:"playerId"`
	Points   int    `json:"points"`
}

type TournamentEndedEvent struct {
	Type              string        `json:"type"`
	PlayerWinnerID    string        `json:"playerWinnerId"`
	GameID            string        `json:"gameId"`
	GameResult        []PlayerPoint `json:"gameResult"`
	TournamentName    string        `json:"tournamentName"`
	TournamentID      string        `json:"tournamentId"`
	ReceivingPlayerID *string       `json:"receivingPlayerId"`
	Timestamp         int           `json:"timestamp"`
}

type InvalidMessage struct {
	Type              string  `json:"type"`
	ErrorMessage      string  `json:"errorMessage"`
	ReceivedMessage   string  `json:"receivedMessage"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}

type HearbeatMessage struct {
	Type              string  `json:"type"`
	ReceivingPlayerID *string `json:"receivingPlayerId"`
	Timestamp         int     `json:"timestamp"`
}
