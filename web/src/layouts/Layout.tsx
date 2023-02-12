import React from 'react';
import { Outlet } from "react-router-dom";
import Navbar from '../components/common/Navbar';


export default function Layout() {
  return (
    <main className="container px-4 mx-auto mb-1">
      <Navbar />
      <Outlet />
    </main>
  )
}
