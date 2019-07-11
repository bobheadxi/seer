import {
  Module, GetterTree, MutationTree, ActionTree, ActionContext,
} from 'vuex';
import { RootState } from '@/store/root';
import Axios from 'axios';

interface Versions {
  item: string;
  rune: string;
  summoner: string;
  champion: string;
  mastery: string;
  map: string;
}

export interface ItemData {
  name: string;
  description: string;
  plaintext: string;
  image: Image;
}

export interface ChampData {
  key: number;
  title: string;
  tags: [string];
  image: Image;
}

export interface RuneData {

}

export interface Image {
  full: string;
  sprite: string;
  group: string;
  x: number;
  y: number;
  w: number;
  h: number;
}

export interface LeagueMetadataState {
  version: string;
  downloaded: string;
  items: { [id: number]: ItemData };
  champs: { [name: string]: ChampData };
};

const leagueMetadataState: LeagueMetadataState = {
  version: '',
  downloaded: '',
  items: {},
  champs: {},
};

export enum LeagueGetters {
  ITEM = 'ITEM',
  ITEM_SPRITE = 'ITEM_SPRITE',
}

const getterTree: GetterTree<LeagueMetadataState, RootState> = {
  [LeagueGetters.ITEM]: (state): ((item: number) => ItemData | undefined) => item => state.items[item],
  [LeagueGetters.ITEM_SPRITE]: (state): ((item: number) => string | undefined) => (item) => {
    const i = state.items[item];
    return `http://ddragon.leagueoflegends.com/cdn/${state.version}/img/sprite/${i.image.sprite}`;
  },
};

async function getAndCommit(context: ActionContext<LeagueMetadataState, RootState>, mutation: string, source: string, path?: string): Promise<any> {
  const resp = await Axios.get(source);
  let { data } = resp.data;
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

  DOWNLOAD_DATA_FOR_VERSIONS: async (context, version) => {
    console.debug(`fetching v${version} data`);
    await getAndCommit(context, 'STORE_ITEMS', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_GB/item.json`);
    await getAndCommit(context, 'STORE_CHAMPS', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_GB/champion.json`);
    await getAndCommit(context, 'STORE_RUNES', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_US/runesReforged.json`);
    await getAndCommit(context, 'STORE_SUMMONERS', `http://ddragon.leagueoflegends.com/cdn/${version}/data/en_GB/summoner.json`);
    context.commit('SET_DOWNLOADED', { version });
  },
};

const mutationTree: MutationTree<LeagueMetadataState> = {
  /* eslint-disable no-param-reassign */
  SET_VERSION: (state, payload: { version: string }) => { state.version = payload.version; },
  SET_DOWNLOADED: (state, payload: { version: string }) => { state.downloaded = payload.version; },
  STORE_ITEMS: (state, payload) => { state.items = payload; },
  STORE_CHAMPS: (state, payload) => { state.champs = payload; },

  // TODO
  STORE_RUNES: (state, payload) => { },
  STORE_SUMMONERS: (state, payload) => { },
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
