package cron

import (
	"log"
	"wordle-with-gRPC/db"

	"github.com/robfig/cron"
)

type Cron struct {
	myDB *db.InMemoryDB
}

func NewCron(myDB *db.InMemoryDB) *Cron {
	return &Cron{myDB}
}

func (c *Cron) Start() error {
	log.Println("starting cron job...")

	cron := cron.New()

	// run everyday at midnight
	cron.AddFunc("0 0 0 * * *", func() {
		c.AssignWord()
		log.Println("cron", c.myDB.RanksHistory.AllRanks)
		c.myDB.StoreTodaysLeaderBoard()
		log.Println("cronAfter", c.myDB.RanksHistory.AllRanks)
	})

	go cron.Start()
	return nil
}

func (c *Cron) AssignWord() {
	c.myDB.TODAY_WORD = c.myDB.WORDS.GetRandomWord()
	log.Println("New word added", c.myDB.TODAY_WORD)
}
