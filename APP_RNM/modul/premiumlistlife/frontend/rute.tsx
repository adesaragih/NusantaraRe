// Rute modul PremiumList Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `App.tsx`, beserta keadaan polis yang sedang dibuka.

import { useEffect, useRef, useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { kasusBisaDibuka, TAHAP_POLIS } from './api'
import { HALAMAN_AWAL_PREMIUMLIST, type HalamanPremiumList } from './menu'
import InboxPremiumList from './pages/InboxPremiumList'
import InputOffer from './pages/InputOffer'
import PremiumListDetail from './pages/PremiumListDetail'
import PremiumListSummary from './pages/PremiumListSummary'
import './premiumlistlife.css'

/**
 * Apakah menu PremiumList Life baru saja dipilih (ulang).
 *
 * ⛔ `halaman` saja tidak cukup: memilih menu halaman yang SEDANG tampil tidak
 * mengubahnya. `ketukMenu` (`PropsRute`) bertambah setiap pilihan menu.
 */
export function menuDipilihUlang(
  ketukLalu: number | undefined,
  ketukBaru: number | undefined,
  halaman: string,
): boolean {
  return ketukBaru !== undefined && ketukBaru !== ketukLalu && halaman === HALAMAN_AWAL_PREMIUMLIST
}

export function RutePremiumList({ halaman, ketukMenu }: PropsRute<HalamanPremiumList>) {
  // Polis yang sedang dibuka, beserta tahapnya - tiket 01 PremiumList.
  const [polis, setPolis] = useState({ id: '', tahap: '' })

  // Memilih menu PremiumList Life (lagi) membawa kembali ke kotak masuk -
  // permintaan work owner 01-10-2026. Rute ini TETAP terpasang (keadaan polis
  // bertahan saat pindah halaman lewat cara lain), jadi kasus yang terbuka
  // ditutup HANYA saat menunya sendiri dipilih.
  const ketukLalu = useRef(ketukMenu)
  useEffect(() => {
    if (menuDipilihUlang(ketukLalu.current, ketukMenu, halaman)) {
      setPolis({ id: '', tahap: '' })
    }
    ketukLalu.current = ketukMenu
  }, [ketukMenu, halaman])

  return (
    // Akar gaya modul: semua aturan `premiumlistlife.css` diawali `.premiumlistlife` (`display: contents`).
    <div className="premiumlistlife">
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
      {/* Kasus tertutup tidak punya layar keputusan (keputusan work owner 02-10-2026). */}
      {halaman === 'premiumlist' && polis.id !== '' && kasusBisaDibuka(polis.tahap) && (
        <InputOffer
          polisID={polis.id}
          tahap={polis.tahap}
          onSelesai={() => {
            setPolis({ id: '', tahap: '' })
          }}
        />
      )}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanPremiumList> = RutePremiumList
