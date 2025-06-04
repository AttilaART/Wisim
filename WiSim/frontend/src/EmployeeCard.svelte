<script lang="ts">
  import { simulation } from "../wailsjs/go/models";

  let { employee_data }: { employee_data?: simulation.Employee | undefined } =
    $props();
</script>

<div
  style="border-radius: var(--border-radius); border: var(--border-thin); overflow: hidden;; display: flex; flex-direction: rows; width: 200px; height: 75px;"
>
  {#if employee_data}
    <div style="flex: 0 0 50px;">img</div>
    <div style="flex: 1 1 150px; padding: 2px 10px;">
      <h5>{employee_data.Name}</h5>
      <table style="width: 100%;">
        <tbody>
          <tr>
            <td> <small>Skill</small></td>
            <td> <small>Motivation</small></td>
          </tr>

          <tr>
            <td>{employee_data.Skill}</td>
            <td>{employee_data.Motivation}</td>
          </tr>
        </tbody>
      </table>
    </div>
  {:else}
    <div style="flex: 0 0 50px; border: var(--border);">img</div>
    <div style="flex: 1 1 150px; padding: 5px 10px;">
      <h1 style="text-align: center;">?</h1>
    </div>
  {/if}
</div>
