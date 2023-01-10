import Navbar from '../NavBar';

type LayoutProps = {
  children: React.ReactNode,
};

export default function Layout(props: LayoutProps) {
  return (
    <>
      <main className="container px-4 mx-auto mb-1">
        <Navbar />
        {props.children}
      </main>
    </>
  )
}
