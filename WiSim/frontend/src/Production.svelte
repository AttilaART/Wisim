<script lang="ts">
  import { format_number } from "./helper";
  import NumberInput from "./number_input.svelte";
  import Sidebar from "./Sidebar.svelte";

  let page: string = $state("machines");
  let production_goal: number = $state();
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
  <div class="grid-container">
    <div>
      <NumberInput
        formatter={(value) => {
          return format_number(value, false, 0);
        }}
        bind:value={production_goal}
      ></NumberInput>
    </div>
    <div></div>
    <div></div>
  </div>
{:else if page == "logistics"}{:else if page == "graphs"}{/if}
