[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-c66648af7eb3fe8bc4f294546bfd86ef473780cde1dea487d3c4ff354943c9ae.svg)](https://classroom.github.com/online_ide?assignment_repo_id=9360321&assignment_repo_type=AssignmentRepo)
# Project 5: building your very own readmapper

In this final project, you will write a complete read mapper.

The read mapper should be able to preprocess a reference genome. To avoid preprocessing each time you need to map reads, you should store the preprocessed data structures on disk. Reference genomes come in Simple-FASTA format, as usual, and reads in Simple-FASTQ format, and your tool must write matches to standard out in Simple-SAM.

Your program, which should be named `readmap`, and should take the following options:

* `readmap -p genome.fa` should preprocess the genome
* `readmap -d k genome.fa reads.fq` should do read-mapping for matches within an edit distance .


## Assembly required

If you have made all project and all exercises you will have most of what goes into a readmapper.

1. You implemented the file format parsers in the first week of the class, and you have been using them in the four previous projects.
2. If you map using a suffix tree, you have implemented it in project 2.
3. If you map using Li & Durbin’s algorithm you implemented most of the necessary data structures in projects 3 and 4.

## Batteries not included

You have not implemented approximative matching, so you have to implement that now.

## Testing

You can use the [gsa] Python package for generating test data and running tests. You can clone it from the GitHub repository or use:

```bash
> python3 -m pip install git+https://github.com/birc-gsa/gsa#egg=gsa
```

Amongst other things, the tool can simulate data. If you run, for example

```bash
> gsa simulate genome 23 100000 > genome.fa
```

you will simulate a genome with 23 chromosomes, each of length 100,000.

After that,

```bash
> gsa simulate reads genome.fa 2000 100
```

will simulate 2000 reads of length 100.

If you then do

```bash
> gsa search genome.fa reads.fq approx -e 1 bwt
```

to find all the hits within one edit distance of a read. If you want it faster, preprocess the genome first with

```bash
> gsa preprocess genome.fa approx-bwt
```

You should notice a speed difference; you want to achieve the same with your own preprocessing.

You can use the tool to test your read mapper as well. This requires a spec file that defines how tools should be tested. It can look like this:

```yaml
tools:
  GSA:
    preprocess: "gsa preprocess {genome} approx-bwt"
    map: "gsa search {genome} {reads} -o {outfile} approx -e {e} bwt"
  readmap:
    preprocess: "{root}/readmap -p {genome}"
    map: "{root}/readmap -d {e} {genome} {reads} > {outfile}"

reference-tool: GSA

genomes:
  length: [1000, 5000, 10000]
  chromosomes: 10

reads:
  number: 10
  length: 10
  edits: [0, 1, 2]
```

The `tools` section is a list of tools to run, each with a `preprocess` and a `map` command line. You can have as many as you like. The `reference-tool` selects which tool to consider “correct”; all other tools are compared against its results. Then `genomes` specify the genome length and number of chromosomes. Lists here will add a test for each combination. Similarly, the `reads` specify the reads, their number and length and how many edits the simulation and the readmapping will use.

The variables in `{...}` are used by `gsa` when you specify command lines. `{root}` refers to the directory where the YAML file sits, so if your tool and the YAML file are in the same directory, your tool is at `{root}/readmap`. The `{genome}` and `{reads}` tags are the input files and `{outfile}` the name of the output file. Don’t get inventive with the command line for your tool, though, I also have a test ready to run, and if you do not implement the interface specified above, the test will fail (and that will be your problem and not mine).

If you put this file in `tests.yml`, and you have the tool `readmap`, you can run the test with

```bash
> gsa -v test tests.yaml
```

The read mapper in `gsa` doesn’t output matches with leading or training deletions. We talk about why, and how you avoid it as well, in the exercises. Keep that in mind when you are developing your own tool.

## Evaluation

Once you have implemented the `readmap` program (and tested it to the best of your abilities) fill out the report below, and notify me that your pull request is ready for review.

# Report

## Algorithm
We decided to build upon project 4 by using BWT where we recursion in order to explore all possible combinations of M I and D operations within the allowed edit distance. More specifically we implemented the Li-Durbin algorithm which reduces the amount of recursive calls needed by using early stopping.

## Insights you may have had while implementing the algorithm
Life lesson: Some algorithms are easy to understand and hard to implement - and some are hard to understand but easy to implement.
## Problems encountered if any
We tried to implement an algorithm for faster preprocessing (skew, SAIS) but without any luck. Especially the skew attempt worked fine on inputs smaller than 500~ but would eventually die to some bug somewhere. Very sad... 
The read mapper part of the project not too demanding to implement - we did however spend some ting trying to remove the beginning and ending deletions.  

## Validation

First we tried to algorithm on some small homemade genomes and reads where we could easily manually validate the output.
In order to validate further it was no longer super optimal to compare it to a previous project, since all the other project assumes an edit length of 0. We did however start by running our implementation with "-d 0" to provide a sanity check.
Later we installed the gsa tool and created larger files. We compared the results from our tool with the results from gsa and made sure that they were identical.

## Running time
The preprocessing used in this project is the same as the one we used in project 3. It is using radix sort and in the tests made back then we showed that it runs O(n^2) as expected...
Until we changed it because we got skew to work - which on all our comparisons clearly outperformed our standard lsd radix. This makes sense since skew is supposed to run linear O(n)

For our read mapping we use the Burrows Wheeler transformation method which runs in O(n) on preprocessed data. 

We wanted to see how much impact the addition of the D array from li-durbin algorithm had on our runtimes, so we ran the algorithm on different size n with edit distance = 4 to see a difference. The results are shown in the image below, and at appears that there were close to no difference - or atleast the difference was dominated by some other more time demanding things at runtime. Our preprocessing limits the scale of our data at the time of writing, but otherwise it would be interesting to see if this pattern would change if we ran the algorithm on way larger genomes. 
![](figs/lidurbin_stopping.png)