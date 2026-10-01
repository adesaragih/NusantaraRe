// modul/endorsementlife/frontend/api.ts - panggilan backend modul Endorsement Life, satu fungsi per rute
// (`backend/handlers/rute_edm.go`). Klien HTTP-nya `inti/klien.ts`.
//
// ⛔ Uang dan angka TEKS sepanjang jalan - tidak pernah `Number` (ADR-0003). Kosong = `""`.

import { kegagalanDari, minta, rakitURL } from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'
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
  /** Riwayat keputusan `T_VIEW_SUGGEST`, `NO DESC` - grid `ConfirmSection` b1856. */
  riwayat: RiwayatEDM[]
  sumber: 'aplikasi' | 'warisan' | ''
}

/** Satu baris riwayat - `models.BarisRiwayat`. */
export interface RiwayatEDM {
  no: number
  date: string
  pic: string
  status: string
  comment: string
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

/** Satu penolakan CSV - `models.PesanCSV` (AC 37: nomor baris dan kolom). */
export interface PesanCSVEDM {
  baris: number
  kolom: string
  pesan: string
}

/** Tinjauan unggahan - `services.HasilPeriksaCSV`. */
export interface PeriksaCSVEDM {
  total: number
  ditolak: number
  pesan: PesanCSVEDM[]
  terpotong: boolean
  diabaikan: string[]
}

/** Hasil `Add CSV Data` - `services.HasilTambahCSV`. */
export interface TambahCSVEDM {
  disimpan: number
  dibuang: number
  rekap: Array<Record<string, string>>
}

/** Jawaban `Add CSV Data`: tersimpan, atau ditolak beserta seluruh penolakannya (422). */
export type JawabanTambahCSVEDM =
  | { jenis: 'tersimpan'; hasil: TambahCSVEDM }
  | { jenis: 'ditolak'; galat: string; periksa: PeriksaCSVEDM }

/**
 * Tenggat unggahan - 10 menit. ⛔ `mintaFormulir` inti memutus pada 30 detik, dan pemutusan
 * membatalkan transaksi `Add CSV Data` di server: unggahan tanpa batas baris (AC 36) butuh tenggat
 * sendiri.
 */
export const BATAS_WAKTU_UNGGAH_MS = 600_000

async function kirimBerkas(jalur: string, berkas: File): Promise<{ status: number; teks: string }> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_UNGGAH_MS)
  try {
    const jawab = await fetch(rakitURL(jalur), { method: 'POST', headers: { ...headerIdentitas() }, body: isi, signal: kendali.signal })
    return { status: jawab.status, teks: await jawab.text() }
  } finally {
    clearTimeout(jam)
  }
}

/** `Upload CSV` b8973 - tinjau: urai dan validasi, nol tulis. */
export async function periksaCSV(id: string, berkas: File): Promise<PeriksaCSVEDM> {
  const j = await kirimBerkas(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/unggah`, berkas)
  if (j.status !== 200) throw kegagalanDari(j.status, j.teks)
  return JSON.parse(j.teks) as PeriksaCSVEDM
}

/** `Add CSV Data` b10405 → `SaveCSVEDMLife`; 422 berdaftar pesan dikembalikan, bukan dilempar. */
export async function tambahCSV(id: string, berkas: File): Promise<JawabanTambahCSVEDM> {
  const j = await kirimBerkas(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/csv`, berkas)
  if (j.status === 200) return { jenis: 'tersimpan', hasil: JSON.parse(j.teks) as TambahCSVEDM }
  if (j.status === 422) {
    const b = JSON.parse(j.teks) as Partial<PeriksaCSVEDM> & { galat?: string }
    if (Array.isArray(b.pesan)) {
      return {
        jenis: 'ditolak',
        galat: b.galat ?? '',
        periksa: { total: b.total ?? 0, ditolak: b.ditolak ?? 0, pesan: b.pesan, terpotong: b.terpotong === true, diabaikan: b.diabaikan ?? [] },
      }
    }
  }
  throw kegagalanDari(j.status, j.teks)
}

/** Keputusan - `services.MasukanPutusan` (radio `EmailTypePL` + `Comment`). */
export interface MasukanPutusanEDM {
  status: string
  comment: string
}

/** Hasil keputusan - `services.HasilPutusan`. */
export interface HasilPutusanEDM {
  status: string
  noEndorsement: string
  peserta: number
  rekapWarisan: number
}

/** `Submit` b37494 / b38109 → `IsLifeAccepted`: Confirm meresmikan versi, Decline menutup kasus. */
export function putuskanKasus(id: string, m: MasukanPutusanEDM): Promise<HasilPutusanEDM> {
  return minta<HasilPutusanEDM>(`${PREFIX_EDM}/kasus/${encodeURIComponent(id)}/putuskan`, { metode: 'POST', badan: m })
}
