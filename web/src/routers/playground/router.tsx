import {
  Box,
  Container,
  Fade,
  Flex,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  VisuallyHidden,
} from "@chakra-ui/react";
import { ConnectError } from "@connectrpc/connect";
import { useEffect } from "react";
import {
  type ActionFunctionArgs,
  type LoaderFunctionArgs,
  useLoaderData,
  useSubmit,
} from "react-router-dom";
import { battleClient } from "../../api/connect";
import {
  ActionRequest,
  DeployRequest,
  LogsRequest,
  type LogsResponse,
  WaitRequest,
} from "../../gen/api/v2/game_pb";
import { GameComponent } from "./game/game";
import { LogsComponent } from "./history/history";
import { ProgressBar } from "./progress";
import { StartingComponent } from "./start/start";

function Home() {
  const logs = useLoaderData() as LogsResponse;

  const submit = useSubmit();
  useEffect(() => {
    if (logs.requireAction) {
      return;
    }
    if (logs.requireDeployAction) {
      return;
    }
    if (logs.gameIsOver) {
      return;
    }
    submit(null, { method: "PATCH" });
  }, [logs.requireAction, logs.requireDeployAction, logs.gameIsOver, submit]);

  return (
    <Container>
      <Tabs index={logs.requireDeployAction ? 0 : 1} p={0}>
        <VisuallyHidden>
          <TabList>
            <Tab />
            <Tab />
          </TabList>
        </VisuallyHidden>

        <TabPanels>
          <TabPanel p={0}>
            <Fade in={logs.requireDeployAction}>
              <Flex
                direction={"column"}
                p={0}
                py={2}
                minH={"100svh"}
                maxH={"100svh"}
              >
                <StartingComponent />
                <Box mt={"auto"}>
                  <ProgressBar
                    callback={() => submit(null, { method: "DELETE" })}
                  />
                </Box>
              </Flex>
            </Fade>
          </TabPanel>
          <TabPanel p={0}>
            <Fade in={!logs.requireDeployAction}>
              <Flex
                direction={"column"}
                p={0}
                py={2}
                minH={"100svh"}
                maxH={"100svh"}
              >
                <LogsComponent />
                <GameComponent />
                <ProgressBar
                  callback={() => submit(null, { method: "DELETE" })}
                />
              </Flex>
            </Fade>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Container>
  );
}

export async function loader({ params }: LoaderFunctionArgs) {
  const { gameId, playerId } = params;
  try {
    const logs = await battleClient.logs(
      new LogsRequest({
        gameId,
        playerId,
      }),
    );
    return logs;
  } catch (err) {
    const connectErr = new ConnectError(err as string);
    console.error(connectErr.message);
  }
}

export async function action({ request, params }: ActionFunctionArgs) {
  const formData = await request.formData();
  const { gameId, playerId } = params;

  try {
    switch (request.method) {
      case "POST": {
        switch (formData.get("type")?.toString()) {
          case "first": {
            const at = formData.get("at")?.toString();
            const [mine1, mine2] = formData
              .get("mines")
              ?.toString()
              .split(",") || ["0", "0"];
            const req = new DeployRequest({
              gameId,
              playerId,
              at: Number(at),
              mines: [Number(mine1), Number(mine2)],
            });
            await battleClient.deploy(req, { signal: request.signal });
            break;
          }
          case "action": {
            const at = formData.get("at")?.toString();
            const strActionType = formData.get("act")?.toString();
            const actionType = Number(strActionType);
            const req = new ActionRequest({
              type: actionType,
              gameId,
              playerId,
              at: Number(at),
            });
            await battleClient.action(req, { signal: request.signal });
            break;
          }
        }
        break;
      }
      case "PATCH": {
        const req = new WaitRequest({
          gameId,
          playerId,
        });
        for await (const _ of battleClient.wait(req, {
          signal: request.signal,
        })) {
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

export default Home;
