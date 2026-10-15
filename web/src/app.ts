import { getGuestJWT, isTokenExpired } from "./auth.js";
import { enterMatchmaking, poolMatchmaking } from "./matchmaking.js";
import {
  attachWebSocketHandlers,
  close,
  send,
  upgradeWebsocket,
} from "./websocket.js";
import { createGrid, getCellCoordinate, handleOpponentGrid } from "./grid.js";
import {
  renderShipPlacement,
  shipPlacementClientValidation,
  shipPlacementServerValidation,
  validatePlacement,
  renderPreview,
  clearPreview,
  clearState,
} from "./ship.js";
import { DragSession, MatchState, Ship, WireMessage } from "./types.js";
import { defaultShips } from "./constants.js";

let path: string;
let jwt_token: string;
let matchState: MatchState = "PREPARATION";
let isMyTurn: boolean;
let playerOccupancy: HTMLElement[][];
let opponentOccupancy: HTMLElement[][];
let shipsIndex: number[][] = [];
let ships: Ship[];
let activeSession: DragSession | null = null;

const playerGrid = document.getElementById("player-grid") as HTMLDivElement;
const opponentGrid = document.getElementById("opponent-grid") as HTMLDivElement;

const gameState = document.getElementById("game-state")!;

const randomLink = document.getElementById("random-link") as HTMLElement;
const friendLink = document.getElementById("friend-link") as HTMLElement;
const haveFriend = document.getElementById("have-friend") as HTMLElement;
const urlLink = document.getElementById("url-link") as HTMLElement;

const matchButton = document.getElementById(
  "match-button",
) as HTMLButtonElement;
const prematchOverlay = document.getElementById(
  "prematch-overlay",
) as HTMLElement;
const leaveLink = document.getElementById("leave-link");

const shipIndicators = document.querySelectorAll(".ship-indicator");

let playerID: string;

// TODO: intiate the fleet -> buildOccupancy -> renderFleet
async function mainApp(): Promise<void> {
  try {
    if (window.location.pathname != "/") {
      handleFriendLink();
    }
    gameState.textContent = "Prepare your Ships";
    jwt_token = await getGuestJWT();
    playerOccupancy = createGrid("player-grid");
    opponentOccupancy = createGrid("opponent-grid");
    ships = defaultShips;
    shipsIndex = renderShipPlacement(playerOccupancy, ships, shipsIndex);
  } catch (error) {
    console.error("unexpected error: ", error);
    throw error;
  }
}

// lookup -> calculate grabOffset -> create active session -> capture -> dim origin ship
playerGrid.addEventListener("pointerdown", (event) => {
  if (!event.isPrimary) return;
  if (event.button != 0) return;

  //get actual grid coordinate
  const [col, row] = getCellCoordinate(playerGrid, event);
  const cell = playerOccupancy[row][col];
  const shipIndex = shipsIndex[row][col];
  if (!cell.dataset.state || shipIndex == -1) return;
  const ship = ships[shipIndex];
  const grabOffset =
    ship.direction == "Horizontal" ? col - ship.x : row - ship.y;

  activeSession = {
    pointerId: event.pointerId,
    shipIndex: shipIndex,
    grabOffset: grabOffset,
    oldHeadX: ship.x,
    oldHeadY: ship.y,
    hasMoved: false,
    isPlacementValid: false,
  };
  playerGrid.setPointerCapture(event.pointerId);
  clearState(
    activeSession.oldHeadX,
    activeSession.oldHeadY,
    ship,
    playerOccupancy,
  );
});

