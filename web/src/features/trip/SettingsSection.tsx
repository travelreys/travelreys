import React, { FC, useEffect, useState } from 'react';
import _get from "lodash/get";
import _isEmpty from "lodash/isEmpty";
import { useTranslation } from 'react-i18next';
import { useDebounce } from 'usehooks-ts';
import { XMarkIcon } from '@heroicons/react/24/solid';
import { MagnifyingGlassCircleIcon } from '@heroicons/react/24/outline';

import Modal from '../../components/common/Modal';
import Avatar from '../../components/common/Avatar';

import AuthAPI, { SearchUsersResponse } from '../../apis/auth';
import {
  DefaultTransportModePref,
  LabelTransportModePref,
  MemberRoleCollaborator,
  MemberRoleParticipant,
  userFromMemberID
} from '../../lib/trips';
import { Auth, LabelUserGoogleImage } from '../../lib/auth';
import { Trips } from '../../lib/trips';
import { makeAddOp, makeReplaceOp } from '../../lib/tripsSync';
import { capitaliseWords } from '../../lib/strings';
import { CommonCss, TripSettingsCss } from '../../assets/styles/global';



////////////////////
// Transportation //
////////////////////

interface TransportSection {
  trip: any
  onSelect: (mode: string) => void
}

const TransportationSection: FC<TransportSection> = (props: TransportSection) => {

  const { t } = useTranslation();
  const [transportPref, setTransportPref] = useState(DefaultTransportModePref);

  useEffect(() => {
    setTransportPref(_get(
      props.trip,
      `labels.${LabelTransportModePref}`,
      DefaultTransportModePref
    ))
  }, [props.trip])


  // Event Handlers
  const transportPrefOnChange = (e: any) => {
    props.onSelect(e.target.value);
  }

  const selectID = "transportation";

  return (
    <div className='mb-8'>
      <h2 className='font-bold text-2xl mb-2'>
        {t("tripPage.settings.transportationTitle")}
      </h2>
      <div className='mb-2'>
        <label id={selectID} className={TripSettingsCss.TransportModeLabel}>
          {t("tripPage.settings.transportationModePreferenceLabel")}
        </label>
        <select
          id={selectID}
          value={transportPref}
          onChange={transportPrefOnChange}
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


/////////////////////
// AddMembersModal //
/////////////////////

interface AddMembersModalProps {
  isOpen: boolean
  tripUsers: { [key: string]: Auth.User }
  onClose: () => void
  onSelect: (id: string, role: string) => void
}

const AddMembersModal: FC<AddMembersModalProps> = (props: AddMembersModalProps) => {
  const { t } = useTranslation();

  const [searchEmail, setSearchEmail] = useState("");
  const [selectedMemberRole, setSelectedMemberRole] = useState(MemberRoleCollaborator);
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

  const memberRoleSelectOnChange = (e: any) => {
    setSelectedMemberRole(e.target.value);
  }

  const memberOnClick = (id: string) => {
    props.onSelect(id, selectedMemberRole)
    props.onClose();
  }

  // Renderers
  const renderHeader = () => {
    return (
      <div className={TripSettingsCss.MemberSearchHeader}>
        <div className={TripSettingsCss.MemberSearchHeaderTxt}>
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

  const renderSearchForm = () => {
    const selectID = "memberRole";
    return (
      <div className=''>
        <div className="relative mb-4">
          <div className={TripSettingsCss.MemberSearchIconCtn}>
            <MagnifyingGlassCircleIcon className={TripSettingsCss.MemberSearchIcon} />
          </div>
          <input
            type="text"
            className={TripSettingsCss.MemberSearchInput}
            placeholder={t('tripPage.settings.searchUsersPlaceholder') || ""}
            value={searchEmail}
            onChange={(e) => { setSearchEmail(e.target.value) }}
          />
        </div>
        <select
          id={selectID}
          value={selectedMemberRole}
          onChange={memberRoleSelectOnChange}
          className={TripSettingsCss.MemberRoleSelect}
        >
          <option value={MemberRoleCollaborator}>
            {t("tripPage.settings.memberRoleCollaborator")}
          </option>
          <option value={MemberRoleParticipant}>
            {t("tripPage.settings.memberRoleParticipant")}
          </option>
        </select>
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
          key={usr.id}
          type="button"
          className={TripSettingsCss.MemberSearchItem}
          disabled={isMember}
          onClick={() => memberOnClick(usr.id)}
        >
          <div className={TripSettingsCss.MemberSearchItemAvatar}>
            <Avatar
              placement="top"
              name={_get(usr, "name", "")}
              imgUrl={_get(usr, `labels.${LabelUserGoogleImage}`)}
            />
          </div>
          <div>
            <p className={TripSettingsCss.MemberSearchItemName}>
              {usr.name}
            </p>
            <p className={TripSettingsCss.MemberSearchItemDesc}>
              {isMember ? t('tripPage.settings.alreadyMember') : `${usr.email}`}
            </p>
          </div>
        </button>
      )
    });
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className='bg-white p-5 rounded-lg'>
        {renderHeader()}
        {renderSearchForm()}
        {renderSearchResults()}
      </div>
    </Modal>
  );
}


////////////////////
// MembersSection //
////////////////////

interface MembersSectionProps {
  trip: any
  tripUsers: { [key: string]: Auth.User }
  onAddUser: (id: string, role: string) => void
}

const MembersSection: FC<MembersSectionProps> = (props: MembersSectionProps) => {

  const { t } = useTranslation();
  const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false);

  const renderMembersAvatar = () => {
    let members = { [props.trip.creator.id]: props.trip.creator } as any;
    members = Object.assign(members, props.trip.members);
    return Object.values(members).map((mem: any) => {
      const usr = userFromMemberID(mem, props.tripUsers);
      return (
        <div key={mem.id} className='flex items-center py-4 border-b border-gray-200'>
          <div className={TripSettingsCss.MemberAvatarDiv}>
            <Avatar
              placement="top"
              name={`${_get(usr, "name", "")} (${capitaliseWords(mem.role)}})`}
              imgUrl={_get(usr, `labels.${LabelUserGoogleImage}`)}
            />
          </div>
          <div>
            <p className={TripSettingsCss.MemberSearchItemDesc}>
              {capitaliseWords(mem.role)}
            </p>
            <p className={TripSettingsCss.MemberSearchItemName}>
              {_get(usr, "name", "")}
            </p>
          </div>
        </div>
      );
    });
  }

  return (
    <>
      <div className={TripSettingsCss.MemberSectionCtn}>
        <div className={TripSettingsCss.MemberSectionHeader}>
          <h2 className={TripSettingsCss.MemberSectionTitle}>
            {t("tripPage.settings.membersTitle")}
          </h2>
          <button
            type="button"
            className={TripSettingsCss.SearchMemberBtn}
            onClick={() => setIsAddMemberModalOpen(true)}
          >
            + {t("tripPage.settings.searchMember")}
          </button>
        </div>
        <div>
          {renderMembersAvatar()}
        </div>
      </div>
      <AddMembersModal
        isOpen={isAddMemberModalOpen}
        tripUsers={props.tripUsers}
        onClose={() => setIsAddMemberModalOpen(false)}
        onSelect={props.onAddUser}
      />
    </>
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

  // Event Handlers

  const transportPrefOnChange = (mode: string) => {
    const opFn = _get(props.trip, `/labels.${LabelTransportModePref}`)
      ? makeReplaceOp : makeAddOp;
    props.tripStateOnUpdate([opFn(`/labels/${LabelTransportModePref}`, mode)]);
  }

  const addNewUser = (id: string, role: string) => {
    const member = { id, role, labels: {} } as Trips.Member;
    props.tripStateOnUpdate([makeAddOp(`/members/${id}`, member)]);
  }



  return (
    <div className='p-5'>
      <TransportationSection
        trip={props.trip}
        onSelect={transportPrefOnChange}
      />
      <MembersSection
        trip={props.trip}
        tripUsers={props.tripUsers}
        onAddUser={addNewUser}
      />
    </div>
  );
}

export default SettingsSection;
