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

export const OpJoinSession = "OpJoinSession";
export const OpLeaveSession = "OpLeaveSession";
export const OpPingSession = "OpPingSession";
export const OpMemberUpdate = "OpMemberUpdate";
export const OpUpdateTrip = "OpUpdateTrip";

export const UpdateTitleAddNewMember = "AddNewMember";


export namespace TripSync {
  export interface Message {
    connID?: string
    tripID: string
    op: "OpJoinSession" | "OpLeaveSession" | "OpPingSession"| "OpMemberUpdate" |"OpUpdateTrip"
    counter?: number
    data: MessageData
  }

  interface MessageData {
    joinSession?: MsgDataJoinSession
    leaveSession?: MsgDataLeaveSession
    memberUpdate?: MsgDataMemberUpdate
    ping?: MsgDataPing
    updateTrip?: MsgDataUpdateTrip
  }

  interface MsgDataJoinSession {
    memberID: string
  }

  interface MsgDataLeaveSession {
    memberID: string
  }

  interface MsgDataMemberUpdate {
    members: any
  }

  interface MsgDataPing {}

  interface MsgDataUpdateTrip {
    title: string
    ops: Array<JSONPatchOp>
  }
};

export const makeMsgJoinSession = (tripID: string, memberID: string): TripSync.Message => {
  return {
    tripID,
    op: OpJoinSession,
    data: { joinSession: {memberID}}
  }
}

export const makeMsgLeaveSession = (tripID: string, memberID: string): TripSync.Message => {
  return {
    tripID,
    op: OpLeaveSession,
    data: { leaveSession: {memberID}}
  }
}

export const makeMsgUpdateTrip = (tripID: string, title: string, ops: Array<JSONPatchOp>): TripSync.Message => {
  return {
    tripID,
    op: OpUpdateTrip,
    data: { updateTrip: {title, ops}}
  };
}
