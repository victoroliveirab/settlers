import { create } from "zustand";

interface PlayerState {
  username: string | null;
}

export const usePlayerStore = create<PlayerState>(() => ({
  username: null,
}));

export const setUsername = (username: string | null) => {
  return usePlayerStore.setState({ username });
};