// calculate new head candidate -> validate placement -> clear preview frame before this(if any) -> renew attribute valid/invalid
playerGrid.addEventListener("pointermove", (event) => {
  //candidateHeadX -> head + current - start
  //or you can just current - grabOffset
  if (activeSession === null || !activeSession) return;
  const [col, row] = getCellCoordinate(playerGrid, event);
  if (activeSession.oldHeadX == col && activeSession.oldHeadY == row) {
    return;
  }
  activeSession.hasMoved = true;
  const ship = ships[activeSession.shipIndex];
  const candidateHeadX =
    ship.direction == "Horizontal" ? col - activeSession.grabOffset : col;
  const candidateHeadY =
    ship.direction == "Vertical" ? row - activeSession.grabOffset : row;
  const result = validatePlacement(
    candidateHeadX,
    candidateHeadY,
    ship,
    shipsIndex,
    activeSession.shipIndex,
  );
  clearPreview(
    activeSession.oldHeadX,
    activeSession.oldHeadY,
    ship,
    playerOccupancy,
  );
  renderPreview(candidateHeadX, candidateHeadY, ship, result, playerOccupancy);
  activeSession.oldHeadX = candidateHeadX;
  activeSession.oldHeadY = candidateHeadY;
  activeSession.isPlacementValid = result;
});

playerGrid.addEventListener("pointerup", (event) => {
  if (activeSession === null || !activeSession) return;

  const ship = ships[activeSession.shipIndex];

  if (!activeSession.hasMoved) {
    ship.direction = ship.direction == "Horizontal" ? "Vertical" : "Horizontal";
    if (
      validatePlacement(
        activeSession.oldHeadX,
        activeSession.oldHeadY,
        ship,
        shipsIndex,
        activeSession.shipIndex,
      )
    ) {
      activeSession.isPlacementValid = true;
    }
  }

  clearPreview(
    activeSession.oldHeadX,
    activeSession.oldHeadY,
    ship,
    playerOccupancy,
  );

  if (!activeSession.hasMoved && !activeSession.isPlacementValid) {
    ship.direction = ship.direction == "Horizontal" ? "Vertical" : "Horizontal";
  }

  if (activeSession.isPlacementValid) {
    ship.x = activeSession.oldHeadX;
    ship.y = activeSession.oldHeadY;
  }
  shipsIndex = renderShipPlacement(playerOccupancy, ships, shipsIndex);
  playerGrid.releasePointerCapture(event.pointerId);
  activeSession = null;
});

playerGrid.addEventListener("lostpointercapture", (event) => {
  playerGrid.releasePointerCapture(event.pointerId);
  activeSession = null;
});
playerGrid.addEventListener("pointercancel", (event) => {
  playerGrid.releasePointerCapture(event.pointerId);
  activeSession = null;
});

randomLink.addEventListener("click", (event) => {
  event.preventDefault();
  path = "";
  switchActiveLink(randomLink, friendLink);
  haveFriend.setAttribute("hidden", "");
  window.history.replaceState(null, "", "/");
});

friendLink.addEventListener("click", (event) => {
  event.preventDefault();
  handleFriendLink();
});

function handleFriendLink() {
  switchActiveLink(friendLink, randomLink);
  haveFriend?.removeAttribute("hidden");
  //random number -> attach to path param -> renew the url-link with same number

  let randPath: string;
  if (window.location.pathname == "/") {
    randPath = Math.random().toString(36).substring(2, 8).toUpperCase();
  } else {
    randPath = window.location.pathname.split("/")[1];
  }
  urlLink.textContent = `${window.location.origin}/${randPath}`;
  window.history.pushState({ path: randPath }, "", randPath);
  path = randPath;
}

function switchActiveLink(target: HTMLElement, other: HTMLElement) {
  target.classList.add("disabled");
  target.setAttribute("aria-disabled", "true");

  other.classList.remove("disabled");
  other.removeAttribute("aria-disabled");
}

matchButton!.addEventListener("click", () => handleMatchButton());

