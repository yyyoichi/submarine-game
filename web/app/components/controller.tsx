import { cn } from "@/lib/utils";
import { Circle, Compass, ShipWheel } from "lucide-react";
import type React from "react";
import { useEffect, useState } from "react";
import { Carousel, type CarouselApi, CarouselContent } from "./ui/carousel";

type TDirection = "North" | "West" | "South" | "East";
type TCommand =
  | "deploy-me"
  | "deploy-mines"
  | "deploy-wait"
  | "fire-torpedo"
  | "trigger-mine"
  | "move"
  | "wait";

type Props = {
  visibleDirection: boolean;
  Direction: Record<TDirection, Pick<DirectionProps, "onClick" | "onKeyDown">>;
  Compass: Pick<React.ComponentProps<"div">, "onClick" | "onKeyDown">;
  AButton: Pick<React.ComponentProps<"svg">, "onClick" | "onKeyDown">;
  BButton: Pick<React.ComponentProps<"svg">, "onClick" | "onKeyDown">;
  CommandWindow: CommandWindowProps;
};

export const ControllerComponent = (props: Props) => {
  const directionsProps: Record<TDirection, DirectionProps> = {
    North: { children: "北" },
    West: { children: "西" },
    South: { children: "南" },
    East: { children: "東" },
  };
  for (const d of Object.keys(directionsProps) as TDirection[]) {
    directionsProps[d] = {
      ...directionsProps[d],
      ...props.Direction[d],
      visible: props.visibleDirection,
    };
  }
  return (
    // コンパスとコマンド
    <div className="p-[10%] h-fit flex flex-col gap-2.5">
      <div className="flex justify-between items-center h-fit w-full">
        <div {...props.Compass}>
          <Compass className="w-10 h-10 stroke-3" />
        </div>
        <div className="w-36 rounded-xs">
          <CommandWindow {...props.CommandWindow} />
        </div>
      </div>
      {/* 操作キー */}
      <div className="flex h-fit gap-5">
        {/* 方向ボタン */}
        <div className="relative">
          <ShipWheel className="w-38 h-38 stroke-[1.5]" />
          {/* overlay buttons */}
          <div className="absolute top-1/2 left-1/2 w-[110%] h-[110%] transform -translate-x-1/2 -translate-y-1/2 text-background/95 text-xl font-bold">
            <div className="flex justify-center h-1/3">
              <Direction {...directionsProps.North} />
            </div>
            <div className="flex justify-between h-1/3">
              <Direction {...directionsProps.West} />
              <Direction {...directionsProps.East} />
            </div>
            <div className="flex justify-center h-1/3">
              <Direction {...directionsProps.South} />
            </div>
          </div>
        </div>
        {/* A,Bボタン */}
        <div className="flex flex-col justify-around w-full">
          <div className="w-full flex justify-end">
            <Circle
              className="w-12 h-12 stroke-4 cursor-pointer"
              {...props.AButton}
            />
          </div>
          <div className="w-full flex justify-center">
            <Circle
              className="w-12 h-12 stroke-4 cursor-pointer"
              {...props.BButton}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

type CommandWindowProps = {
  commands: TCommand[];
  use: TCommand;
};

const commandProps: Record<TCommand, React.ComponentProps<"div">> = {
  "deploy-me": { children: "出撃" },
  "deploy-mines": { children: "機雷敷設" },
  "fire-torpedo": { children: "魚雷発射" },
  "trigger-mine": { children: "機雷作動" },
  "deploy-wait": { children: "..." },
  move: { children: "潜航" },
  wait: { children: "待機" },
};
const CommandWindow = (props: CommandWindowProps) => {
  const [api, setApi] = useState<CarouselApi>();
  useEffect(() => {
    if (!api) return;
    const i = props.commands.indexOf(props.use);
    if (i === -1) return;
    api?.scrollTo(i);
  }, [api, props.use, props.commands]);
  return (
    <div className="text-3xl mx-auto w-full h-9 text-center ring-5 ring-muted-foreground bg-muted-foreground/5 tracking-widest">
      <Carousel orientation="vertical" className="h-9" setApi={setApi}>
        <CarouselContent className="mt-0 h-9">
          {props.commands.map((command) => (
            <div key={command} {...commandProps[command]} />
          ))}
        </CarouselContent>
      </Carousel>
    </div>
  );
};

type DirectionProps = {
  visible?: boolean;
} & Pick<React.ComponentProps<"div">, "children" | "onClick" | "onKeyDown">;

const Direction = ({ children, visible, ...props }: DirectionProps) => {
  return (
    <div
      className={cn(
        "h-full aspect-square flex justify-center items-center rounded-full cursor-pointer",
        visible ? "bg-foreground/10" : "",
      )}
      {...props}
    >
      <div
        className={cn(
          "w-fit h-fit p-0.5 rounded-full",
          visible ? "bg-foreground/20" : "",
        )}
      >
        {visible ? children : <></>}
      </div>
    </div>
  );
};
