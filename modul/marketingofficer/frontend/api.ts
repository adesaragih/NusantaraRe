// Klien Marketing Officer - `/api/marketing-officer` (`modul/marketingofficer/backend/handlers/rute.go`). Rutenya
// hanya terbuka bagi pemegang menu `marketingofficer`; cookie sesi dikirim peramban sendiri.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute - SAMA dengan `handlers.Prefix`. */
export const PREFIX_MO = '/api/marketing-officer'

/** Satu baris `POOLDATA.MARKETINGOFFICER` - `models.MarketingOfficer`. Semua teks; kosong = "". */
export interface MarketingOfficer {
  id: string
  /** Marketing Code - tidak pernah berubah sesudah dibuat. */
  clientId: string
  clientName: string
  /** `LEADER` bila baris ini leader, selain itu ID baris leadernya. */
  clientId2: string
  moLeader: string
  /** `1` aktif, `2` nonaktif. */
  moStatus: string
  branchParent: string
  branchDetailId: string
  branchDetailName: string
  teamGroup: string
  branchStatus: string
  /** `M_LOGIN_GO.LOGIN_ID`. */
  aksesLogin: string
  userUpdate: string
  /** `YYYY-MM-DD HH:MI` jam Oracle. */
  tanggal: string
}

/** Status akun sebuah baris: kosong = tanpa akun; `tidak-ada` = Operator ID lama Pega yang tidak ada di M_LOGIN_GO. */
export type StatusAkun = '' | 'aktif' | 'nonaktif' | 'tidak-ada'

/** Satu baris daftar - `services.BarisMO`. */
export interface BarisMO extends MarketingOfficer {
  statusAkun: StatusAkun
  emailAkun: string
}

export interface Akun {
  loginId: string
  nama: string
  contactId: string
  email: string
  aktif: boolean
}

export interface Cabang {
  id: string
  nama: string
  /** `BRANCHPARENTID`. */
  induk: string
  kanwilGroup: string
}

export interface OpsiLeader {
  id: string
  nama: string
  subBranch: string
}

/** Pilihan form - `services.Pilihan`. */
export interface Pilihan {
  akun: Akun[]
  leader: OpsiLeader[]
  branch: Cabang[]
  subBranch: Cabang[]
}

/** Isian form tambah dan ubah - `models.Isian`. */
export interface Isian {
  aksesLogin: string
  leader: boolean
  leaderId: string
  branchParent: string
  branchDetailId: string
  aktif: boolean
}

export interface Daftar {
  daftar: BarisMO[]
  total: number
}

const e = encodeURIComponent

export function ambilDaftar(): Promise<Daftar> {
  return minta<Daftar>(PREFIX_MO)
}

export function ambilPilihan(): Promise<Pilihan> {
  return minta<Pilihan>(`${PREFIX_MO}/pilihan`)
}

/** Tambah - INSERT; ID dari backend. */
export function tambah(isi: Isian): Promise<MarketingOfficer> {
  return minta<MarketingOfficer>(PREFIX_MO, { metode: 'POST', badan: isi })
}

/** Ubah - UPDATE; nonaktif = `aktif: false` (tidak ada hapus). */
export function ubah(id: string, isi: Isian): Promise<MarketingOfficer> {
  return minta<MarketingOfficer>(`${PREFIX_MO}/${e(id)}`, { metode: 'PUT', badan: isi })
}

/** Satu kolom yang berubah - `services.RuasBerubah`; `kolom` = nama kolom MARKETINGOFFICER. */
export interface RuasBerubah {
  kolom: string
  sebelum: string
  sesudah: string
}

/** Satu UPDATE - `services.Perubahan`. `perkiraan` = baris log lama tanpa LOG_TIME (urutan dan waktu perkiraan). */
export interface Perubahan {
  waktu: string
  oleh: string
  perkiraan: boolean
  ruas: RuasBerubah[]
}

/** Log perubahan satu MO dari MARKETINGOFFICER_LOG, terbaru dulu - `services.Riwayat`. */
export interface Riwayat {
  id: string
  perubahan: Perubahan[]
  jumlahLog: number
}

export function ambilLog(id: string): Promise<Riwayat> {
  return minta<Riwayat>(`${PREFIX_MO}/${e(id)}/log`)
}
