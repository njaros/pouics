package game

const (
	loadPlayersQuery = `
		SELECT id, name FROM players
	`

	loadColoniesQuery = `
		SELECT * from colonies
	`

	loadPouicsQuery = `
		SELECT * from pouics
	`

	upsertColonyQuery = `
		INSERT INTO colonie (id, player_id, max_size, pos_x, pos_y, target_x, target_y, speed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT(id) DO UPDATE
			SET player_id = excluded.player_id,
				max_size = excluded.max_size,
				pos_x = excluded.pos_x,
				pos_y = excluded.pos_y,
				target_x = excluded.target_x,
				target_y = excluded.target_y,
				speed = excluded.speed;
	`

	upsertPouicQuery = `
		INSERT INTO pouics (id, colony_id, player_id, name, health, happyness, stress, stats, brain)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT(id) DO UPDATE
			SET colony_id = excluded.colony_id,
				player_id = excluded.player_id,
				name = excluded.name,
				health = excluded.health,
				happyness = excluded.happyness,
				stress = excluded.stress,
				stats = excluded.stats,
				brain = excluded.brain;
	`
)