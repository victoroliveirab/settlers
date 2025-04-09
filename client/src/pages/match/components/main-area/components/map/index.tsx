import { useEffect, useRef, useState } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { mapRendererFactory } from "@/core/maps";
import { useMatchStore } from "@/state/match";
import { usePlayerStore } from "@/state/player";

import { useOnVertexClick } from "./handlers/useOnVertexClick";
import { useOnTileClick } from "./handlers/useOnTileClick";
import { useOnEdgeClick } from "./handlers/useOnEdgeClick";
import { useUpdateEdges } from "./updaters/useUpdateEdges";
import { useUpdateRoads } from "./updaters/useUpdateRoads";
import { useUpdateSettlements } from "./updaters/useUpdateSettlements";
import { useUpdateVertices } from "./updaters/useUpdateVertices";

import "./map.css";
import { useUpdateCities } from "./updaters/useUpdateCities";
import { useUpdateRobbers } from "./updaters/useUpdateRobbers";
import { useUpdateBlockedTiles } from "./updaters/useUpdateBlockedTiles";

export const SettlersMap = () => {
  const ref = useRef<HTMLDivElement>(null);
  const [instance, setInstance] = useState<BaseMapRenderer | null>(null);

  const username = usePlayerStore((state) => state.username);
  const mapName = useMatchStore((state) => state.mapName);
  const map = useMatchStore((state) => state.map);
  const players = useMatchStore((state) => state.players);

  const onEdgeClick = useOnEdgeClick();
  const onTileClick = useOnTileClick();
  const onVertexClick = useOnVertexClick();

  useEffect(() => {
    if (!ref.current || !username || instance) return;
    const renderer = mapRendererFactory(
      mapName,
      ref.current,
      players.reduce(
        (acc, player) => ({
          ...acc,
          [player.name]: player.color,
        }),
        {},
      ),
      username,
      { onEdgeClick, onTileClick, onVertexClick },
    );
    if (!renderer) return;
    renderer.render(map, []);
    setInstance(renderer);
  }, [instance, map, mapName, onEdgeClick, onTileClick, onVertexClick]);

  useUpdateCities(instance);
  useUpdateEdges(instance);
  useUpdateRoads(instance);
  useUpdateSettlements(instance);
  useUpdateVertices(instance);
  useUpdateRobbers(instance); // TODO: change name to useUpdateTiles
  useUpdateBlockedTiles(instance);

  return <div ref={ref} className="h-full w-full flex items-center justify-center" />;
};
