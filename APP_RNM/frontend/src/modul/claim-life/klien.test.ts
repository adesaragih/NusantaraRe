// Uji klien HTTP — F0.2.
//
// ⛔ Yang diuji di sini adalah PERMINTAAN YANG BENAR-BENAR DIKIRIM, bukan
// bentuk fungsinya. `fetch` ditiru dan argumennya diperiksa: itu satu-satunya
// cara membuktikan header identitas sungguh ikut, dan sungguh hilang saat
// keluar.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { berkasTS } from '../../../../inti/frontend/uji/sumber'
import { PERAN } from '../../../../inti/frontend/labels'
import { klasifikasiGalat } from '../../../../inti/frontend/lib/keadaanGalat'

import { ambilKlaimLife, hapusKlaim, tolakBarisAdjustment } from './api'
import { pesanGalat } from '../../../../inti/frontend/klien'

/** Satu panggilan fetch yang tertangkap. */
interface Tangkapan {
  url: string
  init: RequestInit
}

let tertangkap: Tangkapan[] = []

/** Identitas kini datang dari env, bukan dari penyimpanan (F0.6). */
function pasangIdentitas(akun: string | null, ...peran: string[]): void {
  vi.stubEnv('VITE_AUTH_STUB', akun === null ? 'false' : 'true')
  vi.stubEnv('VITE_STUB_PELAKU', akun ?? '')
  vi.stubEnv('VITE_STUB_PERAN', peran.join(','))
}

/** Memasang `fetch` tiruan yang menjawab `badan` dengan `status`. */
function pasangFetch(status: number, badan: string, tipe = 'application/json'): void {
  vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
    tertangkap.push({ url, init })
    // ⚠️ 204/205/304 WAJIB berbadan null - konstruktor `Response`
    // menolak badan pada status itu, termasuk string kosong.
    const isi = status === 204 || status === 205 || status === 304 ? null : badan
    return Promise.resolve(
      new Response(isi, { status, headers: { 'Content-Type': tipe } }),
    )
  })
}

/** Header sebuah tangkapan sebagai objek biasa. */
function header(t: Tangkapan): Record<string, string> {
  return (t.init.headers ?? {}) as Record<string, string>
}

beforeEach(() => {
  tertangkap = []
  pasangIdentitas('UJI-ADMIN', PERAN.admin)
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.unstubAllEnvs()
})

describe('identitas ikut di setiap permintaan', () => {
  it('sesudah masuk, X-Pelaku dan X-Peran TERKIRIM', async () => {
    pasangIdentitas('UJI-ADMIN', PERAN.admin)
    pasangFetch(200, JSON.stringify({ id: 'RNML-000001' }))

    await ambilKlaimLife('RNML-000001')

    expect(tertangkap).toHaveLength(1)
    const h = header(tertangkap[0]!)
    expect(h['X-Pelaku']).toBe('UJI-ADMIN')
    expect(h['X-Peran']).toBe(PERAN.admin)
  })

  it('peran ganda terkirim dipisah koma', async () => {
    pasangIdentitas('UJI-SEMUA', PERAN.admin, PERAN.spv)
    pasangFetch(200, JSON.stringify({ id: 'X' }))

    await ambilKlaimLife('X')

    expect(header(tertangkap[0]!)['X-Peran']).toBe(`${PERAN.admin},${PERAN.spv}`)
  })

  it('stub MATI → header identitas tidak terkirim', async () => {
    pasangIdentitas(null)
    pasangFetch(200, JSON.stringify({ id: 'X' }))

    await ambilKlaimLife('X')

    const h = header(tertangkap[0]!)
    expect(h['X-Pelaku']).toBeUndefined()
    expect(h['X-Peran']).toBeUndefined()
  })
})

describe('metode HTTP yang benar dikirim', () => {
  it('hapusKlaim memakai DELETE, bukan GET', async () => {
    // ⛔ Rutenya SAMA PERSIS dengan ambilKlaimLife; hanya metodenya yang
    // membedakan. Metode yang hilang membuat tombol Hapus diam-diam MEMBACA
    // klaim lalu melaporkan berhasil.
    pasangFetch(200, JSON.stringify({ total: 3 }))
    await hapusKlaim('RNML-000001')
    expect(tertangkap[0]!.init.method).toBe('DELETE')
  })

  it('ambilKlaimLife memakai GET pada jalur yang sama', async () => {
    pasangFetch(200, JSON.stringify({ id: 'RNML-000001' }))
    await ambilKlaimLife('RNML-000001')
    expect(tertangkap[0]!.init.method).toBe('GET')
  })

  it('tolakBarisAdjustment memakai POST dan membawa Remarks (OQ-M5)', async () => {
    pasangFetch(204, '')
    await tolakBarisAdjustment('RNML-000001', 'ADJ-1', 'UJI alasan')
    expect(tertangkap[0]!.init.method).toBe('POST')
    expect(JSON.parse(String(tertangkap[0]!.init.body))).toEqual({ komentar: 'UJI alasan' })
  })
})

