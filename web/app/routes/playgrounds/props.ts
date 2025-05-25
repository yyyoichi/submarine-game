import { useState } from "react";
import type { PlayingPage, PreparingPage } from "./component";

type InputDirection = "North" | "West" | "South" | "East";

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

  const directionAction = (dir: InputDirection) => {
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

type PlayingActionType = "move" | "fire-torpedo" | "trigger-mine";

interface PlayingPageParams {
  timeout: Date;
  milliSecondPerTurn: number;
  boardWidth: number;
  turn: number;
  inMyTurn: boolean;

  at: number;
  islandSectors: number[];
  enableMoveSectors: number[];
  enableTorpedoSectors: number[];
  enableMineSectors: number[];

  prevAction:
    | {
        mode: "hide";
        me: boolean;
        type: "deploy";
      }
    | {
        mode: "hide";
        me: boolean;
        type: "move";
        direction: "North" | "South" | "East" | "West";
      }
    | {
        mode: "show";
        me: boolean;
        type: "fire-torpedo" | "trigger-mine";
        at: number;
        result: "full-speed-ahead" | "hard-to-starboard" | "hit";
      };

  action: (param: {
    at: number;
    type: PlayingActionType;
  }) => void;
}

type PlayingPageActionState =
  | {
      state: "actiontype";
      selectedActType: PlayingActionType;
    }
  | {
      state: "actionat";
      selectedActType: PlayingActionType;
      selectedSct: number;
    }
  | {
      state: "done";
      selectedActType: PlayingActionType;
      selectedSct: number;
    };

export const usePlayingPageProps = (p: PlayingPageParams) => {
  const randSct = (scts: number[]) =>
    scts[Math.floor(Math.random() * scts.length)];
  const [actState, setActState] = useState<PlayingPageActionState>({
    state: "actiontype",
    selectedActType: "move",
  });

  const directionAction = (dir: InputDirection) => {
    if (!p.inMyTurn) {
      return;
    }
    setActState((prev) => {
      if (prev.state === "done") return prev;
      if (actState.state === "actiontype") {
        const uses: Record<
          InputDirection,
          Record<PlayingActionType, PlayingActionType>
        > = {
          North: {
            move: "trigger-mine",
            "fire-torpedo": "move",
            "trigger-mine": "fire-torpedo",
          },
          South: {
            move: "fire-torpedo",
            "fire-torpedo": "trigger-mine",
            "trigger-mine": "move",
          },
          East: {
            move: "move",
            "fire-torpedo": "fire-torpedo",
            "trigger-mine": "trigger-mine",
          },
          West: {
            move: "move",
            "fire-torpedo": "fire-torpedo",
            "trigger-mine": "trigger-mine",
          },
        };
        return {
          state: "actiontype",
          selectedActType: uses[dir][prev.selectedActType],
        };
      }
      // actionat state
      // TODO
      console.log("move to", dir);
      return prev;
    });
  };
  const aAction = () => {
    if (!p.inMyTurn) {
      return;
    }
    setActState((prev) => {
      if (prev.state === "done") return prev;
      if (prev.state === "actiontype") {
        return {
          state: "actionat",
          selectedActType: prev.selectedActType,
          selectedSct:
            prev.selectedActType === "move"
              ? randSct(p.enableMoveSectors)
              : prev.selectedActType === "fire-torpedo"
                ? randSct(p.enableTorpedoSectors)
                : randSct(p.enableMineSectors),
        };
      }
      p.action({
        at: prev.selectedSct,
        type: prev.selectedActType,
      });
      return {
        state: "done",
        selectedActType: prev.selectedActType,
        selectedSct: prev.selectedSct,
      };
    });
  };
  const bAction = () => {
    if (!p.inMyTurn) {
      return;
    }
    setActState((prev) => {
      if (prev.state === "done") return prev;
      if (prev.state === "actiontype") {
        return prev;
      }
      // actionat state
      return {
        state: "actiontype",
        selectedActType: prev.selectedActType,
      };
    });
  };

  const enableAction = p.inMyTurn && actState.state !== "done";
  const props: React.ComponentProps<typeof PlayingPage> = {
    turn: p.turn,
    Leading: {
      inMyTrun: p.inMyTurn,
      Alerm: {
        useAlerm: p.inMyTurn,
        startPingMilliSec: p.milliSecondPerTurn,
        finishDatetime: p.timeout,
      },
      children: p.inMyTurn
        ? "敵艦A1に魚雷着弾！ヨーソロー！A1付近の敵艦を攻撃せよ！"
        : "A3に魚雷、反響音を確認！",
    },
    Ocean: {
      sectorCount: p.boardWidth,
      islands: p.islandSectors,
      me: p.at,
      enables:
        actState.state === "actiontype"
          ? [...Array(p.boardWidth ** 2)].reduce((acc, x_, i) => {
              if (p.islandSectors.includes(i)) {
                return acc;
              }
              acc.push(i);
              return acc;
            }, [] as number[])
          : actState.selectedActType === "move"
            ? p.enableMoveSectors
            : actState.selectedActType === "fire-torpedo"
              ? p.enableTorpedoSectors
              : p.enableMineSectors,
      sector:
        actState.state === "actiontype" ? undefined : actState.selectedSct,
      type:
        actState.selectedActType === "move"
          ? "move"
          : actState.selectedActType === "fire-torpedo"
            ? "torpedo"
            : "mine",
    },
    Controller: {
      Direction: {
        North: {
          onClick: enableAction ? () => directionAction("North") : undefined,
        },
        South: {
          onClick: enableAction ? () => directionAction("South") : undefined,
        },
        East: {
          onClick: enableAction ? () => directionAction("East") : undefined,
        },
        West: {
          onClick: enableAction ? () => directionAction("West") : undefined,
        },
      },
      AButton: {
        onClick: enableAction ? aAction : undefined,
      },
      BButton: {
        onClick: enableAction ? bAction : undefined,
      },
      CommandWindow: {
        commands: ["move", "fire-torpedo", "trigger-mine"],
        use: actState.selectedActType,
      },
    },
    OverlayedLeading: {
      children: "C1に敵艦魚雷着弾！ヨーソロー！",
    },
    OverlayedOcean: {
      gapSector: p.prevAction.mode === "hide" ? undefined : p.prevAction.at,
      ringSector: p.prevAction.mode === "hide" ? undefined : p.prevAction.at,
      title:
        p.prevAction.type === "deploy"
          ? "作戦開始"
          : p.prevAction.type === "fire-torpedo"
            ? "魚雷攻撃"
            : p.prevAction.type === "trigger-mine"
              ? "機雷作動"
              : // move
                p.prevAction.type === "move" &&
                  p.prevAction.direction === "North"
                ? "北に潜航"
                : p.prevAction.type === "move" &&
                    p.prevAction.direction === "South"
                  ? "南に潜航"
                  : p.prevAction.type === "move" &&
                      p.prevAction.direction === "East"
                    ? "東に潜航"
                    : p.prevAction.type === "move" &&
                        p.prevAction.direction === "West"
                      ? "西に潜航"
                      : "",
      titlePosition: p.prevAction.me ? "right" : "left",
    },
  };

  return props;
};
