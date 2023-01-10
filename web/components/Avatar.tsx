import React, { FC } from 'react';
import _get from 'lodash/get';
import { usePopperTooltip } from 'react-popper-tooltip';

import { stringToColor } from '../utils/strings';

interface AvatarProps {
  name: string
  placement: string
}

const Avatar: FC<AvatarProps> = (props: AvatarProps) => {
  const {
    getTooltipProps,
    setTooltipRef,
    setTriggerRef,
    visible,
  } = usePopperTooltip({placement: props.placement as any});

  const style = {
    backgroundColor: stringToColor(props.name),
  }

  return (
    <>
      <div
        ref={setTriggerRef}
        style={style}
        className="relative inline-flex items-center justify-center w-7 h-7 overflow-hidden rounded-full"
      >
        <span className="font-medium text-white uppercase">{_get(props.name, "0")}</span>
      </div>
      {visible && (
        <div
          ref={setTooltipRef}
          {...getTooltipProps()}
          className="absolute z-10 inline-block px-3 py-2 text-sm font-medium text-white transition-opacity duration-300 bg-gray-900 rounded-lg shadow-sm  tooltip dark:bg-gray-700"
        >
          {props.name}
        </div>
      )}
    </>
  );
}

export default Avatar;
