// modul/nbtreatyin/frontend/api.ts - panggilan backend modul NB Treaty In, satu
// fungsi per rute (`backend/handlers/rute.go`). Klien HTTP-nya
// `inti/frontend/klien.ts`.
//
// ⛔ Nol perhitungan uang di sini. Setiap rumus Pega dijalankan backend
// (`POST .../hitung`, models/hitung.go) - layar hanya mengirim halaman dan
// menampilkan hasilnya (AC 25, 79).

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_NBTREATYIN = '/api/nb-treaty-in'

/** Satu anggota PageList - nilai teks per nama properti. */
export type Baris = Record<string, string>

/** Halaman kerja - SAMA dengan `models.Halaman`. Jalur relatif `pyWorkPage`. */
export interface Halaman {
  nilai: Record<string, string>
  daftar: Record<string, Baris[] | null>
  pesan?: Record<string, string[]>
}

/** Keadaan kerja satu kasus - `models.Kasus`. */
export interface Kasus {
  id: string
  position: string
  statusWork: string
  positionNote: string
  noPolis: string
  generasiTertutup: boolean
  createOp: string
  tglCreate: string
}

/** Satu baris daftar portal - `models.RingkasanKasus`. */
export interface RingkasanKasus {
  id: string
  /** T_GENERAL_POLIS_TREATY.TREATY_GROUP_NAME - kolom "Treaty Business" (perintah work owner 10-10-2026). */
  treatyGroupName: string
  insuredName: string
  marketingName: string
  nbStatus: string
  statusWork: string
  positionNote: string
  noPolis: string
  tglCreate: string
  /** T_WORK_POLIS.CREATE_OP_NAME - kolom "User Create" portal. */
  createOpName: string
  /** T_POLIS_QUOTATION.PROPORTIONAL_TYPE - kolom "Type" portal. */
  proportionalType: string
  /** Kolom daftar kotak masuk Beranda: T_GENERAL_POLIS_TREATY.CEDING_CO_NAME / START_DATE, T_WORK_POLIS.TGL_UPDATE. */
  cedingCoName: string
  startDate: string
  tglUpdate: string
  /** T_GENERAL_POLIS_TREATY.TGL_PROD - kolom "Production Date" tab Resolved portal (keputusan work owner 07-10-2026). */
  productionDate: string
}

/** Tombol submit yang tampil - `models.TombolKirim`. */
export type TombolKirim = '' | 'kirim' | 'konfirmasi-tolak' | 'nomor-polis'

/** Satu kasus siap tampil - `services.Layar`. */
export interface Layar {
  kasus: Kasus
  halaman: Halaman
  bolehKerja: boolean
  tombol: TombolKirim
  medanWajib: string[] | null
  tempat: Record<string, boolean>
  pesan?: string[]
}

export interface Pilihan {
  nilai: string
  label: string
}

/** Daftar pilihan layar - `services.Acuan`. */
export interface Acuan {
  mataUang: Pilihan[] | null
  mo: Pilihan[] | null
  spreading: Pilihan[] | null
  jenisReas: Pilihan[] | null
}

/** Satu baris view kontrak (nama kolom view) - `models.BarisKontrak`. */
export type BarisKontrak = Record<string, string>

/** Satu baris RD `BrowseAgentHierarkiList_RD` (pemilih SOB) - `models.BarisAgen`. */
export interface BarisAgen {
  id: string
  clientName: string
  leader0: string
  childCount: string
  clientId: string
}

export interface NomorPolis {
  id: string
  policyNo: string
}

/** Action set satu sel - `services.PermintaanHitung`, SATU bentuk: `urutan`
 *  berisi satu refresh atau lebih, dijalankan berurutan atas halaman yang sama. */
export interface PermintaanHitung {
  urutan: { aksi: string; param?: string }[]
  indeks?: number
  halaman: Halaman
}

const kasus = (id: string) => `${PREFIX_NBTREATYIN}/kasus/${encodeURIComponent(id)}`

/** Satu baris kotak masuk Beranda - `models.AntreanKotakMasuk`. */
export interface AntreanKotakMasuk {
  workbasket: string
  nama: string
  jumlah: number
}

