package main

import (
	"flag"
	"fmt"
	"strings"
)

func InvMerge(s1, s2 string, ch chan int) string {
	j := 0
	i := 0
	ret := ""
	inversions := 0
	for { //ever
		if i >= len(s1) && j >= len(s2) {
			break
		}

		// case where we hit end of first list
		if i >= len(s1) {
			ret += string(s2[j])
			j++
			continue
		}

		// case where we hit end of second list
		if j >= len(s2) {
			ret += string(s1[i])
			i++
			continue
		}

		if i < len(s1) && s1[i] <= s2[j] {
			ret += string(s1[i])
			i++
		} else {
			ret += string(s2[j])
			j++
			// here's where the magic happens
			inversions += 1
		}

	}
	ch <- inversions
	return ret
}

func InvMergeSplit(s string, ch chan int) string {
	// We've hit an atom
	if len(s) == 1 {
		return s
	}

	// the two halves of the string
	first := s[:len(s)/2]
	secnd := s[len(s)/2:]

	// Sort the two halves
	return InvMerge(InvMergeSplit(first, ch), InvMergeSplit(secnd, ch), ch)
}

// Modified merge sort which counts the number of inversions.
func InvMergeSort(s string) (int, string) {
	invChan := make(chan int)
	done := make(chan string)

	sum := 1
	go func() {
		for {
			sum += <-invChan
		}
	}()

	go func() {
		done <- InvMergeSplit(s, invChan)
	}()

	sortedString := <-done
	return sum, sortedString
}

// for MSPerms
func myFactorial(i int64) int64 {
	if i < 1 {
		return 1
	}
	return i * myFactorial(i-1)
}

// Multiset permutation computation
func MSPerms(word string) int64 {
	m := make(map[rune]int64)

	// Fill the map with frequencies of chars
	for _, char := range word {
		m[char] = m[char] + 1
	}

	// use the frequencies to determine the denominator for MSP
	fmt.Println()
	fmt.Printf("%d!\n", len(word))
	denom := int64(1)
	for _, v := range m {
		if v != 1 {
			fmt.Printf("%d!", v)
		}
		denom *= myFactorial(v)
	}
	fmt.Println()

	numer := myFactorial(int64(len(word)))
	return numer / denom
}

func nextPermutation(s string) string {
	// find the largest decreasing slice ( index i represents the first elem )
	i := len(s) - 1
	for ; i > 0; i-- {
		if s[i] > s[i-1] {
			break
		}
	}

	// copy over the contents of that slice
	flip := make([]byte, len(s[i:]))
	copy(flip, s[i:])

	// this means we are already at the last permutation, just return the input
	if string(flip) == s {
		return s
	}
	// this will be the prefix to the return string
	first := s[:i-1]

	// the value which will have to change
	pivot := s[i-1]

	// find the first smallest value greater than pivot in flip
	i = 0
	for ; i < len(flip); i++ {
		if pivot >= flip[i] {
			break
		}
	}

	// we decrement to take a step back from our final step above
	i--

	// swap the values
	pivot, flip[i] = flip[i], pivot

	// reverse the order of flip
	for x, y := 0, len(flip)-1; x < y; x, y = x+1, y-1 {
		flip[x], flip[y] = flip[y], flip[x]
	}

	// finally return the next permutation
	return first + string(pivot) + string(flip)
}

// gets the distinct values from string
func distinct(s string) string {
	// fill map with chars
	m := make(map[rune]bool)
	for _, char := range s {
		m[char] = true
	}
	// iterate over map return chars
	ret := ""
	for key, _ := range m {
		ret += string(key)
	}
	return ret
}

// given some input s, and char c, prune the char c from s
func pruneRune(s string, c uint8) string {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			if len(s) == i+1 {
				return s[:i]
			} else {
				return s[:i] + s[i+1:]
			}
		}
	}
	return s
}

func combinatoricRank(s string) int64 {
	// Base case, we have a single character string ( or empty input )
	if len(s) <= 1 {
		fmt.Println("+1")
		return 1
	} else {
		// grab the sorted version of the input
		_, sorted := InvMergeSort(s)

		// grab the distinct letters from the
		dist := distinct(s)

		// split up the head and tail of the input
		head := s[0]
		tail := s[1:]

		// possible combinations on this step
		posCombs := int64(0)

		// we only need to loop over the distinct letters
		for _, char := range dist {
			if char < rune(head) {
				// for each character that comes before the head, we figure out how
				// many possible combinations of words come before it, that is for
				// each letter, determine how many combinatoins can be made without it
				prune := pruneRune(sorted, uint8(char))
				posCombs += int64(MSPerms(prune))
				fmt.Printf("Sorted: %s\n", sorted)
				fmt.Printf("Pruned: %s(%c)\n", prune, char)
			}
		}

		// return all the combinations at this step + the combinations sans head
		return posCombs + combinatoricRank(tail)
	}
}

func main() {
	// Grab the word from the commandline
	var word string
	flag.StringVar(&word, "word", "", "The word rank you want to see")
	flag.Parse()

	if word == "" {
		fmt.Println("You need to specify a word using the -word \"word\" flag")
		return
	}

	// In case we were given some non-uppercase string
	word = strings.ToUpper(word)

	// Run the modified merge sort to count inversions
	invDist, sortedWord := InvMergeSort(word)

	// grab the multiset permutations
	perms := MSPerms(word)

	// THIS IS SLOW!, words up to 13 char words, then gets unbearable....
	// determine the rank using next permutation
	//np := sortedWord
	//rankS := 1
	//for np != word {
	//	np = nextPermutation( np )
	//	rankS++
	//}

	// faster way to do this using a combinatoric analysis
	rankF := combinatoricRank(word)

	// Let everyone know the good news
	fmt.Printf("The word \"%s\" has inversion distance of %d\n", word, invDist)
	fmt.Printf("The word \"%s\" is part of %d unique words\n", word, perms)
	fmt.Printf("The word \"%s\" sorted lexicographically is \"%s\"\n", word, sortedWord)
	//fmt.Printf("The word \"%s\" is of rank(SLOW) %d\n", word, rankS)
	fmt.Printf("The word \"%s\" is of rank(FAST) %d\n", word, rankF)
}
