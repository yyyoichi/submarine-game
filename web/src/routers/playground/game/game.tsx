import {
  Box,
  Button,
  Flex,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  VStack,
} from "@chakra-ui/react";
import { type ComponentProps, useEffect, useRef, useState } from "react";
import { Form, useLoaderData } from "react-router-dom";
import {
  ActionType,
  GameOverReason,
  type LogsResponse,
} from "../../../gen/api/v2/game_pb";
import { IconMove, IconTorpedo } from "../components/icon";
import { OceanMap } from "../components/ocean";

export function GameComponent() {
  const formRef = useRef<HTMLFormElement>(null);
  const logs = useLoaderData() as LogsResponse;
  const [isLoading, setIsLoading] = useState(false);
  const [clickSector, setClickSector] = useState<number | null>(null);
  useEffect(() => {
    if (logs.actionLogs.length > 0) {
      setClickSector(null);
    }
  }, [logs.actionLogs.length]);
  const [actionTypeSelection, setActionTypeSelection] =
    useState<ActionType | null>(null);
  const oceanMapProps: ComponentProps<typeof OceanMap> = {
    boardWidth: logs.boardWidth,
    sectors: [],
  };
  const enableActionType: ActionType[] = [];
  for (const sector of logs.sectors) {
    const sectorProps: ComponentProps<typeof OceanMap>["sectors"][number] = {
      isSelf: sector.selfOccupied,
      island: sector.island,
      sector: sector.sector,
      actions: sector.enableActions,
      bg: sector.sector === clickSector ? "blue.500" : undefined,
      onClick:
        logs.requireAction && sector.enableActions.length > 0
          ? () => {
              setClickSector(sector.sector);
              setActionTypeSelection(null);
            }
          : undefined,
    };

    oceanMapProps.sectors.push(sectorProps);
  }
  let gameMainText = "";
  switch (logs.gameIsOver) {
    case true:
      if (logs.win) {
        let reason = "";
        switch (logs.gameOverReason) {
          case GameOverReason.MINE_HIT:
            reason = "魚雷命中";
            break;
          case GameOverReason.TORPEDO_HIT:
            reason = "機雷命中";
            break;
          case GameOverReason.TIMEOUT:
            reason = "タイムアップ";
            break;
        }
        gameMainText = `${reason}！勝利！！`;
      } else {
        let reason = "";
        switch (logs.gameOverReason) {
          case GameOverReason.MINE_HIT:
            reason = "魚雷直撃";
            break;
          case GameOverReason.TORPEDO_HIT:
            reason = "機雷直撃";
            break;
          case GameOverReason.TIMEOUT:
            reason = "タイムアップ";
            break;
        }
        gameMainText = `${reason}！敗北...`;
      }
      break;
    case false:
      if (logs.requireAction) {
        gameMainText = "潜行か行動か";
      } else {
        gameMainText = "相手の行動を待機中";
      }
  }
  return (
    <Box mt={"auto"}>
      <Form
        method="POST"
        onSubmit={() => {
          setIsLoading(true);
        }}
        ref={formRef}
      >
        <input type="hidden" name="type" value="action" />
        <input type="hidden" name="at" value={clickSector || ""} />
        <VStack py={2}>
          <Text fontSize={"x-large"} fontWeight={"bold"} my={2}>
            {gameMainText}
          </Text>
          <OceanMap {...oceanMapProps} />
        </VStack>
        <Box visibility={"hidden"}>
          <input
            type="radio"
            name="act"
            value={ActionType.MOVE}
            checked={actionTypeSelection === ActionType.MOVE}
            readOnly
          />
          <input
            type="radio"
            name="act"
            value={ActionType.FIIRE_TORPEDO}
            checked={actionTypeSelection === ActionType.FIIRE_TORPEDO}
            readOnly
          />
          <input
            type="radio"
            name="act"
            value={ActionType.TRIGGER_MINE}
            checked={actionTypeSelection === ActionType.TRIGGER_MINE}
            readOnly
          />
        </Box>
        <Modal
          isOpen={clickSector !== null}
          onClose={() => setClickSector(null)}
          motionPreset="slideInBottom"
          portalProps={{ appendToParentPortal: true, containerRef: formRef }}
        >
          <ModalOverlay />
          <ModalContent mx={6} px={5} bg={"dark.500"}>
            <ModalHeader>{`海域${clickSector}`}</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <Flex gap={5} flexDirection={"column"}>
                {enableActionType.includes(ActionType.MOVE) && (
                  <Button
                    borderColor={"green.500"}
                    leftIcon={<IconMove fill={"green.500"} />}
                    bg={"dark.500"}
                    borderWidth={"2px 4px 3px 2px"}
                    fontSize={"large"}
                    color={
                      actionTypeSelection === ActionType.MOVE
                        ? undefined
                        : "gray.700"
                    }
                    onClick={() => setActionTypeSelection(ActionType.MOVE)}
                  >
                    潜行
                  </Button>
                )}
                {enableActionType.includes(ActionType.FIIRE_TORPEDO) && (
                  <Button
                    borderColor={"orange.500"}
                    leftIcon={<IconMove fill={"orange.500"} />}
                    bg={"dark.500"}
                    borderWidth={"2px 4px 3px 2px"}
                    fontSize={"large"}
                    color={
                      actionTypeSelection === ActionType.FIIRE_TORPEDO
                        ? undefined
                        : "gray.700"
                    }
                    onClick={() =>
                      setActionTypeSelection(ActionType.FIIRE_TORPEDO)
                    }
                  >
                    魚雷発射
                  </Button>
                )}
                {enableActionType.includes(ActionType.TRIGGER_MINE) && (
                  <Button
                    borderColor={"red.500"}
                    leftIcon={<IconTorpedo fill={"red.500"} />}
                    bg={"dark.500"}
                    borderWidth={"2px 4px 3px 2px"}
                    fontSize={"large"}
                    color={
                      actionTypeSelection === ActionType.TRIGGER_MINE
                        ? undefined
                        : "gray.700"
                    }
                    onClick={() =>
                      setActionTypeSelection(ActionType.TRIGGER_MINE)
                    }
                  >
                    機雷発動
                  </Button>
                )}
              </Flex>
            </ModalBody>

            <ModalFooter gap={2}>
              <Button
                size={"lg"}
                bg={"dark.500"}
                color={"white.500"}
                isLoading={isLoading}
                onClick={() => {
                  setClickSector(null);
                }}
              >
                キャンセル
              </Button>
              <Button
                size={"lg"}
                isLoading={isLoading}
                type={"submit"}
                isDisabled={!actionTypeSelection}
              >
                決定
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      </Form>
    </Box>
  );
}
