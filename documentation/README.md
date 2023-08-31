<div style="text-align:center; margin: 100px;">
    <img src="./resources/logo.png"/>
</div>

# TicTactics
This application is a web-based implementation of the game Tic-Tactics. The game was originaly created by studio [HiddenVariable](https://www.hiddenvariable.com/), but near 2017 it has been [taken down](https://www.hiddenvariable.com/tictactics/). 

This project is an attempt to recreate the game.

## 🏛️ Architecture

A high level architecture diagram of the service is shown below.

<center>
    <img src="./resources/HighLevelArchitecture.svg"/>
</center>

Below is a description of the components of the service.

### IO
This component is responsible for handling all input and output to and from the service. Depending on the type of request, it will forward the request to the appropriate component.

It will be deployed as a REST API with the following endpoint groups:

<center>

| Endpoint          | Description               |
| -                 | -                         |
| `/auth`           | Authentication endpoints  |
| `/user`           | User management endpoints |
| `/game`           | Game management endpoints |

</center>

#### Authentication Endpoints

<center>

| Endpoint          | Description               |
| -                 | -                         |
| `/auth/register`  | Register a new user       |
| `/auth/login`     | Login a user              |
| `/auth/logout`    | Logout a user             |
| `/auth/validate`  | Validate a user's session |

</center>

#### User Management Endpoints

<center>

| Endpoint          | Description               |
| -                 | -                         |
| `/user/profile`   | Get a user's profile      |

</center>

#### Game Management Endpoints

<center>

| Endpoint          | Description               |
| -                 | -                         |
| `/game/create`    | Create a new game         |
| `/game/join`      | Join an existing game     |
| `/game/leave`     | Leave a game              |
| `/game/move`      | Make a move in a game     |

</center>

### Auth
### User Manager
### Game Manager
### Game Logic
### Databse