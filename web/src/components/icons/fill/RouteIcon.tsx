import { IconProps } from "../common";

export default function RouteIcon(props: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className={props.className}
      height="48"
      viewBox="0 96 960 960"
    >
      <path d="M355 936q-65 0-110-45.5T200 781V432q-35-13-57.5-41.5T120 326q0-46 32.5-78t77.5-32q46 0 78 32t32 78q0 36-22.5 64.5T260 432v349q0 39 27.5 67t67.5 28q41 0 68-28t27-67V371q0-65 45-110t110-45q65 0 110 45t45 110v349q35 13 57.5 41.5T840 826q0 45-32 77.5T730 936q-45 0-77.5-32.5T620 826q0-36 22.5-65t57.5-41V371q0-40-27.5-67.5T605 276q-40 0-67.5 27.5T510 371v410q0 64-45 109.5T355 936Z"/>
    </svg>
  )
}
