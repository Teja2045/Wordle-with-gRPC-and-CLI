package main

import (
	"context"
	"log"
	pb "wordle-with-gRPC/pbFiles"
	"wordle-with-gRPC/utils/errors"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcServer := "localhost:8082"
	conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error connecting grpc: ", err)
	}
	client := pb.NewGameServiceClient(conn)

	rootCmd := &cobra.Command{
		Use:   "client",
		Short: "root command for the game",
		Long:  `root command for the game`,
	}

	startCmd := &cobra.Command{
		Use:   "start-game",
		Short: "| USER_NAME |",
		Long:  "command to add your name before starting the game",
		Run: func(cmd *cobra.Command, args []string) {
			userName := &pb.UserName{
				UserName: args[0],
			}
			addUser, err := client.Start(context.Background(), userName)
			if err != nil {
				log.Println("Error after starting the game for user: ", err.Error())
			}
			log.Println(addUser)
		},
	}

	submitCmd := &cobra.Command{
		Use:   "submit-word",
		Short: "| USERNAME | WORD |",
		Long:  "command to submit the word you guessed",
		Run: func(cmd *cobra.Command, args []string) {
			wordGuess := &pb.WordGuess{
				UserName: args[0],
				Word:     args[1],
			}

			wordResponse, err := client.Submit(context.Background(), wordGuess)
			if err != nil {
				log.Println("Error after submitting the word:", err.Error())
			}
			log.Println(wordResponse)
		},
	}

	getGameStatusCmd := &cobra.Command{
		Use:   "game-status",
		Short: "|USERNAME|",
		Long:  "command to view your game status",
		Run: func(cmd *cobra.Command, args []string) {
			userName := &pb.UserName{
				UserName: args[0],
			}
			wordResponse, err := client.GetGameStatus(context.Background(), userName)
			if err != nil {
				log.Println("Error while getting game status: ", err.Error())
			}
			log.Println(wordResponse)
		},
	}

	getMyRankCmd := &cobra.Command{
		Use:   "my-rank",
		Short: "| USERNAME |",
		Long:  "command to get your rank for today's game",
		Run: func(cmd *cobra.Command, args []string) {
			userName := &pb.UserName{
				UserName: args[0],
			}
			myRank, err := client.GetMyRank(context.Background(), userName)
			if err != nil {
				log.Println("Error while getting your rank: ", err.Error())
			}
			log.Println(myRank)
		},
	}

	getTodayLeaderBoardCmd := &cobra.Command{
		Use:   "today-leader-board",
		Short: "<No Input>",
		Long:  "Command to get today's leader board",
		Run: func(cmd *cobra.Command, args []string) {
			todayRanks, err := client.GetTodayRanks(context.Background(), &pb.EmptyMessage{})
			if err != nil {
				if err == errors.Err_RANKS_NO_ONE {
					log.Println("________________TODAY-LEADER-BOARD__________________")
					log.Println()
					log.Println()
					log.Println("-NO RANKS-")
				} else {
					log.Println("Error while getting today's leader board: ", err.Error())
				}
			} else {
				log.Println("________________TODAY-LEADER-BOARD___________________")
				log.Println()
				log.Println()
				for _, rank := range todayRanks.Ranks {
					log.Printf("%5d \t", rank.Rank)
					log.Printf("%15s \t", rank.UserName)
					log.Printf("%5s \t", rank.Time.AsTime().Location())
					log.Println()
				}
				log.Println()
			}
		},
	}

	getLeaderBoardCmd := &cobra.Command{
		Use:   "leader-board",
		Short: "<No input>",
		Long:  "Command to all time Leader board",
		Run: func(cmd *cobra.Command, args []string) {
			leaderBoard, err := client.GetRanksHistory(context.Background(), &pb.EmptyMessage{})
			if err != nil {
				log.Println("Error while getting all time leader board: ", err.Error())
			}
			log.Println("________________LEADER-BOARD__________________")
			log.Println()
			log.Println(leaderBoard.AllRanks)
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(getGameStatusCmd)
	rootCmd.AddCommand(getMyRankCmd)
	rootCmd.AddCommand(getTodayLeaderBoardCmd)
	rootCmd.AddCommand(getLeaderBoardCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
