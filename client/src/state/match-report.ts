import { create } from "zustand";

type Statistics = {
  diceStatsByPlayer: Record<SettlersCore.Player["name"], Record<number, number>>;
  generalDiceStats: Record<number, number>;
  longestRoadEvolution: Record<SettlersCore.Player["name"], number[]>;
  numberOfRobberiesByPlayer: Record<SettlersCore.Player["name"], number>;
  pointsEvolution: Record<SettlersCore.Player["name"], number[]> | null;
};

type MatchReportState = {
  endDatetime: string;
  pointsDistribution: Record<SettlersCore.Player["name"], Record<string, number>> | null;
  roundsPlayed: number;
  roomName: string;
  statistics: Statistics;
  startDatetime: string;
};

export const useMatchReportStore = create<MatchReportState>(() => ({
  roomName: "",
  endDatetime: "",
  pointsDistribution: null,
  roundsPlayed: 0,
  statistics: {
    diceStatsByPlayer: {},
    generalDiceStats: {},
    longestRoadEvolution: {},
    numberOfRobberiesByPlayer: {},
    pointsEvolution: null,
  },
  startDatetime: "",
}));

export const setStatistics = (data: Statistics) => {
  return useMatchReportStore.setState({ statistics: data });
};

export const setEndDatetime = (value: string) => {
  return useMatchReportStore.setState({ endDatetime: value });
};

export const setStartDatetime = (value: string) => {
  return useMatchReportStore.setState({ startDatetime: value });
};

export const setPointsDistribution = (value: Record<string, Record<string, number>>) => {
  return useMatchReportStore.setState({ pointsDistribution: value });
};

export const setRoundsPlayed = (value: number) => {
  return useMatchReportStore.setState({ roundsPlayed: value });
};

export const setRoomName = (value: string) => {
  return useMatchReportStore.setState({ roomName: value });
};
