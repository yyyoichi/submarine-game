import { Box, Grid, GridItem } from "@chakra-ui/react";
import type { ComponentProps } from "react";
import { Sector } from "./cell";

type OceanMapProps = {
  sectors: ComponentProps<typeof Sector>[];
  boardWidth: number;
};
export function OceanMap(props: OceanMapProps) {
  const gridProps: ComponentProps<typeof Grid> = {
    templateColumns: `repeat(${props.boardWidth}, 1fr)`,
    p: 1,
    width: "100%",
    gap: 1,
    bg: "blue.500",
  };
  return (
    <Box px={7} width={"100%"}>
      <Grid {...gridProps}>
        {props.sectors.map((sector) => {
          return (
            <GridItem key={sector.sector}>
              <Sector {...sector} />
            </GridItem>
          );
        })}
      </Grid>
    </Box>
  );
}
