import { IconProps } from "./common";

export default function FlightIconFill(props: IconProps) {
  return (
    <svg
      className={props.className}
      viewBox="0 0 48 48"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M6.05 42v-3h36v3ZM9.2 31.6l-5.15-8.75 2.15-.4 3.5 3.1L21 22.5 12.45 8.15l2.9-.85L29.6 20.15l10.8-2.9q1.35-.4 2.45.475t1.1 2.325q0 .95-.575 1.7t-1.475 1Z"/>
    </svg>
  )
}

export const svgWithStyle = (style: string) => `
  <svg
    xmlns="http://www.w3.org/2000/svg"
    style=${style}
    viewBox="0 0 48 48"
  >
  <path d="M6.05 42v-3h36v3ZM9.2 31.6l-5.15-8.75 2.15-.4 3.5 3.1L21 22.5 12.45 8.15l2.9-.85L29.6 20.15l10.8-2.9q1.35-.4 2.45.475t1.1 2.325q0 .95-.575 1.7t-1.475 1Z"/>
  </svg>
`;
