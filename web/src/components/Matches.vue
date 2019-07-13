<template>
  <div class="matches">
    <div v-if=!matches>
      No matches found for this team.
    </div>
    <div v-if=matches>
      <div v-for="m in matches" :key="m.gameId">
        Game ID: {{ m.details.gameId }}
        <div v-for="p in m.details.participants" :key="p.participantId">
          Participant: {{ p.participantId }}
          <br />
          <div>
            ChampionID: {{ p.championId }}
            {{ champIcon(p.championId) }}
            <img :src="champIcon(p.championId)" />
          </div>
          <br />
          <div>
            Spells:
            {{ spell(p.spell1Id).name }}
            {{ spell(p.spell2Id).name }}
          </div>
          <br />
          <div>
            First Item ({{ p.stats.item0 }}):
            <div v-if=!item(p.stats.item0)>
              No Item
            </div>
            <div v-if=item(p.stats.item0)>
              {{ item(p.stats.item0).name }}
              <img :src="itemIcon(p.stats.item0)" />
            </div>
          </div>
          <br />
          <div>
            <div>
              Primary Perk ({{ p.stats.perkPrimaryStyle }}):
              {{ runes(p.stats.perkPrimaryStyle).key }}
            </div>
            <div>
              Secondary Perk ({{ p.stats.perkSubStyle }}):
              {{ runes(p.stats.perkSubStyle).key }}
            </div>
          </div>
          <hr />
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
import { Match } from '../api/types';

const leagueSpace = { namespace: Namespace.LEAGUE };
const teamsSpace = { namespace: Namespace.TEAMS };

@Component
export default class Matches extends Vue {
  @Prop() teamID!: string;

  @Getter(TeamGetters.MATCHES, teamsSpace) matchesData!: (t: string) => [Match] | undefined;

  @Getter(LeagueGetters.ITEM, leagueSpace) item!: IDGetter<ItemData>;
  @Getter(LeagueGetters.ITEM_ICON, leagueSpace) itemIcon!: IDGetter<string>; // TODO: use ITEM_SPRITE

  @Getter(LeagueGetters.CHAMP, leagueSpace) champ!: IDGetter<ChampData>;
  @Getter(LeagueGetters.CHAMP_ICON, leagueSpace) champIcon!: IDGetter<string>;

  @Getter(LeagueGetters.RUNES, leagueSpace) runes!: IDGetter<RunesData>;

  @Getter(LeagueGetters.SPELL, leagueSpace) spell!: IDGetter<SpellData>;

  get matches(): [Match] | undefined {
    return this.matchesData(this.teamID);
  }
}

</script>
