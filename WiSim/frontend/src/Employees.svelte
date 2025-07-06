<script lang="ts">
  import {
    Get_action_employee,
    Temp_generate_new_emloyee,
  } from "../wailsjs/go/main/App";
  import { simulation } from "../wailsjs/go/models";
  import { format_number } from "./helper.svelte";
  import NumberInput from "./number_input.svelte";
  import { current_decisions } from "./store.svelte";
  import { company_id, month } from "./store.svelte";

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
  ): Promise<simulation.Employee_action[]> {
    const stati = {
      existing_employee: 0,
      new_hire: 1,
      fired: 2,
    };

    let employee_actions: simulation.Employee_action[];
    if (type == 0) {
      employee_actions = current_decisions.Employees.Production_actions;
    } else if (type == 1) {
      employee_actions = current_decisions.Employees.Marketing_actions;
    }

    if (employee_actions == null) employee_actions = [];

    let delta: number = new_num_employees - employee_actions.length;

    if (delta > 0) {
      // hire employees
      while (delta != 0) {
        let employee_action: simulation.Employee_action =
          new simulation.Employee_action();
        employee_action.Employee_id = await Temp_generate_new_emloyee(
          $company_id,
          $month,
          type,
        );

        let employee: simulation.Employee =
          await Get_action_employee(employee_action);

        employee_action.Pay = employee.Pay;
        employee_action.Bonus = employee.Bonus;

        employee_action.Status = stati.new_hire;
        employee_actions.push(employee_action);

        delta -= 1;
      }
    } else if (delta < 0) {
      // fire employees
      for (let i in employee_actions) {
        if (delta == 0) break;

        employee_actions[i].Status = stati.fired;
        delta += 1;
      }
    }
    return employee_actions;
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
