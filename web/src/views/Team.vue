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
import { TeamsState, Actions, Getters } from '../store/teams';
import * as types from '../api/types';

const namespace = Namespace.TEAMS;

@Component({
  components: {
    Matches,
  },
})
export default class Team extends Vue {
  @Action(Actions.FETCH_TEAM, { namespace }) private fetchTeam!: (params: any) => void;

  @Getter(Getters.TEAM, { namespace }) private teamData!: (id: string) => types.Team;

  error: ErrorState;

  loading: boolean;

  constructor() {
    super();
    this.loading = true;
    this.error = { occured: false };
  };

  // fetch on mount
  async mounted() {
    try {
      await this.fetchTeam({ teamID: this.teamID });
    } catch (e) {
      this.error = { occured: true, details: e };
    }
    this.loading = false;
    console.debug('loaded!', this.teamID, this.team);
  }


  get team(): types.Team | undefined {
    return this.teamData(this.teamID);
  }

  get teamID(): string {
    return this.$route.params.team;
  }
}

</script>
