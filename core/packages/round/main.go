package round

type Type int

const (
	SetupSettlement1 Type = iota
	SetupRoad1
	SetupSettlement2
	SetupRoad2
	FirstRound
	Regular
	MoveRobberDue7
	MoveRobberDueKnight
	PickRobbed
	BetweenTurns
	BuildRoad1Development
	BuildRoad2Development
	MonopolyPickResource
	YearOfPlentyPickResources
	DiscardPhase
	GameOver
)

var roundTypeTranslation = [16]string{
	"SettlementSetup#1",
	"RoadSetup#1",
	"SettlementSetup#2",
	"RoadSetup#2",
	"FirstRound",
	"Regular",
	"MoveRobber(7)",
	"MoveRobber(Knight)",
	"ChooseRobbedPlayer",
	"BetweenRounds",
	"BuildRoadDevelopment(1)",
	"BuildRoadDevelopment(2)",
	"MonopolyPickResource",
	"YearOfPlentyPickResources",
	"DiscardPhase",
	"GameOver",
}

type Instance struct {
	dice        [2]int
	roundNumber int
	roundType   Type
}

func New() *Instance {
	return &Instance{
		dice:        [2]int{0, 0},
		roundNumber: 0,
		roundType:   SetupSettlement1,
	}
}

func (r *Instance) GetDice() [2]int {
	return r.dice
}

func (r *Instance) GetRoundType() Type {
	return r.roundType
}

func (r *Instance) GetCurrentRoundTypeDescription() string {
	return roundTypeTranslation[r.roundType]
}

func (r *Instance) GetRoundTypeDescription(roundType Type) string {
	return roundTypeTranslation[roundType]
}

func (r *Instance) GetRoundNumber() int {
	return r.roundNumber
}

func (r *Instance) SetDice(d1, d2 int) {
	r.dice[0] = d1
	r.dice[1] = d2
}

func (r *Instance) SetRoundType(roundType Type) {
	r.roundType = roundType
}

func (r *Instance) IncrementRound() {
	r.roundNumber++
}
