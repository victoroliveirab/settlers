import { create } from "zustand";

type MatchStatisticsState = {
  data: {
    diceStatsByPlayer: Record<SettlersCore.Player["name"], Record<number, number>>;
    generalDiceStats: Record<number, number>;
    longestRoadEvolution: Record<SettlersCore.Player["name"], number[]>;
    numberOfRobberiesByPlayer: Record<SettlersCore.Player["name"], number>;
    pointsEvolution: Record<SettlersCore.Player["name"], number[]> | null;
  };
};

export const useMatchStatisticsStore = create<MatchStatisticsState>(() => ({
  data: {
    diceStatsByPlayer: {},
    generalDiceStats: {},
    longestRoadEvolution: {},
    numberOfRobberiesByPlayer: {},
    pointsEvolution: null,
  },
}));

export const setStatistics = (data: MatchStatisticsState["data"]) => {
  return useMatchStatisticsStore.setState({ data });
};
