import { SeerAPI } from '@/api/api';

export interface RootState {
  version: string;
  client: SeerAPI;
}
