import React, { FC, useState, useEffect } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";

import { useUser } from '../../context/user-context';
import { LabelUserGoogleImage } from '../../lib/auth';

const ProfilePage: FC = () => {

  // UI State
  const { state } = useUser();
  const [isLoading, setIsLoading] = useState(false);

  // Event Handlers

  // Renderers
  const renderUserInfo = () => {
    const profileImgURL = _get(state.user, `labels.${LabelUserGoogleImage}`);
    const name = _get(state.user, "name");
    const email = _get(state.user, "email");
    return (
      <div className='flex items-center mb-4'>
        <div className='mr-4'>
          <img className="w-24 h-24 rounded-full"
            src={profileImgURL}
            alt="profile image"
            referrerPolicy="no-referrer"
          />
        </div>
        <div>
          <p className='font-bold text-2xl'>{name}</p>
          <p className='text-gray-500'>{email}</p>
        </div>
      </div>
    );
  }

  return (
    <div className='p-4 px-2'>
      {renderUserInfo()}
    </div>
  );
}


export default ProfilePage;
