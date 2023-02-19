import axios from 'axios';
import _find from "lodash/find";
import _get from "lodash/get";

import { BASE_URL, makeCommonAxios } from './common';

export interface CreateTripResponse {
  id: string
  error?: string
}

export interface ReadTripResponse {
  tripPlan: any
  users?: any
  error?: string
}

export interface ReadTripsResponse {
  tripPlans: Array<any>
  error?: string
}


const tripsPathPrefix = "/api/v1/trips";

const TripsAPI = {
  createTrip: (
    name: string,
    startDate: Date | undefined,
    endDate: Date | undefined
  ): Promise<CreateTripResponse> => {
    const ax = makeCommonAxios();
    return ax.post(tripsPathPrefix, {name, startDate, endDate})
      .then((res) => {
        const id = _get(res, 'data.tripPlan.id', "");
        return {id: id}
      })
      .catch((err) => {
        return {id: "", error: err.message};
      });
  },

  readTrip: (id: string | undefined): Promise<ReadTripResponse> => {
    const ax = makeCommonAxios();
    return ax.get(`${tripsPathPrefix}/${id}`, { params: { withUsers: "true" } })
      .then(res => {
        const tripPlan = _get(res, "data.tripPlan", {});
        const users = _get(res, "data.users", {});
        return {tripPlan, users}
      })
      .catch((err) => {
        return {tripPlan: {}, error: err.message}
      });
  },

  listTrips: (): Promise<ReadTripsResponse> => {
    const ax = makeCommonAxios();
    return ax.get(tripsPathPrefix)
      .then(res => {
        const tripPlans = _get(res, "data.tripPlans", []);
        return {tripPlans}
      })
      .catch((err) => {
        return {tripPlans: [], error: err.message}
      });
  },
};

export default TripsAPI;
