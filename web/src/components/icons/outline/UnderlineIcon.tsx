import { IconProps } from "../common";

export default function UnderlineIcon(props: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className={props.className}
      viewBox="0 96 960 960"
    >
      <path d="M200 916v-60h560v60H200Zm280-140q-100 0-156.5-58.5T267 559V216h83v343q0 63 34 101t96 38q62 0 96-38t34-101V216h83v343q0 100-56.5 158.5T480 776Z"/>
    </svg>
  )
}