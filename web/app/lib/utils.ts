import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * A1, B1, C1, D1, ...
 * ....
 * A6, B6, C6, D6
 *
 *  のようにセクター位置を文字に変換する
 * @param base 1辺のセクターの数
 * @param i 左上から数えたセクターのインデックス
 * @returns セクター位置を表す文字列
 */
export function sector(base: number, i: number) {
  const x = i % base;
  const y = Math.floor(i / base);
  return String.fromCharCode(65 + x) + (y + 1);
}
