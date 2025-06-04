<script lang="ts">
  import { blur, draw, fly } from "svelte/transition";
  import { Get_current_stock } from "../wailsjs/go/main/App";
  import { simulation } from "../wailsjs/go/models";
  import Info from "./assets/images/Info.svelte";
  import { format_currency, format_number } from "./helper";
  import NumberInput from "./number_input.svelte";
  import Sidebar from "./Sidebar.svelte";
  import { company_id, decisions, external_factors } from "./store.svelte";
  import EmployeeCard from "./EmployeeCard.svelte";

  let page: string = $state("machines");
  let production_goal: number = $derived(decisions.Production.Production_goal);
  let current_stock = $state(0);
  let total_local_storage_capacity = $derived(get_local_storage_capacity());
  let available_storage_capacity = $derived(
    total_local_storage_capacity -
      (decisions.Production.Production_goal -
        decisions.Predictions.Sales_prediction +
        current_stock),
  );
  let reqired_storage_capacity = $derived(
    decisions.Production.Production_goal -
      decisions.Predictions.Sales_prediction +
      current_stock,
  );

  if (decisions.Production.Machines === null) {
    decisions.Production.Machines = [];
  }

  console.log(decisions);

  (async () => {
    current_stock = await Get_current_stock($company_id);
  })();

  if (decisions.Production.Logistics === null)
    decisions.Production.Logistics = [];

  function buy_machine() {
    // TODO: make machines cost money

    let bought_machine = external_factors.Machine_on_offer;
    bought_machine.Assigned_workers_ids = [];

    try {
      decisions.Production.Machines.push(bought_machine);
    } catch (exception) {
      if (exception instanceof TypeError) {
        decisions.Production.Machines = [bought_machine];
      } else {
        throw exception;
      }
    }
    console.log("Machine Bought");
  }
  function buy_warehouse() {
    // TODO: make warehouse cost money
    console.log("Warehouse Bought");
    const warehouse = {
      Id: decisions.Production.Logistics
        ? decisions.Production.Logistics.length
        : 1,
      Capacity: 10000,
      Operating_costs: 7000,
      Value: 10000,
    };
    try {
      decisions.Production.Logistics.push(warehouse);
    } catch (exception) {
      if (exception instanceof TypeError) {
        decisions.Production.Logistics = [warehouse];
      } else {
        throw exception;
      }
    }
  }

  function sell_machine(machine_index: number) {
    // TODO: Get money back when selling
    console.log("machine_sold");
    decisions.Production.Machines.splice(machine_index, 1);
  }

  function sell_warehouse(warehouse_index: number) {
    // TODO: Get money back when selling
    console.log("warehouse sold");
    decisions.Production.Logistics.splice(warehouse_index, 1);
  }

  function find_employee_by_id(id: number): simulation.Employee {
    for (let e of decisions.Employees) {
      if (e.Id == id) {
        return e;
      }
    }
    throw new Error(`no employee with id ${id} found`);
  }

  function calculate_production_capacity(): number {
    let total_production_capacity: number = 0;
    for (let m of decisions.Production.Machines) {
      total_production_capacity += m.Production_capacity;
    }

    return total_production_capacity;
  }

  function get_local_storage_capacity(): number {
    let local_storage_capacity = 0;
    for (let w of decisions.Production.Logistics) {
      local_storage_capacity += w.Capacity;
    }
    return local_storage_capacity;
  }

  function local_storage_cost(): number {
    let local_storage_capacity = get_local_storage_capacity();

    if (local_storage_capacity == 0) return 0;
    let local_storage_cost = 0;
    for (let w of decisions.Production.Logistics) {
      local_storage_cost += w.Operating_costs;
    }

    return local_storage_cost / local_storage_capacity;
  }
</script>

<Sidebar
  buttons={[
    {
      Text: "Machines",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "machines";
      },
      selected_by_default: true,
    },
    {
      Text: "Logistics",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "logistics";
      },
    },
    {
      Text: "Graphs",
      Style: "",
      Show: 0,
      onClick: () => {
        page = "graphs";
      },
    },
  ]}
  expand={true}
  horisontal={true}
  keep_pressed={true}
