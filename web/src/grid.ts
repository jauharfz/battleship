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
      grid.append(cell);
    }
  }
  return cells;
}
