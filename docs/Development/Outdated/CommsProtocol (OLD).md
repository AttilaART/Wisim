> [!warning] Warning
> Information contained in this document is outdated and does not match current information.
> For up to date information, see [[../Technical Info/Json CommsProtocol]]



Protocol for interpreting data sent along WebSocket channels in Wisim

# Syntax Request

```
    Method Data ... (encoded in JSON)
Eg. gExFa "None"
Eg. sDecs {"Machines": 20 ...}
```

All messages start with a **5** letter method identifier^[1] followed by Data formatted as JSON.

[1]: Fixed length method identifiers simplify method parsing, as you can check if the first five characters match a known method

# Syntax Response
```
       Method Status Data (encoded in JSON) or Error
Eg. RE gExFa OK {IntrestRate: 0.122 ...}
Eg. RE sDecs ER "Invalid JSON"```
```

Responses are prefixed `RE` before the method that was called. Similarly the data is prefixed with `OK` or `ER`, with `OK` meaning the request was completed successfully and `ER` that an error occurred.

If message is too short, no response will be sent
# Available Methods:
- `gDecs {}`: Gets latest decisions of company
- `gComp {}`: Gets latest company state
- `gExFa {}`: Gets latest external factors

- `sDecs {Decisions}`: Sets company decisions
- `sRedy {}`: Marks company as ready
- `sURdy {}`: Marks company us unready

- `fProd {Product}`: Calculates product stats
- `bChat {"Message": message}`: sends message in game chat

##Server Only:
- `bSimS {}`: Tells client that simulation is starting
- `bSimD {"Success": boolean, "Message": string}`: Tells client that simulation is finished

>[!Note] Note on Method names
> `g` methods get some data from the server
> `s` methods send data to the server
> `f` methods execute a function call
> `b` methods broadcasts to all clients
