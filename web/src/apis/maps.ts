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
export const placeFields = [
  "address_component",
  "adr_address",
  "business_status",
  "formatted_address",
  "geometry",
  "name",
  "photos",
  "place_id",
  "types",
  "utc_offset",
  "opening_hours",
  "formatted_phone_number",
  "international_phone_number",
  "website",
];
export const EMBED_MAPS_APIKEY = "AIzaSyBaqenQ0nQVtkhnXBn-oWBtlPDL5uHmvNU";
export const PLACE_IMAGE_APIKEY = "AIzaSyBgNwirAT6TSS208emcC0Lbgex6i3EwhR0";