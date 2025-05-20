import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * A1, A2, A3, A4, ...D1, D2, D3, D4 のようにセクター位置を文字に変換する
 * @param base 1辺のセクターの数
 * @param i 左上から数えたセクターのインデックス
 * @returns セクター位置を表す文字列
 */
export function sector(base: number, i: number) {
  const x = Math.floor(i / base);
  const y = i % base;
  return String.fromCharCode(65 + x) + (y + 1);
}
