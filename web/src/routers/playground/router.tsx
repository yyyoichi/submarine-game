import {
  useLoaderData,
  type LoaderFunctionArgs,
  type ActionFunctionArgs,
  useSubmit,
} from "react-router-dom";
import { getGameClient } from "../../api/connect";
import {
  ActionRequest,
  ActionType,
  HistoryRequest,
  type HistoryResponse,
  WaitRequest,
  FirstActionRequest,
} from "../../gen/api/v1/game_pb";
import { ConnectError } from "@connectrpc/connect";
import { useEffect, useRef } from "react";
import { HistoryComponent } from "./history/history";
import { StartingComponent } from "./start/start";
import {
  Container,
  Fade,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  VisuallyHidden,
} from "@chakra-ui/react";
import { GameComponent } from "./game/game";
import { ProgressBar } from "./progress";

function Home() {
  const history = useLoaderData() as HistoryResponse;
  const waitFirstAction = history.histories.length < 2;
  const gameIsOver = history.winner !== "";

  const formRef = useRef<HTMLFormElement>(null);
  const submit = useSubmit();
  useEffect(() => {
    if (gameIsOver) {
      return;
    }
    if (!formRef.current) {
      return;
    }
    if (waitFirstAction || !history.myTurn) {
      submit(formRef.current, { method: "PATCH" });
    }
  }, [gameIsOver, waitFirstAction, history.myTurn, submit]);

  return (
    <Container p={0}>
      <ProgressBar callback={() => submit(null, { method: "DELETE" })} />
      {gameIsOver ? (
        <>
          <HistoryComponent />
          <GameComponent />
        </>
      ) : (
        <Tabs index={waitFirstAction ? 0 : 1}>
          <VisuallyHidden>
            <TabList>
              <Tab />
              <Tab />
            </TabList>
          </VisuallyHidden>

          <TabPanels>
            <TabPanel>
              <Fade in={waitFirstAction}>
                <StartingComponent />
              </Fade>
            </TabPanel>
            <TabPanel>
              <Fade in={!waitFirstAction}>
                <HistoryComponent />
                <GameComponent />
              </Fade>
            </TabPanel>
          </TabPanels>
        </Tabs>
      )}
    </Container>
  );
}

export async function loader({ params }: LoaderFunctionArgs) {
  const { gameId, userId } = params;
  const client = getGameClient();
  try {
    const history = await client.history(
      new HistoryRequest({
        gameId: gameId ?? "",
        userId: userId ?? "",
      }),
    );
    return history;
  } catch (err) {
    const connectErr = new ConnectError(err as string);
    console.error(connectErr.message);
  }
}

export async function action({ request, params }: ActionFunctionArgs) {
  const formData = await request.formData();
  const { gameId, userId } = params;

  try {
    switch (request.method) {
      case "POST": {
        switch (formData.get("type")?.toString()) {
          case "first": {
            const place = formData.get("place")?.toString();
            const [mine1, mine2] = formData
              .get("mines")
              ?.toString()
              .split(",") || ["0", "0"];
            const req = new FirstActionRequest({
              gameId,
              userId,
              camp: Number(place),
              mineCamps: [Number(mine1), Number(mine2)],
            });
            const clinet = getGameClient();
            await clinet.firstAction(req, { signal: request.signal });
            break;
          }
          case "action": {
            const place = formData.get("place")?.toString();
            let actionType = ActionType.UNSPECIFIED;
            const strActionType = formData.get("act")?.toString();
            console.log(strActionType);
            if (strActionType === "1") {
              actionType = ActionType.MOVE;
            } else if (strActionType === "2") {
              actionType = ActionType.BOMB;
            } else if (strActionType === "5") {
              actionType = ActionType.MINE;
            }
            const req = new ActionRequest({
              type: actionType,
              gameId,
              userId,
              camp: Number(place),
            });
            const clinet = getGameClient();
            await clinet.action(req, { signal: request.signal });
            break;
          }
        }
        break;
      }
      case "PATCH": {
        const clinet = getGameClient();
        const req = new WaitRequest({
          gameId,
          userId,
        });
        for await (const _ of clinet.wait(req, { signal: request.signal })) {
        }
        break;
      }

      case "DELETE": {
        const clinet = getGameClient();
        const req = new ActionRequest({
          type: ActionType.LEAVE,
          gameId,
          userId,
          camp: 0,
        });
        await clinet.action(req, { signal: request.signal });
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

export default Home;
