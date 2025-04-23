import type { CartesianScaleTypeRegistry, ScaleOptionsByType } from "chart.js";
import type { Options } from "chartjs-plugin-datalabels/types/options";

export const colors = {
  bar: "#444",
  datalabel: "#444",
};

export const disableGrid: Partial<ScaleOptionsByType<keyof CartesianScaleTypeRegistry>> = {
  grid: {
    display: false,
  },
};

export const defaultDataLabelConfig: Partial<Options> = {
  anchor: "end",
  align: "end",
  color: colors.datalabel,
  offset: -4,
};
