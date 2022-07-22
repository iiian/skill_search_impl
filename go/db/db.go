package db

import (
	"math/rand"
	"sort"
)

type DB struct {
	by_skill     map[string][]string   // map from skill to candidate id, sorted alphanumerically
	by_id        map[string]*Candidate // map from candidate id to candidate
	search_cache SearchCache           // map from serialized skill searches to candidates that matched for that search
}

func New(search_cache SearchCache) DB {
	return DB{
		by_skill:     make(map[string][]string),
		by_id:        make(map[string]*Candidate),
		search_cache: search_cache,
	}
}

// Add a new candidate to the database
func (db *DB) Add(cs ...*Candidate) {
	for _, c := range cs {
		skills := c.Skills
		for _, skill := range skills {
			db.by_skill[skill] = sortedInsert(db.by_skill[skill], c.ID)
		}
		db.by_id[c.ID] = c
		db.search_cache.HandleCacheUpdate(c)
	}
}

// Get a random best candidate from the db
func (db *DB) Search(skills []string) (*Candidate, bool) {
	sort.Strings(skills)
	cache_hit, ok := db.search_cache.Get(skills)
	if !ok {
		// slice the db for the relevant subset of present skills
		// create an iterator out of them
		best_ids, best_score := db.handleSearch(skills)
		cache_hit = &CacheRecord{
			IDs:   best_ids,
			score: best_score,
		}
		db.search_cache.Put(skills, cache_hit)
	}
	if cache_hit.score == 0 {
		return nil, false
	}
	rand_i := rand.Intn(len(cache_hit.IDs))
	c_id := cache_hit.IDs[rand_i]
	candidate := db.by_id[c_id]
	return candidate, true
}

func (db *DB) handleSearch(skills []string) ([]string, int) {
	var slice map[string][]string = make(map[string][]string)

	for _, skill := range skills {
		if _, ok := db.by_skill[skill]; ok {
			slice[skill] = db.by_skill[skill]
		}
	}

	iter := NewIter(slice)

	best_ids := []string{}
	best_score := 0

	for iter.HasNext() {
		c_id := iter.Next()
		c_score, _ := iter.GetScore(c_id)
		if best_score == c_score {
			best_ids = append(best_ids, c_id)
		} else if best_score < c_score {
			best_ids = []string{c_id}
			best_score = c_score
		}
		iter.PushPast(c_id)
	}
	return best_ids, best_score
}
