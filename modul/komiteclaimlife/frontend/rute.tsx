// Rute modul Komite Claim Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `App.tsx`, beserta keadaan kasus komite yang sedang dibuka.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanKomite } from './menu'
import InboxKomite from './pages/InboxKomite'
import KasusKomite from './pages/KasusKomite'
import './komiteclaimlife.css'

export function RuteKomite({ halaman, masuk }: PropsRute<HalamanKomite>) {
  // Kasus komite yang sedang dibuka dari Inbox Komite; kosong = daftar.
  const [kasusKomite, setKasusKomite] = useState('')

  return (
    // Akar gaya modul: semua aturan `komiteclaimlife.css` diawali `.komiteclaimlife` (`display: contents`).
    <div className="komiteclaimlife">
      {/* Komite Claim Life tiket 01 — Inbox Komite, lalu satu kasus dari baris. */}
      {halaman === 'komite' && kasusKomite === '' && <InboxKomite onBuka={setKasusKomite} peran={masuk.peran} />}
      {halaman === 'komite' && kasusKomite !== '' && (
        <KasusKomite
          kasusID={kasusKomite}
          peran={masuk.peran}
          onKembali={() => {
            setKasusKomite('')
          }}
        />
      )}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanKomite> = RuteKomite
