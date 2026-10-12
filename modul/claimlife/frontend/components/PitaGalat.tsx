// PitaGalat — SATU tempat galat layar Claim Life (keputusan work owner 11-10-2026: "untuk semua pesan error ini tolong
// dirapikan, jangan di bawah seperti sekarang, buat rapi seperti premium list").
//
// Acuan (BACA SAJA, tidak diimpor — isolasi modul): premiumlistlife/frontend/pages/PremiumListDetail.tsx dan
// InputOffer.tsx menaruh `<Gagal>` inti di ATAS isi panel (di bawah judul, sebelum tabel / tombol aksi), dan
// FormDataPolis.tsx menulis kolom wajib sebagai daftar ringkas ("Required: …"). Di sini pola yang sama: satu pita di atas
// isi layar, memakai kelas inti `alert alert--error` / `<Gagal>`, ditambah tombol tutup.
//
// Isinya salah satu dari:
//   - kekurangan Save to RNM (judul + satu baris per peserta);
//   - kalimat teks (layar yang galatnya sudah berupa kalimat);
//   - galat apa pun (ApiFailure / TypeError jaringan) lewat `<Gagal>` inti, yang SUDAH menampilkan pesan server — tanpa
//     salinan `<p>` (bug alert ganda 11-10-2026).

import type { ReactNode } from 'react'

import { Gagal, IkonTutup } from '../../../../inti/frontend/components/ui/dasar'
import { OS_PEGA } from '../labels'
import type { KekuranganRNM } from '../api'

/** Judul kotak kekurangan Save to RNM (keputusan work owner 11-10-2026). */
export const JUDUL_KEKURANGAN = 'Please complete the required fields before Save to RNM:'

/** Judul kotak kekurangan Send to Medical Check (keputusan work owner 11-10-2026). */
export const JUDUL_KEKURANGAN_MEDIS =
  'Please upload at least one claim document for every accepted insured before Send to Medical Check:'

/** Baris kekurangan, satu per peserta: "No 1 · NAMA: Date of Loss, Claim RNM"; tingkat klaim = "Register Claim: …". */
export function barisKekurangan(kekurangan: readonly KekuranganRNM[]): string[] {
  const grup = new Map<number, { label: string; kolom: string[] }>()
  for (const k of kekurangan) {
    let x = grup.get(k.nomorPeserta)
    if (x === undefined) {
      const label = k.nomorPeserta === 0 ? OS_PEGA.tabRegister : `No ${k.nomorPeserta} · ${k.nama || k.sertifikat || '—'}`
      x = { label, kolom: [] }
      grup.set(k.nomorPeserta, x)
    }
    if (!x.kolom.includes(k.kolom)) x.kolom.push(k.kolom)
  }
  return [...grup.values()].map((x) => `${x.label}: ${x.kolom.join(', ')}`)
}

export function PitaGalat({
  galat,
  kekurangan,
  judul,
  onTutup,
}: {
  /** Galat aksi: kalimat, atau objek galat (dibaca `<Gagal>`). null / '' = tidak ada. */
  galat?: unknown
  /** Kekurangan kolom wajib Save to RNM; didahulukan dari `galat`. */
  kekurangan?: readonly KekuranganRNM[] | null
  /** Judul kotak kekurangan (bawaan: Save to RNM). */
  judul?: string
  /** Tombol tutup (×); tanpa ini pita tidak dapat ditutup. */
  onTutup?: () => void
}) {
  let isi: ReactNode = null
  if (kekurangan !== undefined && kekurangan !== null && kekurangan.length > 0) {
    isi = (
      <div role="alert" className="alert alert--error cl-pita__isi">
        <strong>{judul ?? JUDUL_KEKURANGAN}</strong>
        <ul>
          {barisKekurangan(kekurangan).map((b) => (
            <li key={b}>{b}</li>
          ))}
        </ul>
      </div>
    )
  } else if (typeof galat === 'string') {
    if (galat !== '') {
      isi = (
        <div role="alert" className="alert alert--error cl-pita__isi">
          {galat}
        </div>
      )
    }
  } else if (galat !== null && galat !== undefined) {
    isi = <Gagal galat={galat} />
  }
  if (isi === null) return null
  return (
    <div className="cl-pita">
      {isi}
      {onTutup !== undefined && (
        <button type="button" className="cl-pita__tutup" aria-label="Dismiss" title="Dismiss" onClick={onTutup}>
          <IkonTutup ukuran={14} />
        </button>
      )}
    </div>
  )
}
