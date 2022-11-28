// You can create modules at this level and they will be
// interpreted as under module birc.au.dk, so to import
// package `shared` you need `import "birc.au.dk/gsa/shared"`

package shared

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func Preprocess(genome string) {

	p_genomes := GeneralParser(genome, Fasta)
	f, err := os.Create(os.Args[2] + "zz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, gen := range p_genomes {
		//lets not consider empty genomes
		if len(gen.Rec) == 0 {
			continue
		}
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
	f, err := os.Open(genome + "zz")
	if err != nil {
		panic(err)
	}
	p_genomes := FMParser(f)
	p_reads := GeneralParser(reads, Fastq)

	//ensure all go routines terminate
	var wg sync.WaitGroup

	for _, gen := range p_genomes {
		//first reconstruct SA
		gen.BS = ReverseBWT(gen.Bwt, gen.C, gen.O)

		for _, read := range p_reads {
			//lets not consider empty reads
			if len(read.Rec) == 0 {
				continue
			}

			/*Search for matches using
			go routine cheese*/
			wg.Add(1)
			go FM_search_approx(gen, read, dist, &wg)
		}
	}
	wg.Wait()

}
