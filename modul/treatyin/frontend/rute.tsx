// Rute modul Treaty In — tiket 14, 15, dan ronde layar 1.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreatyIn } from './menu'
import type { ModeForm } from './mode'
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
  // ⭐ Mode form — `Edit`/`Add` membuka `ubah`, `View` membuka `lihat`.
  const [mode, setMode] = useState<ModeForm>('lihat')
  // ⭐ Tombol `Copy` — pengenal kontrak SUMBER draf salinan (`TreatyInCopy`:
  // `ID = "UnknownId"`, `OLDID` = sumber). `dibuka` tetap TEKS KOSONG —
  // salinan adalah kontrak baru sampai Save melahirkan pengenalnya.
  const [salinDari, setSalinDari] = useState<string | null>(null)

  return (
    // Akar gaya modul: semua aturan `treatyin.css` diawali `.treatyin` (`display: contents`).
    <div className="treatyin">
      {halaman === 'treatyin-kontrak' && dibuka === null && (
        <DaftarKontrakTreatyIn
          onBuka={(id, m) => {
            setSalinDari(null)
            setMode(m)
            setDibuka(id)
          }}
          onSalin={(id) => {
            // `TreatyInCopy` [4]: `ViewState = 0`, `IsEditData = 0` — form dapat disunting.
            setSalinDari(id)
            setMode('ubah')
            setDibuka('')
          }}
          onTambah={() => {
            setMode('ubah')
            setSalinDari(null)
            setDibuka('')
          }}
        />
      )}
      {halaman === 'treatyin-kontrak' && dibuka !== null && (
        <FormKontrakTreatyIn
          idKontrak={dibuka}
          mode={mode}
          salinDari={salinDari ?? undefined}
          onKembali={() => {
            setSalinDari(null)
            setDibuka(null)
          }}
          // ⭐ Kontrak BARU tersimpan — form dibuka ulang dengan pengenal
          // yang baru lahir, tetap di mode Edit.
          onTersimpan={(id) => {
            setSalinDari(null)
            setDibuka(id)
          }}
        />
      )}
      {halaman === 'treatyin-acuan' && <AcuanTreatyIn />}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyIn> = RuteTreatyIn
