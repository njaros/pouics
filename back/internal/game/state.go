package game

import (
	"context"
	"log"
	"pouic/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ColonyPouics struct {
	Colony model.Colony
	Pouics map[string]model.Pouic
}

type PlayerColonies struct {
	Player model.Player
	Colonies map[string]ColonyPouics
}

type State struct {
	Players map[string]PlayerColonies
	pool *pgxpool.Pool
}

func NewState(pool *pgxpool.Pool) *State {
	return &State {
		Players: make(map[string]PlayerColonies),
		pool: pool,
	}
}

func (s *State) LoadFromDb(ctx context.Context) {
	

}

func (s *State)loadPlayers(ctx context.Context) {
	rows, err := s.pool.Query(ctx, loadPlayersQuery)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	defer func() {
		rows.Close()
	}()

	var player model.Player

	if err := rows.Scan(player); err != nil {
		log.Fatal("error: %v", err)
	}
	s.AddPlayer(player)
	for ok := rows.Next(); ok == true; {
		if err := rows.Scan(player); err != nil {
			log.Fatal("error: %v", err)
		}
		s.AddPlayer(player)
	}
}

func (s *State)loadColonies(ctx context.Context) {
	rows, err := s.pool.Query(ctx, loadColoniesQuery)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	defer func() {
		rows.Close()
	}()

	var colony model.Colony

	if err := rows.Scan(colony); err != nil {
		log.Fatal("error: %v", err)
	}
	s.addColony(colony)
	for ok := rows.Next(); ok == true; {
		if err := rows.Scan(colony); err != nil {
			log.Fatal("error: %v", err)
		}
		s.addColony(colony)
	}
}

func (s *State)loadPouics(ctx context.Context) {
	rows, err := s.pool.Query(ctx, loadPouicsQuery)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	defer func() {
		rows.Close()
	}()

	var pouic model.Pouic

	if err := rows.Scan(pouic); err != nil {
		log.Fatal("error: %v", err)
	}
	s.addPouic(pouic)
	for ok := rows.Next(); ok == true; {
		if err := rows.Scan(pouic); err != nil {
			log.Fatal("error: %v", err)
		}
		s.addPouic(pouic)
	}
}

func (s State) SaveToDb(ctx context.Context) error {
	return nil
}

func (s *State) AddPlayer(player model.Player) {
	s.Players[player.Id] = PlayerColonies{
		Player: player,
		Colonies: make(map[string]ColonyPouics),
	}
}

func (s *State) addColony(colony model.Colony) {
	s.Players[colony.PlayerId].Colonies[colony.Id] = ColonyPouics{
		Colony: colony,
		Pouics: make(map[string]model.Pouic),
	}
}

func (s *State) addPouic(pouic model.Pouic) {
	s.Players[pouic.PlayerId].Colonies[pouic.ColonyId].Pouics[pouic.Id] = pouic
}