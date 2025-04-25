import { useMemo, useState } from "react";
import { Bar } from "react-chartjs-2";

import { Checkbox } from "@/components/ui/checkbox";
import { defaultDataLabelConfig, disableGrid } from "@/lib/chartjs";

import { type IStatisticsProps } from "../..";

const diceValues = [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];

export const DiceByPlayer = ({
  data,
  players,
}: {
  data: IStatisticsProps["diceStatsByPlayer"];
  players: SettlersCore.Player[];
}) => {
  const [checkboxesStates, setCheckboxesStates] = useState<Record<string, boolean>>(
    players.reduce(
      (acc, player) => ({
        ...acc,
        [player.name]: true,
      }),
      {},
    ),
  );
  const numberOfChecked = Object.values(checkboxesStates).reduce(
    (sum, curr) => (curr ? sum + 1 : sum),
    0,
  );

  const filteredPlayers = players.filter((player) => checkboxesStates[player.name]);
  const occurancesByPlayer = useMemo(() => {
    return filteredPlayers.map((player) => diceValues.map((value) => data[player.name][value]));
  }, [data, filteredPlayers]);

  const datasets = useMemo(() => {
    return occurancesByPlayer.map((occurance, index) => ({
      data: occurance,
      borderRadius: 4,
      backgroundColor: filteredPlayers[index].color.background,
    }));
  }, [occurancesByPlayer, filteredPlayers]);

  return (
    <div className="flex flex-col gap-4">
      <Bar
        data={{
          labels: diceValues,
          datasets,
        }}
        options={{
          layout: {
            padding: {
              top: 24,
            },
          },
          plugins: {
            datalabels: defaultDataLabelConfig,
            tooltip: {
              enabled: false,
            },
            legend: {
              display: false,
            },
          },
          scales: {
            x: disableGrid,
            y: disableGrid,
          },
        }}
      />
      <ul className="flex gap-2 z-10">
        {players.map((player) => (
          <li key={player.name} className="flex space-x-2">
            <Checkbox
              id={`dice-by-player-${player.name}`}
              checked={checkboxesStates[player.name]}
              disabled={checkboxesStates[player.name] && numberOfChecked === 1}
              onCheckedChange={(state) =>
                setCheckboxesStates((prev) => ({
                  ...prev,
                  [player.name]: !!state,
                }))
              }
            />
            <div className="grid gap-1.5 leading-none">
              <label
                htmlFor={`dice-by-player-${player.name}`}
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                {player.name}
              </label>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
