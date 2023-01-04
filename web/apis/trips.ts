import axios from 'axios'
import useSWR from 'swr';

import { BASE_URL } from './common';

const TripsAPI = {
  createTrip: (name: string, startDate: Date | undefined, endDate: Date | undefined) => {
    const url = `${BASE_URL}/api/v1/trips`;
    return axios.post(url, {name, startDate, endDate})
  },
  readTrips: () => {
    const url = `${BASE_URL}/api/v1/trips`;
    const fetcher = (url: string) => {
      return axios.get(url).then(res => res.data);
    }

    const { data, error, isLoading } = useSWR(url, fetcher);
    return { data, error, isLoading};
  },
  readTrip: (id: string | undefined) => {
    const url = `${BASE_URL}/api/v1/trips/${id}`;
    const fetcher = (url: string) => {
      return axios.get(url).then(res => res.data);
    }

    const { data, error, isLoading } = useSWR(id ? url : null, fetcher);
    return { data, error, isLoading};
  }
};

export default TripsAPI;
