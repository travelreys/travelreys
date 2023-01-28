import React, { FC } from 'react';
import { Link } from 'react-router-dom';

import { GlobeAmericasIcon } from '@heroicons/react/24/solid'

const NavBar: FC = () => {
  return (
    <nav className="container py-5 flex">
      <Link to="/" className="text-2xl sm:text-3xl font-bold text-indigo-500">
        <GlobeAmericasIcon className='inline align-bottom h-8 w-8'/>
        <span className='inline-block pl-1'>tiinyplanet</span>
      </Link>
    </nav>
  );
}

export default NavBar;
