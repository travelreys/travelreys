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
  users?: any
  error?: string
}

export interface ReadsResponse {
  trips: Array<any>
  error?: string
}


const tripsPathPrefix = "/api/v1/trips";

const TripsAPI = {
  createTrip: (
    name: string,
    startDate: Date | undefined,
    endDate: Date | undefined
  ): Promise<CreateResponse> => {
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

  readTrip: (id: string | undefined): Promise<ReadResponse> => {
    const ax = makeCommonAxios();
    return ax.get(`${tripsPathPrefix}/${id}`, { params: { withUsers: "true" } })
      .then(res => {
        const trip = _get(res, "data.trip", {});
        const users = _get(res, "data.users", {});
        return {trip, users}
      })
      .catch((err) => {
        return {trip: {}, error: err.message}
      });
  },

  listTrips: (): Promise<ReadsResponse> => {
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
