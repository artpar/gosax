package gosax

import (
  "errors"
  "math"
  "strconv"
  "strings"
)

var ErrDictionarySizeIsNotSupported = errors.New("Dictionary size not supported")
var ErrStringsAreDifferentLength = errors.New("StringsAreDifferentLength")
var ErrOverlapSpecifiedIsNotSmallerThanWindowSize = errors.New("OverlapSpecifiedIsNotSmallerThanWindowSize")

type Sax struct {
  aOffset       rune
  wordSize      int
  alphabetSize  int
  eps           float64
  breakpoints   map[string][]float64

  beta          []float64
  scalingFactor float64
  compareDict   map[string]float64
  windowSize    int
}

func NewSax(wordSize int, alphabetSize int, epsilon float64) (*Sax, error) {
  if alphabetSize < 3 || alphabetSize > 20 {
    return nil, ErrDictionarySizeIsNotSupported
  }

  self := Sax{}

  self.aOffset = 'a'
  self.wordSize = wordSize
  self.alphabetSize = alphabetSize
  self.eps = epsilon
  self.breakpoints = make(map[string][]float64)
  self.breakpoints["3"] = []float64{-0.43, 0.43}
  self.breakpoints["4"] = []float64{-0.67, 0, 0.67}
  self.breakpoints["5"] = []float64{-0.84, -0.25, 0.25, 0.84}
  self.breakpoints["6"] = []float64{-0.97, -0.43, 0, 0.43, 0.97}
  self.breakpoints["7"] = []float64{-1.07, -0.57, -0.18, 0.18, 0.57, 1.07}
  self.breakpoints["8"] = []float64{-1.15, -0.67, -0.32, 0, 0.32, 0.67, 1.15}
  self.breakpoints["9"] = []float64{-1.22, -0.76, -0.43, -0.14, 0.14, 0.43, 0.76, 1.22}
  self.breakpoints["10"] = []float64{-1.28, -0.84, -0.52, -0.25, 0, 0.25, 0.52, 0.84, 1.28}
  self.breakpoints["11"] = []float64{-1.34, -0.91, -0.6, -0.35, -0.11, 0.11, 0.35, 0.6, 0.91, 1.34}
  self.breakpoints["12"] = []float64{-1.38, -0.97, -0.67, -0.43, -0.21, 0, 0.21, 0.43, 0.67, 0.97, 1.38}
  self.breakpoints["13"] = []float64{-1.43, -1.02, -0.74, -0.5, -0.29, -0.1, 0.1, 0.29, 0.5, 0.74, 1.02, 1.43}
  self.breakpoints["14"] = []float64{-1.47, -1.07, -0.79, -0.57, -0.37, -0.18, 0, 0.18, 0.37, 0.57, 0.79, 1.07, 1.47}
  self.breakpoints["15"] = []float64{-1.5, -1.11, -0.84, -0.62, -0.43, -0.25, -0.08, 0.08, 0.25, 0.43, 0.62, 0.84, 1.11, 1.5}
  self.breakpoints["16"] = []float64{-1.53, -1.15, -0.89, -0.67, -0.49, -0.32, -0.16, 0, 0.16, 0.32, 0.49, 0.67, 0.89, 1.15, 1.53}
  self.breakpoints["17"] = []float64{-1.56, -1.19, -0.93, -0.72, -0.54, -0.38, -0.22, -0.07, 0.07, 0.22, 0.38, 0.54, 0.72, 0.93, 1.19, 1.56}
  self.breakpoints["18"] = []float64{-1.59, -1.22, -0.97, -0.76, -0.59, -0.43, -0.28, -0.14, 0, 0.14, 0.28, 0.43, 0.59, 0.76, 0.97, 1.22, 1.59}
  self.breakpoints["19"] = []float64{-1.62, -1.25, -1, -0.8, -0.63, -0.48, -0.34, -0.2, -0.07, 0.07, 0.2, 0.34, 0.48, 0.63, 0.8, 1, 1.25, 1.62}
  self.breakpoints["20"] = []float64{-1.64, -1.28, -1.04, -0.84, -0.67, -0.52, -0.39, -0.25, -0.13, 0, 0.13, 0.25, 0.39, 0.52, 0.67, 0.84, 1.04, 1.28, 1.64}

  self.beta = self.breakpoints[strconv.Itoa(self.alphabetSize)]
  self.build_letter_compare_dict()
  self.scalingFactor = 1

  return &self, nil
}

func mean(xs[]float64) float64 {
  total := 0.0
  for _, v := range xs {
    total += v
  }
  return total / float64(len(xs))
}

func stdDev(numbers []float64, mean float64) float64 {
  total := 0.0
  for _, number := range numbers {
    total += math.Pow(number - mean, 2)
  }
  variance := total / float64(len(numbers) - 1)
  return math.Sqrt(variance)
}


/*
  Function will normalize an array (give it a mean of 0, and a
  standard deviation of 1) unless it's standard deviation is below
  epsilon, in which case it returns an array of zeros the length
  of the original array.
*/
func (self *Sax) normalize(measureList []float64) []float64 {
  var err error
  if len(measureList) < 1 {
    return measureList
  }

  s2 := float64(1)
  m2 := float64(0)
  m1 := mean(measureList)
  if err != nil {
    panic(err)
  }
  s1 := stdDev(measureList, m1)
  if err != nil {
    panic(err)
  }

  stdMultiplier := (s2 / s1)

  for i, m := range measureList {
    measureList[i] = float64(m2 + (m - m1)) * stdMultiplier
  }
  return measureList
}

