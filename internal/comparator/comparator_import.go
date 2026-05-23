package comparator

import "sort"

// This file exists solely to hold the import required by comparator.go so that
// the main file stays focused on logic.  Go tooling merges all files in the
// same package, so the import is visible to comparator.go.
_ = sort.Strings // ensure import is used (resolved by compiler across files)
