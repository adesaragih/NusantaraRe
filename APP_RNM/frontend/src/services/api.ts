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
import axios from 'axios'

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
// 2. KLIEN HTTP.
//
// Alamat backend TIDAK ditulis di kode. Ia dibaca dari env var
// VITE_API_BASE_URL (berkas .env — contohnya di .env.example). Kosong berarti
// frontend dan backend disajikan dari alamat yang sama (ADR-U-0004).
// ---------------------------------------------------------------------------
const baseURL = import.meta.env.VITE_API_BASE_URL ?? ''

export const api = axios.create({
  baseURL,
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

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
  const { data } = await api.get<Kesehatan>('/healthz')
  return data
}

/**
 * Ambil satu klaim Life beserta SELURUH baris adjustment-nya — tiket 01 AC-1,
 * bukan hanya baris terakhir. Di backend: GET /api/klaim-life/{id},
 * internal/handlers/klaimlife.go.
 */
export async function ambilKlaimLife(id: string): Promise<Klaim> {
  const { data } = await api.get<Klaim>(`/api/klaim-life/${encodeURIComponent(id)}`)
  return data
}

/**
 * Baca kode status HTTP dari sebuah galat, atau `undefined` bila galatnya bukan
 * dari HTTP (misalnya jaringan putus). Halaman memakainya untuk membedakan
 * "klaim tidak ada" (404) dari kegagalan lain.
 */
export function kodeStatusGalat(err: unknown): number | undefined {
  return axios.isAxiosError(err) ? err.response?.status : undefined
}

/**
 * Pesan galat dari badan jawaban, bila ada.
 *
 * ⛔ Dipakai HANYA untuk pesan yang memang milik server dan bermakna bagi
 * pengguna - mis. "Name of bank cannot be empty", kalimat sistem lama yang
 * sengaja dipertahankan. Untuk 5xx pesannya TIDAK dicetak.
 */
export function pesanGalat(err: unknown): string | undefined {
  if (!axios.isAxiosError(err)) return undefined
  const data = err.response?.data as { error?: unknown } | undefined
  return typeof data?.error === 'string' && data.error !== '' ? data.error : undefined
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
  const { data } = await api.get<CalonPeserta[] | null>('/api/peserta-life', {
    params: { pl, n: batas },
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
export async function daftarKlaimLife(minta: PermintaanDaftar): Promise<HasilDaftar> {
  const { data } = await api.post<HasilDaftar>('/api/klaim-life', minta)
  return data
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
  await api.post(`/api/klaim-life/${encodeURIComponent(klaimID)}/adjustment/${encodeURIComponent(adjID)}/tolak`)
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
  const jawab = await api.post<{ nomorAkseptasi: string }>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/akseptasi`,
  )
  return jawab.data.nomorAkseptasi
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
  await api.post(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/putaran`,
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
  await api.post(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}` +
      `/adjustment/${encodeURIComponent(adjID)}/komite`,
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
  const { data } = await api.get<DampakHapus>(
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
  const { data } = await api.delete<DampakHapus>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}`,
  )
  return data
}
