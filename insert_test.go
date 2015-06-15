// Unit tests for the append* functions
package journald

import "testing"

func TestAppRecord(t *testing.T) {
	lg := NewLog()

	// Test for retrieving existing existing record
	rcexp := new(Record)
	lg.recmap["fieldval"] = rcexp
	rcact := lg.append_record("field", "val")
	if rcexp != rcact {
		t.Error("Incorrect existing record retrieved")
	}

	// Test for creating new record with no previous record
	// without the same field already in
	rcact = lg.append_record("field1", "val1")
	if _, ok := lg.fieldmap["field1"]; !ok {
		t.Fail()
	}

	// Test for adding new record when another record with the
	// same field is present
	newrec := lg.append_record("field1", "val2")
	if rcact.next_field != newrec {
		t.Error("Linkage broken between records")
	}
}

// TODO(rdo) add test for appending records
// TODO(rdo) add test for managing entry arrays
func TestCreateEntryArray(t *testing.T) {
	lg := NewLog()

	// Test to add when space is available
	lg.createEntryArray()

	// Test if returns a shared reference to a free array in the
	// backlog entries. No other element should be modified
	if lg.backlogSize != 1 {
		t.Error("Incorrect backlog entry")
	}

	if &lg.backlogEntryArrays[0][0] != &lg.currentEntry[0] {
		t.Error("The current entry does not point to the entry array")
	}

	// Test when two arrays need to be merged
	for i := 0; i < 2*maxMergeSwitch-2; i++ {
		for j := 0; j < newArraySwitch; j++ {
			lg.currentEntry[j] = Entry{Timestamp: int64(j)}
		}
		lg.createEntryArray()

		// Test if a free array has been returned
		if len(lg.currentEntry) != newArraySwitch {
			t.Error("Incorrect number of elements in the returned array")
		}

		for k := 0; k < newArraySwitch; k++ {
			if lg.currentEntry[k].Timestamp != 0 {
				t.Error("Found non-nil entry")
			}
		}
	}

	// Test of the two smallest arrays have been chosen by testing
	// if a resulting array is too large
	for i, ary := range lg.backlogEntryArrays[:maxMergeSwitch-1] {
		if len(ary) != 2*newArraySwitch {
			t.Error("An incorrectly sized array has been found at pos", i, len(ary))
		}
	}
	if len(lg.backlogEntryArrays[maxMergeSwitch-1]) != newArraySwitch {
		t.Error("An incorrectly sized array has been found at pos 15")
	}

	// Test if the content of those arrays are kept, and the final
	// array is still sorted

	// Note: to test for same element, use a map[Entry]int with
	// the int value counting how many elements
}
