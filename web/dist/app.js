var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import { getGuestJWT } from "./auth.js";
import { enterMatchmaking, poolMatchmaking } from "./matchmaking.js";
import { attachWebSocketHandlers, close, send, upgradeWebsocket, } from "./websocket.js";
import { createGrid, handleOpponentGrid } from "./grid.js";
import { renderShipPlacement, renderShipPlacement as initialShipPlacements, shipPlacementClientValidation, shipPlacementServerValidation, } from "./ship.js";
let jwt_token;
let matchState = "PREPARATION";
let turn;
let player_grid;
let opponent_grid;
const gameState = document.getElementById("game-state");
const matchButton = document.getElementById("match-button");
const prematchOverlay = document.getElementById("prematch-overlay");
const leaveLink = document.getElementById("leave-link");
const shipIndicators = document.querySelectorAll(".ship-indicator");
const playerGrid = document.getElementById("player-grid");
const opponentGrid = document.getElementById("opponent-grid");
let playerID;
// createGrid  : initial render
// shipPlacements : PrepState <- initial state just after upgrading ws
// updateCell : InBattleState
// clearGrid : after Finished/Aborted-STATE
function mainApp() {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            gameState.textContent = "Prepare your Ships";
            jwt_token = yield getGuestJWT();
            player_grid = createGrid("player-grid");
            opponent_grid = createGrid("opponent-grid");
            renderShipPlacement(player_grid, defaultShips);
        }
        catch (error) {
            console.error("unexpected error: ", error);
            throw error;
        }
    });
}
const defaultShips = [
    {
        name: "IJN_Taihou",
        type: "aircraft_carrier",
        x: 5,
        y: 9,
        direction: "Horizontal",
    },
    {
        name: "HMS_Belfast",
        type: "light_cruiser",
        x: 0,
        y: 0,
        direction: "Vertical",
    },
    {
        name: "KMS_Bismarck",
        type: "battleship",
        x: 6,
        y: 0,
        direction: "Horizontal",
    },
    {
        name: "KMS_Prinz_Eugen",
        type: "heavy_cruiser",
        x: 0,
        y: 7,
        direction: "Vertical",
    },
    {
        name: "USS_Laffey",
        type: "destroyer",
        x: 5,
        y: 5,
        direction: "Horizontal",
    },
];
matchButton.addEventListener("click", () => handleMatchButton());
// TODO: before enterMatchmaking, should validate the ship first, client validation -> server validation/send -> render
export function handleMatchButton() {
    return __awaiter(this, void 0, void 0, function* () {
        if (!jwt_token || jwt_token == "") {
            console.warn("invalid token");
            return;
        }
        try {
            matchButton.disabled = true;
            // unfortunately bcs dev having some skill issue, we can using defaultShips temporarily
            yield shipPlacementClientValidation(player_grid, defaultShips);
            yield shipPlacementServerValidation(defaultShips, jwt_token);
            renderShipPlacement(player_grid, defaultShips);
            yield enterMatchmaking(jwt_token);
            gameState.textContent = "Matchmaking....";
            const ws_token = yield poolMatchmaking(jwt_token);
            yield upgradeWebsocket(ws_token);
            yield attachWebSocketHandlers(HandleMessage);
            send("CMD_PLACE_SHIPS", defaultShips);
        }
        catch (error) {
            matchButton.disabled = false;
            gameState.textContent = "Prepare your Ships";
            console.error(error);
            throw error;
        }
    });
}
opponentGrid === null || opponentGrid === void 0 ? void 0 : opponentGrid.addEventListener("click", (event) => handleOpponentGrid(event));
leaveLink === null || leaveLink === void 0 ? void 0 : leaveLink.addEventListener("click", () => handleLeaveLink());
function handleLeaveLink() {
    if (matchState != "INBATTLE")
        return;
    if (confirm("are you sure little human??")) {
        send("CMD_SURRENDER");
    }
}
//TODO: add try-catch
function HandleMessage(msg) {
    switch (msg.type) {
        case "EVENT_ROOM_CREATED":
            playerID = msg.data.player_id;
            console.log("Room Created");
            send("CMD_PLACE_SHIPS", defaultShips);
            return;
        case "EVENT_PLACE_SHIPMENT_RESULT":
            if (msg.data.is_valid) {
                initialShipPlacements(player_grid, defaultShips);
            }
            console.log(msg.data.is_valid
                ? "Ships Succesfully Placed"
                : `Invalid Ship Placements, Reason: ${msg.data.reason}`);
            return;
        case "EVENT_PLAYER_READY":
            console.log(`${msg.data.player_id} is${msg.data.is_player_ready ? " Ready" : "n't Ready Yet"}`);
            return;
        case "EVENT_BATTLE_START":
            prematchOverlay.style.display = "none";
            shipIndicators.forEach((elem) => {
                elem.classList.remove("hidden");
            });
            matchState = "INBATTLE";
            updateTurnState(msg.data.turn);
            gameState.textContent = `The Wheel of Fate is Turning!, ${turn}`;
            leaveLink === null || leaveLink === void 0 ? void 0 : leaveLink.removeAttribute("hidden");
            return;
        case "EVENT_FIRE_RESULT":
            // TODO: shloud display/render to player
            handleFireResult(msg.data);
            updateTurnState(msg.data.next_turn);
            gameState.textContent = turn;
            return;
        case "EVENT_ERROR":
            return `${msg.data}`;
        case "EVENT_GAME_OVER":
            gameState.textContent = "Game Over";
            handleGameOver(msg.data);
            return;
        default:
            return `invalid type msg: ${msg}`;
    }
}
function handleFireResult(data) {
    const targetGrid = data.target == playerID ? player_grid : opponent_grid;
    const cell = targetGrid[data.y][data.x];
    cell.dataset.state = data.hit ? "hit" : "miss";
}
// need an turnUpdate
function updateTurnState(turnID) {
    if (playerID == turnID) {
        turn = "Your turn";
        playerGrid.classList.add("disabled");
        opponentGrid.classList.remove("disabled");
    }
    else {
        turn = "Opponent turn";
        playerGrid.classList.remove("disabled");
        opponentGrid.classList.add("disabled");
    }
}
function handleGameOver(data) {
    matchState = "FINISHED";
    let msg;
    switch (data.reason) {
        case "ALL_ENEMY_SHIP_DESTROYED":
            msg = data.winner_id == playerID ? "you win !" : "nice try !";
            break;
        case "OPPONENT_SURRENDERED":
            msg =
                data.winner_id == playerID
                    ? "You Win !, the enemy has a massive skill issue(they leave lmao)"
                    : "looks like skill issue, try again?";
            break;
        case "OPPONENT_DISCONNECTED_TIMEOUT":
            msg =
                data.winner_id == playerID
                    ? "you win !, the enemy flees"
                    : "the enemy flees, but you lose, lmao";
            break;
        default:
            msg = "something strange happened, watchout !";
            break;
    }
    alert(msg);
    close();
    window.location.reload();
}
mainApp();
//save the client id -> if == unlock, otherwise lock
// (window as any).game = {
//   placeShips: (ships: Ship[] = defaultShips) => {
//     send("CMD_PLACE_SHIPS", ships);
//   },
//   cancelReady: () => {
//     send("CMD_CANCEL_READY", {});
//   },
//   fire: (x: number = 8, y: number = 9) => {
//     send("CMD_FIRE", { x, y });
//   },
//   surrender: () => {
//     send("CMD_SURRENDER", {});
//   },
// };
//
//
// TODO: create cancel button
