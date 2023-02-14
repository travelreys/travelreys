import React from 'react';
import { Outlet } from "react-router-dom";
import Navbar from '../components/common/Navbar';

export default function Layout() {
  return (
    <main className='min-h-screen'>
      <Outlet />
    </main>
  )
}