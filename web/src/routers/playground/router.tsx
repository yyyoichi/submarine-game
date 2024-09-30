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
import { useEffect } from "react";
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
  const doneFirstAction = history.histories.length > 0;

  const submit = useSubmit();
  useEffect(() => {
    if (history.myTurn) {
      return;
    }
    submit(null, { method: "PATCH" });
  }, [history.myTurn, submit]);

  return (
    <Container p={0} minH={"100vh"} maxH={"100vh"}>
        <Tabs index={doneFirstAction ? 1 : 0} p={0}>
          <VisuallyHidden>
            <TabList>
              <Tab />
              <Tab />
            </TabList>
          </VisuallyHidden>

          <TabPanels>
            <TabPanel>
              <Fade in={!doneFirstAction}>
                <StartingComponent />
              </Fade>
            </TabPanel>
            <TabPanel p={0}>
              <Fade in={doneFirstAction}>
                <GameComponent />
                <ProgressBar
                  callback={() => submit(null, { method: "DELETE" })}
                />
                <HistoryComponent />
              </Fade>
            </TabPanel>
          </TabPanels>
        </Tabs>
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
