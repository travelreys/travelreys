import _get from "lodash/get";

import { makeCommonAxios } from './common';

export interface CreateResponse {
  id: string
  error?: string
}

export interface ReadResponse {
  trip: any
  members?: any
  error?: string
}

export interface ReadMembersResponse {
  members?: any
  error?: string
}

export interface ReadsResponse {
  trips: Array<any>
  error?: string
}

const tripsPathPrefix = "/api/v1/trips";

const create = (name: string, startDate?: Date, endDate?: Date): Promise<CreateResponse> => {
  return makeCommonAxios().post(tripsPathPrefix, {name, startDate, endDate})
    .then((res) => {
      const id = _get(res, 'data.trip.id', "");
      return {id: id}
    })
    .catch((err) => {
      return {id: "", error: err.message};
    });
}

const read = (id?: string): Promise<ReadResponse> => {
  return makeCommonAxios().get(`${tripsPathPrefix}/${id}`, { params: { withMembers: "true" } })
    .then(res => {
      return {
        trip: _get(res, "data.trip", {}),
        members: _get(res, "data.members", {})
      }
    })
    .catch((err) => {
      return {trip: {}, error: err.message}
    });
}

const readMembers = (id?: string): Promise<ReadMembersResponse> => {
  return makeCommonAxios().get(`${tripsPathPrefix}/${id}/members`)
    .then(res => {
      return {members: _get(res, "data.members", {})}
    })
    .catch((err) => {
      return {trip: {}, error: err.message}
    });
}

const list = (): Promise<ReadsResponse> => {
  return makeCommonAxios().get(tripsPathPrefix)
    .then(res => {
      const trips = _get(res, "data.trips", []);
      return {trips}
    })
    .catch((err) => {
      return {trips: [], error: err.message}
    });
}

export default {
  create,
  read,
  readMembers,
  list,
};
