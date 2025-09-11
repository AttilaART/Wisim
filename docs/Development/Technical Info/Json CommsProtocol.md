>


Protocol for interpreting data sent along Web-Socket channels in Wisim

# Syntax Message


A message is principally a JSON object with the following fields:

```json
{
	"Method": String,
	"IsResponse": Boolean,
	"Error": String,
	"DataType": String,
	"Data": JSON Object
}
```

Using JSON objects simplifies message sending as tooling around JSON is already widely available as opposed to the [[Outdated/CommsProtocol (OLD)|previous protocol]].

Note: `DataType` should match the reciever implementation's type. So if sending to Godot the type would usually be `Dictionary` whereas when sending to Go the type might be `Decisions`
## Eg. Request
```Json
{"Method": "gExternal_factors", "IsResponse": false, "Data": {}}
```

```Json
{"Method": "sDecisions", "IsResponse": false, "Data": {"Machines": 20 ...}} 
```
## Eg. Response
```Json
{"Method": "gExternal_factors", "IsResponse": true, "Error": null, Data: {IntrestRate: 0.122 ...}}
```

```Json
{"Method": "sDecisions", "IsResponse": true, "Error": "Invalid JSON", Data: null}
```
# Available Methods[^1]:
- `gDecisions() -> Decisons`: Gets latest decisions of company
- `gCompany() -> Company`: Gets latest company state
- `gExternal_factors() -> External_factors`: Gets latest external factors
- `gEmployees(type) -> []Employee` Gets current list of employees
- `gUnemployedEmployees(type) -> []Employee` Gets current list of unemployed employees

- `sCompany(company-number) -> bool`: When joining, asks server to set company
- `sDecisions(Decisions) -> void`: Sets company decisions
- `sReady() -> void` Marks company as ready
- `sUnready() -> void`: Marks company us unready

- `fProduct_stats(Product_decisions, Research_decisions) -> Product`: Calculates product stats
- `bChat(message) -> void`: sends message in game chat

## Server Only:
- `bSim_starting() -> void`: Tells client that simulation is starting
- `bSim_done() -> void`: Tells client that simulation is finished

>[!Note] Note on Method names
> `g` methods get some data from the server
> `s` methods send data to the server
> `f` methods execute a function call
> `b` methods broadcasts to all clients

[^1]: The methods are only shown with function notation for simplicity. To the server and client they are just JSON objects with Strings representing the methods