// modul/komiteclaimlife/frontend/api.ts - panggilan backend modul Komite Claim Life, satu fungsi
// per endpoint. Klien HTTP-nya `inti/klien.ts` (refactor bentuk B, 30-09-2026:
// dipecah dari `services/api.ts` tanpa mengubah satu panggilan pun).

import { minta } from '../../../inti/frontend/klien'

// ——— Komite Claim Life tiket 01: Inbox Komite ———

/** Satu baris Inbox Komite. Uang TEKS; status baris KATA. */
export interface BarisInboxKomite {
  kasusId: string
  tglUpdate: string
  statusWork: string
  klaimId: string
  nomorKlaim: string
  tingkatBerjalan: number
  komiteLoop: number
  nilaiKlaim: string
  mataUang: string
  statusBaris: string
}

export interface HalamanInboxKomite {
  baris: BarisInboxKomite[]
  total: number
  halaman: number
  ukuran: number
}

export interface AnggotaKasusKomite {
  urut: number
  jabatan: string
  approval: string
  kataApproval: string
  komentar: string
  tglApprove: string
  saya: boolean
}

/** Satu efek keluar kasus komite — KATA, bukan kode (tiket 08). */
export interface EfekKomite {
  jenis: string
  keadaan: string
  percobaan: number
  sejak: string
}

export interface KasusKomite {
  kasus: BarisInboxKomite
  adjustmentId: string
  tangga: AnggotaKasusKomite[]
  giliranSaya: boolean
  /** Tiket 08 — keadaan ringkas ('' bila belum ada keputusan) dan tiap efek. */
  efek: { keadaan: string; efek: EfekKomite[] }
}

/** Keputusan asli satu tingkat yang tertimpa langkah 5.1 (OQ-K-05). */
export interface KeputusanAsliKomite {
  status: string
  comment: string
}

/** Riwayat tangga — `GET /api/komite/{id}/riwayat` (tiket 09; siapa pun). */
export interface RiwayatKomite {
  kasusId: string
  adjustmentId: string
  tangga: {
    urut: number
    committee: string
    anggota: string
    /** KATA — Setuju / Tolak / Menunggu / Dilewati (eskalasi). */
    status: string
    dateApprove: string
    comment: string
    /**
     * Keputusan tingkat ini SEBELUM langkah 5.1 `KomitePostAdjustment`
     * menimpanya (tingkat akhir menolak → seluruh tingkat `2`, komentar
     * kosong). Dibaca dari jejak — OQ-K-05, GILIRAN-17. Tidak ada bila tidak
     * tertimpa.
     */
    asli?: KeputusanAsliKomite
  }[]
  eskalasi: { dariTingkat: number; keTingkat: number; oleh: string; waktu: string }[]
}

export async function ambilRiwayatKomite(kasusID: string): Promise<RiwayatKomite> {
  return minta<RiwayatKomite>(`/api/komite/${encodeURIComponent(kasusID)}/riwayat`)
}

/** Laporan harian "perlu intervensi" — `GET /api/komite/laporan-harian` (admin). */
export interface LaporanHarianKomite {
  tanggal: string
  kosong: boolean
  baris: (EfekKomite & { kasusId: string })[]
}

export async function ambilLaporanHarianKomite(): Promise<LaporanHarianKomite> {
  return minta<LaporanHarianKomite>('/api/komite/laporan-harian')
}

/** Inbox Komite milik pelaku — `GET /api/komite`. */
export async function ambilInboxKomite(halaman = 1): Promise<HalamanInboxKomite> {
  return minta<HalamanInboxKomite>(`/api/komite?halaman=${String(halaman)}`)
}

/** Jawaban `POST /api/komite/{id}/keputusan`. */
export interface HasilKeputusanKomite {
  tingkatDiputus: number
  keputusan: string
  kataKeputusan: string
  berlanjut: boolean
  tingkatBerikut: number
  akseptasiAkhir: boolean
  tolakAkhir: boolean
  /** Tiket 04a — terisi hanya pada Setuju di tingkat akhir. */
  nomorAkseptasi: string
  /** Tiket 06 — efek keluar yang DIANTRE: tersimpan, belum tuntas. */
  efekTertunda: string[]
}

/**
 * Keputusan satu tingkat — `POST /api/komite/{id}/keputusan`.
 *
 * ⚠️ Menjawab **501** di tingkat akhir sampai akseptasi (tiket 04a/04b) dan
 * penolakan ke baris klaim (tiket 05) dibangun — pesan server tampil apa adanya.
 */
export async function putuskanKomite(
  kasusID: string,
  keputusan: string,
  komentar: string,
): Promise<HasilKeputusanKomite> {
  return minta<HasilKeputusanKomite>(`/api/komite/${encodeURIComponent(kasusID)}/keputusan`, {
    metode: 'POST',
    badan: { keputusan, komentar },
  })
}

/** Eskalasi naik satu tingkat — `POST /api/komite/{id}/eskalasi` (admin). */
export async function eskalasiKomite(
  kasusID: string,
): Promise<{ dariTingkat: number; keTingkat: number }> {
  return minta<{ dariTingkat: number; keTingkat: number }>(
    `/api/komite/${encodeURIComponent(kasusID)}/eskalasi`,
    { metode: 'POST' },
  )
}

/** Satu kasus komite — `GET /api/komite/{id}`. 403 bila bukan anggota tangga. */
export async function ambilKasusKomite(kasusID: string): Promise<KasusKomite> {
  return minta<KasusKomite>(`/api/komite/${encodeURIComponent(kasusID)}`)
}
