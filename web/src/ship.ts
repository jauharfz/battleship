import { shipSizes } from "./constants";
import { Ship } from "./types";

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

export function renderShipPlacement(
  cells: HTMLElement[][],
  ships: Ship[],
  shipsIndex: number[][],
): number[][] {
  shipsIndex = Array.from({ length: 10 }, () => new Array<number>(10).fill(-1));
  traverseShips(cells, ships, shipsIndex, true);
  return shipsIndex;
}

export function shipPlacementClientValidation(
  cells: HTMLElement[][],
  ships: Ship[],
  shipsIndex: number[][],
): boolean {
  return traverseShips(cells, ships, shipsIndex, false);
}

// TODO: add another validation for outbound
function traverseShips(
  cells: HTMLElement[][],
  ships: readonly Ship[],
  shipsIndex: number[][],
  isValid: boolean,
): boolean {
  let index = 0;
  for (let ship of ships) {
    for (let i = 0; i < shipSizes[ship.type]; i++) {
      const currX = ship.direction == "Horizontal" ? ship.x + i : ship.x;
      if (currX > 9 || currX < 0) return false;
      const currY = ship.direction == "Vertical" ? ship.y + i : ship.y;
      if (currY > 9 || currY < 0) return false;
      // buildOccupancy
      const cell = cells[currY][currX];
      if (isValid) {
        shipsIndex[currY][currX] = index;
        cell.dataset.state = "ship";
      } else {
        if (cell.dataset.state || shipsIndex[currY][currX] != -1) {
          return false;
        }
      }
    }
    if (isValid) index++;
  }
  return true;
}

// validate col&row with occupancy? or shipsIndex?
export function validatePlacement(
  col: number,
  row: number,
  ship: Ship,
  shipsIndex: number[][],
  currentShipIndex: number,
): boolean {
  for (let i = 0; i < shipSizes[ship.type]; i++) {
    const currX = ship.direction == "Horizontal" ? col + i : col;
    if (currX > 9 || currX < 0) return false;
    const currY = ship.direction == "Vertical" ? row + i : row;
    if (currY > 9 || currY < 0) return false;
    if (
      shipsIndex[currY][currX] != -1 &&
      shipsIndex[currY][currX] != currentShipIndex
    ) {
      return false;
    }
  }
  return true;
}

export function renderPreview(
  col: number,
  row: number,
  ship: Ship,
  isValid: boolean,
  cells: HTMLElement[][],
) {
  for (let i = 0; i < shipSizes[ship.type]; i++) {
    const currX = ship.direction == "Horizontal" ? col + i : col;
    const currY = ship.direction == "Vertical" ? row + i : row;
    if (currX >= 0 && currX <= 9 && currY >= 0 && currY <= 9) {
      cells[currY][currX].dataset.preview = isValid ? "valid" : "invalid";
    }
  }
}

export function clearPreview(
  oldCol: number,
  oldRow: number,
  ship: Ship,
  cells: HTMLElement[][],
) {
  for (let i = 0; i < shipSizes[ship.type]; i++) {
    const currX = ship.direction == "Horizontal" ? oldCol + i : oldCol;
    const currY = ship.direction == "Vertical" ? oldRow + i : oldRow;
    if (currX >= 0 && currX <= 9 && currY >= 0 && currY <= 9) {
      delete cells[currY][currX].dataset.preview;
    }
  }
}

export function clearState(
  col: number,
  row: number,
  ship: Ship,
  cells: HTMLElement[][],
) {
  for (let i = 0; i < shipSizes[ship.type]; i++) {
    const currX = ship.direction == "Horizontal" ? col + i : col;
    const currY = ship.direction == "Vertical" ? row + i : row;
    if (currX >= 0 && currX <= 9 && currY >= 0 && currY <= 9) {
      delete cells[currY][currX].dataset.state;
    }
  }
}
