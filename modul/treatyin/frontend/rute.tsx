// Rute modul Treaty In — tiket 14, 15, dan ronde layar 1.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreatyIn } from './menu'
import AcuanTreatyIn from './pages/AcuanTreatyIn'
import DaftarKontrakTreatyIn from './pages/DaftarKontrakTreatyIn'
import FormKontrakTreatyIn from './pages/FormKontrakTreatyIn'
import './treatyin.css'

export function RuteTreatyIn({ halaman }: PropsRute<HalamanTreatyIn>) {
  // Kontrak yang sedang dibuka di form; null = daftar.
  //
  // ⛔ Pengenalnya TEKS sejak layar daftar membaca `POOLDATA.TREATY_IN`:
  // kolom `TREATY_IN.ID` adalah `VARCHAR2(100)`. Teks KOSONG berarti kontrak
  // BARU (tombol `Add`), dan ia dibedakan dari `null` supaya "belum memilih"
  // tidak tertukar dengan "memilih yang baru".
  const [dibuka, setDibuka] = useState<string | null>(null)

  return (
    // Akar gaya modul: semua aturan `treatyin.css` diawali `.treatyin` (`display: contents`).
    <div className="treatyin">
      {halaman === 'treatyin-kontrak' && dibuka === null && (
        <DaftarKontrakTreatyIn
          onBuka={setDibuka}
          onTambah={() => {
            setDibuka('')
          }}
        />
      )}
      {halaman === 'treatyin-kontrak' && dibuka !== null && (
        <FormKontrakTreatyIn
          idKontrak={dibuka}
          onKembali={() => {
            setDibuka(null)
          }}
        />
      )}
      {halaman === 'treatyin-acuan' && <AcuanTreatyIn />}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyIn> = RuteTreatyIn
