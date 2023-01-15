import React, { FC } from 'react';

import { ExclamationCircleIcon } from '@heroicons/react/24/solid'

interface AlertProps {
  status: string
  title: string
  message: string
}

const Alert: FC<AlertProps> = (props: AlertProps) => {
  return (
    <div className="bg-red-100 mb-4 border-t-4 border-red-500 rounded-b text-red-900 px-4 py-3 shadow-md" role="alert">
      <div className="flex">
        <div className="mr-2">
          <ExclamationCircleIcon className='inline align-bottom h-6 w-6'/>
        </div>
        <div>
          <p className="font-bold">{props.title}</p>
          <p className="text-sm">{props.message}</p>
        </div>
      </div>
    </div>
  );
}

export default Alert;
