import axios from 'axios';
import _find from "lodash/find";
import _get from "lodash/get";

import { BASE_URL, makeCommonAxios } from './common';

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

const TripsAPI = {
  create: ( name: string, startDate?: Date, endDate?: Date): Promise<CreateResponse> => {
    const ax = makeCommonAxios();
    return ax.post(tripsPathPrefix, {name, startDate, endDate})
      .then((res) => {
        const id = _get(res, 'data.trip.id', "");
        return {id: id}
      })
      .catch((err) => {
        return {id: "", error: err.message};
      });
  },

  read: (id?: string): Promise<ReadResponse> => {
    const ax = makeCommonAxios();
    return ax.get(`${tripsPathPrefix}/${id}`, { params: { withMembers: "true" } })
      .then(res => {
        return {
          trip: _get(res, "data.trip", {}),
          members: _get(res, "data.members", {})
        }
      })
      .catch((err) => {
        return {trip: {}, error: err.message}
      });
  },

  readMembers: (id?: string): Promise<ReadMembersResponse> => {
    const ax = makeCommonAxios();
    return ax.get(`${tripsPathPrefix}/${id}/members`)
      .then(res => {
        return {members: _get(res, "data.members", {})}
      })
      .catch((err) => {
        return {trip: {}, error: err.message}
      });
  },

  list: (): Promise<ReadsResponse> => {
    const ax = makeCommonAxios();
    return ax.get(tripsPathPrefix)
      .then(res => {
        const trips = _get(res, "data.trips", []);
        return {trips}
      })
      .catch((err) => {
        return {trips: [], error: err.message}
      });
  },
};

export default TripsAPI;
