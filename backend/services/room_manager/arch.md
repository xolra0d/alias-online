Architecture
---

| Operation     | Redis Commands                                                                            |
|---------------|-------------------------------------------------------------------------------------------|
| Ping          | `SET worker:{id}:addr EX 30` + update local cache                                         |
| RegisterRoom  | `SADD worker:{id}:rooms {room}` + `SET room:{id}:lock {id} EX N`                          |
| ProlongRoom   | `EXPIRE room:{id}:lock N`                                                                 |
| ReleaseRoom   | `SREM worker:{id}:rooms {room}` + `DEL room:{id}:lock`                                    |
| GetRoomWorker | Choose most free worker, try set `room:{id}:lock` with it, if fails, return current value |
