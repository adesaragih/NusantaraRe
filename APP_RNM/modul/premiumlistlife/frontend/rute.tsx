// Rute modul PremiumList Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `App.tsx`, beserta keadaan polis yang sedang dibuka.

import { useState } from 'react'

import type { PropsRute } from '../../../inti/frontend/modul'
import { TAHAP_POLIS } from './api'
import type { HalamanPremiumList } from './menu'
import InboxPremiumList from './pages/InboxPremiumList'
import InputOffer from './pages/InputOffer'
import PremiumListDetail from './pages/PremiumListDetail'
import PremiumListSummary from './pages/PremiumListSummary'

export function RutePremiumList({ halaman }: PropsRute<HalamanPremiumList>) {
  // Polis yang sedang dibuka, beserta tahapnya - tiket 01 PremiumList.
  const [polis, setPolis] = useState({ id: '', tahap: '' })

  return (
    <>
      {halaman === 'premiumlist' && polis.id === '' && (
        <InboxPremiumList
          onBuka={(caseID, tahap) => {
            setPolis({ id: caseID, tahap })
          }}
        />
      )}
      {/*
        ⛔ DUA LAYAR DI TAHAP YANG SAMA, dan itu bentuk aslinya:
        `ShowLifePremiumDetail` memuat grid peserta DAN tombol keputusannya.
        `Reject` bahkan HANYA punya konektor di tahap ini (`Transition9`
        b2306), jadi memisahkan gridnya dari tombolnya berarti menyembunyikan
        satu-satunya tempat `Reject` dapat ditekan.
      */}
      {halaman === 'premiumlist' &&
        polis.id !== '' &&
        polis.tahap === TAHAP_POLIS.detail && (
          <PremiumListDetail polisID={polis.id} />
        )}
      {/* Tiket 05a bagian 2 — `ShowLifePremiumSummary`, tahap Input Premium Summary. */}
      {halaman === 'premiumlist' &&
        polis.id !== '' &&
        polis.tahap === TAHAP_POLIS.summary && <PremiumListSummary polisID={polis.id} />}
      {halaman === 'premiumlist' && polis.id !== '' && (
        <InputOffer
          polisID={polis.id}
          tahap={polis.tahap}
          onSelesai={() => {
            setPolis({ id: '', tahap: '' })
          }}
        />
      )}
    </>
  )
}
