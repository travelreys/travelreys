import { useRouter } from "next/router";

const Trip = () => {
  const router = useRouter();
  const { id } = router.query;

  return <p>Trip: {id}</p>
}

export default Trip;