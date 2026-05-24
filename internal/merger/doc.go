// Package merger provides utilities for combining the key sets of multiple
// parsed .env files into a single unified reference set.
//
// When envlens is given several reference files (e.g. .env.base and
// .env.defaults), Merge produces the union of all their keys. Keys that do
// not appear in every source are marked as "partial" so downstream reporters
// can emit warnings rather than hard errors.
package merger
