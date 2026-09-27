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

/**
 * Keenam total uang SATU peserta.
 *
 * ⛔ TURUNAN, bukan kolom. Tidak ada `TOTAL_*` di basis data; angkanya
 * dihitung backend saat Detail dibaca, dari SELURUH baris adjustment peserta
 * itu - `SavePesertaClaim.xml` langkah 8.1 b4221 dan 8.2 b4592, yang
 * keduanya berprasyarat kosong sehingga baris DITOLAK pun ikut.
 *
 * Urutan medan mengikuti urutan layar: b20629, b20914, b21201, b21488,
 * b21775, b22063.
 */
export interface TotalPeserta {
  cedingRetention: Uang
  shareNusantaraRe: Uang
  sumInsured: Uang
  sumReasured: Uang
  shareRetro: Uang
  jumlahKlaim: Uang
}

/**
 * Satu dokumen pendukung milik seorang peserta.
 *
 * Sumber: `T_CLAIMLF_DOCUMENT`, dirakit backend saat Detail dibaca — meniru
 * `Activity/LoadDocumentLife_ACT.xml`.
 *
 * ⛔ `urlPublic` TIDAK ADA di sini dan itu disengaja. Di Pega URL-nya
 * datang dari `GetUrlGoogleStorage_Act` (langkah 1.4.1 b1011); penyambungan
 * ke penyimpanan luar belum dilakukan, dan URL yang dikarang adalah pranala
 * yang membawa orang ke tempat yang salah. Yang ada hanya penunjuknya,
 * `tStorageId` — ADR-U-0010: isi berkas tidak pernah masuk basis data.
 */
