import axios from 'axios'
import useSWR from 'swr';

import { BASE_URL } from './common';

const TripsAPI = {
  readTrips: () => {
    const fetcher = (url: string) => axios.get(url).then(res => res.data);
    const url = `${BASE_URL}/api/v1/trips`;

    const { data, error, isLoading } = useSWR(url, fetcher);
    return { trips: data, isLoading, error};
  }
};

export default TripsAPI;
