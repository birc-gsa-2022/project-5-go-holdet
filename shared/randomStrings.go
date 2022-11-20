package shared

import (
	"math/rand"
	"strconv"
	"strings"
)

type Alphabet string

const (
	English Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DNA     Alphabet = "ACTG"
	AB      Alphabet = "ab"
	A       Alphabet = "a"
)

func randString(number int, alphabet Alphabet) string {
	var letters = alphabet
	b := make([]byte, number)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func BuildSomeFastaAndFastq(len_Fasta int, len_Fastq int, amount int, alphabet Alphabet, seed int64) (string, string) {

	rand.Seed(seed)

	var sb strings.Builder
	var sc strings.Builder

	for i := 0; i < amount; i++ {
		sb.WriteString("> chr" + strconv.Itoa(i) + "\n")

		b_str := randString(len_Fasta, alphabet)
		sb.WriteString(b_str + "\n")

		sc.WriteString("@read" + strconv.Itoa(i) + "\n")

		dif := len_Fasta - len_Fastq
		idx := 0
		//allows for fasta and fastq to be same len. rand.intn panics for n=0 input.
		if dif > 0 {
			idx = rand.Intn(len_Fasta - len_Fastq)

		}
		c_str := b_str[idx:(idx + len_Fastq)]

		sc.WriteString(c_str + "\n")
	}

	return sb.String(), sc.String()
}

func BuildARepetitiveFastaAndFastq(repetitions int, len_Fastq int, seed int64) (string, string) {
	rand.Seed(seed)

	var sb strings.Builder
	var sc strings.Builder

	single_rep := strings.Repeat("a", len_Fastq-1)

	sb.WriteString("> chr" + strconv.Itoa(0) + "\n")
	sb.WriteString(strings.Repeat(single_rep+"c", repetitions))

	sc.WriteString("@ read" + strconv.Itoa(0) + "\n")
	sc.WriteString(single_rep + "c")

	return sb.String(), sc.String()

}

func Handin1_ba(parsedGenomes []Recs, parsedReads []Recs) string {
	var rb strings.Builder

	for _, read := range parsedReads {
		for _, gen := range parsedGenomes {
			matches := naive(gen.Rec, read.Rec)
			for _, match := range matches {
				rb.WriteString(SamStub(read.Name, gen.Name, match, read.Rec))
			}
		}
	}
	return rb.String()
}

func naive(x string, p string) (matches []int) {
	if p == "" {
		return matches
	}
outer_loop:
	for i := 0; i < len(x)-len(p)+1; i++ {
		for j, char := range []byte(p) {
			if char != x[i+j] {
				continue outer_loop
			}
		}
		matches = append(matches, i)
	}

	return matches
}
