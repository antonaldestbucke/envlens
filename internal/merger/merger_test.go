package merger_test

import (
	"sort"
	"testing"

	"github.com/user/envlens/internal/loader"
	"github.com/user/envlens/internal/merger"
)

func makeEnv(name string, keys ...string) loader.Env {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = "value"
	}
	return loader.Env{Name: name, Keys: m}
}

func sortedKeys(r merger.Result) []string {
	out := make([]string, len(r.Keys))
	copy(out, r.Keys)
	sort.Strings(out)
	return out
}

func TestMerge_SingleSource(t *testing.T) {
	envs := []loader.Env{makeEnv("base", "A", "B", "C")}
	r := merger.Merge(envs)

	if got := sortedKeys(r); len(got) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(got))
	}
	if len(r.Partial) != 0 {
		t.Errorf("expected no partial keys with single source, got %v", r.Partial)
	}
}

func TestMerge_UnionOfKeys(t *testing.T) {
	envs := []loader.Env{
		makeEnv("a", "X", "Y"),
		makeEnv("b", "Y", "Z"),
	}
	r := merger.Merge(envs)
	got := sortedKeys(r)

	want := []string{"X", "Y", "Z"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("key[%d]: want %q got %q", i, k, got[i])
		}
	}
}

func TestMerge_PartialKeys(t *testing.T) {
	envs := []loader.Env{
		makeEnv("a", "SHARED", "ONLY_A"),
		makeEnv("b", "SHARED", "ONLY_B"),
	}
	r := merger.Merge(envs)

	if !r.IsPartial("ONLY_A") {
		t.Error("expected ONLY_A to be partial")
	}
	if !r.IsPartial("ONLY_B") {
		t.Error("expected ONLY_B to be partial")
	}
	if r.IsPartial("SHARED") {
		t.Error("expected SHARED to not be partial")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	r := merger.Merge([]loader.Env{})
	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestMerge_DeduplicatesKeys(t *testing.T) {
	envs := []loader.Env{
		makeEnv("a", "DUP", "UNIQUE"),
		makeEnv("b", "DUP"),
	}
	r := merger.Merge(envs)

	count := 0
	for _, k := range r.Keys {
		if k == "DUP" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected DUP to appear once, got %d", count)
	}
}
