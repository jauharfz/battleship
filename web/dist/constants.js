export const shipSizes = {
    aircraft_carrier: 5,
    battleship: 4,
    heavy_cruiser: 3,
    light_cruiser: 3,
    destroyer: 2,
};
export const directions = ["Horizontal", "Vertical"];
export const defaultShips = [
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
    {
        name: "IJN_Taihou",
        type: "aircraft_carrier",
        x: 8,
        y: 2,
        direction: "Vertical",
    },
    {
        name: "HMS_Belfast",
        type: "light_cruiser",
        x: 1,
        y: 5,
        direction: "Horizontal",
    },
    {
        name: "KMS_Bismarck",
        type: "battleship",
        x: 2,
        y: 0,
        direction: "Vertical",
    },
    {
        name: "KMS_Prinz_Eugen",
        type: "heavy_cruiser",
        x: 4,
        y: 7,
        direction: "Horizontal",
    },
    {
        name: "USS_Laffey",
        type: "destroyer",
        x: 4,
        y: 2,
        direction: "Vertical",
    },
];
export const GUEST_TOKEN_KEY = "guest_jwt";
export const SHIPS_LAYOUT_KEY = "ships";
