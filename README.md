# Coding Exercise: Simple OAuth2 Server
_(version v1.0.0)_

## Exercise
### Tasks
*    Create a Golang http server that issues JWT Access Tokens (rfc7519) using Client Credentials Grant with Basic Authentication (rfc6749) in endpoint /token
*    Sign the tokens with a self-made RS256 key
*    Provide an endpoint to list the signing keys (rfc7517)
*    Provide deployment manifests to deploy the server in Kubernetes cluster
*    Create an Introspection endpoint (rfc7662) to introspect the issued JWT Access Tokens
### Remarks
* Publish the exercise in a git server and grant us access to review it.
* Avoid a single commit with the whole solution, we want to see how the solution was developed incrementally.
* Provide instructions to execute it

## Setup & Configuration
### Setup
1. Download go greater than v1.20
2. Download all dependencies and run 
```bash 
go mod tidy
```
3. Verify all tests are passed with
```bash 
make test
```
4. Run the Redis
```bash
redis-server
```

### Configuration
1. Environment variables 
* Add your environment name: example `ENV=local`
* Setup all environment variables in [environment](environment) package naming `{env_name}.env` Example:
```
LOG_LEVEL=debug
HOST=localhost
PORT=8080
DURATION=24h
KEY_UUID=18188a8d-7784-462b-bb91-5b3a540e588c
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```
_Note: by default the application runs with local environment variables._

2. Generate local RS256 keys with the endpoint: `/generate-data` (more details see the next section).

## User Guide
## Health Check
#### Endpoint: POST `/health-check`
#### Description:
This endpoint verifies if the server is up.

## Generate local keys
#### Endpoint: POST `/generate-data`
#### Description:
When the endpoint is triggered, it generates two data in Redis:
* `signature_{uuid}`: this is signing key in the RFC7517 format
* `secret_{uuid}`: this is the private key in bytes

The UUID wants use must be same with `KEY_UUID` of env variable.

#### Example from Postman
Request:
```
http://localhost:8080/generate-data
```

Response: `data generated`

Redis result:

![img.png](readme_images/img_2.png)

### Create Token
#### Endpoint: POST `/token`
#### Header: Basic Auth
#### Description:
This endpoint generate token from user and password getting from Authorization header.
#### Example from Postman
Request:
```
http://localhost:8080/token
Basic Auth: Username: user, Password: password
```
![img.png](readme_images/img.png)

Response:
```json
{
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjM5YTUxZDczLWJlMDAtNDIzMy1hODlhLTI3YjlmMDhhMDI1NSIsImJhc2U2NCI6ImRYTmxjanB3WVhOemQyOXlaQT09IiwiaXNzdWVkQXQiOiIyMDIzLTA2LTEzVDE3OjAxOjQ5LjIzNjc1NC0wNTowMCIsImV4cGlyZWRBdCI6IjIwMjMtMDYtMTRUMTc6MDE6NDkuMjM2NzU0LTA1OjAwIn0.JGUKZSlIF7sd65XS0V2z3_UwxH6jdinTHr1TdyCmz1H0dce1twVvYQcw5K4S8zbiFMfBIcx9cpkMcmxHzyVFKMcLk4Pfnd5NH3H_H5RVDG2xNDsMYwDfmUKvTzqDmp88nrNqHvk0NDEFuXIoOmmw_J2aMCXG7pkZD--jCNQA5nQG-WnVvXXL8D_vqzdQrHZogpxqHCp65vgV3cEsBGpxU1uc13o8Foz8NSzt0WQYk7teiK-hsjSoZKPEqCvye1CIhCP-dalB71kFy-FrwvOEpxB4tRl7IYuI_i_Qjt0fyC5hqZF_aD2JmgTpY4Fz9bKf88WUusUd4YEWRq4lHDlvplJDzyO7OcEwdOzbetYB5D8RVZZw2JKr1ET1OwsHEEEP0vglzDNsrjhhNIDz5Tr4WBiNSDQMA7lnmgBLb5P8k15pfdg3wGy47K1tG2RcOfLUdzJxaVnbSRf6Z3gjYjwgiDYAFmbyqn5ZVC4XVFAEiaAXf-Od4zmrg1pByoPefv0s_DDy2Y39WxWbAjmMyWUapi6tK72NJ86xp-pXQ3XeeMEX0X6MozUiiUMLGwpBPROzhb219KD4JjlgbEolaBbc853o_xCpH2KgLyq64nDF-sVrWkJFUETDu_4hILo3pWDJDgNwJremxr8QZYfZhayVHyXeqq7q1fDBgLs0r66LZ5Y",
    "message": "token generated"
}
```

### Verify Token
#### Endpoint: POST `/verify-token`
#### Header: Bearer Token
#### Description:
This endpoint validate the token from Authorization header.
#### Example from Postman
Request:
```
http://localhost:8080/verify-token
Bearer Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjM5YTUxZDczLWJlMDAtNDIzMy1hODlhLTI3YjlmMDhhMDI1NSIsImJhc2U2NCI6ImRYTmxjanB3WVhOemQyOXlaQT09IiwiaXNzdWVkQXQiOiIyMDIzLTA2LTEzVDE3OjAxOjQ5LjIzNjc1NC0wNTowMCIsImV4cGlyZWRBdCI6IjIwMjMtMDYtMTRUMTc6MDE6NDkuMjM2NzU0LTA1OjAwIn0.JGUKZSlIF7sd65XS0V2z3_UwxH6jdinTHr1TdyCmz1H0dce1twVvYQcw5K4S8zbiFMfBIcx9cpkMcmxHzyVFKMcLk4Pfnd5NH3H_H5RVDG2xNDsMYwDfmUKvTzqDmp88nrNqHvk0NDEFuXIoOmmw_J2aMCXG7pkZD--jCNQA5nQG-WnVvXXL8D_vqzdQrHZogpxqHCp65vgV3cEsBGpxU1uc13o8Foz8NSzt0WQYk7teiK-hsjSoZKPEqCvye1CIhCP-dalB71kFy-FrwvOEpxB4tRl7IYuI_i_Qjt0fyC5hqZF_aD2JmgTpY4Fz9bKf88WUusUd4YEWRq4lHDlvplJDzyO7OcEwdOzbetYB5D8RVZZw2JKr1ET1OwsHEEEP0vglzDNsrjhhNIDz5Tr4WBiNSDQMA7lnmgBLb5P8k15pfdg3wGy47K1tG2RcOfLUdzJxaVnbSRf6Z3gjYjwgiDYAFmbyqn5ZVC4XVFAEiaAXf-Od4zmrg1pByoPefv0s_DDy2Y39WxWbAjmMyWUapi6tK72NJ86xp-pXQ3XeeMEX0X6MozUiiUMLGwpBPROzhb219KD4JjlgbEolaBbc853o_xCpH2KgLyq64nDF-sVrWkJFUETDu_4hILo3pWDJDgNwJremxr8QZYfZhayVHyXeqq7q1fDBgLs0r66LZ5Y
```
![img_1.png](readme_images/img_1.png)

Response:
```json
{
  "issuedAt": "2023-06-13T17:01:49.236754-05:00",
  "expiredAt": "2023-06-14T17:01:49.236754-05:00",
  "message": "valid token"
}
```