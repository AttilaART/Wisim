import simulation from "./simulation.js"

/**
 * @callback eventHandler
 * @param {Event} event
 */

/**
 * Message sent to server
 * @typedef {{Method: string, IsResponse: boolean, Error: string, DataType: DataType, Data: any}} Message
 */

/**
 * @typedef {string} DataType
 */

const Methods = {
  Get_decisions: "gDecisions",
  Get_company: "gCompany",
  Get_external_factors: "gExternal_factors",

  Set_decisions: "sDecisions",
  Set_ready: "sReady",
  Set_unready: "sUnready",

  func_calculate_product_stats: "fProduct_stats",
  broadcast_chat: "bChat",
}

/**
 * @param {string} url - URL of simulation server, should be localhost:8000
 * @param {eventHandler} eventHandler - Function that handles every event 
 * @returns {{socket: WebSocket, state: gamestate, gDecisions: ()=>void, gCompany: ()=>void, gExternal_factors: ()=>void, sDecisions: (decisions: simulation.Decisions)=>void, sReady: ()=>void, sUnready: ()=>void, fProduct_stats: (decisions: simulation.Decisions_product, research: simulation.Decisions_research)=>void, bChat: (chat: string)=>void}}
  */
export function newConnection(url, eventHandler) {
  let websocket = new WebSocket(url)

  websocket.addEventListener("message", (event) => {
    let message = JSON.parse(event.data)
    console.log(message)
  })

  /**
   * @typedef {{company: simulation.Company, decisions: simulation.Decisions, external_factors: simulation.External_factors}} gamestate
   * @type {gamestate}
   */

  let state = {
    company: {
      Id: 0,
      Name: "",
      Balance: 0,
      Loans: 0,
      Bridge_loans: 0,
      Decision_history: [],
      Reports: [],
      Global_quality_factor: 0,
      Base_marketing_strength: 0,
      Offer: undefined, // placeholder for Offer
      Orders: 0,
      Warehouses: [],
      Items_in_storage: 0,
      Machines: []
    },
    decisions: {
      Predictions: {
        Sales_prediction: 0
      },
      Finances: {
        Set_bank_loan: 0
      },
      Marketing: {
        Price: 0,
        Product: undefined, // Decisions_product object
        Promotion: {
          Quantity: 0,
          Style_quality: 0,
          Style_ecology: 0,
          Style_ethics: 0,
          Style_durability: 0
        }
      },
      Employees: {
        Production_deltas: [], // array of Delta<Employee>
        Marketing_deltas: [],  // array of Delta<Employee>
        Severance_pay: 0
      },
      Production: {
        Production_goal: 0,
        Machines: [],   // array of Delta<Machine>
        Logistics: []   // array of Delta<Warehouse>
      },
      Research: {
        Quality: 0,
        Durability: 0,
        Ecology: 0,
        Promotion: 0,
        Production_cost: 0
      }
    },
    external_factors: {
      Month: 0,
      Inflation: 0,
      Intrest_rate: 0,
      Bridge_loans_intrest_rate: 0,
      Economic_situation_index: 0,
      Tax_rate: 0,
      Turnover: 0,
      Production_minimum_wage: 0,
      Marketing_minimum_wage: 0,
      Machine_on_offer: undefined, // placeholder for a Machine object
      External_storage_price: 0,
      Energy_price: 0,
      Material_price: 0,
      Machine_depreciation_rate: 0
    }
  }

  return {
    socket: websocket,
    state: state,
    gDecisions: function () {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Get_decisions,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: ""
      }
      websocket.send(JSON.stringify(message))
    },

    gCompany: function () {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Get_company,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: ""
      }
      websocket.send(JSON.stringify(message))
    },

    gExternal_factors: function () {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Get_external_factors,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: ""
      }
      websocket.send(JSON.stringify(message))
    },

    /** 
     * @param {simulation.Decisions} decisions
     */
    sDecisions: function (decisions) {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Set_decisions,
        IsResponse: false,
        Error: "",
        DataType: "Decisions",
        Data: decisions
      }
      websocket.send(JSON.stringify(message))
    },

    sReady: function () {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Set_ready,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: ""
      }
      websocket.send(JSON.stringify(message))
    },

    sUnready: function () {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Set_unready,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: ""
      }
      websocket.send(JSON.stringify(message))
    },

    /**
     * @param {simulation.Decisions_product} decisions
     * @param {simulation.Decisions_research} research
     */
    fProduct_stats: function (decisions, research) {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Set_unready,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: {
          Product_decisions: decisions,
          Research_decisions: research,
        }
      }
      websocket.send(JSON.stringify(message))
    },

    /**
     * @param {string} chat - The message broadcased in the chat
     */
    bChat: function (chat) {
      /**
       * @type {Message}
       */
      let message = {
        Method: Methods.Set_unready,
        IsResponse: false,
        Error: "",
        DataType: "",
        Data: chat
      }
      websocket.send(JSON.stringify(message))
    }
  }
}

