import { FireData } from "./types";
import { send } from "./websocket";

export function createGrid(containerID: string): HTMLElement[][] {
  const grid = document.getElementById(containerID);
  if (!grid || !(grid instanceof HTMLElement)) return [];
  const cells: HTMLElement[][] = [];
  for (let y = 0; y < 10; y++) {
    cells[y] = [];
    for (let x = 0; x < 10; x++) {
      const cell = document.createElement("div");
      cell.classList.add("cell");
      cells[y][x] = cell;
      cell.dataset.y = y.toString();
      cell.dataset.x = x.toString();
      grid.append(cell);
    }
  }
  return cells;
}

export async function handleOpponentGrid(event: Event) {
  const cell = (event.target as HTMLElement).closest("[data-y][data-y]");
  if (!cell) return;
  if (!(cell instanceof HTMLElement)) return;

  send("CMD_FIRE", { x: Number(cell.dataset.x), y: Number(cell.dataset.y) });
}