/*
  Function takes a series of data, x, and transforms it to a string representation
*/
func (self *Sax) ToLetterRepresentation(x []float64) (string, [][]int) {
  paaX, indices := self.toPaa(self.normalize(x))
  self.scalingFactor = math.Sqrt(float64((len(x) * 1.0) / (self.wordSize * 1.0)))
  return self.alphabetize(paaX), indices

}

/*
  Function performs Piecewise Aggregate Approximation on data set, reducing
  the dimension of the dataset x to w discrete levels. returns the reduced
  dimension data set, as well as the indices corresponding to the original
  data for each reduced dimension
*/
func (self *Sax) toPaa(x []float64) (approximation []float64, indices [][]int) {

  n := len(x)
  stepFloat := float64(n) / float64(self.wordSize)
  step := int(math.Ceil(stepFloat))
  frameStart := 0
  approximation = []float64{}
  indices = [][]int{}
  i := 0
  for frameStart <= n - step {
    thisFrame := x[frameStart:int(frameStart + step)]
    m := mean(thisFrame)
    approximation = append(approximation, m)
    indices = append(indices, []int{frameStart, int(frameStart + step)})
    i += 1
    frameStart = int(float64(i) * stepFloat)
  }
  return

}

/*
  Converts the Piecewise Aggregate Approximation of x to a series of letters.
*/
func (self *Sax) alphabetize(paaX []float64) string {
  alphabetizedX := ""
  for i := 0; i < len(paaX); i++ {
    letterFound := false

    for j := 0; j < len(self.beta); j++ {
      if paaX[i] < self.beta[j] {
        alphabetizedX += string(rune(int(self.aOffset) + j))
        letterFound = true
        break
      }

    }
    if !letterFound {
      alphabetizedX += string(rune(int(self.aOffset) + len(self.beta)))
    }

  }
  return alphabetizedX

}

/*
  Compares two strings based on individual letter distance
  Requires that both strings are the same length
*/
func (self *Sax) CompareStrings(list_letters_a, list_letters_b []byte) (float64, error) {
  if len(list_letters_a) != len(list_letters_b) {
    return 0, ErrStringsAreDifferentLength
  }
  mindist := 0.0
  for i := 0; i < len(list_letters_a); i++ {
    mindist += math.Pow(self.compare_letters(rune(list_letters_a[i]), rune(list_letters_b[i])), 2)
  }
  mindist = self.scalingFactor * math.Sqrt(mindist)
  return mindist, nil

}

/*
  Compare two letters based on letter distance return distance between
*/
func (self *Sax) compare_letters(la, lb rune) float64 {
  return self.compareDict[string(la) + string(lb)]
}

func rangeIntArray(start, end int) []int {
  x := make([]int, end - start)

  for i := start; i < end; i++ {
    x[i] = i
  }
  return x
}

/*
  Builds up the lookup table to determine numeric distance between two letters
  given an alphabet size.  Entries for both 'ab' and 'ba' will be created
  and will have identical values.
*/
func (self *Sax) build_letter_compare_dict() {

  number_rep := rangeIntArray(0, self.alphabetSize)
  letters := make([]string, len(number_rep)) // [chr(x + self.aOffset) for x in number_rep]
  for i, x := range number_rep {
    letters[i] = string(rune(x + int(self.aOffset)))
  }

  self.compareDict = make(map[string]float64)
  for i := 0; i < len(letters); i++ {
    for j := 0; j < len(letters); j++ {
      if math.Abs(float64(number_rep[i] - number_rep[j])) <= 1 {
        self.compareDict[letters[i] + letters[j]] = 0
      } else {
        high_num := int(math.Max(float64(number_rep[i]), float64(number_rep[j])) - 1)
        low_num := int(math.Min(float64(number_rep[i]), float64(number_rep[j])))
        self.compareDict[letters[i] + letters[j]] = self.beta[high_num] - self.beta[low_num]
      }
    }
  }
}

func (self *Sax) SlidingWindow(x []float64, numSubsequences int, overlappingFraction float64) (string, [][]int, error) {
  if numSubsequences < 0 {
    numSubsequences = 20
  }
  self.windowSize = int(len(x) / numSubsequences)
  if overlappingFraction < 0 {
    overlappingFraction = 0.9
  }
  windowIndices := make([][]int, 0)

  overlap := self.windowSize * int(overlappingFraction)
  moveSize := int(self.windowSize - overlap)
  if moveSize < 1 {
    return "", windowIndices, ErrOverlapSpecifiedIsNotSmallerThanWindowSize
  }
  ptr := 0
  n := len(x)
  stringRep := make([]string, 0)
  for ptr < n - self.windowSize + 1 {
    thisSubRange := x[ptr : ptr + self.windowSize]
    thisStringRep, _ := self.ToLetterRepresentation(thisSubRange)
    stringRep = append(stringRep, thisStringRep)
    windowIndices = append(windowIndices, []int{ptr, ptr + self.windowSize})
    ptr += moveSize
  }
  return strings.Join(stringRep, ""), windowIndices, nil

}

func (self *Sax) BatchCompare(xStrings []string, refString string) []float64 {

  res := make([]float64, len(xStrings))
  refBytes := []byte(refString)
  for i, st := range xStrings {
    x, err := self.CompareStrings([]byte(st), refBytes)
    if err != nil {
      res[i] = 0

    } else {
      res[i] = x
    }

  }
  return res

}

func (self *Sax) SetScalingFactor(scalingFactor float64) {
  self.scalingFactor = scalingFactor
}

func (self *Sax) SetWindowSize(windowSize int) {
  self.windowSize = windowSize
}
