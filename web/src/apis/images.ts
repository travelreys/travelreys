import axios from 'axios';
import _get from 'lodash/get';

import { BASE_URL } from './common';

const ImagesAPI = {
  search: (query: string) => {
    const url = `${BASE_URL}/api/v1/images/search`;
    return axios.get(url, { params: { query } });
  },

  makeUserURL: (username: string) => {
    return `https://unsplash.com/@${username}?utm_source=tiinyplanet&utm_medium=referral`;
  },

  makeSrcSet: (image: any) => {
    return [
      `${_get(image, "urls.thumb")} 200w`,
      `${_get(image, "urls.small")} 400w`,
      `${_get(image, "urls.regular")} 1080w`
    ].join(", ")
  },

  makeSrc: (image: any) => {
    return _get(image, "urls.full");
  },
};

export default ImagesAPI;
