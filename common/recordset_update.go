package ricommon

import (
    "fmt"

    "golang.org/x/net/context"

    "github.com/dsoprea/go-logging"
)

const (
    TargetReadBufferCount = 100
    SourceReadBufferCount = 100
)

// Misc
var (
    ruLog = log.NewLogger("ri.common.recordset_update")
)

type RecordsetRecord interface {
    // Id Describes a static identifer that identifies the entity that this 
    // record represents. Used for membership comparisons.
    Id() string
    
    // Used to compare two separate versions of the data for the same entity.
    IsUnchanged(olderRecord RecordsetRecord) bool

    String() string
}

type RecordsetDatasource interface {
    ReadSource(sourceSet chan<- interface{}) (err error)
    ReadTarget(targetSet chan<- interface{}) (err error)
    String() string
}

type RecordsetUpdaterByList interface {
    ProcessInsert(record RecordsetRecord) (err error)
    ProcessUpdate(record RecordsetRecord) (err error)
    ProcessDelete(record RecordsetRecord) (err error)
    Flush() (err error)
}

type RecordsetDiff struct {
    New []RecordsetRecord
    Updated []RecordsetRecord
    Deleted []RecordsetRecord
}

func (rd *RecordsetDiff) Count() int {
    return len(rd.New) + len(rd.Updated) + len(rd.Deleted)
}

type RecordsetUpdate struct {
    ctx context.Context
}

func NewRecordsetUpdate(ctx context.Context) *RecordsetUpdate {
    return &RecordsetUpdate{
        ctx: ctx,
    }
}

func (ru *RecordsetUpdate) Diff(rd RecordsetDatasource) (diff *RecordsetDiff, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
            ruLog.Errorf(ru.ctx, nil, "Diff failed: [%s]", err)
        }
    }()

    // Load lookup for existing records.

    target := make(chan interface{}, TargetReadBufferCount)

    if err := rd.ReadTarget(target); err != nil {
        log.Panic(err)
    }

    stored := make(map[string]RecordsetRecord)
    for {
        if x, ok := <-target; ok == true {
            switch t := x.(type) {
                case RecordsetRecord:
                    r := x.(RecordsetRecord)
                    stored[r.Id()] = r
                case error:
                    log.Panic(x)
                default:
                    log.Panic(fmt.Errorf("source value not valid: [%s]", t))
            }
        } else {
            break
        }
    }

    // Calculate deltas.

    diff = new(RecordsetDiff)
    diff.New = make([]RecordsetRecord, 0)
    diff.Updated = make([]RecordsetRecord, 0)
    diff.Deleted = make([]RecordsetRecord, 0)

    source := make(chan interface{}, SourceReadBufferCount)
    if err := rd.ReadSource(source); err != nil {
        log.Panic(err)
    }

    for {
        if x, ok := <-source; ok == true {
            switch t := x.(type) {
                case RecordsetRecord:
                    r := x.(RecordsetRecord)

                    if olderRecord, exists := stored[r.Id()]; exists == false {
                        diff.New = append(diff.New, r)
                    } else if r.IsUnchanged(olderRecord) == false {
                        diff.Updated = append(diff.Updated, r)
                    } else {
                        delete(stored, r.Id())
                    }
                case error:
                    log.Panic(x)
                default:
                    log.Panic(fmt.Errorf("source value not valid: [%s]", t))
            }
        } else {
            break
        }
    }

    for _, record := range stored {
        diff.Deleted = append(diff.Deleted, record)
    }

    ruLog.Infof(ru.ctx, "(%d) changes are required for [%s].", diff.Count(), rd)

    return diff, nil
}

func (ru *RecordsetUpdate) Apply(diff *RecordsetDiff, rul RecordsetUpdaterByList) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
            ruLog.Errorf(ru.ctx, nil, "Could not apply changes: [%s]", err)
        }
    }()

    for _, r := range diff.New {
        ruLog.Infof(ru.ctx, "INSERT [%s]: [%s]", r.Id(), r)
        if err := rul.ProcessInsert(r); err != nil {
            log.Panic(err)
        }
    }

    for _, r := range diff.Updated {
        ruLog.Infof(ru.ctx, "UPDATE [%s]: [%s]", r.Id(), r)
        if err := rul.ProcessUpdate(r); err != nil {
            log.Panic(err)
        }
    }

    for _, r := range diff.Deleted {
        ruLog.Infof(ru.ctx, "DELETE [%s]: [%s]", r.Id(), r)
        if err := rul.ProcessDelete(r); err != nil {
            log.Panic(err)
        }
    }

    if err := rul.Flush(); err != nil {
        log.Panic(err)
    }

    return nil
}
