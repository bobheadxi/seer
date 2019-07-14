<template>
  <div class="matches">
    <div>
      <div v-for="m in team.members" v-bind:key="m.id">
        {{ m.name }} ({{ m.summonerLevel }})
      </div>
      <button v-on:click="copyMembersToClipboard()">copy to clipboard</button>
      <a v-bind:href="'http://na.op.gg/multi/query='+memberNames()" target="_blank">
        <button>open in na.op.gg</button>
      </a>
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
import { Match, Team } from '../api/types';

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

  get matches(): [Match] | undefined {
    return this.matchesData(this.teamID);
  }

  get team(): Team | undefined {
    return this.teamData(this.teamID);
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
