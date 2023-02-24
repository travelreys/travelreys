import React, { FC } from 'react';

// Import Swiper styles
import { Swiper, SwiperSlide } from 'swiper/react';
import 'swiper/css';
import 'swiper/css/navigation';
import 'swiper/css/pagination';
import { PLACE_IMAGE_APIKEY } from '../../apis/maps';


const gMapsPlaceImageSrcURL = (ref: string) => {
  return [
    "https://maps.googleapis.com/maps/api/place/photo",
    "?maxwidth=1024",
    `&photo_reference=${ref}`,
    `&key=${PLACE_IMAGE_APIKEY}`,
  ].join("");
}

// PlacePicturesCarousel

interface PlacePicturesCarouselProps {
  photos: any
}

const PlacePicturesCarousel: FC<PlacePicturesCarouselProps> = (props: PlacePicturesCarouselProps) => {
  if (props.photos.length === 0) {
    return (<></>);
  }

  const landscape = props.photos.filter((photo: any) => photo.width > photo.height);
  return (
    <Swiper slidesPerView={1}>
      {landscape.map((photo: any) => (
        <SwiperSlide key={photo.photo_reference}>
          <img
            className="rounded h-72 w-full"
            alt={"place"}
            src={gMapsPlaceImageSrcURL(photo.photo_reference)} />
        </SwiperSlide>
      ))}
    </Swiper>
  );
}

export default PlacePicturesCarousel;
