import _get from "lodash/get";

export const stockImageSrc = "https://images.unsplash.com/photo-1570913149827-d2ac84ab3f9a?crop=entropy&cs=tinysrgb&fm=jpg&ixid=MnwzOTc1ODV8MHwxfHNlYXJjaHwxfHxhcHBsZXxlbnwwfHx8fDE2NzMzNTc4Njc&ixlib=rb-4.0.3&q=80"

export const makeUserReferURL = (username: string) => {
  return `https://unsplash.com/@${username}?utm_source=tiinyplanet&utm_medium=referral`;
}

export const makeSrcSet = (image: any) => {
  return `${_get(image, "urls.thumb")} 200w, ${_get(image, "urls.small")} 400w, ${_get(image, "urls.regular")} 1080w`
}

export const makeSrc = (image: any) => {
  return _get(image, "urls.full");
}
