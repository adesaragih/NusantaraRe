// ============================================================================
// services/api.ts — SATU-SATUNYA tempat frontend berbicara ke backend Go.
//
// Cara membaca berkas ini bila baru mengenal TypeScript:
//   - `interface` = "bentuk data". Ia hanya ada saat kode diperiksa (tsc),
//     tidak ikut dikirim ke browser. Gunanya: editor menegur lebih dulu bila
//     kita salah mengeja nama medan, misalnya `klaim.nomorKlaimm`.
//   - `export` = boleh dipakai berkas lain lewat `import`.
//   - `async` / `await` = menunggu jawaban server tanpa membekukan layar.
//   - `Promise<Klaim>` = "nanti, kalau sudah datang, isinya Klaim".
// ============================================================================

import { headerIdentitas } from '../store/sesi'

// ---------------------------------------------------------------------------
// 1. BENTUK DATA yang dikirim backend.
//
// Sumbernya: internal/models/klaimlife.go, fungsi MarshalJSON. Nama medan di
// sini HARUS sama persis dengan yang ditulis Go — huruf besar-kecilnya pun —
// sebab nama itulah yang muncul di JSON.
// ---------------------------------------------------------------------------

/** Nilai uang. `amount` sengaja TEKS, bukan angka — alasannya di bab 3. */
export interface Uang {
  amount: string
  currency: string
}

/** Satu baris AdjustmentList milik seorang peserta. */
export interface BarisAdjustment {
  id: string
  /** Kata untuk pengguna: "Outstanding" / "Aksep" / "Ditolak" / "Tidak diketahui". */
  status: string
  /** Kode mentah kolom STS_REJECT, dibawa untuk penelusuran. TIDAK ditampilkan. */
  kodeStatus: string
  statusDiketahui: boolean
  jumlahKlaim: Uang
  nomorAkseptasi: string
  tanggalAkseptasi: string
  komiteId: string
}

/** Satu peserta yang diklaim, beserta baris-barisnya sendiri. */
export interface Peserta {
  id: string
  nomorPremiList: string
  nomorPolis: string
  nomorSertifikat: string
  mataUang: string
  /** Penanda DIPILIH untuk diklaim (`IS_CHECK`). Teks, bukan boolean. */
  isCheck: string
  baris: BarisAdjustment[]
}

/** Satu klaim Life beserta seluruh pesertanya. */
export interface Klaim {
  id: string
  nomorKlaim: string
  nomorPolis: string
  namaBisnis: string
  kodeStatus: string
  /**
   * Status klaim yang DIHITUNG dari seluruh barisnya - "Berjalan", "Selesai",
   * "Ditolak seluruhnya", "Tidak dapat dipastikan", atau "Belum berbaris".
   *
   * ⛔ Ini BUKAN terjemahan `kodeStatus`. `kodeStatus` adalah cerminan mentah
   * baris terakhir di kolom `STS_REJECT`; yang di bawah ini kesimpulan atas
   * seluruh baris. Keduanya dapat berbeda, dan memang boleh berbeda.
   */
  statusTurunan: string
  peserta: Peserta[]
  cacahBaris: number
}

/** Jawaban GET /healthz. */
export interface Kesehatan {
  status: string
  database: string
}

// ---------------------------------------------------------------------------
// 2. KLIEN HTTP — `fetch`, bukan axios (F0.2).
//
// ⛔ RALAT RUJUKAN ADR: blok ini dahulu menyebut ADR-U-0004. ADR itu
// DIGANTIKAN ADR-0013 (lihat `docs/adr/0013…` baris `menggantikan:
// ADR-0004`), dan ADR yang sudah diganti tidak boleh dikutip sebagai
// alasan — pembaca berikutnya akan mencarinya dan menemukan keputusan
// yang sudah dicabut.
//
// Alamat backend TIDAK ditulis di kode. Ia dibaca dari `VITE_API_BASE_URL`;
// kosong berarti frontend dan backend disajikan dari alamat yang sama.
//
// # Kenapa `fetch`, bukan axios
//
// Komponen dasar (`components/ui/dasar.tsx`) dan `lib/keadaanGalat.ts` yang
// diadopsi dari REFERENSI_UI keduanya berbicara dalam `ApiFailure`.
// Mempertahankan axios berarti menulis adaptor yang MEREPRODUKSI kelas itu,
// ditambah satu lapis lagi yang dapat salah. Dua model galat berdampingan
// berarti dua jalan menampilkan kegagalan yang sama, dan yang satu akan
// diam-diam kalah.
// ---------------------------------------------------------------------------
const baseURL: string = import.meta.env.VITE_API_BASE_URL ?? ''

