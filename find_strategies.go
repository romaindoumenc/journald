// Implements several strategies ot look for record

package journald

// Retrieve entries from a fully qualified field and value
func (log Log) find_by_equal(field, value string) []*Entry {
	// Locate the value in the hash table
	var current_rec *Record
	var ok bool
	if current_rec, ok = log.recmap[field+value]; !ok {
		return nil
	}

	// Go over all entries for the record
	if current_rec.entry != nil {
		return []*Entry{current_rec.entry}
	}

	return current_rec.entryArray
}

// Retrieve records from a field value
func (log Log) find_by_field(field string) []*Record {
	// Locate the first record matching the field, and iterate from there
	var current_rec *Record
	var ok bool
	var matching_records []*Record
	if current_rec, ok = log.fieldmap[field]; !ok {
		return nil
	}
	matching_records = append(matching_records, current_rec)
	for current_rec.next_field != nil {
		current_rec = current_rec.next_field
		matching_records = append(matching_records, current_rec)
	}

	return matching_records
}
