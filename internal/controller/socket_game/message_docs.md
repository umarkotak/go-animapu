# Message Type List

## Base structure
{
	"message_type": "base",
	"meta": {},
	"data": {},
	"headers": {},
  "direction": "request/response"
}

## Request join stream
{
  "message_type": "init_stream",
	"meta": {},
	"data": {},
	"headers": {
    "authorization": ""
  },
  "direction": "request"
}

## Response join stream
{
  "message_type": "init_stream",
	"meta": {},
	"data": {
    "username": "jhone.doe",
    "welcome_message": "welcome to socket game"
  },
	"headers": {},
  "direction": "response"
}

## Request player
{
	"message_type": "player_info",
	"meta": {},
	"data": {},
	"headers": {
    "authorization": ""
  },
  "direction": "request"
}

## Response player
{
  "message_type": "player_info",
	"meta": {},
	"data": {
    "username": "jhone.doe",
    "loc_x": 0,
    "loc_y": 0
  },
	"headers": {},
  "direction": "response"
}

## Request world map
{
	"message_type": "world_map_info",
	"meta": {},
	"data": {},
	"headers": {},
  "direction": "request"
}

## Response world map
{
	"message_type": "world_map_info",
	"meta": {},
	"data": {
    "horizontal_length": 10,
    "vertical_length": 10,
    "maps": [
      {
        "row_idx": 0,
        "row_data": [
          {
            "pos_x": 0,
            "pos_y": 0,
            "info": {
              "terrain": "grassland"
            },
            "players": [
              {
                "username": "jhone.doe"
              }
            ]
          },
          {
            "pos_x": 1,
            "pos_y": 0,
            "info": {
              "terrain": "grassland"
            },
            "players": []
          }
        ]
      }
    ]
  },
	"headers": {},
  "direction": "response"
}

## Request broadcast
{
  "message_type": "global_message",
	"meta": {},
	"data": {
    "message": "hello world!"
  },
	"headers": {},
  "direction": "request"
}

## Response broadcast
{
  "message_type": "global_message",
	"meta": {},
	"data": {
    "from": "jhone.doe",
    "message": "hello world!",
    "ts": "2020-01-01T14:00:00.000Z"
  },
	"headers": {},
  "direction": "response"
}