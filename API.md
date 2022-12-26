# Websocket API

## APIv1 Incoming Data
```
type IncomingData struct{
	DataType   string          `json:"datatype"`
	RawData    json.RawMessage `json:"rawdata"`
	APIVersion int             `json:"version"`
}
```


### JSON: datatype
	Type of raw data being sent over JSON field rawdata.
#### Valid JSON Data APIv1
	"minecraft"

### JSON: rawdata
	Raw JSON data, differs according to datatype.
#### Valid JSON Data
	RAW JSON

### APIVersion
	API Version being handled by client connecting to the bot. Data sent over the websocket is handled according to API version negotiated.
#### Valid JSON Data
	Integer


## websocket.datatype: minecraft
	First JSON message to the bot must always be of datatype "authenticate", if first message isn't of datatype "authenticate" then bot closes connection with HTTP Status 401 (Unauthorized). 

### websocket.rawdata:
```
type Data struct {
	DataType string          `json:"datatype"`
	RawData  json.RawMessage `json:"rawdata"`
	AuthKey  string          `json:"authkey"`
}
```
### JSON: datatype
	Type of raw data being sent over JSON field rawdata.
#### Valid JSON Data
	- authenticate
	- playerjoinevent
	- playermessageevent

### JSON: authkey
	Authentication key provided by bot.
#### Valid JSON Data
	String


### websocket.rawdata.datatype: authenticate
	Always sent immediately after establishing connection to websocket. If first message datatype is not authenticate, bot closes connection with HTTP Status 401 (Unauthorized). If authentication is successful, bot responds with HTTP Status 200 (OK) and keeps connection alive. If authentication fails, bot closes connection  with HTTP Status 401 (Unauthorized).

#### websocket.rawdata.rawdata:
```
type Authenticate struct {
	AuthKey        string `json:"authkey"`
	DefaultChannel string `json:"channel"`
	GuildID        string `json:"guild"`
}
```

#### JSON: authkey
	Authentication key provided by bot.
##### Valid JSON Data
	String
