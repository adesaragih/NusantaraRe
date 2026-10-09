// Rute modul Treaty In Adjustment — layar Adjustment, rantai versi (tiket 01
// dan 05), dan Attachment + History kontrak warisan.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import DaftarPilihBawah from './komponen/DaftarPilihBawah'
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

  // ⛔ BAHASA INGGRIS untuk seluruh teks bawaan komponen bersama —
  // permintaan pemilik proses 8 Oktober 2026: *"gunakan bahasa inggris untuk
  // semua label dan text yang ada di treaty in dan treaty in adjustment"*.
  //
  // ⭐ Lewat `BahasaUI.Provider`, mekanisme yang SUDAH dipakai modul lain
  // (Accounts, Aggregate, Treaty Contract Out) — bukan tiruannya. Tanpa itu
  // `Memuat...`, `-- pilih --`, `Type to filter`, `Tidak ada baris`
  // dan teks pager muncul berbahasa Indonesia di tengah layar Inggris.
  return (
    // Akar gaya modul: semua aturan `treatyinadjustment.css` diawali `.treatyinadjustment`
    // (`display: contents`).
    <BahasaUI.Provider value="en">
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
      {/* ⭐ Daftar `<select>` selalu terbuka KE BAWAH (8 Oktober 2026) —
          satu pemasangan untuk seluruh modul; lihat komponennya. */}
      <DaftarPilihBawah akar=".treatyinadjustment" />
      </div>
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyInAdjustment> = RuteTreatyInAdjustment
