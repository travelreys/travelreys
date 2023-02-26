import _get from 'lodash/get';
import { makeCommonAxios } from './common';

const search = (query: string) => {
  return makeCommonAxios().get("/api/v1/images/search", { params: { query } });
}

const makeUserURL = (username: string) => {
  return `https://unsplash.com/@${username}?utm_source=tiinyplanet&utm_medium=referral`;
}

const makeSrc = (image: any) => {
  return _get(image, "urls.full");
}

const makeSrcSet = (image: any) => {
  return [
    `${_get(image, "urls.thumb")} 200w`,
    `${_get(image, "urls.small")} 400w`,
    `${_get(image, "urls.regular")} 1080w`
  ].join(", ")
}

export default {
  search,
  makeUserURL,
  makeSrcSet,
  makeSrc,
};