export interface Dokumen {
  /**
   * ⛔ TEKS, bukan `number`, dan itu bukan pilihan gaya.
   *
   * Nilainya cap waktu `yyyyMMddhhmmssSSS` — **17 angka** ≈ 2,0e16 —
   * sedangkan `Number.MAX_SAFE_INTEGER` ≈ 9,0e15. `JSON.parse`
   * membulatkannya diam-diam: `20260927103000123` menjadi
   * `...124`. Tautan unduh lalu menunjuk dokumen yang tidak ada dan
   * penghapusan mengenai baris yang salah, **tanpa satu pun galat**.
   *
   * Go mengirimnya dengan tag `,string`; dikunci dua sisi
   * (`TestPengenalDokumenMenyeberangSebagaiTeks` di Go, uji tautan di
   * `PanelDokumenPeserta.test.ts` di sini).
   */
  id: string
  pesertaId: string
  /** Nama BERKAS unggahan. Bukan nama orang. */
  namaFile: string
  mime: string
  kategori1: string
  /** Yang ditampilkan kepada manusia di `DocumentLife.xml`. */
  kategori2: string
  /** Penunjuk berkas di penyimpanan luar. Kosong berarti belum terunggah. */
  tStorageId: string
  /**
   * Tanggal dokumen. ⛔ `string | null`, dan `null` BERARTI SESUATU.
   *
   * Go mengirimnya dari `*time.Time`, jadi kolom NULL menyeberang sebagai
   * `null` — bukan `""`. Ronde pertama menulis `tanggal: string`, dan
   * TypeScript tidak dapat menangkapnya: tipe yang berbohong tentang data
   * yang datang dari jaringan tetap dikompilasi. Akibatnya
   * `b.tanggal === ''` tidak pernah menyala, dan tanggal yang memang kosong
   * tampil sebagai sel kosong alih-alih penanda — persis kebalikan dari
   * yang ADR-U-0027 minta.
   *
   * ⚠️ Ini cacat lintas-lapis KELIMA dengan bentuk yang sama di modul ini,
   * dan ia lahir di giliran yang sama ketika keempat pendahulunya
   * didaftarkan. Dikunci dua sisi: `TestNamaJSONDokumenDikunci` di Go
   * memeriksa `"tanggal":null`, dan `PanelDokumenPeserta.test.ts` memeriksa
   * `null` diperlakukan sama dengan kosong.
   */
  tanggal: string | null
  paymentDate: string | null
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
  /**
   * `DATE_OF_LOSS` — tanggal kejadian, milik PESERTA dan bukan klaim.
   *
   * Teks apa adanya, bentuk `YYYY-MM-DD`. Kosong berarti belum diisi.
   */
  tanggalKejadian: string
  baris: BarisAdjustment[]
  /** Keenam total uang peserta ini, dihitung backend. Lihat TotalPeserta. */
  total: TotalPeserta
  /** Dokumen pendukung peserta ini. Daftar kosong, tidak pernah null. */
  dokumen: Dokumen[]
  /**
   * Diagnosa peserta ini — butir bd. Daftar kosong, tidak pernah null.
   *
   * Daftar, bukan sepasang kolom: b3923 menyajikannya RepeatGrid dengan
   * tombol Add b4690 dan Delete b6160.
   */
  diagnosa: Diagnosa[]
  /** Cermin STS_REJECT peserta. Kosong berarti belum diputus. */
  kodeStatus: string
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
  /**
   * Nama assignment VERBATIM dari `T_WORK_CLAIM.TAHAP`.
   *
   * ⚠️ Boleh KOSONG: backend tidak pernah menyisipkan ke `T_WORK_CLAIM`, jadi
   * klaim tanpa baris work adalah keadaan nyata. Kosong = tahap tidak
   * diketahui, dan tombol yang bergantung tahap TIDAK ditawarkan.
   */
  tahap: string
  /**
   * `T_WORK_CLAIM.STATUS_WORK` — butir bb.
   *
   * Kosong berarti kasus BELUM ditutup. Satu-satunya nilai lain yang pernah
   * ditulis adalah `STATUS_WORK_SELESAI`.
   */
  statusWork: string
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
    const o = (isi ?? {}) as { galat?: unknown; penghalang?: unknown }
    throw new ApiFailure(jawab.status, {
      code: 'DITOLAK_BACKEND',
      // ⛔ Kuncinya `galat`, bukan `error`. `handlers.galat` di Go menulis
      // `{"galat": "<kalimat>"}` sejak tiket 01; baris ini sempat membaca
      // `error`, sehingga SETIAP pesan backend - 401, 403, 409, 503 - jatuh
      // ke teks bawaan dan pemakai melihat pita merah yang tidak
      // menyebutkan apa pun.
      //
      // ⛔ `error` sengaja TIDAK ikut diterima. Menerima kedua kunci akan
      // menambal gejalanya dan menyembunyikan sebabnya: sejak itu kedua sisi
      // tidak pernah dipaksa bertemu lagi. Kontraknya dikunci dua sisi -
      // `envelopegalat.test.ts` di sini, `envelopegalat_test.go` di Go.
      message: typeof o.galat === 'string' && o.galat !== '' ? o.galat : undefined,
      // ⛔ SELURUH penghalang dibawa, bukan yang pertama. Pega memasang
      // pesannya di dalam loop, sekali per peserta yang tertandai; melaporkan
      // satu saja memaksa pemakai menutup berulang kali dan menemukan satu
      // penghalang baru setiap kali.
      penghalang: Array.isArray(o.penghalang)
        ? (o.penghalang as PenghalangTutup[])
        : undefined,
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
export interface SaringPeserta {
  /** `SearchPolicyHolder.CARI2` — `CERTIFICATE_NO LIKE %..%`, TIDAK di-uppercase. */
  sertifikat?: string
  /** `SearchPolicyHolder.CARI3` — `UPPER(NAME_OF_INSURED) LIKE %..%` (b405). */
  nama?: string
}

/**
 * Mencari calon peserta satu premium list.
 *
 * Ketiga kriterianya dari `RDBList/GetPesertaClaim_sql1.xml:85`, dipanggil
 * `LoadDataPesertaSpesifik_Act` langkah `RDB-List` b485. `pl` WAJIB: tabel
 * sumbernya 66,8 juta baris dan hanya ber-index pada PL_NUMBER,
 * CERTIFICATE_NO, dan POLICY_NO.
 *
 * ⚠️ Kotak yang dibiarkan kosong TIDAK menyaring. Pega memasang kedua `LIKE`
 * tanpa syarat, dan di Oracle `X LIKE '%'` bernilai FALSE saat X NULL —
 * sehingga kotak kosong di sana diam-diam membuang peserta ber-nama NULL.
 * Penyimpangan sadar, dilaporkan OQ-E.
 */
export async function cariPesertaLife(
  pl: string,
  saring: SaringPeserta = {},
  batas = 50,
): Promise<CalonPeserta[]> {
  const kueri: Record<string, string | number> = { pl, n: batas }
  // Hanya kirim yang terisi: parameter kosong dan parameter tidak dikirim
  // harus berarti hal yang SAMA, dan cara paling aman menjamin itu adalah
  // tidak pernah mengirim yang kosong.
  for (const [nama, isi] of [
    ['sertifikat', saring.sertifikat],
    ['nama', saring.nama],
  ] as const) {
    const bersih = isi?.trim() ?? ''
    if (bersih !== '') kueri[nama] = bersih
  }
  const data = await minta<CalonPeserta[] | null>('/api/peserta-life', { kueri })
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

/**
 * Kode mentah STS_REJECT untuk baris yang diaksep.
 *
 * ⛔ Ia lahir bersama `bolehUbahDiagnosa`, dan sebabnya tajam: gerbang
 * `pyDisabledWhen` b4682 menyebut **dua** nilai dengan `||`, dan layar yang
 * hanya mengenal yang ditolak akan membiarkan diagnosa peserta yang sudah
 * DIAKSEP tetap dapat disunting - separuh gerbang hilang tanpa berbunyi.
 */
export const KODE_AKSEP = '1'

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
// `{"galat": "<kalimat>"}` - satu medan, tanpa `code`, tanpa `fields`; itu
// yang `handlers.galat` tulis. (Komentar ini sempat menyebut `error`, dan
// kekeliruan itulah yang menuntun kode di bawah membaca kunci yang salah.) `ApiFailure` di sini memetakan envelope KITA,
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
  /**
   * Daftar peserta yang menahan penutupan — hanya pada 409 dari
   * `POST /api/klaim-life/{id}/tutup`.
   *
   * ⛔ Ia ikut di amplop yang SAMA, bukan di bentuk kedua. Cacat `galat` vs
   * `error` lahir persis dari dua bentuk amplop yang masing-masing benar
   * menurut dirinya sendiri; menambah bentuk ketiga untuk satu rute akan
   * mengulanginya. Bentuknya `{ galat, penghalang }`, dikunci
   * `TestBadanPenghalangMemakaiAmplopYangSama` di Go.
   */
  penghalang?: PenghalangTutup[]
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

// ---------------------------------------------------------------------------
// A3 Outstanding — perpindahan tahap (butir aw).
// ---------------------------------------------------------------------------

/**
 * Tahap tujuan sebagai KATA — sama dengan `handlers.tahapTujuan`.
 *
 * ⛔ Kata, bukan angka: jalur `/tahap/2` tidak dapat dibaca siapa pun, dan
 * angka yang bergeser bila urutan `models.Tahap` berubah akan memindahkan
 * kasus ke tempat yang salah tanpa satu pun galat.
 */
export const TAHAP_JALUR = {
  inputRegister: 'input-register',
  outstanding: 'outstanding',
  medicalCheck: 'medical-check',
  claimAnalis: 'claim-analis',
} as const

export type TahapJalur = (typeof TAHAP_JALUR)[keyof typeof TAHAP_JALUR]

/**
 * Memindahkan kasus ke tahap lain.
 *
 * Kode jawaban yang mungkin:
 *   204 berhasil
 *   403 pelaku bukan pemegang tahap ASALnya
 *   409 perpindahan itu tidak ada di tangga kerja
 */
export async function pindahTahap(klaimID: string, tujuan: TahapJalur): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/tahap/${tujuan}`,
    { metode: 'POST' },
  )
}

// ---------------------------------------------------------------------------
// TIKET 06 — tanggal kejadian (DOL), dari layar Detail.
// ---------------------------------------------------------------------------

/**
 * Mengubah tanggal kejadian satu peserta.
 *
 * Tombol `Edit Date` `Section/ClaimLifeDetailGCNM.xml` b14115 →
 * `pyLocalAction ShowEditClaimLife` b14144 → `EditDateClaimLife_Section`.
 * Validasinya `ValidasiDOL_Act`, dipanggil section yang sama di b11177 dan
 * b11298.
 *
 * ⚠️ Fungsi ini lahir belakangan daripada rutenya. Rute backendnya ada sejak
 * tiket 06 lengkap dengan `ValidasiDOL` dan ujinya, tetapi sampai kelompok
 * Detail & Tutup NOL pemanggil di React — jadi jalurnya tidak pernah dapat
 * dijalankan siapa pun, dan itu tidak berbunyi di uji mana pun: backend
 * hijau, layar hijau, fiturnya tidak ada.
 *
 * Kode jawaban yang mungkin:
 *   204 berhasil
 *   400 tanggalnya bukan tanggal yang dikenal, atau belum diisi
 *   403 peran tidak mencukupi
 *   422 DOL di luar jendela valuasi — badannya membawa `pesertaId`
 */
export async function ubahTanggalKejadian(
  klaimID: string,
  pesertaID: string,
  tanggalKejadian: string,
): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/tanggal-kejadian`,
    { metode: 'PUT', badan: { tanggalKejadian } },
  )
}

// ---------------------------------------------------------------------------
// Gerbang Close Claim — `CloseClaim_Section.xml` b1081 -> b1101.
// ---------------------------------------------------------------------------

/** Satu peserta yang menahan penutupan klaim. */
export interface PenghalangTutup {
  urutan: number
  nomorSertifikat: string
  /** Kalimat yang dilihat pemakai, disusun SERVER dan verbatim dari rule. */
  pesan: string
}

/** Jawaban gerbang tutup. */
export interface HasilPeriksaTutup {
  boleh: boolean
  penghalang: PenghalangTutup[]
}

/**
 * Memeriksa apakah klaim boleh ditutup.
 *
 * ⛔ GET, dan ia BELUM MENUTUP apa pun. Di Pega tombol `Close Claim`
 * menjalankan dua aksi pada satu klik — `refresh` → `ProtectCloseClaim_act`
 * (b1101) yang memeriksa lalu memanggil `FinishAssignment`, dan
 * `closeContainer` (b1129) yang menutup jendelanya. Yang dibangun di sini
 * baru pemeriksaannya; sisi penyelesaian penugasan menunggu pembacaan
 * `Flow/`, sebab tahap tujuannya belum diketahui.
 */
export async function periksaBolehTutup(klaimID: string): Promise<HasilPeriksaTutup> {
  return minta<HasilPeriksaTutup>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/boleh-tutup`,
  )
}

/**
 * Status kerja kasus yang sudah ditutup — VERBATIM `Register_Flow.xml` b899.
 *
 * ⛔ Satu-satunya nilai yang pernah ditulis. Status Pega untuk kasus yang
 * sedang berjalan tidak ada di ekspor dan tidak dikarang: kasus terbuka
 * berkolom kosong.
 */
export const STATUS_WORK_SELESAI = 'Resolved-Completed'

/** Kedua tahap yang layarnya menawarkan `Close Claim` — cacah berkas. */
export const TAHAP_PENAWAR_TUTUP = ['Outstanding Claim', 'Claim Analis']

/**
 * Apakah tombol `Close Claim` pantas ditawarkan untuk klaim ini.
 *
 * ⛔ Dua syarat, dan keduanya dari XML:
 *
 *   1. tahapnya menawarkannya — `pyLocalAction>CloseClaim` hanya ada di
 *      `InputOSClaimLife.xml` (b22837, b22988) dan
 *      `InputAkseptasiClaimLife.xml` (b21457, b21602);
 *   2. kasusnya belum tertutup.
 *
 * ⚠️ Ini HANYA menentukan tampil atau tidaknya tombol. Gerbang sebenarnya
 * ada di backend dan diperiksa lagi di sana — layar yang menjadi satu-satunya
 * penjaga adalah layar yang dapat dilewati dengan satu permintaan.
 */
export function bolehTutupDiLayar(klaim: Klaim | null): boolean {
  if (klaim === null) return false
  if (klaim.statusWork === STATUS_WORK_SELESAI) return false
  return TAHAP_PENAWAR_TUTUP.includes(klaim.tahap)
}

/** Apakah kasus ini sudah ditutup dan karena itu tidak dapat diubah lagi. */
export function kasusTertutup(klaim: Klaim | null): boolean {
  return klaim !== null && klaim.statusWork === STATUS_WORK_SELESAI
}

/**
 * Daftar penghalang dari sebuah galat 409, bila ada.
 *
 * Kosong berarti galat itu bukan penolakan gerbang — mis. 403 atau 503 —
 * dan pemanggil harus menampilkan pesannya, bukan daftar kosong.
 */
export function penghalangDariGalat(err: unknown): PenghalangTutup[] {
  if (!(err instanceof ApiFailure)) return []
  return err.detail.penghalang ?? []
}

/**
 * Menutup kasus — `POST /api/klaim-life/{id}/tutup`.
 *
 * ⛔ Keputusan bb. Di Pega `Close Claim` adalah LOCAL ACTION, dan alurnya
 * tidak punya konektor senama; perilaku mesinnya untuk `FinishAssignment`
 * dari local action tidak dapat diturunkan dari ekspor (OQ-I). Yang ditiru
 * adalah niat nyatanya — label, konfirmasi b499, dan gerbang "not approved
 * yet" hanya masuk akal bila tombolnya MENUTUP.
 *
 * Jawaban 409 membawa SELURUH penghalang, bukan yang pertama; pemanggil
 * membacanya lewat `penghalangDariGalat`.
 */
export async function tutupKlaim(klaimID: string): Promise<void> {
  await minta<void>(`/api/klaim-life/${encodeURIComponent(klaimID)}/tutup`, {
    metode: 'POST',
  })
}

// ---------------------------------------------------------------------------
// KELOMPOK MEDIS — pencarian diagnosa.
// ---------------------------------------------------------------------------

/** Satu baris `DISEASE_LIFE`. */
export interface Penyakit {
  /** Pengenal sebagai TEKS — "007" bukan "7" (ADR-U-0022). */
  nomor: string
  nama: string
  kodeIcd: string
}

/** `pyMaxRecords` b659 — batas atas mutlak, dari rule. */
export const BATAS_BARIS_PENYAKIT = 500

/** `pyPageSize` b514. */
export const UKURAN_HALAMAN_PENYAKIT = 50

/**
 * Mencari diagnosa — `GET /api/penyakit-life`.
 *
 * ⛔ `DISEASE_LIFE` berisi **97.586 baris**. Batasnya dijepit LAGI di
 * backend ke `pyMaxRecords` 500; yang di sini hanya supaya permintaannya
 * masuk akal sejak berangkat.
 *
 * ⛔ Kedua kata kunci disambung **AND** di backend (b535 `A AND B`), dan
 * keduanya dinaikkan ke huruf besar (`@toUpperCase` b255, b302). Pencarian
 * dengan kedua kata kunci KOSONG sah — itu yang Pega lakukan, dan batasnya
 * yang menahan.
 */
export async function cariPenyakit(
  kodeIcd: string,
  nama: string,
  batas = UKURAN_HALAMAN_PENYAKIT,
): Promise<Penyakit[]> {
  return minta<Penyakit[]>('/api/penyakit-life', {
    kueri: {
      icd: kodeIcd,
      nama,
      batas: String(Math.min(Math.max(batas, 1), BATAS_BARIS_PENYAKIT)),
    },
  })
}


// ---------------------------------------------------------------------------
// Grid diagnosa per peserta — butir bd.
//
// ⛔ Tiga rute, tiga tombol: `Add` b4690, `Choose` b2509, `Delete` b6160.
// Jalurnya BERSARANG di bawah pesertanya — `SetDisease.xml` b389 menutup
// dengan `Obj-Save pyWorkPage`, jadi diagnosa tidak punya hidup di luar
// peserta yang memuatnya.
// ---------------------------------------------------------------------------

/** Satu baris `.DiagnoseList` milik seorang peserta. */
export interface Diagnosa {
  /**
   * Pengenal ANGKA — berbeda dengan pengenal klaim dan peserta, yang teks.
   *
   * ⚠️ Ia milik kita (`SEQ_CLAIMLF_DIAGNOSE`), bukan warisan Pega: di sana
   * baris ini hanya punya subscript di dalam halaman induknya. ADR-U-0022
   * berlaku atas KODE yang datang dari sistem lama, bukan atas identitas
   * yang kita terbitkan sendiri.
   */
  id: number
  pesertaId: string
  /** Posisi di grid, mulai 1. Dirapatkan backend sesudah penghapusan. */
  urutan: number
  /** `ICD_CODE` — read-only di layar (b5566), diisi dari hasil pencarian. */
  kodeIcd: string
  /** `DISEASE` — read-only di layar (b5374). */
  nama: string
  /** `GROUP_DIAGNOSE`. ⛔ Daftar pilihannya belum ada — OQ-L, butir bf. */
  groupDiagnose: string
  /** Cermin `STS_REJECT` PESERTA, bukan keputusan baris ini sendiri. */
  kodeStatus: string
}

/**
 * Apakah diagnosa peserta ini masih boleh disunting.
 *
 * ⛔ VERBATIM `pyDisabledWhen` — satu kalimat di TUJUH tempat pada
 * `ClaimLifeDetailGCNM.xml`, empat di antaranya di grid ini (b4682 `Add`,
 * b5059 `Find Disease`, b5870 `GROUPDIAGNOSE`, b6152 `Delete`):
 *
 *     .STS_REJECT=='1' || .STS_REJECT=='2'
 *
 * ⛔ Yang diuji `STS_REJECT` **PESERTA**. Layar bukan penjaga — backend
 * menolak permintaannya dengan 409 — tetapi tombol yang tetap hidup padahal
 * pasti ditolak adalah tombol yang mengajari orang mengabaikan galat.
 */
export function bolehUbahDiagnosa(peserta: Peserta): boolean {
  return peserta.kodeStatus !== KODE_AKSEP && peserta.kodeStatus !== KODE_DITOLAK
}

/**
 * `Add` b4690 — menambah baris KOSONG di ekor daftar peserta.
 *
 * Mengembalikan baris yang baru lahir, lengkap dengan `id` dan `urutan`,
 * supaya layar tidak perlu membaca ulang seluruh klaim hanya untuk
 * menampilkan satu baris kosong.
 */
export async function tambahDiagnosa(
  klaimID: string,
  pesertaID: string,
): Promise<Diagnosa> {
  return minta<Diagnosa>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/diagnosa`,
    { metode: 'POST' },
  )
}

/**
 * `Choose` b2509 → `SetDisease` — menulis isi satu baris.
 *
 * ⛔ `kodeIcd` dan `nama` datang dari baris hasil `GET /api/penyakit-life`,
 * tidak diketik: keduanya `Read-only` di grid (b5374, b5566). Mengetiknya
 * berarti nama diagnosa yang tidak ada di katalog dapat masuk, dan tidak ada
 * satu pun yang akan membandingkannya lagi.
 */
export async function ubahDiagnosa(
  klaimID: string,
  pesertaID: string,
  diagID: number,
  isi: { kodeIcd: string; nama: string; groupDiagnose: string },
): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/diagnosa/${diagID}`,
    { metode: 'PUT', badan: isi },
  )
}

/** `Delete` b6160 — `deleteRow` b6170 dan `save` b6191: menetap seketika. */
export async function hapusDiagnosa(
  klaimID: string,
  pesertaID: string,
  diagID: number,
): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/diagnosa/${diagID}`,
    { metode: 'DELETE' },
  )
}

