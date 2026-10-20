## Matchmaking w/ REST
### 1. Register Guest (Auth)

A player registers as a temporary guest to acquire a JWT token.

    Request: POST /api/v1/auth/guest

    Response: 201 Created
    JSON Payload:

    {
      "status": "token created succesfully",
      "data": {
        "token": "<guest_jwt_token>"
      }
    }

### 2. Validate Ship Placements (Pre-match Validation)

The client optionally validates their fleet placement on the server before queuing.

    Request: QUERY /api/v1/validate

    Headers: Authorization: Bearer <user_token>

    JSON Payload (Array of ships):

    [
      {
        "name": "IJN_Taihou",
        "type": "aircraft_carrier",
        "x": 5,
        "y": 9,
        "direction": "Horizontal"
      },
      {
        "name": "HMS_Belfast",
        "type": "light_cruiser",
        "x": 0,
        "y": 0,
        "direction": "Vertical"
      }
    ]

    Response (Valid): 200 OK
    JSON Payload:

    {
      "status": "ship validated"
    }

    Response (Invalid): 400 Bad Request
    JSON Payload:

    {
      "status": "error",
      "error": "invalid ship position"
    }

### 3. Matchmaking (Queue)

The player requests to enter the matchmaking queue. For random matchmaking, the body is empty string `""`. For custom room matching with a friend, the body contains the room ID string.

    Request: POST /api/v1/matchmaking

    Headers: Authorization: Bearer <user_token>

    JSON Payload:

    ""

    or for custom room:

    "ROOM123"

    Response: 201 Created
    JSON Payload:

    {
      "status": "queued",
      "data": {
        "message": "Searching for an opponent"
      }
    }

### 4. Check Matchmaking Status (Polling)

The client checks periodically whether an opponent has been found.

    Request: GET /api/v1/matchmaking/status

    Headers: Authorization: Bearer <user_token>

    Response (Searching): 200 OK
    JSON Payload:

    {
      "status": "searching"
    }

    Response (Matched): 200 OK
    JSON Payload:

    {
      "status": "matched",
      "data": {
        "ws_token": "token_uuid_123"
      }
    }

    Response (Timed Out): 200 OK
    JSON Payload:

    {
      "status": "timedout"
    }

## Battle-Phase w/ WebSocket Protocol

Once matched, the client upgrades to WebSocket at `/ws?ticket=<ws_token>`. The server validates and consumes the one-time ticket atomically during the HTTP 101 Switching Protocols upgrade.

> **Note on Message Envelope**: All WebSocket messages (Client-to-Server and Server-to-Client) follow the standard envelope format:
> ```json
> {
>   "type": "<COMMAND_OR_EVENT_NAME>",
>   "data": { ... }
> }
> ```

### 1. Room Created

When both players connect to the room, the server notifies each player of their assigned player ID.

    Direction: Server to Client (S2C)
    Type: EVENT_ROOM_CREATED

    JSON Payload:

    {
      "type": "EVENT_ROOM_CREATED",
      "data": {
        "player_id": "player_uuid_abc"
      }
    }

### 2. Place Ships & Ready State (Setup Phase)

*Note: Server-side validates accepted ship type and count.*

*Accepted Ship Types: `aircraft_carrier` (5), `battleship` (4), `heavy_cruiser` (3), `light_cruiser` (3), `destroyer` (2)*

*Fleet Count: Exactly 5 ships*

*Grid Boundaries: 10x10 grid (0-indexed, x: 0..9, y: 0..9)*

*Direction values: `"Horizontal"`, `"Vertical"`*

The player submits their ship placements to the game server.

    Direction: Client to Server (C2S)
    Type: CMD_PLACE_SHIPS

    JSON Payload:

    {
      "type": "CMD_PLACE_SHIPS",
      "data": [
        {
          "name": "IJN_Taihou",
          "type": "aircraft_carrier",
          "x": 5,
          "y": 9,
          "direction": "Horizontal"
        },
        {
          "name": "HMS_Belfast",
          "type": "light_cruiser",
          "x": 0,
          "y": 0,
          "direction": "Vertical"
        },
        {
          "name": "KMS_Bismarck",
          "type": "battleship",
          "x": 6,
          "y": 0,
          "direction": "Horizontal"
        },
        {
          "name": "KMS_Prinz_Eugen",
          "type": "heavy_cruiser",
          "x": 0,
          "y": 7,
          "direction": "Vertical"
        },
        {
          "name": "USS_Laffey",
          "type": "destroyer",
          "x": 5,
          "y": 5,
          "direction": "Horizontal"
        }
      ]
    }

