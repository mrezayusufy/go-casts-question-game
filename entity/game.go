package entity

type Game struct {
	ID          uint
	CategoryID  uint
	QuestionIDs []uint
	Player      []Player
}

type Player struct {
	ID      uint
	UserID  uint
	GameID  uint
	Score   uint
	Answers []PlayerAnswer
}

type PlayerAnswer struct {
	ID, PlayerID, QuestionID uint
	Choice                   PossibleAnswerChoice
}
