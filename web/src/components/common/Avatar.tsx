import React, { FC } from 'react';
import _get from 'lodash/get';
import { usePopperTooltip } from 'react-popper-tooltip';
import { stringToColor } from '../../lib/strings';

export const css = {
  ctn: "relative inline-flex items-center justify-center h-full w-full overflow-hidden rounded-full",
  text: "text-white font-semibold text-xl uppercase",
  tooltip: "absolute z-10 inline-block p-2 text-sm font-medium text-white bg-gray-900 rounded-lg",
}

interface AvatarProps {
  name: string
  imgurl?: string
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

  // Renderers
  const renderImg = () => {
    if (props.imgurl) {
      return (
        <img className="h-full w-full rounded-full"
          src={props.imgurl}
          alt="profile"
          referrerPolicy="no-referrer"
        />
      );
    }
    return (
      <span className={css.text}>
        {_get(props.name, "0")}
      </span>
    );
  }

  return (
    <>
      <div ref={setTriggerRef} style={style} className={css.ctn}>
        {renderImg()}
      </div>
      {visible && (
        <div ref={setTooltipRef} className={css.tooltip} {...getTooltipProps()}>
          {props.name}
        </div>
      )}
    </>
  );
}

export default Avatar;
