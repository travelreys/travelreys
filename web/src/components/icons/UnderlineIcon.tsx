import { IconProps } from "./common";

export default function UnderlineIcon(props: IconProps) {
  return (
    <svg
      className={props.className}
      width="24"
      height="24"
      viewBox="0 0 24 24"
      stroke-width="2"
      stroke="currentColor"
      fill="none"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path stroke="none" d="M0 0h24v24H0z"/>  <line x1="6" y1="20" x2="18" y2="20" />  <path d="M8 5v6a4 4 0 0 0 8 0v-6" />
    </svg>
  )
}
