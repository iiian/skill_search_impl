import { Candidate, IDatabase, ISearchCache, CacheRecord } from './api';
import { MultiSkillIterator } from './multi-skill-iterator';

export class Database implements IDatabase {
  /// map from skill to candidate id, sorted alphanumerically
  protected by_skill: { [key: string]: string[] };
  /// map from candidate id to candidate
  protected by_id: { [id: string]: Candidate };
  /// map from serialized skill searches to candidates that matched for that search
  protected search_cache: ISearchCache;

  constructor(search_cache: ISearchCache) {
    this.search_cache = search_cache;
    this.by_skill = {};
    this.by_id = {};
  }

  /// Add a new candidate to the database
  add(...candidates: Candidate[]): void {
    for (const candidate of candidates) {
      const { id, skills } = candidate;
      for (const skill of skills) {
        // note: prefix trie would have better performance here
        this.by_skill[skill] = sortedInsert(this.by_skill[skill] ?? [], id);
      }
      this.by_id[id] = candidate;
      this.search_cache.update(candidate);
    }
  }

  /**
   * Get a random best candidate from the db
   * 
   * @param skills the skills to search against
   * 
   * @returns a candidate, if one exists with a non-zero match score, otherwise null
   */
  search(...skills: string[]): Candidate | null {
    skills.sort();
    let cache_hit = this.search_cache.get(skills);
    if (!cache_hit) {
      cache_hit = this._handleSearch(skills);
      this.search_cache.put(skills, cache_hit);
    }
    if (cache_hit.score  == 0) {
      return null;
    }
    const index = Math.floor(Math.random() * cache_hit.ids.length);
    const id = cache_hit.ids[index];
    const candidate = this.by_id[id];
    return candidate;
  }
  /**
   * slice the db for the relevant subset of present skills
   * create an MultiSkillIterator out of them
   * take max score over the iterator
   * 
   * @param skills the skills by which to search the candidates
   * @returns [candidate ids, their common score]
   */
  protected _handleSearch(skills: string[]): CacheRecord {
    // slice of skills to iterate a frontier over
    const slice: { [skill: string]: string[] } = {};

    const target_skills = skills.filter(skill => !!this.by_skill[skill]);
    for (const skill of target_skills) {
      slice[skill] = this.by_skill[skill];
    }
    const iterator = new MultiSkillIterator(slice).as();
    let best_ids: string[] = [];
    let best_score: number = 0;
    for (const [id, score] of iterator) {
      if (best_score == score) {
        best_ids.push(id);
      } else if (best_score < score) {
        best_score = score;
        best_ids = [id];
      }
    }
    
    return {
      ids: best_ids,
      score: best_score,
    };
  }
}

/// Modifies haystack by inserting needle in an alphabetically sorted position
function sortedInsert(haystack: string[], needle: string) {
  let at_index = 0;
  if (haystack.length) {
    while (at_index < haystack.length) {
      const next = haystack[at_index];
      if (needle < next) {
        break
      }
      at_index++;
    }
  }
  haystack.splice(at_index, 0, needle);
  return haystack;
}