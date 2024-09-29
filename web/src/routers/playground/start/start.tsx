import { Text, VStack } from "@chakra-ui/react";
import { Board } from "../components/borad";
import {
  CampStatus,
  HistoryResponse_Camp,
  HistoryResponse_Line,
} from "../../../gen/api/v1/game_pb";

export function StartingComponent() {
  return (
    <VStack>
      <Text>Hoge</Text>
      <Board camps={mock} />
    </VStack>
  );
}

const newCamp = (camp: number, status: CampStatus[]) => {
  return new HistoryResponse_Camp({ camp, status });
};

const mock: HistoryResponse_Line[] = [
  new HistoryResponse_Line({
    camps: [
      newCamp(0, [CampStatus.BOMB, CampStatus.MINE]),
      newCamp(1, [CampStatus.BOMB, CampStatus.MINE, CampStatus.MOVE]),
      newCamp(2, [CampStatus.ISLAND]),
      newCamp(3, []),
      newCamp(4, []),
      newCamp(5, []),
    ],
  }),
  new HistoryResponse_Line({
    camps: [
      newCamp(6, [CampStatus.BOMB, CampStatus.MOVE]),
      newCamp(7, [CampStatus.SUBMARINE]),
      newCamp(8, [CampStatus.BOMB, CampStatus.MOVE]),
      newCamp(9, []),
      newCamp(10, []),
      newCamp(11, []),
    ],
  }),
  new HistoryResponse_Line({
    camps: [
      newCamp(12, [CampStatus.BOMB]),
      newCamp(13, [CampStatus.BOMB, CampStatus.MOVE]),
      newCamp(14, [CampStatus.ISLAND]),
      newCamp(15, []),
      newCamp(16, []),
      newCamp(17, []),
    ],
  }),
  new HistoryResponse_Line({
    camps: [
      newCamp(18, []),
      newCamp(19, []),
      newCamp(20, []),
      newCamp(21, []),
      newCamp(22, []),
      newCamp(23, []),
    ],
  }),
  new HistoryResponse_Line({
    camps: [
      newCamp(24, []),
      newCamp(25, []),
      newCamp(26, []),
      newCamp(27, []),
      newCamp(28, []),
      newCamp(29, []),
    ],
  }),
  new HistoryResponse_Line({
    camps: [
      newCamp(30, []),
      newCamp(31, []),
      newCamp(32, []),
      newCamp(33, []),
      newCamp(34, []),
      newCamp(35, []),
    ],
  }),
];
