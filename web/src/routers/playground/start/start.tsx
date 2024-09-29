import {
	Button,
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
import {
	CampStatus,
	HistoryResponse_Camp,
	HistoryResponse_Line,
} from "../../../gen/api/v1/game_pb";
import { IconMine, IconMyLocation } from "../components/icon";
import { useState } from "react";
import { Form } from "react-router-dom";

export function StartingComponent() {
	const [tabIndex, setTabIndex] = useState(0);
	const [isLoading, setIsLoading] = useState(false);
	const [startCamps, setStartCamps] = useState<{
		place: number;
		mines: number[];
	} | null>(null);
	return (
		<Form>
			<input type="hidden" name="type" value="first" />
			<input type="hidden" name="place" value={startCamps?.place} />
			<input type="hidden" name="mines" value={startCamps?.mines.join(",")|| ""} />
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
					<TabPanels p={1}>
						<TabPanel p={0}>
							<Text py={3}>行動開始する海域を選択</Text>
							<Board camps={mock} />
						</TabPanel>
						<TabPanel p={0}>
							<Text py={3}>機雷を敷設する海域を選択</Text>
							<Board camps={mock} />
							<Flex justifyContent={"center"} py={10}>
								<Button
									size={"lg"}
									isLoading={isLoading}
									// loadingText="対戦相手待ち合わせ..."
									type="submit"
									onClick={() => {
										setIsLoading(true);
									}}
								>
									作戦開始
								</Button>
							</Flex>
						</TabPanel>
					</TabPanels>
				</Tabs>
			</VStack>
		</Form>
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
