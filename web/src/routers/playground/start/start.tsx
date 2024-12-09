import {
  Button,
  Fade,
  Flex,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  VStack,
} from "@chakra-ui/react";
import { type ComponentProps, useState } from "react";
import { Form, useLoaderData } from "react-router-dom";
import { ActionType, type LogsResponse } from "../../../gen/api/v2/game_pb";
import { IconMine, IconMyLocation } from "../components/icon";
import { OceanMap } from "../components/ocean";

export function StartingComponent() {
  const logs = useLoaderData() as LogsResponse;
  const [tabIndex, setTabIndex] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [startPlace, setStartPlace] = useState<number | null>(null);
  const [startMines, setStartMines] = useState<number[]>([]);
  const atOceanMapProps: ComponentProps<typeof OceanMap> = {
    sectors: [],
    boardWidth: logs.boardWidth,
  };
  const minesOceanMapProps: ComponentProps<typeof OceanMap> = {
    sectors: [],
    boardWidth: logs.boardWidth,
  };
  for (const sector of logs.sectors) {
    const atProps: ComponentProps<typeof OceanMap>["sectors"][number] = {
      isSelf: true,
      island: sector.island,
      sector: sector.sector,
      actions: [],
      bg: sector.sector === startPlace ? "blue.500" : undefined,
      onClick: !sector.island
        ? () => {
            setStartPlace(sector.sector);
            // NOTE ほんとはアニメーションで対処したい
            setTimeout(() => {
              setTabIndex(1);
            }, 100);
          }
        : undefined,
    };
    atOceanMapProps.sectors.push(atProps);

    const minesProps: ComponentProps<typeof OceanMap>["sectors"][number] = {
      isSelf: false,
      island: sector.island,
      sector: sector.sector,
      actions: [ActionType.TRIGGER_MINE],
      bg: startMines.includes(sector.sector) ? "orange.500" : undefined,
      onClick: sector.island
        ? undefined
        : startMines.includes(sector.sector)
          ? () => {
              setStartMines((pv) => {
                return pv.filter((x) => x !== sector.sector);
              });
            }
          : () => {
              setStartMines((pv) => {
                return [sector.sector, ...pv].splice(0, 2);
              });
            },
    };
    minesOceanMapProps.sectors.push(minesProps);
  }

  return (
    <Form
      method="POST"
      onSubmit={() => {
        setIsLoading(true);
      }}
    >
      <input type="hidden" name="type" value="first" />
      <input type="hidden" name="at" value={startPlace || ""} />
      <input type="hidden" name="mines" value={startMines.join(",")} />
      <VStack py={2}>
        <Tabs
          index={tabIndex}
          onChange={(index) => setTabIndex(index)}
          width={"100%"}
          border={"0"}
        >
          <TabList border={"0"}>
            <Tab>
              <IconMyLocation
                fill={tabIndex === 0 ? "white.500" : "gray.300"}
              />
              <Text mx={1} color={tabIndex === 0 ? "white.500" : "gray.500"}>
                行動開始
              </Text>
            </Tab>
            <Tab>
              <IconMine fill={tabIndex === 1 ? "white.500" : "gray.300"} />
              <Text mx={1} color={tabIndex === 1 ? "white.500" : "gray.500"}>
                機雷敷設
              </Text>
            </Tab>
          </TabList>
          <TabPanels p={0}>
            <TabPanel p={0} transitionDelay={""}>
              <Fade in={tabIndex === 0} delay={{ exit: 0.1 }}>
                <Text py={3}>潜行開始する海域を選択</Text>
                <OceanMap {...atOceanMapProps} />
              </Fade>
            </TabPanel>
            <TabPanel p={0}>
              <Fade in={tabIndex === 1} delay={{ enter: 0.1 }}>
                <Text py={3}>機雷を敷設する海域を選択</Text>
                <OceanMap {...minesOceanMapProps} />
                <Flex justifyContent={"center"} py={10}>
                  <Button
                    size={"lg"}
                    isLoading={isLoading}
                    type="submit"
                    isDisabled={!startPlace || startMines.length !== 2}
                  >
                    作戦開始
                  </Button>
                </Flex>
              </Fade>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </VStack>
    </Form>
  );
}
