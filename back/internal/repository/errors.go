package repository

import "errors"

// Erreurs métier renvoyées par les repositories, indépendantes du driver.
var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrDuplicate    = errors.New("player name already exists")
)
