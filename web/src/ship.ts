import { ShipType, Ship } from "./types";

export const shipSizes: Map<ShipType, number> = new Map([
  ["aircraft_carrier", 5],
  ["battleship", 4],
  ["heavy_cruiser", 3],
  ["light_cruiser", 3],
  ["destroyer", 2],
]);
// there's need validation each ship
// and each ship need iterates the length/sizes ship itself
// okey, just focus create client validation for each ship first

export async function shipPlacementServerValidation(
  ships: Ship[],
  token: string,
) {
  try {
    const res = await fetch("/api/v1/validate", {
      method: "QUERY",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(ships),
    });
    const payload = await res.json();
    if (payload.status == "error") {
      throw new Error(payload.error);
    }
    return;
  } catch (error) {
    throw error;
  }
}

export function renderShipPlacement(cells: HTMLElement[][], ships: Ship[]) {
  traverseShips(cells, ships, true);
}

export function shipPlacementClientValidation(
  cells: HTMLElement[][],
  ships: Ship[],
): boolean {
  return traverseShips(cells, ships, false);
}

// TODO: add another validation for outbound
function traverseShips(
  cells: HTMLElement[][],
  ships: Ship[],
  isValid: boolean,
): boolean {
  for (let ship of ships) {
    for (let i = 0; i < shipSizes.get(ship.type)!; i++) {
      const currX = ship.direction == "Horizontal" ? ship.x + i : ship.x;
      const currY = ship.direction == "Vertical" ? ship.y + i : ship.y;
      const cell = cells[currY][currX];
      if (cell.dataset.state) {
        return false;
      }
      if (isValid) {
        cell.dataset.state = "ship";
      }
    }
  }
  return true;
}
