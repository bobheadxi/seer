<template>
  <div class="matches">
    <div v-if=matches>
      <div v-for="m in matches" v-bind:key="m.gameId">
        <div>
          {{ teamID }} {{ m.details.gameId }} {{ m.details.participants }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

import { Vue, Component, Prop } from 'vue-property-decorator';
import { State, Action, Getter } from 'vuex-class';

import { Namespace } from '../store';
import { Getters } from '../store/teams';
import * as types from '../api/types';

const namespace = Namespace.TEAMS;

@Component
export default class Matches extends Vue {
  @Prop() teamID!: string;

  @Getter(Getters.MATCHES, { namespace }) private matchesData!: (id: string) => [types.Match];

  get matches(): [types.Match] | undefined {
    return this.matchesData(this.teamID);
  }
}

</script>
