// Asal: pola `modul/nbtreatyin/frontend/api.ts` (06-10-2026), disesuaikan rute EDM Treaty In.
//
// Panggilan backend modul EDM Treaty In, satu fungsi per rute (`backend/handlers/rute.go`, prefix
// `/api/edm-treaty-in`). Klien HTTP-nya `inti/frontend/klien.ts`.
//
// ⛔ Nol perhitungan uang di sini. Setiap rumus Pega dijalankan backend (`POST .../hitung`) - layar hanya mengirim
// halaman dan menampilkan hasilnya.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_EDMTREATYIN = '/api/edm-treaty-in'

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
  prodKe: number
  oldPolisId: string
  generasiTertutup: boolean
  createOp: string
  tglCreate: string
}

/** Satu baris daftar portal - `models.RingkasanKasus` (grid `SFAPortal_Endorsement_Treaty`, RD `InboxEDM_RD2`). */
export interface RingkasanKasus {
  /** `A.pyID` - kolom 1 "EDM Number" (tautan buka kasus). */
  id: string
  noOffer: string
  noPolis: string
  edmNo: string
  sobName: string
  cedingCoName: string
  edmType: string
  proportionalType: string
  marketingName: string
  nbStatus: string
  statusWork: string
  /** Syarat tautan buka kasus (`.pxPages(A).PositionNote != 'ReasTreatyInAdmin'` = nonaktif). */
  positionNote: string
  tglCreate: string
  startDate: string
  tglUpdate: string
  /** `TGL_PROD` (ProductionDate Utility1) - kolom "Production Date" tab Resolved; kosong bila belum selesai. */
  tglProd?: string
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
  jenisReas: Pilihan[] | null
  /** "Source of Change" layar `TreatyCreateEdm` (DT `TreatyEDMListType`, tidak ada di korpus). */
  jenisEdm: Pilihan[] | null
}

/** Jawaban sel "No Polis Treaty" layar Create - `services.HasilPeriksaPolis`. */
export interface HasilPeriksaPolis {
  /** "No Master Treaty" (`DisplayData.CARI2`). */
  noMaster: string
  /** `TrtERR.CARI1` - kosong = tanpa galat. */
  galat: string
}

/** Isian layar `TreatyCreateEdm` - `services.PermintaanBuat`. */
export interface PermintaanBuat {
  noPolis: string
  edmType: string
}

/** Action set satu sel / tombol - `services.PermintaanHitung`, SATU bentuk: `urutan` satu refresh atau lebih,
 *  dijalankan berurutan atas halaman yang sama. `indeks` berbasis 1 (`.pxListSubscript`). */
export interface PermintaanHitung {
  urutan: { aksi: string; param?: string }[]
  indeks?: number
  halaman: Halaman
}

/** Akibat satu submit - `services.HasilKirim`. */
export interface HasilKirim {
  kasus: Kasus
  pesanKonversi?: string
}

/** Satu baris grid popup `BusinessAndSOBListEDM` (RDB `TreatyLoadMasterJoinEdmChooseBusiness` / RD
 *  `BrowseTREATY_IN_EDM`). */
export type BarisBisnis = Record<string, string>

/** Satu baris kotak masuk Beranda - `models.AntreanKotakMasuk`. */
export interface AntreanKotakMasuk {
  workbasket: string
  nama: string
  jumlah: number
}

const kasus = (id: string) => `${PREFIX_EDMTREATYIN}/kasus/${encodeURIComponent(id)}`

/** Kotak masuk Beranda: berkas yang menunggu akun per workbasket yang ia pegang (pola NB Treaty In). */
export function kotakMasuk(): Promise<AntreanKotakMasuk[]> {
  return minta<AntreanKotakMasuk[]>(`${PREFIX_EDMTREATYIN}/kotak-masuk`)
}

/** Daftar berkas yang menunggu akun di `workbasket` (null = semua yang ia pegang) - kotak masuk Beranda. */
export function daftarMenunggu(workbasket: string | null): Promise<RingkasanKasus[]> {
  return minta<RingkasanKasus[]>(`${PREFIX_EDMTREATYIN}/kotak-masuk/kasus`, {
    kueri: { workbasket: workbasket ?? undefined },
  })
}

/** Daftar portal (`InboxEDM_RD2`): `cari` = kotak saring "Policy Number". */
export function daftarKasus(cari: string, selesai = false): Promise<RingkasanKasus[]> {
  // status=selesai - switch Resolved (aturan portal NB, keputusan work owner 07-10-2026)
  return minta<RingkasanKasus[]>(`${PREFIX_EDMTREATYIN}/kasus`, {
    kueri: { cari: cari || undefined, status: selesai ? 'selesai' : undefined },
  })
}

