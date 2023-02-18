import React, { FC } from 'react';
import { ReactComponent as HeroSVG } from '../../assets/images/undraw_floating_re_xtcj.svg'

const LandingPage: FC = () => {

  // UI State

  // Event Handlers

  // Renderers
  return (
    <main>
      <div className="pt-24 md:pt-36 px-6 mx-auto flex flex-wrap flex-col md:flex-row items-center">
        <div className="flex flex-col w-full xl:w-2/5 justify-center lg:items-start overflow-y-hidden">
          <h1 className="my-4 text-3xl md:text-5xl text-gray-700 font-bold leading-tight text-center md:text-left">
            A fun multiplayer social travel app.
          </h1>
          <p className="leading-normal text-base md:text-2xl mb-4 text-center md:text-left">
            Be it solo trips or group vacations, start your next adventure with us
          </p>
          <button
            type="button"
            className='bg-indigo-500 px-6 py-2 text-white font-bold rounded-lg mb-4 hover:bg-indigo-600'
          >
            Join our beta
          </button>
        </div>

        <div className="w-full xl:w-3/5  overflow-y-hidden">
          <HeroSVG  className="w-5/6 mx-auto lg:mr-0" />
        </div>

        <div className="w-full pt-16 pb-6 text-sm text-center md:text-left fade-in">
          <a className="text-gray-500 no-underline hover:no-underline" href="#">&copy; tiinyplanet 2023</a>
        </div>

      </div>
    </main>
  );
}


export default LandingPage;
