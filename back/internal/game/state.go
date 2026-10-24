package game

import (
	"context"
	"log"
	"pouic/internal/model"

	"github.com/jackc/pgx/v5"
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
	s.loadPlayers(ctx)
	s.loadColonies(ctx)
	s.loadPouics(ctx)
}

func (s *State) loadPlayers(ctx context.Context) {
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

func (s *State) loadColonies(ctx context.Context) {
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

func (s *State) loadPouics(ctx context.Context) {
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

func (s State) SaveToDb(ctx context.Context) {
	batch := &pgx.Batch{}
	for _, player := range s.Players {
		for _, colony := range player.Colonies {
			batch.Queue(
				upsertColonyQuery,
				colony.Colony.Id,
				colony.Colony.PlayerId,
				colony.Colony.MaxSize,
				colony.Colony.PosX,
				colony.Colony.PosY,
				colony.Colony.TargetX,
				colony.Colony.TargetY,
				colony.Colony.Speed,
			)
			for _, pouic := range colony.Pouics {
				batch.Queue(
					upsertPouicQuery,
					pouic.Id,
					pouic.ColonyId,
					pouic.PlayerId,
					pouic.Name,
					pouic.Health,
					pouic.Happyness,
					pouic.Stress,
					pouic.Stats,
					pouic.Brain,
				)
			}
		}
	}

	batchResults := s.pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	for _, player := range s.Players {
		for _, colony := range player.Colonies {
			_, err := batchResults.Exec()
			if err != nil {
				log.Printf("Error during upserting the colony: %v. Error: %v", colony, err)
			}
			for _, pouic := range colony.Pouics {
				_, err := batchResults.Exec()
				if err != nil {
					log.Printf("Error during upserting the pouic: %v. Error: %v", pouic, err)
				}
			}
		}
	}
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