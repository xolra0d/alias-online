CREATE TABLE IF NOT EXISTS rooms (
	id TEXT PRIMARY KEY,
	admin UUID NOT NULL references users(id),
	seed INT NOT NULL,
	current_word_index INT NOT NULL DEFAULT 0,
	current_player_id UUID NOT NULL REFERENCES users(id),
	game_state INT DEFAULT 0,
	language TEXT NOT NULL references vocabularies(language),
	rude_words boolean NOT NULL,
	additional_vocabulary TEXT [] NOT NULL DEFAULT '{}',
	clock integer NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT NULL
);
