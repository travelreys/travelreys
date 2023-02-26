import React, { useEffect } from 'react';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import { GoogleOAuthProvider } from '@react-oauth/google';

import HomePage from './pages/home';
import LandingPage from './pages/landing';
import Layout from './layouts/Layout';
import TripPage from './pages/trip';
import ProfilePage from './pages/profile';

import { readAuthUser } from './lib/auth';
import {
  makeSetUserAction,
  UserProvider,
  useUser
} from './context/user-context';

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    errorElement: (<div>Oops! Not found!</div>),
    children: [
      {
        path: "",
        element: <LandingPage />,
      },
      {
        path: "home",
        element: <HomePage />,
      },
      {
        path: "profile",
        element: <ProfilePage />,
      },
    ],
  },
  {
    path: "trips",
    children: [
      {
        path: ":id",
        element: <TripPage />,
      },
    ],
  },
]);

const RouterApp: React.FC = () => {
  const {state, dispatch} = useUser();
  useEffect(() => {
    if (state.user === null) {
      dispatch(makeSetUserAction(readAuthUser()))
    }
  }, [])

  return (<RouterProvider router={router} />);
}

const App: React.FC = () => {
  return (
    <div className="app">
      <GoogleOAuthProvider
        clientId="697392212622-m3mcs1396bu9tuc8joqolrj6uid0u374.apps.googleusercontent.com"
      >
        <UserProvider>
          <RouterApp />
        </UserProvider>
      </GoogleOAuthProvider>
    </div>
  );

}

export default App;
