// Journald: in-memory, efficient JSON storage
//
// Debug helper functions

package journald

import (
	"fmt"
	"os"
)

// Dump a log structure to a graphviz representation
func Dump_to_graphviz(log Log, filename string) {
	gviz, err := os.Create(filename)
	defer gviz.Close()
	if err != nil {
		panic(err)
	}

	// Graph start
	fmt.Fprintf(gviz, "digraph Log {\n")

	// Simply go over existing records
	for _, rec := range log.recmap {
		fmt.Fprintf(gviz, "%p -> %p [label = %s];\n", rec, rec.next_field, rec.Field)
		if rec.entryArray != nil {
			for _, entry := range rec.entryArray {
				fmt.Fprintf(gviz, "%p -> %p;\n", rec, entry)
			}
		}
		if rec.entry != nil {
			fmt.Fprintf(gviz, "%p -> %p;\n", rec, rec.entry)
		}
	}

	// Also go over entries for back pointer

	// Finally, go over fields for back pointers
	fmt.Fprintf(gviz, "}\n")
}
