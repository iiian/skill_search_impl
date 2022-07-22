package db

import "golang.org/x/exp/slices"

type CacheRecord struct {
	IDs   []string
	score int
}

const SEP byte = byte(178)

// Assumes skills is sorted alphabetically
func serializeSkillset(skills []string) string {
	var serialization = []byte{}
	for _, skill := range skills {
		bytes := []byte(skill)
		serialization = append(serialization, bytes...)
		serialization = append(serialization, SEP)
	}
	return string(serialization)
}

func deserializeKey(key []byte) []string {
	var skills []string = []string{}
	bytes := []byte{}
	for _, next := range key {
		if next == SEP {
			skill := string(bytes[:])
			skills = append(skills, skill)
			bytes = []byte{}
		} else {
			bytes = append(bytes, next)
		}
	}
	return skills
}

type searchCache map[string]*CacheRecord

// Updates the db search cache based on the given candidate's skills and a given search cache key.
func (sc searchCache) HandleCacheUpdate(c *Candidate) {
	for key := range sc {
		// key is a compact serialization of skills
		skills := deserializeKey([]byte(key))

		// build score based on set overlap of candidate skills with key skills
		c_score := 0
		for _, c_skill := range c.Skills {
			if slices.Contains(skills, c_skill) {
				c_score += 1
			}
		}

		// update cache if applicable
		cache_hit := sc[key]
		if cache_hit.score == c_score {
			sc[key].IDs = append(sc[key].IDs, c.ID)
		} else if cache_hit.score < c_score {
			sc[key] = &CacheRecord{
				score: c_score,
				IDs:   []string{c.ID},
			}
		}
	}
}

func (sc searchCache) Get(skills []string) (*CacheRecord, bool) {
	key := serializeSkillset(skills)
	rec, ok := sc[key]
	return rec, ok
}

func (sc searchCache) Put(skills []string, rec *CacheRecord) {
	key := serializeSkillset(skills)
	sc[key] = rec
}
