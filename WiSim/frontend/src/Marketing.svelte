<script>
  import Price from "./Price.svelte";
  import Product from "./Product.svelte";
  import Promotion from "./Promotion.svelte";
  import Sidebar from "./Sidebar.svelte";

  let page = $state("product");
</script>

<Sidebar
  buttons={[
    {
      Text: "Product",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "product";
      },
      selected_by_default: true,
    },
    {
      Text: "Price",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "price";
      },
    },
    {
      Text: "Place",
      Style: "",
      Show: 0,
      onClick: () => {
        page = "place";
      },
    },
    {
      Text: "Promotion",
      Style: "",
      Show: 1,
      onClick: () => {
        page = "promotion";
      },
    },
  ]}
  expand={true}
  horisontal={true}
  keep_pressed={true}
></Sidebar>
{#if page == "product"}
  <Product></Product>
{:else if page == "price"}
  <Price></Price>
{:else if page == "place"}
  <h2>Coming Soon</h2>
{:else}
  <Promotion></Promotion>
{/if}

<style>
</style>
