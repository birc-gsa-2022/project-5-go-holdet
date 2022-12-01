package main

import (
	"testing"

	"fmt"

	"birc.au.dk/gsa/shared"
)

func TestSkewMississippi(t *testing.T) {
	x := "mississippi$"
	sa := shared.Skew(x)
	sa2 := shared.LsdRadixSort(x)

	for i, v := range sa {
		if v != sa2[i] {
			t.Error("Suffix array is not identical at idx: ", i)
		}

	}

	fmt.Println(sa, sa2, "hmm")

	for _, v := range sa {
		fmt.Println(x[v:])
	}

}

func TestSkewLargerGenome(t *testing.T) {
	x := "gfjsdnfduuuuccfnsdfjddsfccccccuuauuauauauabbbefdsfndsjgdsnjkfdsbjfdsbfhjdsbfbbbbbbbdbdbdbddbdbdbaSDSAbbbbbbbuuuuubbbbuuubuuuuubbbbbbbbbbbbbbbbbbbuuuuuuuuuuuukdsnuuuuuuuuuuuujkasdjbbbbbbbuu$"
	fmt.Println(len(x))
	sa := shared.Skew(x)
	sa2 := shared.LsdRadixSort(x)

	for i, v := range sa {
		if v != sa2[i] {
			t.Error("Suffix array is not identical at idx: ", i)
			break
		}

	}

}

func TestSkewRandomString(t *testing.T) {
	genSize := 2500

	for i := 0; i < 8; i++ {

		gen_fasta, _ := shared.BuildSomeFastaAndFastq(genSize, 0, 1, shared.DNA, 876)
		x := shared.GeneralParserStub(gen_fasta, shared.Fasta, genSize+1)[0].Rec

		x = x + "$"
		sa := shared.Skew(x)
		sa2 := shared.LsdRadixSort(x)
		for i, v := range sa {
			if v != sa2[i] {
				t.Error("Suffix array is not identical at idx: ", i)
			}

		}
	}
}

/*
func TestMakeDataSearch(t *testing.T) {
	csvFile, err := os.Create("./testdata/search_time.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"x_size", "quadratic"})

	num_of_n := 1000
	num_of_m := 50
	time_sq := 0

	//always use the same genome in order to make the sa process go faster.
	genomes, _ := shared.BuildSomeFastaAndFastq(50000, 0, 1, shared.A, 102)
	parsedGenomes := shared.GeneralParserStub(genomes, shared.Fasta, 50000+1)

	if len(parsedGenomes) != 1 {
		t.Errorf("should only be 1.")
	}
	gen := parsedGenomes[0]
	fmt.Println("creating sa")
	sa := shared.LsdRadixSort(gen.Rec)
	fmt.Println("sa created")
	for i := 1; i < 100; i++ {

		num_of_n += 500
		//num_of_m += 500

		_, reads := shared.BuildSomeFastaAndFastq(50000, num_of_m, 1, shared.A, 102)
		parsedReads := shared.GeneralParserStub(reads, shared.Fastq, 40000*num_of_m+1)

		for i := 0; i < 5; i++ {

			for _, read := range parsedReads {

				var wg sync.WaitGroup
				wg.Add(1)
				shared.FM_search_approx(gen, read, i, &wg)

				fmt.Println("time", int((time_sq)))
				_ = csvwriter.Write([]string{strconv.Itoa(num_of_n), strconv.Itoa(time_sq)})

				csvwriter.Flush()

				time_sq = 0

			}

		}
	}
}
*/