></Sidebar>
{#if page == "machines"}
  <div
    class="grid_container"
    style="grid-template-columns: 70% 30%; height: calc(100% - 50px); grid-template-rows: auto 100%; overflow-y: scroll;"
  >
    <div
      style="text-align: left; grid-column: 1; padding: 10px; border-bottom: var(--border-thin); border-right: var(--border-thin); height: fit-content;"
    >
      Production Goal:

      <h2>
        <NumberInput
          formatter={(value) => {
            return format_number(value, false, 0);
          }}
          bind:value={production_goal}
        ></NumberInput>
      </h2>
      <br />
      Base Production Capacity: {calculate_production_capacity()}
      <br />
      {#if calculate_production_capacity() < production_goal}
        <span style="color: var(--red)"
          >Not enough machines to achieve goal</span
        >
      {/if}
    </div>
    <div style="grid-row: 1 / span 2; grid-column: 2;"></div>
    <div
      style="grid-row: 2; grid-column: 1; display: flex; flex-direction: column; border-right: var(--border-thin); height: fit-content; min-height: 100%;"
    >
      <p>
        <input id="manually_assign_workers" type="checkbox" />
        <label for="manually_assign_workers">Manually Assign Workers</label>
      </p>
      <div style="display: flex;flex-direction: column;">
        <div
          style="display: flex; border: var(--border-thin); border-radius: var(--border-radius); margin: 10px;"
        >
          <div style="  padding: 10px; flex: 1 0 40%; text-align: left;">
            <p>Machine</p>
            <small>Base Production Capacity</small>
            <p>
              {external_factors.Machine_on_offer.Production_capacity} items/month
            </p>
            <small>Min / Max Required Employees</small>
            <p>
              {external_factors.Machine_on_offer.Minimum_workers} / {external_factors
                .Machine_on_offer.Required_workers}
            </p>
            <small>Energy Use</small>
            <p>{external_factors.Machine_on_offer.Energy_use} kWh/month</p>
            <button
              style="width: 100%; background-color: var(--green);"
              onclick={() => {
                let dialogue: HTMLDialogElement =
                  document.getElementById("buy_machine");
                dialogue.showModal();
              }}
              >Buy: {format_currency(
                external_factors.Machine_on_offer.Value,
                0,
              )}</button
            >
          </div>

          <dialog id="buy_machine">
            <article>
              Are you sure you want to buy a new machine for {format_currency(
                external_factors.Machine_on_offer.Value,
              )}
            </article>
            <footer style="display: flex; gap: 10px;">
              <button
                onclick={() => {
                  let dialogue: HTMLDialogElement =
                    document.getElementById("buy_machine");
                  dialogue.close();
                }}>No</button
              ><button
                onclick={() => {
                  let dialogue: HTMLDialogElement =
                    document.getElementById("buy_machine");
                  dialogue.close();
                  buy_machine();
                }}>Yes</button
              >
            </footer>
          </dialog>
        </div>
        {#each decisions.Production.Machines as m, m_idnex}
          <div
            transition:fly
            style="display: flex; border: var(--border-thin); border-radius: var(--border-radius); margin: 10px;"
          >
            <div style="  padding: 10px; flex: 1 0 40%; text-align: left;">
              <p>Machine</p>
              <small>Base Production Capacity</small>
              <p>{m.Production_capacity} items/month</p>
              <small>Min / Max Required Employees</small>
              <p>{m.Minimum_workers} / {m.Required_workers}</p>
              <small>Energy Use</small>
              <p>{m.Energy_use} kWh/month</p>
              <button
                style="width: 100%; background-color: var(--red);"
                onclick={() => {
                  let dialogue: HTMLDialogElement = document.getElementById(
                    `sell_machine_${m_idnex}`,
                  );
                  dialogue.showModal();
                }}>Sell: {format_currency(m.Value, 0)}</button
              >
            </div>

            <!--<div
              style="display: flex; flex-wrap: wrap; gap: 10px; margin: 10px; overflow-y: scroll; max-height: 100%;"
            >
              {#each m.Assigned_workers_ids as w_id}
                <div>
                  <h4>{find_employee_by_id(w_id).Id}</h4>
                </div>
              {/each}
              {#each Array(m.Required_workers - m.Assigned_workers_ids.length).keys() as _}
                <EmployeeCard employee_data={undefined}></EmployeeCard>
              {/each}
            </div>-->
          </div>

          <dialog id="sell_machine_{m_idnex}">
            <article>
              Are you sure you want to sell this machine for {format_currency(
                m.Value,
              )}
            </article>
            <footer style="display: flex; gap: 10px;">
              <button
                onclick={() => {
                  let dialogue: HTMLDialogElement = document.getElementById(
                    `sell_machine_${m_idnex}`,
                  );
                  dialogue.close();
                }}>No</button
              ><button
                onclick={() => {
                  let dialogue: HTMLDialogElement = document.getElementById(
                    `sell_machine_${m_idnex}`,
                  );
                  dialogue.close();
                  sell_machine(m_idnex);
                }}>Yes</button
              >
            </footer>
          </dialog>
        {/each}
      </div>
    </div>
  </div>
{:else if page == "logistics"}
  <div
    class="grid_container"
    style="grid-template-columns: 70% 30% ;height: calc(100% - 50px);overflow-y: scroll;"
  >
    <div
      class="grid_item"
      style="text-align: left; grid-row: 1; grid-column: 1; border-right: var(--border-thin);"
    >
      <p>
        <Info />
        Whatever products your company cannot sell during a month will have to be
        stored in a warehouse.
        <br />
        <br />
        If you have no more available warehouse space your company will rent external
        storage (at a high cost).
      </p>
      <br />
      <small
        >Avg. local storage cost: {format_currency(local_storage_cost())} /item/month
        <br />
        External storage cost: {format_currency(
          external_factors.External_storage_price,
        )} /item/month
      </small>
    </div>
    <div
      class="grid_item"
      style="text-align: left; grid-row: 1; grid-column: 2;"
    >
      <small>
        Production goal: {format_number(decisions.Production.Production_goal)}
        <br />
        - Sales estimate: {format_number(
          decisions.Predictions.Sales_prediction,
        )}
        <br />
        + Existing stock: {format_number(current_stock)}
      </small>
      <p>
        = {format_number(reqired_storage_capacity)}
      </p>
      <small
        >Est. Required warehouse capacity
        <br /> Total capacity: {format_number(total_local_storage_capacity)}
        <br />Available capacity:
        <span
          style="color: {available_storage_capacity >= reqired_storage_capacity
            ? 'var(--green);'
            : 'var(--red);'};">{format_number(available_storage_capacity)}</span
        >
      </small>
    </div>
    <div
      class="grid_item"
      style="display: flex; flex-direction: row; flex-wrap: wrap; grid-column: 1 / span 2; border-top: var(--border-thin); gap: 10px;"
    >
      <div
        style="border: var(--border-thin); border-radius: var(--border-radius); padding: 10px; flex: 1 1 45%; max-width: calc(50% - 30px); height: fit-content;"
      >
        <h3>New Warehouse</h3>
        <br />
        <table>
          <tbody>
            <tr>
              <td style="text-align: left;">Capacity:</td>
              <td style="text-align: right;">10'000</td>
            </tr>
            <tr>
              <td style="text-align: left;">Operating costs:</td>
              <td style="text-align: right;">7'000 CHF / month</td></tr
            >
          </tbody>
        </table>
        <br />
        <button
          onclick={() => {
            let dialogue: HTMLDialogElement =
              document.getElementById("buy_new");
            dialogue.showModal();
          }}
          style="background-color: var(--green);"
        >
          Buy: 10'000 CHF</button
        >
        <dialog id="buy_new">
          <article>
            Are you sure you want to buy a new warehouse for 10'000 CHF
          </article>
          <footer style="display: flex; gap: 10px;">
            <button
              onclick={() => {
                let dialogue: HTMLDialogElement =
                  document.getElementById("buy_new");
                dialogue.close();
              }}>No</button
            ><button
              onclick={() => {
                let dialogue: HTMLDialogElement =
                  document.getElementById("buy_new");
                dialogue.close();
                buy_warehouse();
              }}>Yes</button
            >
          </footer>
        </dialog>
      </div>
      {#each decisions.Production.Logistics as warehouse, w_index}
        <div
          transition:fly
          style="border: var(--border-thin); border-radius: var(--border-radius); padding: 10px; flex: 1 1 45%; max-width: calc(50% - 30px); height: fit-content;"
        >
          <h3>Warehouse</h3>
          <br />
          <table>
            <tbody>
              <tr>
                <td style="text-align: left;">Capacity:</td>
                <td style="text-align: right;"
                  >{format_number(warehouse.Capacity)}</td
                >
              </tr>
              <tr>
                <td style="text-align: left;">Operating costs:</td>
                <td style="text-align: right;"
                  >{format_currency(warehouse.Operating_costs)} / month</td
                ></tr
              >
            </tbody>
          </table>
          <br />
          <button
            onclick={() => {
              let dialogue: HTMLDialogElement = document.getElementById(
                `sell_${w_index}`,
              );
              dialogue.showModal();
            }}
            style="background-color: var(--red);"
            >Sell: {format_currency(warehouse.Value)}</button
          >
          <dialog id="sell_{w_index}">
            <article>
              Are you sure you want to sell this warehouse for {format_currency(
                warehouse.Value,
              )}
            </article>
            <footer style="display: flex; gap: 10px;">
              <button
                onclick={() => {
                  let dialogue: HTMLDialogElement = document.getElementById(
                    `sell_${w_index}`,
                  );
                  dialogue.close();
                }}>No</button
              ><button
                onclick={() => {
                  let dialogue: HTMLDialogElement = document.getElementById(
                    `sell_${w_index}`,
                  );
                  dialogue.close();
                  sell_warehouse(w_index);
                }}>Yes</button
              >
            </footer>
          </dialog>
        </div>
      {/each}
    </div>
  </div>
{:else if page == "graphs"}{/if}

<style>
  small,
  p {
    padding-left: 0px;
  }

  p {
    padding-bottom: 0.5rem;
  }

  table {
    width: 100%;
  }

  div {
    text-align: left;
  }
  button {
    width: 100%;
    padding: 10px;
  }
</style>
