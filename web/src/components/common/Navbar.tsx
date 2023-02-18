import React, { FC, useRef, useState } from 'react';
import {
  Link,
  useNavigate,
  useLocation
} from 'react-router-dom';
import _get from 'lodash/get';
import _find from 'lodash/find';
import { useGoogleLogin } from '@react-oauth/google';
import { useTranslation } from 'react-i18next';
import {
  ArrowLeftOnRectangleIcon,
  ChevronDownIcon,
  GlobeAmericasIcon,
  XMarkIcon
} from '@heroicons/react/24/solid';


import Modal from './Modal';
import Dropdown from './Dropdown';
import GoogleIcon from '../icons/GoogleIcon';

import AuthAPI, {
  LoginResponse,
  makeUpdateUserFilter,
  ReadUserResponse
} from '../../apis/auth';
import {
  deleteAuthToken,
  deleteAuthUser,
  LabelUserGoogleImage,
  persistAuthToken,
  readAuthMetadata,
  LabelCurrency,
  LabelLocale
} from '../../lib/auth';
import {
  NavbarCss,
  CommonCss,
  CurrencyDropdownCss
} from '../../assets/styles/global';
import { makeSetUserAction, useUser } from '../../context/user-context';
import useOutsideAlerter from '../../hooks/useOutsideAlerter';
import currencies from '../../data/currency.json';
import locales from '../../data/locales.json';

////////////////
// LoginModal //
////////////////

interface LoginModalProps {
  isOpen: boolean
  onClose: () => void
}