// ---------------------------------------------------------------------------
// Dokumen pendukung — butir be: unggah / unduh / hapus lewat outbox.
// ---------------------------------------------------------------------------

/**
 * `Add attachment` b1245 → `AttachDocumentLife` b1273.
 *
 * ⛔ `multipart/form-data`, dan `Content-Type` SENGAJA tidak disetel tangan:
 * batas multipart-nya dibangkitkan browser, dan header yang ditulis sendiri
 * akan menyebut batas yang salah. Karena itu jalur ini tidak lewat `minta`,
 * yang selalu memasang `application/json`.
 *
 * ⛔ Batas ukurannya ditegakkan BACKEND saat menyalin (25 MiB). Yang di sini
 * hanya supaya permintaannya masuk akal sejak berangkat — layar bukan
 * penjaga.
 */
export async function unggahDokumen(
  klaimID: string,
  pesertaID: string,
  berkas: File,
  kategori: string,
): Promise<Dokumen> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  isi.append('kategori', kategori)
  return mintaFormulir<Dokumen>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/dokumen`,
    isi,
  )
}

/**
 * `Delete` b4288 → `ConfirmDeleteAttachment` b4317.
 *
 * ⚠️ Barisnya hilang SEKETIKA; penghapusan di penyimpanan menyusul lewat
 * outbox — urutan itu dari rule: `DeleteDocument_Act` b513 `Obj-Delete`
 * berjalan tanpa prasyarat, hanya panggilan penyimpanannya yang bersyarat
 * (b472).
 */
export async function hapusDokumen(klaimID: string, dokID: string): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/dokumen/` +
      encodeURIComponent(dokID),
    { metode: 'DELETE' },
  )
}

