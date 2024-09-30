import {
  TableContainer,
  Table,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
} from "@chakra-ui/react";
import { useLoaderData } from "react-router-dom";
import type { HistoryResponse } from "../../../gen/api/v1/game_pb";

export function HistoryComponent() {
  const history = useLoaderData() as HistoryResponse;
  const histories = history.histories.sort((a, b) => a.turn - b.turn);
  const histProps: { me: string; enemy: string }[] = [];
  for (let i = 0; i < Math.round(histories.length / 2); i++) {
    histProps.push({ me: "", enemy: "" });
  }
  for (const h of histories) {
    const i = Math.floor((h.turn - 1) / 2);
    if (h.userId !== "") {
      histProps[i].me = `${h.description}${h.impact && `> ${h.impact}`}`;
    } else {
      histProps[i].enemy = h.description;
    }
  }
  return (
    <TableContainer h={"40vh"}>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th color={"white.500"}>あなた</Th>
            <Th>あいて</Th>
          </Tr>
        </Thead>
        <Tbody fontSize={"md"}>
          {histProps.map((line, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
            <Tr key={i}>
              <Td py={".5rem"} px={1}>
                {line.me}
              </Td>
              <Td py={".5rem"} px={1}>
                {line.enemy}
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
}
