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
