package main

import (
	"fmt"
	pb "wordle-with-gRPC/pbFiles"
	word_service "wordle-with-gRPC/utils/errors/wordService"
)

func main() {
	str := "astring"
	fmt.Println(int(str[0] - 'a'))

	// fmt.Println(word_service.CheckWord("start", "start"))
	fmt.Println(word_service.CheckWord("start", "stars"))
	//fmt.Println(word_service.CheckWord("start", "aaaaa"))

	maap := map[string]*pb.GameStatus{}
	fmt.Println(maap["teja"] == nil)
}