/** Batas waktu satu permintaan. Sama dengan axios sebelumnya. */
const BATAS_WAKTU_MS = 30_000

interface OpsiMinta {
  metode?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  /** Badan permintaan; di-JSON-kan. `undefined` berarti tanpa badan. */
  badan?: unknown
  /** Parameter kueri; nilai `undefined` dilewati. */
  kueri?: Record<string, string | number | undefined>
}

function rakitURL(jalur: string, kueri?: OpsiMinta['kueri']): string {
  if (kueri === undefined) return baseURL + jalur
  const p = new URLSearchParams()
  for (const [k, v] of Object.entries(kueri)) {
    if (v !== undefined) p.set(k, String(v))
  }
  const teks = p.toString()
  return teks === '' ? baseURL + jalur : baseURL + jalur + '?' + teks
}

/**
 * Satu permintaan HTTP.
 *
 * ⛔ Jawaban yang BUKAN JSON menjadi `BACKEND_TIDAK_TERJANGKAU`, bukan
 * `SyntaxError`. Proxy pengembangan maupun reverse-proxy produksi menjawab
 * HTML atau teks biasa ketika upstream-nya mati; tanpa pemetaan ini keadaan
 * "backend mati" muncul sebagai galat parser yang tidak dapat dibaca siapa
 * pun, dan layar menampilkannya sebagai "belum ada data".
 *
 * ⚠️ Galat JARINGAN (fetch menolak sebelum ada respons) sengaja DILEMPAR
 * apa adanya sebagai `TypeError`: `lib/keadaanGalat.ts` mengenalinya dan
 * memberi petunjuk yang benar. Membungkusnya di sini menghapus tipenya.
 */
async function minta<T>(jalur: string, opsi: OpsiMinta = {}): Promise<T> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  let jawab: Response
  try {
    jawab = await fetch(rakitURL(jalur, opsi.kueri), {
      method: opsi.metode ?? 'GET',
      headers: {
        'Content-Type': 'application/json',
        // ⛔ Identitas ikut di SETIAP permintaan, dari satu tempat.
        // Memasangnya per pemanggil berarti satu pemanggil akan lupa.
        ...headerIdentitas(),
      },
      body: opsi.badan === undefined ? undefined : JSON.stringify(opsi.badan),
      signal: kendali.signal,
    })
  } finally {
    clearTimeout(jam)
  }

  // 204 dan badan kosong: tidak ada yang perlu diurai.
  const teks = await jawab.text()
  let isi: unknown
  if (teks.trim() !== '') {
    try {
      isi = JSON.parse(teks)
    } catch {
      throw new ApiFailure(jawab.status, {
        code: 'BACKEND_TIDAK_TERJANGKAU',
        message:
          'Jawaban dari server bukan JSON; permintaan tampaknya tidak ' +
          'sampai ke backend.',
      })
    }
  }

  if (!jawab.ok) {
    const o = (isi ?? {}) as { error?: unknown }
    throw new ApiFailure(jawab.status, {
      code: 'DITOLAK_BACKEND',
      // Envelope backend kita: `{"error": "<kalimat>"}`.
      message: typeof o.error === 'string' && o.error !== '' ? o.error : undefined,
    })
  }
  return isi as T
}

// ---------------------------------------------------------------------------
// 3. UANG — aturan yang paling sering ditegur (ADR-U-0003, ADR-U-0016).
//
// Uang TIDAK PERNAH menjadi angka JavaScript. Angka JavaScript adalah
// "float64", dan float64 membulatkan diam-diam:
//     Number("1234567890.12345678")  menjadi  1234567890.1234567
// Karena itu backend mengirim uang sebagai TEKS, dan di sini teks itu hanya
// diteruskan ke layar. Jangan pernah membungkusnya dengan Number(),
// parseFloat(), atau tanda + - * /.
// ---------------------------------------------------------------------------

