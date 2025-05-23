import { ControllerComponent } from "@/components/controller";
import {
  LeadingComponent,
  OverlayedLeadingComponent,
} from "@/components/leading";
import { OceanComponent, OverlayedOceanComponent } from "@/components/ocean";
import {
  Carousel,
  type CarouselApi,
  CarouselContent,
  CarouselItem,
} from "@/components/ui/carousel";
import type React from "react";
import { useEffect, useState } from "react";

type PreparingPageProps = {
  preparingStep: "me" | "mines" | "done";
  Ocean: Pick<React.ComponentProps<typeof OceanComponent>, "sectorCount"> & {
    islands: number[];
  };
  Leading: Pick<React.ComponentProps<typeof LeadingComponent>, "Alerm">;
  MeOcean: {
    sectors: number[];
  };
  MinesOcean: {
    sectors: number[];
  };
  Controller: Pick<
    React.ComponentProps<typeof ControllerComponent>,
    "Direction"
  >;
};

export const PreparingPage = (props: PreparingPageProps) => {
  const [viewGuide, setViewGuide] = useState(false);
  const [api, setApi] = useState<CarouselApi>();
  useEffect(() => {
    if (!api) return;
    const p =
      props.preparingStep === "me"
        ? 0
        : props.preparingStep === "mines"
          ? 1
          : 2;
    api.scrollTo(p);
  }, [api, props.preparingStep]);
  const islandSectors: React.ComponentProps<typeof OceanComponent>["Sectors"] =
    {};
  for (const i of props.Ocean.islands) {
    islandSectors[i] = {
      embed: true,
    };
  }

  const overlayedLeadingProps: React.ComponentProps<
    typeof OverlayedLeadingComponent
  > = {
    fadeout: true,
    absolute: true,
    GuideLine: {
      children: "敵作戦海域に到達しました！",
    },
  };
  const overlayedOceanProps: React.ComponentProps<
    typeof OverlayedOceanComponent
  > = {
    ...props.Ocean,
    fadeout: true,
    absolute: true,
    title: "作戦準備",
    titlePosition: "right",
  };

  const meLeadingProps: React.ComponentProps<typeof LeadingComponent> = {
    ...props.Leading,
    inMyTrun: true,
    GuideLine: {
      children: "行動を開始する海域を決定してください。",
    },
  };
  const meOceanProps: React.ComponentProps<typeof OceanComponent> = {
    ...props.Ocean,
    Sectors: props.MeOcean.sectors.reduce(
      (acc, s) => {
        acc[s] = {
          embed: false,
          icon: "me",
        };
        return acc;
      },
      { ...islandSectors },
    ),
    displaySectorName: viewGuide,
  };

  const minesLeadingProps: React.ComponentProps<typeof LeadingComponent> = {
    inMyTrun: true,
    GuideLine: {
      children: "制御機雷を敷設してください。",
    },
    ...props.Leading,
  };
  const minesOceanProps: React.ComponentProps<typeof OceanComponent> = {
    ...props.Ocean,
    Sectors: props.MinesOcean.sectors.reduce(
      (acc, s) => {
        acc[s] = {
          embed: false,
          icon: "mine",
        };
        return acc;
      },
      { ...islandSectors },
    ),
    displaySectorName: viewGuide,
  };

  const waitingOceanProps: React.ComponentProps<
    typeof OverlayedOceanComponent
  > = {
    ...props.Ocean,
    fadeout: false,
    title: "待機中",
    titlePosition: "right",
  };

  const controllerPrpos: React.ComponentProps<typeof ControllerComponent> = {
    ...props.Controller,
    Compass: {
      onClick: () => {
        setViewGuide((v) => !v);
      },
      onKeyDown: () => {}, // TODO
    },
    visibleDirection: viewGuide,
  };

  return (
    <div className="flex flex-col">
      <Carousel setApi={setApi} opts={opts}>
        <CarouselContent className="mx-0 px-0">
          {/* 0 行動開始位置 */}
          <CarouselItem className="pl-0">
            <div className="relative w-full">
              <LeadingComponent {...meLeadingProps} />
              <OverlayedLeadingComponent {...overlayedLeadingProps} />
            </div>
            <div className="relative w-full">
              <OceanComponent {...meOceanProps} />
              <OverlayedOceanComponent {...overlayedOceanProps} />
            </div>
          </CarouselItem>
          {/* 1 制御機雷敷設*/}
          <CarouselItem className="pl-0">
            <LeadingComponent {...minesLeadingProps} />
            <OceanComponent {...minesOceanProps} />
          </CarouselItem>
          {/* 2 待機*/}
          <CarouselItem className="pl-0">
            <div className="relative w-full">
              <OverlayedLeadingComponent />
            </div>
            <div className="relative w-full">
              <OverlayedOceanComponent {...waitingOceanProps} />
            </div>
          </CarouselItem>
        </CarouselContent>
      </Carousel>
      <div>
        <ControllerComponent {...controllerPrpos} />
      </div>
    </div>
  );
};

export const PlayingPage = () => {
  return <div className="w-full h-10 bg-red-500" />;
};

export const FinishedPage = () => {};

const opts: React.ComponentProps<typeof Carousel>["opts"] = {
  duration: 20,
  align: "center",
  slidesToScroll: 1,
  watchDrag: false,
  watchFocus: false,
  watchResize: false,
  watchSlides: false,
};
