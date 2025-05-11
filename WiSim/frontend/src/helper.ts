export function isEqual(a: unknown, b: unknown): boolean {
  return JSON.stringify(a) == JSON.stringify(b)
}

export function format_number(num: number | any, add_plus?: boolean, decimal_places?: number): string {
  if (decimal_places == undefined) {
    decimal_places == 2
  }

  if (typeof num != typeof 1) {
    return num.toString()
  }

  if (!add_plus || num <= 0) {
    return num.toLocaleString("de-CH", {
      maximumFractionDigits: decimal_places,
      minimumFractionDigits: decimal_places,
    });
  }

  return `+${num.toLocaleString("de-CH", { maximumFractionDigits: decimal_places, minimumFractionDigits: decimal_places })}`
}

export function format_currency(num: number, decimal_places?: number) {
  if (decimal_places == undefined) {
    decimal_places == 0
  }

  return `${format_number(num, false, decimal_places)} CHF`
}

export type Series = {
  Name: string;
  Value: number;
  Color: string;
};

export function capitalise_first_letter(val: string): string {
  return String(val).charAt(0).toUpperCase() + String(val).slice(1);
}

export function generateGradient(
  data: Series[],
  type: string,
  direction: string,
  total?: number,
): string {
  let gradient: string = `${type}(${direction}`;

  if (total == undefined) {
    total = 0
    for (let n of data) {
      total += n.Value
    }
  }

  let current_percentage = 0;
  for (let i in data) {
    gradient += `, ${data[i].Color} ${current_percentage}%, ${data[i].Color} ${current_percentage + (data[i].Value / total) * 100}%`;
    current_percentage += (data[i].Value / total) * 100;
  }
  gradient += `, rgba(0, 0, 0, 0) ${current_percentage}%, rgba(0, 0, 0, 0) 100%`;

  gradient += ")";
  return gradient;
}

export const red = "rgb(255, 128, 128)";
export const green = "rgb(128, 255, 128)";
