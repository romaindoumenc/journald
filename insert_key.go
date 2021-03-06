// Journald: in-memory, efficient JSON storage
//
// Routines to insert a new entry in the log

package journald

import (
	"fmt"
)

const (
	// maximum size before creating a new entry array
	newArraySwitch = 1024
	// maximum number of arrays before merging
	maxMergeSwitch = 16
)

// Append an entry to the system. The entry is passed by its
// timestamp, and a map of all field and matching values
func (log *Log) Append(ts int64, records map[string]string) {

	// Append all records to the map
	//fmt.Println("Append records")
	log_records := make([]*Record, 0)
	for field, value := range records {
		log_record := log.append_record(field, value)
		log_records = append(log_records, log_record)
	}

	// Append entry to the LSM-style entries
	//fmt.Println("Append entry")
	last_entry := log.append_entry(ts, log_records)

	// Update pointers in record list
	//fmt.Println("Update pointers")
	for _, rc := range log_records {
		if rc.entryArray != nil {
			rc.entryArray = append(rc.entryArray, last_entry)
		} else {
			if rc.entry == nil {
				rc.entry = last_entry
			} else {
				rc.entryArray = []*Entry{rc.entry, last_entry}
				rc.entry = nil
			}
		}
	}
}

func (log *Log) append_record(field, value string) *Record {
	if log_rec, ok := log.recmap[field+value]; ok {
		// We already have this record, just use the existing one
		return log_rec
	} else {
		// Create a new record, and add it to the table
		new_rec := create_record(field, value)
		// Append to record table
		log.recmap[field+value] = new_rec
		// Append to field table
		if samef_record, ok := log.fieldmap[field]; ok {
			// Go down the chain of records with the same field
			for samef_record.next_field != nil {
				samef_record = samef_record.next_field
			}
			// Update the chain
			samef_record.next_field = new_rec
		} else {
			// Insert the new record
			log.fieldmap[field] = new_rec
		}
		return new_rec
	}
}

func (log *Log) append_entry(ts int64, records []*Record) *Entry {

	// Check if we have enough room in the first array
	if log.currentEntrySize >= newArraySwitch {
		log.createEntryArray()
	}

	entry_array := log.currentEntry

	// Add the entry to the sorted array, making sure it is still
	// sorted. Starting at end since the entries are usually input
	// as chronological order.
	pos := log.currentEntrySize - 1
	for pos >= 0 {
		if entry_array[pos].Timestamp < ts {
			break
		}
		pos--
	}
	pos++

	// Copy the elements after the entry to insert
	if pos < log.currentEntrySize {
		copy(entry_array[pos+1:log.currentEntrySize], entry_array[pos:log.currentEntrySize-1])
	}

	entry_array[pos] = Entry{Timestamp: ts, Records: records}
	log.currentEntrySize++

	return &entry_array[pos]
}

func (lsm *Log) createEntryArray() {

	if lsm.backlogSize < maxMergeSwitch {

		// We still have space for insertion at the end of the array
		lsm.backlogSize++
		lsm.backlogEntryArrays[lsm.backlogSize-1] = make([]Entry, newArraySwitch)
		lsm.currentEntry = lsm.backlogEntryArrays[lsm.backlogSize-1]
		lsm.currentEntrySize = 0

		return
	}

	// If w are out of arrays, merge the two smallest ones

	// Select the two smallest arrays. This is easy since the
	// algorithm will alway merge arrays to the lowest positions
	// will be used first, and have the biggest arrays
	ary1 := lsm.backlogEntryArrays[0]
	var ai1, ai2 int
	for i, ary := range lsm.backlogEntryArrays[1:] {
		// Smaller array => use it
		if len(ary) < len(ary1) {
			ai1 = i + 1
			ary1 = ary
			break
		}
	}
	if ai1 < len(lsm.backlogEntryArrays) - 1 {
		ai2 = ai1 + 1
	} else {
		ai2 = 0
	}
	ary2 := lsm.backlogEntryArrays[ai2]
	fmt.Println("Array Index", ai1, ai2)

	// Merge them
	new_ary := make([]Entry, len(ary1)+len(ary2))
	var i1, i2 int
	for npos := range new_ary {
		if i2 >= len(ary2) || ary1[i1].Timestamp < ary2[i2].Timestamp {
			new_ary[npos] = ary1[i1]
			i1++
		} else {
			new_ary[npos] = ary2[i2]
			i2++
		}
	}

	// Replace them
	lsm.backlogEntryArrays[ai1] = new_ary
	lsm.backlogEntryArrays[ai2] = make([]Entry, newArraySwitch)

	// Write to the newly created array
	lsm.currentEntrySize = 0
	lsm.currentEntry = lsm.backlogEntryArrays[ai2]
}
