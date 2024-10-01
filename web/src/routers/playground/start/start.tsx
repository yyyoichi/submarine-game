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
import { Board } from "../components/borad";
import { CampStatus, type HistoryResponse } from "../../../gen/api/v1/game_pb";
import { IconMine, IconMyLocation } from "../components/icon";
import { type ComponentProps, useState } from "react";
import { Form, useLoaderData } from "react-router-dom";

export function StartingComponent() {
  const history = useLoaderData() as HistoryResponse;
  const [tabIndex, setTabIndex] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [startPlace, setStartPlace] = useState<number | null>(null);
  const [startMines, setStartMines] = useState<number[]>([]);
  const placeBoardProps: ComponentProps<typeof Board> = { camps: [] };
  const minesBoardProps: ComponentProps<typeof Board> = { camps: [] };
  for (let i = 0; i < history.camps.length; i++) {
    const line = history.camps[i];
    placeBoardProps.camps[i] = [];
    minesBoardProps.camps[i] = [];
    for (let j = 0; j < line.camps.length; j++) {
      const camp = line.camps[j];
      placeBoardProps.camps[i][j] = {
        camp: camp.camp ?? 0,
        status: camp.status.filter(
          (x) => x === CampStatus.PLACE || x === CampStatus.ISLAND,
        ),
        onClick: camp.status.includes(CampStatus.PLACE)
          ? () => {
              setStartPlace(camp.camp);
              // NOTE ほんとはアニメーションで対処したい
              setTimeout(() => {
                setTabIndex(1);
              }, 100);
            }
          : undefined,
        bg: startPlace === camp.camp ? "blue.500" : undefined,
      };
      minesBoardProps.camps[i][j] = {
        camp: camp.camp || 0,
        status: camp.status.filter(
          (x) => x === CampStatus.MINE || x === CampStatus.ISLAND,
        ),
        onClick: camp.status.includes(CampStatus.MINE)
          ? () => {
              setStartMines((pv) => {
                return [camp.camp, ...pv].splice(0, 2);
              });
            }
          : undefined,
        bg: startMines.includes(camp.camp) ? "orange.500" : undefined,
      };
    }
  }

  return (
    <Form
      method="POST"
      onSubmit={() => {
        setIsLoading(true);
      }}
    >
      <input type="hidden" name="type" value="first" />
      <input type="hidden" name="place" value={startPlace || ""} />
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
                <Text py={3}>行動開始する海域を選択</Text>
                <Board {...placeBoardProps} />
              </Fade>
            </TabPanel>
            <TabPanel p={0}>
              <Fade in={tabIndex === 1} delay={{ enter: 0.1 }}>
                <Text py={3}>機雷を敷設する海域を選択</Text>
                <Board {...minesBoardProps} />
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
