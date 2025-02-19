package wordfilter

import "strings"

type MessageFilter struct {
	WordFilter map[string]bool
}

func (msgHandler *MessageFilter) FilterWord(data string) string {
	msgSplit := strings.Split(data, " ")
	for index, word := range msgSplit {
		word := strings.ToLower(word)
		ok := msgHandler.WordFilter[word]
		if ok {
			msgSplit[index] = "****"
		}
	}
	return strings.Join(msgSplit, " ")
}
