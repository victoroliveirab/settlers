import { useMemo } from "react";
import { Bar } from "react-chartjs-2";

import { colors, defaultDataLabelConfig, disableGrid } from "@/lib/chartjs";

import { type IStatisticsProps } from "../..";

const diceValues = [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];

export const GeneralDiceStatistics = ({ data }: { data: IStatisticsProps["diceStats"] }) => {
  const occurances = useMemo(() => {
    return diceValues.map((sum) => data[sum]);
  }, [data]);

  return (
    <Bar
      data={{
        labels: diceValues,
        datasets: [
          {
            data: occurances,
            borderRadius: 4,
            backgroundColor: colors.bar,
          },
        ],
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
  );
};
