// Tata letak dua kolom layar kasus, mengikuti screenshot layar Pega (perintah work owner 05-10-2026):
//
//   General : tombol Choose Business / Select SOB / Enable-Disable paling atas; kolom kiri Master ID .. RNM Share,
//             kolom kanan Statement Date .. Payment Type/Layer, Remark melebar di bawah; "RNM Share IDR 163.125.000"
//             satu baris; "Q [..] / [..] U/Y [..]" satu baris; NonProp "layer 1 Of layer 1" satu baris
//             (screenshot work owner 06-10-2026).
//   Uang    : baris Gross di atas (Gross Premium 100% | Claim 100%), kolom kiri OGP lalu klaim/saldo,
//             kolom kanan ONP lalu potongan/pajak. Kode mata uang polis TIDAK diulang di setiap medan (sudah
//             tampil di medan Currency) - "banyak sekali tulisan IDR".
//
// Hanya MENATA medan definisi `medan.ts` - nol medan baru, nol medan hilang, urutan asli dipertahankan.

import { POLIS } from './api'
import { BAGIAN } from './labels'
import { MEDAN_ADMIN_UANG, MEDAN_ATASAN_TOTAL, MEDAN_ATASAN_UANG, type Kelompok, type Medan } from './medan'

/** Medan section General per kolom. */
export interface TataUmum {
  kiri: Medan[]
  kanan: Medan[]
  bawah: Medan[]
}

/** Bagian uang: baris atas dua sel, lalu dua kolom berkelompok. */
export interface TataUang {
  atasKiri: Medan[]
  atasKanan: Medan[]
  kiri: Kelompok[]
  kanan: Kelompok[]
}

const AWAL_KANAN = POLIS + 'StatementDate'
const AWAL_BAWAH = POLIS + 'Remark'
const KODE_SHARE = POLIS + 'ShareCurrency'
const NILAI_SHARE = POLIS + 'ShareValue'
const MATA_UANG_POLIS = POLIS + 'Currency'

/** "RNM Share [IDR] [163.125.000]": baris kode ShareCurrency dilebur ke ShareValue yang memakai labelnya. */
function satuBarisShare(ms: Medan[]): Medan[] {
  const kode = ms.find((m) => m.jalur === KODE_SHARE)
  if (kode === undefined || !ms.some((m) => m.jalur === NILAI_SHARE)) return ms
  return ms.filter((m) => m !== kode).map((m) => (m.jalur === NILAI_SHARE ? { ...m, label: kode.label } : m))
}

/** Pecah medan General (urutan sel Pega) menjadi kolom kiri, kanan, dan baris bawah. */
export function tataUmum(semua: Medan[]): TataUmum {
  const ms = satuBarisShare(semua)
  const i = ms.findIndex((m) => m.jalur === AWAL_KANAN)
  const j = ms.findIndex((m) => m.jalur === AWAL_BAWAH)
  const k = i < 0 ? ms.length : i
  const b = j < 0 ? ms.length : j
  return { kiri: ms.slice(0, k), kanan: ms.slice(k, b), bawah: ms.slice(b) }
}

const DERET_Q = [POLIS + 'Quartal', POLIS + 'YearOfQuartal', POLIS + 'TreatyYear']
/** Medan layer NonProp: sel tanpa label, kecuali LABEL "Of" di depan LayerPartType. */
const DERET_LAYER = [POLIS + 'LayerType', POLIS + 'Layer', POLIS + 'LayerPartType', POLIS + 'LayerPart']
const TANPA_LABEL = new Set([POLIS + 'LayerType', POLIS + 'Layer', POLIS + 'LayerPart'])

/** Deret layer NonProp ("layer 1 Of layer 1"). */
export const deretLayer = (d: Medan[]) => d[0]?.jalur === DERET_LAYER[0]

/** Satukan "Q / U/Y" (tiga medan berurutan) dan layer NonProp (empat medan berurutan, label nama properti dibuang)
 *  menjadi satu deret; medan lain tetap sendiri. */
export function deretQ(ms: Medan[]): (Medan | Medan[])[] {
  const hasil: (Medan | Medan[])[] = []
  for (let i = 0; i < ms.length; i++) {
    const pola = [DERET_Q, DERET_LAYER].find((p) => {
      const potong = ms.slice(i, i + p.length)
      return potong.length === p.length && potong.every((m, n) => m.jalur === p[n])
    })
    if (pola === undefined) {
      hasil.push(ms[i]!)
      continue
    }
    hasil.push(ms.slice(i, i + pola.length).map((m) => (TANPA_LABEL.has(m.jalur) ? { ...m, label: '' } : m)))
    i += pola.length - 1
  }
  return hasil
}

/** Potong `ms` menurut jalur: dari medan `dari` (inklusif) sampai sebelum `sampai`; kode mata uang polis dilepas. */
function antara(ms: Medan[], dari: string, sampai?: string): Medan[] {
  const i = ms.findIndex((m) => m.jalur === POLIS + dari)
  const j = sampai === undefined ? ms.length : ms.findIndex((m) => m.jalur === POLIS + sampai)
  return ms.slice(i, j).map((m) => (m.mataUang === MATA_UANG_POLIS ? { ...m, mataUang: undefined } : m))
}

export const TATA_UANG_ADMIN: TataUang = {
  atasKiri: antara(MEDAN_ADMIN_UANG, 'GrossPremium', 'GrossClaim'),
  atasKanan: antara(MEDAN_ADMIN_UANG, 'GrossClaim', 'PremiOgp'),
  kiri: [
    { judul: BAGIAN.ogp, medan: antara(MEDAN_ADMIN_UANG, 'PremiOgp', 'PremiOnp') },
    { medan: antara(MEDAN_ADMIN_UANG, 'Claim', 'Deduction1') },
  ],
  kanan: [
    { judul: BAGIAN.onp, medan: antara(MEDAN_ADMIN_UANG, 'PremiOnp', 'Claim') },
    { medan: antara(MEDAN_ADMIN_UANG, 'Deduction1') },
  ],
}

export const TATA_UANG_ATASAN: TataUang = {
  atasKiri: antara(MEDAN_ATASAN_UANG, 'GrossPremium', 'PremiOgp'),
  atasKanan: [],
  kiri: [{ judul: BAGIAN.ogp, medan: antara(MEDAN_ATASAN_UANG, 'PremiOgp', 'PremiOnp') }],
  kanan: [{ judul: BAGIAN.onp, medan: antara(MEDAN_ATASAN_UANG, 'PremiOnp') }],
}

/** Total di bawah grid spreading layar atasan: kolom premium | kolom klaim. */
export const TOTAL_ATASAN: [Medan[], Medan[]] = [
  antara(MEDAN_ATASAN_TOTAL, 'TotalSharePercentagePremium', 'TotalSharePercentageClaim'),
  antara(MEDAN_ATASAN_TOTAL, 'TotalSharePercentageClaim'),
]
