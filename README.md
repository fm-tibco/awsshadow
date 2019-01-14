
# AWS Device Shadow
This activity allows you to update/get/delete a device shadow on AWS.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/fm-tibco/awsshadow
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "thingName",
      "type": "string",
      "required": true
    },
    {
      "name": "op",
      "type": "string",
      "required": true,
      "allowed" : ["GET", "UPDATE", "DELETE"]
    }
    
  ],
  "output":[
    {
      "name": "result",
      "type": "object"
    }
  ]
}
```

## Input
| Name     | Required | Description |
|:------------|:---------|:------------|
| thingName   | true     | The name of the "thing" in Aws IoT |         
| op          | true     | The Aws Iot shadow operation to perform  |

## Configuration

### Flow Configuration
Configure a task in flow to update the device shadow of 'raspberry-pi' with a reported temp of "50".

```json
{
  "id": "shadow_update",
  "name": "Update AWS Device Shadow",
  "activity": {
    "ref": "github.com/fm-tibco/awsshadow",
    "input": {
      "thingName": "raspberry-pi",
      "op": "UPDATE",
      "reported": { "temp":"50" }
    }
  }
}
```

To Configure AWS credentials see:
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html