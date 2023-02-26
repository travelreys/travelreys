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
  ReadResponse
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
import { CommonCss } from '../../assets/styles/global';
import { makeSetUserAction, useUser } from '../../context/user-context';
import useOutsideAlerter from '../../hooks/useOutsideAlerter';
import currencies from '../../data/currency.json';
import locales from '../../data/locales.json';


interface LoginModalProps {
  isOpen: boolean
  onClose: () => void
}

const LoginModal: FC<LoginModalProps> = (props: LoginModalProps) => {
  const history = useNavigate();
  const { dispatch } = useUser();
  const { t } = useTranslation();

  const css = {
    ctn: 'p-5 py-8 flex flex-col',
    wrapper: "flex flex-row-reverse mb-2",
    title: 'font-bold text-2xl text-center mb-8',
    btnCtn: "flex justify-around mb-4",
    googleLoginBtn: 'inline-flex items-center rounded-full bg-white border border-gray-200 p-2 px-4 font-semibold',
  }

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
          return AuthAPI.read(metadata!.sub)
        })
        .then((res: ReadResponse) => {
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
      <button className={css.googleLoginBtn} onClick={googleLoginOnClick}>
        <GoogleIcon className={CommonCss.LeftIcon} />
        {t('navbar.loginModal.googleSignIn')}
      </button>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={css.ctn}>
        <div className={css.wrapper}>
          <button type="button" onClick={props.onClose}>
            <XMarkIcon className={CommonCss.Icon} />
          </button>
        </div>
        <h1 className={css.title}>
          {t('navbar.loginModal.title')}
        </h1>
        <div className={css.btnCtn}>
          {renderGoogleLoginBtn()}
        </div>
      </div>
    </Modal>
  );
}

interface CurrencySelectorProps {
  currency?: string
  onSelect: (code: string) => void
}

const CurrencySelector: FC<CurrencySelectorProps> = (props: CurrencySelectorProps) => {

  const [isActive, setIsActive] = useState(false);
  const wrapperRef = useRef(null);
  useOutsideAlerter(wrapperRef, () => {setIsActive(false)});

  const { t } = useTranslation();

  const css = {
    ctn: "absolute z-10 rounded-lg bg-white shadow block right-0",
    wrapper: "p-4 max-h-96 overflow-y-auto",
    title: 'font-bold mb-2',
    btnCtn: "columns-2 sm:columns-4 smgap-4",
    ddBtn: "flex items-center p-2 rounded-lg gap-1 hover:bg-gray-200",
    currency: "font-semibold text-sm",
    btn: "flex rounded-lg p-1 text-sm hover:bg-indigo-100",
    btnCode: "text-gray-400 mr-2",
    btnName: "text-gray-700 text-left",
  }

  // Renderers
  const renderSelection = () => {
    const opts = currencies.map((loc: any) => (
      <button
        key={loc.code}
        type="button"
        className={css.btn}
        onClick={() => {props.onSelect(loc.code)}}
      >
        <div className={css.btnCode}>{loc.code}</div>
        <div className={css.btnName}>{loc.name}</div>
      </button>
    ))

    return (
      <div ref={wrapperRef} className={css.ctn}>
        <div className={css.wrapper}>
          <h3 className={css.title}>
            {t("navbar.currencySelector.title")}
          </h3>
          <div className={css.btnCtn}>{opts}</div>
        </div>
      </div>
    );
  }

  return (
    <div className='relative'>
      <button
        type="button" className={css.ddBtn}
        onClick={() => { setIsActive(!isActive) }}
      >
        <span className={css.currency}>{props.currency}</span>
        <ChevronDownIcon className={CommonCss.Icon} />
      </button>
      {isActive ? renderSelection() : null}
    </div>
  );
}

interface LocaleSelctorProps {
  locale?: string
  onSelect: (locale: string) => void
}

