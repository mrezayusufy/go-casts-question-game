package entity

type Game struct {
	ID          uint
	CategoryID  uint
	QuestionIDs []uint
	PlayerIDs   []uint
}

type Player struct {
	ID         uint
	UserID     uint
	GameID     uint
	Score      uint
	AnswersIDs []uint
}

type PlayerAnswer struct {
	ID, PlayerID, QuestionID uint
	Choice                   PossibleAnswerChoice
}
