import { Participant } from '@/api/types';

// overviews

export interface TeamOverview {
  aggs?: CompiledTeamAggregations,
}

export interface PlayerOverview {
  tier?: string;
  aggs?: CompiledPlayerAggregations;
}

export interface PlayerOverviews { [name: string]: PlayerOverview }

export interface Overviews {
  team: TeamOverview,
  players: PlayerOverviews,
}

// aggregations

export interface PlayerAggregations {
  // average, etc
  vision: number[],
  cs: number[],
  jungle: {
    friendly: number[],
    enemy: number[],
  },
  dealt: number[],
  taken: number[],
  gold: number[],

  // take top 3
  champs: number[],

  // take top 1
  lanes: string[],
  roles: string[],
}

export function newPlayerAggregation(): PlayerAggregations {
  return {
    vision: [],
    cs: [],
    jungle: {
      friendly: [],
      enemy: [],
    },
    dealt: [],
    taken: [],
    gold: [],
    champs: [],
    lanes: [],
    roles: [],
  };
}

export function updatePlayerAggregation(agg: PlayerAggregations, part: Participant) {
  agg.vision.push(part.stats.visionScore);

  // minions + ungle
  agg.cs.push(part.stats.totalMinionsKilled);
  agg.cs.push(part.stats.neutralMinionsKilled);

  agg.jungle.friendly.push(part.stats.neutralMinionsKilledTeamJungle);
  agg.jungle.enemy.push(part.stats.neutralMinionsKilledEnemyJungle);

  agg.dealt.push(part.stats.totalDamageDealtToChampions);
  agg.taken.push(part.stats.totalDamageTaken);
  agg.gold.push(part.stats.goldEarned);
  agg.champs.push(part.championId);
  agg.lanes.push(part.timeline.lane);
  agg.roles.push(part.timeline.role);
}

export interface CompiledPlayerAggregations {
  avg: {
    // rounding turns numbers into... strings
    vision: string,
    cs: string,
    dealt: string,
    taken: string,
    gold: string,
    jungle: {
      friendly: string,
      enemy: string,
    },
  }
  favourite: {
    champs: number[],
    lane: string,
    role: string,
  },
}

export function getAvg(arr: number[]): string {
  const total = arr.reduce((acc, c) => acc + c, 0);
  return (total / arr.length).toFixed(2);
}

// this function drives me nuts. ts-ignore it all. TODO: use my brain
export function getTop(arr: (string | number)[], count?: number): (string | number) | (string | number)[] {
  // @ts-ignore
  const counts = arr.reduce((m, v) => {
    // @ts-ignore
    m[v] = (m[v] || 0) + 1; // eslint-disable-line no-param-reassign
    return m;
  }, {});
  // @ts-ignore
  const sorted = Object.keys(counts).sort((a, b) => counts[b] - counts[a]);
  if (!count) return sorted[0];
  return sorted.slice(0, count);
}

// TODO: more aggregations
export function compilePlayerAggregations(agg: PlayerAggregations): CompiledPlayerAggregations {
  return {
    avg: {
      vision: getAvg(agg.vision),
      cs: getAvg(agg.cs),
      dealt: getAvg(agg.dealt),
      taken: getAvg(agg.taken),
      gold: getAvg(agg.gold),
      jungle: {
        friendly: getAvg(agg.jungle.friendly),
        enemy: getAvg(agg.jungle.enemy),
      },
    },
    favourite: {
      champs: getTop(agg.champs, 5) as number[],
      lane: getTop(agg.lanes) as string,
      role: getTop(agg.roles) as string,
    },
  };
}

export interface TeamAggregations {
  matchTime: number[],
  wins: boolean[],

  towers: number[],
  dragons: number[],
  barons: number[],
}

export interface CompiledTeamAggregations {
  games: number,
  winRate: string,
  avg: {
    matchTime: string,
    towers: string,
    dragons: string,
    barons: string,
  },
}

export function compileTeamAggregations(agg: TeamAggregations): CompiledTeamAggregations {
  const wins = agg.wins.filter(x => x).length;

  const avgGameTime = parseInt(getAvg(agg.matchTime), 10);
  return {
    games: agg.wins.length,
    winRate: (wins / agg.wins.length).toFixed(2),
    avg: {
      matchTime: `${Math.floor(avgGameTime / 60)}:${Math.floor(avgGameTime % 60)}`,
      towers: getAvg(agg.towers),
      dragons: getAvg(agg.dragons),
      barons: getAvg(agg.barons),
    },
  };
}
