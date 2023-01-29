
export const makeHotelPin = (name: string) => {
  const pin = document.createElement("template");
  const template = `
    <div class="absolute cursor-pointer max-h-12 top-0 left-0 -translate-y-full -translate-x-1/2 group">
      <svg
        class="h-12 w-12 stroke-white stroke-2 fill-indigo-300 hover:fill-indigo-600"
        viewBox="0 0 24 24"
      >
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
      </svg>
      <div class="absolute w-fit hidden whitespace-nowrap text-white text-center text-sm w-16 m-2 font-bold font-medium rounded-lg p-1 group-hover:inline-flex group-hover:bg-black">
        ${name}
      </div>
    </div>
  `.trim();
  pin.innerHTML = template;
  return pin.content.firstChild;
}

