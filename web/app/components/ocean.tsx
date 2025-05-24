import { sector } from "@/lib/utils";
import type { ClassValue } from "clsx";
import { ChevronsRight, CircleDot, CircleX, SquareSquare } from "lucide-react";
import { cn } from "src/lib/utils";
import { ConcentricRingsComponent } from "./rings";

type Props = {
  displaySectorName: boolean;
  sectorCount: number; // 1辺のセクター数
  Sectors: Record<number, { embed?: boolean } & Pick<OceanSectorProps, "icon">>;
};

export const OceanComponent = ({
  Sectors,
  displaySectorName,
  sectorCount,
}: Props) => {
  return (
    <div className="grid grid-cols-6 py-2 px-3 gap-2 aspect-square bg-muted-foreground rounded-xs">
      {new Array(sectorCount * sectorCount).fill("").map((_, i) => {
        const sname = sector(sectorCount, i);
        const { embed, ...ps } = Sectors[i] || {};
        const p: OceanSectorProps = {
          type: embed ? "embed" : undefined,
          sectorName: displaySectorName ? sname : undefined,
          ...ps,
        };
        return <OceanSector key={sname} {...p} />;
      })}
    </div>
  );
};

type OverlayProps = {
  fadeout?: boolean; // フェードアウトするか
  absolute?: boolean; // 絶対位置にするか
  gapSector?: number; // 唯一空けるセクター
  ringSector?: number; // ソナーを出すセクター
  sectorCount: number; // 1辺のセクター数
} & (
  | {
      title?: undefined;
      titlePosition?: undefined;
    }
  | {
      title: string;
      titlePosition: "right" | "left";
    }
);

export const OverlayedOceanComponentAbsoluteClassName =
  "absolute top-0 left-0 z-100";

export const OverlayedOceanComponent = (props: OverlayProps) => {
  const titleCn: ClassValue[] = [];
  if (props.title) {
    const topOrBottom =
      (props.gapSector || 0) < props.sectorCount ** 2 / 2 ? "bottom" : "top";
    if (props.titlePosition === "right") {
      titleCn.push("justify-end");
    }
    if (topOrBottom === "bottom") {
      titleCn.push("items-end");
    }
  }
  return (
    <div
      className={cn(
        "grid grid-cols-6 py-2 px-3 gap-2 aspect-square bg-foreground w-full",
        props.fadeout ? "animate-fadeout" : "",
        props.absolute ? OverlayedOceanComponentAbsoluteClassName : "",
      )}
    >
      {[...Array(props.sectorCount * props.sectorCount)].map((_, i) => {
        const sname = sector(props.sectorCount, i);
        const p: OceanSectorProps = {
          type: props.gapSector !== i ? "overlay" : undefined,
          sectorName: props.gapSector === i ? sname : undefined,
        };
        return (
          <div key={sname} className="relative">
            {props.ringSector === i && <ConcentricRingsComponent />}
            <OceanSector {...p} />
          </div>
        );
      })}
      {props.title && (
        <div
          className={cn(
            "absolute w-full h-full top-0 left-0 flex text-background text-4xl",
            ...titleCn,
          )}
        >
          {props.title}
        </div>
      )}
    </div>
  );
};

type OceanSectorProps = {
  type?: "embed" | "overlay";
  icon?: "move" | "torpedo" | "mine" | "me";
  sectorName?: string;
};

const OceanSector = (props: OceanSectorProps) => {
  const addedClass: ClassValue[] = [];
  switch (props.type) {
    case "embed": {
      //  枠を残して背景を薄めで潰す
      addedClass.push("border-backgraund", "bg-muted-foreground");
      break;
    }
    case "overlay":
      // 枠を残さず背景を濃い色で潰す
      addedClass.push("border-foreground", "bg-foreground");
      break;
    default:
      // 枠を背景をと同じ色
      addedClass.push("border-backgraund", "bg-background");
      break;
  }
  const SectorIcon = () => {
    switch (props.icon) {
      case "move":
        return <ChevronsRight className="w-1/2 h-1/2 stroke-[3]" />;
      case "torpedo":
        return <CircleX className="stroke-[3]" />;
      case "mine":
        return <SquareSquare className="stroke-[3]" />;
      case "me":
        return <CircleDot className="text-foreground" />;
    }
  };
  return (
    <div
      className={cn(
        "w-full h-full relative border-1 rounted-xs",
        ...addedClass,
      )}
    >
      {/* iconとセクター位置補助をレイヤーする */}
      {props.icon && (
        <div
          className={
            "absolute w-full h-full top-0 left-0 z-10 flex justify-center items-center bg-background"
          }
        >
          <SectorIcon />
        </div>
      )}
      {props.sectorName && (
        <div
          className={
            "absolute top-0 w-full h-full flex justify-center items-center text-muted-foreground/50 z-0"
          }
        >
          {props.sectorName}
        </div>
      )}
    </div>
  );
};
