// Pesan medan "This name is on Agent Negative List" — padanan
// `Activity/TreatyInCheckCedingBlacklist` (Property-Set-Messages pada
// `TreatyIn.Ceding` [2] dan `TreatyIn.LeadingReinsSource` [3]).
//
// ⭐ KAPAN — dibaca dari ekspor: tombol `Edit` layar daftar menjalankannya
// SESUDAH `SetTreatyIn_Act` (`Section/InputTreatyInOffer.xml`; `View`, `Copy`,
// `Revision` tidak). Jadi di sini: SEKALI, saat kontrak lama selesai dimuat
// di mode ubah — bukan di mode lihat, bukan untuk kontrak baru.
//
// ⛔ TIDAK dijalankan ulang saat memilih lewat "Choose …": tombol `Choose`
// jendela `TreatyInSearchReinsured` / `TreatyInSearchSoB` hanya menjalankan
// DataTransform `TreatyInSetReinsured`, tanpa activity ini. Pesan sebuah
// medan HILANG begitu pengenalnya berganti — daftar pilihan itu sendiri
// hanya menawarkan agen `STATUSACTIVE = '1'` (`BrowseAgentNusaRe_RD`
// saringan A), jadi agen yang baru dipilih tidak pernah berpesan.
//
// ⭐ TIDAK menghalangi apa pun: `TreatyIn.pyErrMsg = "error"` (langkah 4)
// tidak dibaca aturan mana pun di ekspor. Aturan lengkap dan temuan
// "tidak pernah menyala" ada di `backend/services/agen_daftar_negatif.go`.

import { useEffect, useRef, useState } from 'react'

import { periksaDaftarNegatifAgen, type HasilDaftarNegatifAgen } from '../api'

/** Pesan per medan — `undefined` = tanpa pesan. */
export interface PesanAgen {
  cedant?: string
  asalBisnis?: string
}

/**
 * Pesan yang tampil untuk pengenal KINI. Murni.
 *
 * Pesan hanya milik pengenal yang diperiksa; begitu medan berganti agen,
 * pesannya lepas.
 */
export function pesanUntukPengenal(hasil: HasilDaftarNegatifAgen | null, idCedant: string, idAsalBisnis: string): PesanAgen {
  if (hasil === null) return {}
  const keluar: PesanAgen = {}
  if (hasil.cedant.daftarNegatif && hasil.cedant.id === idCedant) keluar.cedant = hasil.cedant.pesan
  if (hasil.asalBisnis.daftarNegatif && hasil.asalBisnis.id === idAsalBisnis) keluar.asalBisnis = hasil.asalBisnis.pesan
  return keluar
}

/**
 * Menjalankan pemeriksaan sekali tiap kali `aktif` menjadi benar (atau
 * `kunci` — pengenal kontrak — berganti selagi aktif), atas pengenal yang
 * sedang dipegang form saat itu.
 */
export function useDaftarNegatifAgen(aktif: boolean, kunci: string, idCedant: string, idAsalBisnis: string): PesanAgen {
  const [hasil, setHasil] = useState<HasilDaftarNegatifAgen | null>(null)
  const pengenal = useRef({ idCedant, idAsalBisnis })
  pengenal.current = { idCedant, idAsalBisnis }
  useEffect(() => {
    setHasil(null)
    if (!aktif) return
    const { idCedant: c, idAsalBisnis: s } = pengenal.current
    if (c.trim() === '' && s.trim() === '') return
    let dibuang = false
    periksaDaftarNegatifAgen(c, s)
      .then((h) => {
        if (!dibuang) setHasil(h)
      })
      // ⚠️ Pemeriksaan PERINGATAN, bukan gerbang: gagal membaca statusnya
      // tidak boleh menghalangi form. Tanpa jawaban = tanpa pesan.
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [aktif, kunci])
  return pesanUntukPengenal(hasil, idCedant, idAsalBisnis)
}

/** Pesan medan di bawah kotaknya — `Property-Set-Messages` tingkat medan. */
export function PesanMedanAgen({ pesan }: { pesan?: string }) {
  if (!pesan) return null
  return (
    <p className="field__error" role="alert">
      {pesan}
    </p>
  )
}
