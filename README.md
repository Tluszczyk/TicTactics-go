![Coverage](https://img.shields.io/badge/Coverage-67.25-yellow)

<div style="text-align:center; margin: 100px;">
    <img src="documentation/resources/logo.png" alt=""/>
</div>

# TicTactics
This application is a web-based implementation of the game Tic-Tactics. The game was originaly created by studio [HiddenVariable](https://www.hiddenvariable.com/), but near 2017 it has been [taken down](https://www.hiddenvariable.com/tictactics/). 

This project is an attempt to recreate the game.

## 🏛️ Architecture ![](https://progress-bar.dev/100)

A high level architecture diagram of the service is shown below.

<div style="text-align: center;">
    <img src="documentation/resources/HighLevelArchitecture.svg" alt=""/>
</div>

Below is a description of the components of the service.

### 🚀 API 
It will be deployed as a REST API with the following endpoint groups:

| Endpoint | Description               | Progress                        |
|----------|---------------------------|---------------------------------|
| `/auth`  | Authentication endpoints  |![](https://progress-bar.dev/66) |
| `/user`  | User management endpoints |![](https://progress-bar.dev/100)|
| `/game`  | Game management endpoints |![](https://progress-bar.dev/0)  |

#### Authentication Endpoints

| Endpoint         | Method | Description        | Progress                        |              
|------------------|------  |--------------------|---------------------------------|
| `/auth/session`  | GET    | Validate a session |![](https://progress-bar.dev/100)|
|                  | POST   | Create a session   |![](https://progress-bar.dev/100)|
|                  | DELETE | Remove session     |![](https://progress-bar.dev/0)  |

#### User Management Endpoints

| Endpoint | Method | Description      | Progress                        |
|----------|------  |------------------|---------------------------------|
| `/user`  | GET    | Get a user       |![](https://progress-bar.dev/100)|
|          | POST   | Create a user    |![](https://progress-bar.dev/100)|

#### Game Management Endpoints

| Endpoint           | Method | Description      | Progress                      |
|--------------------|--------|------------------|-------------------------------|
| `/game/create`     | POST   | Create a game    |![](https://progress-bar.dev/0)|
| `/game/join`       | PUT    | Join a game      |![](https://progress-bar.dev/0)|
| `/game/leave`      | PUT    | Leave a game     |![](https://progress-bar.dev/0)|
| `/game/list-games` | GET    | List games       |![](https://progress-bar.dev/0)|
| `/game`            | GET    | Get a game       |![](https://progress-bar.dev/0)|
|                    | PUT    | Update a game    |![](https://progress-bar.dev/0)|

### 💾 Database
#### 📜 Schema ![](https://progress-bar.dev/100)
The database will follow this relationship schema:

<div style="text-align: center;">
    <img src="documentation/resources/Database.svg" alt=""/>
</div>

Since one of the implementation options is DynamoDB, before designing the database, it is important to understand the [data model](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html#HowItWorks.CoreComponents.DataModel) of DynamoDB and how to [model relationships](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-modeling-nosql-B.html) between entities. First we have to describe the access patterns that will occur.

#### Access Patterns
1. User Access Patterns
    1. Get a user by username
    1. Get a user by email
    1. Get a user by uid

1. Password Access Patterns
    1. Get a password by uid

1. Game Access Patterns
    1. Get a game by player id

1. Session Access Patterns
    1. Get a session by uid

#### 💾 DynamoDB Data Model ![](https://progress-bar.dev/100)
#### Design
Based on the access patterns described in the code, we can design a DynamoDB table with the following attributes:

* PK: Partition key
* SK: Sort key
* GSI1PK: Global secondary index partition key
* GSI1SK: Global secondary index sort key

Here is a possible design for the DynamoDB table:

#### Table: TicTactics
| Attribute | Type   | Description                                                                 |
|-----------|--------|-----------------------------------------------------------------------------|
| PK        | String | Partition key. Possible values: USER#{username}, EMAIL#{email}, UID#{uid}, PASSWORD#{uid}, GAME#{gameId}, SESSION#{sessionId} |
| SK        | String | Sort key. Possible values: USER, PASSWORD, GAME, SESSION                     |
| GSI1PK    | String | Global secondary index partition key. Possible values: GAME#{playerId}      |
| GSI1SK    | String | Global secondary index sort key. Possible values: GAME                       |

#### Accessing the table
User Access Patterns
1. Get a user by username

    Query the table with ```PK = USER#{username} ```and``` SK = USER.```

1. Get a user by email

    Query the table with ```PK = EMAIL#{email} ```and``` SK = USER.```

1. Get a user by uid

    Query the table with ```PK = UID#{uid} ```and``` SK = USER.```

Password Access Patterns
1. Get a password by uid

    Query the table with ```PK = PASSWORD#{uid} ```and``` SK = PASSWORD.```

Game Access Patterns
1. Get a game by player id

    Query the table with ```GSI1PK = GAME#{playerId} ```and``` GSI1SK = GAME.```

Session Access Patterns
1. Get a session by uid

    Query the table with ```PK = SESSION#{sessionId} ```and``` SK = SESSION.```
## Possible Deployment Options
### 🐳 Docker  ![](https://progress-bar.dev/100)
One docker image contains the API and the other contains the database. The API will be deployed as a REST API and the database will be deployed as a MongoDB instance.

<div style="text-align: center;">
    <img src="documentation/resources/DockerDeploymentOption.svg" alt=""/>
</div>

### 🌐 AWS ![](https://progress-bar.dev/0)
The API will be deployed as a REST API, invoking kambda functions and the database will be deployed as a DynamoDB instance.

<div style="text-align: center;">
    <img src="documentation/resources/LambdaDeploymentOption.svg" alt=""/>
</div>

### <img src="https://raw.githubusercontent.com/kubernetes/kubernetes/dbd3f3564ac6cca9a152a3244ab96257e5a4f00c/logo/logo.svg" alt="Kubernetes" height="18em"/> Kubernetes ![](https://progress-bar.dev/0)
The API will be deployed as a REST API, services will be deployed as a pod in a running kubernetes and the database will be deployed as a MongoDB instance.

<div style="text-align: center;">
    <img src="documentation/resources/K8sDeploymentOption.svg" alt=""/>
</div>