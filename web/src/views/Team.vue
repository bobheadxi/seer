<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <h1>Team {{ teamID }}</h1>
    <p>These stats are collected only from games where at least 4 members from this team played together.</p>

    <div v-if="team && !loading">
      <div>
        <button v-on:click="copyMembersToClipboard()">copy to clipboard</button>
        <a v-bind:href="'http://na.op.gg/multi/query='+memberNames()" target="_blank">
          <button>open in na.op.gg</button>
        </a>
        <button v-on:click="forceFetchTeam()">refresh</button>
      </div>

      <Overview v-bind:teamID=teamID />

      <br />

      <Matches v-bind:teamID=teamID />

      <br />

      <div v-if=updateTriggered>
        Matches sync queued
      </div>
      <button v-if=!updateTriggered v-on:click="syncMatches();">
        Sync Matches
      </button>
    </div>

    <div v-if=error.occured>
      Oops an error occured! {{ error.details }}
    </div>

    <div v-if=loading>
      Loading...
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Watch } from 'vue-property-decorator';
import { State, Action, Getter } from 'vuex-class';
import { AxiosError } from 'axios';

import Matches from '@/components/Matches.vue';
import Overview from '@/components/Overview.vue';

import { ErrorState } from '../primitives';
import { Namespace } from '../store';
import {
  TeamsState, TeamActions, TeamGetters, FetchTeamPayload, UpdateTeamPayload,
} from '../store/teams';
import { LeagueActions } from '../store/league';
import * as types from '../api/types';

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

const teamsSpace = { namespace: Namespace.TEAMS };
const leagueSpace = { namespace: Namespace.LEAGUE };

@Component({
  components: {
    Matches, Overview,
  },
})
export default class Team extends Vue {
  @Action(TeamActions.FETCH_TEAM, teamsSpace)
  private fetchTeam!: (params: FetchTeamPayload) => void;
  @Action(TeamActions.UPDATE_TEAM, teamsSpace)
  private updateTeam!: (params: UpdateTeamPayload) => void;

  @Action(LeagueActions.DOWNLOAD_METADATA, leagueSpace)
  private fetchLeagueData!: (params: any) => void;

  @Getter(TeamGetters.TEAM, { namespace: Namespace.TEAMS })
  private teamData!: (id: string) => types.Team;

  @Getter(TeamGetters.MATCHES, teamsSpace)
  matchesData!: (t: string) => [types.Match] | undefined;

  error: ErrorState = { occured: false };
  loading: boolean = true;
  updateTriggered: boolean = false;

  // fetch on mount
  async mounted() {
    try {
      await this.fetchTeam({ teamID: this.teamID });
      await this.fetchLeagueData({});
    } catch (e) {
      this.error = { occured: true, details: e };
    }
    this.loading = false;
  }

  get team(): types.Team | undefined {
    return this.teamData(this.teamID);
  }

  get teamID(): string {
    return this.$route.params.team;
  }

  syncMatches() {
    this.updateTriggered = true;
    this.error = { occured: false };
    try {
      this.updateTeam({ teamID: this.teamID });
    } catch (e) {
      this.error = { occured: true, details: e };
      this.updateTriggered = false;
    }
  }

  memberNames(): string {
    if (!this.team) return '';
    return this.team.members.map(m => m.name).join(',');
  }

  copyMembersToClipboard() {
    const teamStr = this.memberNames();
    copyStringToClipboard(teamStr);
  }

  async forceFetchTeam() {
    this.loading = true;
    try {
      await this.fetchTeam({ teamID: this.teamID, force: true });
    } catch (e) {
      this.error = { occured: true, details: e };
    }
    this.loading = false;
  }
}

</script>
