import _get from "lodash/get";

import { svgWithStyle as svgWithStyleCameraIconFill } from '../icons/CameraIconFill';
import { svgWithStyle as svgWithStyleCupIconFill } from '../icons/CupIconFill';
import { svgWithStyle as svgWithStyleDiningIconFill } from '../icons/DiningIconFill';
import { svgWithStyle as svgWithStyleMapPinIconFill } from '../icons/MapPinIconFill';
import { svgWithStyle as svgWithStyleNatureIconFill } from '../icons/NatureIconFill';
import { svgWithStyle as svgWithStyleShoppingIconFill } from '../icons/ShoppingIconFill';

const iconSvgMap = {
  "camera": svgWithStyleCameraIconFill,
  "coffee": svgWithStyleCupIconFill,
  "forkspoon": svgWithStyleDiningIconFill,
  "nature": svgWithStyleMapPinIconFill,
  "pin": svgWithStyleNatureIconFill,
  "shopping": svgWithStyleShoppingIconFill,
}


export const makeHotelPin = (name: string) => {
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group hover:z-50">
      <span class="absolute right-3.5 top-2.5 text-base font-bold pointer-events-none">
        <svg
          class="h-5 w-5 stroke-black fill-white stroke-2"
          viewBox="0 0 48 48"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M2 38V8.75h3v19.7h17.65V13h16.1q3 0 5.125 2.125T46 20.25V38h-3v-6.55H5V38Zm11.5-12.45q-2.25 0-3.775-1.525T8.2 20.25q0-2.25 1.525-3.775T13.5 14.95q2.25 0 3.775 1.525T18.8 20.25q0 2.25-1.525 3.775T13.5 25.55Zm12.15 2.9H43v-8.2q0-1.75-1.25-3t-3-1.25h-13.1Zm-12.15-5.9q.95 0 1.625-.675t.675-1.625q0-.95-.675-1.625T13.5 17.95q-.95 0-1.625.675T11.2 20.25q0 .95.675 1.625t1.625.675Zm0 0q-.95 0-1.625-.675T11.2 20.25q0-.95.675-1.625t1.625-.675q.95 0 1.625.675t.675 1.625q0 .95-.675 1.625t-1.625.675ZM25.65 16h13.1q1.75 0 3 1.25t1.25 3v8.2H25.65Z"/>
        </svg>
      </span>
      <svg
        class="h-12 w-12 stroke-white stroke-2 fill-gray-800 hover:fill-black "
        viewBox="0 0 24 24"
      >
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
      </svg>
      <div class="absolute w-fit hidden whitespace-nowrap text-white text-center text-sm m-2 font-bold font-medium rounded-lg p-1 group-hover:inline-flex group-hover:bg-black">
        ${name}
      </div>
    </div>
  `.trim();
  pin.innerHTML = template;
  return pin.content.firstChild;
}


export const makeActivityPin = (name: string, color: string, icon: string) => {
  const iconStyle = `"fill:white;stroke:black;stroke-width:2;height:1.25rem;width:1.25rem;"`;
  let iconSvg = "";
  const svgFn = _get(iconSvgMap, icon);
  if (svgFn !== undefined) {
    iconSvg = `
    <span class="absolute right-3.5 top-2.5 text-base font-bold pointer-events-none">
      ${svgFn(iconStyle)}
    </span>
    `
  }

  const pinStyle =`"fill:${color}"`;
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group hover:z-50">
      ${iconSvg}
      <svg
        class="h-12 w-12 stroke-white stroke-2"
        style=${pinStyle}
        viewBox="0 0 24 24"
      >
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
      </svg>
      <div class="absolute w-fit hidden whitespace-nowrap text-white text-center text-sm m-2 font-bold font-medium rounded-lg p-1 group-hover:inline-flex group-hover:bg-black">
        ${name}
      </div>
    </div>
  `.trim();
  pin.innerHTML = template;
  return pin.content.firstChild;
}

