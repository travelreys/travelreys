import { useRouter } from "next/router";

import TripsAPI from '../../apis/trips';


const Trip = () => {
  const router = useRouter();
  const { id } = router.query;

  let {data, error, isLoading} = TripsAPI.readTrip(id as string);
  console.log(data)


  return <p>Trip: {id}</p>
}

export default Trip;