/**
 * @callback eventHandler
 * @param {MessageEvent} event
 */

/**
 * Message sent to server
 * @typedef {{Method: string, IsResponse: boolean, Error: string, Data: any}} Message
 */

export const Methods = {
	Get_decisions: 'gDecisions',
	Get_company: 'gCompany',
	Get_external_factors: 'gExternal_factors',
	Get_employees: 'gEmployees',
	Get_unemployed_employees: 'gUnemployedEmployees',

	Set_company: 'sCompany',
	Set_decisions: 'sDecisions',
	Set_ready: 'sReady',
	Set_unready: 'sUnready',
	Get_product_components: 'gProductComponents',

	Broadcast_chat: 'bChat',

	Sim_starting: 'bSim_starting',
	Sim_done: 'bSim_done'
};

/** @type {import('./simulation').clientState}*/
export const baseState = {
	Company: {
		ID: 0,
		Name: '',
		Balance: 0,
		Loans: 0,
		BridgeLoans: 0,
		DecisionHistory: [],
		Reports: [],
		Tech: {
			Quality: 1,
			MaterialUse: 1,
			ProductionCost: 1,
			Ecology: 1,
			Durability: 1
		},
		BaseMarketingStrength: 0,
		Offers: {}, // placeholder for Offer
		Warehouses: [],
		ProductsInStorage: {},
		Machines: []
	},
	Decisions: {
		Predictions: {
			SalesPrediction: 0
		},
		Finances: {
			SetBankLoan: 0
		},
		Products: {}, // Decisions_product object
		Employees: {
			ProductionDeltas: [], // array of Delta<Employee>
			MarketingDeltas: [], // array of Delta<Employee>
			SeverancePay: 0
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
	ExternalFactors: {
		Month: 0,
		Inflation: 0,
		IntrestRate: 0,
		BridgeLoansIntrestRate: 0,
		EconomicSituationIndex: 0,
		TaxRate: 0,
		ProductionMinimumWage: 0,
		MarketingMinimumWage: 0,
		MachineOnOffer: {
			ID: -1,
			ProductionCapacity: 0,
			RequiredWorkers: 1,
			MinimumWorkers: 1,
			AssignedWorkersIds: [],
			EnergyUse: 0,
			Value: 1,
			MaintananceCost: 1,
			AssignedProductID: '-1'
		}, // placeholder for a Machine object
		ExternalStoragePrice: 0,
		EnergyPrice: 0,
		MaterialPrice: 0,
		MachineDepreciationRate: 0
	},
	Employees: {
		production: [],
		marketing: []
	},
	Unemployed: {
		production: [],
		marketing: []
	},
	productComponents: null
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
 * @property {(Type: employeeType)=>void} Connection.gUnemployedEmployees
 * @property {()=>void} Connection.gProductComponents
 *
 * @property {(ID: number)=>void} Connection.sCompany
 * @property {(decisions: import('./simulation').Decisions)=>void} Connection.sDecisions
 * @property {()=>void} Connection.sReady
 * @property {()=>void} Connection.sUnready
 *
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
				Data: { Type }
			};
			websocket.send(JSON.stringify(message));
		},
		/** @param {employeeType} Type */
		gUnemployedEmployees: (Type) => {
			/** @type {Message} */
			let message = {
				Method: Methods.Get_unemployed_employees,
				IsResponse: false,
				Error: '',
				Data: { Type }
			};
			websocket.send(JSON.stringify(message));
		},

		gProductComponents: () => {
			let message = {
				Method: Methods.Get_product_components,
				IsResponse: false,
				Error: '',
				Data: null
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
				Data: ''
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
