<template>
  <div class="overview">
    <h2>Overview</h2>
    <div>
      <button v-on:click="copyMembersToClipboard()">copy to clipboard</button>
      <a v-bind:href="'http://na.op.gg/multi/query='+memberNames()" target="_blank">
        <button>open in na.op.gg</button>
      </a>
    </div>
    <div>
      <div v-for="m in team.members" v-bind:key="m.id">
        <h3>
          {{ m.name }} ({{ playerOverviews[m.name].tier }}, lv{{ m.summonerLevel }})
          <!-- TODO: regions? --->
          <a v-bind:href="'https://na.op.gg/summoner/userName='+m.name" target="_blank">
            <img width="16" height="16"
              src="https://lh3.googleusercontent.com/UdvXlkugn0bJcwiDkqHKG5IElodmv-oL4kHlNAklSA2sdlVWhojsZKaPE-qFPueiZg" />
          </a>
        </h3>
        <div>
          <h5>Most played lane and role</h5>
          {{ playerOverviews[m.name].aggs.favourite.lane }}
          ({{ playerOverviews[m.name].aggs.favourite.role }})
        </div>
        <div>
          <h5>Most played champions</h5>
          <img
            v-for="c in playerOverviews[m.name].aggs.favourite.champs"
            v-bind:key="'fav-'+m.name+'-'+c"
            v-bind:src="champIcon(c)" />
        </div>
        <div>
          <h5>Average stats</h5>
          {{ playerOverviews[m.name].aggs.avg }}
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
import { Match, Team, Participant } from '../api/types';

const leagueSpace = { namespace: Namespace.LEAGUE };
const teamsSpace = { namespace: Namespace.TEAMS };

function copyStringToClipboard(str: string) {
  const el = document.createElement('textarea');
  el.value = str;
  el.setAttribute('readonly', '');
  document.body.appendChild(el);
  el.select();
  // Copy text to clipboard
  document.execCommand('copy');
  document.body.removeChild(el);
}

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

  get playerOverviews(): PlayerOverviews {
    const matches = this.matchesData(this.teamID);
    if (!matches) return {};
    const mapping = this.idToName();
    const aggs: { [name: string]: PlayerAggregations } = {};
    const data: PlayerOverviews = {}; // finalized data

    matches.forEach((m) => {
      m.details.participantIdentities.forEach((p) => {
        let id = p.player.accountId;
        if (p.player.currentAccountId) id = p.player.currentAccountId;
        const name = mapping[id];
        if (!name) return; // untracked player

        const part = m.details.participants.find(v => v.participantId === p.participantId);
        if (!part) {
          console.debug(`could not find participant ${p.participantId}`);
          return;
        };

        // collect straight data
        if (!data[name]) data[name] = {};
        data[name].tier = part.highestAchievedSeasonTier; // TODO: compare to get highest tier

        // collect data for aggregation
        if (!aggs[name]) aggs[name] = newPlayerAggregation();
        updatePlayerAggregation(aggs[name], part);
      });
    });

    Object.keys(aggs).forEach((k) => {
      data[k].aggs = compilePlayerAggregations(aggs[k]);
    });

    return data;
  }

  memberNames(): string {
    if (!this.team) return '';
    return this.team.members.map(m => m.name).join(',');
  }

  copyMembersToClipboard() {
    const teamStr = this.memberNames();
    copyStringToClipboard(teamStr);
  }
}

</script>
