// Rute modul Treaty In Adjustment — layar Adjustment, rantai versi (tiket 01
// dan 05), dan Attachment + History kontrak warisan.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreatyInAdjustment } from './menu'
import LampiranKontrak from './pages/LampiranKontrak'
import PenyesuaianKontrak from './pages/PenyesuaianKontrak'
import RantaiVersi from './pages/RantaiVersi'
import './treatyinadjustment.css'

export function RuteTreatyInAdjustment({ halaman, onPindah }: PropsRute<HalamanTreatyInAdjustment>) {
  /**
   * Pengenal kontrak WARISAN yang sedang dipilih.
   *
   * ⛔ Hidup di rute, bukan di salah satu layar — sebab dua layar
   * memakainya dan hanya satu yang memilihnya. Menaruhnya di `RantaiVersi`
   * berarti `LampiranKontrak` tidak dapat membacanya; menaruhnya di
   * keduanya berarti dua kebenaran.
   *
   * ⚠️ Kosong sampai pemakai memilih. Itu keadaan yang JUJUR, bukan cacat:
   * tabel `KONTRAK` model baru masih nol baris sampai tiket 59 memindahkan
   * kepala kontrak warisan, jadi pemilihnya memang belum punya isi.
   */
  const [masterID, setMasterID] = useState('')

  return (
    // Akar gaya modul: semua aturan `treatyinadjustment.css` diawali `.treatyinadjustment`
    // (`display: contents`).
    <div className="treatyinadjustment">
      {halaman === 'treatyinadjustment-penyesuaian' && <PenyesuaianKontrak />}
      {halaman === 'treatyinadjustment-rantai-versi' && (
        <RantaiVersi
          onPilihWarisan={(n) => {
            setMasterID(n)
            // Memilih kontrak langsung membuka lampirannya — satu pemilih,
            // dan nol layar daftar kedua.
            if (n !== '') onPindah('treatyinadjustment-lampiran')
          }}
        />
      )}
      {halaman === 'treatyinadjustment-lampiran' && <LampiranKontrak masterID={masterID} />}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyInAdjustment> = RuteTreatyInAdjustment
