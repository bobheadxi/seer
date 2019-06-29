import Vuex, { StoreOptions, GetterTree } from 'vuex';

import { SeerAPI } from '@/api/api';
import teams from '@/store/teams';

// TODO import client, use on actions
// https://medium.com/dailyjs/mastering-vuex-zero-to-hero-e0ca1f421d45

export interface RootState {
  version: string;
  client: SeerAPI;
}

const rootState: RootState = {
  client: new SeerAPI(process.env.API_ADDR || 'localhost:8080'),
  version: '1.0.0',
}

export enum Namespace {
  TEAMS = 'teams',
}

const store: StoreOptions<RootState> = {
  state: rootState,
  modules: {
    [Namespace.TEAMS]: teams,
  },
}

const $store = new Vuex.Store<RootState>(store);

export default $store;