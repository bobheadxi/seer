
export enum Region {
  NA = "NA1",
}

export interface Team {
  Region: Region;
  Members: Member[];
}

export interface Match {
  seasonId: number;
  queueId: number;
  gameId: number;
  participantIdentities: ParticipantIdentity[];
  gameVersion: string;
  platformId: string;
  gameMode: string;
  mapId: number;
  gameType: string;
  teams: Team[];
  participants: Participant[];
  gameDuration: number;
  gameCreation: number;
}

// ------------------------------------

export interface Member {
  name: string;
  summonerLevel: number;
  profileIconId: number;
  accountId: string;
  id: string;
}

export interface Player {
  accountId: string;
}

export interface ParticipantIdentity {
  player: Player;
  participantId: number;
}

export interface Ban {
  pickTurn: number;
  championId: number;
}

export interface Team {
  firstDragon: boolean;
  bans: Ban[];
  firstInhibitor: boolean;
  win: string;
  firstRiftHerald: boolean;
  firstBaron: boolean;
  baronKills: number;
  riftHeraldKills: number;
  firstBlood: boolean;
  teamId: number;
  firstTower: boolean;
  vilemawKills: number;
  inhibitorKills: number;
  towerKills: number;
  dominionVictoryScore: number;
  dragonKills: number;
}

export interface Stats {
  neutralMinionsKilledTeamJungle: number;
  visionScore: number;
  magicDamageDealtToChampions: number;
  largestMultiKill: number;
  totalTimeCrowdControlDealt: number;
  longestTimeSpentLiving: number;
  perk1Var1: number;
  perk1Var3: number;
  perk1Var2: number;
  tripleKills: number;
  perk5: number;
  perk4: number;
  playerScore9: number;
  playerScore8: number;
  kills: number;
  playerScore1: number;
  playerScore0: number;
  playerScore3: number;
  playerScore2: number;
  playerScore5: number;
  playerScore4: number;
  playerScore7: number;
  playerScore6: number;
  perk5Var1: number;
  perk5Var3: number;
  perk5Var2: number;
  totalScoreRank: number;
  neutralMinionsKilled: number;
  statPerk1: number;
  statPerk0: number;
  damageDealtToTurrets: number;
  physicalDamageDealtToChampions: number;
  damageDealtToObjectives: number;
  perk2Var2: number;
  perk2Var3: number;
  totalUnitsHealed: number;
  perk2Var1: number;
  perk4Var1: number;
  totalDamageTaken: number;
  perk4Var3: number;
  wardsKilled: number;
  largestCriticalStrike: number;
  largestKillingSpree: number;
  quadraKills: number;
  magicDamageDealt: number;
  firstBloodAssist: boolean;
  item2: number;
  item3: number;
  item0: number;
  item1: number;
  item6: number;
  item4: number;
  item5: number;
  perk1: number;
  perk0: number;
  perk3: number;
  perk2: number;
  perk3Var3: number;
  perk3Var2: number;
  perk3Var1: number;
  damageSelfMitigated: number;
  magicalDamageTaken: number;
  perk0Var2: number;
  firstInhibitorKill: boolean;
  trueDamageTaken: number;
  assists: number;
  perk4Var2: number;
  goldSpent: number;
  trueDamageDealt: number;
  participantId: number;
  physicalDamageDealt: number;
  sightWardsBoughtInGame: number;
  totalDamageDealtToChampions: number;
  physicalDamageTaken: number;
  totalPlayerScore: number;
  win: boolean;
  objectivePlayerScore: number;
  totalDamageDealt: number;
  neutralMinionsKilledEnemyJungle: number;
  deaths: number;
  wardsPlaced: number;
  perkPrimaryStyle: number;
  perkSubStyle: number;
  turretKills: number;
  firstBloodKill: boolean;
  trueDamageDealtToChampions: number;
  goldEarned: number;
  killingSprees: number;
  unrealKills: number;
  firstTowerAssist: boolean;
  firstTowerKill: boolean;
  champLevel: number;
  doubleKills: number;
  inhibitorKills: number;
  firstInhibitorAssist: boolean;
  perk0Var1: number;
  combatPlayerScore: number;
  perk0Var3: number;
  visionWardsBoughtInGame: number;
  pentaKills: number;
  totalHeal: number;
  totalMinionsKilled: number;
  timeCCingOthers: number;
  statPerk2: number;
}

export interface Deltas {
  [range: string]: number;
}

export interface Timeline {
  participantId: number;
  lane: string;
  role: string;
  csDiffPerMinDeltas: Deltas;
  goldPerMinDeltas: Deltas;
  xpDiffPerMinDeltas: Deltas;
  creepsPerMinDeltas: Deltas;
  xpPerMinDeltas: Deltas;
  damageTakenDiffPerMinDeltas: Deltas;
  damageTakenPerMinDeltas: Deltas;
}

export interface Participant {
  stats: Stats;
  spell1Id: number;
  participantId: number;
  spell2Id: number;
  teamId: number;
  timeline: Timeline;
  championId: number;
  highestAchievedSeasonTier: string;
}

