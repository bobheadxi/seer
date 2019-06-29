import axios, { AxiosInstance, AxiosResponse } from 'axios';

import { Region, Team, Match } from '@/api/types';

export class SeerAPI {
  net: AxiosInstance;

  constructor(addr: string) {
    this.net = axios.create({
      baseURL: addr,
      timeout: 1000,
    });
  }

  async createTeam(region: Region, members: [string]): Promise<CreateTeamResponse> {
    const resp = await this.net.post('/team', {
      region: region,
      members: members,
    })
    return resp.data as CreateTeamResponse;
  }

  async updateTeam(teamID: string): Promise<UpdateTeamResponse> {
    const resp = await this.net.post(`/team/${teamID}`)
    return resp.data as UpdateTeamResponse
  }

  async getTeam(teamID: string): Promise<GetTeamResponse> {
    const resp = await this.net.get(`/team/${teamID}`);
    return resp.data as GetTeamResponse;
  }
}

export interface CreateTeamResponse {
  teamID: string;
}

export interface UpdateTeamResponse {
  status: string;
}

export interface GetTeamResponse {
  team: Team;
  matches: [Match];
}
