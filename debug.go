// Journald: in-memory, efficient JSON storage
//
// Debug helper functions

package journald

import (
	"fmt"
	"io"
	"os"
)

func (rec *Record) Dump(gviz io.Writer) {
	fmt.Fprintf(gviz, "\"%p\" [label = %s];\n", rec, rec.Field)
	fmt.Fprintf(gviz, "\"%p\" -> \"%p\";\n", rec, rec.next_field)
	if rec.entryArray != nil {
		fmt.Fprintf(gviz, "\"%p\" -> {\n", rec)
		for _, entry := range rec.entryArray {
			fmt.Fprintf(gviz, "\"%p\"\n", entry)
		}
		fmt.Fprintf(gviz, "} [color=\"blue\"];\n")
	}
	if rec.entry != nil {
		fmt.Fprintf(gviz, "\"%p\" -> \"%p\" [color=\"blue\"];\n", rec, rec.entry)
	}
}

func (entry *Entry) Dump(gviz io.Writer) {
	fmt.Fprintf(gviz, "\"%p\" [label = \"%x\", shape=\"square\"];\n", entry, entry.Timestamp)
	for _, rc := range entry.Records {
		fmt.Fprintf(gviz, "\"%p\" -> \"%p\" [color=\"red\"];\n", entry, rc)
	}
}

// Dump a log structure to a graphviz representation
func (log Log) Dump(filename string) {
	gviz, err := os.Create(filename)
	defer gviz.Close()
	if err != nil {
		panic(err)
	}

	// Graph start
	fmt.Fprintf(gviz, "digraph Log {\n")

	// Simply go over existing records
	for _, rec := range log.recmap {
		rec.Dump(gviz)
	}

	// Also go over entries for back pointer
	for i := range log.currentEntry {
		if log.currentEntry[i].Timestamp == 0 {
			continue
		}
		log.currentEntry[i].Dump(gviz)
	}

	for _, bl := range log.backlogEntryArrays {
		for i := range bl {
			log.currentEntry[i].Dump(gviz)
		}
	}

	// Finally, go over fields for back pointers
	for field, rec := range log.fieldmap {
		fmt.Fprintf(gviz, "\"%s\" [shape=\"diamond\"];\n", field)
		fmt.Fprintf(gviz, "\"%s\" -> \"%p\";\n", field, rec)
	}

	fmt.Fprintf(gviz, "}\n")
}
