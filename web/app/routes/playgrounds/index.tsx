import {
  ActionResult,
  ActionType,
  GameOverReason,
  type LogsResponse_Action,
} from "@/gen/api/v2/game_pb";
import { battleClient } from "@/lib/connectapi";
import { ConnectError } from "@connectrpc/connect";
import { useEffect, useState } from "react";
import { redirect, useFetcher, useSubmit } from "react-router";
import type { Route } from "./+types";
import { FinishedPage, PlayingPage, PreparingPage } from "./component";
import {
  useFinishedPageProps,
  usePlayingPageProps,
  usePreparingPageProps,
} from "./props";

export default function Playground(props: Route.ComponentProps) {
  // biome-ignore lint/style/noNonNullAssertion: <explanation>
  const logs = props.loaderData!;
  const [usePageName, setUsePageName] = useState<"deploy" | "play" | "finish">(
    "deploy",
  );
  useEffect(() => {
    if (!logs) return;
    // ページ選定
    if (logs.gameIsOver) {
      setUsePageName("finish");
      return;
    }
    if (logs.requireDeployAction) {
      setUsePageName("deploy");
      return;
    }
    setUsePageName("play");
  }, [logs]);
  // wait request
  const submit = useSubmit();
  useEffect(() => {
    if (!logs) return;
    if (logs.gameIsOver) {
      return;
    }
    if (logs.requireAction) {
      return;
    }
    submit(null, { method: "PATCH" });
  }, [submit, logs]);
  // timout
  useEffect(() => {
    if (!logs) return;
    if (logs.gameIsOver) {
      return;
    }
    if (!logs.requireAction) {
      return;
    }
    const r = Number(logs.timeout) - Date.now();
    const timeout = setTimeout(() => {
      submit(null, { method: "DELETE" });
    }, r);
    return () => {
      clearTimeout(timeout);
    };
  }, [submit, logs]);

  const islandSectors = logs.sectors
    .filter((x) => x.island)
    .map((x) => x.sector);
  const fetcher = useFetcher();
  const deployingPageProps = usePreparingPageProps({
    timeout: new Date(Number(logs.timeout)),
    milliSecondPerTurn: Number(logs.millSecondPerTurn),
    boardWidth: logs.boardWidth,
    mineCount: 2,
    islandSectors: islandSectors,
    enableDeploySectors: logs.sectors
      .filter((x) => !x.island)
      .map((x) => x.sector),
    enableMineSectors: logs.sectors
      .filter((x) => !x.island)
      .map((x) => x.sector),
    deployAction: ({ deployAt, mines }) => {
      fetcher.submit(
        { type: "deploy", at: deployAt, mines: mines.join(",") },
        { method: "POST" },
      );
    },
    loading: fetcher.state !== "idle",
  });
  const filterSectors = (type: ActionType) => {
    return logs.sectors
      .filter((x) => x.enableActions.includes(type))
      .map((x) => x.sector);
  };
  type PrevAction = Parameters<typeof usePlayingPageProps>[0]["prevAction"];
  const getPrevAction = (): PrevAction => {
    const me = !logs.requireAction;
    const latestAction = logs.actionLogs.at(-1);
    const defaultAction: PrevAction = {
      mode: "hide",
      me: me,
      type: "deploy",
    };
    if (!latestAction) return defaultAction;
    const latestLog = me ? latestAction.me : latestAction.enemy;
    if (!latestLog) return defaultAction;
    if (latestLog.type === ActionType.UNSPECIFIED) return defaultAction;
    if (latestLog.type === ActionType.MOVE)
      return {
        mode: "hide",
        me: me,
        type: "move",
        direction:
          latestLog.direction === 0
            ? "North"
            : latestLog.direction === 1
              ? "East"
              : latestLog.direction === 2
                ? "South"
                : "West",
      };
    return {
      mode: "show",
      me: me,
      type:
        latestLog.type === ActionType.FIIRE_TORPEDO
          ? "fire-torpedo"
          : "trigger-mine",
      at: latestLog.at,
      result:
        latestLog.result === ActionResult.HIT
          ? "hit"
          : latestLog.result === ActionResult.HARD_TO_STARBOARD
            ? "hard-to-starboard"
            : "full-speed-ahead",
    };
  };
  const playingPageProps = usePlayingPageProps({
    timeout: new Date(Number(logs.timeout)),
    milliSecondPerTurn: Number(logs.millSecondPerTurn),
    boardWidth: logs.boardWidth,
    turn: logs.numTurn,
    inMyTurn: logs.requireAction,

    at: logs.sectors.filter((x) => x.selfOccupied)[0].sector,
    islandSectors: islandSectors,
    enableMoveSectors: filterSectors(ActionType.MOVE),
    enableTorpedoSectors: filterSectors(ActionType.FIIRE_TORPEDO),
    enableMineSectors: filterSectors(ActionType.TRIGGER_MINE),
    action: ({ at, type }) => {
      fetcher.submit({ type: "action", at: at, act: type }, { method: "POST" });
    },
    prevAction: getPrevAction(),
    loading: fetcher.state !== "idle",
  });

  type ActionHistoryLog = Parameters<
    typeof useFinishedPageProps
  >[0]["myActionLogs"][number];
  const makeHistoryLog = (a?: LogsResponse_Action): ActionHistoryLog => {
    const dummy: ActionHistoryLog = {
      mode: "dummy",
    };
    if (!a) return dummy;
    if (a.type === ActionType.MOVE)
      return {
        mode: "move",
        me: a.form,
        sector: a.to,
        type: "move",
      };
    return {
      mode: "attack",
      me: a.form,
      sector: a.to,
      type:
        a.type === ActionType.FIIRE_TORPEDO ? "fire-torpedo" : "trigger-mine",
      result:
        a.result === ActionResult.HIT
          ? "hit"
          : a.result === ActionResult.HARD_TO_STARBOARD
            ? "hard-to-starboard"
            : "full-speed-ahead",
    };
  };
  const myActionLogs = logs.actionLogs.map((x) => x.me).map(makeHistoryLog);
  const enemyActionLogs = logs.actionLogs
    .map((x) => x.enemy)
    .map(makeHistoryLog);
  const finishedPageProps = useFinishedPageProps({
    winner: logs.win ? "me" : "enemy",
    boardWidth: logs.boardWidth,
    islandSectors: islandSectors,
    gameOverReason:
      logs.gameOverReason === GameOverReason.TIMEOUT
        ? "timeout"
        : GameOverReason.TORPEDO_HIT
          ? "torpedo-hit"
          : "mine-hit",

    firstAction: logs.isFirstAction ? "me" : "enemy",
    myActionLogs: myActionLogs.length > 0 ? myActionLogs : [makeHistoryLog()],
    enemyActionLogs:
      enemyActionLogs.length > 0 ? enemyActionLogs : [makeHistoryLog()],
    exit: () => {
      fetcher.submit({ type: "exit" }, { method: "POST" });
    },
  });

  if (!logs) {
    window.alert("Unexpected Error! Please reload.");
    setTimeout(() => {
      redirect(`/playgrounds/${props.params.gameId}/${props.params.playerId}`);
    }, 1000 * 3);
    return <></>;
  }
  switch (usePageName) {
    case "deploy":
      return <PreparingPage {...deployingPageProps} />;
    case "play":
      return <PlayingPage {...playingPageProps} />;
    default:
      return <FinishedPage {...finishedPageProps} />;
  }
}

