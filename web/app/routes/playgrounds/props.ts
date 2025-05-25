import { useState } from "react";
import type { PreparingPage } from "./component";

interface PreparingPageParams {
  timeout: Date;
  milliSecondPerTurn: number;
  boardWidth: number;
  mineCount: number;

  islandSectors: number[];
  enableDeploySectors: number[];
  enableMineSectors: number[];

  deployAction: (param: { deployAt: number; mines: number[] }) => void;
}

type PreparingPageSctState =
  | {
      state: "selectme";
      cursorSct: number;
    }
  | {
      state: "selectmines";
      cursorSct: number;
      selectedMeSct: number;
      selectedMineScts: number[];
    }
  | {
      state: "done";
      selectedMeSct: number;
      selectedMineScts: number[];
    };

export const usePreparingPageProps = (p: PreparingPageParams) => {
  const randSct = () =>
    p.enableDeploySectors[
      Math.floor(Math.random() * p.enableDeploySectors.length)
    ];
  const [sctState, setSctState] = useState<PreparingPageSctState>({
    state: "selectme",
    cursorSct: randSct(),
  });

  const directionAction = (dir: "North" | "West" | "South" | "East") => {
    switch (sctState.state) {
      case "selectme": {
        console.log("me move to", dir);
        break;
      }
      case "selectmines": {
        console.log("mines move to", dir);
        break;
      }
    }
  };
  const aAction = () => {
    if (sctState.state === "done") {
      return;
    }
    setSctState((prev) => {
      switch (prev.state) {
        case "selectme": {
          if (!p.enableDeploySectors.includes(prev.cursorSct)) {
            return prev;
          }
          return {
            state: "selectmines",
            cursorSct: randSct(),
            selectedMeSct: prev.cursorSct,
            selectedMineScts: [],
          };
        }
        case "selectmines": {
          if (!p.enableMineSectors.includes(prev.cursorSct)) {
            return prev;
          }
          if (prev.selectedMineScts.length >= p.mineCount) {
            return prev;
          }
          if (prev.selectedMineScts.length + 1 !== p.mineCount) {
            return {
              state: "selectmines",
              cursorSct: randSct(),
              selectedMeSct: prev.selectedMeSct,
              selectedMineScts: [prev.cursorSct, ...prev.selectedMineScts],
            };
          }
          // 最後の一つ state更新,deploy
          p.deployAction({
            deployAt: prev.selectedMeSct,
            mines: [...prev.selectedMineScts, prev.cursorSct],
          });
          return {
            state: "done",
            selectedMeSct: prev.selectedMeSct,
            selectedMineScts: [...prev.selectedMineScts, prev.cursorSct],
          };
        }
      }
      return prev;
    });
  };
  const bAction = () => {
    if (sctState.state === "done") {
      return;
    }
    setSctState((prev) => {
      switch (prev.state) {
        case "selectme": {
          return {
            state: "selectme",
            cursorSct: randSct(),
          };
        }
        case "selectmines": {
          if (!prev.selectedMineScts.length) {
            return {
              state: "selectme",
              cursorSct: prev.selectedMeSct,
            };
          }
          const [last, ...scts] = prev.selectedMineScts;
          return {
            state: "selectmines",
            cursorSct: last,
            selectedMeSct: prev.selectedMeSct,
            selectedMineScts: scts,
          };
        }
      }
      return prev;
    });
  };

  const displayMeSectors = (): number[] => {
    switch (sctState.state) {
      case "selectme": {
        return [sctState.cursorSct];
      }
      default: {
        return [sctState.selectedMeSct];
      }
    }
  };
  const displayMineSectors = (): number[] => {
    switch (sctState.state) {
      case "selectme": {
        return [];
      }
      case "selectmines": {
        return [...sctState.selectedMineScts, sctState.cursorSct];
      }
      case "done":
        return sctState.selectedMineScts;
    }
  };

  const props: React.ComponentProps<typeof PreparingPage> = {
    preparingStep:
      sctState.state === "selectme"
        ? "me"
        : sctState.state === "selectmines"
          ? "mines"
          : "done",
    Leading: {
      Alerm: {
        useAlerm: true,
        startPingMilliSec: p.milliSecondPerTurn,
        finishDatetime: p.timeout,
      },
    },
    Ocean: {
      sectorCount: p.boardWidth,
      islands: p.islandSectors,
    },
    MeOcean: {
      sectors: displayMeSectors(),
    },
    MinesOcean: {
      sectors: displayMineSectors(),
    },
    Controller: {
      Direction: {
        North: {
          onClick: () => directionAction("North"),
        },
        West: {
          onClick: () => directionAction("West"),
        },
        South: {
          onClick: () => directionAction("South"),
        },
        East: {
          onClick: () => directionAction("East"),
        },
      },
      AButton: {
        onClick: sctState.state === "done" ? undefined : aAction,
      },
      BButton: {
        onClick: sctState.state === "done" ? undefined : bAction,
      },
    },
  };

  return props;
};
