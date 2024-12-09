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
import {
  ActionResult,
  ActionType,
  type LogsResponse,
} from "../../../gen/api/v2/game_pb";

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
  const jpResult = (r?: ActionResult) => {
    switch (r) {
      case ActionResult.HIT:
        return "\n>> 命中！";
      case ActionResult.FULL_SPEED_AHEAD:
        return "\n>> ヨーソロー";
      case ActionResult.HARD_TO_STARBOARD:
        return "\n>> 面舵一杯！";
    }
    return "";
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
                {me &&
                  (me.type === ActionType.MOVE
                    ? `海域${me.to}: ${jpType(me.type)}`
                    : `海域${me.at}: ${jpType(me.type)}${jpResult(me.result)}`)}
              </Td>
              <Td py={".5rem"} px={1} whiteSpace={"pre-line"}>
                {enemy &&
                  (enemy.turn === 0
                    ? `海域${enemy.at === -1 ? "?" : enemy.at}: ${jpType(enemy.type)}`
                    : enemy.type === ActionType.MOVE
                      ? `${jpDirection(enemy.direction)}方向: ${jpType(enemy.type)}`
                      : `海域${enemy.to}: ${jpType(enemy.type)}${jpResult(enemy.result)}`)}
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
}
