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

export class SectorDirection {
  /**
   * @param base 1辺のセクターの数
   */
  constructor(private base: number) {}
  /**
   * 起点から北方向にもっとも妥当なセクターを返す。存在しない場合は、-1を返す。
   * @param at 起点となるセクターインデックス
   * @param enables 利用可能なセクター
   */
  north(at: number, enables: number[]) {
    // 一つ上のセクターから探索を開始して、-方向に探索する
    return this.updown(at - this.base, -1, (q) => enables.includes(q));
  }
  /**
   * 起点から南方向にもっとも妥当なセクターを返す。存在しない場合は、-1を返す。
   * @param at 起点となるセクターインデックス
   * @param enables 利用可能なセクター
   */
  south(at: number, enables: number[]) {
    // 一つ下のセクターから探索を開始して、+方向に探索する
    return this.updown(at + this.base, 1, (q) => enables.includes(q));
  }

  /**
   * 起点から東方向にもっとも妥当なセクターを返す。存在しない場合は、-1を返す。
   * @param at 起点となるセクターインデックス
   * @param enables 利用可能なセクター
   */
  east(at: number, enables: number[]) {
    // 一つ右のセクターから探索を開始して、+方向に探索する
    return this.rightleft(Math.floor(at / this.base), at + 1, 1, (q) =>
      enables.includes(q),
    );
  }
  /**
   * 起点から西方向にもっとも妥当なセクターを返す。存在しない場合は、-1を返す。
   * @param at 起点となるセクターインデックス
   * @param enables 利用可能なセクター
   * */
  west(at: number, enables: number[]) {
    // 一つ左のセクターから探索を開始して、-方向に探索する
    return this.rightleft(Math.floor(at / this.base), at - 1, -1, (q) =>
      enables.includes(q),
    );
  }

  private updown(to: number, step: number, is: (q: number) => boolean): number {
    const y = Math.floor(to / this.base); // 縦方向の位置
    if (y < 0 || this.base <= y) return -1; // 域外
    if (is(to)) return to;
    // 横方向に探索
    for (let i = 1; i < this.base; i++) {
      if (Math.floor((to + i) / this.base) === y) {
        if (is(to + i)) return to + i; // 東
      }
      if (Math.floor((to - i) / this.base) === y) {
        if (is(to - i)) return to - i; // 西
      }
    }
    // 縦方向に再探索
    return this.updown(step * this.base + to, step, is);
  }

  private rightleft(
    y: number,
    to: number,
    step: number,
    is: (q: number) => boolean,
  ): number {
    if (Math.floor(to / this.base) !== y) return -1; // y座標が違う
    if (is(to)) return to;
    // 縦方向に探索
    for (let i = 1; i < this.base; i++) {
      if (to - i * this.base > 0) {
        if (is(to - i * this.base)) return to - i * this.base; // 北
      }
      if (to + i * this.base < this.base * this.base) {
        if (is(to + i * this.base)) return to + i * this.base; // 北
      }
    }
    // 横方向に再探索
    return this.rightleft(y, step + to, step, is);
  }
}
