import { ControllerComponent } from "@/components/controller";
import {
  LeadingComponent,
  OverlayedLeadingComponent,
} from "@/components/leading";
import { OceanComponent, OverlayedOceanComponent } from "@/components/ocean";
import type React from "react";

type PreparingPageProps = {
  Leading: React.ComponentProps<typeof LeadingComponent>;
  OverlayedLeading: React.ComponentProps<typeof OverlayedLeadingComponent>;
  Ocean: React.ComponentProps<typeof OceanComponent>;
  OverlayedOcean: React.ComponentProps<typeof OverlayedOceanComponent>;
  Controller: React.ComponentProps<typeof ControllerComponent>;
};

export const PreparingPage = (props: PreparingPageProps) => {
  return (
    <div className="flex flex-col">
      <div className="relative w-full">
        <LeadingComponent {...props.Leading} />
        <OverlayedLeadingComponent {...props.OverlayedLeading} />
      </div>
      <div className="relative w-full">
        <OceanComponent {...props.Ocean} />
        <OverlayedOceanComponent {...props.OverlayedOcean} />
      </div>
      <ControllerComponent {...props.Controller} />
      <div className="w-[200px] h-[200px]" />
    </div>
  );
};

export const PlayingPage = () => {};

export const FinishedPage = () => {};
