import { CacheRecord, Candidate, ISearchCache } from './api';

const SEPARATOR = String.fromCharCode(178);

export class SearchCache implements ISearchCache {
  protected constructor(protected cache: { [ser: string]: CacheRecord}) {}

  get(skills: string[]): CacheRecord | null {
    return this.cache[SearchCache.serialize(skills)] ?? null;
  }
  put(skills: string[], cache_record: CacheRecord): void {
    const key = SearchCache.serialize(skills);
    this.cache[key] = cache_record;
  }
  update(candidate: Candidate): void {
    for (const key in this.cache) {
      const skills = SearchCache.deserialize(key);
      let score = 0;
      for (const skill of candidate.skills) {
        if (skills.includes(skill)) {
          score++;
        }
      }

      const cache_hit = this.cache[key];
      if (cache_hit.score === score) {
        this.cache[key].ids.push(candidate.id);
      } else if (cache_hit.score < score) {
        this.cache[key] = {
          ids: [candidate.id],
          score,
        }
      }
    }
  }
  static new(): SearchCache {
    return new SearchCache({});
  }
  protected static serialize(skills: string[]): string {
    return skills.join(SEPARATOR);
  }
  protected static deserialize(key: string): string[] {
    return key.split(SEPARATOR);
  }
}