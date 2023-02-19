import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { useTranslation } from 'react-i18next';
import { useDebounce } from 'usehooks-ts';
import { XMarkIcon } from '@heroicons/react/24/solid';
import { MagnifyingGlassCircleIcon } from '@heroicons/react/24/outline';

import Modal from '../../components/common/Modal';
import Avatar from '../../components/common/Avatar';

import {
  DefaultTransportationPreference,
  LabelTransportationPreference,
  userFromMemberID
} from '../../lib/trips';
import { Auth, LabelUserGoogleImage } from '../../lib/auth';
import { makeAddOp, makeReplaceOp } from '../../lib/tripsSync';
import { CommonCss, TripSettingsCss } from '../../assets/styles/global';
import AuthAPI, { SearchUsersResponse } from '../../apis/auth';




/////////////////////
// AddMembersModal //
/////////////////////

interface AddMembersModalProps {
  isOpen: boolean
  tripUsers: { [key: string]: Auth.User }
  onClose: () => void
  onSelect: () => void
}

const AddMembersModal: FC<AddMembersModalProps> = (props: AddMembersModalProps) => {
  const { t } = useTranslation();

  const [searchEmail, setSearchEmail] = useState("");
  const [foundUsers, setFoundUsers] = useState<Array<Auth.User>>([]);
  const debouncedValue = useDebounce<string>(searchEmail, 500);

  // API
  const searchUsers = (email: string) => {
    if (_isEmpty(email)) {
      return;
    }
    AuthAPI.searchUsers(searchEmail)
      .then((res: SearchUsersResponse) => {
        setFoundUsers(res.users || []);
      })
  }

  // Event Handlers

  useEffect(() => {
    if (!_isEmpty(debouncedValue)) {
      searchUsers(debouncedValue);
    }
  }, [debouncedValue]);

  // Renderers
  const renderHeader = () => {
    return (
      <div className='flex justify-between items-center mb-4'>
        <div className='text-gray-800 font-bold text-lg'>
          {t("tripPage.settings.addMemberTitle")}
        </div>
        <button
          type="button"
          onClick={() => { props.onClose() }}
        >
          <XMarkIcon className={CommonCss.Icon} />
        </button>
      </div>
    );
  }

  const renderSearchInput = () => {
    return (
      <div className="relative mb-2">
        <div className={TripSettingsCss.SearchIconCtn}>
          <MagnifyingGlassCircleIcon className={TripSettingsCss.SearchIcon} />
        </div>
        <input
          type="text"
          className={TripSettingsCss.SearchInput}
          placeholder={t('tripPage.settings.searchUsersPlaceholder') || ""}
          value={searchEmail}
          onChange={(e) => { setSearchEmail(e.target.value) }}
        />
      </div>
    );
  }

  const renderSearchResults = () => {
    if (_isEmpty(searchEmail) && _isEmpty(foundUsers)) {
      return null;
    }
    if (_isEmpty(foundUsers)) {
      return <div>{t('tripPage.settings.noUsersFound')}</div>;
    }
    return foundUsers.map((usr: Auth.User) => {
      const isMember = Object.hasOwn(props.tripUsers, usr.id);
      return (
        <button
          type="button"
          className='flex items-center p-2 mb-4 text-left'
          disabled={isMember}
          onClick={props.onSelect}
        >
          <div key={usr.id} className='inline-block h-10 w-10 mr-4'>
            <Avatar
              name={_get(usr, "name", "")}
              imgUrl={_get(usr, `labels.${LabelUserGoogleImage}`)}
              placement="top"
            />
          </div>
          <div>
            <p className='font-semibold'>
              {usr.name}
            </p>
            <p className='text-gray-500'>
            {isMember ? "Already a member" : ""}
            </p>
          </div>
        </button>
      )
    });
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className='bg-white p-5'>
        {renderHeader()}
        {renderSearchInput()}
        {renderSearchResults()}
        <div className='flex justify-around items-center'>
          <button
            type="button"
            className='bg-indigo-500 px-4 py-2 rounded-lg font-bold text-sm text-white'
            onClick={props.onSelect}
          >
            {t("common.submit")}
          </button>
        </div>
      </div>
    </Modal>
  );
}

/////////////////////
// SettingsSection //
/////////////////////

interface SettingsSectionProps {
  trip: any
  tripUsers: { [key: string]: Auth.User }
  tripStateOnUpdate: any
}

const SettingsSection: FC<SettingsSectionProps> = (props: SettingsSectionProps) => {
  const { t } = useTranslation();
  const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false);
  const [transportationPreference, setTransportationPreference] = useState(DefaultTransportationPreference);

  useEffect(() => {
    setTransportationPreference(_get(
      props.trip,
      `labels.${LabelTransportationPreference}`,
      DefaultTransportationPreference
    ))
  }, [props.trip])

  // Event Handlers

  const transportationPreferenceOnChange = (e: any) => {
    setTransportationPreference(e.target.value)

    const op = _get(props.trip, `/labels.${LabelTransportationPreference}`)
      ? makeReplaceOp : makeAddOp;
    props.tripStateOnUpdate([
      op(`/labels/${LabelTransportationPreference}`, e.target.value)
    ]);
  }

  //  Renderers

  const renderTransportationMode = () => {
    const selectID = "transportation"
    return (
      <div className='mb-4'>
        <h2 className='font-bold text-xl mb-2'>
          {t("tripPage.settings.transportationTitle")}
        </h2>
        <div className='mb-2'>
          <label id={selectID} className={TripSettingsCss.TransportModeLabel}>
            {t("tripPage.settings.transportationModePreferenceLabel")}
          </label>
          <select
            id={selectID}
            value={transportationPreference}
            onChange={transportationPreferenceOnChange}
            className={TripSettingsCss.TransportModeSelect}
          >
            <option value="walk+drive">
              {t("tripPage.settings.transportationModePreferenceWalk+Drive")}
            </option>
            <option value="walk+transit">
              {t("tripPage.settings.transportationModePreferenceWalk+Transit")}
            </option>
          </select>
        </div>
      </div>
    );
  }

  const renderMembers = () => {
    let members = { [props.trip.creator.id]: props.trip.creator } as any;
    members = Object.assign(members, props.trip.members);
    members = Object.values(members).map((mem: any) => {
      const user = userFromMemberID(mem, props.tripUsers);
      return (
        <div key={mem.id} className='inline-block h-12 w-12'>
          <Avatar
            name={_get(user, "name", "")}
            imgUrl={_get(user, `labels.${LabelUserGoogleImage}`)}
            placement="top"
          />
        </div>
      );
    });
    return (
      <div className='mb-4'>
        <div className='flex justify-between items-center'>
          <h2 className='font-bold text-xl mb-2'>
            {t("tripPage.settings.membersTitle")}
          </h2>
          <button
            type="button"
            className='font-semibold text-gray-500'
            onClick={() => setIsAddMemberModalOpen(true)}
          >
            + {t("tripPage.settings.searchMember")}
          </button>
        </div>
        <div className='grid grid-cols-8 gap-4'>
          {members}
        </div>
      </div>
    );
  }

  return (
    <>
      <div className='p-5'>
        {renderTransportationMode()}
        {renderMembers()}
      </div>
      <AddMembersModal
        isOpen={isAddMemberModalOpen}
        tripUsers={props.tripUsers}
        onClose={() => setIsAddMemberModalOpen(false)}
        onSelect={() => { }}
      />
    </>
  );
}

export default SettingsSection;
