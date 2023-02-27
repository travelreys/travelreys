import { Op } from "./jsonpatch";

export const OpJoinSession = "OpJoinSession";
export const OpLeaveSession = "OpLeaveSession";
export const OpPingSession = "OpPingSession";
export const OpUpdateTrip = "OpUpdateTrip";

export const MsgUpdateTripTitleAddNewMember = "AddNewMember";
export const MsgUpdateTripTitleReorderItinerary = "ReorderItinerary";

export interface Message {
  connID?: string
  tripID: string
  op: "OpJoinSession" | "OpLeaveSession" | "OpPingSession" | "OpUpdateTrip"
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
  id: string
  members: any
}

interface MsgDataLeaveSession {
  id: string
  members: any
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
    data: { joinSession: {id: memberID, members: []}}
  }
}

export const makeMsgLeaveSession = (tripID: string, memberID: string): Message => {
  return {
    tripID,
    op: OpLeaveSession,
    data: { leaveSession: {id: memberID, members: []}}
  }
}

export const makeMsgUpdateTrip = (tripID: string, title: string, ops: Array<Op>): Message => {
  return {
    tripID,
    op: OpUpdateTrip,
    data: { updateTrip: {title, ops}}
  };
}