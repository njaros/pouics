CREATE TABLE colonies (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	player_id UUID REFERENCES players(id) ON DELETE CASCADE,
	max_size INT,
	pos_x FLOAT,
	pos_y FLOAT,
	target_x FLOAT,
	target_y FLOAT,
	speed FLOAT
);