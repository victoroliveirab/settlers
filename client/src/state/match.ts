import { create } from "zustand";

type QuantityByPlayer = Record<SettlersCore.Player["name"], number>;

type MatchState = {
  actions: {
    buyDevCard: boolean;
    pass: boolean;
    trade: boolean;
  };
  activeTradeOffers: {
    creator: SettlersCore.Player["name"];
    finalized: boolean;
    id: number;
    offer: Record<SettlersCore.Resource, number>;
    parent: number;
    request: Record<SettlersCore.Resource, number>;
    requester: SettlersCore.Player["name"];
    responses: Record<
      SettlersCore.Player["name"],
      { status: "Open" | "Accepted" | "Declined"; blocked: boolean }
    >;
    status: "Open" | "Closed";
    timestamp: number;
  }[];
  blockedTiles: number[];
  cities: SettlersCore.Cities;
  currentRoundPlayer: {
    deadline: string;
    player: string;
    serverNow: string;
    subDeadline: string | null;
  } | null;
  devHand: SettlersCore.DevHand;
  devHandCount: Record<SettlersCore.Player["name"], number>;
  devHandPermissions: Record<SettlersCore.DevelopmentCard, boolean>;
  dice: {
    enabled: boolean;
    value: [number, number];
  };
  discard: {
    discardAmounts: Record<SettlersCore.Player["name"], number>;
    enabled: boolean;
  };
  edges: {
    available: number[];
    enabled: boolean;
    highlight: boolean;
  };
  hand: SettlersCore.Hand;
  knightUsages: QuantityByPlayer;
  longestRoadSize: QuantityByPlayer;
  logs: string[];
  map: SettlersCore.Map;
  mapName: string;
  monopoly: {
    enabled: boolean;
  };
  ownedPorts: SettlersCore.PortType[];
  points: QuantityByPlayer;
  ports: SettlersCore.Ports;
  players: SettlersCore.Player[];
  resourceCount: QuantityByPlayer;
  roads: SettlersCore.Roads;
  robbablePlayers: {
    enabled: boolean;
    options: SettlersCore.Player["name"][] | null;
  };
  robber: {
    availableTiles: number[];
    enabled: boolean;
    highlight: boolean;
  };
  settlements: SettlersCore.Settlements;
  vertices: {
    availableForCity: number[];
    availableForSettlement: number[];
    enabled: boolean;
    highlight: boolean;
  };
  yearOfPlenty: {
    enabled: boolean;
  };
};

export const useMatchStore = create<MatchState>(() => ({
  actions: {
    buyDevCard: false,
    pass: false,
    trade: false,
  },
  activeTradeOffers: [],
  blockedTiles: [],
  cities: {},
  currentRoundPlayer: null,
  devHand: {
    Knight: 0,
    Monopoly: 0,
    "Road Building": 0,
    "Victory Point": 0,
    "Year of Plenty": 0,
  },
  devHandCount: {},
  devHandPermissions: {
    Knight: false,
    Monopoly: false,
    "Road Building": false,
    "Victory Point": false,
    "Year of Plenty": false,
  },
  dice: {
    enabled: false,
    value: [0, 0],
  },
  discard: {
    discardAmounts: {},
    enabled: false,
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
  knightUsages: {},
  longestRoadSize: {},
  logs: [],
  map: [],
  mapName: "",
  monopoly: {
    enabled: false,
  },
  ownedPorts: [],
  points: {},
  roads: {},
  robber: {
    availableTiles: [],
    enabled: false,
    highlight: false,
  },
  robbablePlayers: {
    enabled: false,
    options: [],
  },
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
  yearOfPlenty: {
    enabled: false,
  },
}));

export const setBuyDevCardAction = (state: boolean) => {
  const { actions } = useMatchStore.getState();
  return useMatchStore.setState({
    actions: {
      ...actions,
      buyDevCard: state,
    },
  });
};

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

export const setCurrentRoundPlayer = (
  player: string,
  deadline: string,
  subDeadline: string | null,
  serverNow: string,
) => {
  return useMatchStore.setState({
    currentRoundPlayer: {
      deadline,
      player,
      serverNow,
      subDeadline,
    },
  });
};

export const setDevHand = (value: SettlersCore.DevHand) => {
  return useMatchStore.setState({ devHand: value });
};

export const setDevHandPermissions = (value: Record<SettlersCore.DevelopmentCard, boolean>) => {
  return useMatchStore.setState({ devHandPermissions: value });
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

export const setBlockedTiles = (value: number[]) => {
  return useMatchStore.setState({ blockedTiles: value });
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

export const setDevHandCount = (value: Record<string, number>) => {
  return useMatchStore.setState({ devHandCount: value });
};

export const setPoints = (value: Record<string, number>) => {
  return useMatchStore.setState({ points: value });
};

export const setLongestRoadSizes = (value: Record<string, number>) => {
  return useMatchStore.setState({ longestRoadSize: value });
};

export const setKnightUsages = (value: Record<string, number>) => {
  return useMatchStore.setState({ knightUsages: value });
};

export const setPlayerPorts = (value: SettlersCore.PortType[]) => {
  return useMatchStore.setState({ ownedPorts: value });
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

export const setDiscard = (amounts: Record<string, number>, enabled: boolean) => {
  return useMatchStore.setState({
    discard: {
      discardAmounts: amounts,
      enabled,
    },
  });
};

export const setRobber = (value: MatchState["robber"]) => {
  return useMatchStore.setState({
    robber: value,
  });
};

export const setRobbablePlayers = (value: MatchState["robbablePlayers"]) => {
  return useMatchStore.setState({
    robbablePlayers: value,
  });
};

export const setMonopoly = (value: boolean) => {
  return useMatchStore.setState({ monopoly: { enabled: value } });
};

export const setYearOfPlenty = (value: boolean) => {
  return useMatchStore.setState({ yearOfPlenty: { enabled: value } });
};