/** Kotak masuk Beranda: berkas yang menunggu akun per workbasket yang ia pegang (keputusan work owner 06-10-2026). */
export function kotakMasuk(): Promise<AntreanKotakMasuk[]> {
  return minta<AntreanKotakMasuk[]>(`${PREFIX_NBTREATYIN}/kotak-masuk`)
}

/** Daftar berkas yang menunggu akun di `workbasket` (null = semua yang ia pegang) - kotak masuk Beranda. */
export function daftarMenunggu(workbasket: string | null): Promise<RingkasanKasus[]> {
  return minta<RingkasanKasus[]>(`${PREFIX_NBTREATYIN}/kotak-masuk/kasus`, {
    kueri: { workbasket: workbasket ?? undefined },
  })
}

/** Daftar portal: hanya berkas buatan akun ini (keputusan work owner 06-10-2026); `selesai` = switch Resolved. */
export function daftarKasus(cari: string, posisi = '', selesai = false): Promise<RingkasanKasus[]> {
  return minta<RingkasanKasus[]>(`${PREFIX_NBTREATYIN}/kasus`, {
    kueri: { cari: cari || undefined, posisi: posisi || undefined, status: selesai ? 'selesai' : undefined },
  })
}

export function buatKasus(): Promise<Kasus> {
  return minta<Kasus>(`${PREFIX_NBTREATYIN}/kasus`, { metode: 'POST' })
}

export function bukaKasus(id: string): Promise<Layar> {
  return minta<Layar>(kasus(id))
}

export function simpanKasus(id: string, halaman: Halaman): Promise<Layar> {
  return minta<Layar>(kasus(id), { metode: 'PUT', badan: { halaman } })
}

export function hitung(id: string, p: PermintaanHitung): Promise<Layar> {
  return minta<Layar>(`${kasus(id)}/hitung`, { metode: 'POST', badan: p })
}

export function pilihBisnis(id: string, idDetail: string, halaman: Halaman): Promise<Layar> {
  return minta<Layar>(`${kasus(id)}/pilih-bisnis`, { metode: 'POST', badan: { idDetail, halaman } })
}

/** Jawaban klik satu baris popup `SOB` - `models.SumberBisnisPostDT`: keempat medan
 *  `Quotation.*` yang ditulis `SearchHierarkiSourceBizAgent_PostDT` (langkah 1.1, 1.2, 4, 5). */
export interface HasilSumberBisnis {
  sourceOfBusiness: string
  sobName: string
  sobLeader0: string
  sobLeader1: string
}

/** Klik satu baris popup `SOB` - pra-proses `SearchHierarkiSourceBizAgent_PostDT`,
 *  TANPA simpan (F4): hasilnya dipegang layar (`pegangSumberBisnis`) dan ikut
 *  terkirim pada Save/Submit/refresh; server menerimanya hanya bila cocok dengan
 *  RD `BrowseAgentHierarkiList_RD` yang dijalankan ulang. */
export function pilihSumberBisnis(id: string, idAgen: string, halaman: Halaman): Promise<HasilSumberBisnis> {
  return minta<HasilSumberBisnis>(`${kasus(id)}/pilih-sumber-bisnis`, { metode: 'POST', badan: { idAgen, halaman } })
}

export function terbitkanNomor(id: string, halaman: Halaman): Promise<NomorPolis> {
  return minta<NomorPolis>(`${kasus(id)}/nomor-polis`, { metode: 'POST', badan: { halaman } })
}

/** Akibat satu submit - `services.HasilKirim`. `pesanKonversi` = `FlagErrorKonversi`
 *  bila konversi Arasapas sesudah selesai gagal (penyimpanan tetap berhasil). */
export interface HasilKirim {
  kasus: Kasus
  pesanKonversi?: string
}

export function kirimKasus(id: string, halaman: Halaman): Promise<HasilKirim> {
  return minta<HasilKirim>(`${kasus(id)}/kirim`, { metode: 'POST', badan: { halaman } })
}

