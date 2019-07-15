import axios, { AxiosInstance, AxiosError } from 'axios';

import { Region, Team, Match } from '@/api/types';

// https://godoc.org/go.bobheadxi.dev/res#ErrResponse
interface ResError {
  message: string,
  error?: string,
  data?: any,
}

function isResError(arg: any): arg is ResError {
  return arg && arg.message;
}

export class SeerAPI {
  private net: AxiosInstance;

  constructor(addr: string) {
    this.net = axios.create({
      baseURL: addr,
      timeout: 10000,
      responseType: 'json',
    });
  }

  private formatError(e: AxiosError): Error {
    if (!e.isAxiosError) return e;

    console.debug('request error caught', {
      name: e.name, config: e.config, code: e.code, resp: e.response,
    });
    if (e.response && e.response.data) {
      const res = e.response.data;
      if (isResError(res)) {
        e.message = res.message;
        // TODO: append error, data
      }
    }
    return e;
  }

  async createTeam(region: Region, members: string[]): Promise<CreateTeamResponse> {
    console.debug('making POST /team request', { region, members });
    try {
      const resp = await this.net.post('/team', {
        region,
        members,
      });
      const { data } = resp.data;
      return data as CreateTeamResponse;
    } catch (e) {
      throw this.formatError(e as AxiosError);
    }
  }

  async updateTeam(teamID: string): Promise<UpdateTeamResponse> {
    console.debug('making POST /team/update/{teamID} request', { teamID });
    try {
      const resp = await this.net.post(`/team/update/${teamID}`);
      const { data } = resp.data;
      return data as UpdateTeamResponse;
    } catch (e) {
      throw this.formatError(e as AxiosError);
    }
  }

  async getTeam(teamID: string): Promise<GetTeamResponse> {
    console.debug('making GET /team/{teamID} request', { teamID });
    try {
      const resp = await this.net.get(`/team/${teamID}`);
      const { data } = resp.data;
      return data as GetTeamResponse;
    } catch (e) {
      throw this.formatError(e as AxiosError);
    }
  }
}

export interface CreateTeamResponse {
  teamID: string;
}

export interface UpdateTeamResponse {
  status: string;
  jobID: string;
}

export interface GetTeamResponse {
  team: Team;
  matches: [Match];
}
