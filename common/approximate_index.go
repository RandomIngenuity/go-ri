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

    // AreGroupIdsEqual Return whether the two group-IDs are equal.
    AreGroupIdsEqual(a, b interface{}) bool

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
        mru: make([]interface{}, 0, maxGroups),
    }
}

// updateMru Limit the data that we keep based on usage order.
func (ai *ApproximateIndex) updateMru(newGroupId interface{}) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    len_ := len(ai.mru)
    
    p := ai.provider.String()
    aiLog.Debugf(ai.ctx, "Update: [%s] G=[%v] MAX=(%d) CUR=(%d)", p, newGroupId, ai.maxGroups, len_)

    if ai.maxGroups == 0 {
        return nil
    }

    // Update MRU.
    doPrune := len_ == ai.maxGroups

    var oldestGroupId interface{}
    if len_ > 0 {
        oldestGroupId = ai.mru[len_ - 1]
        ai.mru = append([]interface{} { newGroupId }, ai.mru[1:len_ - 1]...)
    } else {
        ai.mru = []interface{} { newGroupId }
    }

    // Nothing was shifted off the end.
    if doPrune == false {
        return nil
    }

    // The group that was shifted off the end happened to be the one that was 
    // most recent used. False positive (it's a happy day).
    if ai.provider.AreGroupIdsEqual(oldestGroupId, newGroupId) == true {
        return nil
    }

    aiLog.Debugf(ai.ctx, "Forgetting: [%s] [%v]", p, oldestGroupId)

    // Prune the group that has aged-out from the dataset.
    delete(ai.groups, oldestGroupId)

    return nil
}

// Find Either return a match item or the next smallest one.
func (ai *ApproximateIndex) Find(key interface{}) (aie ApproximateIndexEntry, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
            aiLog.Errorf(ai.ctx, nil, "Find failed: [%v]", key)
        }
    }()

    aiLog.Debugf(ai.ctx, "Find: [%v]", key)
    providerDescription := ai.provider.String()

    groupId := ai.provider.Group(key)
    aiLog.Debugf(ai.ctx, "Resolved key to a group: [%s] => [%s]", key, groupId)

    list, isLoaded := ai.groups[groupId]
    if isLoaded == false {
        aiLog.Debugf(ai.ctx, "Fault. We need data for group: [%s] [%s]", providerDescription, groupId)

        from, to, err := ai.provider.GroupRange(groupId)
        log.PanicIf(err)

        aiLog.Debugf(ai.ctx, "Loading data: [%s] K=[%s] G=[%v] FROM=[%v] TO=[%v]", providerDescription, key, groupId, from, to)

        if list, err = ai.provider.ReadEntries(from, to); err != nil {
            log.Panic(err)
        } else {
            if err := ai.updateMru(groupId); err != nil {
                log.Panic(err)
            }

            ai.groups[groupId] = list
        }

        aiLog.Debugf(ai.ctx, "Records returned: [%s] (%d)", providerDescription, len(list))
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
