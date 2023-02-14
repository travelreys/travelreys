import React, { ChangeEvent, FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { ModalCss } from '../../styles/global';

// Modal

interface ModalProps {
  isOpen: boolean
  children: React.ReactNode
}

const Modal: FC<ModalProps> = (props: ModalProps) => {

  if (!props.isOpen) {
    return <></>;
  }

  return (
    <div className={ModalCss.Container}>
      <div className={ModalCss.Inset}></div>
      <div className={ModalCss.Content}>
        <div className={ModalCss.ContentContainer}>
          <div className={ModalCss.ContentCard}>
            {props.children}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Modal;
