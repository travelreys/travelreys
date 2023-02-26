import _filter from 'lodash/filter';
import { makeCommonAxios } from './common';
import airports from '../data/airports.json';

const search = (
  origin: string,
  destination: string,
  departDate: string,
  returnDate: string | undefined,
  cabinClass: string
) => {
  return makeCommonAxios().get(`/api/v1/flights/search`, {
    params: {
      adults: '1',
      currency: 'SGD',
      origin: origin.toUpperCase(),
      destination: destination.toUpperCase(),
      departDate,
      returnDate,
      cabinClass,
    }
  });
}


const searchAirport = (q: string) => {
  return _filter(airports, (a) => a.airport.includes(q.toLowerCase()))
    .slice(0, 10);
}

export default {
  search,
  searchAirport,
};
