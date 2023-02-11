import { IconProps } from "./common";

export default function CupIconFill(props: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className={props.className}
      viewBox="0 0 48 48"
    >
      <path d="M8 42v-3h31.95v3Zm7.55-6q-3.15 0-5.35-2.175Q8 31.65 8 28.5V6h33q1.25 0 2.125.875T44 9v8q0 1.25-.875 2.125T41 20h-4.8v8.5q0 3.15-2.2 5.325Q31.8 36 28.65 36ZM36.2 17H41V9h-4.8Z"/>
    </svg>
  )
}

export const svgWithStyle = (style: string) => `
  <svg
    xmlns="http://www.w3.org/2000/svg"
    style=${style}
    viewBox="0 0 48 48"
  >
    <path d="M8 42v-3h31.95v3Zm7.55-6q-3.15 0-5.35-2.175Q8 31.65 8 28.5V6h33q1.25 0 2.125.875T44 9v8q0 1.25-.875 2.125T41 20h-4.8v8.5q0 3.15-2.2 5.325Q31.8 36 28.65 36ZM36.2 17H41V9h-4.8Z"/>
  </svg>
`;
