<script lang="ts">
  import { fade, fly } from "svelte/transition";
  import {
    month_counter,
    company,
    decisions as decisions_store,
  } from "./store";
  import {
    Get_past_decisions,
    Get_availible_decisions,
  } from "../wailsjs/go/main/App";
  import { simulation } from "../wailsjs/go/models";

  type Decision = {
    name: string;
    info: string;
    two_months_ago: number;
    last_month: number;
    value: number;
  };

  let decisions: Decision[] = $state([]);
  let is_error: boolean = $state(false);
  let error = $state();

  get_decisions();

  async function get_decisions() {
    is_error = false;
    let availible_decisions = $decisions_store;
    let last_month_decisions;
    try {
      if ($month_counter - 1 >= 0) {
        last_month_decisions = await Get_past_decisions(
          $company,
          $month_counter - 1,
        );
      } else {
        last_month_decisions = availible_decisions;
      }
    } catch (exception) {
      is_error = true;
      error = exception;
    }

    try {
      let two_months_ago_decisions;
      if ($month_counter - 2 >= 0) {
        two_months_ago_decisions = await Get_past_decisions(
          $company,
          $month_counter - 2,
        );
      } else {
        two_months_ago_decisions = availible_decisions;
      }
    } catch (exception) {
      is_error = true;
      error = exception;
    }

    let key = Object.keys(availible_decisions);
    let current = Object.values(availible_decisions);
    let last_month = Object.values(last_month_decisions);
    let two_months_ago = Object.values(last_month_decisions);

    for (let i in key) {
      decisions.push({
        name: key[i],

        info: "",
        two_months_ago: Number(two_months_ago[i]),
        last_month: Number(last_month[i]),
        value: Number(current[i]),
      });
    }
  }

  function submit_decisions(decisions: Decision[]) {
    let to_submit = new simulation.Decisions();
    console.log(decisions);
    for (let i in decisions) {
      to_submit[decisions[i].name] = decisions[i].value;
    }
    $decisions_store = to_submit;
  }
</script>

<div
  out:fade={{ duration: 300 }}
  in:fly={{ duration: 300, delay: 300, y: -40 }}
>
  <div class="grid_item">
    <h2>Decisions</h2>
    <p>Here you have all the instruments needed to control your company</p>
  </div>
  <div class="grid_item">
    <table>
      <thead>
        <tr>
          <th style="width: auto; text-align: left;">Decision</th>
          <th style="width: auto; text-align: left;">Info</th>
          <th style="width: 20%; text-align: right;">
            {$month_counter % 13} / {Math.floor(($month_counter - 1) / 12)}
          </th>
          <th style="width: 20%; text-align: right;"
            >{($month_counter + 1) % 13} / {Math.floor($month_counter / 12)}</th
          >
          <th style="width: 20%; text-align: right;"
            >{($month_counter + 2) % 13} / {Math.floor(
              ($month_counter + 1) / 12,
            )}
          </th>
        </tr>
      </thead>
      <tbody>
        {#each decisions as _, decision_index}
          {@render decisions_row(decision_index)}
        {/each}
      </tbody>
    </table>
  </div>
</div>

{#snippet decisions_row(decision_index: number)}
  <tr>
    <td style="text-align: left;"
      >{decisions[decision_index].name.replaceAll("_", " ")}</td
    >
    <td style="text-align: left;">{decisions[decision_index].info}</td>
    <td style="text-align: right;"
      >{decisions[decision_index].two_months_ago}</td
    >
    <td style="text-align: right;">{decisions[decision_index].last_month}</td>
    <td style="text-align: right;">
      <input
        class="input_field"
        type="number"
        bind:value={decisions[decision_index].value}
        onchange={() => submit_decisions(decisions)}
      />
    </td>
  </tr>
{/snippet}

<style>
  .input_field {
    -webkit-appearance: none;
    background-color: rgba(0, 0, 0, 0.2);
    border-width: 0;
    border-radius: 3px;
    color: white;
    padding: 2px 2px 2px 8px;
    size: 10;
  }

  .custom-arrow {
    position: relative;
    top: 0;
    left: -7px;
    display: block;
    height: 100%;
    pointer-events: none;
  }

  .custom-arrow::before,
  .custom-arrow::after {
    content: "";
    --size: 7px;
    position: absolute;
    width: 0;
    height: 0;
    top: calc(50% - 0px);
  }

  .custom-arrow.down::before {
    border-left: var(--size) solid transparent;
    border-right: var(--size) solid transparent;
    border-top: var(--size) solid rgba(255, 255, 2255, 0.5);
    transform: translateY(calc(50% - 6px));
  }

  .custom-arrow.up::after {
    border-left: var(--size) solid transparent;
    border-right: var(--size) solid transparent;
    border-bottom: var(--size) solid rgba(255, 255, 2255, 0.5);
    transform: translateY(calc(-50% - 0px));
  }
</style>