/**
 * Bentuk uang yang BELUM dipercaya — apa pun yang baru datang dari jaringan.
 * `unknown` artinya "belum tahu tipenya; periksa dulu sebelum dipakai".
 * Tanda `?` artinya medan itu boleh tidak ada.
 */
export interface UangMasuk {
  amount?: unknown
  currency?: unknown
}

/** Ambil jumlah uang sebagai teks, apa adanya. Kosong tetap kosong, bukan "0". */
export function jumlahUang(uang: UangMasuk | null | undefined): string {
  if (uang === null || uang === undefined) return ''
  if (typeof uang.amount === 'number') {
    // Kalau ini sampai terjadi, backend-nya yang salah. Lebih baik layar gagal
    // terang-terangan daripada menampilkan angka yang sudah dibulatkan.
    throw new TypeError('jumlah uang datang sebagai angka JSON; ia harus berupa teks desimal')
  }
  return typeof uang.amount === 'string' ? uang.amount : ''
}

/** Tampilkan uang beserta mata uangnya, misalnya "250000 IDR". */
export function tampilUang(uang: UangMasuk | null | undefined): string {
  const jumlah = jumlahUang(uang)
  if (jumlah === '') return ''
  const mataUang = typeof uang?.currency === 'string' ? uang.currency : ''
  return `${jumlah} ${mataUang}`.trim()
}

// ---------------------------------------------------------------------------
// 4. PANGGILAN KE BACKEND — satu fungsi per endpoint.
// ---------------------------------------------------------------------------

export async function cekKesehatan(): Promise<Kesehatan> {
  const data = await minta<Kesehatan>('/healthz')
  return data
}

/**
 * Ambil satu klaim Life beserta SELURUH baris adjustment-nya — tiket 01 AC-1,
 * bukan hanya baris terakhir. Di backend: GET /api/klaim-life/{id},
 * internal/handlers/klaimlife.go.
 */
export async function ambilKlaimLife(id: string): Promise<Klaim> {
  const data = await minta<Klaim>(`/api/klaim-life/${encodeURIComponent(id)}`)
  return data
}

/**
 * Baca kode status HTTP dari sebuah galat, atau `undefined` bila galatnya bukan
 * dari HTTP (misalnya jaringan putus). Halaman memakainya untuk membedakan
 * "klaim tidak ada" (404) dari kegagalan lain.
 */
export function kodeStatusGalat(err: unknown): number | undefined {
  return err instanceof ApiFailure ? err.status : undefined
}

/**
 * Pesan galat dari badan jawaban, bila ada.
 *
 * ⛔ Dipakai HANYA untuk pesan yang memang milik server dan bermakna bagi
 * pengguna - mis. "Name of bank cannot be empty", kalimat sistem lama yang
 * sengaja dipertahankan. Untuk 5xx pesannya TIDAK dicetak.
 */
export function pesanGalat(err: unknown): string | undefined {
  if (!(err instanceof ApiFailure)) return undefined
  // ⚠️ Hanya pesan yang BENAR-BENAR datang dari backend. `ApiFailure`
  // yang dibuat klien (`BACKEND_TIDAK_TERJANGKAU`) punya pesannya
  // sendiri, dan itu bukan kalimat milik server.
  if (err.detail.code !== 'DITOLAK_BACKEND') return undefined
  const pesan = err.detail.message
  return pesan !== undefined && pesan !== '' ? pesan : undefined
}

// ---------------------------------------------------------------------------
// TIKET 02 — pendaftaran klaim.
// ---------------------------------------------------------------------------

/** Satu calon peserta hasil pencarian. Sumber: repository.CalonPeserta. */
export interface CalonPeserta {
  nomorPremiList: string
  nomorPolis: string
  nomorSertifikat: string
  namaTertanggung: string
  mataUang: string
  edmStatus: string
}

