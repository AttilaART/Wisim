import { Get_Decisions, Get_External_Factors, Get_reports, Submit_decisions, Trigger_simulation } from "../wailsjs/go/main/App";
import { simulation } from "../wailsjs/go/models";
import { company_id, latest_reports, decisions, external_factors, month } from "./store.svelte";
import { get } from "svelte/store"

export async function trigger_simulation(force?: boolean) {
  await Submit_decisions(get(company_id), decisions)
  month.set(await Trigger_simulation(Boolean(force)))
  update_external_factors(await Get_External_Factors());
  update_decisions(await Get_Decisions(get(company_id), get(month)));
  update_reports(await Get_reports(get(company_id), get(month) - 1))
}

export function isEqual(a: unknown, b: unknown): boolean {
  return JSON.stringify(a) == JSON.stringify(b)
}

export function update_external_factors(new_value: simulation.External_factors) {
  for (let field of Object.keys(new_value)) {
    external_factors[field] = new_value[field]
  }
}

export function update_decisions(new_value: simulation.Decisions) {
  for (let field of Object.keys(new_value)) {
    console.log(field)
    decisions[field] = new_value[field]
  }
}

export function update_reports(new_value: simulation.Report) {
  for (let field of Object.keys(new_value)) {
    latest_reports[field] = new_value[field]
  }
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

export function format_currency(num: number, decimal_places?: number, add_plus?: boolean): string {
  if (decimal_places == undefined) {
    decimal_places == 0
  }
  if (add_plus === undefined) {
    add_plus = false
  }

  return `${format_number(num, add_plus, decimal_places)} CHF`
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
