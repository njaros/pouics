CREATE TABLE sets (
	id SERIAL PRIMARY KEY,
	player_id SERIAL REFERENCES players(id) ON DELETE CASCADE,
	max_size INT,
	pos_x INT,
	pos_y INT,
	target_x INT,
	target_y INT,
	speed FLOAT
);