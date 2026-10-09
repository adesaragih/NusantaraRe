// Klien Claim Fac In - `/api/claim-fac-in` (`modul/claimfacin/backend/handlers/rute.go`). Layar datang dari server
// sebagai pohon tata (`models.Tata`): tampil / hanya-baca / nonaktif / wajib sudah dievaluasi di server; layar hanya
// merender dan mengirim balik nilai medan terbuka bersama setiap aksi. Disalin dari pola
// `modul/claimnonprop/frontend/api.ts` (bukan impor). Beda dengan Claim Non Prop: panel baris grid masterDetail umum
// (`Layar.panel`, kunci `models.KunciPanel`) dan alamat aksi berkonteks (`PermintaanAksi.konteks`).

import {
  BATAS_WAKTU_MS,
  kegagalanDari,
  minta,
  mintaFormulir,
  rakitURL,
  unduhBerkasBeridentitas,
} from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

export const PREFIX_CFI = '/api/claim-fac-in'

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
  /** Kosong (tidak dikirim) = tombol penutup modal / panel (Cancel), atau tombol OQ yang selalu nonaktif. */
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
  /**
   * Grid masterDetail (expandPane): prefiks panel baris. Baris ke-n membuka `Layar.panel[rincian + ':' + jalur(n)]`
   * (`models.KunciPanel`); panel tidak dikirim = baris itu tidak dapat dibuka.
   */
  rincian?: string
}

export interface Kasus {
  id: string
  tahap: string
  posisi: string
  statusWork: string
  pembuatId: string
  pembuatNama: string
  tglCreate: string
  tglUpdate?: string
  sumber?: string
}

/** Layar satu kasus - `services.Layar`. */
export interface Layar {
  kasus: Kasus
  label: string
  halaman: Halaman
  bolehKerja: boolean
  tata: Tata[]
  /**
   * Panel baris grid masterDetail, kunci `models.KunciPanel` - mis. `est:ClaimData.ObjectList(1)`,
   * `estitem:ClaimData.ObjectList(1).ObjectItemList(2)`, `adjdtl:...ObjectItemList(i).Adjustment(a)`. Panel bersarang:
   * grid di dalam panel boleh ber-`rincian` sendiri.
   */
  panel?: Record<string, Tata[]>
  /** Local action / harness: `pilihPolis`, `protectDOL`, `pla:<o>`, `dla:<o>`, `cedant:<adj>`, `komite:<adj>`, `tutup`, `tolak`. */
  modal?: Record<string, Tata[]>
  pesan?: string[]
  pesanMedan?: Record<string, string[]>
  info?: string
  /** Modal yang dibuka layar sesudah aksi. */
  bukaModal?: string
  /** Penanda mode layar (`models.ModeLayar`) - dikembalikan di setiap aksi. */
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
  policyNo: string
  insuredName: string
  dateOfLoss: string
}

export interface Pilihan {
  nilai: string
  label: string
  tambahan?: Record<string, string>
}

/** `services.AcuanStatis`. */
export interface AcuanStatis {
  mataUang: Pilihan[]
  /** BrowseReinsuranceType_RD (Treaty Type grid spreading). */
  jenisReas: Pilihan[]
  kode: Record<string, string[]>
  /** Label tampilan kode; kode tanpa label tampil apa adanya. */
  labelKode?: Record<string, Record<string, string>>
}

/** `services.PermintaanAksi`. */
export interface PermintaanAksi {
  aksi: string
  /** Kunci panel baris atau modal tempat unsur berada; '' = layar utama. */
  konteks?: string
  /** Baris grid (1..n) DI DALAM konteksnya tempat sel / tombol berada; 0 = bukan sel grid (juga tombol Add kepala grid). */
  indeks: number
  param?: string
  tahap?: string
  masukan?: Record<string, string>
  mode?: Record<string, string>
}

export type JenisDaftar = 'saya' | 'workbasket' | 'selesai'

export function daftarKasus(daftar: JenisDaftar, cari: string): Promise<RingkasanKasus[]> {
  return minta(`${PREFIX_CFI}/kasus`, { kueri: { daftar, cari: cari.trim() === '' ? undefined : cari.trim() } })
}

export function buatKasus(): Promise<Layar> {
  return minta(`${PREFIX_CFI}/kasus`, { metode: 'POST' })
}

/** `lihat` = tampilan saja walau pemegang. */
export function bukaKasus(id: string, lihat = false): Promise<Layar> {
  return minta(`${PREFIX_CFI}/kasus/${encodeURIComponent(id)}${lihat ? '?lihat=1' : ''}`)
}

/**
 * Validasi layar gagal (422 `{galat, pesan, layar?}`, `services.GalatValidasi`): aksi dibatalkan server; `layar` =
 * layar dengan isian dan pesannya (tidak disimpan).
 */
export class GalatValidasiAksi extends Error {
  constructor(
    readonly pesan: string[],
    readonly layar?: Layar,
  ) {
    super(pesan.join('; '))
    this.name = 'GalatValidasiAksi'
  }
}

/**
 * Aksi layar. Dibaca langsung (bukan `minta`) karena jawaban 422 membawa `pesan` dan `layar` yang dibuang amplop galat
 * inti; jawaban gagal lainnya lewat `kegagalanDari` seperti `minta`.
 */
export async function aksiKasus(id: string, r: PermintaanAksi): Promise<Layar> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  let jawab: Response
  try {
    jawab = await fetch(rakitURL(`${PREFIX_CFI}/kasus/${encodeURIComponent(id)}/aksi`), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...headerIdentitas() },
      body: JSON.stringify(r),
      signal: kendali.signal,
    })
  } finally {
    clearTimeout(jam)
  }
  const teks = await jawab.text()
  if (!jawab.ok) throw galatAksi(jawab.status, teks)
  try {
    return JSON.parse(teks) as Layar
  } catch {
    throw kegagalanDari(jawab.status, teks)
  }
}

