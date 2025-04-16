package entities

import (
	"fmt"
	"time"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/logger"
)

var phaseDurationSpeed30 = map[int]time.Duration{
	core.SetupSettlement1:          15 * time.Second,
	core.SetupRoad1:                15 * time.Second,
	core.SetupSettlement2:          15 * time.Second,
	core.SetupRoad2:                15 * time.Second,
	core.FirstRound:                10 * time.Second,
	core.Regular:                   30 * time.Second,
	core.MoveRobberDue7:            10 * time.Second,
	core.MoveRobberDueKnight:       10 * time.Second,
	core.PickRobbed:                10 * time.Second,
	core.BetweenTurns:              10 * time.Second,
	core.BuildRoad1Development:     10 * time.Second,
	core.BuildRoad2Development:     10 * time.Second,
	core.MonopolyPickResource:      10 * time.Second,
	core.YearOfPlentyPickResources: 10 * time.Second,
	core.DiscardPhase:              10 * time.Second,
}

var phaseDurationSpeed45 = map[int]time.Duration{
	core.SetupSettlement1:          20 * time.Second,
	core.SetupRoad1:                20 * time.Second,
	core.SetupSettlement2:          20 * time.Second,
	core.SetupRoad2:                20 * time.Second,
	core.FirstRound:                15 * time.Second,
	core.Regular:                   45 * time.Second,
	core.MoveRobberDue7:            15 * time.Second,
	core.MoveRobberDueKnight:       15 * time.Second,
	core.PickRobbed:                15 * time.Second,
	core.BetweenTurns:              15 * time.Second,
	core.BuildRoad1Development:     15 * time.Second,
	core.BuildRoad2Development:     15 * time.Second,
	core.MonopolyPickResource:      15 * time.Second,
	core.YearOfPlentyPickResources: 15 * time.Second,
	core.DiscardPhase:              15 * time.Second,
}

var phaseDurationSpeed60 = map[int]time.Duration{
	core.SetupSettlement1:          30 * time.Second,
	core.SetupRoad1:                30 * time.Second,
	core.SetupSettlement2:          30 * time.Second,
	core.SetupRoad2:                30 * time.Second,
	core.FirstRound:                20 * time.Second,
	core.Regular:                   60 * time.Second,
	core.MoveRobberDue7:            20 * time.Second,
	core.MoveRobberDueKnight:       20 * time.Second,
	core.PickRobbed:                20 * time.Second,
	core.BetweenTurns:              20 * time.Second,
	core.BuildRoad1Development:     20 * time.Second,
	core.BuildRoad2Development:     20 * time.Second,
	core.MonopolyPickResource:      20 * time.Second,
	core.YearOfPlentyPickResources: 20 * time.Second,
	core.DiscardPhase:              20 * time.Second,
}

var phaseDurationSpeed75 = map[int]time.Duration{
	core.SetupSettlement1:          40 * time.Second,
	core.SetupRoad1:                40 * time.Second,
	core.SetupSettlement2:          40 * time.Second,
	core.SetupRoad2:                40 * time.Second,
	core.FirstRound:                25 * time.Second,
	core.Regular:                   75 * time.Second,
	core.MoveRobberDue7:            25 * time.Second,
	core.MoveRobberDueKnight:       25 * time.Second,
	core.PickRobbed:                25 * time.Second,
	core.BetweenTurns:              25 * time.Second,
	core.BuildRoad1Development:     25 * time.Second,
	core.BuildRoad2Development:     25 * time.Second,
	core.MonopolyPickResource:      25 * time.Second,
	core.YearOfPlentyPickResources: 25 * time.Second,
	core.DiscardPhase:              25 * time.Second,
}

var phaseDurationSpeed90 = map[int]time.Duration{
	core.SetupSettlement1:          45 * time.Second,
	core.SetupRoad1:                45 * time.Second,
	core.SetupSettlement2:          45 * time.Second,
	core.SetupRoad2:                45 * time.Second,
	core.FirstRound:                30 * time.Second,
	core.Regular:                   90 * time.Second,
	core.MoveRobberDue7:            30 * time.Second,
	core.MoveRobberDueKnight:       30 * time.Second,
	core.PickRobbed:                30 * time.Second,
	core.BetweenTurns:              30 * time.Second,
	core.BuildRoad1Development:     30 * time.Second,
	core.BuildRoad2Development:     30 * time.Second,
	core.MonopolyPickResource:      30 * time.Second,
	core.YearOfPlentyPickResources: 30 * time.Second,
	core.DiscardPhase:              30 * time.Second,
}

var phaseDurationsBySpeed = map[int]map[int]time.Duration{
	30: phaseDurationSpeed30,
	45: phaseDurationSpeed45,
	60: phaseDurationSpeed60,
	75: phaseDurationSpeed75,
	90: phaseDurationSpeed90,
}

func newRoundManager(speed int, onRegularTimeout func(), onExpireFuncs map[int]func()) *roundManager {
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

	rm.remaining = phaseDurationsBySpeed[rm.speed][core.Regular]
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

func (rm *roundManager) startPhaseTimer(phase int) {
	logger.LogSystemMessage("room.roundManager.startPhaseTimer", fmt.Sprintf("phase = %s", core.RoundTypeTranslation[phase]))
	rm.Lock()
	defer rm.Unlock()

	rm.cancelSubTimer()

	dur := phaseDurationsBySpeed[rm.speed][phase]

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
