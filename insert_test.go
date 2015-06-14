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


	// Test when two arrays need to be merged

	// Test if a free array has been returned

	// Test of the two smallest arrays have been chosen

	// Test if the content of those arrays are kept, and the final
	// array is still sorted

	// Note: to test for same element, use a map[Entry]int with
	// the int value counting how many elements
}
