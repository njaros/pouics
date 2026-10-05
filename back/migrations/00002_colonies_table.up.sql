CREATE TABLE colonies (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	player_id UUID REFERENCES players(id) ON DELETE CASCADE,
	max_size INT,
	pos_x INT,
	pos_y INT,
	target_x INT,
	target_y INT,
	speed FLOAT
);