/** Isi satu pendaftaran klaim. */
export interface PermintaanDaftar {
  nomorPremiList: string
  nomorPolis: string
  type: string
  kodeBisnis: string
  mataUang: string
  sertifikat: string[]
}

/** Jawaban pendaftaran yang berhasil. */
export interface HasilDaftar {
  id: string
  nomorKlaim: string
}

/**
 * Mencari calon peserta satu premium list.
 *
 * ⛔ `pl` wajib: tabel peserta berisi 66,8 juta baris dan hanya ber-index pada
 * PL_NUMBER. Backend menolak permintaan tanpa itu dengan 400, dan itu memang
 * yang benar — pencarian tanpa penyaring bukan "pencarian luas", melainkan
 * pemindaian penuh yang menahan basis data.
 */
export async function cariPesertaLife(pl: string, batas = 50): Promise<CalonPeserta[]> {
  const data = await minta<CalonPeserta[] | null>('/api/peserta-life', {
    kueri: { pl, n: batas },
  })
  // Go menulis slice kosong sebagai null; layar menginginkan daftar kosong.
  return data ?? []
}

/**
 * Mendaftarkan klaim baru.
 *
 * ⚠️ Backend dapat menjawab 501 bila cara membentuk nomor klaim belum
 * diputuskan work owner. Itu BUKAN kerusakan, dan pesannya menyebut apa yang
 * ditunggu — layar meneruskannya apa adanya, tidak menggantinya dengan
 * "terjadi kesalahan".
 */
// ⚠️ Parameternya DIGANTI NAMA menjadi `permintaan`: nama lamanya `minta`
// menutupi fungsi klien `minta`, dan pemanggilan di dalamnya menjadi
// memanggil OBJEK, bukan klien. tsc menangkapnya; tanpa tsc ia akan
// menjadi galat saat jalan pada satu-satunya jalur pendaftaran klaim.
export async function daftarKlaimLife(
  permintaan: PermintaanDaftar,
): Promise<HasilDaftar> {
  return await minta<HasilDaftar>('/api/klaim-life', {
    metode: 'POST',
    badan: permintaan,
  })
}

/**
 * Kata status baris yang berarti "masih menunggu keputusan".
 *
 * ⛔ Kata, bukan kode mentah. Ronde pertama membandingkan `kodeStatus === '0'`
 * - salinan KETIGA dari `models.KodeOutstanding` yang tidak terhubung apa pun
 * ke Go, pada medan yang berkas ini sendiri tandai "tidak ditampilkan".
 * Kata ini datang dari `models.StatusBaris.String()` dan berubah bersamanya.
 */
/** Kode mentah STS_REJECT untuk baris yang ditolak. */
export const KODE_DITOLAK = '2'

export const STATUS_OUTSTANDING = 'Outstanding'

/**
 * Menolak satu baris adjustment yang masih Outstanding.
 *
 * ⛔ Penolakan membatalkan BARIS itu saja - klaimnya tetap hidup dan baris
 * pengganti dapat diinput sesudahnya.
 *
 * Jawaban yang mungkin, dan artinya berbeda-beda:
 *   403 bukan ReasLifeAdmin
 *   409 baris sudah diputus
 *   422 klaim belum bernomor (butir o belum diputuskan)
 *   501 tempat jejak audit belum diputuskan (butir am)
 */
export async function tolakBarisAdjustment(klaimID: string, adjID: string): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/adjustment/${encodeURIComponent(adjID)}/tolak`,
    { metode: 'POST' },
  )
}

/**
 * Menyimpan akseptasi baris adjustment TERAKHIR seorang peserta.
 *
 * ⭐ Jalur ini LAHIR dari audit XML: `SaveAdjustment_Act` mengaksep di Claim
 * Life sendiri, dan sebelumnya kami hanya mengenal jalur Komite.
 *
 * Yang boleh menekannya adalah PEMEGANG TAHAP klaim — bukan satu peran tetap.
 * Layar tidak tahu peran siapa pun; yang menolak adalah services (403).
 *
 * Jawaban yang mungkin:
 *   403 bukan pemegang tahap klaim ini
 *   409 peserta belum dipilih, baris sudah bernomor, atau bukan Outstanding
 *   501 kode bisnis belum tersimpan (temuan audit A0), atau tempat jejak belum ada
 */
export async function simpanAdjustment(
  klaimID: string,
  pesertaID: string,
): Promise<string> {
  const jawab = await minta<{ nomorAkseptasi: string }>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/akseptasi`,
    { metode: 'POST' },
  )
  return jawab.nomorAkseptasi
}

