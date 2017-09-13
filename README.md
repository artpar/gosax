# GoSax

### High performance golang implementation of Symbolic Aggregate approXimation

### Time Series Classification and Clustering with Golang

Based on the paper [A Symbolic Representation of Time Series, with Implications for Streaming Algorithms](http://www.cs.ucr.edu/~eamonn/SAX.pdf)

```
A Symbolic Representation of Time Series, with Implications for
Streaming Algorithms
Jessica Lin Eamonn Keogh Stefano Lonardi Bill Chiu
University of California - Riverside
Computer Science & Engineering Department
Riverside, CA 92521, USA
{jessica, eamonn, stelo, bill}@cs.ucr.edu

Keywords
Time Series, Data Mining, Data Streams, Symbolic, Discretize
```


## Usage

```

package main

import "github.com/artpar/gosax"


func main(){
  long_arr := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 6, 6, 6, 10, 100}
  // sax, _ := gosax.NewSax(wordSize, alphabetSize, epsilon)
  sax, err := gosax.NewSax(6, 5, 1e-6)
  // check for error before proceding

  letters, _ := sax.ToLetterRepresentation(long_arr)
  t.Logf("%v == bbbbce", letters)
}

```


## Benchmark

```
BenchmarkLongToLetterRep-4   	  500000	      3586 ns/op
BenchmarkSlidingWindow-4       	  200000	     11074 ns/op
```

## Parameters

- WordSize
- AlphabetSize
- Epsilon


If you want to compare x1 and x2 (lists of values):

```
x1String, x1Indices = sax.ToLetterRepresentation(x1)
x2String, x2Indices = sax.ToLetterRepresentation(x2)

x1x2ComparisonScore = s.CompareStrings(x1String, x2String)
```

If you want to use the sliding window functionality:

(say you want to break x3 into a lot of subsequences)

can optionally specify the number of subsequences and how much each subsequence overlaps with the previous subsequence

```
x3Strings, x3Indices = sax.SlidingWindow(x3, numSubsequences, overlappingFraction)
```

Then if you wanted to compare each subsequence to another string (say x2):

```
x3x2ComparisonScores = s.BatchCompare(x3,x2)
```


#### Note:

If you haven't generated the strings through the same SAX object, the scaling factor (square root of the length of the input vector over the word size) will be incorrect, you can correct it using:

```
sax.SetScalingFactor(scalingFactor)
```

#### Tests

```
go test
```

#### Benchmark:

```
go test -bench=.
```
