import {
  Module, GetterTree, MutationTree, ActionTree, ActionContext,
} from 'vuex';
import Axios from 'axios';

import { RootState } from '@/store/root';
import * as LeagueTypes from './types';
import { findChamp, findSpell } from './query';

export interface LeagueMetadataState {
  version: string;
  downloaded: string;
  items: { [id: number]: LeagueTypes.ItemData };
  champs: { [name: string]: LeagueTypes.ChampData };
  runes: LeagueTypes.RunesData[];
  spells: { [name: string]: LeagueTypes.SpellData };
};

const leagueMetadataState: LeagueMetadataState = {
  version: '',
  downloaded: '',
  items: {},
  champs: {},
  runes: [],
  spells: {},
};

export interface IDGetter<T> { (key: number): T | undefined };

export enum LeagueGetters {
  ITEM = 'ITEM',
  ITEM_ICON = 'ITEM_ICON',
  ITEM_SPRITE = 'ITEM_SPRITE',
  CHAMP = 'CHAMP',
  CHAMP_ICON = 'CHAMP_ICON',
  RUNES = 'RUNES',
  SPELL = 'SPELL',
}

const getterTree: GetterTree<LeagueMetadataState, RootState> = {
  [LeagueGetters.ITEM]: (state): IDGetter<LeagueTypes.ItemData> => item => state.items[item],
  [LeagueGetters.ITEM_ICON]: (state): IDGetter<string> => (item) => {
    const { version } = state;
    return `http://ddragon.leagueoflegends.com/cdn/${version}/img/item/${item}.png`;
  },
  [LeagueGetters.ITEM_SPRITE]: (state): IDGetter<string> => (item) => {
    const { version } = state;
    const i = state.items[item];
    return `http://ddragon.leagueoflegends.com/cdn/${version}/img/sprite/${i.image.sprite}`;
  },

  [LeagueGetters.CHAMP]: (state): IDGetter<LeagueTypes.ChampData> => id => findChamp(state.champs, id.toString()),
  [LeagueGetters.CHAMP_ICON]: (state): IDGetter<string> => (id) => {
    const { version } = state;
    const champ = findChamp(state.champs, id.toString());
    if (!champ) return '';
    return `http://ddragon.leagueoflegends.com/cdn/${version}/img/champion/${champ.name.replace(' ', '')}.png`;
  },
  // TODO: sprites

  [LeagueGetters.RUNES]: (state): IDGetter<LeagueTypes.RunesData> => perk => state.runes.find(v => v.id === perk),
  // TODO: icons, sprites

  [LeagueGetters.SPELL]: (state): IDGetter<LeagueTypes.SpellData> => spell => findSpell(state.spells, spell.toString()),
  // TODO: icons, spirtes
};

async function getAndCommit(
  context: ActionContext<LeagueMetadataState, RootState>,
  mutation: string,
  source: string,
  path?: string,
): Promise<any> {
  const resp = await Axios.get(source);
  let { data } = resp;
  if (path) data = data[path];
  console.debug(`fetched data for ${mutation}`, { found: data });
  context.commit(mutation, data);
  return data;
}

export enum LeagueActions {
  DOWNLOAD_METADATA = 'DOWNLOAD_METADATA',
}

const actionTree: ActionTree<LeagueMetadataState, RootState> = {
  [LeagueActions.DOWNLOAD_METADATA]: async (context, force = true) => {
    const version = await getAndCommit(context, 'SET_VERSION', 'https://ddragon.leagueoflegends.com/realms/na.json', 'dd');
    if (context.state.downloaded === version && !force) return;

    context.dispatch('DOWNLOAD_DATA_FOR_VERSION', version);
  },

  DOWNLOAD_DATA_FOR_VERSION: async (context, version) => {
    console.debug(`fetching v${version} data`);
    await getAndCommit(context, 'STORE_ITEMS', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_GB/item.json`, 'data');
    await getAndCommit(context, 'STORE_CHAMPS', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_GB/champion.json`, 'data');
    await getAndCommit(context, 'STORE_RUNES', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_US/runesReforged.json`);
    await getAndCommit(context, 'STORE_SUMMONERS', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_GB/summoner.json`, 'data');
    context.commit('SET_DOWNLOADED', { version });
  },
};

const mutationTree: MutationTree<LeagueMetadataState> = {
  /* eslint-disable no-param-reassign */
  SET_VERSION: (state, version) => { state.version = version; console.debug('downloaded version set', state.version); },
  SET_DOWNLOADED: (state, payload: { version: string }) => { state.downloaded = payload.version; },
  STORE_ITEMS: (state, payload) => { state.items = payload; },
  STORE_CHAMPS: (state, payload) => { state.champs = payload; },
  STORE_RUNES: (state, payload) => { state.runes = payload; },
  STORE_SUMMONERS: (state, payload) => { state.spells = payload; },
  /* eslint-enable no-param-reassign */
};

const leagues: Module<LeagueMetadataState, RootState> = {
  namespaced: true,
  state: leagueMetadataState,
  getters: getterTree,
  actions: actionTree,
  mutations: mutationTree,
};

export default leagues;
