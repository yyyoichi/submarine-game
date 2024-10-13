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
    const interval = setInterval(() => {
      let v = 0;
      setValue(() => {
        v = calc();
        return v;
      });
      if (v > 0) return;
      props.callback();
      clearInterval(interval);
    }, 100); // 毎秒更新

    return () => clearInterval(interval);
  }, [calc, props.callback]);

  return <Progress hasStripe color={"red.500"} value={value} />;
}
