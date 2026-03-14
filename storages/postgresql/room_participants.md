```sql
CREATE TABLE IF NOT EXISTS room_participants (
   room_id TEXT REFERENCES rooms(id),
   user_id UUID REFERENCES users(id),
   score INT DEFAULT 0,
   online BOOLEAN DEFAULT FALSE,
   joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (room_id, user_id)
);

```
