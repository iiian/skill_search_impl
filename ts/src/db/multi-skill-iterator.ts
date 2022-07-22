export class MultiSkillIterator {
  /// Map from skill to candidate ids
  protected slice: { [skill: string]: string[] };
  /// Map from skill to cursor into `by_skill`
  protected cursors: { [skill: string]: number };
  /// Map from candidate id to score, mutates over life of iterator
  protected scores: { [id: string]: number };
  
  constructor(slice: { [skill: string]: string[] }) {
    this.slice = slice;
    this.cursors = {};
    this.scores = {};
  }
  
  *as(): Iterable<[string, number]> {
    for (const skill in this.slice) {
      this.cursors[skill] = 0;
      const first: string = this.slice[skill][0];
      this.scores[first] = (this.scores[first] ?? 0) + 1;
    }

    while (true) {
      const frontier = Object.keys(this.scores);
      if (!frontier.length) {
        return;
      }
      let best: string = frontier[0];
      for (const id of frontier) {
        if (id < best) {
          best = id;
        }
      }
      yield [best, this.scores[best]];
      this._pushPast(best);
    }
  }

  protected _pushPast(id: string) {
    for (const skill in this.cursors) {
      const haystack = this.slice[skill];
      let cursor = this.cursors[skill];
      if (haystack[cursor] === id) {
        this.cursors[skill] = (++cursor);
        if (cursor == haystack.length) {
          delete this.cursors[skill];
          delete this.slice[skill];
          continue;
        }
        const next_id = haystack[cursor];
        this.scores[next_id] = (this.scores[next_id] ?? 0) + 1;
      }
    }
    delete this.scores[id];
  }
}