const LoginModal: FC<LoginModalProps> = (props: LoginModalProps) => {
  const history = useNavigate();
  const { dispatch } = useUser();
  const { t } = useTranslation();

  // Event Handlers
  const googleLoginOnClick = useGoogleLogin({
    // hint: "",
    flow: 'auth-code',
    onSuccess: codeResponse => {
      AuthAPI.login(codeResponse.code)
        .then((res: LoginResponse) => {
          if (res.error) {
            // do smth with error
          }
          persistAuthToken(res.jwtToken!);
          return readAuthMetadata();
        })
        .then((metadata) => {
          return AuthAPI.readUser(metadata!.sub)
        })
        .then((res: ReadUserResponse) => {
          if (res.error) {
            // do smth with error
          }
          dispatch(makeSetUserAction(res.user!));
          history(`/home`);
          props.onClose();
        });
    },
  });

  // Renderers
  const renderGoogleLoginBtn = () => {
    return (
      <button
        className='inline-flex items-center rounded-full bg-white border border-gray-200 p-2 px-4 font-semibold'
        onClick={googleLoginOnClick}
      >
        <GoogleIcon className={CommonCss.LeftIcon} />
        {t('navbar.loginModal.googleSignIn')}
      </button>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className='p-5 py-8 flex flex-col'>
        <div className='flex flex-row-reverse mb-2'>
          <button
            type="button"
            className=''
            onClick={props.onClose}
          >
            <XMarkIcon className={CommonCss.Icon} />
          </button>
        </div>
        <h1 className='font-bold text-2xl text-center mb-8'>
          {t('navbar.loginModal.title')}
        </h1>
        <div className='flex justify-around mb-4'>
          {renderGoogleLoginBtn()}
        </div>
      </div>
    </Modal>
  );
}


//////////////////////
// CurrencySelector //
//////////////////////

interface CurrencySelectorProps {
  currency?: string
  onSelect: (code: string) => void
}

const CurrencySelector: FC<CurrencySelectorProps> = (props: CurrencySelectorProps) => {

  const [isActive, setIsActive] = useState(false);
  const wrapperRef = useRef(null);
  useOutsideAlerter(wrapperRef, () => {setIsActive(false)});
  const { t } = useTranslation();


  // Renderers
  const renderSelection = () => {
    const opts = currencies.map((loc: any) => (
      <button
        key={loc.code}
        type="button"
        className='flex rounded-lg p-1 text-sm hover:bg-indigo-100'
        onClick={() => {props.onSelect(loc.code)}}
      >
        <div className='text-gray-400 mr-2'>{loc.code}</div>
        <div className='text-gray-700 text-left '>{loc.name}</div>
      </button>
    ))

    return (
      <div
        ref={wrapperRef}
        className={CurrencyDropdownCss.Ctn}
      >
        <div className={CurrencyDropdownCss.Wrapper}>
          <h3 className='font-bold mb-2'>
            {t("navbar.currencySelector.title")}
          </h3>
          <div className='columns-2 sm:columns-4 smgap-4'>
            {opts}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className='relative'>
      <button
        type="button"
        className='flex items-center p-2 rounded-lg gap-1 hover:bg-gray-200'
        onClick={() => { setIsActive(!isActive) }}
      >
        <span className='font-semibold text-sm'>{props.currency}</span>
        <ChevronDownIcon className={CommonCss.Icon} />
      </button>
      {isActive ? renderSelection() : null}
    </div>
  );
}

////////////////////
// LocaleSelector //
////////////////////

interface LocaleSelctorProps {
  locale?: string
  onSelect: (locale: string) => void
}

const LocaleSelector: FC<LocaleSelctorProps> = (props: LocaleSelctorProps) => {
  const [isActive, setIsActive] = useState(false);
  const wrapperRef = useRef(null);
  const { t } = useTranslation();

  useOutsideAlerter(wrapperRef, () => {setIsActive(false)});

  // Renderers
  const renderSelection = () => {
    const opts = locales.map((loc: any) => (
      <button
        key={loc.locale}
        type="button"
        className='flex rounded-lg p-2 text-sm hover:bg-indigo-100'
        onClick={() => {props.onSelect(loc.locale)}}
      >
        <div className='text-gray-700 text-left '>{loc.name}</div>
      </button>
    ))

    return (
      <div
        ref={wrapperRef}
        className={CurrencyDropdownCss.Ctn}
      >
        <div className={CurrencyDropdownCss.Wrapper}>
          <h3 className='font-bold mb-2'>
            {t('navbar.localeSelector.title')}
          </h3>
          <div className='columns-3 sm:columns-4 smgap-4'>
            {opts}
          </div>
        </div>
      </div>
    );
  }

  const renderSelectedLocale = () => {
    return _get(
      _find(locales, (loc) => loc.locale === props.locale),
      "name",
      props.locale
    );
  }

  return (
    <div className='relative'>
      <button
        type="button"
        className='flex items-center p-2 rounded-lg gap-1 hover:bg-gray-200'
        onClick={() => { setIsActive(!isActive) }}
      >
        <span className='font-semibold text-sm'>
          {renderSelectedLocale()}
        </span>
        <ChevronDownIcon className={CommonCss.Icon} />
      </button>
      {isActive ? renderSelection() : null}
    </div>
  );
}


////////////
// Navbar //
////////////


interface LandingPageActionsProps {
  onLoginClick: () => void
}

const LandingPageActions: FC<LandingPageActionsProps> = (props: LandingPageActionsProps) => {
  const {t} = useTranslation();
  return (
    <div>
      <button
        type="button"
        className='font-bold py-2 px-6 rounded-full hover:text-indigo-500'
        onClick={props.onLoginClick}
      >
        {t('navbar.landingPageActions.login')}
      </button>
    </div>
  );
}


interface AppPageActionProps { }

const AppPageActions: FC<AppPageActionProps> = (props: AppPageActionProps) => {
  const history = useNavigate();
  const { state, dispatch } = useUser();
  const { t, i18n } = useTranslation();

  // Event Handlers
  const logoutOnClick = () => {
    deleteAuthToken();
    deleteAuthUser();
    history('/')
  }

  const currencyOnSelect = (cur: string) => {
    const newUser = Object.assign({}, state.user);
    newUser.labels[LabelCurrency] = cur;
    dispatch(makeSetUserAction(newUser));

    AuthAPI.updateUser(
      state.user?.id || "", makeUpdateUserFilter(newUser.labels));
  }

  const localeOnSelect = (loc: string) => {
    const newUser = Object.assign({}, state.user);
    newUser.labels[LabelLocale] = loc;
    dispatch(makeSetUserAction(newUser));

    AuthAPI.updateUser(
      state.user?.id || "", makeUpdateUserFilter(newUser.labels))
    .then(() => {
      i18n.changeLanguage(loc);
    })
  }

  // Renderers
  const renderProfileImage = () => {
    const profileImgURL = _get(state.user, `labels.${LabelUserGoogleImage}`);
    return (
      <img className={NavbarCss.ProfileImg}
        src={profileImgURL}
        alt="profile image"
        referrerPolicy="no-referrer"
      />
    );
  }

  const renderProfileDropdown = () => {
    const opts = [
      <button
        type='button'
        className={NavbarCss.LogoutBtn}
        onClick={logoutOnClick}
      >
        <ArrowLeftOnRectangleIcon className={CommonCss.LeftIcon} />
        {t('navbar.appPageActions.logout')}
      </button>,
    ];
    const menu = renderProfileImage();
    return <Dropdown menu={menu} opts={opts} />
  }

  const currency = _get(state.user, `labels.${LabelCurrency}`, "USD");
  const locale = _get(state.user, `labels.${LabelLocale}`, "en");

  return (
    <div className='flex items-center gap-2'>
      <LocaleSelector
        locale={locale}
        onSelect={localeOnSelect}
      />
      <CurrencySelector
        currency={currency}
        onSelect={currencyOnSelect}
      />
      {renderProfileDropdown()}
    </div>
  );
}



const NavBar: FC = () => {
  const location = useLocation()

  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);


  const isLandingPage = () => {
    return location.pathname === "/";
  }

  const isAppPage = () => {
    return location.pathname !== "/";
  }

  const logoHref = () => {
    return isLandingPage() ? "/" : "/home"
  }

  // Renderers
  const renderNavbarActions = () => {
    if (isLandingPage()) {
      return (<LandingPageActions onLoginClick={() => setIsLoginModalOpen(true)} />);
    }
    if (isAppPage()) {
      return (<AppPageActions />);
    }
    return null;
  }

  return (
    <nav className={NavbarCss.Ctn}>
      <Link to={logoHref()} className={NavbarCss.Link}>
        <GlobeAmericasIcon className={NavbarCss.LogoIcon} />
        <span className={NavbarCss.LogoTxt}>
          tiinyplanet
        </span>
      </Link>
      {renderNavbarActions()}
      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={() => setIsLoginModalOpen(false)}
      />

    </nav>
  );
}

export default NavBar;
