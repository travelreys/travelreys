export interface JSONPatchOp {
  op: "add" | "remove" | "replace"
  path: string
  value: string
};

/**
 * Should avoid replacing entire list
 */
export const makeReplaceOp =  (path: string, value: any): JSONPatchOp => {
  return {op: "replace", path, value}
};

export const makeAddOp = (path: string, value: any): JSONPatchOp => {
  return {op: "add", path, value}
};

export const makeRemoveOp = (path: string, value: any): JSONPatchOp => {
  return {op: "remove", path, value}
};

export const SyncOpJoinSession = "SyncOpJoinSession";
export const SyncOpLeaveSession = "SyncOpLeaveSession";
export const SyncOpJoinSessionBroadcast = "SyncOpJoinSessionBroadcast";
export const SyncOpUpdateTrip = "SyncOpUpdateTrip";

export namespace TripSync {
  export interface Message {
    id?: string
    counter?: number
    tripPlanID: string
    opType: string

    syncDataJoinSession?: SyncDataJoinSession
    syncDataLeaveSession?: SyncDataLeaveSession
    syncDataPing?: SyncDataPing
    syncDataUpdateTrip?: SyncDataUpdateTrip
  }

  interface SyncDataJoinSession {
    memberID: string
  }

  interface SyncDataLeaveSession {
    memberID: string
  }

  interface SyncDataPing {}

  interface SyncDataUpdateTrip {
    ops: Array<JSONPatchOp>
  }

};

export const makeSyncMsgJoinSession = (tripPlanID: string, memberID: string): TripSync.Message => {
  return {
    tripPlanID,
    opType: SyncOpJoinSession,
    syncDataJoinSession: {memberID}
  }
}

export const makeSyncMsgLeaveSession = (tripPlanID: string, memberID: string): TripSync.Message => {
  return {
    tripPlanID,
    opType: SyncOpLeaveSession,
    syncDataLeaveSession: {memberID}
  }
}

export const makeSyncMsgUpdateTrip = (tripPlanID: string, ops: Array<JSONPatchOp>): TripSync.Message => {
  return {
    tripPlanID,
    opType: SyncOpUpdateTrip,
    syncDataUpdateTrip: { ops },
  };
}