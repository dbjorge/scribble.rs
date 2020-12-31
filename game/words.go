package game

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gobuffalo/packr/v2"
)

var (
	wordListCache = make(map[string][]string)
	languageMap   = map[string]string{
		"english": "en",
		"italian": "it",
		"german":  "de",
		"french":  "fr",
		"dutch":   "nl",
	}
	wordBox = packr.New("words", "../resources/words")
)

func readWordList(chosenLanguage string) ([]string, error) {
	langFileName := languageMap[chosenLanguage]
	list, available := wordListCache[langFileName]
	if available {
		return list, nil
	}

	wordListFile, pkgerError := wordBox.FindString(langFileName)
	if pkgerError != nil {
		panic(pkgerError)
	}

	tempWords := strings.Split(wordListFile, "\n")
	var words []string
	for _, word := range tempWords {
		word = strings.TrimSpace(word)

		//Newlines will just be empty strings
		if word == "" {
			continue
		}

		//The "i" was "impossible", as in "impossible to draw", tag initially supplied.
		if strings.HasSuffix(word, "#i") {
			continue
		}

		//Since not all words use the tag system, we can just instantly return for words that don't use it.
		lastIndexNumberSign := strings.LastIndex(word, "#")
		if lastIndexNumberSign == -1 {
			words = append(words, word)
		} else {
			words = append(words, word[:lastIndexNumberSign])
		}
	}

	wordListCache[langFileName] = words

	return words, nil
}

// GetRandomWords gets 3 random words for the passed Lobby. The words will be
// chosen from the custom words and the default dictionary, depending on the
// settings specified by the Lobby-Owner.
func GetRandomWords(lobby *Lobby) []string {
	rand.Seed(time.Now().Unix())
	wordsNotToPick := lobby.alreadyUsedWords
	wordOne := getRandomWordWithCustomWordChance(lobby, wordsNotToPick)
	wordsNotToPick = append(wordsNotToPick, wordOne)
	wordTwo := getRandomWordWithCustomWordChance(lobby, wordsNotToPick)
	wordsNotToPick = append(wordsNotToPick, wordTwo)
	wordThree := getRandomWordWithCustomWordChance(lobby, wordsNotToPick)

	return []string{
		wordOne,
		wordTwo,
		wordThree,
	}
}

func getRandomWordWithCustomWordChance(lobby *Lobby, wordsAlreadyUsed []string) string {
	if lobby.CustomWordsChance > 0 && rand.Intn(100)+1 <= lobby.CustomWordsChance {
		unusedCustomWords := filterOutAlreadyUsed(lobby.CustomWords, wordsAlreadyUsed)
		if len(unusedCustomWords) > 0 {
			return getRandomWord(unusedCustomWords)
		}
	}

	unusedStandardWords := filterOutAlreadyUsed(lobby.Words, wordsAlreadyUsed)
	if len(unusedStandardWords) > 0 {
		return getRandomWord(unusedStandardWords)
	}
	return getRandomWord(lobby.Words)
}

func filterOutAlreadyUsed(candidateWords []string, wordsAlreadyUsed []string) []string {
	filteredWords := make([]string, 0, len(candidateWords))
OUTER_LOOP:
	for _, candidateWord := range candidateWords {
		for _, wordAlreadyUsed := range wordsAlreadyUsed {
			if candidateWord == wordAlreadyUsed {
				continue OUTER_LOOP
			}
		}
		filteredWords = append(filteredWords, candidateWord)
	}
	return filteredWords
}

func getRandomWord(candidateWords []string) string {
	return candidateWords[rand.Intn(len(candidateWords))]
}
