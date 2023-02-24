import React, { FC } from 'react';
import _get from 'lodash/get';
import { usePopperTooltip } from 'react-popper-tooltip';

import { stringToColor } from '../../lib/strings';
import { AvatarCss } from '../../assets/styles/global';

interface AvatarProps {
  name: string
  imgUrl?: string
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
    return (
      <img className="h-full w-full rounded-full"
        src={props.imgUrl}
        alt="profile"
        referrerPolicy="no-referrer"
      />
    );
  }

  const renderName = () => {
    return (
      <span className={AvatarCss.InitialsTxt}>
        {_get(props.name, "0")}
      </span>
    );
  }

  return (
    <>
      <div ref={setTriggerRef} style={style} className={AvatarCss.Ctn}>
        {props.imgUrl ? renderImg() : renderName()}
      </div>
      {visible && (
        <div
          ref={setTooltipRef}
          className={AvatarCss.Tooltip}
          {...getTooltipProps()}
        >
          {props.name}
        </div>
      )}
    </>
  );
}

export default Avatar;
