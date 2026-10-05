CREATE TABLE players (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(20) UNIQUE,
	password VARCHAR(64)
);