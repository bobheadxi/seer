import axios, { AxiosInstance } from 'axios';

import { Region, Team, Match } from '@/api/types';

export class SeerAPI {
  net: AxiosInstance;

  constructor(addr: string) {
    this.net = axios.create({
      baseURL: addr,
      timeout: 5000,
      responseType: 'json',
    });
  }

  async createTeam(region: Region, members: [string]): Promise<CreateTeamResponse> {
    const resp = await this.net.post('/team', {
      region,
      members,
    });
    const { data } = resp.data;
    return data as CreateTeamResponse;
  }

  async updateTeam(teamID: string): Promise<UpdateTeamResponse> {
    const resp = await this.net.post(`/team/${teamID}`);
    const { data } = resp.data;
    return data as UpdateTeamResponse;
  }

  async getTeam(teamID: string): Promise<GetTeamResponse> {
    const resp = await this.net.get(`/team/${teamID}`);
    const { data } = resp.data;
    return data as GetTeamResponse;
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
