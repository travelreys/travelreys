import axios from 'axios';
import _get from 'lodash/get';

import { BASE_URL } from './common';

const MapsAPI = {
  placeAutocomplete: (query: string, types: Array<string>, sessiontoken: string) => {
    const url = `${BASE_URL}/api/v1/maps/place/autocomplete`;
    const typesParam = types.join(",");
    return axios.get(url, {
      params: { query, types: typesParam, sessiontoken }
    });
  },
  placeDetails: (placeID: string, fields: Array<string>, sessiontoken: string | undefined) => {
    const url = `${BASE_URL}/api/v1/maps/place/details`;
    const fieldsParam = fields.join(",");
    return axios.get(url, {
      params: { placeID, fields: fieldsParam, sessiontoken }
    });
  }
};

export default MapsAPI;

export const EMBED_MAPS_APIKEY = "AIzaSyBaqenQ0nQVtkhnXBn-oWBtlPDL5uHmvNU";
