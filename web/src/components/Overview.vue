<template>
  <div class="overview">
    <h2>Overview</h2>
    <div>
      <div>
        {{ overviews.team }}
      </div>
      <div v-for="m in team.members" v-bind:key="m.id">
        <h3>
          {{ m.name }}
          ({{ overviews.players && overviews.players[m.name] && overviews.players[m.name].tier ? overviews.players[m.name].tier + ', ' : ''}}
          lv{{ m.summonerLevel }})
          <!-- TODO: regions? --->
          <a v-bind:href="'https://na.op.gg/summoner/userName='+m.name" target="_blank">
            <img width="16" height="16"
              src="https://lh3.googleusercontent.com/UdvXlkugn0bJcwiDkqHKG5IElodmv-oL4kHlNAklSA2sdlVWhojsZKaPE-qFPueiZg" />
          </a>
        </h3>
        <div v-if="overviews.players && overviews.players[m.name] && overviews.players[m.name].aggs">
          <div>
            <h5>Most played lane and role</h5>
            {{ overviews.players[m.name].aggs.favourite.lane }}
            ({{ overviews.players[m.name].aggs.favourite.role }})
          </div>
          <div>
            <h5>Most played champions</h5>
            <img
              v-for="c in overviews.players[m.name].aggs.favourite.champs"
              v-bind:key="'fav-'+m.name+'-'+c"
              v-bind:src="champIcon(c)" />
          </div>
          <div>
            <h5>Average stats</h5>
            {{ overviews.players[m.name].aggs.avg }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

import { Vue, Component, Prop } from 'vue-property-decorator';
import { State, Action, Getter } from 'vuex-class';

import { Namespace } from '../store';
import teams, { TeamGetters } from '../store/teams';
import { LeagueGetters, IDGetter } from '../store/league';
import {
  ItemData, ChampData, RunesData, SpellData,
} from '../store/league/types';
import {
  Match, Team, Participant, MatchTeam,
} from '../api/types';

const leagueSpace = { namespace: Namespace.LEAGUE };
const teamsSpace = { namespace: Namespace.TEAMS };

// TODO: move math stuff into a module

interface PlayerOverviews { [name: string]: PlayerOverview }

interface PlayerOverview {
  tier?: string;
  aggs?: CompiledPlayerAggregations;
}

interface PlayerAggregations {
  // average, etc
  vision: number[],
  cs: number[],
  dealt: number[],
  taken: number[],
  gold: number[],

  // take top 3
  champs: number[],

  // take top 1
  lanes: string[],
  roles: string[],
}

function newPlayerAggregation(): PlayerAggregations {
  return {
    vision: [],
    cs: [],
    dealt: [],
    taken: [],
    gold: [],
    champs: [],
    lanes: [],
    roles: [],
  };
}

function updatePlayerAggregation(agg: PlayerAggregations, part: Participant) {
  agg.vision.push(part.stats.visionScore);
  agg.cs.push(part.stats.totalMinionsKilled);
  agg.dealt.push(part.stats.totalDamageDealtToChampions);
  agg.taken.push(part.stats.totalDamageTaken);
  agg.gold.push(part.stats.goldEarned);
  agg.champs.push(part.championId);
  agg.lanes.push(part.timeline.lane);
  agg.roles.push(part.timeline.role);
}

interface CompiledPlayerAggregations {
  avg: {
    // rounding turns numbers into... strings
    vision: string,
    cs: string,
    dealt: string,
    taken: string,
    gold: string,
  }
  favourite: {
    champs: number[],
    lane: string,
    role: string,
  }
}

function getAvg(arr: number[]): string {
  const total = arr.reduce((acc, c) => acc + c, 0);
  return (total / arr.length).toFixed(2);
}

// this function drives me nuts. ts-ignore it all. TODO: use my brain
function getTop(arr: (string | number)[], count?: number): (string | number) | (string | number)[] {
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
function compilePlayerAggregations(agg: PlayerAggregations): CompiledPlayerAggregations {
  return {
    avg: {
      vision: getAvg(agg.vision),
      cs: getAvg(agg.cs),
      dealt: getAvg(agg.dealt),
      taken: getAvg(agg.taken),
      gold: getAvg(agg.gold),
    },
    favourite: {
      champs: getTop(agg.champs, 5) as number[],
      lane: getTop(agg.lanes) as string,
      role: getTop(agg.roles) as string,
    },
  };
}

interface TeamOverview {
  aggs?: CompiledTeamAggregations,
}

interface TeamAggregations {
  matchTime: number[],
  wins: boolean[],

  towers: number[],
  dragons: number[],
  barons: number[],
}

interface CompiledTeamAggregations {
  games: number,
  winRate: string,
  avg: {
    matchTime: string,
    towers: string,
    dragons: string,
    barons: string,
  },
}

function compileTeamAggregations(agg: TeamAggregations): CompiledTeamAggregations {
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

interface Overviews {
  team: TeamOverview,
  players: PlayerOverviews,
}

@Component
export default class Overview extends Vue {
  @Prop() teamID!: string;

  @Getter(TeamGetters.MATCHES, teamsSpace) matchesData!: (t: string) => [Match] | undefined;
  @Getter(TeamGetters.TEAM, { namespace: Namespace.TEAMS }) teamData!: (id: string) => Team;

  @Getter(LeagueGetters.ITEM, leagueSpace) item!: IDGetter<ItemData>;
  @Getter(LeagueGetters.ITEM_ICON, leagueSpace) itemIcon!: IDGetter<string>; // TODO: use ITEM_SPRITE

  @Getter(LeagueGetters.CHAMP, leagueSpace) champ!: IDGetter<ChampData>;
  @Getter(LeagueGetters.CHAMP_ICON, leagueSpace) champIcon!: IDGetter<string>;

  @Getter(LeagueGetters.RUNES, leagueSpace) runes!: IDGetter<RunesData>;

  @Getter(LeagueGetters.SPELL, leagueSpace) spell!: IDGetter<SpellData>;

  get team(): Team {
    return this.teamData(this.teamID);
  }

  idToName(): { [id: string]: string } {
    const mapping: { [id: string]: string} = {};
    this.teamData(this.teamID).members.forEach((m) => {
      mapping[m.accountId] = m.name;
    });
    return mapping;
  }

  get overviews(): Overviews {
    const matches = this.matchesData(this.teamID);
    if (!matches) return { team: {}, players: {} };
    const mapping = this.idToName();
    console.debug('generating overviews', {
      matches: matches.length,
    });

    const aggs: { [name: string]: PlayerAggregations } = {};
    const teamAggs: TeamAggregations = {
      matchTime: [],
      wins: [],
      towers: [],
      dragons: [],
      barons: [],
    };
    const data: Overviews = { team: {}, players: {} }; // finalized data

    matches.forEach((m) => {
      // collect participant data
      let team: MatchTeam | undefined;
      m.details.participantIdentities.forEach((p) => {
        let id = p.player.accountId;
        if (p.player.currentAccountId && p.player.currentAccountId !== id) id = p.player.currentAccountId;
        const name = mapping[id];
        if (!name) return; // untracked player

        const part = m.details.participants.find(v => v.participantId === p.participantId);
        if (!part) {
          console.debug(`could not find participant ${p.participantId}`);
          return;
        };
        const partTeam = m.details.teams.find(v => v.teamId === part.teamId);
        if (partTeam) team = partTeam;

        // collect straight data
        if (!data.players[name]) data.players[name] = {};
        data.players[name].tier = part.highestAchievedSeasonTier; // TODO: compare to get highest tier

        // collect data for aggregation
        if (!aggs[name]) aggs[name] = newPlayerAggregation();
        updatePlayerAggregation(aggs[name], part);
      });

      // collect team data
      teamAggs.matchTime.push(m.details.gameDuration);
      if (team) {
        teamAggs.wins.push(team.win === 'Win'); // rito why is this variable a string?
        teamAggs.towers.push(team.towerKills);
        teamAggs.dragons.push(team.dragonKills);
        teamAggs.barons.push(team.baronKills);
      }
    });

    // compile player aggs
    Object.keys(aggs).forEach((k) => {
      data.players[k].aggs = compilePlayerAggregations(aggs[k]);
    });
    // compile team aggs
    data.team.aggs = compileTeamAggregations(teamAggs);

    console.log('generated overviews', data);
    return data;
  }
}

</script>
