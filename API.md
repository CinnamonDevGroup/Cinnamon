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
	- minecraft

### JSON: rawdata
Raw JSON data, differs according to datatype.
#### Valid JSON Data
RAW JSON

### APIVersion
API Version being handled by client connecting to the bot. Data sent over the websocket is handled according to API version negotiated.
#### Valid JSON Data
1


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

### JSON: rawdata
Raw JSON data, differs according to datatype.
#### Valid JSON Data
RAW JSON

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

#### JSON: channel
Default channel (ID) where to send messages sent to when receiving events that require the bot to send a message to the specified channel.
##### Valid JSON Data
String

#### JSON: GuildID
Default guild (ID) that the bot is linked to.
##### Valid JSON Data
String


### websocket.rawdata.datatype: playermessageevent
Sent whenever a player sends a public message to the Minecraft chat.

#### websocket.rawdata.data: 
```
type ChatMessage struct {
	UUID  string `json:"uuid"`
	Message string `json:"message"`
	Mention string `json:"mention"`
	Channel string `json:"channel"`
}
```

#### JSON: uuid
Minecraft UUID of player sending the message.
##### Valid JSON Data
String

#### JSON: message
Message content of chat message sent by player.
##### Valid JSON Data
String

#### JSON: mention
(Optional) Message ID sent on discord to reply to.
##### Valid JSON Data
String

#### JSON: channel
(Optional) Channel ID where the Message replying to was sent, only filled if replying to message.
##### Valid JSON Data
String


### websocket.rawdata.datatype: playerjoinevent

#### websocket.rawdata.data
```
type PlayerJoin struct {
	UUID     string `json:"uuid"`
}
```

#### JSON uuid:
Minecraft UUID of player that joined Minecraft server.
##### Valid JSON Data
String