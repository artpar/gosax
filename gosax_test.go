package gosax

import "testing"

func TestToLetterRep(t *testing.T) {
  arr := []float64{7, 1, 4, 4, 4, 4}
  self, _ := NewSax(6, 5, 1e-6)
  letters, _ := self.ToLetterRepresentation(arr)
  if letters != "eacccc" {
    panic("hi")
  }
}

func TestToLetterRepWindow(t *testing.T) {
  arr := []float64{1, 9, 4, 7, 3, 6, 1, 11, 4, 15, 23, 7, 3, 1, 11, 4, 15, 23, 7, 6, 10, 100}

  for ws := 1; ws < 20; ws ++ {
    for i := 3; i < 20; i++ {
      self, err := NewSax(ws, i, 1e-6)
      if err != nil {
        panic(err)
      }
      //self.set_window_size(i)
      letters, _ := self.ToLetterRepresentation(arr)
      t.Logf("Letters [%v][%d]: %v", ws, i, letters)
    }
  }
}

func TestSlidingWindow(t *testing.T) {
  arr := []float64{1, 9, 4, 7, 3, 6, 1, 11, 4, 15, 23, 7, 3, 1, 11, 4, 15, 23, 7, 6, 10, 100}

  for ws := 1; ws < 20; ws ++ {
    for i := 3; i < 20; i++ {
      self, err := NewSax(ws, i, 1e-6)
      if err != nil {
        panic(err)
      }
      //self.set_window_size(i)
      letters, _, _ := self.SlidingWindow(arr, 4, 0.2)
      t.Logf("Letters [%v][%d]: %v", ws, i, letters)
    }
  }
}



func TestLongToLetterRep(t *testing.T) {
  long_arr := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 6, 6, 6, 10, 100}
  self, _ := NewSax(6, 5, 1e-6)
  letters, _ := self.ToLetterRepresentation(long_arr)
  t.Logf("%v == bbbbce", letters)
  if letters != "bbbbce" {
    panic("2")
  }
}

func TestCompareStrings(t *testing.T) {
  self, _ := NewSax(6, 5, 1e-6)
  base_string := "aaabbc"
  similar_string := "aabbbc"
  dissimilar_string := "ccddbc"
  similar_score, _ := self.CompareStrings([]byte(base_string), []byte(similar_string))
  dissimilar_score, _ := self.CompareStrings([]byte(base_string), []byte(dissimilar_string))
  t.Logf("Similarity score: %v", similar_score)
  t.Logf("DisSimilarity score: %v", dissimilar_score)
  if similar_score >= dissimilar_score {
    panic("3")
  }
}