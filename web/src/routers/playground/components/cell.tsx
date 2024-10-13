import { Flex, Square, type SquareProps } from "@chakra-ui/react";
import { ActionType } from "../../../gen/api/v2/game_pb";
import {
  IconLandscape,
  IconMine,
  IconMove,
  IconMyLocation,
  IconTorpedo,
} from "./icon";

type SectorProps = {
  sector: number;
  actions: ActionType[];
  island: boolean;
  isSelf: boolean;
  bg?: string;
} & Pick<SquareProps, "onClick">;

export function Sector(props: SectorProps) {
  const CellWrap = (sps: SquareProps) => {
    const p: SquareProps = {
      aspectRatio: 1,
      bg: "dark.500",
      position: "relative",
      pt: 1,
      bgColor: props.bg,
      onClick: props.onClick,
      ...sps,
    };
    return <Square {...p}>{sps.children}</Square>;
  };

  if (props.isSelf) {
    return (
      <CellWrap>
        <IconMyLocation fill={"white.500"} width={"50%"} height={"50%"} />
      </CellWrap>
    );
  }
  if (props.island) {
    return (
      <CellWrap>
        <IconLandscape fill={"white.500"} width={"50%"} height={"50%"} />
      </CellWrap>
    );
  }
  return (
    <CellWrap>
      {props.sector}
      <Flex position={"absolute"} top={0} left={0} p={0} gap={0}>
        {props.actions.sort().map((s) => {
          switch (s) {
            case ActionType.MOVE:
              return <IconMove key={s} fill={"green.500"} />;
            case ActionType.FIIRE_TORPEDO:
              return <IconTorpedo key={s} fill={"orange.500"} />;
            case ActionType.TRIGGER_MINE:
              return <IconMine key={s} fill={"red.500"} />;
          }
        })}
      </Flex>
    </CellWrap>
  );
}
