import { send } from "./websocket";
export function createGrid(containerID) {
    const grid = document.getElementById(containerID);
    if (!grid || !(grid instanceof HTMLElement))
        return [];
    const cells = [];
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
export async function handleOpponentGrid(event, isMyTurn, state) {
    if (state != "INBATTLE")
        return;
    if (!isMyTurn)
        return;
    const cell = event.target.closest("[data-y][data-x]");
    if (!cell)
        return;
    if (!(cell instanceof HTMLElement))
        return;
    if (cell.dataset.state == "hit" || cell.dataset.state == "miss") {
        console.warn("grid already fired, pls fire another grid");
        return;
    }
    send("CMD_FIRE", { x: Number(cell.dataset.x), y: Number(cell.dataset.y) });
}
export function getCellCoordinate(grid, event) {
    const style = getComputedStyle(grid);
    const rect = grid.getBoundingClientRect();
    const borderLeft = parseFloat(style.borderLeft);
    const contentWidth = rect.width - borderLeft - parseFloat(style.borderRight);
    const pitch = (contentWidth + parseFloat(style.gap)) / 10;
    const col = (event.clientX - rect.left - borderLeft) / pitch;
    const row = (event.clientY - rect.top - parseFloat(style.borderTop)) / pitch;
    return [Math.floor(col), Math.floor(row)];
}
export function renderPreviewOpponentGrid(opponentGrid, cells, event) {
    const [col, row] = getCellCoordinate(opponentGrid, event);
    if (cells[row][col]) {
        cells[row][col].dataset.preview = "valid";
    }
    return [col, row];
}
export function clearPreviewOpponentGrid(col, row, cells) {
    if (cells[row][col]) {
        delete cells[row][col].dataset.preview;
    }
}
