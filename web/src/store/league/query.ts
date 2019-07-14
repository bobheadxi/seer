import * as LeagueTypes from './types';

interface KeyedObject { key: string }

export function findByKey<T extends KeyedObject>(
  src: { [name: string]: T },
  match: string,
): T | undefined {
  const found = Object.keys(src)
    .find(name => src[name].key === match);
  if (!found) return undefined;
  return src[found];
}
