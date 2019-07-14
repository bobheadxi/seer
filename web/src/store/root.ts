import { SeerAPI } from '@/api';

export interface RootState {
  version: string;
  client: SeerAPI;
}
