## Matchmaking w/ REST
### 1. Matchmaking (Queue)

The player requests to be put in the matchmaking queue.

    Request: POST /api/v1/matchmaking

    Headers: Authorization: Bearer <user_token>

    Response: 202 Accepted
    JSON Payload:

    {
      "status": "queued",
      "message": "Searching for an opponent..."
    }

### 2. Check Matchmaking Status (Polling / SSE)

The client checks periodically whether the opponent has been found.

    Request: GET /api/v1/matchmaking/status

    Response (if not found): 200 OK
    JSON Payload:

    {
      "status": "searching",
      "room_id": null,
      "ws_token": null
    }

    Response (if found an opponent): 200 OK
    JSON Payload:

    {
      "status": "matched",
      "room_id": "67",
      "ws_url": "ws://gameserver:8080/room/67",
      "ws_token": "token_p1_abc123" 
    }

## Battle-Phase w/ WebSocket Protocol

Once the WebSocket connection is established using the `ws_url` and `ws_token` from the matchmaking status, communication is handled via Query Parameter at Handshake usign One-Time Token .

### 1. Place Ships (Setup Phase)

*Note: server-side must validate accepted ship type and max ship for each player*

*Accepted Ship Type: {"heavy_cruiser",3}, {"aircraft_carrier",5}, {"light_cruiser",3}, {"destroyer",2}, {"battleship",4}*

*Max Placement: 5*

*Grid Boundaries: using 10x10 grid w/ 0-indexed *

The player submits their ship placements to the game server.

    Direction: Client to Server (C2S)
    Type: CMD_PLACE_SHIPS

    JSON Payload:

    {
      "type": "CMD_PLACE_SHIPS",
      "payload": {
        "ships": [
          {
            "name": "IJN_Taihou",
            "type": "aircraft_carrier",
            "x": 1,
            "y": 3,
            "direction": "horizontal"
          },
          {
            "name": "KMS_Prinz_Eugen",
            "type": "heavy_cruiser",
            "x": 3,
            "y": 1,
            "direction": "vertical"
          }
        ]
      }
    }

The server responds indicating whether the placement is valid.

    Direction: Server to Client (S2C)
    Type: EVT_PLACE_SHIPS_RESULT

    JSON Payload:

    {
      "type": "EVT_PLACE_SHIPS_RESULT",
      "payload": {
        "is_valid": true,
        "reason": ""
      }
    }

The server notifies the other player when their opponent has successfully placed their ships.

    Direction: Server to Client (S2C)
    Type: EVT_OPPONENT_READY

    JSON Payload:

    {
      "type": "EVT_OPPONENT_READY",
      "payload": {}
    }

### 2. Battle Start

The server notifies both players that the setup phase is complete if the both player is ready, the game has started, and indicates who takes the first turn.

    Direction: Server to Client (S2C)
    Type: EVT_BATTLE_START

    JSON Payload:

    {
      "type": "EVT_BATTLE_START",
      "payload": {
        "first_turn": "666"
      }
    }

### 3. Fire (Combat Phase)

The active player fires at a coordinate on the opponent's grid.

    Direction: Client to Server (C2S)
    Type: CMD_FIRE

    JSON Payload:

    {
      "type": "CMD_FIRE",
      "payload": {
        "x": 4,
        "y": 5
      }
    }

The server returns the result of the shot.

*ship_destroyed only return name or type of the ship if the ship completely destroyed, if its not->null*

    Direction: Server to Client (S2C)
    Type: EVT_FIRE_RESULT

    JSON Payload:

    {
      "type": "EVT_FIRE_RESULT",
      "payload": {
        "fired_by": "666",
        "target": "777",
        "x": 4,
        "y": 5,
        "hit": true,
        "ship_destroyed": null,
        "next_turn": "666"
      }
    }

One of the active player surrend in the middle of battle

    Direction: Client to Server (C2S)
    Type: CMD_SURRENDER

    JSON Payload:

    {
      "type": "CMD_SURRENDER",
      "payload": {
      }
    }

### 4. Error Handling

The server sends an error event if the client performs an invalid move (e.g., shooting the same coordinates twice).

    Direction: Server to Client (S2C)
    Type: EVT_ERROR

    JSON Payload:

    {
      "type": "EVT_ERROR",
      "payload": {
        "code": "ALREADY_SHOT",
        "message": "you've already shot the coordinates (4,5) before."
      }
    }
*Error Code List: "ALREADY_SHOT", "NOT_YOUR_TURN", "INVALID_PLACEMENT", "INVALID_PHASE", "UNKNOWN_ERROR"*

### 5. Reconnection Sync

*Note: This event can be implemented later to handle client reconnections.*

The server sends the current state of the game to sync the client upon a successful reconnection.

    Direction: Server to Client (S2C)
    Type: EVT_RECONNECT_SYNC

    JSON Payload:

    {
      "type": "EVT_RECONNECT_SYNC",
      "payload": {
        "your_ships": [
          {
            "name": "KMS_Prinz_Eugen",
            "type": "heavy_cruiser",
            "x": 3,
            "y": 1,
            "direction": "vertical"
          }
        ],
        "your_shots_history": [
          {
            "x": 4,
            "y": 5,
            "hit": true
          },
          {
            "x": 5,
            "y": 4,
            "hit": false
          }
        ],
        "enemy_shots_history": [
          {
            "x": 6,
            "y": 7,
            "hit": false
          }
        ],
        "next_turn": "777"
      }
    }

### 6. Game Over

The server notifies both players that the game has ended.

    Direction: Server to Client (S2C)
    Type: EVT_GAME_OVER

    JSON Payload:

    {
      "type": "EVT_GAME_OVER",
      "payload": {
        "status": "finished",
        "winner_id": "666",
        "reason": "all_enemy_ships_destroyed"
      }
    }

    * Possible values for "status": "finished", "aborted"
    * Possible values for "reason": "all_enemy_ships_destroyed", "opponent_disconnected_timeout", "opponent_surrendered"
