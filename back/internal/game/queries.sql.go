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
)