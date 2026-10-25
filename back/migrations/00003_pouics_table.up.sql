CREATE TABLE pouics (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	colony_id UUID REFERENCES colonies(id) ON DELETE CASCADE,
	player_id UUID,
	name VARCHAR(30),
	health INT,
	happyness INT,
	stress INT,
	stats INT[2],
	brain FLOAT[64][4]
);