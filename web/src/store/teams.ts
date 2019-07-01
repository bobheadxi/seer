import Vue from 'vue';
import {
  Module, GetterTree, MutationTree, ActionTree,
} from 'vuex';
import { RootState } from '@/store/root';

import { SeerAPI } from '@/api/api';
import * as types from '@/api/types';

// TODO NO MAPS????


export interface TeamsState {
  teams: { id: string, data: types.Team }[];
  matches: { id: string, data: [types.Match] }[];
  updateStatus: { id: string, data: any }[];
};

const teamState: TeamsState = {
  teams: [],
  matches: [],
  updateStatus: [],
};

export enum Getters {
  CLIENT = 'CLIENT',
  TEAM = 'TEAM',
  MATCHES = 'MATCHES',
  UPDATE_STATE = 'UPDATE_STATE',
}

const getterTree: GetterTree<TeamsState, RootState> = {
  [Getters.CLIENT]: (state, getters, rootState): SeerAPI => rootState.client,

  [Getters.TEAM]: (state): ((teamID: string) => types.Team | undefined) => (teamID) => {
    const found = state.teams.find(v => v.id === teamID);
    return found ? found.data : undefined;
  },

  [Getters.MATCHES]: (state): ((teamID: string) => [types.Match] | undefined) => (teamID) => {
    const found = state.matches.find(v => v.id === teamID);
    return found ? found.data : undefined;
  },

  [Getters.UPDATE_STATE]: (state): ((teamID: string) => any | undefined) => (teamID) => {
    const found = state.updateStatus.find(v => v.id === teamID);
    return found ? found.data : undefined;
  },
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
    state.teams.push({ id: payload.teamID, data: payload.team });
  },

  STORE_MATCHES: (state, payload: { teamID: string, matches: [types.Match] }) => {
    state.matches.push({ id: payload.teamID, data: payload.matches });
  },

  SET_TEAM_UPDATE_STATUS: (state, payload: { teamID: string, status: any }) => {
    state.updateStatus.push({ id: payload.teamID, data: payload.status });
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