/** Sel "No Polis Treaty" berubah: `TrtEdmCheckPolicyError` lalu `CheckNopolisAvailability`. */
export function periksaPolis(nopolis: string): Promise<HasilPeriksaPolis> {
  return minta<HasilPeriksaPolis>(`${PREFIX_EDMTREATYIN}/periksa-polis`, { kueri: { nopolis } })
}

/** Tombol "Create" layar `TreatyCreateEdm` (`CreateEDMT`). */
export function buatKasus(p: PermintaanBuat): Promise<Kasus> {
  return minta<Kasus>(`${PREFIX_EDMTREATYIN}/kasus`, { metode: 'POST', badan: p })
}

export function bukaKasus(id: string): Promise<Layar> {
  return minta<Layar>(kasus(id))
}

/** Tombol Save (admin). */
export function simpanKasus(id: string, halaman: Halaman): Promise<Layar> {
  return minta<Layar>(kasus(id), { metode: 'PUT', badan: { halaman } })
}

export function hitung(id: string, p: PermintaanHitung): Promise<Layar> {
  return minta<Layar>(`${kasus(id)}/hitung`, { metode: 'POST', badan: p })
}

/** Isi grid popup `BusinessAndSOBListEDM` - isian layar dikirim (`pySubmitData`), tanpa simpan. */
export function daftarBisnis(id: string, halaman: Halaman): Promise<BarisBisnis[]> {
  return minta<BarisBisnis[]>(`${kasus(id)}/bisnis`, { metode: 'POST', badan: { halaman } })
}

/** Tombol Choose popup `BusinessAndSOBListEDM` (`EDMChooseBusiness_Act`). XML tidak membaca baris terpilih -
 *  badan hanya halaman. */
export function pilihBisnis(id: string, halaman: Halaman): Promise<Layar> {
  return minta<Layar>(`${kasus(id)}/pilih-bisnis`, { metode: 'POST', badan: { halaman } })
}

/** finishAssignment. */
export function kirimKasus(id: string, halaman: Halaman): Promise<HasilKirim> {
  return minta<HasilKirim>(`${kasus(id)}/kirim`, { metode: 'POST', badan: { halaman } })
}

export function ambilAcuan(): Promise<Acuan> {
  return minta<Acuan>(`${PREFIX_EDMTREATYIN}/acuan`)
}

// ------------------------------------------------------------------ Copy Old (perintah work owner 07-10-2026)

/** Hak layar portal akun - `services.Hak`. `copyOld` = superadmin (Kelola User) dengan menu EDM ber-hak penuh. */
export interface HakPortal {
  copyOld: boolean
}

/** Satu baris popup Copy Old - `models.DokumenLama`; `alasan` menyebut sebab, tidak pernah nilai dokumen. */
export interface DokumenLama {
  id: string
  noPolis: string
  edmNo: string
  prodKe: number
  edmType: string
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
  return minta<HakPortal>(`${PREFIX_EDMTREATYIN}/hak`)
}

/** Isi popup Copy Old: dokumen endorsemen JSON_POLIS lama yang belum ada di tabel flat, urut generasi. */
export function ambilDokumenLama(): Promise<DokumenLama[]> {
  return minta<DokumenLama[]>(`${PREFIX_EDMTREATYIN}/lama`)
}

/** `Process Copy`: salin ID yang dicentang - satu transaksi per dokumen, urut generasi, hasil per ID. */
export function salinDokumenLama(ids: string[]): Promise<JawabanSalinLama> {
  return minta<JawabanSalinLama>(`${PREFIX_EDMTREATYIN}/lama/salin`, { metode: 'POST', badan: { ids } })
}

// ------------------------------------------------------------------ halaman

/** Awalan jalur halaman polis - `models.HalamanPolis` + ".". */
export const POLIS = 'PolicyTreatyIn.'

/** Posisi admin (`pyWorkPage.PositionNote = 'ReasTreatyInAdmin'`). */
export const POSISI_ADMIN = 'ReasTreatyInAdmin'

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

/** Kunci daftar bersarang `<induk>(<n>).<anak>`, n berbasis 1 - `models.JalurAnak`. */
export function jalurAnak(induk: string, n: number, anak: string): string {
  return `${induk}(${n}).${anak}`
}
