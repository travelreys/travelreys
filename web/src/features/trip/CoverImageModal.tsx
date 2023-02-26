import React, { FC, useState } from 'react';
import _get from "lodash/get";
import {
  MagnifyingGlassIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline'

import Modal from '../../components/common/Modal';
import Spinner from '../../components/common/Spinner';
import ImagesAPI from '../../apis/images';

interface CoverImageModalProps {
  isOpen: boolean
  onClose: any
  onCoverImageSelect: any
}

const CoverImageModal: FC<CoverImageModalProps> = (props: CoverImageModalProps) => {

  const [query, setQuery] = useState("");
  const [imageList, setImageList] = useState([] as any);
  const [isLoading, setIsLoading] = useState(false);

  // API
  const searchImage = () => {
    setIsLoading(true);
    ImagesAPI.search(query)
    .then(res => {
      const images = _get(res, "data.images", []);
      setImageList(images);
      setIsLoading(false);
    });
  }

  // Renderers
  const css = {
    figure: "relative max-w-sm transition-all rounded-lg duration-300 mb-2 group",
    figureImg: "block rounded-lg max-w-full group-hover:grayscale",
    figureBtn: "text-white m-2 py-2 px-3 rounded-full bg-green-500 hover:bg-green-700",
    figureBtnCtn: "absolute group-hover:opacity-100 opacity-0 top-2 right-0",
    figureCaption: "absolute px-1 text-sm text-white rounded-b-lg bg-slate-800/50 w-full bottom-0",
    searchImageCard: "bg-white px-4 pt-5 pb-4 sm:p-8 sm:pb-4 rounded-lg",
    searchImageTitle: "text-lg sm:text-2xl font-bold leading-6 text-slate-900",
    searchImageWebTitle: "text-sm font-medium text-indigo-500 sm:text-xl text-slate-700 mb-2 ml-1",
    searchImageInput: "bg-gray-50 block border-gray-300 border focus:border-blue-500 focus:ring-blue-500 min-w-0 p-2.5 rounded-lg text-gray-900 text-sm w-5/6 mr-2",
    searchImageBtn: "flex-1 inline-flex text-white bg-indigo-500 hover:bg-indigo-800 rounded-2xl p-2.5 text-center items-center justify-around",
    searchImageIcon: "h-5 w-5 stroke-2 stroke-white",
  }

  const renderImageThumbnails = () => {
    if (isLoading) {
      return <Spinner />
    }
    return (
      <div className='columns-2 md:columns-3'>
        { imageList.map((image: any) => (
          <figure
            key={image.id}
            className={css.figure}
          >
            <button type="button">
              <img
                srcSet={ImagesAPI.makeSrcSet(image)}
                src={ImagesAPI.makeSrc(image)}
                alt={"cover"}
                className={css.figureImg}
              />
              <div className={css.figureBtnCtn}>
                <button
                  type="button"
                  className={css.figureBtn}
                  onClick={() => {props.onCoverImageSelect(image)}}
                >
                  Select
                </button>
              </div>
            </button>
            <figcaption className={css.figureCaption}>
              <a
                target="_blank"
                href={ImagesAPI.makeUserURL(_get(image, "user.username"))}
                rel="noreferrer"
              >
                @{_get(image, "user.username")}, Unsplash
              </a>
            </figcaption>
          </figure>
        ))}
      </div>);
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className={css.searchImageCard}>
        <div className='flex justify-between mb-6'>
          <h2 className={css.searchImageTitle}>
            Change cover image
          </h2>
          <button type="button" onClick={props.onClose}>
            <XMarkIcon className='h-6 w-6 text-slate-700' />
          </button>
        </div>
        <h2 className={css.searchImageWebTitle}>
          Search the web
        </h2>
        <div className="flex mb-4 justify-between">
          <input
            type="text"
            className={css.searchImageInput}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" ? searchImage() : ""}
            placeholder="destination, theme ..."
          />
          <button
            type='button'
            className={css.searchImageBtn}
            onClick={searchImage}
          >
            <MagnifyingGlassIcon className={css.searchImageIcon} />
          </button>
        </div>
        {renderImageThumbnails()}
      </div>
    </Modal>
  );
}

export default CoverImageModal;
