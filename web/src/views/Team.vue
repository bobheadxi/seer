<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <p>Team {{ teamID }}</p>

    <div v-if=team>
      <div>
        {{ team.Region }}
        <br />
        {{ team.Members }}
      </div>
      <Matches v-bind:teamID=teamID />
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

import { ErrorState } from '../primitives';
import { Namespace } from '../store';
import {
  TeamsState, TeamActions, TeamGetters, FetchTeamPayload,
} from '../store/teams';
import { LeagueActions } from '../store/league';
import * as types from '../api/types';

const teamsSpace = { namespace: Namespace.TEAMS };
const leagueSpace = { namespace: Namespace.LEAGUE };

@Component({
  components: {
    Matches,
  },
})
export default class Team extends Vue {
  @Action(TeamActions.FETCH_TEAM, teamsSpace)
  private fetchTeam!: (params: FetchTeamPayload) => void;

  @Action(LeagueActions.DOWNLOAD_METADATA, leagueSpace)
  private fetchLeagueData!: (params: any) => void;

  @Getter(TeamGetters.TEAM, { namespace: Namespace.TEAMS })
  private teamData!: (id: string) => types.Team;

  error: ErrorState = { occured: false };
  loading: boolean = true;

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
}

</script>
