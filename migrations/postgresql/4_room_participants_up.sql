CREATE TABLE IF NOT EXISTS room_participants (
   room_id TEXT REFERENCES rooms(id),
   user_login TEXT REFERENCES users(login),
   words_tried INT DEFAULT 0,
   words_guessed INT DEFAULT 0,
   turn_order SERIAL,
   PRIMARY KEY (room_id, user_login)
);
