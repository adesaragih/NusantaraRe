// Klien Claim Prop - `/api/claim-prop` (`modul/claimprop/backend/handlers/rute.go`). Layar datang dari server sebagai
// pohon tata (`models.Tata`): tampil / hanya-baca / nonaktif / wajib sudah dievaluasi di server; layar hanya merender
// dan mengirim balik nilai medan terbuka bersama setiap aksi.

import {
  BATAS_WAKTU_MS,
  kegagalanDari,
  minta,
  mintaFormulir,
  rakitURL,
  unduhBerkasBeridentitas,
} from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

export const PREFIX_CP = '/api/claim-prop'

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
  /** Jalur teks yang ditampilkan medan ber-sumber; nilai `jalur` tetap yang disimpan (Consultant / Adjuster: nama). */
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
  adjustment?: Record<string, Tata[]>
  modal?: Record<string, Tata[]>
  pesan?: string[]
  pesanMedan?: Record<string, string[]>
  info?: string
  /** Penanda mode layar (Edit Catastrophe, Edit RNM Share) - tanpa kolom di tabel; dikembalikan di setiap aksi. */
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
  jenisReas: Pilihan[]
  jenisReas4: Pilihan[]
  kode: Record<string, string[]>
  /** Label tampilan kode (Report Type, Reporter Status - dari work owner 08-10-2026); kode tanpa label tampil apa adanya. */
  labelKode?: Record<string, Record<string, string>>
}

export interface PermintaanAksi {
  aksi: string
  indeks?: number
  param?: string
  tahap?: string
  masukan?: Record<string, string>
  /** Penanda mode layar terakhir (`Layar.mode`) - dikembalikan di setiap aksi. */
  mode?: Record<string, string>
}

export type JenisDaftar = 'saya' | 'workbasket' | 'selesai'

export function daftarKasus(daftar: JenisDaftar, cari: string): Promise<RingkasanKasus[]> {
  return minta(`${PREFIX_CP}/kasus`, { kueri: { daftar, cari: cari.trim() === '' ? undefined : cari.trim() } })
}

export function buatKasus(): Promise<Layar> {
  return minta(`${PREFIX_CP}/kasus`, { metode: 'POST' })
}

/** `lihat` = tampilan saja walau pemegang (View more details Komite Claim Prop, keputusan work owner 09-10-2026). */
export function bukaKasus(id: string, lihat = false): Promise<Layar> {
  return minta(`${PREFIX_CP}/kasus/${encodeURIComponent(id)}${lihat ? '?lihat=1' : ''}`)
}

export function aksiKasus(id: string, r: PermintaanAksi): Promise<Layar> {
  return minta(`${PREFIX_CP}/kasus/${encodeURIComponent(id)}/aksi`, { metode: 'POST', badan: r })
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
  return minta(`${PREFIX_CP}/kasus/${encodeURIComponent(id)}/pilihan/${encodeURIComponent(jenis)}`, { kueri })
}

export function ambilAcuan(): Promise<AcuanStatis> {
  return minta(`${PREFIX_CP}/acuan`)
}

/** Tombol View: berkas NB / EDM Treaty In untuk nomor polis (dibuka di tab baru). 404 = polis belum punya berkas. */
export interface BerkasPolis {
  modul: string
  kasus: string
}

export function berkasPolis(nopolis: string): Promise<BerkasPolis> {
  return minta(`${PREFIX_CP}/berkas-polis`, { kueri: { nopolis } })
}

/**
 * Tombol "+" Consultant / Adjuster (Pega MstAdjusterConsultant, keputusan work owner 08-10-2026): master baru disimpan
 * lewat API modul Adjuster Consultant - rute pinjaman `POST /api/adjuster-consultant` (`cmd/api/rakit.go`). ID dibuat
 * modul itu; namanya tidak boleh kembar (422).
 */
export interface AdjusterBaru {
  id: string
  name: string
}

export function tambahAdjuster(isi: { name: string; address: string; telpNo: string }): Promise<AdjusterBaru> {
  return minta('/api/adjuster-consultant', { metode: 'POST', badan: { id: '', ...isi } })
}

/**
 * Hak halaman awal: switch Teknik aktif hanya bagi anggota workbasket ReasKlaimTeknik; `komite` = tabel komite tampil
 * (akun memegang workbasket roster EMAILKOMITE PROP).
 */
export interface HakPelaku {
  workbasketTeknik: boolean
  komite: boolean
}

/**
 * Satu kasus komite yang menunggu workbasket / akun pelaku - `BarisKerja` modul Komite Claim Prop. Menu Komite Claim
 * Prop dibuang (keputusan work owner 09-10-2026): daftarnya dibaca lewat rute pinjaman `GET /api/komite-claim-prop/kasus`
 * (`cmd/api/rakit.go`), kasusnya dibuka di tempat (`onBukaModul`).
 */