The server responds to the sender indicating whether the placement is valid.

    Direction: Server to Client (S2C)
    Type: EVENT_PLACE_SHIPMENT_RESULT

    JSON Payload:

    {
      "type": "EVENT_PLACE_SHIPMENT_RESULT",
      "data": {
        "is_valid": true,
        "reason": "its true ofc"
      }
    }

The server notifies players (broadcast to both) when a player's ready state changes (after placing ships or canceling ready).

    Direction: Server to Client (S2C)
    Type: EVENT_PLAYER_READY

    JSON Payload:

    {
      "type": "EVENT_PLAYER_READY",
      "data": {
        "player_id": "player_uuid_abc",
        "is_player_ready": true
      }
    }

A player can cancel their ready status during the setup phase:

    Direction: Client to Server (C2S)
    Type: CMD_CANCEL_READY

    JSON Payload:

    {
      "type": "CMD_CANCEL_READY",
      "data": {}
    }

### 3. Battle Start

The server notifies both players that the setup phase is complete (both players ready), the battle has started, and indicates who takes the turn.

    Direction: Server to Client (S2C)
    Type: EVENT_BATTLE_START

    JSON Payload:

    {
      "type": "EVENT_BATTLE_START",
      "data": {
        "turn": "player_uuid_abc"
      }
    }

### 4. Fire (Combat Phase)

The active player fires at a coordinate on the opponent's grid.

    Direction: Client to Server (C2S)
    Type: CMD_FIRE

    JSON Payload:

    {
      "type": "CMD_FIRE",
      "data": {
        "x": 4,
        "y": 5
      }
    }

The server broadcasts the result of the shot to both players.

*`ship_destroyed` only returns the name of the ship if the ship was completely destroyed, otherwise omitted or empty.*

    Direction: Server to Client (S2C)
    Type: EVENT_FIRE_RESULT

    JSON Payload:

    {
      "type": "EVENT_FIRE_RESULT",
      "data": {
        "fired_by": "player_uuid_abc",
        "target": "player_uuid_xyz",
        "x": 4,
        "y": 5,
        "hit": true,
        "ship_destroyed": "USS_Laffey",
        "next_turn": "player_uuid_abc"
      }
    }

One of the active players surrenders in the middle of battle:

    Direction: Client to Server (C2S)
    Type: CMD_SURRENDER

    JSON Payload:

    {
      "type": "CMD_SURRENDER",
      "data": {}
    }

### 5. Error Handling

The server sends an error event if the client performs an invalid move or command.

    Direction: Server to Client (S2C)
    Type: EVENT_ERROR

    JSON Payload:

    {
      "type": "EVENT_ERROR",
      "data": "ALREADY_SHOT"
    }

*Error strings derived from internal game error states: `"ALREADY_SHOT"`, `"NOT_YOUR_TURN"`, `"INVALID_PLACEMENT"`, `"INVALID_PHASE"`, `"UNKNOWN_ERROR"`, `"Invalid Type Command"`.*

### 6. Game Over

The server notifies players that the match has concluded.

    Direction: Server to Client (S2C)
    Type: EVENT_GAME_OVER

    JSON Payload:

    {
      "type": "EVENT_GAME_OVER",
      "data": {
        "status": "finished",
        "winner_id": "player_uuid_abc",
        "reason": "ALL_ENEMY_SHIP_DESTROYED"
      }
    }

* Possible values for `status`: `"finished"`, `"aborted"`
* Possible values for `reason`: `"ALL_ENEMY_SHIP_DESTROYED"`, `"OPPONENT_DISCONNECTED_TIMEOUT"`, `"OPPONENT_SURRENDERED"`
