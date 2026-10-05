// Klien Template Manager - `/api/templat` (`inti/backend/templat/rute`). Kelola: pemegang menu Template Manager;
// unduh versi aktif: juga pemegang menu pemilik slot.

import { minta, mintaFormulir, unduhBerkasBeridentitas } from '../klien'

export const PREFIX_TEMPLAT = '/api/templat'

/** Satu versi berkas yang pernah diunggah - `templat.Versi`. */
export interface Versi {
  versi: number
  namaBerkas: string
  ukuran: number
  jumlahKolom: number
  catatan: string
  aktif: boolean
  diunggahOleh: string
  /** `DD-MM-YYYY HH24:MI`. */
  tglUnggah: string
}

/** Satu slot templat - `templat.RingkasanSlot`. */
export interface Slot {
  kode: string
  menu: string
  grup: string
  nama: string
  dipakaiDi: string
  ekstensi: string
  pemisah: string
  jumlahKolom: number
  namaUnduhan: string
  /** Versi aktif; null = berkas bawaan. */
  aktif: Versi | null
}

export interface Perbedaan {
  kolom: number
  lama: string
  baru: string
}

/** Hasil pemeriksaan berkas sebelum disimpan - `templat.HasilPeriksa`. */
export interface HasilPeriksa {
  galat: string[]
  jumlahKolom: number
  perbedaan: Perbedaan[]
}

const jalurSlot = (kode: string) => `${PREFIX_TEMPLAT}/${encodeURIComponent(kode)}`

export function ambilDaftar(): Promise<{ daftar: Slot[] }> {
  return minta<{ daftar: Slot[] }>(PREFIX_TEMPLAT)
}

export function ambilRiwayat(kode: string): Promise<{ versi: Versi[] }> {
  return minta<{ versi: Versi[] }>(`${jalurSlot(kode)}/riwayat`)
}

function formulir(berkas: File, catatan?: string): FormData {
  const f = new FormData()
  f.append('berkas', berkas)
  if (catatan !== undefined) f.append('catatan', catatan)
  return f
}

export function periksa(kode: string, berkas: File): Promise<HasilPeriksa> {
  return mintaFormulir<HasilPeriksa>(`${jalurSlot(kode)}/periksa`, formulir(berkas))
}

export function unggah(kode: string, berkas: File, catatan: string): Promise<{ versi: number }> {
  return mintaFormulir<{ versi: number }>(`${jalurSlot(kode)}/unggah`, formulir(berkas, catatan))
}

/** Versi 0 = kembali ke berkas bawaan. */
export function aktifkan(kode: string, versi: number): Promise<{ versi: number }> {
  return minta<{ versi: number }>(`${jalurSlot(kode)}/aktifkan`, { metode: 'POST', badan: { versi } })
}

/**
 * Unduh templat - dipakai tombol Template di setiap menu (mis. Bordereaux, Aggregate) dan di Template Manager.
 * Tanpa `versi` = versi aktif (atau bawaan); dengan `versi` = satu versi riwayat (hanya pengelola).
 */
export function unduhTemplat(kode: string, namaBerkas: string, versi?: number): Promise<void> {
  const kueri = versi === undefined ? '' : `?versi=${versi}`
  return unduhBerkasBeridentitas(`${jalurSlot(kode)}/unduh${kueri}`, namaBerkas)
}
