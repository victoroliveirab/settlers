package entities

import (
	"fmt"
	"time"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/logger"
)

// FIXME: temporary copy
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

var phaseDurationSpeed30 = map[round.Type]time.Duration{
	round.SetupSettlement1:          15 * time.Second,
	round.SetupRoad1:                15 * time.Second,
	round.SetupSettlement2:          15 * time.Second,
	round.SetupRoad2:                15 * time.Second,
	round.FirstRound:                10 * time.Second,
	round.Regular:                   30 * time.Second,
	round.MoveRobberDue7:            10 * time.Second,
	round.MoveRobberDueKnight:       10 * time.Second,
	round.PickRobbed:                10 * time.Second,
	round.BetweenTurns:              10 * time.Second,
	round.BuildRoad1Development:     10 * time.Second,
	round.BuildRoad2Development:     10 * time.Second,
	round.MonopolyPickResource:      10 * time.Second,
	round.YearOfPlentyPickResources: 10 * time.Second,
	round.DiscardPhase:              10 * time.Second,
}

var phaseDurationSpeed45 = map[round.Type]time.Duration{
	round.SetupSettlement1:          20 * time.Second,
	round.SetupRoad1:                20 * time.Second,
	round.SetupSettlement2:          20 * time.Second,
	round.SetupRoad2:                20 * time.Second,
	round.FirstRound:                15 * time.Second,
	round.Regular:                   45 * time.Second,
	round.MoveRobberDue7:            15 * time.Second,
	round.MoveRobberDueKnight:       15 * time.Second,
	round.PickRobbed:                15 * time.Second,
	round.BetweenTurns:              15 * time.Second,
	round.BuildRoad1Development:     15 * time.Second,
	round.BuildRoad2Development:     15 * time.Second,
	round.MonopolyPickResource:      15 * time.Second,
	round.YearOfPlentyPickResources: 15 * time.Second,
	round.DiscardPhase:              15 * time.Second,
}

var phaseDurationSpeed60 = map[round.Type]time.Duration{
	round.SetupSettlement1:          30 * time.Second,
	round.SetupRoad1:                30 * time.Second,
	round.SetupSettlement2:          30 * time.Second,
	round.SetupRoad2:                30 * time.Second,
	round.FirstRound:                20 * time.Second,
	round.Regular:                   60 * time.Second,
	round.MoveRobberDue7:            20 * time.Second,
	round.MoveRobberDueKnight:       20 * time.Second,
	round.PickRobbed:                20 * time.Second,
	round.BetweenTurns:              20 * time.Second,
	round.BuildRoad1Development:     20 * time.Second,
	round.BuildRoad2Development:     20 * time.Second,
	round.MonopolyPickResource:      20 * time.Second,
	round.YearOfPlentyPickResources: 20 * time.Second,
	round.DiscardPhase:              20 * time.Second,
}

var phaseDurationSpeed75 = map[round.Type]time.Duration{
	round.SetupSettlement1:          40 * time.Second,
	round.SetupRoad1:                40 * time.Second,
	round.SetupSettlement2:          40 * time.Second,
	round.SetupRoad2:                40 * time.Second,
	round.FirstRound:                25 * time.Second,
	round.Regular:                   75 * time.Second,
	round.MoveRobberDue7:            25 * time.Second,
	round.MoveRobberDueKnight:       25 * time.Second,
	round.PickRobbed:                25 * time.Second,
	round.BetweenTurns:              25 * time.Second,
	round.BuildRoad1Development:     25 * time.Second,
	round.BuildRoad2Development:     25 * time.Second,
	round.MonopolyPickResource:      25 * time.Second,
	round.YearOfPlentyPickResources: 25 * time.Second,
	round.DiscardPhase:              25 * time.Second,
}

var phaseDurationSpeed90 = map[round.Type]time.Duration{
	round.SetupSettlement1:          45 * time.Second,
	round.SetupRoad1:                45 * time.Second,
	round.SetupSettlement2:          45 * time.Second,
	round.SetupRoad2:                45 * time.Second,
	round.FirstRound:                30 * time.Second,
	round.Regular:                   90 * time.Second,
	round.MoveRobberDue7:            30 * time.Second,
	round.MoveRobberDueKnight:       30 * time.Second,
	round.PickRobbed:                30 * time.Second,
	round.BetweenTurns:              30 * time.Second,
	round.BuildRoad1Development:     30 * time.Second,
	round.BuildRoad2Development:     30 * time.Second,
	round.MonopolyPickResource:      30 * time.Second,
	round.YearOfPlentyPickResources: 30 * time.Second,
	round.DiscardPhase:              30 * time.Second,
}

