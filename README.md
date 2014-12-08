# Interoper

**Interoper** is a simple interop testing tool for Telehash. It uses docker under the hood and makes json reports.


## Test file format

Test files are Markdown files which start with an (indent formatted) code block of some JSON followed by a description of the test.

The JSON header can have the following entries:

* `"timeout"` defines the maximum runtime for a test. (examples: `"3ms"`, `"3s"`, `"3m"`, `"3h"`)
* `"containers"` defines container settings for the individual roles. It must be an object. Each entry can have the `"command"` entry.

```markdown
    {
      "timeout": "3m"
    }

# `net-link` Test basic link establishment.

## Scenario

The SUT must start an endpoint and write its keys and paths to `/shared/id_sut.json`.

The Driver must start an endpoint, load the keys and paths from `/shared/id_sut.json` and establish a link with the SUT. The driver must close the link after 2.5 minutes.

## Failure conditions

* The link fails to open (after 1 minute)
* The link breaks before it is closed.

## Success conditions

* The remains open for at least 2.5 minutes (It keeps the exchange open)
* The link closes cleanly
```

## Creating docker images

Interoper can work with any docker image as long as it has a `th-test` executable in its `PATH`.

The `th-test` executable should accept two arguments. The first is the name of the test that should be run and the second argument is the role the process should assume (ex. `th-test net-link driver` asks the `th-test` executable to run the `net-link` test as the `driver`).

`interoper test` will build the `Dockerfile` in the current directory and run all implementations against this image.

## Test Control Protocol.

Test processes can use regular logging to save informational data. In addition to the regular logging processes can also send (and receive) JSON formatted events over STDOUT.

All process types must send a `{"ty":"ready"}` when they are ready with there setup and can start running the test. The Driver will then run its test and when it is done it sends the `{"ty":"done"}` command. When The Driver is done running the tests it must exit. When the driver exits all other processes will be killed.

### Event JSON structure

All events must have the following format.
```json5
{
  "id": 0,
  // integer; unique event id. (must be unique withing the process)

  "ty": "endpoint.new",
  // string;  type of the event.

  "in": {}
  // object;  info object. contains additional information on the event.
}
```

### Event Types

#### `"ready"`

Signals that the sending process is ready to start the test. No extra information must be send with this event.

#### `"done"`

Signals that the sending process is done running the test. No extra information must be send with this event.

#### `"exited"`

Signals that the sending process has exited. This event is emitted by interoper and should not be emitted by the processes.

| param | type | description |
| ----- | ---- | ----------- |
| `"exit_code"` | integer | The exit code of the process |

#### `"log"`

All lines printed by the process will be transformed into `"log"` events. This event is emitted by interoper and should not be emitted by the processes.

| param | type | description |
| ----- | ---- | ----------- |
| `"line"` | string | The logged line |

#### `"exec"`

Emitted by interoper at the start of each test.

| param | type | description |
| ----- | ---- | ----------- |
| `"sut"` | string | The name of the tested implementation |
| `"driver"` | string | The name of the driver implementation  |

#### `"status"`

Emitted by interoper at the end of each test.

| param | type | description |
| ----- | ---- | ----------- |
| `"success"` | bool | `true` when the test was run successfully |

#### `"endpoint.new"`

Should be emitted every time a new endpoint is created.

| param | type | description |
| ----- | ---- | ----------- |
| `"endpoint_id"` | integer | unique id assigned to this endpoint |
| `"hashname"` | string | the hashname associated with this endpoint |

#### `"endpoint.started"`

Should be emitted every time an endpoint is started.

| param | type | description |
| ----- | ---- | ----------- |
| `"endpoint_id"` | integer | id of the endpoint |

#### `"endpoint.error"`

Should be emitted every time an endpoint encounters an error.

| param | type | description |
| ----- | ---- | ----------- |
| `"endpoint_id"` | integer | id of the endpoint |
| `"error"` | string | description of the error |

#### `"endpoint.rcv.packet"`

