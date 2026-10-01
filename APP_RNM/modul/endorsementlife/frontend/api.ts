// modul/endorsementlife/frontend/api.ts - panggilan backend modul Endorsement Life, satu fungsi per rute
// (`backend/handlers/rute_edm.go`). Klien HTTP-nya `inti/klien.ts`.
//
// ⛔ Uang dan angka TEKS sepanjang jalan - tidak pernah `Number` (ADR-0003). Kosong = `""`.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_EDM = '/api/endorsement-life'

/** Satu halaman grid - `services.Halaman`. */
export interface HalamanEDM<T> {
  baris: T[]
  total: number
  halaman: number
  ukuran: number
}

/** Satu baris kotak masuk - `models.BarisInbox`. */
export interface BarisInbox {
  caseId: string
  endorsementNo: string
  type: string
  edmType: string
  policyNo: string
  sob: string
  ceding: string
  policyHolder: string
  marketingName: string
  createDate: string
  createOperator: string
  status: string
}

/** Kepala kasus - `models.Kasus`. `kepala` berkunci nama kolom `T_PREMIUM_LIST`. */
export interface KasusEDM {
  id: string
  policyNo: string
  plNumber: string
  plNumberEdm: string
  prodKe: number
  edmType: string
  description: string
  edmDate: string
  status: string
  createDate: string
  createOperator: string
  kepala: Record<string, string>
  sudahSimpan: boolean
  csvTerkunci: boolean
  cacah: Record<string, number>
  sumber: 'aplikasi' | 'warisan' | ''
}

/** Satu peserta - `models.Peserta`. `nilai` berkunci nama kolom `T_PREMIUM_LIST_DETAIL`. */
export interface PesertaEDM {
  id: string
  parentId: string
  edmStatus: string
  terkunci: boolean
  nilai: Record<string, string>
}

/** Kotak masuk - `GET /inbox`. */
export function ambilInbox(halaman: number): Promise<HalamanEDM<BarisInbox>> {
  return minta<HalamanEDM<BarisInbox>>(`${PREFIX_EDM}/inbox`, { kueri: { halaman } })
}

/** Kepala kasus - `GET /kasus/{id}`. */
export function ambilKasus(id: string): Promise<KasusEDM> {
  return minta<KasusEDM>(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}`)
}

/** Peserta kasus - `GET /kasus/{id}/peserta`. */
export function ambilPeserta(id: string, halaman: number): Promise<HalamanEDM<PesertaEDM>> {
  return minta<HalamanEDM<PesertaEDM>>(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/peserta`, { kueri: { halaman } })
}
