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
1. Download all dependencies and run 
```bash 
go mod tidy
```
2. Verify all tests are passed with
```bash 
make test
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
SECRET_KEY_PATH=./jwtRS256.key
```
_Note: by default the application runs with local environment variables._

2. Generate your RS256 keys with the commands: 
```bash
ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key
openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
```

## User Guide
### Create Token
#### Endpoint: POST `/token`
#### Header: Basic Auth
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