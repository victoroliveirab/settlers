export const sumOfResources = (resources: SettlersCore.ResourceCollection) => {
  let sum = 0;
  for (const resource of Object.keys(resources)) {
    sum += resources[resource as SettlersCore.Resource];
  }
  return sum;
};
