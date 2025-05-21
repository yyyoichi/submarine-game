import { ControllerComponent } from "@/components/controller";
import { OceanComponent } from "@/components/ocean";
import type React from "react";

type PreparingPageProps = {
  Ocean: React.ComponentProps<typeof OceanComponent>;
  Controller: React.ComponentProps<typeof ControllerComponent>;
};

export const PreparingPage = (props: PreparingPageProps) => {
  return (
    <div className="flex flex-col">
      <OceanComponent {...props.Ocean} />
      <ControllerComponent {...props.Controller} />
      <div className="w-[200px] h-[200px]" />
    </div>
  );
};

export const PlayingPage = () => {};

export const FinishedPage = () => {};
