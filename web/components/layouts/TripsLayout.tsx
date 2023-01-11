import Navbar from '../NavBar';

type LayoutProps = {
  children: React.ReactNode,
};

export default function Layout(props: LayoutProps) {
  return (
    <main className='min-h-screen'>
      {props.children}
    </main>
  )
}