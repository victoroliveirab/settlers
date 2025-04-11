import { resourcesOrder } from "@/core/constants";

export function isDirty({
  currentState,
  initialState,
}: {
  currentState: {
    given: SettlersCore.ResourceCollection;
    requested: SettlersCore.ResourceCollection;
  };
  initialState: {
    given: SettlersCore.ResourceCollection;
    requested: SettlersCore.ResourceCollection;
  };
}): boolean {
  for (const resource of resourcesOrder) {
    const initialGiven = initialState.given[resource];
    const currentGiven = currentState.given[resource];
    const initialRequested = initialState.requested[resource];
    const currentRequested = currentState.requested[resource];
    if (initialGiven !== currentGiven || initialRequested !== currentRequested) return true;
  }
  return false;
}

export function hasSameResourceInOfferAndRequest({
  given,
  requested,
}: {
  given: SettlersCore.ResourceCollection;
  requested: SettlersCore.ResourceCollection;
}): boolean {
  for (const resource of resourcesOrder) {
    const givenQuantity = given[resource];
    const requestedQuantity = requested[resource];
    if (givenQuantity > 0 && requestedQuantity > 0) {
      return true;
    }
  }
  return false;
}
