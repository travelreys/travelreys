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
  },
  directions: (originPlaceID: string, destPlaceID: string, mode: string) => {
    const url = `${BASE_URL}/api/v1/maps/place/directions`;
    return axios.get(url, {
      params: { originPlaceID, destPlaceID, mode }
    });
  },
  optimizeRoute: (originPlaceID: string, destPlaceID: string, waypointsPlaceID: string) => {
    const url = `${BASE_URL}/api/v1/maps/place/optimize-route`;
    return axios.get(url, {
      params: { originPlaceID, destPlaceID, waypointsPlaceID }
    });
  },

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
  "url",
];
export const placeAtmosphereFields = placeFields.concat([
  "editorial_summary",
  "price_level",
  "rating",
  "reviews",
  "user_ratings_total"
]);
export const ModeDriving = "driving";

export const EMBED_MAPS_APIKEY = "AIzaSyBaqenQ0nQVtkhnXBn-oWBtlPDL5uHmvNU";
export const PLACE_IMAGE_APIKEY = "AIzaSyBgNwirAT6TSS208emcC0Lbgex6i3EwhR0";
