import { cn } from "@/lib/utils";

export const ConcentricRingsComponent = () => {
  const ss = [
    { b: "border-2", d: -1200 },
    { b: "border-4", d: -800 },
    { b: "border-2", d: -500 },
    { b: "border-4", d: -200 },
    { b: "border-1", d: 0 },
  ];
  return (
    <>
      {ss.map((s, i) => (
        <div
          // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
          key={i}
          className={cn(
            "absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2",
            "w-10 h-10 rounded-full border-teal-500 z-150",
            s.b,
          )}
          style={{
            animation: `ten-ping 3000ms forwards ease-in-out ${s.d}ms 1`,
          }}
        />
      ))}
    </>
  );
};
