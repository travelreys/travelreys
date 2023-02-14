import _get from "lodash/get";

import { svgWithStyle as svgWithStyleCameraIconFill } from '../icons/CameraIconFill';
import { svgWithStyle as svgWithStyleCupIconFill } from '../icons/CupIconFill';
import { svgWithStyle as svgWithStyleDiningIconFill } from '../icons/DiningIconFill';
import { svgWithStyle as svgWithStyleMapPinIconFill } from '../icons/MapPinIconFill';
import { svgWithStyle as svgWithStyleNatureIconFill } from '../icons/NatureIconFill';
import { svgWithStyle as svgWithStyleShoppingIconFill } from '../icons/ShoppingIconFill';
import { svgWithStyle as svgWithStyleHotelIcon } from '../icons/HotelIcon';

const iconSvgMap = {
  "camera": svgWithStyleCameraIconFill,
  "coffee": svgWithStyleCupIconFill,
  "forkspoon": svgWithStyleDiningIconFill,
  "pin": svgWithStyleMapPinIconFill,
  "nature": svgWithStyleNatureIconFill,
  "shopping": svgWithStyleShoppingIconFill,
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
