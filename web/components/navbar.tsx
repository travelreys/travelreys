import React, { FC } from 'react';
import Link from 'next/link';

import { GlobeAmericasIcon } from '@heroicons/react/24/solid'

const NavBar: FC = () => {
  return (
    <nav className="container py-5 flex">
      <h1 className="text-2xl sm:text-3xl font-bold text-indigo-500">
         <Link href="/">
            <GlobeAmericasIcon className='inline align-bottom h-8 w-8'/>
            <span className='inline-block pl-1'>tiinyplanet</span>
          </Link>
      </h1>
    </nav>
  );
}

export default NavBar;
