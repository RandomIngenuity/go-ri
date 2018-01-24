package ricommon

import (
    "sort"
)

type Uint64Slice []uint64

func (u64s Uint64Slice) Len() int {
    return len(u64s)
}

func (u64s Uint64Slice) Swap(i, j int) {
    u64s[i], u64s[j] = u64s[j], u64s[i]
}

func (u64s Uint64Slice) Less(i, j int) bool {
    return u64s[i] < u64s[j]
}

func (u64s Uint64Slice) Sort() {
    sort.Sort(u64s)
}
