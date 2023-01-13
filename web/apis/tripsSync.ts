import { WebsocketBuilder } from 'websocket-ts/lib';
import { BASE_WS_URL } from './common';

export interface SyncMessage {
  id: string
  counter: number
  tripPlanID: string
  opType: string

  syncDataJoinSession?: SyncDataJoinSession
  syncDataLeaveSession?: SyncDataLeaveSession
  syncDataPing?: SyncDataPing
  syncDataUpdateTrip?: SyncDataUpdateTrip
}

export interface SyncDataJoinSession {
  memberID: string
  memberEmail: string
}

export interface SyncDataLeaveSession {
  memberID: string
  memberEmail: string
}

export interface SyncDataPing {}

export interface SyncDataUpdateTrip {
  op: string
  path: string
  value: string
}

const TripsSyncAPI = {
  startTripSyncSession: () => {
    return new WebsocketBuilder(BASE_WS_URL).build();
  },

  makeSyncMsgJoinSession: (tripPlanID: string, memberID: string, memberEmail: string): SyncMessage => {
    return {
      tripPlanID,
      opType: "SyncOpJoinSession",
      syncDataJoinSession: {memberID, memberEmail}
    } as SyncMessage
  },

  makeSyncMsgLeaveSession: (tripPlanID: string, memberID: string, memberEmail: string): SyncMessage => {
    return {
      tripPlanID,
      opType: "SyncOpLeaveSession",
      syncDataLeaveSession: {memberID, memberEmail}
    } as SyncMessage
  },

  makeSyncMsgUpdateTrip: (tripPlanID: string, op: string, path: string, value: string): SyncMessage => {
    return {
      tripPlanID,
      opType: "SyncOpUpdateTrip",
      syncDataUpdateTrip: { op, path, value },
    } as SyncMessage;
  }
};

export default TripsSyncAPI;

