
export const makeHotelPin = (name: string) => {
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group">
      <span class="absolute right-3.5 top-2.5 text-base font-bold pointer-events-none">
        <svg class="h-5 w-5 stroke-black fill-white stroke-2" fill="none" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            strokeWidth="2"
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
          />
        </svg>
      </span>
      <svg
        class="h-12 w-12 stroke-white stroke-2 fill-indigo-300 hover:fill-indigo-600"
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

export const makeActivityPin = (name: string) => {
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group">
      <svg
        class="h-12 w-12 stroke-white stroke-2 fill-orange-300 hover:fill-indigo-600"
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

