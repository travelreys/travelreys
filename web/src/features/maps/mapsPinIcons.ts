import _get from "lodash/get";

import { svgWithStyle as svgWithStyleCameraIcon } from '../../components/icons/fill/CameraIcon';
import { svgWithStyle as svgWithStyleCupIcon } from '../../components/icons/fill/CupIcon';
import { svgWithStyle as svgWithStyleDiningIcon } from '../../components/icons/fill/DiningIcon';
import { svgWithStyle as svgWithStyleMapPinIcon } from '../../components/icons/fill/MapPinIcon';
import { svgWithStyle as svgWithStyleNatureIcon } from '../../components/icons/fill/NatureIcon';
import { svgWithStyle as svgWithStyleShoppingIcon } from '../../components/icons/fill/ShoppingIcon';
import { svgWithStyle as svgWithStyleHotelIcon } from '../../components/icons/fill/HotelIcon';

const iconSvgMap = {
  "camera": svgWithStyleCameraIcon,
  "coffee": svgWithStyleCupIcon,
  "forkspoon": svgWithStyleDiningIcon,
  "pin": svgWithStyleMapPinIcon,
  "nature": svgWithStyleNatureIcon,
  "shopping": svgWithStyleShoppingIcon,
  "hotel": svgWithStyleHotelIcon,
}


const makePinTooltip = (name: string) => {
  return `
  <div class="absolute w-fit hidden whitespace-nowrap text-white text-center text-sm m-2 font-bold font-medium rounded-lg p-1 group-hover:inline-flex group-hover:bg-black">
    ${name}
  </div>
  `
}

const makeIcon = (icon: string) => {
  const iconStyle = `"fill:white;stroke-width:2;height:1.25rem;width:1.25rem;"`;
  let iconSvg = "";
  const svgFn = _get(iconSvgMap, icon);
  if (svgFn !== undefined) {
    iconSvg = `
    <span class="absolute right-3.5 top-2.5 text-base font-bold pointer-events-none">
      ${svgFn(iconStyle)}
    </span>
    `
  }
  return iconSvg;
}

export const makePinWithTooltip = (name: string, color: string, icon: string) => {
  const pinStyle =`"fill:${color}"`;
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group hover:z-50">
      ${makeIcon(icon)}
      <svg
        class="h-12 w-12 stroke-white stroke-2"
        style=${pinStyle}
        viewBox="0 0 24 24"
      >
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
      </svg>
      ${makePinTooltip(name)}
    </div>
  `.trim();
  pin.innerHTML = template;
  return pin.content.firstChild;
}

const makeNumber = (num: string) => {
  let right = "right-5 text-base"
  if (num.length === 2) {
    right = "right-4 text-base"
  }
  if (num.length === 3) {
    right = "right-2 text-sm"
  }

  return `
    <span class="absolute ${right} text-sm top-2.5  font-bold text-white pointer-events-none">
      ${num}
    </span>
  `
}

export const makeNumberPin = (name: string, color: string, number: string) => {
  const pinStyle =`"fill:${color}"`;
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group hover:z-50">
      ${makeNumber(number)}
      <svg
        class="h-12 w-12 stroke-white stroke-2"
        style=${pinStyle}
        viewBox="0 0 24 24"
      >
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
      </svg>
      ${makePinTooltip(name)}
    </div>
  `.trim();
  pin.innerHTML = template;
  return pin.content.firstChild;
}
