# API Docs

## Push

Trigger a push notification manually. This endpoint is mainly used for debugging.

### **GET** - /push/{{device-udid}}

<details>

#### CURL

```shell
curl -X GET "https://dev.micromdm.io/push/<device-udid>" \
-u "micromdm":"<your-api-token>"
```

#### Security

- Basic Authentication
  - **username**: micromdm
  - **password**: your-api-token

</details>

## MDM Commands

### DeviceInformation

<details>
### **POST** - /v1/commands

#### CURL

```sh
curl -X POST "https://dev.micromdm.io/v1/commands" \
    -H "Content-Type: application/json; charset=utf-8" \
    --data-raw "$body" \
-u "micromdm":"<your-api-token>"
```

#### Header Parameters

- **Content-Type** should respect the following schema:

```shell
{
  "type": "string",
  "enum": [
    "application/json; charset=utf-8"
  ],
  "default": "application/json; charset=utf-8"
}
```

#### Body Parameters

- **body** should respect the following schema:

```shell
{
  "type": "string",
  "default": "{\"request_type\":\"DeviceInformation\",\"udid\":\"564DA875-35DE-E49B-7FCF-6A8FFBE52EF7\"}"
}
```

#### Security

- Basic Authentication
  - **username**: micromdm
  - **password**: your-api-token

</details>

## InstallApplication

<details>
### **POST** - /v1/commands

### CURL

```sh
curl -X POST "https://dev.micromdm.io/v1/commands" \
    -H "Content-Type: application/json; charset=utf-8" \
    --data-raw "$body" \
-u "micromdm":"<your-api-token>"
```

#### Header Parameters

- **Content-Type** should respect the following schema:

```shell
{
  "type": "string",
  "enum": [
    "application/json; charset=utf-8"
  ],
  "default": "application/json; charset=utf-8"
}
```

#### Body Parameters

- **body** should respect the following schema:

```shell
{
  "type": "string",
  "default": "{\"request_type\":\"InstallApplication\",\"manifest_url\":\"https://mdm.yourserver/repo/signed-app.plist\",\"udid\":\"564DA875-35DE-E49B-7FCF-6A8FFBE52EF7\"}"
}
```

#### Security

- Basic Authentication
  - **username**: micromdm
  - **password**: your-api-token

</details>
