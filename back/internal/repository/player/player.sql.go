package player

const (
	queryCreate = `
		INSERT INTO players (name, password)
		VALUES ($1, $2)
		RETURNING id, name`

	queryGetById = `
		SELECT id, name
		FROM players
		WHERE id = $1`

	queryGetByName = `
		SELECT id, name
		FROM players
		WHERE id = $1`

	queryGetByNameFull = `
		SELECT *
		FROM players
		WHERE id = $1`
)