/** Jawaban gagal aksi: 422 ber-`pesan` = GalatValidasiAksi; selainnya amplop galat inti. */
export function galatAksi(status: number, teks: string): Error {
  if (status === 422) {
    try {
      const o = JSON.parse(teks) as { pesan?: unknown; layar?: Layar }
      if (Array.isArray(o.pesan)) {
        return new GalatValidasiAksi(
          o.pesan.filter((p): p is string => typeof p === 'string'),
          o.layar ?? undefined,
        )
      }
    } catch {
      // bukan JSON: amplop galat inti yang menamainya
    }
  }
  return kegagalanDari(status, teks)
}

/** Kueri pilihan pop-up / autocomplete (`services.PermintaanPilihan`). */
export interface KueriPilihan {
  /** Panel baris tempat sel berada (sumber tingkat item / adjustment); '' = layar utama. */
  konteks?: string
  /** Baris grid di konteks (item untuk sumber tingkat item; objek untuk Outstanding). */
  indeks?: number
  cari?: string
  /** Search Type pop-up Choose Polis. */
  jenisCari?: string
}

export function pilihanKasus<T>(id: string, jenis: string, q: KueriPilihan = {}): Promise<T> {
  const isi = (v: string | undefined) => (v === undefined || v.trim() === '' ? undefined : v.trim())
  return minta(`${PREFIX_CFI}/kasus/${encodeURIComponent(id)}/pilihan/${encodeURIComponent(jenis)}`, {
    kueri: {
      konteks: isi(q.konteks),
      indeks: q.indeks !== undefined && q.indeks > 0 ? String(q.indeks) : undefined,
      cari: isi(q.cari),
      jenisCari: isi(q.jenisCari),
    },
  })
}

export function ambilAcuan(): Promise<AcuanStatis> {
  return minta(`${PREFIX_CFI}/acuan`)
}

/** Hak halaman awal (`services.HakPelaku`): switch Teknik aktif hanya bagi anggota workbasket Choose Surveyor. */
export interface HakPelaku {
  workbasketSurveyor: boolean
}

export function ambilHak(): Promise<HakPelaku> {
  return minta(`${PREFIX_CFI}/hak`)
}

/** Satu baris grid pop-up Choose Polis - `models.BarisPolisCari` (GetPolisForClaim_SQL). */
export interface BarisPolisCari {
  policyNo: string
  customerName: string
  sourceOfBusinessName: string
  cedingCoName: string
  qq: string
  startDateTime: string
  endDateTime: string
  prodke: string
  businessName: string
}

/** Satu baris grid "List Cause of Loss" - `repository.BarisSebab`. */
export interface BarisSebab {
  id: string
  description: string
}

/** Satu baris grid katastrofe - `repository.BarisKatastrofe`. */
export interface BarisKatastrofe {
  id: string
  stsKatastrofe: string
  nonKatastrofeType: string
  note: string
}

/** Satu baris grid Outstanding_SC - `models.BarisRingkasan`. */
export interface BarisRingkasan {
  noClaim: string
  noPolis: string
  acceptedNo: string
  nama: string
  coverage: string
  currency: string
  tsiRnm: string
  nilai: string
  kurs: string
  konversi: string
}

/** Isi pop-up Outstanding - `models.Ringkasan` (GetAllData_Act). */
export interface RingkasanOutstanding {
  labelNama: string
  baris: BarisRingkasan[] | null
  total: string
  mataUang: string
}

/** `AttachCategory.pxResults` - kategori master FAC dan cacah berkas klaim ini. */
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
  return minta(`${PREFIX_CFI}/kasus/${encodeURIComponent(id)}/lampiran`)
}

/** `GCNMSaveAttachments`: `kategori` bersama (Upload File baris) ATAU `kategoriBerkas` per berkas (Add attachment). */
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
  return mintaFormulir(`${PREFIX_CFI}/kasus/${encodeURIComponent(id)}/lampiran`, isi)
}

/** View File: isi satu dokumen klaim - fetch beridentitas. */
export function unduhLampiran(id: string, a: Lampiran): Promise<void> {
  return unduhBerkasBeridentitas(`${jalurLampiran(id, a)}/isi`, a.namaFile)
}

function jalurLampiran(id: string, a: Lampiran): string {
  return `${PREFIX_CFI}/kasus/${encodeURIComponent(id)}/lampiran/${encodeURIComponent(a.id)}`
}

/** Isi satu dokumen klaim sebagai Blob - `View` pdf / gambar di popup penampil. */
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

/** View Office Online - URL bertanda tangan untuk penampil kantor. */
export function tautanOfficeLampiran(id: string, a: Lampiran): Promise<{ url: string }> {
  return minta(`${jalurLampiran(id, a)}/office`)
}

/** Change Category - dokumen terpilih ke kategori lain. */
export function pindahKategoriLampiran(id: string, kategori: string, ids: string[]): Promise<{ ok: boolean }> {
  return minta(`${PREFIX_CFI}/kasus/${encodeURIComponent(id)}/lampiran/kategori`, {
    metode: 'POST',
    badan: { kategori, ids },
  })
}

/** Delete. */
export function hapusLampiran(id: string, a: Lampiran): Promise<{ ok: boolean }> {
  return minta(`${jalurLampiran(id, a)}/hapus`, { metode: 'POST' })
}
