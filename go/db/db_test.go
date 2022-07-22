package db

import (
	"testing"

	slices "golang.org/x/exp/slices"
)

func TestBestMatch1(t *testing.T) {
	sc := make(searchCache)
	db := New(sc)
	db.Add(
		&Candidate{
			ID:     "a",
			Skills: []string{"nestjs", "golang"},
		},
		&Candidate{
			ID:     "b",
			Skills: []string{"golang", "k8s", "nestjs"},
		},
		&Candidate{
			ID:     "c",
			Skills: []string{"ruby", "k8s"},
		},
	)
	want_any_of := []string{"a", "b"}
	got, ok := db.Search([]string{"nestjs", "golang", "k8s"})
	if !ok {
		t.Fatalf("Couldn't find a match for nestjs, but wanted candidate 'a' or 'b")
	}
	if !slices.Contains(want_any_of, got.ID) {
		t.Fatalf("Expected any of %s, but got %s", want_any_of, got.ID)
	}
	want := "c"
	got, ok = db.Search([]string{"ruby", "k8s"})
	if !ok {
		t.Fatalf("Couldn't find a match for ruby|k8s, but wanted candidate 'c'")
	}
	if want != got.ID {
		t.Fatalf("Expected c but got %s", got.ID)
	}
}

func TestBestMatch2(t *testing.T) {
	sc := make(searchCache)
	db := New(sc)
	db.Add(
		&Candidate{
			ID:     "a",
			Skills: []string{"golang", "k8s", "nestjs"},
		},
		&Candidate{
			ID:     "b",
			Skills: []string{"nestjs", "golang"},
		},
		&Candidate{
			ID:     "c",
			Skills: []string{"ruby", "k8s"},
		},
	)
	want := "c"
	got, ok := db.Search([]string{"ruby", "k8s"})
	if !ok {
		t.Fatalf("Couldn't find a match for ruby, k8s, but wanted candidate 'c'")
	}
	if got.ID != want {
		t.Fatalf("Expected %s but got %s", want, got.ID)
	}
}

func TestBestMatch3(t *testing.T) {
	sc := make(searchCache)
	db := New(sc)
	db.Add(
		&Candidate{
			ID:     "a",
			Skills: []string{"golang", "k8s", "nestjs"},
		},
		&Candidate{
			ID:     "b",
			Skills: []string{"nestjs", "golang"},
		},
		&Candidate{
			ID:     "c",
			Skills: []string{"ruby", "k8s"},
		},
	)
	want := "a"
	got, ok := db.Search([]string{"golang", "k8s"})
	if !ok {
		t.Fatalf("Couldn't find a match for golang, k8s, but wanted candidate 'a'")
	}
	if got.ID != want {
		t.Fatalf("Expected %s but got %s", want, got.ID)
	}
}

func TestBadMatch(t *testing.T) {
	sc := make(searchCache)
	db := New(sc)
	db.Add(
		&Candidate{
			ID:     "a",
			Skills: []string{"golang", "k8s", "nestjs"},
		},
		&Candidate{
			ID:     "b",
			Skills: []string{"nestjs", "golang"},
		},
		&Candidate{
			ID:     "c",
			Skills: []string{"ruby", "k8s"},
		},
	)
	best, ok := db.Search([]string{"java"})
	if ok {
		t.Fatalf("Expected no match but got %s", best.ID)
	}
}

type mockSearchCache struct {
	searchCache
	was_get_called          bool
	was_update_cache_called bool
}

func (msc *mockSearchCache) Get(skills []string) (*CacheRecord, bool) {
	msc.was_get_called = true
	return msc.searchCache.Get(skills)
}

func (msc *mockSearchCache) HandleCacheUpdate(c *Candidate) {
	msc.was_update_cache_called = true
	msc.searchCache.HandleCacheUpdate(c)
}

// It should leverage a search cache when applicable
func TestUseSearchCache(t *testing.T) {
	sc := mockSearchCache{
		searchCache: make(searchCache),
	}
	db := New(&sc)
	db.Add(
		&Candidate{
			ID:     "a",
			Skills: []string{"golang", "k8s", "nestjs"},
		},
		&Candidate{
			ID:     "b",
			Skills: []string{"nestjs", "golang"},
		},
		&Candidate{
			ID:     "c",
			Skills: []string{"ruby", "k8s"},
		},
	)
	db.Search([]string{"kotlin", "k8s"})
	// duplicate call should trigger a cache hit
	res, ok := db.Search([]string{"kotlin", "k8s"})
	if !ok || !sc.was_get_called {
		t.Fatalf("Expected cache to be hit")
	}
	want_any_of := []string{"a", "c"}
	if !slices.Contains(want_any_of, res.ID) {
		t.Fatalf("Expected cache hit to be used")
	}
}

// It should update the search cache when adding new candidates
func TestUpdateSearchCache(t *testing.T) {
	sc := mockSearchCache{
		searchCache: make(searchCache),
	}
	db := New(&sc)
	db.Add(
		&Candidate{
			ID:     "a",
			Skills: []string{"golang", "k8s", "nestjs"},
		},
		&Candidate{
			ID:     "b",
			Skills: []string{"nestjs", "golang"},
		},
		&Candidate{
			ID:     "c",
			Skills: []string{"ruby", "k8s"},
		},
	)
	db.Search([]string{"kotlin", "k8s"})
	db.Add(
		&Candidate{
			ID:     "d",
			Skills: []string{"golang", "k8s", "kotlin"},
		},
	)
	db.Add(
		&Candidate{
			ID:     "e",
			Skills: []string{"golang", "k8s", "kotlin"},
		},
	)
	res, ok := db.Search([]string{"kotlin", "k8s"})
	if !ok || !sc.was_get_called {
		t.Fatalf("Expected cache to be hit")
	}
	want_any_of := []string{"d", "e"}
	if !slices.Contains(want_any_of, res.ID) {
		t.Fatalf("Expected cache hit to be used")
	}
}
