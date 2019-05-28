# go-generate-qslcard
generate QSL Card from JSON fomat log.

In amateur radio, QSL cards are exchanged as proof of communication. This command created a CLI command to generate QSL Card from JSON format log.

# Usage
予定

```sh
qslcard -config example-config.json example-logs.json
```

# Example
[QSL Card](example-logs-qslcards.pdf "Example")

## exapmle-logs.json
``` exapmle-logs.json
[
  {
    "id": 1,
    "date": {
      "year": "2018",
      "month": "12",
      "day": "07"
    },
    "time": "12:22",
    "his_call": "JJ1HGP",
    "mode": "SSB",
    "rst": "59",
    "band": "50"
	},
  {
    "id": 2,
    "date": {
      "year": "2018",
      "month": "12",
      "day": "07"
      },
      "time": "12:25",
      "his_call": "JA1BCD",
      "mode": "SSB",
      "rst": "59",
      "band": "21"
   }
]
```

## example-config.json
``` example-config.json
{
  "CardSize": {
    "H": 283.5,
    "W": 419.5
  },
  "StationData": {
     "Call": "JJ1HGP",
     "QRA": "Your Name",
     "Rig": "AA000"
  },
  "UrCallSign": {
    "Size": 42,
    "Name": "migu-1m",
    "Location": "./fonts/migu-1m-regular.ttf"
  },
  "Body": {
    "Size": 12,
    "Name": "migu-1m",
    "Location": "./fonts/migu-1m-regular.ttf"
  }
}
```
