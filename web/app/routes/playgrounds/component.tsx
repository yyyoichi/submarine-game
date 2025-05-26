import { ControllerComponent } from "@/components/controller";
import { FadeInOutTrigger } from "@/components/fadeinout";
import {
  LeadingComponent,
  OverlayedLeadingComponent,
  OverlayedLeadingComponentAbsoluteClassName,
} from "@/components/leading";
import {
  OceanComponent,
  OverlayedOceanComponent,
  OverlayedOceanComponentAbsoluteClassName,
} from "@/components/ocean";
import {
  Carousel,
  type CarouselApi,
  CarouselContent,
  CarouselItem,
} from "@/components/ui/carousel";
import { cn } from "@/lib/utils";
import React, { useEffect, useLayoutEffect, useState } from "react";

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

type PlayingPageProps = {
  turn: number;
  // ゲーム中のガイドはUIの責務対象外。
  Leading: React.ComponentProps<typeof LeadingComponent>;
  Ocean: Pick<React.ComponentProps<typeof OceanComponent>, "sectorCount"> & {
    islands: number[];
    me: number;
    enables: number[];
    sector?: number;
    type: "torpedo" | "mine" | "move";
  };
  // ターン切替時表示。ガイドはUIの責務対象外。
  OverlayedLeading: Pick<
    React.ComponentProps<typeof OverlayedLeadingComponent>,
    "children"
  >;
  // ターン切替時表示。タイトル、アニメーションはUIの責務対象外。
  OverlayedOcean: Pick<
    React.ComponentProps<typeof OverlayedOceanComponent>,
    "ringSector" | "gapSector" | "title" | "titlePosition"
  >;
  Controller: Pick<
    React.ComponentProps<typeof ControllerComponent>,
    "Direction" | "AButton" | "BButton" | "CommandWindow"
  >;
};

export const PlayingPage = (props: PlayingPageProps) => {
  const [viewGuide, setViewGuide] = useState(false);
  const oceanProps: React.ComponentProps<typeof OceanComponent> = {
    displaySectorName: viewGuide,
    sectorCount: props.Ocean.sectorCount,
    Sectors: {},
  };
  for (let i = 0; i < props.Ocean.sectorCount ** 2; i++) {
    const sct = oceanProps.Sectors[i] || {};
    if (props.Ocean.islands.includes(i) || !props.Ocean.enables.includes(i)) {
      sct.embed = true;
    }
    if (props.Ocean.me === i) {
      sct.icon = "me";
    }
    if (props.Ocean.sector === i) {
      sct.icon = props.Ocean.type;
    }
    oceanProps.Sectors[i] = sct;
  }
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

  const overlayedLeadingProps: React.ComponentProps<
    typeof OverlayedLeadingComponent
  > = {
    ...props.OverlayedLeading,
    absolute: true,
    fadeout: false, // FadeInOutTriggerで制御する
  };
  const overlayedOceanProps: React.ComponentProps<
    typeof OverlayedOceanComponent
  > = {
    // biome-ignore lint/style/noNonNullAssertion: <explanation>
    title: props.OverlayedOcean.title!,
    // biome-ignore lint/style/noNonNullAssertion: <explanation>
    titlePosition: props.OverlayedOcean.titlePosition!,
    ...props.OverlayedOcean,
    absolute: true,
    fadeout: false, // FadeInOutTriggerで制御する
    sectorCount: props.Ocean.sectorCount,
  };
  const fadeInOutTriggerProps: React.ComponentProps<typeof FadeInOutTrigger> = {
    trigger: `${props.turn}-${props.Leading.inMyTrun}`,
    fadein: 50,
    fadeout: 1000,
  };
  return (
    <div className="flex flex-col">
      <div className="relative w-full">
        <LeadingComponent {...props.Leading} />
        <FadeInOutTrigger
          {...fadeInOutTriggerProps}
          className={cn(OverlayedLeadingComponentAbsoluteClassName, "w-full")}
        >
          <OverlayedLeadingComponent {...overlayedLeadingProps} />
        </FadeInOutTrigger>
      </div>
      <div className="relative w-full">
        <OceanComponent {...oceanProps} />
        <FadeInOutTrigger
          {...fadeInOutTriggerProps}
          className={cn(OverlayedOceanComponentAbsoluteClassName, "w-full")}
        >
          <OverlayedOceanComponent {...overlayedOceanProps} />
        </FadeInOutTrigger>
      </div>
      <div>
        <ControllerComponent {...controllerPrpos} />
      </div>
    </div>
  );
};

type FinishedPageProps = {
  Leading: Pick<
    React.ComponentProps<typeof LeadingComponent>,
    "inMyTrun" | "children"
  >;
  Ocean: Pick<React.ComponentProps<typeof OceanComponent>, "sectorCount"> & {
    islands: number[];
    show: "top" | "bottom";
    showIndex: number;
  };
  OverlayedLeading: Pick<
    React.ComponentProps<typeof OverlayedLeadingComponent>,
    "children"
  >;
  OverlayedOcean: Pick<
    React.ComponentProps<typeof OverlayedOceanComponent>,
    "ringSector" | "gapSector" | "title" | "titlePosition"
  >;
  TopOcean: Array<FinishedPageOceanProps>;
  BottomOcean: Array<FinishedPageOceanProps>;
  Controller: Pick<
    React.ComponentProps<typeof ControllerComponent>,
    "Direction" | "AButton" | "BButton" | "CommandWindow"
  >;
};
type FinishedPageOceanProps =
  | {
      mode: "action";
      me: number;
      sector: number;
      type: "fire-torpedo" | "trigger-mine" | "move";
    }
  | {
      mode: "dummy";
    };

