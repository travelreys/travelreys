import axios from 'axios';
import _find from "lodash/find";
import _get from "lodash/get";

import { BASE_URL } from './common';

const TripsAPI = {
  createTrip: (name: string, startDate: Date | undefined, endDate: Date | undefined) => {
    const url = `${BASE_URL}/api/v1/trips`;
    return axios.post(url, {name, startDate, endDate})
  },

  readTrips: () => {
    const url = `${BASE_URL}/api/v1/trips`;
    return axios.get(url).then(res => res.data);
  },

  readTrip: (id: string | undefined) => {
    const url = `${BASE_URL}/api/v1/trips/${id}`;
    return axios.get(url).then(res => res.data);
  },
};

export default TripsAPI;
