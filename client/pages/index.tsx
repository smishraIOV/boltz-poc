import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import { GetInfoResponse, requestProvider, WebLNProvider } from 'webln';
import React, { useEffect, useState }  from 'react';
import BoltzPocClient from '../src/client';

declare global {
  interface Window {
    webln: WebLNProvider
  }
}

const Home: NextPage = () => {

  const [lnwallet, setLnwallet] = useState<GetInfoResponse>()
  const [isLoading, setLoading] = useState(false)
  const [isEnabled, setEnabled] = useState(false)

  const isProviderEnabled = async () => {
    setTimeout(async () => {
      try {
        await requestProvider()
        setLnwallet(await window.webln.getInfo())
        setEnabled(true);
        console.log('webln enabled', lnwallet )
      } catch (e) {
        setEnabled(false);
        console.error('webln NOT enabled', e)
      }
      setLoading(false)
    }, 100)
  }

  useEffect(() => {
    setLoading(true)
    isProviderEnabled()
  }, [])

  if (isLoading) return <div>Loading...</div>


  const onClick = async (e: any) => {
    e.preventDefault()
    if (isEnabled) {
      try {
        const lninvoice = await BoltzPocClient.CreatePayment()
        await window.webln.sendPayment(lninvoice);
      } catch (error) {
        // ignore
        console.error(error);
      }
    } else {
      //show QR code.
      //linkFallback(this.paymentRequest);
    }
  }

  return (
    <div className={styles.container}>
      <Head>
        <title>PoC - Boltz Merchant </title>
        <meta name="description" content="Boltz Merchant PoC" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>
          Boltz Merchant PoC
        </h1>

        
        <div className={styles.grid}>
          <div>
            
            <div>
              {!isEnabled ? 
                  <button onClick={(e)=>isProviderEnabled()}> Connect LN wallet<i aria-label="icon: thunderbolt" className="anticon anticon-thunderbolt"> <svg viewBox="64 64 896 896" data-icon="thunderbolt" width="1em" height="1em" fill="currentColor" aria-hidden="true"> <path d="M848 359.3H627.7L825.8 109c4.1-5.3.4-13-6.3-13H436c-2.8 0-5.5 1.5-6.9 4L170 547.5c-3.1 5.3.7 12 6.9 12h174.4l-89.4 357.6c-1.9 7.8 7.5 13.3 13.3 7.7L853.5 373c5.2-4.9 1.7-13.7-5.5-13.7z"></path> </svg> </i></button> 
                  : 
                  <p>Wallet connected: {lnwallet?.node.alias} ({lnwallet?.node.pubkey.slice(0, 5)}...)</p>
              }
            </div>
          </div>
        </div>

        <div className={styles.grid}>
          <div className={styles.card}>
            <h2>Product 1</h2>
            <p>Click to buy.</p>
            <p></p>
            <div>
                <button className={styles.button} onClick={(e)=>onClick(e)} disabled={!isEnabled}>
                  Pay with&nbsp;  <i aria-label="icon: thunderbolt" className="anticon anticon-thunderbolt"> <svg viewBox="64 64 896 896" data-icon="thunderbolt" width="1em" height="1em" fill="currentColor" aria-hidden="true"> <path d="M848 359.3H627.7L825.8 109c4.1-5.3.4-13-6.3-13H436c-2.8 0-5.5 1.5-6.9 4L170 547.5c-3.1 5.3.7 12 6.9 12h174.4l-89.4 357.6c-1.9 7.8 7.5 13.3 13.3 7.7L853.5 373c5.2-4.9 1.7-13.7-5.5-13.7z"></path> </svg> </i>
                </button>
              <button className={styles.button} onClick={(e)=>onClick(e)}>
                Pay with DOC
              </button>
            </div>
          </div>
          <div className={styles.card}>
            <h2>Product 2</h2>
            <p>Click to buy.</p>
            <p></p>
            <div>
                <button className={styles.button} onClick={(e)=>onClick(e)} disabled={!isEnabled}>
                  Pay with&nbsp;  <i aria-label="icon: thunderbolt" className="anticon anticon-thunderbolt"> <svg viewBox="64 64 896 896" data-icon="thunderbolt" width="1em" height="1em" fill="currentColor" aria-hidden="true"> <path d="M848 359.3H627.7L825.8 109c4.1-5.3.4-13-6.3-13H436c-2.8 0-5.5 1.5-6.9 4L170 547.5c-3.1 5.3.7 12 6.9 12h174.4l-89.4 357.6c-1.9 7.8 7.5 13.3 13.3 7.7L853.5 373c5.2-4.9 1.7-13.7-5.5-13.7z"></path> </svg> </i>
                </button>
              <button className={styles.button} onClick={(e)=>onClick(e)}>
                Pay with DOC
              </button>
            </div>
          </div>
          <div className={styles.card}>
            <h2>Product 3</h2>
            <p>Click to buy.</p>
            <p></p>
            <div>
                <button className={styles.button} onClick={(e)=>onClick(e)} disabled={!isEnabled}>
                  Pay with&nbsp;  <i aria-label="icon: thunderbolt" className="anticon anticon-thunderbolt"> <svg viewBox="64 64 896 896" data-icon="thunderbolt" width="1em" height="1em" fill="currentColor" aria-hidden="true"> <path d="M848 359.3H627.7L825.8 109c4.1-5.3.4-13-6.3-13H436c-2.8 0-5.5 1.5-6.9 4L170 547.5c-3.1 5.3.7 12 6.9 12h174.4l-89.4 357.6c-1.9 7.8 7.5 13.3 13.3 7.7L853.5 373c5.2-4.9 1.7-13.7-5.5-13.7z"></path> </svg> </i>
                </button>
              <button className={styles.button} onClick={(e)=>onClick(e)}>
                Pay with DOC
              </button>
            </div>
          </div>
        </div>
      </main>

      <footer className={styles.footer}>
        <a
          href="https://vercel.com?utm_source=create-next-app&utm_medium=default-template&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          Powered by{' '}
          <span className={styles.logo}>
            <Image src="/vercel.svg" alt="Vercel Logo" width={72} height={16} />
          </span>
        </a>
      </footer>
    </div>
  )
}

export default Home