const LocaleSelector: FC<LocaleSelctorProps> = (props: LocaleSelctorProps) => {
  const [isActive, setIsActive] = useState(false);
  const wrapperRef = useRef(null);
  useOutsideAlerter(wrapperRef, () => {setIsActive(false)});

  const { t } = useTranslation();

  const css = {
    ctn: "absolute z-10 rounded-lg bg-white shadow block right-0",
    wrapper: "p-4 max-h-96 overflow-y-auto",
    ddBtn: "flex items-center p-2 rounded-lg gap-1 hover:bg-gray-200",
    ddBtnText: "font-semibold text-sm",
    btn: "flex rounded-lg p-2 text-sm hover:bg-indigo-100",
    btnText: "text-gray-700 text-left"
  }

  // Renderers
  const renderSelection = () => {
    const opts = locales.map((loc: any) => (
      <button
        key={loc.locale}
        type="button"
        className={css.btn}
        onClick={() => {props.onSelect(loc.locale)}}
      >
        <div className={css.btnText}>{loc.name}</div>
      </button>
    ))

    return (
      <div
        ref={wrapperRef}
        className={css.ctn}
      >
        <div className={css.wrapper}>
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
    const loc = _find(locales, (loc) => loc.locale === props.locale);
    return _get(loc, "name", props.locale);
  }

  return (
    <div className='relative'>
      <button
        type="button" className={css.ddBtn}
        onClick={() => { setIsActive(!isActive) }}
      >
        <span className={css.ddBtnText}>
          {renderSelectedLocale()}
        </span>
        <ChevronDownIcon className={CommonCss.Icon} />
      </button>
      {isActive ? renderSelection() : null}
    </div>
  );
}

interface LandingPageActionsProps {
  onLoginClick: () => void
}

const LandingPageActions: FC<LandingPageActionsProps> = (props: LandingPageActionsProps) => {
  const {t} = useTranslation();

  const css = {
    btn: "font-bold py-2 px-6 rounded-full hover:text-indigo-500"
  }

  return (
    <div>
      <button
        type="button"
        className={css.btn}
        onClick={props.onLoginClick}
      >
        {t('navbar.landingPageActions.login')}
      </button>
    </div>
  );
}


interface AppPageActionProps {}

const AppPageActions: FC<AppPageActionProps> = (props: AppPageActionProps) => {
  const history = useNavigate();
  const { state, dispatch } = useUser();
  const { t, i18n } = useTranslation();

  const css = {
    ctn: "flex items-center gap-2",
    profileImg: "h-8 w-8 rounded-full",
    logoutBtn: "flex items-center w-full hover:text-indigo-500",
  }

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

    AuthAPI.update(state.user?.id || "", newUser.labels);
  }

  const localeOnSelect = (loc: string) => {
    const newUser = Object.assign({}, state.user);
    newUser.labels[LabelLocale] = loc;
    dispatch(makeSetUserAction(newUser));

    AuthAPI.update(
      state.user?.id || "", newUser.labels)
    .then(() => {
      i18n.changeLanguage(loc);
    })
  }

  // Renderers
  const renderProfileImage = () => {
    const profileImgURL = _get(state.user, `labels.${LabelUserGoogleImage}`);
    return (
      <img className={css.profileImg}
        src={profileImgURL}
        alt="profile"
        referrerPolicy="no-referrer"
      />
    );
  }

  const renderProfileDropdown = () => {
    const opts = [
      <button
        type='button'
        className={css.logoutBtn}
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
    <div className={css.ctn}>
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

interface NavbarLogoProps {
  href: string
}

export const NavbarLogo: FC<NavbarLogoProps> = (props:NavbarLogoProps) => {
  const css = {
    link: "text-2xl sm:text-3xl font-bold text-indigo-500",
    logoIcon: "inline align-bottom h-8 w-8 mr-2",
    logoTxt: "inline-block",
  }

  return (
    <Link to={props.href} className={css.link}>
      <GlobeAmericasIcon className={css.logoIcon} />
      <span className={css.logoTxt}>
        tiinyplanet
      </span>
    </Link>
  );
}

const NavBar: FC = () => {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const location = useLocation();

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
  const css = {
    ctn: "container py-5 flex justify-between items-center",
  }

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
    <nav className={css.ctn}>
      <NavbarLogo href={logoHref()} />
      {renderNavbarActions()}
      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={() => setIsLoginModalOpen(false)}
      />
    </nav>
  );
}

export default NavBar;