Should be emitted every time an endpoint received a packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"endpoint_id"` | integer | id of the endpoint |
| `"packet_id"` | integer | unique id for this packet |
| `"packet"."src"` | string | the source address the packet was received from |
| `"packet"."msg"` | string | the content of the packet |

#### `"endpoint.drop.packet"`

Should be emitted every time an endpoint dropped a packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"endpoint_id"` | integer | id of the endpoint |
| `"packet_id"` | integer | unique id for this packet |
| `"reason"` | string | the reason for dropping the packet |
| `"packet"."src"` | string | the source address the packet was received from |
| `"packet"."msg"` | string | the content of the packet |

#### `"exchange.new"`

Should be emitted every time a new exchange is created.

| param | type | description |
| ----- | ---- | ----------- |
| `"endpoint_id"` | integer | id of the endpoint |
| `"exchange_id"` | integer | unique id for this exchange |

#### `"exchange.started"`

Should be emitted every time an exchange is opened.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the exchange |
| `"peer"` | string | hashname of the remote endpoint |

#### `"exchange.stopped"`

Should be emitted every time an exchange is closed.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the exchange |

#### `"exchange.rcv.handshake"`

Should be emitted every time an exchange received a handshake.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the endpoint |
| `"packet_id"` | integer | id of the packet |
| `"handshake"."csid"` | integer | the CSID of the key in the body of the handshake |
| `"handshake"."public_key"` | string | the key in the body |
| `"handshake"."at"` | integer | the `"at"` header of the handshake |
| `"handshake"."parts"` | object | the parts in the handshake header |

#### `"exchange.drop.handshake"`

Should be emitted every time an exchange dropped a handshake.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the endpoint |
| `"packet_id"` | integer | id of the packet |
| `"reason"` | string | the reason for dropping the handshake |
| `"handshake"."csid"` | integer | the CSID of the key in the body of the handshake |
| `"handshake"."public_key"` | string | the key in the body |
| `"handshake"."at"` | integer | the `"at"` header of the handshake |
| `"handshake"."parts"` | object | the parts in the handshake header |

#### `"exchange.rcv.packet"`

Should be emitted every time an exchange received a packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the exchange |
| `"packet_id"` | integer | unique id for this packet |
| `"packet"."header"` | object | the header of the packet |
| `"packet"."body"` | string | the body of the packet |

#### `"exchange.drop.packet"`

Should be emitted every time an exchange dropped a packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the exchange |
| `"packet_id"` | integer | unique id for this packet |
| `"reason"` | string | the reason for dropping the packet |
| `"packet"."header"` | object | the header of the packet |
| `"packet"."body"` | string | the body of the packet |

#### `"channel.new"`

Should be emitted every time a new channel is created.

| param | type | description |
| ----- | ---- | ----------- |
| `"exchange_id"` | integer | id of the exchange |
| `"channel_id"` | integer | unique id for this channel (not the `"c"` header) |
| `"channel"."type"` | string | the type of the channel |
| `"channel"."reliable"` | bool | `true` when this is a reliable channel |
| `"channel"."cid"` | integer | the channel id as defined by telehash (the `"c"` header) |

#### `"channel.write"`

Should be emitted every time a channel writes packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"channel_id"` | integer | id of the channel |
| `"packet_id"` | integer | unique id for this packet |
| `"path"` | string | the target path (if specified) |
| `"packet"."header"` | object | the header of the packet |
| `"packet"."body"` | string | the body of the packet |

#### `"channel.write.error"`

Should be emitted every time a channel write failed.

| param | type | description |
| ----- | ---- | ----------- |
| `"channel_id"` | integer | id of the channel |
| `"packet_id"` | integer | unique id for this packet |
| `"reason"` | string | the reason the write failed |
| `"path"` | string | the target path (if specified) |
| `"packet"."header"` | object | the header of the packet |
| `"packet"."body"` | string | the body of the packet |


#### `"channel.rcv.packet"`

Should be emitted every time a channel received a packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"channel_id"` | integer | id of the channel |
| `"packet_id"` | integer | unique id for this packet |
| `"packet"."header"` | object | the header of the packet |
| `"packet"."body"` | string | the body of the packet |

#### `"channel.drop.packet"`

Should be emitted every time a channel dropped a packet.

| param | type | description |
| ----- | ---- | ----------- |
| `"channel_id"` | integer | id of the channel |
| `"packet_id"` | integer | unique id for this packet |
| `"reason"` | string | the reason for dropping the packet |
| `"packet"."header"` | object | the header of the packet |
| `"packet"."body"` | string | the body of the packet |
