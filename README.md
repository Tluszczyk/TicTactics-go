![Coverage](https://img.shields.io/badge/Coverage-67.25-yellow)

<div style="text-align:center; margin: 100px;">
    <img src="documentation/resources/logo.png" alt=""/>
</div>

# TicTactics
This application is a web-based implementation of the game Tic-Tactics. The game was originaly created by studio [HiddenVariable](https://www.hiddenvariable.com/), but near 2017 it has been [taken down](https://www.hiddenvariable.com/tictactics/). 

This project is an attempt to recreate the game.

## 📜 Game Rules
*TODO replace diagrams with screenshots from frontend*

The game is a variation of the classic Tic-Tac-Toe game. The game is played on a 3x3 grid of Tiles, where each Tile contains a 3x3 grid of Cells.

**Starting Moves:**

The first player can make their initial move on any Tile on the board.
The second player responds by playing in the corresponding Tile.

<div style="text-align: center;">
    <img src="documentation/resources/Grid_with_first_move.svg" alt=""/>
</div>

**Subsequent Moves:**

Each subsequent move must be played in the corresponding Tile on the board based on the opponent's move.
For example, if X played in the upper left Cell of a Tile, O must play in the upper left Tile of the board. You cannot however force your opponent to play on a Tile that they've just played on.

<div style="text-align: center;">
    <img src="documentation/resources/Grid_with_second_move.svg" alt=""/>
</div>

**Tile Winning Conditions:**

A player wins a Tile when they achieve three of their symbols in a row (horizontally, vertically, or diagonally).

**Game Winning Conditions:**

The game is won by a player when they secure victories in three Tiles in a row (horizontally, vertically, or diagonally).

<div style="text-align: center;">
    <img src="documentation/resources/Grid_O_won.svg" alt=""/>
</div>

**Continuous Play on Won Tiles:**

Players can continue to place moves on Tiles that have already been won.

**Filled Tile Flexibility:**

When a player is compelled to play on a fully filled Tile, they are permitted to make their move anywhere else on the board.

**Game End Conditions:**

The game concludes when either player wins or when all Cells on the board have been played.
If no player has won three Tiles in a row, the game results in a draw.

## 🏛️ Architecture ![](https://progress-bar.dev/100)

A high level architecture diagram of the service is shown below.

<div style="text-align: center;">
    <img src="documentation/resources/HighLevelArchitecture.svg" alt=""/>
</div>

Below is a description of the components of the service.

### 🚀 API 
It will be deployed as a REST API with the following endpoint groups:

| Endpoint | Description               | Implemented                     | Covered by tests                |
|----------|---------------------------|---------------------------------|---------------------------------|
| `/auth`  | Authentication endpoints  |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/user`  | User management endpoints |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game`  | Game management endpoints |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |

#### Authentication Endpoints

| Endpoint         | Method | Description        | Implemented                     | Covered by tests                |              
|------------------|------  |--------------------|---------------------------------|---------------------------------|
| `/auth/session`  | GET    | Validate a session |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
|                  | POST   | Create a session   |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
|                  | DELETE | Remove session     |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |

#### User Management Endpoints

| Endpoint | Method | Description      | Implemented                     | Covered by tests                |
|----------|------  |------------------|---------------------------------|---------------------------------|
| `/user`  | GET    | Get a user       |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
|          | POST   | Create a user    |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
|          | DELETE | Delete a user    |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |

#### Game Management Endpoints

| Endpoint           | Method | Description      | Implemented                     | Covered by tests                |
|--------------------|--------|------------------|---------------------------------|---------------------------------|
| `/game/create`     | POST   | Create a game    |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game/join`       | PUT    | Join a game      |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game/leave`      | PUT    | Leave a game     |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game/leaveAll`   | PUT    | Leave all games  |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game/listGames`  | GET    | List games       |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game`            | GET    | Get a game       |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |
| `/game/move`       | PUT    | Put a move       |![](https://progress-bar.dev/100)|![](https://progress-bar.dev/0)  |

### 💾 Database
#### <img src="https://static.vecteezy.com/system/resources/previews/029/345/981/non_2x/database-icon-data-analytics-icon-monitoring-big-data-analysis-containing-database-free-png.png" alt="" height="20em"/> Schema ![](https://progress-bar.dev/100)
The database will follow this relationship schema:

<div style="text-align: center;">
    <img src="documentation/resources/Database.svg" alt=""/>
</div>

<u>This is not the database's diagram. Each database implementation has it's own database diagram</u>

The *Sessions* table has a Time-to-Live (TTL) set via an environment variable, which means it will be automatically deleted after a certain time. Additionally, sessions are deleted when the user logs out.

#### Implementation Options
1. MongoDB <img src="https://www.mongodb.com/assets/images/global/favicon.ico" alt="Kubernetes" height="18em"/>
1. DynamoDB <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/fd/DynamoDB.png/120px-DynamoDB.png" alt="Kubernetes" height="18em"/>

#### 📦 Data Model
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

## Ideas for the future
* AI for the game
* Authentication cache
* Automated tests
* Board record compression
* Game rules customistaion
* GRPC connections
* Archive for old games

## Testing *In progress* 
![hippo](https://media3.giphy.com/media/aUovxH8Vf9qDu/giphy.gif)

* Unit tests
* Integration tests
* Deploymemnt tests
* Automated tests