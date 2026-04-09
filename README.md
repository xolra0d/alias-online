Alias Online
---

Explain word without saying it. 

Built with `net/http`, `Redis`, `PostgreSQL`, and `ReactJS`.

## Try it out

It's currently hosted at https://alias.xolra0d.com

## How to run yourself
0. Apply all migrations in `storages/`
1. Build:
```shell
cd backend/
docker build -f services/main/Dockerfile -t ghcr.io/xolra0d/alias-main:latest .
docker build -f services/auth/Dockerfile -t ghcr.io/xolra0d/alias-auth:latest .
docker build -f services/vocab_manager/Dockerfile -t ghcr.io/xolra0d/alias-vocab-manager:latest .
docker build -f services/room_manager/Dockerfile -t ghcr.io/xolra0d/alias-room-manager:latest .
docker build -f services/room_worker/Dockerfile -t ghcr.io/xolra0d/alias-room-worker:latest .
cd ../frontend/
docker build --build-arg VITE_BACKEND_URL=https://YOUR_URL.com -t ghcr.io/xolra0d/alias-frontend:latest .
```
2. Check out which settings are applicable to your situation for each service, and run them.

## Docs

### Backend

[Protos](backend/shared/proto/README.md)

[Shared](backend/shared/pkg/README.md)

[Service/MainGateway](backend/services/main/README.md)

[Service/VocabManager](backend/services/vocab_manager/README.md)

[Service/Auth](backend/services/auth/README.md)

[Service/RoomManager](backend/services/room_manager/README.md)

[Service/RoomWorker](backend/services/room_worker/README.md)
