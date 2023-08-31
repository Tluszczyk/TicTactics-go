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

### IO
This component is responsible for handling all input and output to and from the service. Depending on the type of request, it will forward the request to the appropriate component.

It will be deployed as a REST API with the following endpoint groups:

<div style="text-align: center;">

| Endpoint          | Description               |
| -                 | -                         |
| `/auth`           | Authentication endpoints  |
| `/user`           | User management endpoints |
| `/game`           | Game management endpoints |

</div>

#### Authentication Endpoints

<div style="text-align: center;">

| Endpoint          | Description               |
| -                 | -                         |
| `/auth/register`  | Register a new user       |
| `/auth/login`     | Login a user              |
| `/auth/logout`    | Logout a user             |
| `/auth/validate`  | Validate a user's session |

</div>

#### User Management Endpoints

<div style="text-align: center;">

| Endpoint          | Description               |
| -                 | -                         |
| `/user/profile`   | Get a user's profile      |

</div>

#### Game Management Endpoints

<div style="text-align: center;">

| Endpoint          | Description               |
| -                 | -                         |
| `/game/create`    | Create a new game         |
| `/game/join`      | Join an existing game     |
| `/game/leave`     | Leave a game              |
| `/game/move`      | Make a move in a game     |

</div>

### Auth
### User Manager
### Game Manager
### Game Logic
### Databse