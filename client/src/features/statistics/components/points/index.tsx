import { useMemo } from "react";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { type IStatisticsProps } from "../..";

export const Points = ({
  data,
  players,
}: {
  data: NonNullable<IStatisticsProps["pointsDistribution"]>;
  players: SettlersCore.Player[];
}) => {
  const rows = useMemo(() => {
    // FIXME: the order should come from the server
    const sortedPlayers = [...players].sort((playerA, playerB) => {
      const aName = playerA.name;
      const bName = playerB.name;
      if (data[aName].total !== data[bName].total) {
        return data[aName].total < data[bName].total ? 1 : -1;
      }
      return 0;
    });
    return sortedPlayers.map((player) => ({
      data: data[player.name],
      player,
    }));
  }, [data, players]);

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]"></TableHead>
          <TableHead>Player</TableHead>
          <TableHead>🏠</TableHead>
          <TableHead>🏢</TableHead>
          <TableHead>🎖️</TableHead>
          <TableHead>⚔️</TableHead>
          <TableHead>🛤️</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((row, index) => (
          <TableRow key={row.player.name}>
            <TableCell className="font-medium">{index + 1}</TableCell>
            <TableCell>{row.data.total}</TableCell>
            <TableCell>{row.data.settlements}</TableCell>
            <TableCell>{row.data.cities}</TableCell>
            <TableCell>{row.data.victoryPoints}</TableCell>
            <TableCell>{row.data.largestArmy}</TableCell>
            <TableCell>{row.data.longestRoad}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};
