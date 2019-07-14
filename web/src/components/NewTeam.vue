<template>
  <div class="new-team">
    <div>
      <h1>Create a Team</h1>
      <form class="new-team-form" @submit.prevent="onSubmit">
        <p>
          <label for="inputRegion">Region:</label>
          <select id="inputRegion" v-model=inputRegion>
            <option v-for="r in regions" v-bind:key="r">
              {{r}}
            </option>
          </select>
        </p>

        <fieldset>
          <legend>Members:</legend>
          <input
            type="text"
            v-for="(m, i) in inputMembers"
            v-model="inputMembers[i]"
            v-bind:key="'teamMember'+i" />
          <button
            type="button"
            v-if="inputMembers.length < 7"
            v-on:click="inputMembers.push(null);" >
            Add another member
          </button>
        </fieldset>

        <p>
          <input type="submit" value="Submit">
        </p>
      </form>
    </div>

    <div v-if="data.createdTeam">
      Team {{ data.createdTeam.teamID }} created!
    </div>

    <div v-if=loading>
      Creating team...
    </div>

    <div v-if=error.occured>
      Oops an error occured! {{ error.details }}
    </div>
  </div>
</template>

<script lang="ts">

import { Vue, Component, Prop } from 'vue-property-decorator';
import { State, Action, Getter } from 'vuex-class';

import { Namespace } from '../store';
import teams, { TeamActions, CreateTeamPayload } from '../store/teams';
import { Region } from '../api/types';
import { CreateTeamResponse } from '../api';
import { ErrorState } from '../primitives';

const teamsSpace = { namespace: Namespace.TEAMS };

@Component
export default class NewTeam extends Vue {
  inputRegion: Region = Region.NA1;
  inputMembers: (string | null)[] = [null, null, null, null, null];
  data: {
    createdTeam?: CreateTeamResponse,
  } = {};

  @Action(TeamActions.CREATE_TEAM, teamsSpace)
  private createTeam!: (p: CreateTeamPayload) => CreateTeamResponse;

  error: ErrorState = { occured: false };
  loading: boolean = false;

  async onSubmit() {
    this.error = { occured: false };
    this.loading = true;
    try {
      this.data.createdTeam = await this.createTeam({
        region: this.inputRegion,
        members: this.inputMembers.filter((x): x is string => x !== null && x !== '' && x !== undefined),
      });
      if (this.data.createdTeam) {
        console.debug(`created team ${this.inputRegion}.${this.data.createdTeam.teamID}`,
          { members: this.inputMembers });
        this.resetForm();
      } else {
        this.error = { occured: true, details: 'Team creation did not fail, but no team was returned' };
      }
    } catch (e) {
      this.error = { occured: true, details: e };
    }
    this.loading = false;
  }

  resetForm() {
    this.inputRegion = Region.NA1;
    this.inputMembers = [null, null, null, null, null];
  }

  get regions(): string[] {
    return Object.keys(Region);
  }
}

</script>
