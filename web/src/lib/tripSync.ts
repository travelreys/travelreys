import { Op } from "./jsonpatch";

export const OpJoinSession = "OpJoinSession";
export const OpLeaveSession = "OpLeaveSession";
export const OpPingSession = "OpPingSession";
export const OpMemberUpdate = "OpMemberUpdate";
export const OpUpdateTrip = "OpUpdateTrip";

export const UpdateTitleAddNewMember = "AddNewMember";

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
  ID: string
}

interface MsgDataLeaveSession {
  ID: string
}

interface MsgDataMemberUpdate {
  members: any
}

interface MsgDataPing {}

interface MsgDataUpdateTrip {
  title: string
  ops: Array<Op>
}

export const makeMsgJoinSession = (tripID: string, memberID: string): Message => {
  return {
    tripID,
    op: OpJoinSession,
    data: { joinSession: {ID: memberID}}
  }
}

export const makeMsgLeaveSession = (tripID: string, memberID: string): Message => {
  return {
    tripID,
    op: OpLeaveSession,
    data: { leaveSession: {ID: memberID}}
  }
}

export const makeMsgUpdateTrip = (tripID: string, title: string, ops: Array<Op>): Message => {
  return {
    tripID,
    op: OpUpdateTrip,
    data: { updateTrip: {title, ops}}
  };
}