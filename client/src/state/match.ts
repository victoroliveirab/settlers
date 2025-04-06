import { create } from "zustand";

type MatchState = {
  actions: {
    pass: boolean;
    trade: boolean;
  };
  activeTradeOffers: {
    counters: number[];
    finalized: boolean;
    id: number;
    offer: Record<SettlersCore.Resource, number>;
    opponents: Record<
      SettlersCore.Player["name"],
      { status: "Open" | "Accepted" | "Declined"; blocked: boolean }
    >;
    parent: number;
    player: SettlersCore.Player["name"];
    request: Record<SettlersCore.Resource, number>;
    status: "Open" | "Closed";
    timestamp: number;
  }[];
  cities: SettlersCore.Cities;
  currentRoundPlayer: string;
  devHand: SettlersCore.DevHand;
  dice: {
    enabled: boolean;
    value: [number, number];
  };
  edges: {
    available: number[];
    enabled: boolean;
    highlight: boolean;
  };
  hand: SettlersCore.Hand;
  logs: string[];
  map: SettlersCore.Map;
  mapName: string;
  ports: SettlersCore.Ports;
  players: SettlersCore.Player[];
  resourceCount: Record<SettlersCore.Player["name"], number>;
  roads: SettlersCore.Roads;
  settlements: SettlersCore.Settlements;
  vertices: {
    availableForCity: number[];
    availableForSettlement: number[];
    enabled: boolean;
    highlight: boolean;
  };
};

export const useMatchStore = create<MatchState>(() => ({
  actions: {
    pass: false,
    trade: false,
  },
  activeTradeOffers: [],
  cities: {},
  currentRoundPlayer: "",
  devHand: {
    Knight: 0,
    Monopoly: 0,
    "Road Building": 0,
    "Victory Point": 0,
    "Year of Plenty": 0,
  },
  dice: {
    enabled: false,
    value: [0, 0],
  },
  edges: {
    available: [],
    enabled: false,
    highlight: false,
  },
  hand: {
    Lumber: 0,
    Brick: 0,
    Sheep: 0,
    Grain: 0,
    Ore: 0,
  },
  logs: [],
  map: [],
  mapName: "",
  roads: {},
  settlements: {},
  ports: [],
  players: [],
  resourceCount: {},
  vertices: {
    availableForCity: [],
    availableForSettlement: [],
    enabled: false,
    highlight: false,
  },
}));

export const setPassAction = (state: boolean) => {
  const { actions } = useMatchStore.getState();
  return useMatchStore.setState({
    actions: {
      ...actions,
      pass: state,
    },
  });
};

export const setTradeAction = (state: boolean) => {
  const { actions } = useMatchStore.getState();
  return useMatchStore.setState({
    actions: {
      ...actions,
      trade: state,
    },
  });
};

export const setActiveTradeOffers = (value: MatchState["activeTradeOffers"]) => {
  return useMatchStore.setState({ activeTradeOffers: value });
};

export const setCities = (value: SettlersCore.Cities) => {
  return useMatchStore.setState({ cities: value });
};

export const setCurrentRoundPlayer = (name: string) => {
  return useMatchStore.setState({ currentRoundPlayer: name });
};

export const setDevHand = (value: SettlersCore.DevHand) => {
  return useMatchStore.setState({ devHand: value });
};

export const setDice = (value: { enabled: boolean; value: [number, number] }) => {
  return useMatchStore.setState({ dice: value });
};

export const setEdges = (available: number[], highlight: boolean, enabled: boolean) => {
  return useMatchStore.setState({
    edges: {
      available,
      enabled,
      highlight,
    },
  });
};

export const setHand = (value: SettlersCore.Hand) => {
  return useMatchStore.setState({ hand: value });
};

export const setMap = (value: SettlersCore.Map) => {
  return useMatchStore.setState({ map: value });
};

export const setMapName = (mapName: string) => {
  return useMatchStore.setState({ mapName });
};

export const setLogs = (logs: string[]) => {
  // This will overwrite logs everytime on purpose.
  return useMatchStore.setState({ logs });
};

export const setRoads = (value: SettlersCore.Roads) => {
  return useMatchStore.setState({ roads: value });
};

export const setSettlements = (value: SettlersCore.Settlements) => {
  return useMatchStore.setState({ settlements: value });
};

export const setPorts = (value: SettlersCore.Ports) => {
  return useMatchStore.setState({ ports: value });
};

export const setPlayers = (value: SettlersCore.Player[]) => {
  return useMatchStore.setState({ players: value });
};

export const setResourceCount = (value: Record<string, number>) => {
  return useMatchStore.setState({ resourceCount: value });
};

export const setVertices = (
  availableForSettlement: number[],
  availableForCity: number[],
  highlight: boolean,
  enabled: boolean,
) => {
  return useMatchStore.setState({
    vertices: {
      availableForCity,
      availableForSettlement,
      enabled,
      highlight,
    },
  });
};
