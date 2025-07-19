import { SectorDirection, sector } from "@/lib/utils";
import type React from "react";
import { useState } from "react";
import type { FinishedPage, PlayingPage, PreparingPage } from "./component";

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
  loading: boolean;
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
    if (p.loading) {
      return;
    }
    switch (sctState.state) {
      case "selectme": {
        const d = new SectorDirection(p.boardWidth);
        const next =
          dir === "North"
            ? d.north(sctState.cursorSct, p.enableDeploySectors)
            : dir === "South"
              ? d.south(sctState.cursorSct, p.enableDeploySectors)
              : dir === "East"
                ? d.east(sctState.cursorSct, p.enableDeploySectors)
                : d.west(sctState.cursorSct, p.enableDeploySectors);
        if (next === -1) {
          return;
        }
        console.log("me move to", next);
        setSctState((prev) => ({
          ...prev,
          cursorSct: next,
        }));
        return;
      }
      case "selectmines": {
        const d = new SectorDirection(p.boardWidth);
        const next =
          dir === "North"
            ? d.north(sctState.cursorSct, p.enableDeploySectors)
            : dir === "South"
              ? d.south(sctState.cursorSct, p.enableDeploySectors)
              : dir === "East"
                ? d.east(sctState.cursorSct, p.enableDeploySectors)
                : d.west(sctState.cursorSct, p.enableDeploySectors);
        if (next === -1) {
          return;
        }
        console.log("mine move to", next);
        setSctState((prev) => ({
          ...prev,
          cursorSct: next,
        }));
        return;
      }
    }
  };
  const aAction = () => {
    if (p.loading) {
      return;
    }
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
    if (p.loading) {
      return;
    }
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
  const enableAction = !p.loading && sctState.state !== "done";
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
          onClick: enableAction ? () => directionAction("North") : undefined,
        },
        West: {
          onClick: enableAction ? () => directionAction("West") : undefined,
        },
        South: {
          onClick: enableAction ? () => directionAction("South") : undefined,
        },
        East: {
          onClick: enableAction ? () => directionAction("East") : undefined,
        },
      },
      AButton: {
        onClick: enableAction ? aAction : undefined,
      },
      BButton: {
        onClick: enableAction ? bAction : undefined,
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
  loading: boolean;
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
    if (p.loading) {
      return;
    }
    if (!p.inMyTurn) {
      return;
    }
    setActState((prev) => {
      if (prev.state === "done") return prev;
      if (prev.state === "actiontype") {
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
      if (prev.state === "actionat") {
        const d = new SectorDirection(p.boardWidth);
        const enables =
          prev.selectedActType === "move"
            ? p.enableMoveSectors
            : prev.selectedActType === "fire-torpedo"
              ? p.enableTorpedoSectors
              : p.enableMineSectors;
        const next =
          dir === "North"
            ? d.north(prev.selectedSct, enables)
            : dir === "South"
              ? d.south(prev.selectedSct, enables)
              : dir === "East"
                ? d.east(prev.selectedSct, enables)
                : d.west(prev.selectedSct, enables);
        if (next === -1) {
          return prev;
        }
        console.log("move to", dir);
        return {
          state: "actionat",
          selectedActType: prev.selectedActType,
          selectedSct: next,
        };
      }
      return prev;
    });
  };
  const aAction = () => {
    if (p.loading) {
      return;
    }
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
    if (p.loading) {
      return;
    }
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

  const enableAction = !p.loading && p.inMyTurn && actState.state !== "done";
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

interface FinishedPageParams {
  winner: "me" | "enemy";
  boardWidth: number;
  islandSectors: number[];

  firstAction: "me" | "enemy";
  gameOverReason: "torpedo-hit" | "mine-hit" | "timeout";
  myActionLogs: Array<FinishedPageParamsActionLog>;
  enemyActionLogs: Array<FinishedPageParamsActionLog>;

  exit: () => void;
}
type FinishedPageParamsActionLog =
  | {
      mode: "attack";
      me: number; // 艦位置
      sector: number; // 行動先
      type: "fire-torpedo" | "trigger-mine";
      result: "full-speed-ahead" | "hard-to-starboard" | "hit";
    }
  | {
      mode: "move";
      me: number; // 艦位置
      sector: number; // 行動先
      type: "move";
    }
  | {
      mode: "dummy";
    };
export const useFinishedPageProps = (p: FinishedPageParams) => {
  type Props = React.ComponentProps<typeof FinishedPage>;
  // 初期表示は最後の行動
  const [showSt, setShowSt] = useState<{
    show: "top" | "bottom";
    showIndex: number;
  }>({
    show:
      (p.winner === "me" && p.firstAction === "me") ||
      (p.winner === "enemy" && p.firstAction === "enemy")
        ? "top"
        : "bottom",
    showIndex: p.myActionLogs.length - 1,
  });
  const [commandUse, setCommandUse] = useState<"win" | "lose" | "exit">(
    p.winner === "me" ? "win" : "lose",
  );

  const makeOceanProps = (
    log: FinishedPageParamsActionLog,
  ): Props["TopOcean"][number] => {
    if (log.mode === "dummy") return { mode: "dummy" };
    return {
      mode: "action",
      me: log.me,
      sector: log.sector,
      type: log.type,
    };
  };

  const directionAction = (dir: InputDirection) => {
    setShowSt((prev) => {
      switch (dir) {
        case "North": {
          if (prev.show === "bottom") {
            return { ...prev, show: "top" };
          }
          // 前のページがあればbottomに移動
          if (prev.showIndex > 0) {
            return { show: "bottom", showIndex: prev.showIndex - 1 };
          }
          return prev;
        }
        case "South": {
          if (prev.show === "top") {
            return { ...prev, show: "bottom" };
          }
          // 次のページがあればtopに移動
          if (prev.showIndex < p.myActionLogs.length - 1) {
            return { show: "top", showIndex: prev.showIndex + 1 };
          }
          return prev;
        }
        case "West": {
          if (prev.showIndex > 0) {
            return { ...prev, showIndex: prev.showIndex - 1 };
          }
          return { ...prev, showIndex: p.myActionLogs.length - 1 };
        }
        case "East": {
          if (prev.showIndex < p.myActionLogs.length - 1) {
            return { ...prev, showIndex: prev.showIndex + 1 };
          }
          return { ...prev, showIndex: 0 };
        }
      }
    });
  };

  const makeOverlayedLoadingChildren = (
    index: number,
    log: FinishedPageParamsActionLog,
  ) => {
    switch (log.mode) {
      case "dummy":
        return "";
      case "attack": {
        const sct = sector(p.boardWidth, log.sector);
        const actTxt =
          log.type === "fire-torpedo" ? "魚雷が着弾" : "制御機雷が作動";
        switch (log.result) {
          case "full-speed-ahead":
            return `${index + 1}. ${actTxt}しました。${sct}`;
          case "hard-to-starboard":
            return `${index + 1}. 相手艦の至近に${actTxt}しました。${sct}`;
          default:
            return `${index + 1}. ${sct}に${actTxt}し、敵艦を撃沈しました。${sct}`;
        }
      }
      case "move": {
        const sct = sector(p.boardWidth, log.sector);
        return `${index + 1}. ${sct}に潜航しました。`;
      }
    }
  };
  let hitSector: number | undefined;
  for (const l of [...p.myActionLogs, ...p.enemyActionLogs]) {
    if (l.mode === "attack" && l.result === "hit") {
      hitSector = l.sector;
      break;
    }
  }

  const props: Props = {
    Leading: {
      inMyTrun:
        (showSt.show === "top" && p.firstAction === "me") ||
        (showSt.show === "bottom" && p.firstAction === "enemy"),
      children:
        (showSt.show === "top" && p.firstAction === "me") ||
        (showSt.show === "bottom" && p.firstAction === "enemy")
          ? makeOverlayedLoadingChildren(
              showSt.showIndex,
              p.myActionLogs[showSt.showIndex],
            )
          : makeOverlayedLoadingChildren(
              showSt.showIndex,
              p.enemyActionLogs[showSt.showIndex],
            ),
    },
    Ocean: {
      sectorCount: p.boardWidth,
      islands: p.islandSectors,
      show: showSt.show,
      showIndex: showSt.showIndex,
    },
    OverlayedLeading: {
      children:
        p.winner === "me"
          ? makeOverlayedLoadingChildren(
              showSt.showIndex,
              p.myActionLogs[showSt.showIndex],
            )
          : makeOverlayedLoadingChildren(
              showSt.showIndex,
              p.enemyActionLogs[showSt.showIndex],
            ),
    },
    OverlayedOcean: {
      ringSector: hitSector,
      gapSector: hitSector,
      title:
        p.gameOverReason === "torpedo-hit"
          ? "魚雷着弾"
          : p.gameOverReason === "mine-hit"
            ? "機雷直撃"
            : "逃走",
      titlePosition:
        // timeoutした方にtitle表示、その他攻撃の場合は勝者に表示
        p.gameOverReason === "timeout" && p.winner === "me"
          ? "left"
          : p.gameOverReason === "timeout" && p.winner === "enemy"
            ? "right"
            : p.winner === "me"
              ? "right"
              : "left",
    },
    TopOcean: (p.firstAction === "me" ? p.myActionLogs : p.enemyActionLogs).map(
      makeOceanProps,
    ),
    BottomOcean: (p.firstAction === "me"
      ? p.enemyActionLogs
      : p.myActionLogs
    ).map(makeOceanProps),
    Controller: {
      AButton: {
        onClick: () => {
          setCommandUse((prev) => {
            if (prev === "exit") {
              p.exit();
            }
            return "exit";
          });
        },
      },
      BButton: {
        onClick: () => {
          setCommandUse(p.winner === "me" ? "win" : "lose");
        },
      },
      CommandWindow: {
        commands: [p.winner === "me" ? "win" : "lose", "exit"],
        use: commandUse,
      },
      Direction: {
        North: {
          onClick: () => directionAction("North"),
        },
        South: {
          onClick: () => directionAction("South"),
        },
        East: {
          onClick: () => directionAction("East"),
        },
        West: {
          onClick: () => directionAction("West"),
        },
      },
    },
  };

  return props;
};
