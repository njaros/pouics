package game

import (
	"context"
	"fmt"
	"log"
	"pouic/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ColonyPouics struct {
	Colony *model.Colony
	Pouics map[string]*model.Pouic
}

type PlayerColonies struct {
	Player *model.Player
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

// GETTERS

func (s *State) GetColonies() []*model.Colony {
	collect := []*model.Colony{}
	for _, player := range s.Players {
		for _, colony := range player.Colonies {
			collect = append(collect, colony.Colony)
		}
	}
	return collect
}

// DB STUFFS

func (s *State) LoadFromDb(ctx context.Context) {
	s.loadPlayers(ctx)
	s.loadColonies(ctx)
	s.loadPouics(ctx)
	logString := "Players loaded: ["
	for id, player := range s.Players {
		logString += fmt.Sprintf("%s: %v, ", id, player)
	}
	logString += "]\n"
	log.Print(logString)
}

func (s *State) loadPlayers(ctx context.Context) {
	rows, err := s.pool.Query(ctx, loadPlayersQuery)
	if err != nil {
		log.Fatalf("State.loadPlayers Query error: %v", err)
	}

	defer rows.Close()

	for ok := rows.Next(); ok == true; {
		var player model.Player

		if err := rows.Scan(&player.Id, &player.Name); err != nil {
			log.Fatalf("State.loadPlayers Scan error: %v", err)
		}
		s.AddPlayer(&player)
		ok = rows.Next()
	}
	if err:= rows.Err(); err != nil {
		log.Fatalf("State.loadPlayers rows.Err error: %v", err)
	}
}

func (s *State) loadColonies(ctx context.Context) {
	rows, err := s.pool.Query(ctx, loadColoniesQuery)
	if err != nil {
		log.Fatalf("State.loadColonies Query error: %v", err)
	}

	defer rows.Close()

	for ok := rows.Next(); ok == true; {
		var colony model.Colony
		if err := rows.Scan(
			&colony.Id, &colony.PlayerId, &colony.MaxSize,
			&colony.PosX, &colony.PosY, &colony.TargetX,
			&colony.TargetY, &colony.Speed); err != nil {
			log.Fatalf("State.loadColonies Scan error: %v", err)
		}
		s.addColony(&colony)
		ok = rows.Next()
	}
	if err:= rows.Err(); err != nil {
		log.Fatalf("State.loadColonies rows.Err error: %v", err)
	}
}

func (s *State) loadPouics(ctx context.Context) {
	rows, err := s.pool.Query(ctx, loadPouicsQuery)
	if err != nil {
		log.Fatalf("State.loadPouics Query error: %v", err)
	}

	defer rows.Close()

	for ok := rows.Next(); ok == true; {
		var pouic model.Pouic
		if err := rows.Scan(
			&pouic.Id, &pouic.ColonyId, &pouic.PlayerId,
			&pouic.Name, &pouic.Health, &pouic.Happyness,
			&pouic.Stress, &pouic.Stats, &pouic.Brain); err != nil {
			log.Fatalf("State.loadPouics Scan error: %v", err)
		}
		s.addPouic(&pouic)
		ok = rows.Next()
	}
	if err:= rows.Err(); err != nil {
		log.Fatalf("State.loadPouics rows.Err error: %v", err)
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

func (s *State) AddPlayer(player *model.Player) {
	s.Players[player.Id] = PlayerColonies{
		Player: player,
		Colonies: make(map[string]ColonyPouics),
	}
}

func (s *State) addColony(colony *model.Colony) {
	s.Players[colony.PlayerId].Colonies[colony.Id] = ColonyPouics{
		Colony: colony,
		Pouics: make(map[string]*model.Pouic),
	}
}

func (s *State) addPouic(pouic *model.Pouic) {
	s.Players[pouic.PlayerId].Colonies[pouic.ColonyId].Pouics[pouic.Id] = pouic
}