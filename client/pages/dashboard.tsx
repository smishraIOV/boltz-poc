import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import { GetInfoResponse, requestProvider, WebLNProvider } from 'webln';
import React, { useEffect, useState }  from 'react';

declare global {
  interface Window {
    webln: WebLNProvider
  }
}

const Dashboard: NextPage = () => {

  const [data, setData] = useState<[any]>()
  const [isLoading, setLoading] = useState(false)

  useEffect(() => {
    setLoading(true)
    fetch('http://localhost:8080/payment/')
      .then((res) => res.json())
      .then((data) => {
        setData(data)
        setLoading(false)
      })
  }, [])

  if (isLoading) return <p>Loading...</p>

  return (
    <div className={styles.container}>
      <Head>
        <title>PoC - Boltz Merchant </title>
        <meta name="description" content="Boltz Merchant PoC" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>
        Payments
        </h1>

        
        <div className={styles.grid}>
          <table>
            <thead>
              <tr>
                <th>
                  #
                </th>
                <th>
                  Sats
                </th>
                <th>
                  Preimage
                </th>
                <th>
                  Hash
                </th>
              </tr>
            </thead>
            <tbody>
            {data ? (
              data.length ? (
                data.map((todo, i) => {
                  return <tr key={todo.PreimageHash}>
                      <td>{i} </td>
                      <td>{todo.Amount} </td>
                      <td>{todo.Preimage}</td>
                      <td>{todo.PreimageHash}</td>
                  </tr>
                })
              ) : (
                <tr>
                  <td colSpan={4}>No payments yet.</td> 
                </tr>
              )
            ) : (
              <tr>
                <td colSpan={4}>No payments yet.</td> 
              </tr>
            )}
            </tbody>
          </table>
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

export default Dashboard
