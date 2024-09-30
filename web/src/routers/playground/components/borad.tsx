import { Box, Grid, GridItem } from "@chakra-ui/react";
import { Cell } from "./cell";
import type { ComponentProps } from "react";

type BoardProps = {
  camps: ComponentProps<typeof Cell>[][];
};
export function Board({ camps }: BoardProps) {
  return (
    <Box px={7} width={"100%"}>
      <Grid p={1} width={"100%"} gap={1} bg={"blue.500"}>
        {camps.map((line) => {
          return (
            <GridItem key={line.at(0)?.camp}>
              <Grid templateColumns="repeat(6, 1fr)" width={"100%"} gap={1}>
                {line.map((cell) => {
                  return (
                    <GridItem key={cell.camp}>
                      <Cell {...cell} />
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
