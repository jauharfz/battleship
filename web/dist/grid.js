var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
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
export function handleOpponentGrid(event) {
    return __awaiter(this, void 0, void 0, function* () {
        const cell = event.target.closest("[data-y][data-y]");
        if (!cell)
            return;
        if (!(cell instanceof HTMLElement))
            return;
        send("CMD_FIRE", { x: Number(cell.dataset.x), y: Number(cell.dataset.y) });
    });
}
