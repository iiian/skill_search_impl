package db

type MultiSkillIter struct {
	// map from skill to candidate ids, sorted alphanumerically
	by_skill map[string][]string
	// map from skill to index ptr into `by_skill`
	ptrs map[string]int
	// map from candidate id to score, mutates over life of iterator
	scores map[string]int
}

func NewIter(by_skill map[string][]string) MultiSkillIter {
	ptrs := make(map[string]int)
	scores := make(map[string]int)
	for skill := range by_skill {
		ptrs[skill] = 0
		first := by_skill[skill][0]
		scores[first] += 1
	}
	return MultiSkillIter{by_skill, ptrs, scores}
}

// Gets the Next sort-order candidate id
func (iter *MultiSkillIter) Next() string {
	var best string
	var assigned = false
	for c_id := range iter.scores {
		// init best
		if !assigned {
			best = c_id
			assigned = true
			continue
		}

		// there's a better best
		if c_id < best {
			best = c_id
		}
	}

	return best
}

// Does the iterator have more candidates to yield?
func (iter *MultiSkillIter) HasNext() bool {
	var has_next bool = false
	for skill := range iter.by_skill {
		ptr := iter.ptrs[skill]
		max := len(iter.by_skill[skill])
		if ptr < max {
			has_next = true
		}
	}
	return has_next
}

// After draining a candidate from the next
// this must be called as a clean up step to correctly maintain state of
// the iterator. This "pushes forward" the cursors into each skill index
func (iter *MultiSkillIter) PushPast(id string) {
	for skill := range iter.ptrs {
		haystack := iter.by_skill[skill]
		curs := iter.ptrs[skill]
		if haystack[curs] == id {
			iter.ptrs[skill] += 1
			curs = iter.ptrs[skill]
			if curs == len(haystack) {
				delete(iter.ptrs, skill)
				delete(iter.by_skill, skill)
				continue
			}
			next_id := haystack[curs]
			// https://staticcheck.io/docs/checks#S1036 :: golang maps act like python defaultdicts, so if the int DNE here, it'll init to 0
			iter.scores[next_id] += 1
		}
	}
	delete(iter.scores, id)
}

func (iter *MultiSkillIter) GetScore(id string) (int, bool) {
	score, ok := iter.scores[id]

	return score, ok
}