describe('kegagalan dikenali menurut jenisnya', () => {
  it('jawaban BUKAN JSON menjadi BACKEND_TIDAK_TERJANGKAU', async () => {
    // Proxy yang gagal menghubungi upstream mengirim HTML, bukan JSON.
    pasangFetch(502, '<html><body>Bad Gateway</body></html>', 'text/html')

    const galat = await ambilKlaimLife('X').then(
      () => null,
      (e: unknown) => e,
    )

    const keadaan = klasifikasiGalat(galat)
    expect(keadaan?.jenis).toBe('backend-mati')
    // ⛔ Petunjuknya menyebut langkah konkret, bukan "terjadi kesalahan".
    expect(keadaan?.petunjuk).toContain('muat-env.ps1')
  })

  it('galat JARINGAN tetap TypeError dan terbaca sebagai backend mati', async () => {
    vi.stubGlobal('fetch', () => Promise.reject(new TypeError('Failed to fetch')))

    const galat = await ambilKlaimLife('X').then(
      () => null,
      (e: unknown) => e,
    )

    expect(galat).toBeInstanceOf(TypeError)
    expect(klasifikasiGalat(galat)?.jenis).toBe('backend-mati')
  })

  it('penolakan backend meneruskan kalimatnya APA ADANYA', async () => {
    // ⛔ Kuncinya `galat`, dan uji ini SEMPAT MENGUNCI YANG SALAH.
    //
    // Sampai 27-09-2026 baris di bawah menyuapkan `{error: ...}` lalu
    // menuntut kalimatnya lolos - dan ia hijau, sebab klien memang membaca
    // `error`. Yang diuji bukan kontrak dengan backend, melainkan kontrak
    // klien dengan dirinya sendiri; backend Go menulis `galat` sejak tiket
    // 01 dan tidak pernah menulis `error`.
    //
    // Itulah sebab cacatnya bertahan: ADA ujinya, dan ujinya ikut keliru.
    // Sejak sekarang kontraknya dikunci dua sisi - `envelopegalat.test.ts`
    // di sini dan `envelopegalat_test.go` di Go.
    pasangFetch(422, JSON.stringify({ galat: 'Name of bank cannot be empty' }))

    const galat = await ambilKlaimLife('X').then(
      () => null,
      (e: unknown) => e,
    )

    expect(pesanGalat(galat)).toBe('Name of bank cannot be empty')
  })

  it('badan KOSONG pada 204 bukan kegagalan', async () => {
    pasangFetch(204, '')
    await expect(tolakBarisAdjustment('K', 'A', 'UJI alasan')).resolves.toBeUndefined()
  })
})

describe('axios benar-benar dilepas', () => {
  it('nol impor axios di seluruh src/', () => {
    // ⛔ Penjaga arah-balik. Dua klien HTTP berdampingan berarti dua model
    // galat, dan yang satu akan diam-diam kalah - layar menampilkan "terjadi
    // kesalahan" untuk penolakan yang sebenarnya membawa kalimat server.
    // Seluruh akar kode frontend (`inti/frontend/uji/sumber.ts`), bukan satu
    // folder: sejak struktur tim satu folder per modul kodenya tersebar.
    const berkas = berkasTS()
    expect(berkas.length).toBeGreaterThan(100)
    const tertuduh = berkas.filter((p) => /from ['"]axios['"]/.test(readFileSync(p, 'utf8')))
    expect(tertuduh).toEqual([])
  })

  it('nol axios di package.json', () => {
    // package.json tinggal di APP_RNM/ sejak struktur tim satu folder per modul.
    const pkg = readFileSync(join(__dirname, '..', '..', '..', '..', 'package.json'), 'utf8')
    expect(JSON.parse(pkg) as { dependencies?: Record<string, string> }).not.toHaveProperty(
      'dependencies.axios',
    )
  })
})
