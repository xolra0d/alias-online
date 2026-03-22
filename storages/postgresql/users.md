```sql
CREATE TABLE IF NOT EXISTS users (
	id uuid PRIMARY KEY,
	name TEXT NOT NULL,
	secret_hash text NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_seen timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
