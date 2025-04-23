
export function format_number(num: number | any, add_plus?: boolean, decimal_places: number = 2): string {
  if (typeof num != typeof 1) {
    return num.toString()
  }

  if (!add_plus || num <= 0) {
    return num.toLocaleString("de-CH", {
      maximumFractionDigits: decimal_places,
      minimumFractionDigits: decimal_places,
    });
  }

  return `+${num.toLocaleString("de-CH", { maximumFractionDigits: 2, minimumFractionDigits: 2 })}`
}


export const red = "rgb(255, 128, 128)";
export const green = "rgb(128, 255, 128)";
