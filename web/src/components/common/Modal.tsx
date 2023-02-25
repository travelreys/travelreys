import React, { FC } from 'react';

export const css = {
  container: "relative z-20",
  inset: "fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity",
  content: "fixed inset-0 z-10 overflow-y-auto",
  contentContainer: "flex min-h-full flex-col p-4 text-center sm:items-center sm:p-0",
  contentCard: "bg-white relative transform rounded-lg text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-2xl",
}

interface ModalProps {
  isOpen: boolean
  children: React.ReactNode
}

const Modal: FC<ModalProps> = (props: ModalProps) => {
  if (!props.isOpen) {
    return <></>;
  }
  return (
    <div className={css.container}>
      <div className={css.inset}></div>
      <div className={css.content}>
        <div className={css.contentContainer}>
          <div className={css.contentCard}>
            {props.children}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Modal;
