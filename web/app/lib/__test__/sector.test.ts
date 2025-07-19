import { describe, expect, it } from "vitest";
import { SectorDirection } from "../utils";

describe("SectorDirection", () => {
  it("should return correct directions", () => {
    const base6 = new SectorDirection(6);
    const enables = Array.from({ length: 6 * 6 }, (_, i) => i);
    // 中央から
    expect(base6.north(10, enables)).toBe(4);
    expect(base6.south(10, enables)).toBe(16);
    expect(base6.east(10, enables)).toBe(11);
    expect(base6.west(10, enables)).toBe(9);
    // 折り返しあり
    expect(base6.north(0, enables)).toBe(-1);
    expect(base6.south(0, enables)).toBe(6);
    expect(base6.east(0, enables)).toBe(1);
    expect(base6.west(0, enables)).toBe(-1);
    expect(base6.north(35, enables)).toBe(29);
    expect(base6.south(35, enables)).toBe(-1);
    expect(base6.east(35, enables)).toBe(-1);
    expect(base6.west(35, enables)).toBe(34);
  });
});
