// ============================================================================
// inti/klien.ts - KLIEN HTTP bersama: satu-satunya tempat frontend berbicara
// ke backend Go (`minta`, `mintaFormulir`, amplop galat, identitas).
//
// Refactor bentuk B (30-09-2026): dipecah dari `services/api.ts`. Panggilan
// per endpoint kini hidup di `modul/<nama>/api.ts` dan memakai berkas ini;
// yang tinggal di sini hanya yang dipakai LEBIH DARI SATU modul.
//
// Cara membaca berkas ini bila baru mengenal TypeScript:
//   - `interface` = "bentuk data". Ia hanya ada saat kode diperiksa (tsc),
//     tidak ikut dikirim ke browser. Gunanya: editor menegur lebih dulu bila
//     kita salah mengeja nama medan, misalnya `klaim.nomorKlaimm`.
//   - `export` = boleh dipakai berkas lain lewat `import`.
//   - `async` / `await` = menunggu jawaban server tanpa membekukan layar.
//   - `Promise<Klaim>` = "nanti, kalau sudah datang, isinya Klaim".
// ============================================================================

import { headerIdentitas } from './store/sesi'

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
export const BATAS_WAKTU_MS = 30_000

interface OpsiMinta {
  metode?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  /** Badan permintaan; di-JSON-kan. `undefined` berarti tanpa badan. */
  badan?: unknown
  /** Parameter kueri; nilai `undefined` dilewati. */
  kueri?: Record<string, string | number | undefined>
}

export function rakitURL(jalur: string, kueri?: OpsiMinta['kueri']): string {
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
export async function minta<T>(jalur: string, opsi: OpsiMinta = {}): Promise<T> {
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
  if (!jawab.ok) throw kegagalanDari(jawab.status, teks)
  return uraiJSON(jawab.status, teks) as T
}

/** Badan sukses sebagai JSON; kosong menjadi `undefined`, bukan JSON menjadi galat. */
function uraiJSON(status: number, teks: string): unknown {
  if (teks.trim() === '') return undefined
  try {
    return JSON.parse(teks)
  } catch {
    throw bukanJSON(status)
  }
}

/** `BACKEND_TIDAK_TERJANGKAU` — jawaban yang tidak sampai dari backend. */
function bukanJSON(status: number): ApiFailure {
  return new ApiFailure(status, {
    code: 'BACKEND_TIDAK_TERJANGKAU',
    message:
      'Jawaban dari server bukan JSON; permintaan tampaknya tidak ' +
      'sampai ke backend.',
  })
}

/**
 * Kegagalan dari jawaban yang TIDAK ok — SATU pembaca amplop `{galat}` untuk
 * `minta`, `mintaFormulir`, dan `ambilIsiDokumen` (temuan /code-review
 * GILIRAN-12: tiga salinan dalam satu berkas lolos dari penjaga per berkas).
 */
export function kegagalanDari(status: number, teks: string): ApiFailure {
  let isi: unknown
  if (teks.trim() !== '') {
    try {
      isi = JSON.parse(teks)
    } catch {
      return bukanJSON(status)
    }
  }
  const o = (isi ?? {}) as { galat?: unknown; penghalang?: unknown }
  return new ApiFailure(status, {
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

// ---------------------------------------------------------------------------
// 4. PANGGILAN KE BACKEND — satu fungsi per endpoint.
// ---------------------------------------------------------------------------

export async function cekKesehatan(): Promise<Kesehatan> {
  const data = await minta<Kesehatan>('/healthz')
  return data
}

/**
 * Modul yang dipasang backend - GET /api/modul-aktif, cmd/api/rakit.go
 * (refactor bentuk B, `MODUL_AKTIF`). Menu modul yang tidak disebut tidak
 * tampil (`modul/daftar.ts` `halamanAktif`).
 *
 * ⛔ Bentuk yang tidak dikenal menjadi `null` (= semua menu tampil), BUKAN
 * daftar kosong: badan yang tak terbaca bukan pernyataan "tidak ada modul".
 */
export async function ambilModulAktif(): Promise<readonly string[] | null> {
  const data = await minta<{ modul?: unknown }>('/api/modul-aktif')
  const modul = data?.modul
  if (!Array.isArray(modul) || !modul.every((m): m is string => typeof m === 'string')) return null
  return modul
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
// Gerbang Close Claim — `CloseClaim_Section.xml` b1081 -> b1101.
// ---------------------------------------------------------------------------

/** Satu peserta yang menahan penutupan klaim. */
export interface PenghalangTutup {
  urutan: number
  nomorSertifikat: string
  /** Kalimat yang dilihat pemakai, disusun SERVER dan verbatim dari rule. */
  pesan: string
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
export async function mintaFormulir<T>(jalur: string, isi: FormData): Promise<T> {
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
  if (!jawab.ok) throw kegagalanDari(jawab.status, teks)
  return uraiJSON(jawab.status, teks) as T
}

/**
 * Mengunduh berkas lewat `fetch` BERHEADER IDENTITAS lalu menyerahkannya ke
 * peramban sebagai Blob.
 *
 * ⛔ BUKAN `<a href>`. Tautan biasa tidak membawa `X-Pelaku`, dan rute unduh
 * bergerbang identitas seperti rute lain — tautan biasa dijawab 401. Jalur ini
 * satu-satunya cara layar modul ini mengunduh.
 *
 * ⛔ Jawaban galat diurai sebagai amplop `{"galat": ...}` yang sama dengan
 * `minta`, sehingga 409 "rekam tanpa berkas" tampil dengan kalimatnya.
 */
export async function unduhBerkasBeridentitas(jalur: string, namaBerkas: string): Promise<void> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  let jawab: Response
  try {
    jawab = await fetch(rakitURL(jalur), { headers: { ...headerIdentitas() }, signal: kendali.signal })
  } finally {
    clearTimeout(jam)
  }
  if (!jawab.ok) {
    let pesan: string | undefined
    try {
      const o = (await jawab.json()) as { galat?: unknown }
      pesan = typeof o.galat === 'string' && o.galat !== '' ? o.galat : undefined
    } catch {
      pesan = undefined
    }
    throw new ApiFailure(jawab.status, { code: 'DITOLAK_BACKEND', message: pesan })
  }
  const blob = await jawab.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = namaBerkas
  document.body.appendChild(a)
  a.click()
  a.remove()
  // Temuan /code-review: mencabut URL di tik yang sama dengan klik dapat membatalkan
  // unduhan (Firefox, Safari) - dicabut sesudah peramban mulai mengunduh.
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}
