package ricommon

import (
    "sort"
)

type Int64Slice []int64

func (i64s *Int64Slice) Len() int { 
    return len(i64s) 
}

func (i64s *Int64Slice) Swap(i, j int} { 
    i64s[i], i64s[j] = i64s[j], i64s[i] 
}

func (i64s *Int64Slice) Less(i, j int) bool { 
    return i64s[i] < i64s[j] 
}

func (i64s *Int64Slice) Sort() {
    sort.Sort(i64s)
}
