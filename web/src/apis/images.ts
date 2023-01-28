import axios from 'axios';
import _get from 'lodash/get';

import { BASE_URL } from './common';

const ImagesAPI = {
  search: (query: string) => {
    const url = `${BASE_URL}/api/v1/images/search`;
    return axios.get(url, {
      params: { query }
    });
  },
  makeUserReferURL: (username: string) => {
    return `https://unsplash.com/@${username}?utm_source=tiinyplanet&utm_medium=referral`;
  },
  makeSrcSet: (image: any) => {
    return `${_get(image, "urls.thumb")} 200w, ${_get(image, "urls.small")} 400w, ${_get(image, "urls.regular")} 1080w`
  },
  makeSrc: (image: any) => {
    return _get(image, "urls.full");
  },
};

export default ImagesAPI;

export const stockImageSrc = "https://images.unsplash.com/photo-1570913149827-d2ac84ab3f9a?crop=entropy&cs=tinysrgb&fm=jpg&ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHxhcHBsZXxlbnwwfHx8fDE2NzMzNTc4Njc&ixlib=rb-4.0.3&q=80"
