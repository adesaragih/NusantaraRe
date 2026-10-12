// modul/claimlife/frontend/api.ts - panggilan backend modul Claim Life, satu fungsi per
// endpoint. Klien HTTP-nya `inti/klien.ts` (refactor bentuk B, 30-09-2026:
// dipecah dari `services/api.ts` tanpa mengubah satu panggilan pun).

import {
  ApiFailure,
  BATAS_WAKTU_MS,
  kegagalanDari,
  minta,
  mintaFormulir,
  type PenghalangTutup,
  rakitURL,
} from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'
import { TAHAP } from './labels'

/**
 * Nama tahap sebagai TIPE, bukan string mentah (GILIRAN-11 paket 4) — nilainya
 * VERBATIM `pyTaskName` `Register_Flow`, satu sumber di `TAHAP`.
 */
export type NamaTahap = (typeof TAHAP)[keyof typeof TAHAP]

/** Penjaga tipe: teks tahap dari backend (boleh kosong atau asing) -> NamaTahap. */
export function apakahNamaTahap(t: string): t is NamaTahap {
  return (Object.values(TAHAP) as string[]).includes(t)
}

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
  /** Medan dialog AdjustmentDetail (11-10-2026). persenKlaim dan klaimDibayar DITURUNKAN server dari Claim RNM / Share RNM. */
  shareNusantaraRe?: Uang
  cedingRetention?: Uang
  sumInsured?: Uang
  sumReasured?: Uang
  shareRetro?: Uang
  retrocededShare?: Uang
  persenKlaim?: string
  klaimDibayar?: Uang
  /** OutStanding Date — diturunkan server dari jejak transisi '' → '0' (Save to RNM); kosong sebelumnya. */
  tanggalOutstanding?: string
  /** Data bank (Payable To baca saja; tiga lainnya isian teks sebelum Save to RNM, 11-10-2026). */
  payableTo?: string
  namaBank?: string
  cabangBank?: string
  nomorRekening?: string
}

/** Rasio (persen/rate) dari server: teks desimal, tidak pernah Number. */
export interface Rasio {
  value: string
  scale?: number
}

