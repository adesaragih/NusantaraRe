// modul/treatyinadjustment/frontend/api.ts - panggilan backend modul Treaty In
// Adjustment (`backend/handlers/rute_treaty_in_adjustment.go`).
//
// ⛔ Dua rute, keduanya BACA. Jalur simpan lahir bersama spesifikasi layarnya.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_TREATYINADJUSTMENT = '/api/treaty-in-adjustment'

/** Kepala kontrak - lapisan BEKU yang seluruh versinya bagi (ADR-0040). */
export interface Kontrak {
  id: number
  nomorKontrakWarisan: string
  sifatProporsi: string
  tanggalMulai: string
  tanggalBerakhir: string
}

/**
 * Satu baris rantai versi - SAMA dengan `models.Versi`.
 *
 * `nomorUrutVersi` null = baris warisan yang belum dinomori ulang (tiket 10),
 * BUKAN nol. `idVersiDasar` null = versi PERTAMA kontrak itu, dan itu keadaan
 * yang benar, bukan data yang hilang.
 */
export interface Versi {
  id: number
  idKontrak: number
  nomorUrutVersi: number | null
  keadaanSiklusHidup: string
  jenisAddendum: string
  sifatMaterialAddendum: string
  tanggalBerlakuAddendum: string
  idVersiDasar: number | null
  namaKontrak: string
}

/** `GET /api/treaty-in-adjustment/kontrak` - kepala seluruh kontrak. */
export async function ambilKontrak(): Promise<Kontrak[]> {
  return minta<Kontrak[]>(`${PREFIX_TREATYINADJUSTMENT}/kontrak`)
}

/** `GET /api/treaty-in-adjustment/kontrak/{id}/versi` - rantai versi satu kontrak. */
export async function ambilRantaiVersi(idKontrak: number): Promise<Versi[]> {
  return minta<Versi[]>(`${PREFIX_TREATYINADJUSTMENT}/kontrak/${idKontrak}/versi`)
}