export async function clientLoader({ params }: Route.ClientLoaderArgs) {
  const { gameId, playerId } = params;
  try {
    const logs = await battleClient.logs({
      gameId: gameId,
      playerId: playerId,
    });
    return logs;
  } catch (e) {
    if (e instanceof ConnectError) {
      console.error(e.message);
    } else if (e instanceof Error) {
      const ce = new ConnectError(e.message);
      console.error(ce.message);
    } else {
      console.error(e);
    }
    return null;
  }
}

export async function clientAction({
  request,
  params,
}: Route.ClientActionArgs) {
  const formData = await request.formData();
  const { gameId, playerId } = params;

  try {
    switch (request.method) {
      case "POST": {
        switch (formData.get("type")?.toString()) {
          case "deploy": {
            const at = formData.get("at")?.toString();
            const [mine1, mine2] = formData
              .get("mines")
              ?.toString()
              .split(",") || ["0", "0"];
            await battleClient.deploy(
              {
                gameId: gameId,
                playerId: playerId,
                at: Number(at),
                mines: [Number(mine1), Number(mine2)],
              },
              { signal: request.signal },
            );
            break;
          }
          case "action": {
            const at = formData.get("at")?.toString();
            const strActionType = formData.get("act")?.toString();
            const actionType = Number(strActionType);
            await battleClient.action(
              {
                type: actionType,
                gameId: gameId,
                playerId: playerId,
                at: Number(at),
              },
              { signal: request.signal },
            );
            break;
          }
          case "exit":
            redirect("/");
            break;
        }
        break;
      }
      case "PATCH": {
        for await (const _ of battleClient.wait(
          {
            gameId: gameId,
            playerId: playerId,
          },
          {
            signal: request.signal,
          },
        )) {
        }
        break;
      }

      case "DELETE": {
        break;
      }
    }
  } catch (e) {
    if (e instanceof ConnectError) {
      console.error(e.message);
    } else if (e instanceof Error) {
      const ce = new ConnectError(e.message);
      console.error(ce.message);
    } else {
      console.error(e);
    }
  }
  return null;
}
