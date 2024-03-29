import {
  Module, GetterTree, MutationTree, ActionTree,
} from 'vuex';
import { RootState } from '@/store/root';

import { SeerAPI, CreateTeamResponse } from '@/api';
import * as types from '@/api/types';

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

export enum TeamGetters {
  CLIENT = 'CLIENT',
  TEAM = 'TEAM',
  MATCHES = 'MATCHES',
  UPDATE_STATE = 'UPDATE_STATE',
}

const getterTree: GetterTree<TeamsState, RootState> = {
  [TeamGetters.CLIENT]: (state, getters, rootState): SeerAPI => rootState.client,

  [TeamGetters.TEAM]: (state): ((teamID: string) => types.Team | undefined) => (teamID) => {
    const found = state.teams.find(v => v.id === teamID);
    return found ? found.data : undefined;
  },

  [TeamGetters.MATCHES]: (state): ((teamID: string) => [types.Match] | undefined) => (teamID) => {
    const found = state.matches.find(v => v.id === teamID);
    return found ? found.data : undefined;
  },

  [TeamGetters.UPDATE_STATE]: (state): ((teamID: string) => any | undefined) => (teamID) => {
    const found = state.updateStatus.find(v => v.id === teamID);
    return found ? found.data : undefined;
  },
};

export enum TeamActions {
  CREATE_TEAM = 'CREATE_TEAM',
  FETCH_TEAM = 'FETCH_TEAM',
  UPDATE_TEAM = 'UPDATE_TEAM',
}

export interface CreateTeamPayload { region: types.Region, members: string[] }
export interface FetchTeamPayload { teamID: string, force?: boolean }
export interface UpdateTeamPayload { teamID: string }

const actionTree: ActionTree<TeamsState, RootState> = {
  [TeamActions.CREATE_TEAM]: async (context, payload: CreateTeamPayload): Promise<CreateTeamResponse> => {
    const { client } = context.rootState;
    const { region, members } = payload;
    const resp = await client.createTeam(region, members);
    return resp;
  },

  [TeamActions.FETCH_TEAM]: async (context, payload: FetchTeamPayload) => {
    const { teamID, force } = payload;
    if (context.state.teams.find(v => v.id === teamID) && !force) return;
    if (force) context.commit('REMOVE_TEAM_DATA', { teamID });

    const { client } = context.rootState;
    const { team, matches } = await client.getTeam(teamID);
    context.commit('STORE_TEAM', { teamID, team });
    context.commit('STORE_MATCHES', { teamID, matches });
    console.debug('stored team and matches', { teamID, team, matches });
  },

  [TeamActions.UPDATE_TEAM]: async (context, payload: UpdateTeamPayload) => {
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
    state.updateStatus = state.updateStatus.filter(v => v.id !== payload.teamID);
    state.updateStatus.push({ id: payload.teamID, data: payload.status });
  },

  REMOVE_TEAM_DATA: (state, payload: { teamID: string }) => {
    state.teams = state.teams.filter(v => v.id !== payload.teamID);
    state.matches = state.matches.filter(v => v.id !== payload.teamID);
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