/**
 * Tautan unduh satu dokumen — `View Office Online` b3502 dan tautan baris.
 *
 * ⛔ Jalurnya TIDAK menyebut klaim, dan itu meniru aslinya: `URLPUBLIC`
 * dicari dengan `imageid` saja (`GetLinkStorage_SQL.xml` b91). Batas klaimnya
 * tetap ditegakkan backend.
 */
export function tautanDokumen(dokID: string): string {
  return rakitURL(`/api/dokumen/${encodeURIComponent(dokID)}/isi`)
}

/**
 * Apakah dokumen ini sudah benar-benar terunggah.
 *
 * ⛔ `tStorageId` kosong berarti efek outbox-nya BELUM selesai — bukan
 * berkasnya hilang. Layar membedakan keduanya supaya pemakai tahu menunggu,
 * bukan mengunggah ulang.
 */
export function dokumenTerunggah(d: Dokumen): boolean {
  return d.tStorageId.trim() !== ''
}


/**
 * Satu permintaan `multipart/form-data`.
 *
 * ⛔ Terpisah dari `minta` hanya pada SATU hal: ia tidak memasang
 * `Content-Type`. Batas multipart-nya dibangkitkan browser, dan header yang
 * ditulis tangan akan menyebut batas yang salah — permintaan lalu ditolak
 * server dengan galat yang tidak menyebutkan sebabnya.
 *
 * ⚠️ Amplop galatnya SAMA (`ApiFailure`, kunci `galat`). Dua model
 * galat berdampingan berarti dua jalan menampilkan kegagalan yang sama, dan
 * yang satu akan diam-diam kalah — pelajaran yang sudah tertulis di kepala
 * bagian ini.
 */
async function mintaFormulir<T>(jalur: string, isi: FormData): Promise<T> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  let jawab: Response
  try {
    jawab = await fetch(rakitURL(jalur), {
      method: 'POST',
      headers: { ...headerIdentitas() },
      body: isi,
      signal: kendali.signal,
    })
  } finally {
    clearTimeout(jam)
  }
  const teks = await jawab.text()
  let hasil: unknown
  if (teks.trim() !== '') {
    try {
      hasil = JSON.parse(teks)
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
    // ⛔ Kuncinya `galat`, DIBACA DENGAN CARA YANG SAMA seperti `minta`.
    // Ronde pertama fungsi ini meneruskan badan mentahnya sebagai
    // `IsiGalatApi` - dan setiap pesan backend pada jalur unggah akan jatuh
    // ke teks bawaan, persis cacat `galat` vs `error` yang pertama.
    const o = (hasil ?? {}) as { galat?: unknown }
    throw new ApiFailure(jawab.status, {
      code: 'DITOLAK_BACKEND',
      message: typeof o.galat === 'string' && o.galat !== '' ? o.galat : undefined,
    })
  }
  return hasil as T
}