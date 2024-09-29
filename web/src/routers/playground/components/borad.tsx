import { Box, Grid, GridItem } from "@chakra-ui/react";
import type { HistoryResponse_Line } from "../../../gen/api/v1/game_pb";
import { Cell } from "./cell";

type BoardProps = {
  camps: HistoryResponse_Line[];
};
export function Board({ camps }: BoardProps) {
  return (
    <Box px={7} width={"100%"}>
      <Grid p={1} width={"100%"} gap={1} bg={"blue.500"}>
        {camps.map((line) => {
          return (
            <GridItem key={line.camps.at(0)?.camp}>
              <Grid templateColumns="repeat(6, 1fr)" width={"100%"} gap={1}>
                {line.camps.map((cell) => {
                  return (
                    <GridItem key={cell.camp}>
                      <Cell camp={cell.camp} status={cell.status} />
                    </GridItem>
                  );
                })}
              </Grid>
            </GridItem>
          );
        })}
      </Grid>
    </Box>
  );
}
