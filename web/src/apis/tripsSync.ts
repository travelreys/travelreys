import { WebsocketBuilder } from 'websocket-ts/lib';
import { BASE_WS_URL } from './common';
import { SyncMessage, JSONPatchOp } from '../lib/tripsSync';


const TripsSyncAPI = {
  startTripSyncSession: () => {
    return new WebsocketBuilder(BASE_WS_URL).build();
  },

  makeJSONPatchOp: (op: string, path: string, value: any): JSONPatchOp => {
    return {op, path, value}
  },

  /**
   * Should avoid replacing entire list
   *
   */
  newReplaceOp: (path: string, value: any): JSONPatchOp => {
    return {op: "replace", path, value}
  },

  makeAddOp: (path: string, value: any) => {
    return {op: "add", path, value}
  },

  makeRemoveOp: (path: string, value: any) => {
    return {op: "remove", path, value}
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

