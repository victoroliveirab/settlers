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
import { useUpdateCities } from "./updaters/useUpdateCities";
import { useUpdateRobbers } from "./updaters/useUpdateRobbers";
import { useUpdateBlockedTiles } from "./updaters/useUpdateBlockedTiles";
import { useDarkenTilesAfterDiceRoll } from "./updaters/useDarkenTilesAfterDiceRoll";

import "./map.css";

export const SettlersMap = () => {
  const ref = useRef<HTMLDivElement>(null);
  const [instance, setInstance] = useState<BaseMapRenderer | null>(null);
  const [tick, setTick] = useState(0); // force re-render on resize

  const username = usePlayerStore((state) => state.username);
  const mapName = useMatchStore((state) => state.mapName);
  const map = useMatchStore((state) => state.map);
  const ports = useMatchStore((state) => state.ports);
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
    renderer.render(map, ports);
    setInstance(renderer);
    setTick(1);

    const main = document.getElementsByTagName("main")[0];
    const observer = new ResizeObserver(() => {
      renderer.render(map, ports);
      setTick((prev) => prev + 1);
    });
    observer.observe(main);
  }, [instance, map, mapName, onEdgeClick, onTileClick, onVertexClick]);

  useUpdateCities(instance, tick);
  useUpdateEdges(instance, tick);
  useUpdateRoads(instance, tick);
  useUpdateSettlements(instance, tick);
  useUpdateVertices(instance, tick);
  useUpdateRobbers(instance, tick); // TODO: change name to useUpdateTiles
  useUpdateBlockedTiles(instance, tick);

  useDarkenTilesAfterDiceRoll(instance, tick);

  return <div ref={ref} className="h-full w-full flex items-center justify-center" />;
};