/**
 * Apakah layar menampilkan kontrol "Save Adjustment" bagi seorang peserta.
 *
 * ⛔ Meniru prasyarat XML `SaveAdjustment_Act` (pecahan 854, 1048): peserta
 * DIPILIH, baris terakhir belum bernomor akseptasi, dan masih Outstanding.
 *
 * ⚠️ Peran TIDAK diperiksa di sini, dan itu sesuai XML: pohon section
 * membuktikan tombolnya tidak dibungkus gerbang peran mana pun. Yang
 * menggerbanginya pemegang tahap, dan itu hanya services yang tahu.
 */
export function bolehSimpanAdjustment(
  peserta: { isCheck: string },
  baris: BarisAdjustment[],
): boolean {
  if (peserta.isCheck.toLowerCase() !== 'true') return false
  const terakhir = baris[baris.length - 1]
  return (
    terakhir !== undefined &&
    terakhir.nomorAkseptasi === '' &&
    terakhir.status === STATUS_OUTSTANDING
  )
}

/**
 * Membuka putaran adjustment berikutnya bagi seorang peserta.
 *
 * Inilah yang membuat klaim TIDAK TERMINAL: penolakan Komite menghasilkan
 * putaran berikutnya, bukan akhir (ADR-U-0011). Yang terminal adalah BARIS.
 *
 * Jawaban yang mungkin:
 *   403 bukan ReasLifeSPV
 *   409 baris terakhir belum ditolak
 *   501 tempat jejak audit belum diputuskan (butir am)
 */
export async function tambahPutaran(klaimID: string, pesertaID: string): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/putaran`,
    { metode: 'POST' },
  )
}

/**
 * Apakah layar menampilkan kontrol "Putaran berikutnya" bagi seorang peserta.
 *
 * ⛔ Hanya keadaan BARIS TERAKHIR yang menentukan - dan hanya penolakan yang
 * membuka putaran baru. Seperti kontrol lain di layar ini, ia kenyamanan
 * bukan pagar: yang menolak adalah services (403/409).
 */
export function bolehPutaranBaru(baris: BarisAdjustment[]): boolean {
  const terakhir = baris[baris.length - 1]
  return terakhir !== undefined && terakhir.kodeStatus === KODE_DITOLAK
}

/**
 * Apakah layar menampilkan kontrol "Send ke Komite" untuk sebuah baris.
 *
 * ⛔ Sengaja TIDAK memeriksa rekening pembayaran, dan tidak memeriksa peran.
 * Layar tidak tahu peran siapa pun, dan gerbang rekening adalah aturan dagang
 * - aturan dagang yang ditegakkan di layar dapat dilewati siapa pun yang
 * memanggil API langsung. Kontrol ini kenyamanan, BUKAN pagar; yang menolak
 * adalah services (403/422).
 */
export function bolehSerahkanDiLayar(b: BarisAdjustment): boolean {
  return b.status === STATUS_OUTSTANDING && b.komiteId === ''
}

/**
 * Menyerahkan satu baris adjustment ke Komite Life.
 *
 * Komite adalah SISTEM LUAR. Yang diserahkan adalah barisnya, bukan klaimnya:
 * unit keputusan tetap baris (ADR-U-0011).
 *
 * Jawaban yang mungkin, dan artinya berbeda-beda:
 *   403 peran tidak berwenang untuk Type ini (QP/QR hanya SPV)
 *   409 baris sudah diserahkan, atau bukan lagi Outstanding
 *   422 rekening pembayaran belum lengkap, mata uang campur, atau roster kosong
 *   501 tempat roster/kasus komite belum diputuskan work owner (butir af)
 */
export async function serahkanKeKomite(
  klaimID: string,
  pesertaID: string,
  adjID: string,
): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}` +
      `/adjustment/${encodeURIComponent(adjID)}/komite`,
    { metode: 'POST' },
  )
}

