package journald

import "testing"

func TestAddRetrieve(t *testing.T) {
	lg := NewLog()
	rc1 := map[string]string{
		"val1": "test1",
		"val2": "test2",
		"val3": "test3",
		"val4": "test3",
		"val5": "test4",
		"val6": "test5",
	}
	lg.Append(42, rc1)
	rc2 := map[string]string{
		"val1":  "test1",
		"val2":  "test2",
		"val7":  "test3",
		"val4":  "test3",
		"val5":  "test8",
		"val16": "test5",
	}
	lg.Append(43, rc2)

	rc3 := map[string]string{
		"val5": "test12",
		"val42": "This is a very long string that should get compressed eventually. Heck yeah, we want to compress long strings that can for example be large queries sent to the database. To do so, we use the snappy library, that should provide a reasonably fast implementation of a compression for relatively good performance. Another level of compression is of course the record de-duplication offered by the way this database is modelled. Interestingly, the code in the debug.go file generate a wonderfull graphviz-compatible dot file, very usefull for visual debugging. I have no clue how this approach will scale, though, since graphs are notoriously bad when many elements are present — they are called spaghetti chart for a reason :)",
	}
	lg.Append(44, rc3)

	ets := lg.find_by_equal("val3", "test3")
	if ets[0].Timestamp != 42 {
		t.Error("Val3 was not correctly located")
	}

	ets = lg.find_by_equal("val1", "test1")
	if len(ets) != 2 || ets[0].Timestamp + ets[1].Timestamp != 85 {
		t.Error("Val1 not correctly located")
	}

	ets = lg.find_by_equal("val7", "nonsense")
	if len(ets) != 0 {
		t.Error("Some entries were created")
	}

	rcs := lg.find_by_field("val5")
	// Flaky test: order does not matter
	if len(rcs) != 3 || rcs[0].Value != "test4" || rcs[1].Value != "test8" || rcs[2].Value != "test12" {
		t.Error("Some associated values of val5 missing")
	}

	rcs = lg.find_by_field("val38")
	if len(rcs) != 0 {
		t.Error("Entries were created for val5")
	}
}
