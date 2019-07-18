<template>
  <div class="overview">
    <h2>Overview</h2>
    <div>
      <div>
        {{ overviews.team }}
        <SimpleAggregatePie :aggKeys="['avg', 'dealt']" :aggs=overviews.players />
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

import SimpleAggregatePie from './vis/SimpleAggregatePie';

import { Namespace } from '../store';
import teams, { TeamGetters } from '../store/teams';
import { LeagueGetters, IDGetter } from '../store/league';
import {
  ItemData, ChampData, RunesData, SpellData,
} from '../store/league/types';
import {
  Match, Team, Participant, MatchTeam,
} from '../api/types';
import {
  Overviews,
  PlayerAggregations, TeamAggregations,
  CompiledPlayerAggregations, CompiledTeamAggregations,
  newPlayerAggregation, updatePlayerAggregation,
  compileTeamAggregations, compilePlayerAggregations,
} from '../math/aggregates';

const leagueSpace = { namespace: Namespace.LEAGUE };
const teamsSpace = { namespace: Namespace.TEAMS };

@Component({
  components: {
    SimpleAggregatePie,
  },
})
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
