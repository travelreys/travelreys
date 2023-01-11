import _get from "lodash/get";

export const makeUserReferURL = (username: string) => {
  return `https://unsplash.com/@${username}?utm_source=tiinyplanet&utm_medium=referral`;
}

export const makeSrcSet = (image: any) => {
  return `${_get(image, "urls.thumb")} 200w, ${_get(image, "urls.small")} 400w, ${_get(image, "urls.regular")} 1080w`
}

export const makeSrc = (image: any) => {
  return _get(image, "urls.full");
}