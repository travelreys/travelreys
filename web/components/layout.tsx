import Navbar from './navbar';

type LayoutProps = {
  children: React.ReactNode,
};

export default function Layout(props: LayoutProps) {
  return (
    <>
      <Navbar />
      <main>{props.children}</main>
    </>
  )
}