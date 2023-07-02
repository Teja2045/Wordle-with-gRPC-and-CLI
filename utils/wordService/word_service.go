package word_service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"
	pb "wordle-with-gRPC/pbFiles"
)

type Words struct {
	Words []string
}

func (words *Words) Initialize() error {
	_, err := os.Stat("../utils/wordService/words.txt")

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
	} else {
		fmt.Println("file exists")
	}
	file, err := os.Open("../utils/wordService/words.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	newWords := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		newWords = append(newWords, word)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	words.Words = newWords
	return nil
}

func (words *Words) GetRandomWord() string {
	year, month, day := time.Now().Date()
	refDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	targetDate := time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.UTC)
	duration := targetDate.Sub(refDate)
	days := int(duration.Hours() / 24)
	randomIndex := days % len(words.Words)
	return words.Words[randomIndex]
}

func CheckWord(targetWord, submitWord string) (bool, *pb.WordStatus) {
	if targetWord == submitWord {
		return true, &pb.WordStatus{
			FirstCharacter:  pb.CharacterStatus_CORRECT,
			SecondCharacter: pb.CharacterStatus_CORRECT,
			ThirdCharacter:  pb.CharacterStatus_CORRECT,
			FourthCharacter: pb.CharacterStatus_CORRECT,
			FifthCharacter:  pb.CharacterStatus_CORRECT,
		}
	}
	targetWordCharacterCount := [26]int{}
	for i := 0; i < 5; i++ {
		value := int(targetWord[i] - 'a')
		targetWordCharacterCount[value]++
	}
	result := &pb.WordStatus{}
	if submitWord[0] == targetWord[0] {
		result.FirstCharacter = pb.CharacterStatus_CORRECT
		value := int(submitWord[0] - 'a')
		targetWordCharacterCount[value]--
	}
	if submitWord[1] == targetWord[1] {
		result.SecondCharacter = pb.CharacterStatus_CORRECT
		value := int(submitWord[1] - 'a')
		targetWordCharacterCount[value]--
	}
	if submitWord[2] == targetWord[2] {
		result.ThirdCharacter = pb.CharacterStatus_CORRECT
		value := int(submitWord[2] - 'a')
		targetWordCharacterCount[value]--
	}
	if submitWord[3] == targetWord[3] {
		result.FourthCharacter = pb.CharacterStatus_CORRECT
		value := int(submitWord[3] - 'a')
		targetWordCharacterCount[value]--
	}
	if submitWord[4] == targetWord[4] {
		result.FifthCharacter = pb.CharacterStatus_CORRECT
		value := int(submitWord[4] - 'a')
		targetWordCharacterCount[value]--
	}

	if submitWord[0] != targetWord[0] {
		value := int(submitWord[0] - 'a')
		if targetWordCharacterCount[value] > 0 {
			targetWordCharacterCount[value]--
			result.FirstCharacter = pb.CharacterStatus_PRESENT_BUT_MISPLACED
		} else {
			result.FirstCharacter = pb.CharacterStatus_NOT_PRESENT
		}
	}
	if submitWord[1] != targetWord[1] {
		value := int(submitWord[1] - 'a')
		if targetWordCharacterCount[value] > 0 {
			targetWordCharacterCount[value]--
			result.SecondCharacter = pb.CharacterStatus_PRESENT_BUT_MISPLACED
		} else {
			result.SecondCharacter = pb.CharacterStatus_NOT_PRESENT
		}
	}
	if submitWord[2] != targetWord[2] {
		value := int(submitWord[2] - 'a')
		if targetWordCharacterCount[value] > 0 {
			targetWordCharacterCount[value]--
			result.ThirdCharacter = pb.CharacterStatus_PRESENT_BUT_MISPLACED
		} else {
			result.ThirdCharacter = pb.CharacterStatus_NOT_PRESENT
		}
	}
	if submitWord[3] != targetWord[3] {
		value := int(submitWord[3] - 'a')
		if targetWordCharacterCount[value] > 0 {
			targetWordCharacterCount[value]--
			result.FourthCharacter = pb.CharacterStatus_PRESENT_BUT_MISPLACED
		} else {
			result.FourthCharacter = pb.CharacterStatus_NOT_PRESENT
		}
	}
	if submitWord[4] != targetWord[4] {
		value := int(submitWord[4] - 'a')
		if targetWordCharacterCount[value] > 0 {
			targetWordCharacterCount[value]--
			result.FifthCharacter = pb.CharacterStatus_PRESENT_BUT_MISPLACED
		} else {
			result.FifthCharacter = pb.CharacterStatus_NOT_PRESENT
		}
	}
	return false, result
}