// TODO: before enterMatchmaking, should validate the ship first, client validation -> server validation/send -> render
export async function handleMatchButton() {
  if (!jwt_token || jwt_token == "") {
    console.warn("invalid token");
    return;
  }
  try {
    jwt_token = await getGuestJWT();
    matchButton.disabled = true;
    playerGrid.classList.add("disabled");

    shipPlacementClientValidation(playerOccupancy, ships, shipsIndex);
    await shipPlacementServerValidation(ships, jwt_token);
    shipsIndex = renderShipPlacement(playerOccupancy, ships, shipsIndex);

    await enterMatchmaking(jwt_token, path);
    gameState.textContent = "Matchmaking....";
    const ws_token = await poolMatchmaking(jwt_token);
    await upgradeWebsocket(ws_token);
    gameState.textContent = "Matched, waiting for opponent....";
    await attachWebSocketHandlers(HandleMessage);
  } catch (error) {
    matchButton.disabled = false;
    playerGrid.classList.remove("disabled");
    gameState.textContent = "Prepare your Ships";
    console.error(error);
    throw error;
  }
}

opponentGrid?.addEventListener("click", (event) =>
  handleOpponentGrid(event, isMyTurn, matchState),
);

leaveLink?.addEventListener("click", () => handleLeaveLink());

function handleLeaveLink() {
  if (matchState != "INBATTLE") return;
  if (confirm("are you sure to surrender??")) {
    send("CMD_SURRENDER");
  }
}

//TODO: add try-catch, and handle event_error
function HandleMessage(msg: WireMessage) {
  switch (msg.type) {
    case "EVENT_ROOM_CREATED":
      playerID = msg.data.player_id;
      console.log("Room Created");
      send("CMD_PLACE_SHIPS", defaultShips);
      return;
    case "EVENT_PLACE_SHIPMENT_RESULT":
      if (msg.data.is_valid) {
        shipsIndex = renderShipPlacement(
          playerOccupancy,
          defaultShips,
          shipsIndex,
        );
      }
      console.log(
        msg.data.is_valid
          ? "Ships Succesfully Placed"
          : `Invalid Ship Placements, Reason: ${msg.data.reason}`,
      );
      return;
    case "EVENT_PLAYER_READY":
      console.log(
        `${msg.data.player_id} is${msg.data.is_player_ready ? " Ready" : "n't Ready Yet"}`,
      );
      return;
    case "EVENT_BATTLE_START":
      prematchOverlay.style.display = "none";
      shipIndicators.forEach((elem) => {
        elem.classList.remove("hidden");
      });
      matchState = "INBATTLE";
      updateTurnState(msg.data.turn);
      gameState.textContent = `The Wheel of Fate is Turning!, ${isMyTurn ? "Your " : "Opponent "}Turn`;
      leaveLink?.removeAttribute("hidden");
      return;
    case "EVENT_FIRE_RESULT":
      // TODO: shloud display/render to player
      handleFireResult(msg.data);
      updateTurnState(msg.data.next_turn);
      gameState.textContent = isMyTurn ? "Your Turn" : "Opponent Turn";
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

function handleFireResult(data: any) {
  const targetGrid =
    data.target == playerID ? playerOccupancy : opponentOccupancy;
  const cell = targetGrid[data.y][data.x];
  cell.dataset.state = data.hit ? "hit" : "miss";
}

// need an turnUpdate
function updateTurnState(turnID: string) {
  if (playerID == turnID) {
    // isMyTurn = "Your turn";
    isMyTurn = true;
    opponentGrid.classList.remove("disabled");
  } else {
    // isMyTurn = "Opponent turn";
    isMyTurn = false;
    opponentGrid.classList.add("disabled");
  }
}

function handleGameOver(data: any) {
  let msg: string;
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
        matchState == "PREPARATION"
          ? "the enemy flees, try again?"
          : "draw, the enemy flees";
      // data.winner_id == playerID
      //   ? "you win !, the enemy flees"
      //   : "the enemy flees, but you lose, lmao";
      break;
    default:
      msg = "something strange happened, watchout !";
      break;
  }
  matchState = "FINISHED";
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
