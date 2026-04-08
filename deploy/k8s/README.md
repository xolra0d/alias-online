# Kubernetes deploy (autoscaling + gRPC load balancing)

This setup deploys all services with horizontal autoscaling and gRPC client-side balancing.

## What this deploy gives you

1. **Autoscaling** via HPA for `main`, `auth`, `vocab-manager`, and `room-worker`.
2. **gRPC balancing** with `round_robin` policy + `dns:///...` service endpoints.
3. **Worker discovery in Redis**: room workers are automatically registered in Redis list `workers` on heartbeat (`PingWorker`).

## Prerequisites

1. Kubernetes cluster with:
   - Ingress controller (nginx)
   - Metrics server (for HPA)
2. Image registry access (example below uses `ghcr.io/xolra0d/...`).
3. PostgreSQL migrations applied from `migrations/postgresql`.

## Build and push images

Run from `backend/`:

```bash
docker build -f services/main/Dockerfile -t ghcr.io/xolra0d/alias-main:latest .
docker build -f services/auth/Dockerfile -t ghcr.io/xolra0d/alias-auth:latest .
docker build -f services/vocab_manager/Dockerfile -t ghcr.io/xolra0d/alias-vocab-manager:latest .
docker build -f services/room_manager/Dockerfile -t ghcr.io/xolra0d/alias-room-manager:latest .
docker build -f services/room_worker/Dockerfile -t ghcr.io/xolra0d/alias-room-worker:latest .
```

Run from `frontend/`:

```bash
docker build --build-arg VITE_BACKEND_URL=https://alias.example.com -t ghcr.io/xolra0d/alias-frontend:latest .
```

Push all images to your registry.

## Configure secrets and env

1. Copy `deploy/k8s/config-and-secrets.example.yaml`.
2. Fill:
   - `POSTGRES_URL`
   - `REDIS_PASSWORD` (if used)
   - `private.pem` and `public.pem`
   - domain-related values (`ALLOWED_ORIGINS`, `WS_ORIGIN_PATTERNS`, `JWT_COOKIE_DOMAIN`)

## Apply manifests

```bash
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/config-and-secrets.example.yaml
kubectl apply -f deploy/k8s/infra.yaml
kubectl apply -f deploy/k8s/apps.yaml
kubectl apply -f deploy/k8s/hpa.yaml
kubectl apply -f deploy/k8s/ingress.yaml
```

## Apply PostgreSQL migrations

Apply:

```text
migrations/postgresql/1_users_up.sql
migrations/postgresql/2_vocabs_up.sql
migrations/postgresql/3_rooms_up.sql
migrations/postgresql/4_room_participants_up.sql
migrations/postgresql/vocab_watcher.sql
```

## Verify worker discovery and scaling

1. Scale room workers (or wait for HPA):
   ```bash
   kubectl -n alias-online scale deployment room-worker --replicas=5
   ```
2. Check Redis workers list:
   ```bash
   kubectl -n alias-online exec deploy/redis -- redis-cli LRANGE workers 0 -1
   ```
3. Check HPAs:
   ```bash
   kubectl -n alias-online get hpa
   ```

## Important networking note (direct worker mode)

This project currently returns each worker address directly to the browser.  
In this mode, users must be able to reach each worker address (`WORKER_PUBLIC_ADDR`).
