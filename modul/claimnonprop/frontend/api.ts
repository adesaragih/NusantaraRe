// Klien Claim Non Prop - `/api/claim-non-prop` (`modul/claimnonprop/backend/handlers/rute.go`). Layar datang dari server
// sebagai pohon tata (`models.Tata`): tampil / hanya-baca / nonaktif / wajib sudah dievaluasi di server; layar hanya
// merender dan mengirim balik nilai medan terbuka bersama setiap aksi. Disalin dari pola `modul/claimprop/frontend/api.ts`
// (bukan impor); nol rute lampiran (korpus Non Prop tanpa section unggahan).

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_CNP = '/api/claim-non-prop'

/** Satu baris PageList - nilai teks per properti anggota. */
export type Baris = Record<string, string>

/** Halaman kerja - `models.Halaman`. */
export interface Halaman {
  nilai: Record<string, string>
  daftar: Record<string, Baris[] | undefined>
  pesan?: Record<string, string[]>
}

/** Keadaan satu sel grid pada satu baris. */
export interface SelTata {
  tampil: boolean
  hanyaBaca?: boolean
  nonaktif?: boolean
}

/** Satu unsur tata - `models.Tata`. */
export interface Tata {
  jenis: 'bagian' | 'medan' | 'label' | 'tombol' | 'grid'
  id?: string
  label?: string
  jalur?: string
  kendali?: string
  sumber?: string
  /** Jalur teks yang ditampilkan medan ber-sumber; nilai `jalur` tetap yang disimpan. */
  tampilan?: string
  aksi?: string
  catatan?: string
  hanyaBaca?: boolean
  nonaktif?: boolean
  wajib?: boolean
  anak?: Tata[]
  kolom?: Tata[]
  baris?: SelTata[][]
  kaki?: Tata[]
  tambah?: Tata
  bernomor?: boolean
  /** Format layout Pega: dua = Inline grid double, sebaris = Inline, tab = layout group Tab, judul = kepala layar. */
  letak?: 'dua' | 'sebaris' | 'tab' | 'judul' | 'tabel'
  /** Ikon tombol dari XML (pi-plus, pi-trash, pi-pencil, pi-check). */
  ikon?: 'tambah' | 'hapus' | 'ubah' | 'simpan'
  /** Paging grid (pyGridPaginator). */
  perHalaman?: number
}

export interface Kasus {
  id: string
  tahap: string
  posisi: string
  statusWork: string
  pembuatId: string
  pembuatNama: string
  tglCreate: string
}

/** Layar satu kasus - `services.Layar`. */
export interface Layar {
  kasus: Kasus
  label: string
  halaman: Halaman
  bolehKerja: boolean
  tata: Tata[]
  /** `AdjustmentDetailNP_Section` per baris Acceptation List (kunci "1".."n"). */
  adjustment?: Record<string, Tata[]>
  /**
   * Local action / harness / expand pane: `tutupKlaim`, `cwp`, `pla`, `komite:<n>`, `interest:<i>` (expand pane
   * InputDtlInterest), `reinstatement:<n>:<i>` (expand pane ShowDetailXOL).
   */
  modal?: Record<string, Tata[]>
  pesan?: string[]
  pesanMedan?: Record<string, string[]>
  info?: string
  /** Modal yang dibuka layar sesudah aksi (harness KomiteCNP, local action CloseClaimNP). */
  bukaModal?: string
  /** Penanda mode layar (Edit Catastrophe) - dikembalikan di setiap aksi. */
  mode?: Record<string, string>
}

/** Satu baris daftar kerja - `repository.RingkasanKasus`. */
export interface RingkasanKasus {
  id: string
  tahap: string
  label: string
  statusWork: string
  pembuatNama: string
  tglCreate: string
  noClaim: string
  claimNoTemp: string
  policyNo: string
  treatyName: string
  insuredName: string
}

export interface Pilihan {
  nilai: string
  label: string
  tambahan?: Record<string, string>
}

export interface AcuanStatis {
  mataUang: Pilihan[]
  /** `ListLossAllocation.pxResults` (Treaty Name Loss Allocation). */
  lossAlloc: Pilihan[]
  kode: Record<string, string[]>
  /** Label tampilan kode (properti sekelas dengan label Claim Prop); kode tanpa label tampil apa adanya. */
  labelKode?: Record<string, Record<string, string>>
}

export interface PermintaanAksi {
  aksi: string
  /** Baris grid tingkat klaim, atau nomor akseptasi untuk aksi di panel detail akseptasi. */
  indeks?: number
  /** Baris grid DI DALAM panel detail akseptasi; 0 bila bukan sel grid panel. */
  baris?: number
  param?: string
  tahap?: string
  masukan?: Record<string, string>
  mode?: Record<string, string>
}

export type JenisDaftar = 'saya' | 'workbasket' | 'selesai'

export function daftarKasus(daftar: JenisDaftar, cari: string): Promise<RingkasanKasus[]> {
  return minta(`${PREFIX_CNP}/kasus`, { kueri: { daftar, cari: cari.trim() === '' ? undefined : cari.trim() } })
}

export function buatKasus(): Promise<Layar> {
  return minta(`${PREFIX_CNP}/kasus`, { metode: 'POST' })
}

/** `lihat` = tampilan saja walau pemegang (untuk tahap 2, pola Claim Prop). */
export function bukaKasus(id: string, lihat = false): Promise<Layar> {
  return minta(`${PREFIX_CNP}/kasus/${encodeURIComponent(id)}${lihat ? '?lihat=1' : ''}`)
}

export function aksiKasus(id: string, r: PermintaanAksi): Promise<Layar> {
  return minta(`${PREFIX_CNP}/kasus/${encodeURIComponent(id)}/aksi`, { metode: 'POST', badan: r })
}

/** `saring` = filter per kolom popup master (kunci = parameter kueri server; kosong dilewati). */
export function pilihanKasus<T>(
  id: string,
  jenis: string,
  indeks = 0,
  cari = '',
  saring: Record<string, string> = {},
): Promise<T> {
  const kueri: Record<string, string | undefined> = {
    indeks: indeks > 0 ? String(indeks) : undefined,
    cari: cari.trim() === '' ? undefined : cari.trim(),
  }
  for (const [k, v] of Object.entries(saring)) kueri[k] = v.trim() === '' ? undefined : v.trim()
  return minta(`${PREFIX_CNP}/kasus/${encodeURIComponent(id)}/pilihan/${encodeURIComponent(jenis)}`, { kueri })
}

export function ambilAcuan(): Promise<AcuanStatis> {
  return minta(`${PREFIX_CNP}/acuan`)
}

/** Tombol View polis: berkas NB / EDM Treaty In untuk nomor polis. 404 = polis belum punya berkas. */
export interface BerkasPolis {
  modul: string
  kasus: string
}

export function berkasPolis(nopolis: string): Promise<BerkasPolis> {
  return minta(`${PREFIX_CNP}/berkas-polis`, { kueri: { nopolis } })
}

/** Hak halaman awal: switch Teknik aktif hanya bagi anggota workbasket ReasKlaimTeknik. */
export interface HakPelaku {
  workbasketTeknik: boolean
}

export function ambilHak(): Promise<HakPelaku> {
  return minta(`${PREFIX_CNP}/hak`)
}
