// modul/endorsementlife/frontend/api.ts - panggilan backend modul Endorsement Life, satu fungsi per rute
// (`backend/handlers/rute_edm.go`). Klien HTTP-nya `inti/klien.ts`.
//
// ⛔ Uang dan angka TEKS sepanjang jalan - tidak pernah `Number` (ADR-0003). Kosong = `""`.

import { minta } from '../../../inti/frontend/klien'
import type { PilihanHapusEDM } from './tampilan'

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
  /** Rekap mata uang `T_PREMIUM_LIST_SUMMARY` - kolom `models.KolomRekapKasus`. */
  rekap: Array<Record<string, string>>
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

/** Jawaban gerbang - `models.Kelayakan`. */
export interface KelayakanEDM {
  pesan: string[]
  boleh: boolean
}

/** Isian `EndorsmentLife_Section` - `services.MasukanKasus`. Tanggal `YYYY-MM-DD`. */
export interface MasukanKasusEDM {
  policyNo: string
  edmType: string
  edmDate: string
  description: string
}

/** Kasus baru - `services.HasilBuat`. */
export interface HasilBuatEDM {
  id: string
  peserta: number
  spreading: number
  spreadingRetro: number
}

/** Gerbang kelayakan tanpa tulis - `POST /kelayakan` (`SetErrorBatalEndorsement_Act`). */
export function cekKelayakan(policyNo: string, edmType: string): Promise<KelayakanEDM> {
  return minta<KelayakanEDM>(`${PREFIX_EDM}/kelayakan`, { metode: 'POST', badan: { policyNo, edmType } })
}

/** Buat kasus - `POST /kasus` (`Submit` b4226 → `MappingEDMLife`). */
export function buatKasus(m: MasukanKasusEDM): Promise<HasilBuatEDM> {
  return minta<HasilBuatEDM>(`${PREFIX_EDM}/kasus`, { metode: 'POST', badan: m })
}

/** Satu baris spreading retro - `models.SpreadingRetro`. */
export interface SpreadingRetroEDM {
  reinsurerName: string
  percentShare: string
  amount: string
  rate: string
  premiumSpreadedGross: string
  commision: string
  ovrComm: string
  premiumSpreadedNet: string
}

/** Satu baris spreading - `models.Spreading`. */
export interface SpreadingEDM {
  id: string
  treatyTypeName: string
  retrocadedShare: string
  retro: SpreadingRetroEDM[]
}

/** Rincian peserta - `models.RincianPeserta`. */
export interface RincianPesertaEDM {
  peserta: PesertaEDM
  spreading: SpreadingEDM[]
}

/** Isi popup polis lama - `models.PolisLama`. */
export interface PolisLamaEDM {
  sumber: 'aplikasi' | 'warisan' | ''
  prodKe: number
  type: string
  peserta: PesertaEDM[]
  total: number
  rekap: Record<string, string>[]
}

/** Rincian peserta - `GET /kasus/{id}/peserta/{pid}`. */
export function ambilRincian(id: string, pesertaId: string): Promise<RincianPesertaEDM> {
  return minta<RincianPesertaEDM>(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/peserta/${encodeURIComponent(pesertaId)}`)
}

/** Popup polis lama - `GET /kasus/{id}/polis-lama`. */
export function ambilPolisLama(id: string, halaman: number): Promise<PolisLamaEDM> {
  return minta<PolisLamaEDM>(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/polis-lama`, { kueri: { halaman } })
}

/** Hasil `Save` - `services.HasilSimpan`. */
export interface HasilSimpanEDM {
  ditandai: number
  status: string
  rekap: Array<Record<string, string>>
}

/** `Save` b37202 → `SetPremi_EDM`; centang `.EdmBatal` dikirim bersamanya. */
export function simpanKasus(id: string, pilihan: PilihanHapusEDM): Promise<HasilSimpanEDM> {
  return minta<HasilSimpanEDM>(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/simpan`, { metode: 'POST', badan: pilihan })
}
