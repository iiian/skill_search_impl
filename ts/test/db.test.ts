import { Database } from '../src/db/db';
import { SearchCache } from '../src/db/search-cache';
import { candidate } from '../src/db/api';
describe('Database', () => { 
  let spy: jest.SpyInstance;
  let cache: SearchCache;
  let db: Database;
  beforeEach(() => {
    cache = SearchCache.new();
    db = new Database(cache);
  });
  describe('add', () => {
    it('should store candidates', () => {
      db.add(
        candidate('a', 'nestjs', 'golang'),
        candidate('b', 'golang', 'k8s', 'nestjs'),
        candidate('c', 'ruby', 'k8s'),
      );
      expect(db.search('golang', 'k8s', 'nestjs')!.id).toBe('b');
      expect(db.search('nestjs', 'golang')!.id).toMatch(/a|b/);
      expect(db.search('ruby', 'k8s')!.id).toBe('c');
    });
    describe('given a candidate with skills related to a cached search query', () => {
      beforeEach(() => {
        spy = jest.spyOn(cache, "update");
      });
      afterEach(() => {
        spy.mockRestore();
      })
      it('should update the search cache entry', () => {
        db.search('nestjs', 'golang'); // a cached query
        db.add(
          candidate('a', 'nestjs', 'golang'),
          candidate('b', 'golang', 'k8s', 'nestjs'),
          candidate('c', 'ruby', 'k8s'),
        );
        expect(spy).toHaveBeenCalledTimes(3);
      });
    });
  });

  describe('search', () => {
    describe('given a search request that is cached', () => {
      beforeEach(() => {
        spy = jest.spyOn(cache, "get");
      });
      afterEach(() => {
        spy.mockRestore();
      });
      it('should yield a candidate from the cached collection', () => {
        db.add(
          candidate('a', 'nestjs', 'golang'),
          candidate('b', 'golang', 'k8s', 'nestjs'),
          candidate('c', 'ruby', 'k8s'),
        );
        db.search('k8s');
        expect(spy).toHaveReturnedWith(null);
        db.search('k8s');
        expect(spy).not.toHaveLastReturnedWith(null);
      });
    });
    describe('otherwise', () => {
      it('should return a random candidate with the highest number of matching skills', () => {
        db.add(
          candidate('a', 'nestjs', 'golang'),
          candidate('b', 'golang', 'k8s', 'nestjs'),
          candidate('c', 'ruby', 'k8s'),
        );
        expect(db.search('ruby', 'sql', 'f-sharp')!.id).toBe('c');
      });
      it('should return null if none could be found', () => {
        db.add(
          candidate('a', 'nestjs', 'golang'),
          candidate('b', 'golang', 'k8s', 'nestjs'),
          candidate('c', 'ruby', 'k8s'),
        );
        expect(db.search('sql', 'haskell')).toBeFalsy();
      });
      beforeEach(() => {
        spy = jest.spyOn(cache, "put");
      });
      afterEach(() => {
        spy.mockRestore();
      });
      it('should update the cache with the result', () => {
        db.search('k8s', 'golang');
        expect(spy).toHaveBeenCalled();
      });
    });
  });
});