// Unduh dokumen membawa identitas pelaku — GILIRAN-12 paket 0.
//
// ⛔ CACAT YANG DITUTUP: `View Office Online` (`DocumentLife.xml` b3502)
// dulu `<a href={tautanDokumen(id)} download>`. Browser yang mengikuti pranala
// TIDAK membawa header `X-Pelaku`/`X-Peran`, sedangkan backend membaca
// identitas HANYA dari header itu (`handlers/pelaku.go`) - setiap unduhan
// dijawab 401. Kedua sisi benar menurut dirinya sendiri; pertemuannya yang
// salah. Sisi Go dikunci `handlers/dokumen_identitas_test.go`.

import { readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { PERAN } from '../../inti/labels'

import { ambilIsiDokumen } from './api'
import { pesanGalat } from '../../inti/klien'

interface Tangkapan {
  url: string
  init: RequestInit
}

let tertangkap: Tangkapan[] = []

function pasangFetch(status: number, badan: BodyInit, tipe: string): void {
  vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
    tertangkap.push({ url, init })
    return Promise.resolve(new Response(badan, { status, headers: { 'Content-Type': tipe } }))
  })
}

beforeEach(() => {
  tertangkap = []
  vi.stubEnv('VITE_AUTH_STUB', 'true')
  vi.stubEnv('VITE_STUB_PELAKU', 'UJI-ADMIN')
  vi.stubEnv('VITE_STUB_PERAN', PERAN.admin)
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.unstubAllEnvs()
})

describe('ambilIsiDokumen', () => {
  it('memakai GET ke rute isi dokumen DENGAN header identitas', async () => {
    pasangFetch(200, 'isi-uji', 'application/pdf')
    await ambilIsiDokumen('20250101120000123')
    expect(tertangkap).toHaveLength(1)
    const t = tertangkap[0]!
    expect(t.url).toBe('/api/dokumen/20250101120000123/isi')
    expect(t.init.method).toBe('GET')
    const h = (t.init.headers ?? {}) as Record<string, string>
    expect(h['X-Pelaku']).toBe('UJI-ADMIN')
    expect(h['X-Peran']).toBe(PERAN.admin)
  })

  it('mengembalikan isi berkas sebagai Blob, bukan teks yang diurai', async () => {
    pasangFetch(200, 'isi-uji', 'application/pdf')
    const b = await ambilIsiDokumen('1')
    expect(b).toBeInstanceOf(Blob)
    expect(b.type).toBe('application/pdf')
    expect(await b.text()).toBe('isi-uji')
  })

  it('penolakan backend memakai amplop galat yang SAMA', async () => {
    pasangFetch(401, JSON.stringify({ galat: 'permintaan tanpa identitas pelaku ditolak' }),
      'application/json')
    const e = await ambilIsiDokumen('1').catch((x: unknown) => x)
    expect(pesanGalat(e)).toBe('permintaan tanpa identitas pelaku ditolak')
  })
})

/** Seluruh berkas sumber non-uji di bawah src/. */
function sumberSrc(): { jalur: string; isi: string }[] {
  const hasil: { jalur: string; isi: string }[] = []
  const jelajah = (dir: string): void => {
    for (const e of readdirSync(dir, { withFileTypes: true })) {
      const p = join(dir, e.name)
      if (e.isDirectory()) jelajah(p)
      else if (/\.(ts|tsx)$/.test(e.name) && !/\.test\.(ts|tsx)$/.test(e.name)) {
        hasil.push({ jalur: p, isi: readFileSync(p, 'utf8') })
      }
    }
  }
  jelajah(join(__dirname, '..', '..'))
  return hasil
}

describe('nol pranala ke /api/ — setiap permintaan lewat fetch yang membawa identitas', () => {
  it('tidak ada pranala, src, atau navigasi yang dapat menuju backend', () => {
    const berkas = sumberSrc()
    expect(berkas.length).toBeGreaterThan(30)
    // ⛔ `href={...}` DINAMIS dilarang seluruhnya, bukan hanya yang memuat
    // "/api/": ronde pertama penjaga ini mencari "/api/" harfiah dan LOLOS
    // terhadap `href={tautanDokumen(b.id)}` - cacat yang justru dijaganya.
    const terlarang = [
      /\b(href|src)=\{/,
      /\b(href|src)\s*=\s*["'`][^"'`]*\/api\//,
      /window\.open\(/,
      /location\.(href|assign|replace)\b/,
    ]
    const pelanggar = berkas
      .flatMap((b) => b.isi.split('\n').map((baris, i) => ({ b, baris, i })))
      .filter(({ baris }) => !baris.trim().startsWith('//') && !baris.trim().startsWith('*'))
      .filter(({ b, baris }) => {
        if (terlarang.some((p) => p.test(baris))) return true
        // `a.href = X` hanya bila X sendiri objek URL lokal:
        // `const X = URL.createObjectURL(...)` di berkas yang sama. Diperiksa
        // PER PENUGASAN (ronde pertama per berkas), dan `===` bukan penugasan.
        const m = /\.href\s*=(?!=)\s*([A-Za-z_$][\w$]*)/.exec(baris)
        if (m === null) return /\.href\s*=(?!=)/.test(baris)
        return !new RegExp(`\\b${m[1]}\\s*=\\s*URL\\.createObjectURL\\(`).test(b.isi)
      })
      .map(({ b, baris, i }) => `${b.jalur}:${i + 1}: ${baris.trim()}`)
    expect(pelanggar).toEqual([])
  })

  it('alamat backend hanya dirakit DI DALAM panggilan fetch', () => {
    // ⛔ Fungsi yang mengembalikan URL backend untuk dipasang di pranala
    // adalah jalan pintas melewati header identitas - persis `tautanDokumen`.
    // ⚠️ Diperiksa per BARIS: `fetch(` dan `rakitURL(` wajib satu baris. Itu
    // aturan bentuk yang disengaja - pembungkus baris yang memisahkannya
    // membuat uji ini merah, bukan lolos.
    // Refactor bentuk B: `rakitURL` kini DIEKSPOR `inti/klien.ts` supaya klien
    // tiap modul dapat memanggil fetch-nya sendiri. Karena itu yang dibaca
    // SELURUH sumber, bukan satu berkas api: pemakai `rakitURL` di luar fetch
    // di berkas mana pun menjadi merah.
    const berkas = sumberSrc()
    const pemakaian = berkas
      .flatMap((b) => b.isi.split('\n'))
      .filter((b) => b.includes('rakitURL(') && !b.includes('function rakitURL'))
    expect(pemakaian.length).toBeGreaterThan(0)
    for (const b of pemakaian) expect(b).toContain('fetch(rakitURL(')
    for (const b of berkas) expect(b.isi).not.toContain('export function tautanDokumen')
  })
})
