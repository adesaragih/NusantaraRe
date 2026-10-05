// Klien Aggregate - `/api/aggregate` (`modul/aggregate/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `aggregate`; cookie sesi dikirim peramban sendiri.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute - SAMA dengan `handlers.Prefix`. */
export const PREFIX_AG = '/api/aggregate'

/** Batas waktu unggah dan simpan: berkas CSV sampai 5.000 baris. */
export const BATAS_WAKTU_UNGGAH_MS = 120_000

/** Satu baris AGGREGATE: kolom Oracle -> teks (angka bertitik desimal tanpa pemisah ribuan, tanggal DD-MM-YYYY). */
export type Baris = Record<string, string>

/** Kunci satu baris daftar - `models.Kunci`. */
export interface Kunci {
  tanggalInput: string
  cedingCode: string
  cedingName: string
  treatyType: string
  asAt: string
  uwYear: string
}

/** Satu baris daftar - `models.Kelompok`. */
export interface Kelompok extends Kunci {
  jumlahBaris: number
  inputTerakhir: string
}

export interface Halaman {
  daftar: Kelompok[]
  total: number
  halaman: number
  ukuran: number
}

/** Satu baris view Master ID - `models.MasterTreaty`. */
export interface MasterTreaty {
  id: string
  treatyId: string
  cedingId: string
  ceding: string
  rnmShare: string
  treatyYear: string
  proportionType: string
  treatyContractName: string
  treatyGroup: string
  sob: string
}

/** Satu daun chart - `models.IrisanRingkasan`: jumlah RNM Value (USD) satu Ceding, Treaty Type, Coverage. */
export interface IrisanRingkasan {
  cedingCode: string
  cedingName: string
  treatyType: string
  coverage: string
  rnmValueInUsd: string
}

/** Isi chart bertingkat - `services.Ringkasan`: seluruh As At dijumlah. */
export interface Ringkasan {
  irisan: IrisanRingkasan[]
}

/** Hasil Upload CSV - `services.Pratinjau`. */
export interface Pratinjau {
  baris: Baris[]
  masterTreaty: MasterTreaty[]
}

export interface HasilSimpan {
  disimpan: number
  pesan: string
}

function kueriKunci(k: Kunci): Record<string, string> {
  return { ...k }
}

export function ambilDaftar(q: string, halaman: number, ukuran?: number): Promise<Halaman> {
  return minta<Halaman>(PREFIX_AG, { kueri: { q: q === '' ? undefined : q, halaman, ukuran } })
}

/** Chart: seluruh As At (keputusan work owner 04-10-2026). */
export function ambilRingkasan(): Promise<Ringkasan> {
  return minta<Ringkasan>(`${PREFIX_AG}/ringkasan`)
}

export function ambilRincian(k: Kunci): Promise<{ baris: Baris[] }> {
  return minta<{ baris: Baris[] }>(`${PREFIX_AG}/rincian`, { kueri: kueriKunci(k) })
}

export function hapus(k: Kunci): Promise<{ dihapus: number }> {
  return minta<{ dihapus: number }>(`${PREFIX_AG}/hapus`, { metode: 'POST', badan: k })
}

/** Popup Master ID; kata kosong = nol baris. */
export function cariMasterTreaty(q: string): Promise<{ daftar: MasterTreaty[] }> {
  return minta<{ daftar: MasterTreaty[] }>(`${PREFIX_AG}/master-treaty`, { kueri: { q } })
}

/** Upload CSV: isi berkas (teks) dan Master ID terpilih (ID baris view, urutan pilihan). */
export function pratinjau(csv: string, masterTreaty: string[]): Promise<Pratinjau> {
  return minta<Pratinjau>(`${PREFIX_AG}/pratinjau`, {
    metode: 'POST',
    badan: { csv, masterTreaty },
    batasWaktuMs: BATAS_WAKTU_UNGGAH_MS,
  })
}

/** Save: grid pratinjau apa adanya (baris "Total :" dilewati backend). */
export function simpan(baris: Baris[]): Promise<HasilSimpan> {
  return minta<HasilSimpan>(`${PREFIX_AG}/simpan`, {
    metode: 'POST',
    badan: { baris },
    batasWaktuMs: BATAS_WAKTU_UNGGAH_MS,
  })
}
