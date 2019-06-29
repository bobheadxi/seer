import { Module, GetterTree, MutationTree, ActionTree } from 'vuex';
import { RootState } from '@/store';

import { SeerAPI } from '@/api/api';
import * as types from '@/api/types';

export interface TeamsState {
  teams: Map<string, types.Team>;
  matches: Map<string, [types.Match]>;
  updateStatus: Map<string, any>;
};

const state: TeamsState = {
  teams: new Map<string, types.Team>(),
  matches: new Map<string, [types.Match]>(),
  updateStatus: new Map<string, any>(),
}

export enum Getters {
  CLIENT = 'CLIENT',
  TEAM = 'TEAM',
  MATCHES = 'MATCHES',
  UPDATE_STATE = 'UPDATE_STATE',
}

const getterTree: GetterTree<TeamsState, RootState> = {
  [Getters.CLIENT]: (state, getters, rootState): SeerAPI => rootState.client,

  [Getters.TEAM]: (state): ((teamID: string) => types.Team | undefined) => {
    return (teamID) => state.teams.get(teamID);
  },

  [Getters.MATCHES]: (state): ((teamID: string) => [types.Match] | undefined) => {
    return (teamID) => state.matches.get(teamID);
  },

  [Getters.UPDATE_STATE]: (state): ((teamID: string) => any | undefined) => {
    return (teamID) => state.updateStatus.get(teamID);;
  },
}

export enum Actions {
  CREATE_TEAM = 'CREATE_TEAM',
  GET_TEAM = 'GET_TEAM',
  UPDATE_TEAM = 'UPDATE_TEAM',
}

const actionTree: ActionTree<TeamsState, RootState> = {
  [Actions.CREATE_TEAM]: async (context, payload: { region: types.Region, members: [string] }) => {
    const { client } = context.rootState;
    const { region, members } = payload;
    client.createTeam(region, members);
  },

  [Actions.GET_TEAM]: async (context, payload: { teamID: string }) => {
    const { client } = context.rootState;
    const { teamID } = payload;
    const { team, matches } = await client.getTeam(teamID);
    context.commit('STORE_TEAM', { teamID, team });
    context.commit('STORE_MATCHES', { teamID, matches });
  },

  [Actions.UPDATE_TEAM]: async (context, payload: { teamID: string }) => {
    const { client } = context.rootState;
    const { teamID } = payload;
    const resp = await client.updateTeam(teamID);
    context.commit('SET_TEAM_UPDATE_STATUS', resp.status);
  }
}

const mutationTree: MutationTree<TeamsState> = {
  STORE_TEAM: (state, payload: { teamID: string, team: types.Team }) => {
    state.teams.set(payload.teamID, payload.team);
  },

  STORE_MATCHES: (state, payload: { teamID: string, matches: [types.Match] }) => {
    state.matches.set(payload.teamID, payload.matches);
  },

  SET_TEAM_UPDATE_STATUS: (state, payload: { teamID: string, status: any }) => {
    state.updateStatus.set(payload.teamID, payload.status);
  }
}

const teams: Module<TeamsState, RootState> = {
  namespaced: true,
  state,
  getters: getterTree,
  actions: actionTree,
  mutations: mutationTree,
};

export default teams;
