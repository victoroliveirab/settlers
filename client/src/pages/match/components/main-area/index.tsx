import { ActiveTrades } from "./components/active-trades";
import { DevHand } from "./components/dev-hand";
import { Hand } from "./components/hand";
import { SettlersMap } from "./components/map";
import { Players } from "./components/players";
import { Statistics } from "./components/statistics";
import { Timer } from "./components/timer";

export const MainArea = () => {
  return (
    <section className="h-full w-3/4 flex flex-col gap-2 xl:gap-4 justify-between relative">
      <div>
        <Players />
      </div>
      <hr className="w-full h-px bg-border" />
      <div className="overflow-hidden flex-1 relative">
        <SettlersMap />
        <div className="absolute top-0 right-0 w-80 flex flex-col gap-1">
          <ActiveTrades />
        </div>
        <div className="absolute bottom-0 left-0 flex flex-col">
          <Timer />
        </div>
        <div className="absolute top-0 left-0">
          <Statistics />
        </div>
      </div>
      <hr className="w-full h-px bg-border" />
      <div className="h-20 flex">
        <div className="flex-1">
          <Hand />
        </div>
        <div>
          <DevHand />
        </div>
      </div>
    </section>
  );
};