export const FinishedPage = (props: FinishedPageProps) => {
  const [viewGuide, setViewGuide] = useState(false);
  const [verticalApi, setVerticalApi] = useState<CarouselApi>();
  const [topOceanApi, setTopOceanApi] = useState<CarouselApi>();
  const [bottomOceanApi, setBottomOceanApi] = useState<CarouselApi>();
  const oceanRef = React.useRef<HTMLDivElement>(null);
  const [oceanContentHight, setOceanContentHeight] = useState(0);
  useEffect(() => {
    if (!verticalApi || !topOceanApi || !bottomOceanApi) return;
    if (props.Ocean.show === "top") {
      verticalApi.scrollTo(0);
    } else {
      verticalApi.scrollTo(1);
    }
    topOceanApi.scrollTo(props.Ocean.showIndex);
    bottomOceanApi.scrollTo(props.Ocean.showIndex);
  }, [
    verticalApi,
    topOceanApi,
    bottomOceanApi,
    props.Ocean.show,
    props.Ocean.showIndex,
  ]);
  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useLayoutEffect(() => {
    if (!oceanRef.current) return;
    setOceanContentHeight(oceanRef.current.scrollHeight);
  }, [oceanRef.current]);
  const leadingProps: React.ComponentProps<typeof LeadingComponent> = {
    ...props.Leading,
    Alerm: {
      useAlerm: false,
    },
  };
  const overlayedLeadingProps: React.ComponentProps<
    typeof OverlayedLeadingComponent
  > = {
    ...props.OverlayedLeading,
    fadeout: true,
    absolute: true,
  };
  const overlayedOceanProps: React.ComponentProps<
    typeof OverlayedOceanComponent
  > = {
    ...props.OverlayedOcean,
    // biome-ignore lint/style/noNonNullAssertion: <explanation>
    title: props.OverlayedOcean.title!,
    // biome-ignore lint/style/noNonNullAssertion: <explanation>
    titlePosition: props.OverlayedOcean.titlePosition!,
    absolute: true,
    fadeout: true,
    sectorCount: props.Ocean.sectorCount,
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
  // Ocean
  const islandSectors: React.ComponentProps<typeof OceanComponent>["Sectors"] =
    {};
  for (const i of props.Ocean.islands) {
    islandSectors[i] = {
      embed: true,
    };
  }
  const makeOceans = (
    p: FinishedPageOceanProps,
  ): React.ComponentProps<typeof OceanComponent> => {
    const sectors: React.ComponentProps<typeof OceanComponent>["Sectors"] = {
      ...islandSectors,
    };
    if (p.mode === "action") {
      sectors[p.me] = {
        embed: false,
        icon: "me",
      };
      sectors[p.sector] = {
        embed: false,
        icon:
          p.type === "fire-torpedo"
            ? "torpedo"
            : p.type === "trigger-mine"
              ? "mine"
              : "move",
      };
    }
    return {
      displaySectorName: viewGuide,
      sectorCount: props.Ocean.sectorCount,
      Sectors: sectors,
    };
  };
  return (
    <div className="flex flex-col">
      <div className="relative w-full">
        <LeadingComponent {...leadingProps} />
        <OverlayedLeadingComponent {...overlayedLeadingProps} />
      </div>
      <div className="relative w-full">
        {/* 縦方向 */}
        <Carousel setApi={setVerticalApi} orientation="vertical" opts={opts}>
          <CarouselContent
            className="mt-0"
            style={{ height: `${oceanContentHight}px` }}
          >
            <CarouselItem className="pt-0">
              {/* 1. Top */}
              <Carousel setApi={setTopOceanApi} opts={opts}>
                <CarouselContent className="mx-0 px-0">
                  {props.TopOcean.map(makeOceans).map((ocean, index) => (
                    <CarouselItem
                      className="relative w-full pl-0 h-fit"
                      // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
                      key={index}
                    >
                      <OceanComponent {...ocean} />
                    </CarouselItem>
                  ))}
                </CarouselContent>
              </Carousel>
            </CarouselItem>
            <CarouselItem className="pt-0">
              {/* 2. Bottom */}
              <Carousel setApi={setBottomOceanApi} opts={opts}>
                <CarouselContent className="mx-0 px-0">
                  {props.BottomOcean.map(makeOceans).map((ocean, index) => (
                    // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
                    <CarouselItem className="relative w-full pl-0" key={index}>
                      <OceanComponent {...ocean} />
                    </CarouselItem>
                  ))}
                </CarouselContent>
              </Carousel>
            </CarouselItem>
          </CarouselContent>
        </Carousel>
        <OverlayedOceanComponent {...overlayedOceanProps} ref={oceanRef} />
      </div>
      <div>
        <ControllerComponent {...controllerPrpos} />
      </div>
    </div>
  );
};

const opts: React.ComponentProps<typeof Carousel>["opts"] = {
  duration: 20,
  align: "center",
  slidesToScroll: 1,
  watchDrag: false,
  watchFocus: false,
  watchResize: false,
  watchSlides: false,
};
