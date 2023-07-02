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
		Use:     "start-game",
		Short:   "| USER_NAME |",
		Long:    "command to add your name before starting the game",
		Example: "client start-game teja",
		Args:    cobra.ExactArgs(1),
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
		Use:     "submit-word",
		Short:   "| USERNAME | WORD |",
		Long:    "command to submit the word you guessed",
		Args:    cobra.ExactArgs(2),
		Example: "client submit-word teja silly",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args[1]) != 5 {
				log.Panic("the word should have 5 letters")
			}
			wordGuess := &pb.WordGuess{
				UserName: args[0],
				Word:     args[1],
			}

			wordResponse, err := client.Submit(context.Background(), wordGuess)
			if err != nil {
				log.Panic("Error after submitting the word:", err.Error())
			}
			log.Println(wordResponse)
		},
	}

	getGameStatusCmd := &cobra.Command{
		Use:     "game-status",
		Short:   "|USERNAME|",
		Long:    "command to view your game status",
		Args:    cobra.ExactArgs(1),
		Example: "client game-status teja",
		Run: func(cmd *cobra.Command, args []string) {
			userName := &pb.UserName{
				UserName: args[0],
			}
			wordResponse, err := client.GetGameStatus(context.Background(), userName)
			if err != nil {
				log.Panic("Error while getting game status: ", err.Error())
			}
			log.Println(wordResponse)
		},
	}

	getMyRankCmd := &cobra.Command{
		Use:     "my-rank",
		Short:   "| USERNAME |",
		Long:    "command to get your rank for today's game",
		Example: "client my-rank teja",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			userName := &pb.UserName{
				UserName: args[0],
			}
			myRank, err := client.GetMyRank(context.Background(), userName)
			if err != nil {
				log.Panic("Error while getting your rank: ", err.Error())
			}
			log.Println(myRank)
		},
	}

	getTodayLeaderBoardCmd := &cobra.Command{
		Use:     "today-leader-board",
		Short:   "<No Input>",
		Long:    "Command to get today's leader board",
		Example: "client today-leader-board",
		Run: func(cmd *cobra.Command, args []string) {
			todayRanks, err := client.GetTodayRanks(context.Background(), &pb.EmptyMessage{})
			if err != nil {
				if err == errors.Err_RANKS_NO_ONE {
					log.Println("________________TODAY-LEADER-BOARD__________________")
					log.Println("-NO RANKS-")
				} else {
					log.Panic("Error while getting today's leader board: ", err.Error())
				}
			} else {
				log.Println("________________TODAY-LEADER-BOARD___________________")
				log.Println()
				for _, rank := range todayRanks.Ranks {
					log.Printf("Rank:%d\t Name:%s\t %s", rank.Rank, rank.UserName, rank.Time.AsTime().Local())
				}
				log.Println()
			}
		},
	}

	getLeaderBoardCmd := &cobra.Command{
		Use:     "leader-board",
		Short:   "<No input>",
		Long:    "Command to all time Leader board",
		Example: "client leader-board",
		Run: func(cmd *cobra.Command, args []string) {
			leaderBoard, err := client.GetRanksHistory(context.Background(), &pb.EmptyMessage{})
			if err != nil {
				log.Panic("Error while getting all time leader board: ", err.Error())
			}
			log.Println("________________LEADER-BOARD__________________")
			log.Println()
			if len(leaderBoard.AllRanks) == 0 {
				log.Println("-NO RANKS-")
			} else {
				for _, dayRank := range leaderBoard.AllRanks {
					if len(dayRank.Ranks) == 0 {
						continue
					}
					day := dayRank.Ranks[0].Time.AsTime().Local().Day()
					month := dayRank.Ranks[0].Time.AsTime().Local().Month()
					year := dayRank.Ranks[0].Time.AsTime().Local().Year()
					log.Printf("DATE: %d-%d-%d", day, month, year)
					for _, rank := range dayRank.Ranks {
						log.Printf("Rank:%d\t Name:%s\t Time:%s", rank.Rank, rank.UserName, rank.Time.AsTime().Local())
					}
					log.Println()
				}
			}
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
