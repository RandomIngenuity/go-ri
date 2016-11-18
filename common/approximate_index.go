// Provides an index that allows you to find nearest items. This differs from 
// the built-in search functionality in that it 1) takes a load function, 2) 
// hosts the actual data, and 3) and lazily populates the internal data-
// structure by batches whose size is based on a pre-configured modulous.
package ricommon

import (
    "sort"

    "golang.org/x/net/context"

    "github.com/dsoprea/go-logging"
)

// Other
var (
    aiLog = log.NewLogger("ri.common.approximate_index")
)


type ApproximateIndexEntry interface {
    // Key Return an identifier for the item that we're searching and sorting 
    // by (e.g. time.Time).
    Key() (key interface{})
}

type ApproximateIndexProvider interface {
    // ReadEntries Return a list of entries in the given range (inclusive of 
    // the start but not of the stop).
    ReadEntries(from, to interface{}) (entries []ApproximateIndexEntry, err error)

    // HashIndex return the identifier of the load-group that this entry 
    // belongs to. It should be the same type as Value() (e.g. time.Time)
    Group(value interface{}) (groupId interface{})

    // GroupRange Return the start and stop items for the group with the given
    // key. It must be inclusive of the first and exclusive of the second.
    GroupRange(groupId interface{}) (from, to interface{}, err error)

    // String Return a description of the provider.
    String() string
}

type ApproximateIndex struct {
    ctx context.Context
    provider ApproximateIndexProvider
    groups map[interface{}][]ApproximateIndexEntry
    maxGroups int
    mru []interface{}
    mruLookup map[interface{}]int
}

func NewApproximateIndex(ctx context.Context, provider ApproximateIndexProvider, maxGroups int) *ApproximateIndex {
    groups := make(map[interface{}][]ApproximateIndexEntry)

    return &ApproximateIndex{
        ctx: ctx,
        provider: provider,
        groups: groups,
        maxGroups: maxGroups,
        mru: make([]interface{}, maxGroups),
    }
}

// updateMru Limit the data that we keep based on usage order.
func (ai *ApproximateIndex) updateMru(newGroupId interface{}) {
    aiLog.Debugf(ai.ctx, "Update: [%s] [%v]", ai.provider.String(), newGroupId)

    if ai.maxGroups == 0 {
        return
    }

    len_ := len(ai.mru)

    // Update MRU.
    doPrune := len_ == ai.maxGroups
    oldestGroupId := ai.mru[len_ - 1]

    ai.mru = append([]interface{} { newGroupId }, ai.mru[1:len_]...)

    // Nothing was shifted off the end.
    if doPrune == false {
        return
    }

    // The group that was shifted off the end happened to be the one that was 
    // most recent used. False positive (it's a happy day).
    if oldestGroupId == newGroupId {
        return
    }

    aiLog.Debugf(ai.ctx, "Forgetting: [%s] [%v]", ai.provider.String(), oldestGroupId)

    // Prune the group that has aged-out from the dataset.
    delete(ai.groups, oldestGroupId)
}

// Find Either return a match item or the next smallest one.
func (ai *ApproximateIndex) Find(key interface{}) (aie ApproximateIndexEntry, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    providerDescription := ai.provider.String()

    groupId := ai.provider.Group(key)

    list, isLoaded := ai.groups[groupId]
    if isLoaded == false {
        from, to, err := ai.provider.GroupRange(groupId)
        log.Panic(err)

        aiLog.Debugf(ai.ctx, "Loading data: [%s] [%s] [%v] FROM=[%v] TO=[%v]", providerDescription, key, groupId, from, to)

        if list, err = ai.provider.ReadEntries(from, to); err != nil {
            log.Panic(err)
        } else {
            ai.updateMru(groupId)
            ai.groups[groupId] = list
        }
    } else {
        aiLog.Debugf(ai.ctx, "Data already available: [%s] [%v] [%v]", providerDescription, key, groupId)
    }

    p := func(i int) bool {
        x := list[i]
        currentKey := x.Key()
        return currentKey == key
    }

    len_ := len(list)
    j := sort.Search(len_, p)
    if j == len_ {
        log.Panic(ErrNotFound)
    }

    nearestItem := list[j]

    aiLog.Debugf(ai.ctx, "Nearest item: [%s] [%v] [%v] [%v]", providerDescription, key, groupId, nearestItem)

    return nearestItem, nil
}
