import { create } from "zustand";

type Color = {
  background: string;
  foreground: string;
};

type RoomParticipant = {
  player: {
    name: string;
    color: Color;
  } | null;
  ready: boolean;
  bot: boolean;
};

type Room = {
  id: string;
  capacity: number;
  map: string;
  participants: RoomParticipant[];
  private: boolean;
  owner: string;
  status: string;
  colors: Color[];
};

export type RoomParam = {
  description: string;
  key: string;
  label: string;
  value: number;
  values: number[];
};

interface RoomState {
  minMaxPlayers: [number, number];
  room: Room;
  params: RoomParam[];
}

export const useRoomStore = create<RoomState>(() => ({
  minMaxPlayers: [0, 0],
  room: {
    id: "",
    capacity: 0,
    map: "",
    participants: [],
    private: true,
    owner: "",
    status: "",
    colors: [],
  },
  params: [],
}));

export const setRoom = (room: Room) => {
  return useRoomStore.setState({ room });
};

export const setRoomCapacity = (minMax: [number, number]) => {
  return useRoomStore.setState({ minMaxPlayers: minMax });
};

export const setRoomParams = (params: RoomParam[]) => {
  return useRoomStore.setState({ params });
};

export const setRoomStatus = (status: string) => {
  const state = useRoomStore.getState();
  return useRoomStore.setState({
    room: {
      ...state.room,
      status,
    },
  });
};
