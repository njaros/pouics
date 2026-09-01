CREATE TABLE pouics (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	set_id UUID REFERENCES sets(id) ON DELETE CASCADE,
	name VARCHAR(30),
	health INT,
	happyness INT,
	stress INT,
	stats INT ARRAY[2],
	brain FLOAT ARRAY[64][4]
);