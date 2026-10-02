// modul/treatyin/frontend/api.ts - panggilan backend modul Treaty In, satu
// fungsi per rute (`backend/handlers/rute_treaty_in.go`). Klien HTTP-nya
// `inti/frontend/klien.ts`.
//
// ⛔ SATU rute, dan itu disengaja. Papan tiket modul ini menyatakan `L-4`:
// "tidak ada spesifikasi layar di mana pun". Rute kontrak, versi, dan
// persetujuan lahir bersama spesifikasinya.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_TREATYIN = '/api/treaty-in'

/** Keenam himpunan acuan tiket 15 - SAMA dengan `models.Himpunan`. */
export const HIMPUNAN_ACUAN = [
  'mata-uang',
  'jenis-potongan',
  'kelas-bisnis',
  'kelompok-treaty',
  'bahaya',
  'jenis-reasuransi',
] as const

export type HimpunanAcuan = (typeof HIMPUNAN_ACUAN)[number]

/**
 * Satu baris tabel acuan - SAMA dengan `models.Acuan`.
 *
 * `idInduk` hanya terisi pada `jenis-reasuransi`, yang bersusun. Pada kelima
 * himpunan lain ia SELALU kosong, dan kosong di sana berarti "tabel ini memang
 * tidak bersusun" - bukan "induknya belum diisi".
 */
export interface Acuan {
  id: number
  kode: string
  nama: string
  aktif: string
  idInduk?: number
}

/** `GET /api/treaty-in/acuan/{himpunan}` - isi satu tabel acuan, urut KODE. */
export async function ambilAcuan(himpunan: HimpunanAcuan): Promise<Acuan[]> {
  return minta<Acuan[]>(`${PREFIX_TREATYIN}/acuan/${himpunan}`)
}
