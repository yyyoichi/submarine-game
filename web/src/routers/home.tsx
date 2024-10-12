import { ArrowRightIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Container,
  Heading,
  Highlight,
  ListItem,
  StackDivider,
  type SystemStyleObject,
  Text,
  UnorderedList,
  VStack,
} from "@chakra-ui/react";
import { ConnectError } from "@connectrpc/connect";
import { useState } from "react";
import { Form, redirect } from "react-router-dom";
import { matchingClient } from "../api/connect";
import { JoinRequest, LeaveRequest } from "../gen/api/v2/game_pb";

const highlightStyle: SystemStyleObject = {
  textDecoration: "underline",
  fontWeight: "bold",
  color: "white.500",
};

function Home() {
  const [isLoading, setIsLoading] = useState(false);
  return (
    <Container p={0}>
      <Heading
        as="h1"
        size="3xl"
        fontFamily={"Train One"}
        fontStyle={"system-ui"}
      >
        🪖潜水艦ゲーム
      </Heading>
      <VStack py={8}>
        <Form
          method="post"
          onSubmit={() => {
            setIsLoading(true);
          }}
        >
          <Button
            size={"lg"}
            isLoading={isLoading}
            loadingText="対戦相手待ち合わせ..."
            type="submit"
            rightIcon={<ArrowRightIcon />}
          >
            {"開始"}
          </Button>
        </Form>
      </VStack>
      <VStack
        divider={<StackDivider borderColor="gray.200" />}
        spacing={4}
        align="stretch"
        px={1}
        pb={3}
      >
        <Heading
          as="h2"
          size="xl"
          fontFamily={"Train One"}
          fontStyle={"system-ui"}
        >
          {"遊び方"}
        </Heading>
        <Box>
          <Text>
            「潜水艦ゲーム」は、2人で通信対戦する戦略ボードゲームです。
          </Text>
          <Text>
            6x6のボード上で、潜水艦を操作して相手の潜水艦を魚雷で撃沈することを目指します。
          </Text>
        </Box>
        <Box>
          <UnorderedList spacing={1}>
            <Heading
              as={"h3"}
              size="md"
              fontFamily={"Train One"}
              fontStyle={"system-ui"}
            >
              {"基本ルール"}
            </Heading>
            <Text py={1}>
              ゲームはターン制で、プレイヤー次のような行動を選択します。
            </Text>
            <ListItem>
              <Highlight query="行動" styles={highlightStyle}>
                行動: 潜水艦を上下左右に一マス移動させることができます。
              </Highlight>
            </ListItem>
            <ListItem>
              <Highlight query="魚雷攻撃" styles={highlightStyle}>
                魚雷攻撃:
                同じ位置にとどまり、隣接する上下左右斜めのマスに魚雷を撃ちます。
              </Highlight>
            </ListItem>
            <ListItem>
              <Highlight query="機雷作動" styles={highlightStyle}>
                機雷作動:
                同じ位置にとどまり、あらかじめ敷設した機雷を作動させます。
              </Highlight>
            </ListItem>
          </UnorderedList>
        </Box>
        <Box>
          <UnorderedList spacing={1}>
            <Heading
              as={"h3"}
              size="md"
              fontFamily={"Train One"}
              fontStyle={"system-ui"}
            >
              勝敗の分かれ目
            </Heading>
            <Text py={1}>
              次のような魚雷や機雷攻撃の結果から、相手の位置を推測して勝利を目指します。
            </Text>
            <ListItem>
              <Highlight query="面舵一杯！" styles={highlightStyle}>
                面舵一杯！:
                魚雷・機雷攻撃の上下左右斜めの範囲に相手が潜行しています。
              </Highlight>
            </ListItem>
            <ListItem>
              <Highlight query="ヨーソロー！" styles={highlightStyle}>
                ヨーソロー！:
                少なくとも魚雷・機雷攻撃の上下左右斜めの範囲に相手は潜行していません。
              </Highlight>
            </ListItem>
          </UnorderedList>
        </Box>
      </VStack>
    </Container>
  );
}

export async function action() {
  let playerId = "";
  const leaveFromGame = async () => {
    if (!playerId) {
      return;
    }
    matchingClient.leave(new LeaveRequest());
  };
  window.addEventListener("beforeunload", leaveFromGame);

  try {
    const stream = matchingClient.join(new JoinRequest());
    for await (const resp of stream) {
      playerId = resp.playerId;
      if (resp.gameId === "") {
        continue;
      }
      return redirect(`/playground/${resp.gameId}/${resp.playerId}`);
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
  } finally {
    window.removeEventListener("beforeunload", leaveFromGame);
  }
  return null;
}

export default Home;
