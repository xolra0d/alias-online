Alias Online
---

Explain word without saying it. 

Built with `net/http`, `PostgreSQL`, and `ReactJS`.

## Try it out

It's currently hosted at https://alias.xolra0d.com

## How to run yourself
0. Apply all migrations in `storages/`
1. Run backend:
```shell
cd backend/
docker build -t alias-backend .
docker run -d -e POSTGRES_URL='YOUR POSTGRES LINK' -e ALLOWED_ORIGINS='YOUR FRONTEND URL' -e WS_ORIGIN_PATTERNS='YOUR FRONTEND URL' -p 27323:8080 --name=alias-backend alias-backend
```
2. Run frontend:
```shell
cd frontend/
docker build --build-arg VITE_BACKEND_URL='YOUR BACKEND URL' -t alias-frontend .
docker run -d -p 27324:80 --name=alias-frontend alias-frontend
```
3. Go to your frontend url and enjoy.

## How it works

Communication:
- FRONTEND - GET /api/available-vocabs > BACKEND > cached
- FRONTEND - POST /api/create-user IP RATE LIMITED > BACKEND > postgres 
- FRONTEND - POST /api/protected/create-room ID RATE LIMITED > BACKEND > postgres 
- FRONTEND <- GET /api/ws/:roomId -> BACKEND > cached // communication over Websockets

### WS messages

Client message types:
- GetState = 0 // REQUEST ask for state
- StartRound = 1 // REQUEST start my round
- GetWord // REQUEST get my word
- TryGuess // REQUEST guess other player's word
- FinishGame // REQUEST finish game
- GetNewWord // REQUEST skip current word, get new word

Server message types:
- NewUpdate = 0 // REQUEST there is new update 
- CurrentState = 1 // RESPONSE here is game state
- YourWord // RESPONSE your word is ..
- WordGuessed // RESPONSE you word is guessed
- RightGuess // RESPONSE your guess is correct
- WrongGuess // RESPONSE your guess is wrong
 
## Docs

### Backend

[Protos](backend/shared/proto/README.md)
[Shared](backend/shared/pkg/README.md)
[Service/MainGateway](backend/services/main/README.md)
[Service/VocabManager](backend/services/vocab_manager/README.md)
[Service/Auth](backend/services/auth/README.md)
[Service/RoomManager](backend/services/room_manager/README.md)
[Service/RoomWorker](backend/services/vocab_worker/README.md)
