import { directions } from "./constants";

export type MatchState = "PREPARATION" | "INBATTLE" | "FINISHED" | "ABORTED";

export interface WireMessage {
  type: string;
  data?: any;
}

export type Direction = (typeof directions)[number];

export type Ship = {
  name: string;
  type: ShipType;
  x: number;
  y: number;
  direction: Direction;
};

export type ShipType =
  | "aircraft_carrier"
  | "battleship"
  | "heavy_cruiser"
  | "light_cruiser"
  | "destroyer";

export type DragSession = {
  pointerId: number;
  shipIndex: number;
  grabOffset: number;
  oldHeadX: number;
  oldHeadY: number;
  hasMoved: boolean;
  isPlacementValid: boolean;
};

export type PreviewSession = {
  col: number;
  row: number;
};
