CREATE TABLE pouics (
	id SERIAL PRIMARY KEY,
	set_id SERIAL REFERENCES sets(id) ON DELETE CASCADE,
	name VARCHAR(30),
	health INT,
	happyness INT,
	stress INT,
	stats INT ARRAY[2],
	brain FLOAT ARRAY[64][4]
);