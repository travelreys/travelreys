import { WebsocketBuilder } from 'websocket-ts/lib';
import { BASE_WS_URL } from './common';

const TripsSyncAPI = {
  startTripSyncSession: () => {
    return new WebsocketBuilder(BASE_WS_URL).build();
  },
};

export default TripsSyncAPI;