/** Cacah baris yang akan ikut terhapus bersama sebuah klaim. */
export interface DampakHapus {
  header: number
  peserta: number
  adjustment: number
  spreading: number
  spreadingRetro: number
  dokumen: number
  barisWork: number
  /**
   * Baris `OS_AKSEPTASI_KLAIM_LIFE` ber-CASEID sama.
   *
   * ⚠️ Nasib baris ini saat klaim dihapus BELUM diputuskan work owner.
   * Angkanya ditampilkan justru supaya keputusan yang belum diambil itu
   * terlihat oleh yang menekan tombol, bukan tersembunyi di dalam kaskade.
   */
  barisDatarWarisan: number
  /** ⛔ SENGAJA tidak memuat `barisDatarWarisan` - lihat di atas. */
  total: number
}

/**
 * Menghitung dampak penghapusan - TANPA menghapus apa pun.
 *
 * Inilah yang membuat "Batal" benar-benar membatalkan: jalur ini tidak punya
 * satu pun tulisan untuk dibatalkan.
 */
export async function dampakHapusKlaim(klaimID: string): Promise<DampakHapus> {
  const data = await minta<DampakHapus>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/dampak-hapus`,
  )
  return data
}

/**
 * Menghapus klaim.
 *
 * ⛔ Saat ini SELALU menjawab 501: ADR-U-0031 menetapkan penghapusan berupa
 * PENANDA, bukan hapus fisik, dan kolom penandanya belum diputuskan.
 */
export async function hapusKlaim(klaimID: string): Promise<DampakHapus> {
  // ⛔ DELETE, bukan GET. Rutenya sama persis dengan `ambilKlaimLife`;
  // yang membedakan hanya METODEnya, dan metode yang hilang membuat
  // tombol Hapus diam-diam MEMBACA klaim lalu melaporkan berhasil.
  const data = await minta<DampakHapus>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}`,
    { metode: 'DELETE' },
  )
  return data
}

// ---------------------------------------------------------------------------
// F0.1 — ENVELOPE GALAT, bentuk yang dipakai komponen dasar.
//
// ⛔ PILIHAN YANG DINYATAKAN (brief lanjutan 6 §2 aturan 2): `api.ts` kita
// DIMIGRASI UTUH ke klien `fetch` referensi beserta `ApiFailure`/`bolehUlang`,
// dan `axios` dilepas. Alasannya bukan selera:
//
//   - `components/ui/dasar.tsx` dan `lib/keadaanGalat.ts` yang diadopsi dari
//     referensi keduanya berbicara dalam `ApiFailure`. Mempertahankan axios
//     berarti menulis adaptor yang MEREPRODUKSI `ApiFailure` - yaitu menulis
//     kelas yang sama, ditambah satu lapis lagi yang dapat salah.
//   - Dua model galat berdampingan berarti dua jalan menampilkan kegagalan
//     yang sama, dan yang satu akan diam-diam kalah.
//
// ⚠️ Migrasi transportnya dilakukan F0.2. Yang lahir di F0.1 hanya BENTUK
// galatnya, supaya komponen dasar dapat diperiksa tipe-nya.
//
// ⚠️ BENTUKNYA TIDAK DISALIN MENTAH dari referensi. Envelope backend kita
// `{"error": "<kalimat>"}` - satu medan, tanpa `code`, tanpa `fields`; itu
// yang `handlers.galat` tulis. `ApiFailure` di sini memetakan envelope KITA,
// bukan envelope aplikasi lain. Menyalin `ErrorCode` beserta kedua belas
// kodenya berarti menjanjikan kode yang backend kita tidak pernah kirim.
// ---------------------------------------------------------------------------

