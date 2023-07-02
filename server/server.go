package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"wordle-with-gRPC/cron"
	"wordle-with-gRPC/db"
	pb "wordle-with-gRPC/pbFiles"
	word_service "wordle-with-gRPC/utils/wordService"

	"google.golang.org/grpc"
)

var myDB db.InMemoryDB

type WordleServer struct {
	pb.UnimplementedGameServiceServer
}

func main() {

	myDB = db.NewInMemoryDB()
	if err := myDB.Start(); err != nil {
		log.Fatal("db start error", err)
	}
	log.Println("starting the game server...")

	cron := cron.NewCron(&myDB)
	if err := cron.Start(); err != nil {
		log.Fatal("cron start error", err)
	}

	listen, err := net.Listen("tcp", ":8082")

	if err != nil {
		log.Fatal("listen failed", err)
	}

	server := grpc.NewServer()

	pb.RegisterGameServiceServer(server, &WordleServer{})
	log.Println("Starting grpc server on port :8082")
	server.Serve(listen)
}

func (server *WordleServer) Submit(ctx context.Context, wordGuess *pb.WordGuess) (*pb.WordGuessResponse, error) {
	todayWord := myDB.TODAY_WORD
	isCorrect, wordStatus := word_service.CheckWord(todayWord, wordGuess.Word)
	if err := myDB.SubmitWord(wordGuess.UserName, isCorrect); err != nil {
		return nil, err
	}
	return &pb.WordGuessResponse{
		WordMatch:  wordStatus,
		GameStatus: *myDB.UserGameStatus[wordGuess.UserName],
	}, nil
}

func (server *WordleServer) Start(ctx context.Context, userName *pb.UserName) (*pb.StartResponse, error) {
	fmt.Println("servr1")
	if err := myDB.User(userName.UserName); err != nil {
		return nil, err
	}
	fmt.Println("servr2")
	return &pb.StartResponse{
		StartResponse: "User is added",
	}, nil
}

func (server *WordleServer) GetGameStatus(ctx context.Context, userName *pb.UserName) (*pb.GameStatus, error) {
	status, err := myDB.GetGameStatus(userName.UserName)
	if err != nil {
		return nil, err
	}
	return &pb.GameStatus{
		GameStatus: *status,
	}, nil
}

func (server *WordleServer) GetMyRank(ctx context.Context, userName *pb.UserName) (*pb.Rank, error) {
	return myDB.GetRank(userName.UserName)
}

func (server *WordleServer) GetTodayRanks(ctx context.Context, empty *pb.EmptyMessage) (*pb.DayRanks, error) {
	return myDB.GetTodayRanks()
}

func (server *WordleServer) GetRanksHistory(ctx context.Context, empty *pb.EmptyMessage) (*pb.RanksHistory, error) {
	return myDB.GetRanksHistory()
}
