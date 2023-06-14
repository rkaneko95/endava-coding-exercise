# Coding Exercise: Simple OAuth2 Server
_(version v1.0.0)_

## Content
* [Exercise](#exercise)
  * [Tasks](#tasks)
  * [Remarks](#remarks)
* [Get Started](#get-started)
  * [Tools](#tools)
  * [Configuration](#configuration)
  * [Deployment](#deployment)
* [User Guide](#user-guide)
  * [Health Check](#health-check)
  * [Generate local keys](#generate-local-keys)
  * [Create Token](#create-token)
  * [Verify Token](#verify-token)
  * [List Signing Keys](#list-signing-keys)
* [Packages](#packages)
* [Redis](#redis)
  * [Private Key](#private-key)
  * [Public Key](#public-key)
* [Quick Demo](#quick-demo)

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

## Get Started
### Tools
* Docker
* Kubernetes
* IDE (Golang, VSCode, etc)

### Configuration
1. Environment variables
* Add your environment name in Dockerfile: example `ENV=local`
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

2. Verify all tests are passed with
```bash 
make test
```

### Deployment
1. Pull Redis
```bash
docker pull redis
```
2. Deploy the application
```bash
docker build -t endava-coding-exercise .

kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```
3. Generate port forwarding
```bash
kubectl get pods
kubectl port-forward {pod_name} 8080:8080
```
4. Verify the server calling `http://localhost:8080/health-check`
5. Generate local RS256 keys with the endpoint: `/generate-keys` (more details see the next section).

## User Guide
## Health Check
#### Endpoint: POST `/health-check`
#### Description:
This endpoint verifies if the server is up.

## Generate local keys
#### Endpoint: POST `/generate-keys`
#### Description:
When the endpoint is triggered, it generates two data in Redis:
* `signature_{uuid}`: this is signing key in the RFC7517 format
* `secret_{uuid}`: this is the private key in bytes

The UUID wants use must be same with `KEY_UUID` of env variable.

#### Example from Postman
Request:
```
http://localhost:8080/generate-keys
```

Response:
```json
{
    "UUID": "18188a8d-7784-462b-bb91-5b3a540e588c",
    "Message": "Keys generated"
}
```

Redis result:

![img.png](images/img_2.png)

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
![img.png](images/img.png)

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
![img_1.png](images/img_1.png)

Response:
```json
{
  "issuedAt": "2023-06-13T17:01:49.236754-05:00",
  "expiredAt": "2023-06-14T17:01:49.236754-05:00",
  "active": true
}
```

### List Signing Keys
#### Endpoint: GET `/signing-keys`
#### Description:
This endpoint list all signing keys in RFC7662 format. Important that keys were already generated in Redis.

_Note: The JSON format is parsing from `github.com/golang-jwt/jwt/v4` package._
#### Example from Postman
Request:
```
http://localhost:8080/signing-keys
```
Response:
```json
[
    {
        "kty": "RSA",
        "n": "vWCPW9T7nqfngfAVxBqBvMLOpmYcnV841-0ZEmNTze_-0fm9zhulQ6BtDvpsVoItf0KUKvfGF5ba_8vJLsuLDuuYhjrJv3JbcanM62eQijA4wKRuu-an0-VAaGLCqXQZZAYkoq003ykRc0RIfBfmj4GP_a9zFf44XcwtoZ4zRaGymMf4J2Vl5SPfgs3ChZxpjjiXu9mCJgiouKRfGIQ0cnqjIzY1F1okRYE47mrLc3N-hqu40r_IBrQViPSZsSlA-n3Gyg-fiYG6PWVE7Cc-UJySfvcP-iM1-5tHFIx_zYxv_NRCfTydh0xef5pDYOY1mgpuB2lQ8RC54bYqrfiFPQ",
        "e": "AQAB"
    },
    {
        "kty": "RSA",
        "n": "v0vTxIFPtDq1gFBatvAmjEKBQs65WNUSuZeTODFersSnMVsdqK5E-tHNYJqrEWYVIKIbUbeMCL2-t7hiUK5BqJcnPXNvgtXZN-xPJ-ua8zYa7Zmv_rcbSbSWNj5dsbd3JfPUXUsu_gtztKLtdwmL-sB2hlPXJXHBihRPrF2Y0CKarH6T2b4pVde-E9Nv08uEq8tMoGbgSGpXBg41-H47gxeKg2lyM-9sP9JMOoW_lZZFYoSyhi5Pgjckskw0CSID4l0GueDptPbgQCIBO2zEzHgPOooLomDzJsQh9amvgU8VGeyLj1-7CnTHtqbnOYfbbKLSA7puLmMVG-6Az3q4NQ",
        "e": "AQAB"
    },
    {
        "kty": "RSA",
        "n": "u1S184t15CXqL5mvF_1kRlv4pbNhb2N-vhod4rR8OlABBGoaJdjrxcPrmuhaa8sS6oueodJk5QfKBq3vOr1nTU8IXX5Eku07RZvfe_JkBVJFWrHWf1OBeUfl8bjNm_xKrcnuH57V37JOn65smcYJsQJJQFtKN2UQOeaRfnTmMHBWrNYk2IIsNn1JqYTEkMv5Z-ChCRSX33vapkynFv2Sh9pJCwcC3EhMk-vHOUIeahn-pWPrfuWY_IYm7TAmhfMPzSs55tXCbFur6KDhibVRILbz-zSugf83D65GmUnHlfS7E-yuxkT6v3h390fBEPRCv1hTsN6ramTCeJFNUsHCJw",
        "e": "AQAB"
    },
]
```

## Packages
```
endava-coding-exercise/
├── api/
│   ├── handler.go
│   ├── keys.go
│   ├── route.go
│   ├── service.go
│   ├── token_claims.go
│   └── utils.go
├── config/
│   ├── env_variables.go
│   └── logs.go
├── db/
│   └── redis.go
├── environment/
│   └── handler.go
├── images/
├── resources/
├── main.go
├── go.mod
│   └── go.sum
├── deployment.yaml
├── service.yaml
├── Dockerfile
├── Makefile
├── .gitignore
└── README.md
```
* **api:** packages handle all functionality of API such as handlers and logics (service)
* **config:** load all env variables saved in [environment](environment) when the app starts.
* **db:** packages handle all functionality of Redis.
* **environment:** contains all env variables ver environment.
* **images:** contains images use in README.
* **resources:** mock resources use for unit tests.

![workflow.png](images/workflow.png)

## Redis
The Redis saves private and public keys.

## Private Key
It storages as bytes with name `secret_{uuid}`

Example:

![img.png](images/img_3.png)

In the function `getPrivateKey`, it read the bytes and encode to PEM as `RSA PRIVATE KEY`
```go
privateKeyPEM := pem.EncodeToMemory(&pem.Block{
	Type:  "RSA PRIVATE KEY", 
	Bytes: data,
})
```
Then, it parses RSA to get the private key
```go
privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
```

## Public Key
It storages as JSON in format string with name `signature_{uuid}`


The JSON format is parsing from `github.com/golang-jwt/jwt/v4` package.

Example:

![img.png](images/img_4.png)

In the function `getPublicKey`, it read the JSON and get the key using x509 package.
```go
pubKeyBytes, err := x509.MarshalPKIXPublicKey(jwk.Key)
```
Then encode to PEM as `PUBLIC KEY`
```go
publicKeyPEM := pem.EncodeToMemory(&pem.Block{
	Type:  "PUBLIC KEY", 
	Bytes: pubKeyBytes,
})
```
Then, it parses RSA to get the public key
```go
publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
```

## Quick Demo
![code-demo.gif](images/code-demo.gif)