/** Kode galat yang benar-benar dapat muncul di jalur kita. */
export type KodeGalatApi =
  /** Dipasang KLIEN, bukan backend: jawaban bukan JSON sama sekali.
   *
   *  ⛔ Ia berdiri di sini supaya keadaan "backend mati" punya NAMA alih-alih
   *  menjadi SyntaxError yang tak terbaca siapa pun. */
  | 'BACKEND_TIDAK_TERJANGKAU'
  /** Backend menjawab JSON, dan medannya `error`. */
  | 'DITOLAK_BACKEND'

/** Satu galat per medan. Backend kita belum mengirimnya; bentuknya disiapkan
 *  supaya layar tidak perlu diubah ketika ia mengirimkannya. */
export interface GalatMedan {
  field: string
  message: string
}

/** Isi envelope galat, sesudah dinormalkan dari `{"error": …}`. */
export interface IsiGalatApi {
  code?: KodeGalatApi
  message?: string
  fields?: GalatMedan[]
}

/**
 * Galat yang membawa envelope backend apa adanya.
 *
 * ⚠️ Nama medannya (`status`, `detail`) mengikuti referensi DENGAN SENGAJA:
 * `lib/keadaanGalat.ts` mengenali bentuk ini secara struktural, tanpa mengimpor
 * kelasnya, supaya modul itu tetap murni dan dapat diuji tanpa DOM.
 */
export class ApiFailure extends Error {
  constructor(
    readonly status: number,
    readonly detail: IsiGalatApi,
  ) {
    super(detail.message ?? 'Permintaan ditolak backend')
    this.name = 'ApiFailure'
  }

  /** Pesan untuk satu medan, bila backend melaporkannya. */
  fieldMessage(nama: string): string | undefined {
    return this.detail.fields?.find((f) => f.field === nama)?.message
  }

  /** SELURUH medan yang gagal - bukan yang pertama saja. */
  get fields(): GalatMedan[] {
    return this.detail.fields ?? []
  }
}

// ---------------------------------------------------------------------------
// F0.4 — KOTAK MASUK per tahap.
//
// Bentuknya dari `Claim Life/ReportDefinition/InboxPremiumList.xml`:
// 50 baris per halaman (592), maksimum 500 (942), urut waktu BUAT menurun
// (733 · 736).
// ---------------------------------------------------------------------------

/** Satu baris antrian. Sumber: handlers.barisInboxJSON. */
export interface BarisInbox {
  /** Label `Case ID` — `InboxPremiumList.xml:721`. */
  caseId: string
  /** Pengenal work; dipakai membuka kasusnya. */
  id: string
  /** Nama assignment VERBATIM (butir at). */
  tahap: string
  /** Label `Create Operator Name` — baris 751. */
  createOpName: string
  /** Label `Create Date/Time` — baris 736. RFC3339, atau kosong. */
  tglCreate: string
  /** Label `Work Status` — baris 765. KATA, bukan kode mentah. */
  status: string
  nomorKlaim: string
  nomorPolis: string
  namaBisnis: string
  mataUang: string
}

export interface HalamanInbox {
  baris: BarisInbox[]
  /** Cacah SELURUH kasus pada tahap itu — lencana tab memakainya. */
  total: number
  tahap: number
  halaman: number
  ukuran: number
  namaTahap: string
}

/**
 * Nomor tahap — HARUS sama dengan `models.Tahap` di Go.
 *
 * ⛔ Angka di kabel, kata di layar. Nama tahapnya datang dari backend
 * (`namaTahap`), bukan dikarang di sini: ia nama assignment VERBATIM dan
 * mengetiknya ulang di dua tempat berarti dua tempat untuk salah.
 */
export const TAHAP_NOMOR = {
  inputRegister: 1,
  outstanding: 2,
  medicalCheck: 3,
  claimAnalis: 4,
} as const

export type NomorTahap = (typeof TAHAP_NOMOR)[keyof typeof TAHAP_NOMOR]

/** Membaca satu halaman antrian sebuah tahap. */
export async function ambilKotakMasuk(
  tahap: NomorTahap,
  halaman = 1,
  ukuran?: number,
): Promise<HalamanInbox> {
  return await minta<HalamanInbox>('/api/klaim-life', {
    kueri: { tahap, halaman, ukuran },
  })
}
