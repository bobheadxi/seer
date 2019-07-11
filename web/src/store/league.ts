import Vue from 'vue';
import {
  Module, GetterTree, MutationTree, ActionTree,
} from 'vuex';
import { RootState } from '@/store/root';
import Axios from 'axios';

function getItemSpriteSource(sprite: string): string {
  return `http://ddragon.leagueoflegends.com/cdn/6.24.1/img/sprite/${sprite}`;
}

export interface ItemData {
  name: string;
  description: string;
  plaintext: string;
  image: ItemImage;
}

export interface ItemImage {
  sprite: string;
  group: string;
  x: number;
  y: number;
  w: number;
  h: number;
}

export interface LeagueMetadataState {
  items: { [id: number]: ItemData };
};

const leagueMetadataState: LeagueMetadataState = {
  items: {},
};

export enum LeagueGetters {
  ITEM = 'ITEM',
}

const getterTree: GetterTree<LeagueMetadataState, RootState> = {
  [LeagueGetters.ITEM]: (state): ((item: number) => ItemData | undefined) => item => state.items[item],
};

export enum LeagueActions {
  DOWNLOAD_ITEMS = 'DOWNLOAD_ITEMS',
}

const actionTree: ActionTree<LeagueMetadataState, RootState> = {
  [LeagueActions.DOWNLOAD_ITEMS]: async (context, force = false) => {
    if (context.state.items && !force) return;

    const resp = await Axios.get('http://ddragon.leagueoflegends.com/cdn/6.24.1/data/en_US/item.json');
    const { data } = resp.data;
    context.commit('STORE_ITEMS', data);
  },
};

const mutationTree: MutationTree<LeagueMetadataState> = {
  STORE_ITEMS: (state, payload: Map<number, ItemData>) => {
    Vue.set(state, 'items', payload);
    console.debug('stored item metadata', { payload });
  },
};

const leagues: Module<LeagueMetadataState, RootState> = {
  namespaced: true,
  state: leagueMetadataState,
  getters: getterTree,
  actions: actionTree,
  mutations: mutationTree,
};

export default leagues;
