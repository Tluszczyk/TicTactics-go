<div style="text-align:center; margin: 100px;">
    <img src="documentation/resources/logo.png" alt=""/>
</div>

# TicTactics
This application is a web-based implementation of the game Tic-Tactics. The game was originaly created by studio [HiddenVariable](https://www.hiddenvariable.com/), but near 2017 it has been [taken down](https://www.hiddenvariable.com/tictactics/). 

This project is an attempt to recreate the game.

## 🏛️ Architecture

A high level architecture diagram of the service is shown below.

<div style="text-align: center;">
    <img src="documentation/resources/HighLevelArchitecture.svg" alt=""/>
</div>

Below is a description of the components of the service.

### API
It will be deployed as a REST API with the following endpoint groups:

| Endpoint | Description               |
|----------|---------------------------|
| `/auth`  | Authentication endpoints  |
| `/user`  | User management endpoints |
| `/game`  | Game management endpoints |

#### Authentication Endpoints

| Endpoint         | Method | Description        |
|------------------|------  |--------------------|
| `/auth/session`  | GET    | Validate a session |
|                  | POST   | Create a session   |
|                  | DELETE | Remove session     |

#### User Management Endpoints

| Endpoint | Method | Description      |
|----------|------  |------------------|
| `/user`  | GET    | Get a user       |
|          | POST   | Create a session |

#### Game Management Endpoints

| Endpoint           | Method | Description      |
|--------------------|--------|------------------|
| `/game/create`     | POST   | Create a game    |
| `/game/join`       | PUT    | Join a game      |
| `/game/leave`      | PUT    | Leave a game     |
| `/game/list-games` | GET    | List games       |
| `/game`            | GET    | Get a game       |
|                    | PUT    | Update a game    |

### Database
The database will follow this relationship schema:

<div style="text-align: center;">
    <img src="documentation/resources/Database.svg" alt=""/>
</div>

Since one of the implementation options is DynamoDB, before designing the database, it is important to understand the [data model](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html#HowItWorks.CoreComponents.DataModel) of DynamoDB and how to [model relationships](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-modeling-nosql-B.html) between entities. First we have to describe the access patterns that will occur.

#### User Access Patterns
TODO

#### Password Access Patterns
TODO

#### Game Access Patterns
TODO
