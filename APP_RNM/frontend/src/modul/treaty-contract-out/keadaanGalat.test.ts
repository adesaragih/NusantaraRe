import { afterEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, pesanGalat } from '../../inti/klien'
import { ambilKursTahun, konversiKurs } from './api'

import { KODE_BACKEND_MATI, klasifikasiGalat } from '../../inti/lib/keadaanGalat'

// 503 DARI backend lawan 503 dari proxy — lanjutan 6 Treaty Contract Out.
//
// ⛔ Sebab uji ini ada, dan itu sungguh terjadi (laporan work owner
// 29-09-2026): `List Description` → `Show TreatyDesc` menampilkan "Backend
// tidak terhubung (127.0.0.1:8080)" padahal backend menyala. Rute kurs menjawab
// 503 `{"galat": "services: master kurs … tidak dapat dipakai: …"}` — penolakan
// backend sendiri — dan klien memperlakukan SETIAP 503 sebagai backend mati.
// Kalimat backend yang tepat tertelan, dan pemakai disuruh menyalakan backend
// yang sudah menyala.
//
// Pasangan uji sisi backend: `internal/handlers/galat503_test.go`.

const KALIMAT = 'services: master kurs atau mata uang tidak dapat dipakai: UJI sebab'

/** Menjawab setiap permintaan dengan `badan` dan `status`. */
function jawab(status: number, badan: string, tipe = 'application/json'): void {
  vi.stubGlobal('fetch', () =>
    Promise.resolve(new Response(badan, { status, headers: { 'Content-Type': tipe } })),
  )
}

/** Galat yang dilempar sebuah panggilan, atau `null` bila ia berhasil. */
function galatDari(p: Promise<unknown>): Promise<unknown> {
  return p.then(
    () => null,
    (e: unknown) => e,
  )
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('503 yang MEMBAWA galat backend', () => {
  it('kurs tahun: kalimat backend sampai ke pemakai, bukan "backend tidak terhubung"', async () => {
    jawab(503, JSON.stringify({ galat: KALIMAT }))

    const k = klasifikasiGalat(await galatDari(ambilKursTahun('1000001')))

    expect(k?.jenis).toBe('galat-api')
    expect(k?.pesan).toBe(KALIMAT)
    // Petunjuk "jalankan backend" tidak boleh ikut: backend-nya menyala.
    expect(k?.petunjuk).toBeUndefined()
  })

  it('konversi kurs: rute kedua layar yang sama', async () => {
    jawab(503, JSON.stringify({ galat: KALIMAT }))

    const k = klasifikasiGalat(await galatDari(konversiKurs('1000001', 'Rp', '1', '8')))

    expect(k?.jenis).toBe('galat-api')
    expect(k?.pesan).toBe(KALIMAT)
  })

  it.each([502, 504])('%i yang membawa galat backend pun bukan backend mati', async (status) => {
    jawab(status, JSON.stringify({ galat: KALIMAT }))

    const k = klasifikasiGalat(await galatDari(ambilKursTahun('1000001')))

    expect(k?.jenis).toBe('galat-api')
    expect(k?.pesan).toBe(KALIMAT)
  })
})

describe('jawaban TANPA galat backend tetap backend mati', () => {
  it.each([502, 503, 504])('%i berbadan kosong (reverse-proxy)', async (status) => {
    jawab(status, '', 'text/plain')

    const k = klasifikasiGalat(await galatDari(ambilKursTahun('1000001')))

    expect(k?.jenis).toBe('backend-mati')
    expect(k?.petunjuk).toContain('muat-env.ps1')
  })

  it('503 berbadan HTML (reverse-proxy) → BACKEND_TIDAK_TERJANGKAU', async () => {
    jawab(503, '<html><body>Service Unavailable</body></html>', 'text/html')

    const k = klasifikasiGalat(await galatDari(ambilKursTahun('1000001')))

    expect(k?.jenis).toBe('backend-mati')
  })

  it('503 dengan galat KOSONG bukan kalimat backend', async () => {
    jawab(503, JSON.stringify({ galat: '' }))

    expect(klasifikasiGalat(await galatDari(ambilKursTahun('1000001')))?.jenis).toBe('backend-mati')
  })

  it('503 dengan kunci `error` bukan kalimat backend (amplopnya `galat`)', async () => {
    jawab(503, JSON.stringify({ error: KALIMAT }))

    expect(klasifikasiGalat(await galatDari(ambilKursTahun('1000001')))?.jenis).toBe('backend-mati')
  })

  it('kode BACKEND_TIDAK_TERJANGKAU menang atas pesan buatan klien', () => {
    // Pesan kode ini ditulis KLIEN (`bukanJSON`), bukan backend — ia tidak
    // boleh terbaca sebagai galat backend pada status apa pun.
    const g = { status: 200, detail: { code: KODE_BACKEND_MATI, message: 'Jawaban dari server bukan JSON' } }

    expect(klasifikasiGalat(g)?.jenis).toBe('backend-mati')
  })

  it('galat JARINGAN tetap backend mati', () => {
    expect(klasifikasiGalat(new TypeError('Failed to fetch'))?.jenis).toBe('backend-mati')
  })
})

describe('satu aturan "kalimat dari backend", dua tempat', () => {
  // ⛔ `inti/lib/keadaanGalat.ts` murni dan tidak mengimpor `inti/klien.ts`, jadi
  // aturannya DISALIN dari `pesanGalat`. Uji ini mengunci keduanya sepakat:
  // 503 menjadi galat-api TEPAT bila `pesanGalat` mengakui pesannya milik
  // backend. Aturan yang menyimpang akan menampilkan kalimat klien sebagai
  // kalimat server, atau menyuruh menyalakan backend yang menyala.
  it.each([
    ['DITOLAK_BACKEND + pesan', new ApiFailure(503, { code: 'DITOLAK_BACKEND', message: KALIMAT })],
    ['DITOLAK_BACKEND + pesan kosong', new ApiFailure(503, { code: 'DITOLAK_BACKEND', message: '' })],
    ['DITOLAK_BACKEND tanpa pesan', new ApiFailure(503, { code: 'DITOLAK_BACKEND' })],
    ['BACKEND_TIDAK_TERJANGKAU + pesan klien', new ApiFailure(503, { code: 'BACKEND_TIDAK_TERJANGKAU', message: 'x' })],
    ['tanpa kode + pesan', new ApiFailure(503, { message: KALIMAT })],
  ])('%s', (_nama, e) => {
    const dariBackend = pesanGalat(e) !== undefined
    expect(klasifikasiGalat(e)?.jenis).toBe(dariBackend ? 'galat-api' : 'backend-mati')
  })
})
