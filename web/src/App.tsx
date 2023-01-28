import React from 'react';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import Layout from './layouts/Layout';
import TripsLayout from './layouts/TripsLayout'
import HomePage from './pages/home';
import TripPage from './pages/trips';

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    errorElement: (<div>Oops! Not found!</div>),
    children: [
      {
        path: "",
        element: <HomePage />,
      },
    ],
  },
  {
    path: "trips",
    element: <TripsLayout />,
    children: [
      {
        path: ":id",
        element: <TripPage />,
      },
    ],
  },
]);


function App() {
  return (
    <div className="App">
      <RouterProvider router={router} />
    </div>
  );
}

export default App;
