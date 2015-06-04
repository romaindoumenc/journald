// Unit test for structures
package journald

import "testing"

func TestNonCompressedRecord(t *testing.T) {
	rc := create_record("field", "val")
	if rc.Field != "field" {
		t.Error()
	}

	if rc.Value != "val" {
		t.Error()
	}
}

func TestCompressedRecord(t *testing.T) {
	rc := create_record("field", "This is a very long string that should get compressed eventually. Heck yeah, we want to compress long strings that can for example be large queries sent to the database. To do so, we use the snappy library, that should provide a reasonably fast implementation of a compression for relatively good performance. Another level of compression is of course the record de-duplication offered by the way this database is modelled. Interestingly, the code in the debug.go file generate a wonderfull graphviz-compatible dot file, very usefull for visual debugging. I have no clue how this approach will scale, though, since graphs are notoriously bad when many elements are present — they are called spaghetti chart for a reason :)")

	if rc.Field != "field" {
		t.Error()
	}
	if rc.cmpval == nil {
		t.Error()
	}
}

func TestGetSingleEntry(t *testing.T) {
	et := Entry{42, nil}
	rc := Record{entry: &et}

	act_entries := rc.GetEntries()
	if len(act_entries) != 1 {
		t.Fail()
	}
	if act_entries[0] != &et {
		t.Fail()
	}
}

func TestGetMultipleEntry(t *testing.T) {
	et1 := Entry{42, nil}
	et2 := Entry{43, nil}
	rc := Record{entryArray: []*Entry{&et1, &et2}}

	act_entries := rc.GetEntries()
	if len(act_entries) != 2 {
		t.Fail()
	}
	if act_entries[0] != &et1 || act_entries[1] != &et2 {
		t.Fail()
	}
}
