var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export const shipSizes = new Map([
    ["aircraft_carrier", 5],
    ["battleship", 4],
    ["heavy_cruiser", 3],
    ["light_cruiser", 3],
    ["destroyer", 2],
]);
// there's need validation each ship
// and each ship need iterates the length/sizes ship itself
// okey, just focus create client validation for each ship first
export function shipPlacementServerValidation(ships, token) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const res = yield fetch("/api/v1/validate", {
                method: "QUERY",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(ships),
            });
            const payload = yield res.json();
            if (payload.status == "error") {
                throw new Error(payload.error);
            }
            return;
        }
        catch (error) {
            throw error;
        }
    });
}
export function renderShipPlacement(cells, ships) {
    traverseShips(cells, ships, true);
}
export function shipPlacementClientValidation(cells, ships) {
    return traverseShips(cells, ships, false);
}
// TODO: add another validation for outbound
function traverseShips(cells, ships, isValid) {
    for (let ship of ships) {
        for (let i = 0; i < shipSizes.get(ship.type); i++) {
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
