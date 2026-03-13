```sql
CREATE TABLE IF NOT EXISTS vocabularies (
	language text PRIMARY KEY, -- e.g., "English"
	primary_words text [] NOT NULL,
	rude_words text [] NOT NULL DEFAULT '{}',
	available boolean NOT NULL DEFAULT TRUE,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT NULL
);
```

```sql
INSERT INTO vocabularies (language, primary_words) VALUES ('Own vocabulary', '{}');
```
