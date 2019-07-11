<template>
  <div class="matches">
    <div v-if=matches>
      <div v-for="m in matches" v-bind:key="m.gameId">
        <div>
          {{ teamID }} {{ m.details.gameId }} {{ item(m.details.participants[0].stats.item0) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

import { Vue, Component, Prop } from 'vue-property-decorator';
import { State, Action, Getter } from 'vuex-class';

import { Namespace } from '../store';
import { TeamGetters } from '../store/teams';
import { LeagueGetters, ItemData } from '../store/league';
import * as types from '../api/types';

const namespace = Namespace.TEAMS;

@Component
export default class Matches extends Vue {
  @Prop() teamID!: string;

  @Getter(TeamGetters.MATCHES, { namespace: Namespace.TEAMS })
  private matchesData!: (id: string) => [types.Match];

  @Getter(LeagueGetters.ITEM, { namespace: Namespace.LEAGUE })
  private itemData!: (id: string) => ItemData | undefined;

  get matches(): [types.Match] | undefined {
    return this.matchesData(this.teamID);
  }

  item(id: string): ItemData | undefined {
    return this.itemData(id);
  }
}

</script>
