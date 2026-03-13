```sql
CREATE TABLE IF NOT EXISTS rooms (
	id TEXT PRIMARY KEY,
	admin UUID NOT NULL references users(id),
	language TEXT NOT NULL references vocabularies(language),
	rude_words boolean NOT NULL,
	additional_vocabulary TEXT [] NOT NULL DEFAULT '{}',
	clock integer NOT NULL,
	finished boolean NOT NULL DEFAULT FALSE,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT NULL
);
```
