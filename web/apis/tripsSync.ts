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
  ops: Array<JSONPatchOp>
}

export interface JSONPatchOp {
  op: string
  path: string
  value: string
}

const TripsSyncAPI = {
  startTripSyncSession: () => {
    return new WebsocketBuilder(BASE_WS_URL).build();
  },

  makeJSONPatchOp: (op: string, path: string, value: string): JSONPatchOp => {
    return {op, path, value}
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

  makeSyncMsgUpdateTrip: (tripPlanID: string, ops: Array<JSONPatchOp>): SyncMessage => {
    return {
      tripPlanID,
      opType: "SyncOpUpdateTrip",
      syncDataUpdateTrip: { ops },
    } as SyncMessage;
  }
};

export default TripsSyncAPI;

