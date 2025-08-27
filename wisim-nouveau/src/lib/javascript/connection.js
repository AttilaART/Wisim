/**
 * @callback eventHandler
 * @param {MessageEvent} event
 */

/**
 * Message sent to server
 * @typedef {{Method: string, IsResponse: boolean, Error: string, DataType: DataType, Data: any}} Message
 */

/**
 * @typedef {string} DataType
 */

export const Methods = {
	Get_decisions: 'gDecisions',
	Get_company: 'gCompany',
	Get_external_factors: 'gExternal_factors',
	Get_employees: 'gEmployees',

	Set_company: 'sCompany',
	Set_decisions: 'sDecisions',
	Set_ready: 'sReady',
	Set_unready: 'sUnready',

	Func_calculate_product_stats: 'fProduct_stats',
	Broadcast_chat: 'bChat',

	Sim_starting: 'bSim_starting',
	Sim_done: 'bSim_done'
};

/** @type {import('./simulation').clientState}*/
export const baseState = {
	company: {
		Id: 0,
		Name: '',
		Balance: 0,
		Loans: 0,
		Bridge_loans: 0,
		Decision_history: [],
		Reports: [],
		Global_quality_factor: 0,
		Base_marketing_strength: 0,
		Offer: {
			Price: 0,
			Promotion_quality: 0,
			Promotion_goal: {
				Style_quality: 0,
				Style_durability: 0,
				Style_ecology: 0,
				Style_ethics: 0,
				Quantity: 0
			},
			Product: {
				Id: 0,
				Name: '',
				Weight: 0,
				Material_use: 0,
				Production_cost: 0,

				Base_material_use: 0,
				Base_production_cost: 0,
				Base_quality: 0,
				Base_ecology: 0,
				Base_durability: 0,

				Ethics_factor: 0,
				Quality_factor: 0,
				Ecology_factor: 0,
				Durabilty: 0
			}
		}, // placeholder for Offer
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
			Product: {
				Materials: {
					Quality: 0,
					Ecology: 0,
					Ethical_sourcing: 0
				},
				Manufacturing: {
					Quality: 0,
					Ecological_energy: 0,
					Material_efficiency: 0,
					Durability: 0,
					Max_durability: 0
				}
			}, // Decisions_product object
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
			Marketing_deltas: [], // array of Delta<Employee>
			Severance_pay: 0
		},
		Production: {
			Production_goal: 0,
			Machines: [], // array of Delta<Machine>
			Logistics: [] // array of Delta<Warehouse>
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
	},
	employees: {
		production: []
	}
};

/** @typedef {string} employeeType */
/**
 * @typedef {Object} Connection
 * @property {WebSocket} Connection.socket
 * @property {import('./simulation').clientState} Connection.state
 * @property {()=>void} Connection.gDecisions
 * @property {()=>void} Connection.gCompany
 * @property {()=>void} Connection.gExternal_factors
 * @property {(Type: employeeType)=>void} Connection.gEmployees
 *
 * @property {(ID: number)=>void} Connection.sCompany
 * @property {(decisions: import('./simulation').Decisions)=>void} Connection.sDecisions
 * @property {()=>void} Connection.sReady
 * @property {()=>void} Connection.sUnready
 *
 * @property {(decisions: import('./simulation').Decisions_product, research: import('./simulation').Decisions_research)=>void} Connection.fProduct_stats
 * @property {(chat: string)=>void} Connection.bChat
 *
 */

/**
 * @param {string} url - URL of simulation server, should be localhost:8000
 * @param {eventHandler} eventHandler - Function that handles every event
 * @param {(event: CloseEvent)=>void} onClose - Do something when websocket closes
 * @param {(this: WebSocket, event: Event)=>void} onError - Do something when websocket errors
 * @returns {Promise<{connection: Connection, clientState: import('./simulation').clientState}>}
 */
export async function newConnection(url, eventHandler, onClose, onError) {
	let websocket = new WebSocket(url);

	websocket.addEventListener('message', (event) => {
		eventHandler(event);
	});
	websocket.addEventListener('close', onClose);
	websocket.addEventListener('error', onError);

	/**
	 * @type {import('./simulation').clientState}
	 */

	let state = JSON.parse(JSON.stringify(baseState));

	const methods = {
		socket: websocket,
		state: state,
		gDecisions: function () {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Get_decisions,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: ''
			};
			websocket.send(JSON.stringify(message));
		},

		gCompany: function () {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Get_company,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: ''
			};
			websocket.send(JSON.stringify(message));
		},

		gExternal_factors: function () {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Get_external_factors,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: ''
			};
			websocket.send(JSON.stringify(message));
		},
		/** @param {employeeType} Type */
		gEmployees: (Type) => {
			/** @type {Message} */
			let message = {
				Method: Methods.Get_employees,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: { Type }
			};
			websocket.send(JSON.stringify(message));
		},
		/**
		 * @param {number} ID
		 */
		sCompany: (ID) => {
			let message = {
				Method: Methods.Set_company,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: { ID }
			};

			websocket.send(JSON.stringify(message));
		},

		/**
		 * @param {import('./simulation').Decisions} decisions
		 */
		sDecisions: function (decisions) {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Set_decisions,
				IsResponse: false,
				Error: '',
				DataType: 'Decisions',
				Data: decisions
			};
			websocket.send(JSON.stringify(message));
		},

		sReady: function () {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Set_ready,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: ''
			};
			websocket.send(JSON.stringify(message));
		},

		sUnready: function () {
			/** @type {Message} */
			let message = {
				Method: Methods.Set_unready,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: ''
			};
			websocket.send(JSON.stringify(message));
		},

		/**
		 * @param {import('./simulation').Decisions_product} decisions
		 * @param {import('./simulation').Decisions_research} research
		 */
		fProduct_stats: function (decisions, research) {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Func_calculate_product_stats,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: {
					Product: decisions,
					Research: research
				}
			};
			websocket.send(JSON.stringify(message));
		},

		/**
		 * @param {string} chat - The message broadcased in the chat
		 */
		bChat: function (chat) {
			/**
			 * @type {Message}
			 */
			let message = {
				Method: Methods.Broadcast_chat,
				IsResponse: false,
				Error: '',
				DataType: '',
				Data: chat
			};
			websocket.send(JSON.stringify(message));
		}
	};

	/**
	 * @param {number} timeout
	 */

	const asyncWait = async (timeout) => {
		return new Promise((res) => setTimeout(res, timeout));
	};

	let waitTime = 0;
	const maxWait = 10000;
	while (waitTime <= maxWait) {
		if (websocket.readyState == websocket.OPEN) {
			return { connection: methods, clientState: state };
		}

		await asyncWait(100);
		waitTime += 100;
	}

	throw 'connection timed out';
}

/**
 * @type {eventHandler}
 */
export function handleConnection(event) {
	console.log(event.data);
}
