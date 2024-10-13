import {
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { useLoaderData } from "react-router-dom";
import { ActionType, type LogsResponse } from "../../../gen/api/v2/game_pb";

export function LogsComponent() {
  const logs = useLoaderData() as LogsResponse;
  const jpType = (t?: ActionType) => {
    switch (t) {
      case ActionType.MOVE:
        return "潜行";
      case ActionType.FIIRE_TORPEDO:
        return "魚雷攻撃";
      case ActionType.TRIGGER_MINE:
        return "機雷作動";
    }
    return "潜行開始";
  };
  const jpDirection = (d?: number) => {
    switch (d) {
      case 0:
        return "北";
      case 1:
        return "東";
      case 2:
        return "南";
      case 3:
        return "西";
    }
  };
  return (
    <TableContainer maxH={"100%"} overflowY={"auto"}>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th color={"white.500"}>あなた</Th>
            <Th>あいて</Th>
          </Tr>
        </Thead>
        <Tbody fontSize={"md"}>
          {logs.actionLogs.map(({ me, enemy }, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
            <Tr key={i}>
              <Td py={".5rem"} px={1} whiteSpace={"pre-line"}>
                {`海域${me?.to}: ${jpType(me?.type)}`}
              </Td>
              <Td py={".5rem"} px={1} whiteSpace={"pre-line"}>
                {enemy?.turn === 0
                  ? `海域?: ${enemy.type}`
                  : enemy?.type === ActionType.MOVE
                    ? `${jpDirection(enemy.direction)}方向: ${jpType(me?.type)}`
                    : `海域${enemy?.at}: ${jpType(me?.type)}`}
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
}
