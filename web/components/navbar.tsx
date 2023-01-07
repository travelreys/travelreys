import React, { FC } from 'react';
import Link from 'next/link';

import { GlobeAmericasIcon } from '@heroicons/react/24/solid'

const NavBar: FC = () => {
  return (
    <div className="container flex mb-1 px-24 py-5">
      <h1 className="text-3xl font-bold text-indigo-500">
         <Link href="/">
            <GlobeAmericasIcon className='inline align-bottom h-8 w-8'/>
            <span className='inline-block pl-1'>tiinyplanet</span>
          </Link>
      </h1>
    </div>
  );
}

export default NavBar;
