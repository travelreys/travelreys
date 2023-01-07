import type { AppProps } from 'next/app';
import Head from 'next/head';

import Layout from "../components/layout";

import "../styles/global.css";
import 'react-day-picker/dist/style.css';

interface MyAppProps extends AppProps {
}

export default function App(props: MyAppProps) {
  const {
    Component,
    pageProps
  } = props;
  return (
    <>
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <Layout>
        <Component {...pageProps} />
      </Layout>
    </>
  );
}
