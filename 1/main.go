package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"sort"
	"strings"
	"errors"
)

const SRC string = "http://norvig.com/ngrams/word.list"
const MIN_LEN int = 2

func main() {
	// from some reasons norgig.com was returning 403 forbidden for my requests with net/http
	// fetchRes, err := fetchWords(SRC)

	// here is a local copy of list
	fetchRes, err := fetchWordsLocally()

	if err != nil {
		fmt.Printf("Fetch Error: %s \n", err)
		return
	}
	
	words := strings.Split(fetchRes, "\r\n")
	result := findLongestAnagramm(words)
	if (len(result) == 0) {
		fmt.Println("Unfortunately no anagramms were found in provided list")
		return
	}

	fmt.Printf("Longest anagram is \"%s\"\n", result)
}

func fetchWordsLocally () (string, error) {
	res, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		return "", err
	}
	return string(res[:]), nil
}

func fetchWords(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", errors.New(response.Status)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body[:]), nil
}

func findLongestAnagramm(words []string) string {
	sort.Sort(ByLen(words))
	for _, word := range words {
		if len(word) <= MIN_LEN {
			continue
		}
		if checkIfAnagram(word, words) {
			return word
		}
	}
	return ""
}

// recursively check if a word is an anagramm composed from words from list
// it checks all possible substrings that has length not less than minimal one
func checkIfAnagram (word string, list []string) bool {
	l := len(word)
	if l < MIN_LEN {
		return false
	}
	for i := MIN_LEN; i < l - MIN_LEN; i++ {
		prefix := word[:i]
		suffix := word[i:]
	
		if isWordInArray(prefix, list) {
			if isWordInArray(suffix, list) || checkIfAnagram(suffix, list) {
				return true
			}
		}
	}
	return false
}

func isWordInArray (word string, words []string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

// for sort strings from longests to shortests
type ByLen []string

func (a ByLen) Len() int { return len(a) }
func (a ByLen) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLen) Less (i, j int) bool { return len(a[i]) > len(a[j]) }