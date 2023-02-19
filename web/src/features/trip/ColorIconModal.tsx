import React, {
  FC,
  useState,
} from 'react';
import {
  CheckIcon,
  XMarkIcon
} from '@heroicons/react/24/outline';

import Modal from '../../components/common/Modal';
import { ContentIconOpts } from '../../lib/trips';
import { CommonCss } from '../../assets/styles/global';

interface IconProps {
  icon: any
  selected?: boolean
  onClick: () => void
}

const Icon: FC<IconProps> = (props: IconProps) => {
  let css = 'p-2 rounded-full mr-2 inline-flex items-center justify-around'
  css = css + (props.selected ? " bg-indigo-200" : " bg-slate-200")

  return (
    <button type="button" className={css} onClick={props.onClick}>
      <props.icon className={CommonCss.Icon} />
    </button>
  )
}

interface ColorIconModalProps {
  isOpen: boolean
  defaultSelectedColor?: string
  defaultSelectedIcon?: string
  colors: Array<string>
  icons: Array<string>

  onClose: () => void
  onSubmit: (color: string | undefined, icon: string | undefined) => void
}

const ColorIconModal: FC<ColorIconModalProps> = (props: ColorIconModalProps) => {

  const [selectedColor, setSelectedColor] = useState(props.defaultSelectedColor);
  const [selectedIcon, setSelectedIcon] = useState(props.defaultSelectedIcon);

  // Event Handlers
  const submitBtnOnClick = () => {
    props.onSubmit(selectedColor, selectedIcon);
    props.onClose();
  }


  // Renderers
  const renderHeader = () => {
    return (
      <div className='flex justify-between items-center mb-4'>
        <div className='text-gray-800 font-bold text-lg'>
          Choose color and icon
        </div>
        <button
          type="button"
          onClick={() => {props.onClose()}}
        >
          <XMarkIcon className={CommonCss.Icon} />
        </button>
      </div>
    );
  }

  const renderColorOpts = () => {
    return (
      <div>
        <div className='text-gray-800 font-bold text-sm'>Choose Color</div>
        <div className='w-full p-2 pl-0 mb-2'>
          {props.colors.map((c) => (
            <button
              key={c}
              className='rounded-full mr-2 inline-flex items-center justify-around'
              onClick={() =>  {setSelectedColor(c)}}
              style={{ backgroundColor: c }}
            >
              {selectedColor === c
                ? <span className='p-2'>
                    <CheckIcon className="h-4 w-4 stroke-white stroke-[4]" />
                  </span>
                : <span className='p-4'/>}
            </button>
          ))}
        </div>
      </div>
    );
  }

  const renderIconOpts = () => {
    return (
      <div>
        <div className='text-gray-800 font-bold text-sm'>Choose Icon</div>
        <div className='w-full p-2 pl-0 mb-4'>
          {props.icons.map((icon: string) => (
            <Icon
              key={icon}
              icon={ContentIconOpts[icon]}
              selected={selectedIcon === icon}
              onClick={() => {setSelectedIcon(icon)}}
            />
          ))}
        </div>
      </div>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className='p-5'>
        {renderHeader()}
        {renderColorOpts()}
        {renderIconOpts()}
        <button
          type="button"
          className='bg-indigo-500 px-4 py-2 rounded-lg font-bold text-sm text-white'
          onClick={submitBtnOnClick}
        >
          Submit
        </button>
      </div>
    </Modal>
  );
}


export default ColorIconModal;




