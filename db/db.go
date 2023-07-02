package db

import (
	"log"
	"sync"
	pb "wordle-with-gRPC/pbFiles"
	errors "wordle-with-gRPC/utils/errors"
	word_service "wordle-with-gRPC/utils/wordService"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var Db InMemoryDB

type InMemoryDB struct {
	TODAY_WORD         string
	WORDS              *word_service.Words
	TodaysRanks        *pb.DayRanks
	RanksHistory       *pb.RanksHistory
	UserGameStatus     map[string]*pb.Status
	UserWordStatus     map[string]int32
	TodayRanksLock     sync.Mutex
	RanksHistoryLock   sync.Mutex
	UserGameStatusLock sync.Mutex
	UserWordStatusLock sync.Mutex
}

func NewInMemoryDB() InMemoryDB {
	log.Println("database is created...")
	return InMemoryDB{
		TODAY_WORD: "start",
		WORDS: &word_service.Words{
			Words: []string{},
		},
		TodaysRanks: &pb.DayRanks{
			Ranks: []*pb.Rank{},
		},
		RanksHistory: &pb.RanksHistory{
			AllRanks: []*pb.DayRanks{},
		},
		UserGameStatus:     make(map[string]*pb.Status),
		UserWordStatus:     make(map[string]int32),
		TodayRanksLock:     sync.Mutex{},
		RanksHistoryLock:   sync.Mutex{},
		UserGameStatusLock: sync.Mutex{},
		UserWordStatusLock: sync.Mutex{},
	}
}

func (db *InMemoryDB) Start() error {
	if err := db.WORDS.Initialize(); err != nil {
		return err
	}
	// Todo : write logic to get data from a persistent database
	// Todo : need to start a cron JOB
	db.TODAY_WORD = db.WORDS.GetRandomWord()
	log.Println("db is started...")
	return nil
}

func (db *InMemoryDB) User(username string) error {
	if db == nil {
		return errors.Err_DB_EMPTY
	}
	log.Println("user1")
	userGameStatusDB := db.UserGameStatus
	userWordStatusDB := db.UserWordStatus
	if userGameStatusDB[username] != nil {
		return errors.Err_USER_ALREADY_EXIST
	}
	log.Println("user2")
	db.UserGameStatusLock.Lock()
	userGameStatusDB[username] = pb.Status_NOT_ATTEMPTED.Enum()
	db.UserGameStatusLock.Unlock()
	log.Println("user3")
	db.UserWordStatusLock.Lock()
	userWordStatusDB[username] = 0
	db.UserWordStatusLock.Unlock()
	log.Println("user4")
	return nil
}

func (db *InMemoryDB) GetTodayWord() (string, error) {
	if db == nil {
		return "", errors.Err_DB_EMPTY
	}
	return db.TODAY_WORD, nil
}

func (db *InMemoryDB) SubmitWord(userName string, isCorrectGuess bool) error {
	if db == nil {
		return errors.Err_DB_EMPTY
	}

	if db.UserGameStatus[userName] == nil {
		return errors.Err_USER_DOES_NOT_EXIST
	}
	if db.UserWordStatus[userName] == 6 {
		return errors.Err_TRIES_EXHAUSTED
	}
	if db.UserGameStatus[userName] == pb.Status_WON.Enum() {
		return errors.Err_GAME_ALREADY_WON
	}
	db.UserWordStatusLock.Lock()
	db.UserWordStatus[userName]++
	db.UserWordStatusLock.Unlock()

	if db.UserWordStatus[userName] == 1 {
		db.UserGameStatusLock.Lock()
		db.UserGameStatus[userName] = pb.Status_PENDING.Enum()
		db.UserGameStatusLock.Unlock()
	}
	if isCorrectGuess {
		db.UserGameStatusLock.Lock()
		db.UserGameStatus[userName] = pb.Status_WON.Enum()
		db.UserGameStatusLock.Unlock()
		db.TodayRanksLock.Lock()
		db.TodaysRanks.Ranks = append(db.TodaysRanks.Ranks, &pb.Rank{
			UserName: userName,
			Time:     timestamppb.Now(),
			Rank:     int64(len(db.TodaysRanks.Ranks) + 1),
		})
		db.TodayRanksLock.Unlock()
	}
	if db.UserWordStatus[userName] == 6 {
		db.UserGameStatusLock.Lock()
		db.UserGameStatus[userName] = pb.Status_LOST.Enum()
		db.UserGameStatusLock.Unlock()
	}

	return nil
}

func (db *InMemoryDB) GetGameStatus(userName string) (*pb.Status, error) {
	if db == nil {
		return nil, errors.Err_DB_EMPTY
	}
	if db.UserGameStatus[userName] == nil {
		return nil, errors.Err_USER_DOES_NOT_EXIST
	}
	return db.UserGameStatus[userName], nil
}

func (db *InMemoryDB) GetRank(userName string) (*pb.Rank, error) {
	if db == nil {
		return nil, errors.Err_DB_EMPTY
	}
	if db.UserGameStatus[userName] == nil {
		return nil, errors.Err_USER_DOES_NOT_EXIST
	}
	switch db.UserGameStatus[userName] {
	case pb.Status_LOST.Enum():
		return nil, errors.Err_RANK_GAME_LOST
	case pb.Status_NOT_ATTEMPTED.Enum():
		return nil, errors.Err_RANK_GAME_NOT_STARTED
	case pb.Status_PENDING.Enum():
		return nil, errors.Err_RANK_GAME_PENDING
	}
	for _, rank := range db.TodaysRanks.Ranks {
		if rank.UserName == userName {
			return rank, nil
		}
	}
	return nil, errors.Err_UNKNOWN_ERROR
}

func (db *InMemoryDB) GetTodayRanks() (*pb.DayRanks, error) {
	if db == nil {
		return nil, errors.Err_DB_EMPTY
	}
	if len(db.TodaysRanks.Ranks) == 0 {
		return nil, errors.Err_RANKS_NO_ONE
	}
	return db.TodaysRanks, nil
}

func (db *InMemoryDB) GetRanksHistory() (*pb.RanksHistory, error) {
	if db == nil {
		return nil, errors.Err_DB_EMPTY
	}
	return db.RanksHistory, nil
}

func (db *InMemoryDB) StoreTodaysLeaderBoard() {

	db.RanksHistoryLock.Lock()
	db.RanksHistory.AllRanks = append(db.RanksHistory.AllRanks, &pb.DayRanks{Ranks: db.TodaysRanks.Ranks})
	db.RanksHistoryLock.Unlock()

	log.Println("history: ", db.RanksHistory.AllRanks)
	db.TodayRanksLock.Lock()
	db.TodaysRanks.Ranks = []*pb.Rank{}
	db.TodayRanksLock.Unlock()
	log.Println("LeaderBoard reset!")
}
