import Vue from 'vue';
import {
  Module, GetterTree, MutationTree, ActionTree,
} from 'vuex';
import { RootState } from '@/store/root';

import { SeerAPI } from '@/api/api';
import * as types from '@/api/types';

// TODO NO MAPS???? UGHHHH


export interface TeamsState {
  teams: Map<string, types.Team>;
  matches: Map<string, [types.Match]>;
  updateStatus: Map<string, any>;
};

const teamState: TeamsState = {
  teams: new Map<string, types.Team>(),
  matches: new Map<string, [types.Match]>(),
  updateStatus: new Map<string, any>(),
};

export enum Getters {
  CLIENT = 'CLIENT',
  TEAM = 'TEAM',
  MATCHES = 'MATCHES',
  UPDATE_STATE = 'UPDATE_STATE',
}

const getterTree: GetterTree<TeamsState, RootState> = {
  [Getters.CLIENT]: (state, getters, rootState): SeerAPI => rootState.client,

  [Getters.TEAM]: (state): ((teamID: string) => types.Team | undefined) => teamID => state.teams.get(teamID),

  [Getters.MATCHES]: (state): ((teamID: string) => [types.Match] | undefined) => teamID => state.matches.get(teamID),

  [Getters.UPDATE_STATE]: (state): ((teamID: string) => any | undefined) => teamID => state.updateStatus.get(teamID),
};

export enum Actions {
  CREATE_TEAM = 'CREATE_TEAM',
  FETCH_TEAM = 'FETCH_TEAM',
  UPDATE_TEAM = 'UPDATE_TEAM',
}

const actionTree: ActionTree<TeamsState, RootState> = {
  [Actions.CREATE_TEAM]: async (context, payload: { region: types.Region, members: [string] }) => {
    const { client } = context.rootState;
    const { region, members } = payload;
    await client.createTeam(region, members);
  },

  [Actions.FETCH_TEAM]: async (context, payload: { teamID: string }) => {
    const { client } = context.rootState;
    const { teamID } = payload;
    const { team, matches } = await client.getTeam(teamID);
    context.commit('STORE_TEAM', { teamID, team });
    context.commit('STORE_MATCHES', { teamID, matches });
    console.log('stored', { teamID, team, matches });
  },

  [Actions.UPDATE_TEAM]: async (context, payload: { teamID: string }) => {
    const { client } = context.rootState;
    const { teamID } = payload;
    const resp = await client.updateTeam(teamID);
    context.commit('SET_TEAM_UPDATE_STATUS', resp.status);
  },
};

const mutationTree: MutationTree<TeamsState> = {
  STORE_TEAM: (state, payload: { teamID: string, team: types.Team }) => {
    Vue.set(state.teams, payload.teamID, payload.team);
  },

  STORE_MATCHES: (state, payload: { teamID: string, matches: [types.Match] }) => {
    Vue.set(state.matches, payload.teamID, payload.matches);
  },

  SET_TEAM_UPDATE_STATUS: (state, payload: { teamID: string, status: any }) => {
    Vue.set(state.updateStatus, payload.teamID, payload.status);
  },
};

const teams: Module<TeamsState, RootState> = {
  namespaced: true,
  state: teamState,
  getters: getterTree,
  actions: actionTree,
  mutations: mutationTree,
};

export default teams;
