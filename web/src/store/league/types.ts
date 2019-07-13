

// Champions
// http://ddragon.leagueoflegends.com/cdn/9.13.1/data/en_US/champion.json

export interface ChampData {
  key: string;
  name: string;
  title: string;
  tags: [string];
  image: Image;
}

// Items
// http://ddragon.leagueoflegends.com/cdn/9.13.1/data/en_US/item.json

export interface ItemData {
  name: string;
  description: string;
  plaintext: string;
  image: Image;
}

// Runes reforged
// http://ddragon.leagueoflegends.com/cdn/9.13.1/data/en_US/runesReforged.json

export interface RunesData {
  id: number;
  key: string;
  icon: string;
  name: string;
  slots: RuneSlot[];
}

export interface RuneSlot {
  runes: Rune[];
}

export interface Rune {
  id: number;
  key: string;
  icon: string;
  name: string;
  shortDesc: string;
  longDesc: string;
}

// Summoner spells data
// http://ddragon.leagueoflegends.com/cdn/9.13.1/data/en_US/summoner.json

export interface SpellData {
  id: string;
  name: string;
  description: string;
  tooltip: string;
  maxrank: number;
  cooldown: number[];
  cooldownBurn: string;
  cost: number[];
  costBurn: string;
  datavalues: {};
  effect: number[][];
  effectBurn: string[];
  vars: any[];
  key: string;
  summonerLevel: number;
  modes: string[];
  costType: string;
  maxammo: string;
  range: number[];
  rangeBurn: string;
  image: Image;
  resource: string;
}

// Primitives

export interface Image {
  full: string;
  sprite: string;
  group: string;
  x: number;
  y: number;
  w: number;
  h: number;
}
