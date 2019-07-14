import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';
import createPersistedState from 'vuex-persistedstate';

import { SeerAPI } from '@/api';
import { RootState } from '@/store/root';
import teams from '@/store/teams';
import league from '@/store/league';

Vue.use(Vuex);

export enum Namespace {
  TEAMS = 'teams',
  LEAGUE = 'league',
}

const store: StoreOptions<RootState> = {
  state: {
    client: new SeerAPI(process.env.VUE_APP_API_ADDR || 'localhost:8080'),
    version: '1.0.0',
  },
  modules: {
    [Namespace.TEAMS]: teams,
    [Namespace.LEAGUE]: league,
  },
  plugins: [
    createPersistedState({
      paths: [Namespace.LEAGUE, Namespace.TEAMS],
    }),
  ],
};

export default new Vuex.Store<RootState>(store);
