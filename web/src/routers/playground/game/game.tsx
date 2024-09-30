import { Form, useLoaderData } from "react-router-dom";
import {
  ActionType,
  CampStatus,
  type HistoryResponse,
} from "../../../gen/api/v1/game_pb";
import { useEffect, useRef, useState, type ComponentProps } from "react";
import { Board } from "../components/borad";
import {
  VStack,
  Text,
  Flex,
  Box,
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from "@chakra-ui/react";
import { IconBomb, IconMove } from "../components/icon";

export function GameComponent() {
  const formRef = useRef<HTMLFormElement>(null);
  const history = useLoaderData() as HistoryResponse;
  const [isLoading, setIsLoading] = useState(false);
  const [clickCamp, setClickCamp] = useState<number | null>(null);
  useEffect(() => {
    setIsLoading(history.histories.length < 0);
    setClickCamp(null);
  }, [history.histories.length]);
  const [actionTypeSelection, setActionTypeSelection] = useState<
    ActionType.MOVE | ActionType.BOMB | ActionType.MINE | null
  >(null);
  const boardProps: ComponentProps<typeof Board> = { camps: [] };
  const enableActionType: ActionType[] = [];
  for (let i = 0; i < history.camps.length; i++) {
    const line = history.camps[i];
    boardProps.camps[i] = [];
    for (let j = 0; j < line.camps.length; j++) {
      const camp = line.camps[j];
      // 島か自分の位置以外でステータスがあれば行動可能位置
      const canClick =
        !camp.status.includes(CampStatus.ISLAND) &&
        !camp.status.includes(CampStatus.SUBMARINE) &&
        history.winner === "" &&
        camp.status.length;
      boardProps.camps[i][j] = {
        camp: camp.camp ?? 0,
        status: camp.status,
        onClick: canClick
          ? () => {
              setClickCamp(camp.camp ?? 0);
              setActionTypeSelection(null);
            }
          : undefined,
        bg: clickCamp === camp.camp ? "blue.500" : undefined,
      };
      if (clickCamp !== camp.camp) continue;
      if (clickCamp === null) continue;
      if (camp.status.includes(CampStatus.MOVE)) {
        enableActionType.push(ActionType.MOVE);
      }
      if (camp.status.includes(CampStatus.BOMB)) {
        enableActionType.push(ActionType.BOMB);
      }
      if (camp.status.includes(CampStatus.MINE)) {
        enableActionType.push(ActionType.MINE);
      }
    }
  }
  return (
    <Form
      method="POST"
      onSubmit={() => {
        setIsLoading(true);
      }}
      ref={formRef}
    >
      <input type="hidden" name="type" value="action" />
      <input type="hidden" name="place" value={clickCamp || ""} />
      <VStack py={2}>
        <Text fontSize={"x-large"} fontWeight={"bold"} my={2}>
          {history.description}
        </Text>
        <Board {...boardProps} />
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
          value={ActionType.BOMB}
          checked={actionTypeSelection === ActionType.BOMB}
          readOnly
        />
        <input
          type="radio"
          name="act"
          value={ActionType.MINE}
          checked={actionTypeSelection === ActionType.MINE}
          readOnly
        />
      </Box>
      <Modal
        isOpen={clickCamp !== null}
        onClose={() => setClickCamp(null)}
        motionPreset="slideInBottom"
        portalProps={{ appendToParentPortal: true, containerRef: formRef }}
      >
        <ModalOverlay />
        <ModalContent mx={6} px={5} bg={"dark.500"}>
          <ModalHeader>{`海域${clickCamp}`}</ModalHeader>
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
                  行動
                </Button>
              )}
              {enableActionType.includes(ActionType.BOMB) && (
                <Button
                  borderColor={"orange.500"}
                  leftIcon={<IconMove fill={"orange.500"} />}
                  bg={"dark.500"}
                  borderWidth={"2px 4px 3px 2px"}
                  fontSize={"large"}
                  color={
                    actionTypeSelection === ActionType.BOMB
                      ? undefined
                      : "gray.700"
                  }
                  onClick={() => setActionTypeSelection(ActionType.BOMB)}
                >
                  魚雷発射
                </Button>
              )}
              {enableActionType.includes(ActionType.MINE) && (
                <Button
                  borderColor={"red.500"}
                  leftIcon={<IconBomb fill={"red.500"} />}
                  bg={"dark.500"}
                  borderWidth={"2px 4px 3px 2px"}
                  fontSize={"large"}
                  color={
                    actionTypeSelection === ActionType.MINE
                      ? undefined
                      : "gray.700"
                  }
                  isDisabled={!enableActionType.includes(ActionType.MINE)}
                  onClick={() => setActionTypeSelection(ActionType.MINE)}
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
                setClickCamp(null);
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
  );
}
