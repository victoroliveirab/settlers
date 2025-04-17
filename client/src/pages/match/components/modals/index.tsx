import { DiscardModal } from "./components/discard-modal";
import { MonopolyModal } from "./components/monopoly-modal";
import { PickRobbedModal } from "./components/pick-robbed-modal";
import { YearOfPlentyModal } from "./components/year-of-plenty-modal";

export const Modals = () => {
  return (
    <>
      <DiscardModal />
      <PickRobbedModal />
      <YearOfPlentyModal />
      <MonopolyModal />
    </>
  );
};
