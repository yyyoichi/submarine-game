import { Progress } from "@chakra-ui/react";
import { useCallback, useEffect, useState } from "react";
import { useLoaderData } from "react-router-dom";
import type { LogsResponse } from "../../gen/api/v2/game_pb";

type ProgressBarProps = {
  callback: () => void;
};
export function ProgressBar(props: ProgressBarProps) {
  const logs = useLoaderData() as LogsResponse;
  const calc = useCallback(() => {
    const r = Number(logs.timeout) - Date.now();
    if (r < 0) {
      return 0;
    }
    return (r / Number(logs.millSecondPerTurn)) * 100;
  }, [logs.timeout, logs.millSecondPerTurn]);
  const [value, setValue] = useState(calc());
  useEffect(() => {
    if (value === 0) {
      return;
    }
    if (logs.gameIsOver) {
      return;
    }
    const id = setTimeout(() => {
      const v = calc();
      if (v === 0) {
        props.callback();
      }
      setValue(v);
    }, 100);
    return () => clearTimeout(id);
  }, [logs.gameIsOver, calc, props.callback, value]);

  return <Progress hasStripe color={"red.500"} value={value} />;
}
