<script lang="ts">
  import {
    Get_unemployed,
    Temp_generate_new_emloyee,
  } from "../wailsjs/go/main/App";
  import { simulation } from "../wailsjs/go/models";
  import { format_number } from "./helper.svelte";
  import NumberInput from "./number_input.svelte";
  import { current_decisions } from "./store.svelte";

  let unapplied_changes: boolean = $state(false);

  console.log(current_decisions);

  let num_production_employees: number = $state(
    null_to_zero(current_decisions.Employees.Production_actions),
  );
  let num_marketing_employees: number = $state(
    null_to_zero(current_decisions.Employees.Marketing_actions),
  );

  function null_to_zero(item: null | any[]) {
    if (item == null) return 0;
    else return item.length;
  }

  async function temp_hire_fire(
    new_num_employees: number,
    type: number,
  ): Promise<simulation.Delta_WiSim_simulation_Employee_[]> {
    if (new_num_employees == 0) {
      return [];
    } else if (new_num_employees > 0) {
      let unemployed_pool: simulation.Employee[] = await Get_unemployed(0);
      let new_hires: simulation.Delta_WiSim_simulation_Employee_[] = [];
      for (let i in unemployed_pool) {
        if (Number(i) > new_num_employees) {
          break;
        }

        let new_hire: simulation.Delta_WiSim_simulation_Employee_ =
          new simulation.Delta_WiSim_simulation_Employee_();
        new_hire.Change = 0;
        new_hire.Item = unemployed_pool[i];
        new_hires.push(new_hire);
      }
    }
  }

  async function apply() {
    current_decisions.Employees.Production_actions = await temp_hire_fire(
      num_production_employees,
      0,
    );

    current_decisions.Employees.Marketing_actions = await temp_hire_fire(
      num_marketing_employees,
      1,
    );

    console.log(current_decisions.Employees.Production_actions);
    console.log(current_decisions.Employees.Marketing_actions);
  }

  $effect(() => {
    if (
      null_to_zero(current_decisions.Employees.Production_actions) !=
      num_production_employees
    )
      unapplied_changes = true;
    else if (
      null_to_zero(current_decisions.Employees.Marketing_actions) !=
      num_marketing_employees
    )
      unapplied_changes = true;
    else unapplied_changes = false;
  });
</script>

<div>
  <p>Num production employees</p>
  <h2>
    <NumberInput
      formatter={(value) => {
        return format_number(value, false, 0);
      }}
      bind:value={num_production_employees}
    ></NumberInput>
  </h2>

  <p>Num marketing employees</p>
  <h2>
    <NumberInput
      formatter={(value) => {
        return format_number(value, false, 0);
      }}
      bind:value={num_marketing_employees}
    ></NumberInput>
  </h2>

  <button class={unapplied_changes ? "" : "greyed_out"} onclick={apply}>
    Apply
  </button>
</div>
