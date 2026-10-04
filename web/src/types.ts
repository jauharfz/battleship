export type MatchState = "PREPARATION" | "INBATTLE" | "FINISHED" | "ABORTED";

export interface WireMessage {
  type: string;
  data?: any;
}

export type Ship = {
  name: string;
  type: ShipType;
  x: number;
  y: number;
  direction: string;
};

export type ShipType =
  | "aircraft_carrier"
  | "battleship"
  | "heavy_cruiser"
  | "light_cruiser"
  | "destroyer";

export type FireData = {
  x: number;
  y: number;
};
