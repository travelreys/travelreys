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