var phaseDurationsBySpeed = map[int]map[round.Type]time.Duration{
	30: phaseDurationSpeed30,
	45: phaseDurationSpeed45,
	60: phaseDurationSpeed60,
	75: phaseDurationSpeed75,
	90: phaseDurationSpeed90,
}

func newRoundManager(speed int, onRegularTimeout func(), onExpireFuncs map[round.Type]func()) *roundManager {
	_, ok := phaseDurationsBySpeed[speed]
	if !ok {
		speed = 60
	}
	fmt.Println("SPEED", speed)

	return &roundManager{
		speed:         speed,
		onTimeout:     onRegularTimeout,
		onExpireFuncs: onExpireFuncs,
	}
}

func (rm *roundManager) start() {
	logger.LogSystemMessage("room.roundManager.start", "start()")
	rm.Lock()
	defer rm.Unlock()

	rm.remaining = phaseDurationsBySpeed[rm.speed][round.Regular]
	deadline := time.Now().UTC().Add(rm.remaining)
	rm.deadline = &deadline
	rm.subPhaseDeadline = nil

	rm.timer = time.AfterFunc(rm.remaining, rm.onTimeout)
	logger.LogSystemMessage("room.roundManager.start", "starting")
}

func (rm *roundManager) pause() {
	logger.LogSystemMessage("room.roundManager.pause", "pause()")
	rm.Lock()
	defer rm.Unlock()

	if rm.timer != nil {
		if rm.timer.Stop() {
			logger.LogSystemMessage("room.roundManager.pause", "pausing")
			rm.remaining = time.Until(*rm.deadline)
		}
	}
}

func (rm *roundManager) resume() {
	logger.LogSystemMessage("room.roundManager.resume", "resume()")
	rm.Lock()
	defer rm.Unlock()

	newDeadline := time.Now().UTC().Add(rm.remaining)
	rm.deadline = &newDeadline
	rm.timer = time.AfterFunc(rm.remaining, rm.onTimeout)
	logger.LogSystemMessage("room.roundManager.resume", "resuming")
}

func (rm *roundManager) cancel() {
	logger.LogSystemMessage("room.roundManager.cancel", "cancel()")
	rm.Lock()
	defer rm.Unlock()

	if rm.timer != nil {
		logger.LogSystemMessage("room.roundManager.cancel", "stopping timer")
		rm.timer.Stop()
	}
	rm.cancelSubTimer()

	rm.deadline = nil
	rm.subPhaseDeadline = nil
}

func (rm *roundManager) startPhaseTimer(phase round.Type) {
	logger.LogSystemMessage("room.roundManager.startPhaseTimer", fmt.Sprintf("phase = %s", roundTypeTranslation[phase]))
	rm.Lock()
	defer rm.Unlock()

	rm.cancelSubTimer()

	dur := phaseDurationsBySpeed[rm.speed][round.Type(phase)]

	onExpire := rm.onExpireFuncs[phase]
	rm.subTimer = time.AfterFunc(dur, onExpire)
	subPhaseDeadline := time.Now().UTC().Add(dur)
	rm.subPhaseDeadline = &subPhaseDeadline
}

func (rm *roundManager) cancelSubTimer() {
	logger.LogSystemMessage("room.roundManager.cancelSubTimer", "cancelSubTimer()")
	if rm.subTimer != nil {
		rm.subTimer.Stop()
		rm.subTimer = nil
	}
	rm.subPhaseDeadline = nil
}

func (rm *roundManager) Deadline() *time.Time {
	rm.Lock()
	defer rm.Unlock()
	return rm.deadline
}

func (rm *roundManager) SubPhaseDeadline() *time.Time {
	rm.Lock()
	defer rm.Unlock()
	return rm.subPhaseDeadline
}

func (rm *roundManager) Now() time.Time {
	return time.Now().UTC()
}
