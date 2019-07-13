import * as LeagueTypes from './types';

export function findChamp(
  champs: { [name: string]: LeagueTypes.ChampData },
  id: string,
): LeagueTypes.ChampData | undefined {
  const found = Object.keys(champs)
    .find(name => champs[name].key === id);
  if (!found) return undefined;
  return champs[found];
}

export function findSpell(
  spells: { [name: string]: LeagueTypes.SpellData },
  id: string,
): LeagueTypes.SpellData | undefined {
  const found = Object.keys(spells)
    .find(name => spells[name].key === id);
  if (!found) return undefined;
  return spells[found];
}
