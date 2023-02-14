import { IconProps } from "./common";

export default function NatureIconFill(props: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 48 48"
      className={props.className}
    >
      <path d="M11.5 44v-3h11v-9H18q-4.15 0-7.075-2.925T8 22q0-3 1.65-5.525Q11.3 13.95 14.1 12.8q.45-3.75 3.275-6.275Q20.2 4 24 4t6.625 2.525Q33.45 9.05 33.9 12.8q2.8 1.15 4.45 3.675Q40 19 40 22q0 4.15-2.925 7.075T30 32h-4.5v9H37v3Z"/>
    </svg>
  )
}

export const svgWithStyle = (style: string) => `
  <svg
    xmlns="http://www.w3.org/2000/svg"
    class="h-5 w-5 fill-white stroke-2"
    viewBox="0 0 48 48"
  >
  <path d="M11.5 44v-3h11v-9H18q-4.15 0-7.075-2.925T8 22q0-3 1.65-5.525Q11.3 13.95 14.1 12.8q.45-3.75 3.275-6.275Q20.2 4 24 4t6.625 2.525Q33.45 9.05 33.9 12.8q2.8 1.15 4.45 3.675Q40 19 40 22q0 4.15-2.925 7.075T30 32h-4.5v9H37v3Z"/>
  </svg>
`;