export interface BarisKomite {
  kasusId: string
  klaimId: string
  noKlaim: string
  tingkat: number
  komiteLoop: number
  jabatan: string
  nilai: string
  mataUang: string
  tglUpdate: string
}

export function daftarKomite(): Promise<BarisKomite[]> {
  return minta('/api/komite-claim-prop/kasus')
}

/** Nama modul layar kasus komite (`onBukaModul`). */
export const MODUL_KOMITE = 'komiteclaimprop'

export function ambilHak(): Promise<HakPelaku> {
  return minta(`${PREFIX_CP}/hak`)
}

/** `AttachCategory.pxResults` - kategori master PROP dan cacah berkas klaim ini. */
export interface KategoriLampiran {
  id: string
  label: string
  countAttach: number
}

/** Satu dokumen klaim. */
export interface Lampiran {
  id: string
  namaFile: string
  kategori: string
  mime: string
  /** KATEGORI_2 (kolom Note popup NB). */
  note: string
  /** Upload Date `DD-MM-YYYY HH:mm`. */
  tanggal: string
  /** PXCREATEOPERATOR - username pengunggah. */
  operator: string
  /** Objek penyimpanan tercatat (syarat View / View Office Online). */
  adaObjek: boolean
}

export interface LampiranKasus {
  kategori: KategoriLampiran[]
  lampiran: Lampiran[]
  bolehUnggah: boolean
}

/** Kategori + dokumen klaim kasus. */
export function ambilLampiran(id: string): Promise<LampiranKasus> {
  return minta(`${PREFIX_CP}/kasus/${encodeURIComponent(id)}/lampiran`)
}

/** `GCNMSaveAttachments`: `kategori` bersama (TempInputParam.pyCategory, Upload File baris) ATAU `kategoriBerkas` per
 *  berkas (`.pyCategory`, Add attachment), satu `InsertDocument_Act` per berkas. */
export function unggahLampiran(
  id: string,
  kategori: string,
  berkas: File[],
  kategoriBerkas?: string[],
): Promise<{ lampiran: Lampiran[] }> {
  const isi = new FormData()
  if (kategori !== '') isi.append('kategori', kategori)
  berkas.forEach((b, i) => {
    isi.append('berkas', b)
    if (kategoriBerkas) isi.append('kategoriBerkas', kategoriBerkas[i] ?? '')
  })
  return mintaFormulir(`${PREFIX_CP}/kasus/${encodeURIComponent(id)}/lampiran`, isi)
}

/** View File: isi satu dokumen klaim (`GetBase64Attachment` -> `GetUrlGoogleStorage_Act`) - fetch beridentitas. */
export function unduhLampiran(id: string, a: Lampiran): Promise<void> {
  return unduhBerkasBeridentitas(`${jalurLampiran(id, a)}/isi`, a.namaFile)
}

function jalurLampiran(id: string, a: Lampiran): string {
  return `${PREFIX_CP}/kasus/${encodeURIComponent(id)}/lampiran/${encodeURIComponent(a.id)}`
}

/** Isi satu dokumen klaim sebagai Blob - `View` pdf / gambar di popup penampil (pola NB Treaty In). */
export async function ambilIsiLampiran(id: string, a: Lampiran): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`${jalurLampiran(id, a)}/isi`), {
      method: 'GET',
      headers: { ...headerIdentitas() },
      signal: kendali.signal,
    })
    if (!jawab.ok) throw kegagalanDari(jawab.status, await jawab.text())
    return await jawab.blob()
  } finally {
    clearTimeout(jam)
  }
}

/** View Office Online - URL bertanda tangan untuk penampil kantor (pola NB `DownloadDocumentPolis` ViewOffice). */
export function tautanOfficeLampiran(id: string, a: Lampiran): Promise<{ url: string }> {
  return minta(`${jalurLampiran(id, a)}/office`)
}

/** Change Category - dokumen terpilih ke kategori lain (layar Pega View File). */
export function pindahKategoriLampiran(id: string, kategori: string, ids: string[]): Promise<{ ok: boolean }> {
  return minta(`${PREFIX_CP}/kasus/${encodeURIComponent(id)}/lampiran/kategori`, {
    metode: 'POST',
    badan: { kategori, ids },
  })
}

/** Delete (pola NB `DeleteDocumentPolis_Act`). */
export function hapusLampiran(id: string, a: Lampiran): Promise<{ ok: boolean }> {
  return minta(`${jalurLampiran(id, a)}/hapus`, { metode: 'POST' })
}
