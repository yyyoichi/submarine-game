import { ControllerComponent } from "@/components/controller";
import { FadeInOutTrigger } from "@/components/fadeinout";
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
import React, { useEffect, useState } from "react";

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
    "Direction" | "AButton" | "BButton"
  >;
};

export const PreparingPage = (props: PreparingPageProps) => {
  const [viewGuide, setViewGuide] = useState(false);
  const [api, setApi] = useState<CarouselApi>();
  const [leadingApi, setLeadingApi] = useState<CarouselApi>();
  const [oceanApi, setOceanApi] = useState<CarouselApi>();
  const meGuideRef = React.useRef<HTMLDivElement>(null);
  const minesGuideRef = React.useRef<HTMLDivElement>(null);
  const [guideContentWidth, setGuideContentWidth] = useState<
    number | undefined
  >(undefined);
  useEffect(() => {
    if (!api || !leadingApi || !oceanApi) return;
    switch (props.preparingStep) {
      case "me":
        api.scrollTo(0);
        leadingApi.scrollTo(0);
        oceanApi.scrollTo(0);
        break;
      case "mines":
        leadingApi.scrollTo(1);
        oceanApi.scrollTo(1);
        break;
      case "done":
        api.scrollTo(1);
        break;
    }
  }, [api, leadingApi, oceanApi, props.preparingStep]);
  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (!meGuideRef.current || !minesGuideRef.current) return;
    switch (props.preparingStep) {
      case "me":
        setGuideContentWidth(meGuideRef.current.scrollWidth);
        break;
      case "mines":
        setGuideContentWidth(minesGuideRef.current.scrollWidth);
        break;
      case "done":
        setGuideContentWidth(undefined);
        break;
    }
  }, [meGuideRef.current, minesGuideRef.current, props.preparingStep]);
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
    children: "敵作戦海域に到達しました！",
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

  const leadingProps: React.ComponentProps<typeof LeadingComponent> = {
    ...props.Leading,
    inMyTrun: true,
    contentWidth: guideContentWidth,
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
    CommandWindow: {
      commands: ["deploy-me", "deploy-mines", "deploy-wait"],
      use:
        props.preparingStep === "me"
          ? "deploy-me"
          : props.preparingStep === "mines"
            ? "deploy-mines"
            : "deploy-wait",
    },
  };

  return (
    <div className="flex flex-col">
      <Carousel setApi={setApi} opts={opts}>
        <CarouselContent className="mx-0 px-0">
          {/* api.0 行動開始位置 */}
          <CarouselItem className="pl-0">
            {/* Leading */}
            <div className="relative w-full">
              <Carousel setApi={setLeadingApi} opts={{ ...opts }}>
                <LeadingComponent {...leadingProps}>
                  <CarouselContent className="mx-0 px-0">
                    {/* leadingApi.0 */}
                    <CarouselItem className="pl-0">
                      <div className="w-fit" ref={meGuideRef}>
                        {"行動開始海域を決定してください。"}
                      </div>
                    </CarouselItem>
                    {/* oceanApi.1 */}
                    <CarouselItem className="pl-0">
                      <div className="w-fit" ref={minesGuideRef}>
                        {"制御機雷を敷設してください。"}
                      </div>
                    </CarouselItem>
                  </CarouselContent>
                </LeadingComponent>
              </Carousel>
              <OverlayedLeadingComponent {...overlayedLeadingProps} />
            </div>
            {/* Ocean */}
            <div className="relative w-full">
              <Carousel setApi={setOceanApi} opts={opts}>
                <CarouselContent className="mx-0 px-0">
                  {/* oceanApi.1 */}
                  <CarouselItem className="pl-0">
                    <OceanComponent {...meOceanProps} />
                  </CarouselItem>
                  {/* oceanApi.2 */}
                  <CarouselItem className="pl-0">
                    <OceanComponent {...minesOceanProps} />
                  </CarouselItem>
                </CarouselContent>
              </Carousel>
              <OverlayedOceanComponent {...overlayedOceanProps} />
            </div>
          </CarouselItem>
          {/* api.1 待機*/}
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
  const [trigger, setTrigger] = useState("i");
  useEffect(() => {
    const timer = setInterval(() => {
      setTrigger(String(Math.random() * 100));
    }, 4000);
    return () => {
      clearTimeout(timer);
    };
  }, []);
  return (
    <>
      <FadeInOutTrigger trigger={trigger} duration={1000}>
        <div className="w-full h-10 bg-red-500" />
      </FadeInOutTrigger>
    </>
  );
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