/** Isi grid popup `BusinessAndSOBList` (RD `BrowseTreatyJoinEDM`): showHarness `pySubmitData=Yes`
 *  mengirim isian layar - server menyaring dengan `QuotationData.ProportionalType`-nya, tanpa simpan.
 *  `saringan` (nama kolom view -> teks) dicari di SERVER sebelum batas 500 baris (keputusan work owner
 *  06-10-2026), supaya kontrak di luar 500 baris pertama dapat ditemukan. */
export function daftarBisnis(
  id: string,
  halaman: Halaman,
  saringan: Record<string, string> = {},
): Promise<BarisKontrak[]> {
  return minta<BarisKontrak[]>(`${kasus(id)}/bisnis`, { metode: 'POST', badan: { halaman, saringan } })
}

/** Isi TreeGrid popup `SOB` (`Section/SourceHierarki`). */
export function daftarSumberBisnis(): Promise<BarisAgen[]> {
  return minta<BarisAgen[]>(`${PREFIX_NBTREATYIN}/sumber-bisnis`)
}

export function ambilAcuan(): Promise<Acuan> {
  return minta<Acuan>(`${PREFIX_NBTREATYIN}/acuan`)
}

// ------------------------------------------------------------------ Copy Old (perintah work owner 07-10-2026)

/** Hak layar portal akun - `services.Hak`. `copyOld` = superadmin (Kelola User) dengan menu NB ber-hak penuh. */
export interface HakPortal {
  copyOld: boolean
}

/** Satu baris popup Copy Old - `models.DokumenLama`; `alasan` menyebut sebab, tidak pernah nilai dokumen. */
export interface DokumenLama {
  id: string
  noOffer: string
  noPolis: string
  insuredName: string
  businessName: string
  sobName: string
  cedingCoName: string
  tglProd: string
  bolehDisalin: boolean
  alasan: string[]
}

export type StatusSalinLama = 'disalin' | 'sudahAda' | 'ditolak' | 'gagal'

/** Hasil satu ID sesudah `Process Copy` - `models.HasilSalinLama`. */
export interface HasilSalinLama {
  id: string
  status: StatusSalinLama
  pesan: string[]
}

export interface JawabanSalinLama {
  hasil: HasilSalinLama[]
  disalin: number
}

export function ambilHak(): Promise<HakPortal> {
  return minta<HakPortal>(`${PREFIX_NBTREATYIN}/hak`)
}

/** Isi popup Copy Old: dokumen polis NB lama (JSON_POLIS generasi 0) yang belum ada di tabel flat. */
export function ambilDokumenLama(): Promise<DokumenLama[]> {
  return minta<DokumenLama[]>(`${PREFIX_NBTREATYIN}/lama`)
}

/** `Process Copy`: salin ID yang dicentang - satu transaksi per dokumen, hasil per ID. */
export function salinDokumenLama(ids: string[]): Promise<JawabanSalinLama> {
  return minta<JawabanSalinLama>(`${PREFIX_NBTREATYIN}/lama/salin`, { metode: 'POST', badan: { ids } })
}

// ------------------------------------------------------------------ halaman

/** Awalan jalur halaman polis - SATU-SATUNYA salinan `models.HalamanPolis` + ".". */
export const POLIS = 'PolicyTreatyIn.'

/** Awalan jalur halaman master kontrak - `models.HalamanMaster` + ".". */
export const MASTER = 'TreatyIn.'

/** Nilai `.ClaimType` XOL Retro - `models.KlaimXOLRetro`. */
export const KLAIM_XOL_RETRO = 'XOL Retro'

/** Nilai satu jalur halaman ("" bila tidak ada). */
export function nilai(h: Halaman, jalur: string): string {
  return h.nilai?.[jalur] ?? ''
}

/** Salinan halaman dengan satu nilai diganti. */
export function setel(h: Halaman, jalur: string, v: string): Halaman {
  return { ...h, nilai: { ...h.nilai, [jalur]: v } }
}

/** Baris satu daftar (kosong bila tidak ada). */
export function daftar(h: Halaman, jalur: string): Baris[] {
  return h.daftar?.[jalur] ?? []
}

/** Salinan halaman dengan satu daftar diganti. */
export function setelDaftar(h: Halaman, jalur: string, b: Baris[]): Halaman {
  return { ...h, daftar: { ...h.daftar, [jalur]: b } }
}
