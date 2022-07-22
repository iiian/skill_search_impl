export type Candidate = {
  id: string;
  skills: string[];
}

export function candidate(id: string, ...skills: string[]): Candidate {
  return { id, skills };
}

export interface IDatabase {
  add(...candidates: Candidate[]): void;
  search(...skills: string[]): Candidate|null;
}

export type CacheRecord = {
  score: number;
  ids: string[];
}

export interface ISearchCache {
  get(skills: string[]): CacheRecord|null;
  put(skills: string[], cache_record: CacheRecord): void;
  update(candidate: Candidate): void;
}
