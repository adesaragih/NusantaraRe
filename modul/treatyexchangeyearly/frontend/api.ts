// Klien Treaty Exchange Yearly - `/api/treaty-exchange-yearly` (`modul/treatyexchangeyearly/backend/handlers/rute.go`).
// Rutenya hanya terbuka bagi pemegang menu `treatyexchangeyearly`; tulis ditolak bagi akses View only. Tanpa hapus.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_TEY = '/api/treaty-exchange-yearly'

/** Satu baris - `models.Kurs`. */
export interface Kurs {
  /** ROWID baris - ID warisan tidak unik, jadi Edit memakai ini. */
  kunci: string
  id: string
  treatyYear: string
  idCurrency: string
  currency: string
  currencyName: string
  /** Teks format Pega `YYYYMMDDTHHMMSS.mmm GMT`. */
  startDate: string
  endDate: string
  /** Tampilan `DD-MM-YYYY`. */
  mulai: string
  akhir: string
  toIdr: string
  toUsd: string
  quarter: string
  userId: string
  dateIu: string
  dateIn: string
  /** Tampilan DATEIU WIB `DD-MM-YYYY HH:MM`. */
  diubah: string
}

/** Satu mata uang - `models.MataUang` (view CURRENCY). */
export interface MataUang {
  id: string
  kode: string
  nama: string
}

/** Isian form - `models.Isian`. Kunci kosong = Add. Tanggal `YYYY-MM-DD`. */
export interface Isian {
  kunci: string
  treatyYear: string
  idCurrency: string
  startDate: string
  endDate: string
  toIdr: string
  toUsd: string
  quarter: string
}

export function ambilDaftar(q: string, tahun: string): Promise<{ daftar: Kurs[] }> {
  return minta(PREFIX_TEY, {
    kueri: { q: q.trim() === '' ? undefined : q.trim(), tahun: tahun === '' ? undefined : tahun },
  })
}

export function ambilPilihan(): Promise<{ mataUang: MataUang[]; tahun: string[] }> {
  return minta(`${PREFIX_TEY}/pilihan`)
}

/** Add - ID dibuat backend (situs aktif + TREATYEXCHANGE_SEQ). */
export function tambah(isi: Isian): Promise<Kurs> {
  return minta(PREFIX_TEY, { metode: 'POST', badan: { ...isi, kunci: '' } })
}

/** Edit - baris dikenali `kunci` (ROWID) di badan. */
export function ubah(isi: Isian): Promise<Kurs> {
  return minta(PREFIX_TEY, { metode: 'PUT', badan: isi })
}
