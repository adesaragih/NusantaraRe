// Klien API layar master generik (inti, keputusan work owner 04-10-2026 K3 "Pindah ke inti"). Setiap modul master
// memanggilnya dengan PREFIX rutenya sendiri (`/api/master-nation`, `/api/master-province`, …) - kontrak "kontrak rute
// 8-modul" sesi c3:
//   GET P/meta · GET P?q=&status=&halaman= · POST P · PUT P/{id} · PUT P/{id}/status · GET P/rujukan/{kunci}?q=&halaman=
// Seluruh nilai baris = teks; status aktif = boolean `aktif`.

import { minta } from '../klien'

/** Satu kolom master. */
export interface KolomMaster {
  /** Kunci JSON baris. */
  kunci: string
  /** Nama kolom tabel. */
  kolom: string
  /** Lebar kolom (byte). */
  lebar: number
  wajib: boolean
  /** Diisi backend (rujukan / jejak ubah) - tidak menjadi isian. */
  turunan: boolean
}

/** Satu kolom rujukan: isian `kunci` diisi `nilai` baris master `judul`; saran menampilkan `nama`. */
export interface RujukanMaster {
  kunci: string
  judul: string
  nilai: string
  nama: string
}

/** `GET P/meta`. */
export interface MetaMaster {
  kunci: string
  judul: string
  /** ID dibuat backend (Accumulation). */
  idOtomatis: boolean
  kolom: KolomMaster[]
  rujukan: RujukanMaster[]
}

/** Satu baris master: kunci kolom -> teks, plus status. */
export type BarisMaster = { aktif: boolean } & Record<string, string | boolean>

export interface HalamanMaster {
  baris: BarisMaster[]
  total: number
  halaman: number
  ukuran: number
}

/** Saringan status daftar. */
export type StatusSaring = '' | 'aktif' | 'nonaktif'

/** Klien satu modul master. */
export interface KlienMaster {
  meta(): Promise<MetaMaster>
  cari(q: string, status: StatusSaring, halaman: number): Promise<HalamanMaster>
  tambah(isi: Record<string, string>): Promise<{ id: string }>
  ubah(id: string, isi: Record<string, string>): Promise<{ id: string }>
  ubahStatus(id: string, aktif: boolean): Promise<{ id: string; aktif: boolean }>
  /** Baris master yang dirujuk kolom `kunci` (aktif saja). */
  rujukan(kunci: string, q: string, halaman: number): Promise<HalamanMaster>
}

/** Klien untuk prefix rute satu modul (mis. `/api/master-province`). */
export function klienMaster(prefix: string): KlienMaster {
  const id = (v: string) => `${prefix}/${encodeURIComponent(v)}`
  return {
    meta: () => minta<MetaMaster>(`${prefix}/meta`),
    cari: (q, status, halaman) => {
      const kueri: Record<string, string | number> = { halaman }
      if (q.trim() !== '') kueri.q = q.trim()
      if (status !== '') kueri.status = status
      return minta<HalamanMaster>(prefix, { kueri })
    },
    tambah: (isi) => minta<{ id: string }>(prefix, { metode: 'POST', badan: isi }),
    ubah: (v, isi) => minta<{ id: string }>(id(v), { metode: 'PUT', badan: isi }),
    ubahStatus: (v, aktif) => minta<{ id: string; aktif: boolean }>(`${id(v)}/status`, { metode: 'PUT', badan: { aktif } }),
    rujukan: (kunci, q, halaman) => {
      const kueri: Record<string, string | number> = { halaman }
      if (q.trim() !== '') kueri.q = q.trim()
      return minta<HalamanMaster>(`${prefix}/rujukan/${encodeURIComponent(kunci)}`, { kueri })
    },
  }
}
