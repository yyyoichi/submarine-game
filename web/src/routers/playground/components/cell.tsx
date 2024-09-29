import { Flex, Square, type SquareProps } from "@chakra-ui/react";
import { CampStatus } from "../../../gen/api/v1/game_pb";
import {
  IconBomb,
  IconLandscape,
  IconMine,
  IconMove,
  IconMyLocation,
} from "./icon";

// Chakra UIのBoxコンポーネントをmotion対応にする

type CellProps = {
  camp: number;
  status: CampStatus[];
};

export function Cell(props: CellProps) {
  if (props.status.includes(CampStatus.SUBMARINE)) {
    return (
      <CellWrap>
        <IconMyLocation fill={"white.500"} width={"50%"} height={"50%"} />
      </CellWrap>
    );
  }
  if (props.status.includes(CampStatus.ISLAND)) {
    return (
      <CellWrap>
        <IconLandscape fill={"white.500"} width={"50%"} height={"50%"} />
      </CellWrap>
    );
  }
  return (
    <CellWrap>
      {props.camp}
      <Flex position={"absolute"} top={0} left={0} p={0} gap={0}>
        {props.status.sort().map((s) => {
          switch (s) {
            case CampStatus.MOVE:
              return <IconMove key={s} fill={"green.500"} />;
            case CampStatus.BOMB:
              return <IconBomb key={s} fill={"orange.500"} />;
            case CampStatus.MINE:
              return <IconMine key={s} fill={"red.500"} />;
          }
        })}
      </Flex>
    </CellWrap>
  );
}

function CellWrap(props: SquareProps) {
  const p: SquareProps = {
    aspectRatio: 1,
    bg: "dark.500",
    position: "relative",
    pt: 1,
    ...props,
  };
  return <Square {...p}>{props.children}</Square>;
}
