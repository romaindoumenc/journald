// Journald: in-memory, efficient JSON storage
//
// In-memory structures

package journald

// A record is the base element to consider: it simply consists of a
// field name (e.g. process id), and corresponding value (e.g. 1534).
type Record struct {
	Field, Value string  // Actual payload
	cmpval       []byte  // Possibly compressed value
	next_field   *Record // Next record with same field

	// Entry where this record occurs
	// Optimization: most of the time, one entry
	entry      *Entry
	entryArray []*Entry
}

// An entry is a group of records at the same time stamp
type Entry struct {
	// The time stamp attached to the entry
	Timestamp int64
	// The records making this entry
	Records []*Record
}

// A log consists of a set of entries
type Log struct {
	fieldmap map[string]*Record
	recmap   map[string]*Record

	currentEntrySize   int
	currentEntry       []Entry
	backlogSize        int
	backlogEntryArrays [][]Entry
}

func NewLog() Log {
	var bea [][]Entry
	for i := 0; i < maxMergeSwitch; i++ {
		bea = append(bea, nil)
	}
	return Log {
		fieldmap: make(map[string]*Record),
		recmap:make(map[string]*Record),
		currentEntry:make([]Entry,newArraySwitch),
		currentEntrySize:0,
		backlogSize:0,
		backlogEntryArrays:bea,
	}
}
