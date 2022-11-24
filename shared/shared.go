// You can create modules at this level and they will be
// interpreted as under module birc.au.dk, so to import
// package `shared` you need `import "birc.au.dk/gsa/shared"`

package shared

import (
	"fmt"
	"os"
	"strings"
)

func Preprocess(genome string) {

	//	fmt.Println("Preprocessing:", genome)

	p_genomes := GeneralParser(genome, Fasta)
	f, err := os.Create(os.Args[2] + "zz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, gen := range p_genomes {
		var sb strings.Builder
		//add sentinel if missing
		if gen.Rec[len(gen.Rec)-1] != '$' {
			sb.WriteString(gen.Rec)
			sb.WriteRune('$')
			gen.Rec = sb.String()
		}
		sa := LsdRadixSort(gen.Rec)
		bwt, c := FM_build(sa, gen.Rec)

		//preprocessing for RO
		r_gen := ReverseStr(gen.Rec)
		r_gen = r_gen + "$"
		if r_gen[0] == '$' {
			r_gen = r_gen[1:]
		}
		r_sa := LsdRadixSort(r_gen)
		rbwt, _ := FM_build(r_sa, r_gen)

		//write to file
		f.WriteString(">" + gen.Name + "\n")
		f.WriteString("@" + string(bwt) + "\n")
		f.WriteString("_" + string(rbwt) + "\n")
		for k, v := range c {
			f.WriteString("*" + string(k) + fmt.Sprint(v) + "\n")
		}
	}
}

func Readmap(genome, reads string, dist int) {
	//fmt.Println("Redmap genome", genome, "with", reads, "within distance", dist)

	f, err := os.Open(genome + "zz")
	if err != nil {
		panic(err)
	}
	p_genomes := FMParser(f)
	p_reads := GeneralParser(reads, Fastq)

	/*
		fo, err := os.Create("./data/output.txt")
		if err != nil {
			panic(err)
		}*/

	for _, gen := range p_genomes {
		for _, read := range p_reads {

			//first reconstruct SA
			gen.BS = ReverseBWT(gen.Bwt, gen.C, gen.O)

			/*start, end := FM_search(gen, read.Rec)*/

			FM_search_approx(gen, read, dist)

			/*if start != end {

				for i := start; i < end; i++ {

					Sam(read.Name, gen.Name, gen.BS[i], read.Rec)

					res := SamStub(read.Name, gen.Name, gen.BS[i], read.Rec)

					fo.Write([]byte(res))

				}
			}*/
		}
	}

	//fmt.Println(kogglobal)
}