export interface SpreadingBaris {
  lapisan: {
    treatyTypeName: string
    treatyYearLife: string
    retrocededShare: Uang
    reinsurer: {
      nama: string
      pctShare: Rasio
      amount: Uang
      rate: Rasio
      komisi: Rasio
      ovrComm: Rasio
      gross: Uang
      net: Uang
    }[]
  }[]
  rekap: { reinsurer: string; currency: string; gross: Uang; net: Uang }[]
  /** Baris TOTAL per mata uang di bawah Spreading Summary (11-10-2026), dijumlahkan server. */
  total?: { currency: string; gross: Uang; net: Uang }[]
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
 * diminta saat UNDUH - `DownloadDocumentClaim.xml` langkah 2 b369
 * `Call GetUrlGoogleStorage_Act` (langkah pemuat 1.4.1 b1011 ter-remark,
 * sensus 28-09-2026); penyambungan
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

/**
 * Satu peserta yang diklaim, beserta baris-barisnya sendiri.
 *
 * ⛔ Tiga tanggal klaim dialog Edit Date datang dari `TanggalKlaim` — satu
 * tipe, bukan tiga medan lepas (GILIRAN-11 paket 4). Saat DIBACA isinya
 * `YYYY-MM-DD HH24:MI:SS` (`fmtTanggalOracle`); kotak tanggal memakai
 * `keIsianTanggal`.
 */
export interface Peserta extends TanggalKlaim {
  id: string
  nomorPremiList: string
  /** IDPEGA baris sumber (diturunkan server lewat SOURCE_ID) - kunci data polis Outstanding (10-10-2026). */
  idPega?: string
  nomorPolis: string
  nomorSertifikat: string
  mataUang: string
  /** Penanda DIPILIH untuk diklaim (`IS_CHECK`). Teks, bukan boolean. */
  isCheck: string
  /**
   * Checkbox Accept grid Participant Details layar Outstanding (`.IsAccept` InputOSClaimLife b17582, IS_ACCEPT migrasi
   * 026): 'true' / 'false' / ''. ⛔ Bukan isCheck. Dibaca gerbang dokumen Save to RNM.
   */
  isAccept?: string
  /** NAME OF INSURED / DATE OF BIRTH (YYYY-MM-DD) grid Participant Details - dibaca server dari baris sumber, tidak disimpan. */
  namaTertanggung?: string
  tanggalLahir?: string
  /** Rincian peserta layar Outstanding (ClaimLifeDetailGCNM, screenshot Pega 11-10-2026). Tanggal teks apa adanya. */
  jenisKelamin?: string
  plan?: string
  umur?: string
  tanggalMulai?: string
  tanggalEfektif?: string
  tanggalExpired?: string
  valuasiMulai?: string
  valuasiSelesai?: string
  /** Jendela RETRO (TR/TP) — validasi DOL langsung di layar (11-10-2026). */
  valuasiRetroMulai?: string
  valuasiRetroSelesai?: string
  tanggalTerima?: string
  stnc?: string
  maksKlaimKadaluarsa?: string
  /** ADMIN NOTES (migrasi 027) dan CLAIM ANALYST (NOTES). */
  catatanAdmin?: string
  catatanAnalis?: string
  /** RECOMMENDATION (migrasi 003) — kotak MEDICAL ADVISOR, diisi Medical Advisor di tahap Medical Check (11-10-2026). */
  rekomendasi?: string
  /** MAX DATA RECEIVED (migrasi 028), diisi pemakai sebelum Save to RNM. */
  tanggalMaksData?: string
  /**
   * `DATE_OF_LOSS` — tanggal kejadian, milik PESERTA dan bukan klaim.
   *
   * Teks apa adanya, bentuk `YYYY-MM-DD HH24:MI:SS` (`fmtTanggalOracle`).
   * Kosong berarti belum diisi.
   */
  tanggalKejadian: string
  /**
   * Butir bk — `.MAXCLAIM_RECEIVED` (`MAX CLAIM RECEIVED` b12131), dihitung
   * backend saat baca. Kosong berarti sah — KECUALI `penandaTerimaKlaimAlasan`
   * terisi: di situ penandanya tidak dapat dihitung, dan alasannya itulah.
   */
  penandaTerimaKlaim: string
  penandaTerimaKlaimAlasan: string
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
  /** `T_GENERAL_CLAIM.TYPE` (QP/QR/TP/TR) - tombol Send to Medical Check hanya QP/QR (InputOSClaimLife b22123). */
  type?: string
  /** Tab Register Claim layar Outstanding (10-10-2026): Catastrophe b3504, Is KPR b3810, STNC Treaty b4757, Note b4934. */
  stsKatastrofe?: string
  isKpr?: string
  stncClaim?: string
  katastrofeNote?: string
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

/**
 * Ambil satu klaim Life beserta SELURUH baris adjustment-nya — tiket 01 AC-1,
 * bukan hanya baris terakhir. Di backend: GET /api/klaim-life/{id},
 * internal/handlers/klaimlife.go.
 */
export async function ambilKlaimLife(id: string): Promise<Klaim> {
  const data = await minta<Klaim>(`/api/klaim-life/${encodeURIComponent(id)}`)
  return data
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
  /** M_LIFE_PREMIUM_DETAIL.IDPEGA - kunci header data polis (10-10-2026). */
  idPega: string
  /** DOB teks ISO YYYY-MM-DD ("" = kosong) - kolom Date of Birth layar Register (10-10-2026). */
  tanggalLahir: string
}

/** Pilihan peserta satu premium list (10-10-2026). Mata uang hanya keterangan: server memakai mata uang sumber. */
export interface KelompokDaftar {
  nomorPremiList: string
  nomorPolis: string
  mataUang: string
  sertifikat: string[]
}

/**
 * Isi satu pendaftaran klaim - BANYAK premium list sekaligus (keputusan work owner 10-10-2026): server membuat satu
 * case per PL number + mata uang, semua atau tidak sama sekali. Type dan kode bisnis TIDAK dikirim: server membacanya
 * dari polis PL tiap case.
 */
export interface PermintaanDaftar {
  kelompok: KelompokDaftar[]
}

/** Satu case yang terbentuk. */
export interface KasusTerdaftar {
  id: string
  nomorKlaim: string
  nomorPremiList: string
  mataUang: string
  jumlahPeserta: number
}

/** Jawaban pendaftaran yang berhasil: daftar case (satu per PL + mata uang). */
export interface HasilDaftar {
  kasus: KasusTerdaftar[]
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
export async function tolakBarisAdjustment(
  klaimID: string,
  adjID: string,
  komentar: string,
): Promise<void> {
  // OQ-M5 (GILIRAN-17): `Remarks` dialog Reject Outstanding (b1687, wajib)
  // disimpan backend di `T_CLAIMLF_JEJAK.KOMENTAR`; kosong dijawab 400.
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/adjustment/${encodeURIComponent(adjID)}/tolak`,
    { metode: 'POST', badan: { komentar } },
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
 * `Add` grid adjustment — putaran adjustment berikutnya bagi seorang peserta.
 *
 * Inilah yang membuat klaim TIDAK TERMINAL: penolakan Komite menghasilkan
 * putaran berikutnya, bukan akhir (ADR-U-0011). Yang terminal adalah BARIS.
 *
 * Jawaban yang mungkin:
 *   403 bukan ReasLifeSPV
 *   409 baris terakhir belum ditolak (atau peserta tanpa baris — klaim lama);
 *       kasus tertutup
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
 * Apakah layar menawarkan `Add` b17937 — tombol grid `.AdjustmentList`.
 *
 * ⛔ RALAT 29-09-2026 (GILIRAN-14 butir bp, meralat bo): `Add` adalah jalur
 * PUTARAN saja. Baris pertama lahir saat Submit Register (`SavePesertaClaim`
 * 7.8), jadi grid peserta terpilih tidak pernah kosong — GILIRAN-13 sempat
 * menjadikan `Add` pada grid kosong jalan lahir baris pertama.
 *
 * ⛔ Tiga syarat:
 *
 *   1. tahap Claim Analis — syarat tampil `pyWorkPage.pyPosition
 *      =='ReasLifeSPV'` b18160;
 *   2. baris terakhir DITOLAK (`bolehPutaranBaru`) — di keadaan lain
 *      layanan menjawab 409, dan tombol yang selalu ditolak bukan tombol.
 *      ⚠️ PENYIMPANGAN SADAR: b18160 menampilkan `Add` tanpa syarat baris;
 *      syarat ini milik layanan (`BarisLanjutan`), bukan XML;
 *   3. kasus belum tertutup (butir bb).
 *
 * ⚠️ Kenyamanan, bukan pagar: services menolak dengan 403/409.
 */
export function bolehAddAdjustment(klaim: Klaim, peserta: Peserta): boolean {
  return (
    !kasusTertutup(klaim) && klaim.tahap === TAHAP.claimAnalis && bolehPutaranBaru(peserta.baris)
  )
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
  remarks: string,
): Promise<void> {
  // 12-10-2026 (opsi A): Remarks pop-up WAJIB, disimpan server di KOMENTAR jejak penyerahan.
  await minta<void>(jalurKomite(klaimID, pesertaID, adjID), { metode: 'POST', badan: { remarks } })
}

/** Jalur penyerahan / pratinjau komite satu baris adjustment. */
function jalurKomite(klaimID: string, pesertaID: string, adjID: string): string {
  return (
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
    `/peserta/${encodeURIComponent(pesertaID)}` +
    `/adjustment/${encodeURIComponent(adjID)}/komite`
  )
}

/** Isi pop-up Send Claim to Committee sebelum dikirim (GET; gerbang Serahkan tanpa menulis). */
export interface PratinjauKomite {
  /** Waktu server ISO — Date, baca saja. */
  waktu: string
  pic: string
  /** Claim RNM baris (teks desimal) dan mata uangnya. */
  jumlahKlaim: string
  mataUang: string
  anggota: { namaKomite: string; urut: number }[]
}

export async function ambilPratinjauKomite(klaimID: string, pesertaID: string, adjID: string): Promise<PratinjauKomite> {
  return minta<PratinjauKomite>(jalurKomite(klaimID, pesertaID, adjID))
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
 * ⛔ Saat ini SELALU menjawab 405 (+ header `Allow`): ADR-U-0031 menetapkan
 * penghapusan berupa PENANDA, bukan hapus fisik, dan kolom penandanya belum
 * diputuskan. (Komentar lama menyebut 501 - layar pun menunggu 501, dan
 * menampilkan "Gagal menghapus klaim."; diralat GILIRAN-11.)
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
  /** "PL Number" (11-10-2026, pengganti Policy No) — PL peserta aktif klaim; kosong untuk klaim tanpa peserta. */
  plNumber: string
  /** "Update Operator Name" — akun pelaku jejak terakhir, cadangan pembuat (T_WORK_CLAIM tak punya kolomnya). */
  updateOpName: string
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

/** Panjang kata kunci kotak cari inbox — sama dengan repository.MinCariInbox / MaksCariInbox (12-10-2026). */
export const MIN_CARI_INBOX = 2
export const MAKS_CARI_INBOX = 100

/** Kata kunci yang dikirim: dipangkas; kurang dari 2 karakter = tanpa saringan (layar tidak mengirim yang pasti ditolak). */
export function kataCariInbox(kata: string): string {
  const k = kata.trim()
  return k.length >= MIN_CARI_INBOX ? k.slice(0, MAKS_CARI_INBOX) : ''
}

/**
 * Membaca satu halaman antrian sebuah tahap. `cari` (12-10-2026): Case ID, Claim No, PL Number, atau nama tertanggung —
 * disaring di SERVER (inbox berpaginasi); kosong = tanpa saringan.
 */
export async function ambilKotakMasuk(
  tahap: NomorTahap,
  halaman = 1,
  ukuran?: number,
  cari = '',
): Promise<HalamanInbox> {
  return await minta<HalamanInbox>('/api/klaim-life', {
    kueri: cari === '' ? { tahap, halaman, ukuran } : { tahap, halaman, ukuran, cari },
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
  // 11-10-2026: lewat fetch sendiri (pola simpanKeRNM) supaya 422 berdaftar (peserta Accept tanpa dokumen saat Send to
  // Medical Check) membawa `kekurangan` — dilempar sebagai GagalSimpanRNM (penolakan berkekurangan yang sama bentuknya).
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jalur = `/api/klaim-life/${encodeURIComponent(klaimID)}/tahap/${tujuan}`
    const jawab = await fetch(rakitURL(jalur), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...headerIdentitas() },
      signal: kendali.signal,
    })
    const teks = await jawab.text()
    if (!jawab.ok) {
      const dasar = kegagalanDari(jawab.status, teks)
      const kurang = jawab.status === 422 ? kekuranganDari(teks) : []
      throw kurang.length > 0 ? new GagalSimpanRNM(dasar, kurang) : dasar
    }
  } finally {
    clearTimeout(jam)
  }
}

/** Hasil `Save to RNM` yang berhasil — `services.HasilSimpanRNM`. */
export interface HasilSimpanRNM {
  nomorKlaim: string
  /** Nomor diterbitkan saat simpan ini (langkah 16-20), bukan saat pendaftaran. */
  nomorBaru: boolean
  /** Baris tanpa status yang kini Outstanding (langkah 22.1.3.2). */
  barisDitandai: number
  /**
   * Keadaan efek Arasapas langkah 28, sebagai kata: `terkirim`, `gagal, …`, atau
   * `dilewati: …` (lingkungan bukan produksi, atau kode retro langkah 27 — dibaca
   * sesudah tukar dua syarat, OQ-N5 ditutup GILIRAN-17).
   */
  arasapas: string
}

/**
 * `Save to RNM` — tombol layar Outstanding (`InputOSClaimLife.xml` b21102 →
 * `SaveOutStandingLife_Act`). Seluruh gerbang XML diperiksa di server.
 *
 * Kode jawaban yang mungkin:
 *   200 tersimpan
 *   400 permintaan tidak sah (mis. peserta tanpa SOURCE_ID)
 *   401 tanpa identitas · 403 bukan pemegang tahap Outstanding
 *   409 bukan tahap Outstanding Claim, atau kasusnya sudah ditutup
 *   422 gerbang XML menolak — `galat` membawa kalimat XML apa adanya,
 *       `langkah` langkah asalnya; atau data untuk memeriksanya belum lengkap.
 *       Sejak 11-10-2026 kolom wajib yang kosong dikumpulkan SEMUA di
 *       `kekurangan` (dilempar sebagai GagalSimpanRNM)
 *   503 basis data belum dikonfigurasi
 *
 * ⚠️ Lewat `fetch` sendiri (pola ambilIsiDokumen), bukan `minta`: amplop
 * `ApiFailure` inti hanya membawa `galat`/`penghalang`, dan `kekurangan`
 * milik modul ini — inti tidak diubah.
 */
export async function simpanKeRNM(klaimID: string): Promise<HasilSimpanRNM> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`/api/klaim-life/${encodeURIComponent(klaimID)}/outstanding`), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...headerIdentitas() },
      signal: kendali.signal,
    })
    const teks = await jawab.text()
    if (!jawab.ok) {
      const dasar = kegagalanDari(jawab.status, teks)
      const kurang = jawab.status === 422 ? kekuranganDari(teks) : []
      throw kurang.length > 0 ? new GagalSimpanRNM(dasar, kurang) : dasar
    }
    return JSON.parse(teks) as HasilSimpanRNM
  } finally {
    clearTimeout(jam)
  }
}

/** Satu kolom wajib Save to RNM yang belum terisi — `services.KekuranganRNM`. nomorPeserta 0 = tingkat klaim (Is KPR). */
export interface KekuranganRNM {
  nomorPeserta: number
  pesertaId: string
  sertifikat: string
  nama: string
  kolom: string
  pesan: string
}

/** Penolakan Save to RNM yang membawa SELURUH kekurangan kolom wajib (422, keputusan work owner 11-10-2026). */
export class GagalSimpanRNM extends ApiFailure {
  constructor(
    dasar: ApiFailure,
    readonly kekurangan: KekuranganRNM[],
  ) {
    super(dasar.status, dasar.detail)
    this.name = 'GagalSimpanRNM'
  }
}

/** Daftar `kekurangan` dari badan 422 Save to RNM; bentuk asing = []. */
export function kekuranganDari(teks: string): KekuranganRNM[] {
  try {
    const o = JSON.parse(teks) as { kekurangan?: unknown }
    return Array.isArray(o.kekurangan) ? (o.kekurangan as KekuranganRNM[]) : []
  } catch {
    return []
  }
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
 *   401 tanpa identitas
 *   403 bukan Admin (butir bj, 28-09-2026)
 *   409 bukan tahap Outstanding Claim, atau kasusnya sudah ditutup
 *   422 DOL di luar jendela valuasi — badannya membawa `pesertaId`; atau
 *       tahap kasus tidak dikenal
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

/**
 * Tiga tanggal klaim dialog Edit Date selain DOL — `CLAIM_RECEIVED_DATE`
 * b1076, `COMPLETE_DATE` b1387, `CONFIRMATION_DATE` b1626. Dikirim `Save`
 * sebagai `YYYY-MM-DD` (kosong = dikosongkan); dibaca lewat `Peserta`.
 */
export interface TanggalKlaim {
  tanggalTerimaKlaim: string
  tanggalDokumenLengkap: string
  tanggalKonfirmasi: string
}

/**
 * Menyimpan tiga tanggal klaim satu peserta SEKALIGUS.
 *
 * Tombol `Save` `EditDateClaimLife_Section.xml` b1910 →
 * `UpdateDateClaimLife_Act` b1929. Satu permintaan, sebab satu tombol.
 *
 * Kode jawaban yang mungkin:
 *   204 berhasil
 *   400 salah satu tanggal bukan tanggal yang dikenal
 *   401 tanpa identitas pelaku
 *   403 bukan Admin
 *   409 bukan tahap Outstanding Claim, atau kasusnya sudah ditutup
 *   422 tahap kasus tidak dikenal
 *   503 basis data belum dikonfigurasi
 */
export async function ubahTanggalKlaim(
  klaimID: string,
  pesertaID: string,
  tanggal: TanggalKlaim,
): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}` +
      `/peserta/${encodeURIComponent(pesertaID)}/tanggal-klaim`,
    { metode: 'PUT', badan: tanggal },
  )
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
export const TAHAP_PENAWAR_TUTUP: readonly NamaTahap[] = [TAHAP.outstandingClaim, TAHAP.claimAnalis]

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
  const t = klaim.tahap
  return apakahNamaTahap(t) && TAHAP_PENAWAR_TUTUP.includes(t)
}

/** Apakah kasus ini sudah ditutup dan karena itu tidak dapat diubah lagi. */
export function kasusTertutup(klaim: Klaim | null): boolean {
  return klaim !== null && klaim.statusWork === STATUS_WORK_SELESAI
}

/**
 * Apakah klaim ini sudah pernah di-Save to RNM — padanan `repository.SudahSaveRNM`
 * (OQ-M1, GILIRAN-17): ada baris adjustment yang `STS_REJECT`-nya terisi.
 * Baris lahir kosong; hanya Save to RNM (atau jalur yang menuntut '0') mengisinya.
 */
export function sudahSaveRNM(klaim: Klaim): boolean {
  return klaim.peserta.some((p) => p.baris.some((b) => b.kodeStatus !== ''))
}

/**
 * Apakah tombol `DELETE` peserta (`InputOSClaimLife` b17865) ditawarkan —
 * padanan `models.BolehCabutPeserta` (OQ-M6, GILIRAN-17): layar Outstanding,
 * dan `CLAIM_NO == ''` (b18082) = belum pernah Save to RNM.
 *
 * ⚠️ Kenyamanan tampilan saja; peran Admin dan gerbangnya ditegakkan backend.
 */
export function bolehCabutPeserta(klaim: Klaim | null): boolean {
  return (
    klaim !== null &&
    !kasusTertutup(klaim) &&
    klaim.tahap === TAHAP.outstandingClaim &&
    !sudahSaveRNM(klaim)
  )
}

/**
 * Mencabut seorang peserta dari klaimnya — PENANDA `STS_HAPUS` (migrasi 022),
 * bukan hapus baris (ADR-U-0031). Tanpa konfirmasi, seperti b18021.
 */
export async function cabutPeserta(klaimID: string, pesertaID: string): Promise<void> {
  await minta<void>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/peserta/${encodeURIComponent(pesertaID)}/cabut`,
    { metode: 'POST' },
  )
}

/**
 * Apakah tiga tanggal klaim dialog Edit Date terbuka untuk diubah.
 *
 * Padanan `models.BolehUbahTanggalKlaim`: isiannya baca-saja bila
 * `pyPosition!='ReasLifeAdmin' || CLAIM_NO!=''` (b1000, b1313, b1550, b1788).
 * Tombol `Edit Date` berdiri di grid peserta — irisannya hanya Outstanding
 * Claim — dan sesudah Save to RNM pertama tanggalnya terkunci (OQ-M1).
 *
 * ⚠️ Ini HANYA menentukan terbuka atau tidaknya kotak. Gerbang sebenarnya
 * (termasuk peran Admin) ada di backend dan diperiksa lagi di sana.
 */
export function bolehUbahTanggalKlaim(klaim: Klaim | null): boolean {
  return (
    klaim !== null &&
    !kasusTertutup(klaim) &&
    klaim.tahap === TAHAP.outstandingClaim &&
    !sudahSaveRNM(klaim)
  )
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
/**
 * Daftar pilihan GROUP DIAGNOSE (keputusan work owner 12-10-2026) — satu sumber di backend (models.KelompokDiagnosa),
 * urutan WO. Layar tidak menyalin daftarnya.
 */
export async function ambilKelompokDiagnosa(): Promise<string[]> {
  const h = await minta<{ daftar: string[] }>('/api/klaim-life/kelompok-diagnosa')
  return h.daftar
}

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
 * "View" berbatas klaim (11-10-2026): GET /api/klaim-life/{id}/dokumen/{dokId}/isi — dokumen milik klaim lain ditolak 404;
 * selama penyimpanan belum berjalan berkas dilayani dari salinan lokal server. Lewat fetch (identitas ikut), pola
 * ambilIsiDokumen.
 */
export async function ambilIsiDokumenKlaim(klaimID: string, dokID: string): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jalur = `/api/klaim-life/${encodeURIComponent(klaimID)}/dokumen/${encodeURIComponent(dokID)}/isi`
    const jawab = await fetch(rakitURL(jalur), {
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

/**
 * Daftar kategori dropdown "Category" (11-10-2026) = Document Claim produk polis klaim (Master Product Name Life), urutan
 * master — sumber yang sama dengan pemeriksa unggahan.
 */
export async function ambilKategoriDokumen(klaimID: string): Promise<string[]> {
  const data = await minta<{ kategori: string[] | null }>(`/api/klaim-life/${encodeURIComponent(klaimID)}/dokumen/kategori`)
  return data.kategori ?? []
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
 * Isi satu dokumen sebagai Blob — `View Office Online` b3502
 * (runActivity `DownloadDocumentClaim` b3519).
 *
 * ⛔ LEWAT `fetch`, bukan pranala (GILIRAN-12 paket 0). Ronde sebelumnya
 * mengembalikan URL untuk `<a href download>`; browser yang mengikuti pranala
 * tidak membawa `X-Pelaku`/`X-Peran`, dan backend membaca identitas HANYA
 * dari header itu - setiap unduhan dijawab 401.
 *
 * ⛔ Jalurnya TIDAK menyebut klaim, dan itu meniru aslinya: `URLPUBLIC`
 * dicari dengan `imageid` saja (`GetLinkStorage_SQL.xml` b91). Batas klaimnya
 * tetap ditegakkan backend.
 */
export async function ambilIsiDokumen(dokID: string): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`/api/dokumen/${encodeURIComponent(dokID)}/isi`), {
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

// ——— Butir pl4/av: PolicyDataLife dari PremiumList Life ———

/**
 * Data polis yang mengisi `.PolicyDataLife.*` di layar Claim Life.
 *
 * ⛔ SUMBERNYA TABEL RELASIONAL modul PremiumList Life, bukan cermin JSON-nya
 * — `[keputusan work owner, butir av]`.
 *
 * ⚠️ Nama kuncinya sengaja sama dengan nama properti Pega, supaya siapa pun
 * yang membandingkan layar dengan rule tidak perlu menerjemahkan.
 */
export interface PolicyDataLife {
  nomorPolis: string
  type: string
  marketingName: string
  cedingCoName: string
  policyHolderName: string
  businessName: string
  /** `YYYY-MM-DD`; kosong berarti kolomnya belum diisi. */
  dateReceived: string
  status: string
  statusUpdate: string
  /** Kunci ambang batas hari (butir ba). */
  productNameId: string
  productName: string
  /**
   * Versi polis yang terbaca — `PROD_KE` TERBESAR.
   *
   * ⛔ Satu nomor polis punya banyak versi. Layar yang tidak dapat menyebut
   * versi mana yang ditampilkannya membuat selisih angka mustahil ditelusuri.
   */
  prodKe: number
  /** Medan layar lama yang belum punya kolom di mana pun. */
  medanTanpaSumber: string[]
  /** av-2 — `.TypeCeding` kode; `typeCedingName` katanya (yang ditampilkan). */
  typeCeding: string
  typeCedingName: string
  /** av-2 — `.ProRateType` "Premium Method". */
  proRateType: string
  /** av-2 — `.WPC`, `YYYY-MM-DD`; kosong berarti belum diisi. */
  wpc: string
  /** av-2 — tampil HANYA bagi Type TP/TR (`pyCondition` b10541/b10824). */
  retroName: string
  securityReinsurer: string
  /** av-2 — `.SobName` "SOB". */
  sobName: string
}

/**
 * Membaca data polis versi berjalan — `GET /api/polis-life/ringkas`.
 *
 * ⚠️ Menjawab **404** bila nomornya tidak ada di PremiumList Life. Itu bukan
 * kerusakan: klaim dapat didaftarkan atas polis yang belum ada di modul itu,
 * dan layar menyatakannya alih-alih pecah.
 */
export async function ambilDataPolis(nomorPolis: string): Promise<PolicyDataLife> {
  return minta<PolicyDataLife>(
    `/api/polis-life/ringkas?nomorPolis=${encodeURIComponent(nomorPolis)}`,
  )
}

/**
 * Membaca header polis dari IDPEGA peserta — `GET /api/polis-life/ringkas?idpega=` (keputusan work owner 10-10-2026:
 * "ambil ke tabel T_PREMIUM_LIST berdasarkan idpega peserta nya saja"). 404 = header tidak ditemukan.
 */
export async function ambilDataPolisIdPega(idPega: string): Promise<PolicyDataLife> {
  return minta<PolicyDataLife>(`/api/polis-life/ringkas?idpega=${encodeURIComponent(idPega)}`)
}

// ---------------------------------------------------------------------------
// LAYAR OUTSTANDING mengikuti Pega (10-10-2026) - `Section/InputOSClaimLife.xml` (korpus kelvin/InboxClaimLife (Done)).
// ---------------------------------------------------------------------------

/** Nilai radio Catastrophe - `[keputusan work owner]` claimprop grilling ronde 2 P12; bawaan Non-Catastrophe b3728. */
export const PILIHAN_KATASTROFE = ['Catastrophe', 'Non-Catastrophe'] as const
/** Nilai radio Is KPR - keputusan work owner 11-10-2026: teks disimpan persis (pola Catastrophe / Non-Catastrophe). */
export const PILIHAN_KPR = ['KPR', 'Non-KPR'] as const

/** Isi tombol `Save` layar Outstanding (b19413): medan tab Register Claim + Accept per peserta. */
export interface IsianOutstanding {
  stsKatastrofe: string
  isKpr: string
  katastrofeNote: string
  /** Pengenal peserta -> 'true' / 'false'. */
  accept: Record<string, string>
  /** Pengenal peserta -> ADMIN NOTES (migrasi 027, 11-10-2026). */
  catatanAdmin?: Record<string, string>
  /** Pengenal peserta -> lima tanggal yang diisi pemakai sebelum Save to RNM (11-10-2026). */
  tanggal?: Record<string, IsianTanggal>
  /** Pengenal BARIS adjustment -> Claim RNM + data bank (11-10-2026: dilebur ke tombol Save grid). */
  adjustment?: Record<string, IsianBarisAdjustment>
}

/** Lima tanggal peserta, teks YYYY-MM-DD ('' = kosong). DOL, Claim Received, Confirmation WAJIB sebelum Save to RNM. */
export interface IsianTanggal {
  kejadian: string
  terimaKlaim: string
  dokumenLengkap: string
  maksData: string
  konfirmasi: string
}

/** Tombol `Save` layar Outstanding - PUT /api/klaim-life/{id}/outstanding/isian. */
export async function simpanIsianOutstanding(klaimID: string, isian: IsianOutstanding): Promise<void> {
  await minta<void>(`/api/klaim-life/${encodeURIComponent(klaimID)}/outstanding/isian`, { metode: 'PUT', badan: isian })
}

/** Isian tombol Save layar Medical Check — Recommendation per pengenal peserta (11-10-2026). */
export interface IsianMedis {
  rekomendasi: Record<string, string>
}

/** Batas RECOMMENDATION VARCHAR2(1000) — server memeriksa byte, layar membatasi karakter. */
export const BATAS_REKOMENDASI = 1000

/** Tombol Save layar Medical Check (`MedicalCheckClaimLife`). Gerbang tahap / Medical Advisor / STS_REJECT di server. */
export async function simpanMedis(klaimID: string, isian: IsianMedis): Promise<void> {
  await minta<void>(`/api/klaim-life/${encodeURIComponent(klaimID)}/medis`, { metode: 'PUT', badan: isian })
}

/** Batas CLAIM ANALYST (NOTES VARCHAR2(1000)) — server memeriksa byte. */
export const BATAS_CATATAN_ANALIS = 1000

/** Tombol Save layar Claim Analis (12-10-2026): catatan CLAIM ANALYST per peserta. Gerbang tahap / SPV di server. */
export async function simpanAnalis(klaimID: string, isian: { catatan: Record<string, string> }): Promise<void> {
  await minta<void>(`/api/klaim-life/${encodeURIComponent(klaimID)}/analis`, { metode: 'PUT', badan: isian })
}

/** Kasus ada di tahap Outstanding dan belum ditutup - isian dan tombol layar Outstanding hidup. */
export function bolehIsianOutstanding(klaim: Klaim | null): boolean {
  return klaim !== null && !kasusTertutup(klaim) && klaim.tahap === TAHAP.outstandingClaim
}

/** Klaim sudah bernomor (`CLAIM_NO != ''`) - nomornya lahir di Save to RNM. */
export function bernomor(klaim: Klaim): boolean {
  return klaim.nomorKlaim.trim() !== ''
}

/**
 * Tombol bawah layar Outstanding - kondisi Pega, ditegakkan juga di server (409):
 *   Save to RNM b21102: nonaktif sesudah tersimpan (`pyWorkPage.Save = 1` b21095; padanannya sudahSaveRNM, OQ-N1).
 *   Send Back to Register b21404: tampil hanya bila CLAIM_NO kosong (b21699).
 *   Send to Medical Check b21839: tampil bila Type QP/QR (b22123), nonaktif bila CLAIM_NO kosong (b21832).
 */
export function aturanTombolOutstanding(klaim: Klaim): {
  simpanRNMAktif: boolean
  kirimBalikTampil: boolean
  kirimMedisTampil: boolean
  kirimMedisAktif: boolean
} {
  const tipe = (klaim.type ?? '').trim().toUpperCase()
  return {
    simpanRNMAktif: !sudahSaveRNM(klaim),
    kirimBalikTampil: !bernomor(klaim),
    kirimMedisTampil: tipe === 'QP' || tipe === 'QR',
    kirimMedisAktif: bernomor(klaim),
  }
}

/** Tombol DELETE grid Participant Details: tampil hanya bila CLAIM_NO kosong (b18082) dan cabut diizinkan. */
export function bolehDeletePeserta(klaim: Klaim | null): boolean {
  return klaim !== null && !bernomor(klaim) && bolehCabutPeserta(klaim)
}

/** Satu baris grid View Detail peserta POLIS (DetailPolisLife b4537; 11-10-2026). Tanggal teks YYYY-MM-DD. */
export interface BarisPesertaPolis {
  urut: number
  nomorPolis: string
  nomorSertifikat: string
  namaTertanggung: string
  jenisKelamin: string
  tanggalLahir: string
  umur: string
  plan: string
  tanggalMulai: string
  tanggalEfektif: string
  tanggalExpired: string
}

/** Satu halaman View Detail - 500 baris per halaman (LoadDataPeserta_Act b572). */
export interface HalamanPesertaPolis {
  baris: BarisPesertaPolis[]
  halaman: number
  ukuran: number
  total: number
}

/** View Detail tab Policy Details - GET /api/klaim-life/{id}/peserta-polis?halaman=N (baca saja). */
export async function ambilPesertaPolis(klaimID: string, halaman: number): Promise<HalamanPesertaPolis> {
  return minta<HalamanPesertaPolis>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/peserta-polis?halaman=${encodeURIComponent(String(halaman))}`,
  )
}

/** Isi dialog AdjustmentDetail yang dapat diisi sebelum Save to RNM (11-10-2026). */
export interface IsianBarisAdjustment {
  jumlah: string
  /** Name of Bank — SATU-SATUNYA data bank yang dapat diisi ("semua read only kecuali name of bank", 11-10-2026). */
  namaBank: string
}

/**
 * Status baris sebagai teks grid adjustment. '' (belum Save to RNM) → '—'; '0' → 'OutStanding' (screenshot Pega
 * 11-10-2026); '1' / '2' → 'Accepted' / 'Rejected' seperti status inbox. ⚠️ Daftar nilai Pega (`pyListSource associated`
 * ClaimLifeDetailGCNM b18258) tidak ada di korpus.
 */
export function kataStatusBaris(kode: string): string {
  switch (kode) {
    case '':
      return '—'
    case '0':
      return 'Outstanding'
    case '1':
      return 'Accepted'
    case '2':
      return 'Rejected'
    default:
      return kode
  }
}

/** Claim RNM dapat diisi: klaim belum bernomor, baris belum berstatus, tahap Outstanding, kasus terbuka. */
export function bolehSuntingClaimRNM(klaim: Klaim, b: BarisAdjustment): boolean {
  return bolehIsianOutstanding(klaim) && !bernomor(klaim) && b.kodeStatus === ''
}

/**
 * Spreading PremiumList Life satu peserta klaim — keputusan work owner 11-10-2026 ("ambil saja dari premium list, tidak perlu
 * perhitungan lagi"). GET /api/klaim-life/{id}/peserta/{pesertaId}/spreading-polis, baca saja; lapisan kosong = tidak ada.
 */
export async function ambilSpreadingPolis(klaimID: string, pesertaID: string): Promise<SpreadingBaris> {
  return minta<SpreadingBaris>(
    `/api/klaim-life/${encodeURIComponent(klaimID)}/peserta/${encodeURIComponent(pesertaID)}/spreading-polis`,
  )
}

// ——— Tabel komite di bawah inbox (perintah work owner 09-10-2026: menu Komite Claim Life dihapus, komite digabung ke
// menu Claim Life). Dibaca lewat rute pinjaman `GET /api/komite` (`cmd/api/rakit.go`); kasusnya dibuka di tempat
// (`onBukaModul`). Bentuk = `BarisInboxKomite` modul Komite Claim Life (disalin, bukan diimpor). ———

/** Satu kasus komite yang menunggu pelaku. Uang TEKS; status baris KATA. */
export interface BarisKomite {
  kasusId: string
  tglUpdate: string
  statusWork: string
  klaimId: string
  nomorKlaim: string
  tingkatBerjalan: number
  komiteLoop: number
  nilaiKlaim: string
  mataUang: string
  statusBaris: string
}

export interface HalamanKomite {
  baris: BarisKomite[]
  total: number
  halaman: number
  ukuran: number
}

/** Laporan harian "perlu intervensi" — `GET /api/komite/laporan-harian` (admin). */
export interface LaporanHarianKomite {
  tanggal: string
  kosong: boolean
  baris: { kasusId: string; jenis: string; keadaan: string; percobaan: number; sejak: string }[]
}

export async function daftarKomite(halaman = 1): Promise<HalamanKomite> {
  return minta<HalamanKomite>(`/api/komite?halaman=${String(halaman)}`)
}

export async function laporanHarianKomite(): Promise<LaporanHarianKomite> {
  return minta<LaporanHarianKomite>('/api/komite/laporan-harian')
}

/** Nama modul layar kasus komite (`onBukaModul`). */
export const MODUL_KOMITE = 'komiteclaimlife'
