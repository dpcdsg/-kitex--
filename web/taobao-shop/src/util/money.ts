/** 后端金额字段为「分」 */
export function fenToYuan(fen: number): string {
  return (fen / 100).toFixed(2);
}
