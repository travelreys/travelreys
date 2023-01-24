import axios from 'axios';
import _get from 'lodash/get';
import _filter from 'lodash/filter';

import { BASE_URL } from './common';

import airports from '../public/airports.json';

const FlightsAPI = {
  search: (
    origIATA: string,
    destIATA: string,
    departDate: string,
    returnDate: string | undefined,
    cabinClass: string
  ) => {
    const url = `${BASE_URL}/api/v1/flights/search`;
    return axios.get(url, {
      params: {
        numAdults: '1',
        currency: 'SGD',
        origIATA: origIATA.toUpperCase(),
        destIATA: destIATA.toUpperCase(),
        departDate,
        returnDate,
        cabinClass,
       }
    });
  },
  airportAutocomplete: (q: string) => {
    return _filter(airports, (a) => a.airport.includes(q.toLowerCase())).slice(0, 10);
  }
};

export default FlightsAPI;
