// Rute modul Claim Life - refactor bentuk B (30-09-2026): dipindah apa adanya
// dari `App.tsx`, beserta keadaan kasus yang sedang dibuka.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanClaimLife } from './menu'
import InboxClaimLife from './pages/InboxClaimLife'
import KlaimLife from './pages/KlaimLife'
import OutstandingClaimLife from './pages/OutstandingClaimLife'
import RegisterKlaim from './pages/RegisterKlaim'

export function RuteClaimLife({ halaman, masuk, onPindah }: PropsRute<HalamanClaimLife>) {
  // Kasus yang sedang dibuka. Kosong berarti belum ada yang dipilih.
  const [kasus, setKasus] = useState('')

  return (
    <>
      {halaman === 'inbox' && (
        <InboxClaimLife
          peran={masuk.peran}
          onBuka={(workID) => {
            // ⚠️ Baris Inbox membuka layar TAHAPnya. Tab Outstanding
            // membuka `OSClaimLife`; tahap lain menyusul bersama
            // kelompok A3 masing-masing.
            setKasus(workID)
            onPindah('outstanding')
          }}
          onRegister={() => {
            onPindah('register')
          }}
        />
      )}
      {halaman === 'outstanding' && (
        <OutstandingClaimLife
          klaimID={kasus}
          onPindah={() => {
            onPindah('inbox')
          }}
          onDetail={() => {
            onPindah('detail')
          }}
        />
      )}
      {halaman === 'register' && <RegisterKlaim />}
      {halaman === 'detail' && <KlaimLife />}
    </>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanClaimLife> = RuteClaimLife
