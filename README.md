# Godis

Godis is a simple, probably naive, in-memory key-value database written in Go,
for learning purposes. It is inspired by Redis, but it is not a Redis clone.

### Protocol

Godis uses a websocket protocol to communicate between a client and the
server. It exchanges JSON messages for commands and replies.

For example, sending a command to the server to set the key "example" to the
string "hello, world" will respond with the reply of type ACK (for acknowledge)
with the value "OK".

```json
{ "o": "SET", [ "example" "hello, world" ] }
{ "t": "ACK", "v": "OK" }
```

### Commands

Godis supports the following commands:

| Syntax    | Description                                       |
| :-------- | :------------------------------------------------ |
| `SET`     | Sets a key to hold a string value                 |
| `GET`     | Returns the string value of a key                 |
| `MGET`    | Returns the string values of one or more key      |
| `DEL`     | Deletes one or more keys                          |
| `EXISTS`  | Determines if one or more keys exists             |
| `INCR`    | Increments the integer value of a key by one      |
| `DECR`    | Decrements the integer value of a key by one      |
| `KEYS`    | Returns all key names that matches a glob pattern |
| `DBSIZE`  | Returns the number of keys in the database        |
| `FLUSHDB` | Removes all the keys in the database              |

