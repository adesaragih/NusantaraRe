import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'

import ts from 'typescript'
import { describe, expect, it } from 'vitest'

import { akarSumberFrontend, berkasTS, relatifAplikasi } from './uji/sumber'

// Penjaga LAPISAN frontend - refactor bentuk B paket 8 (30-09-2026).
//
// Padanan `inti/backend/penjaga/impor_lintas_modul_test.go` di backend.
// Aturannya, atas SETIAP berkas .ts/.tsx kode frontend (termasuk berkas uji):
//
//   1. `inti/**` hanya mengimpor `inti/**` - inti tidak mengenal modul.
//   2. `modul/X/**` hanya mengimpor `inti/**` dan `modul/X/**`. Yang mengenal
//      semua modul hanya lapisan aplikasi: `App.tsx`, `Beranda.tsx`,
//      `main.tsx`, dan daftar modul `frontend/daftar.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): jalur dibaca relatif
// APP_RNM - `inti/frontend/**` adalah inti, `modul/<nama>/frontend/**` adalah
// modul, dan selebihnya (`frontend/**`, perakit) adalah aplikasi.
//
// ⚠️ Impor dibaca pengurai TypeScript (`preProcessFile`), bukan pola teks:
// teks yang MENYEBUT sebuah impor di dalam string atau komentar tidak
// dihitung, dan impor dinamis `import('...')` ikut dihitung.

/** Lapisan sebuah berkas (jalur relatif APP_RNM, bergaris-miring). */
function lapisan(rel: string): string {
  if (rel.startsWith('inti/frontend/')) return 'inti'
  const m = /^modul\/([^/]+)\/frontend\//.exec(rel)
  if (m) return `modul:${m[1] ?? ''}`
  return 'aplikasi'
}

/** Alasan `dari` tidak boleh mengimpor `ke`; `null` = boleh. */
function pelanggaranLapisan(dari: string, ke: string): string | null {
  const a = lapisan(dari)
  const b = lapisan(ke)
  if (a === 'aplikasi') return null
  if (a === 'inti') return b === 'inti' ? null : 'inti tidak mengenal modul maupun lapisan aplikasi'
  if (b === 'inti' || b === a) return null
  return 'modul hanya mengimpor inti/ dan dirinya sendiri; yang merakit modul hanya frontend/daftar.ts dan App.tsx'
}

/** Setiap impor RELATIF tiap berkas, sebagai jalur relatif APP_RNM. */
function imporRelatif(): { dari: string; ke: string; spec: string }[] {
  const hasil: { dari: string; ke: string; spec: string }[] = []
  for (const f of berkasTS()) {
    const dari = relatifAplikasi(f)
    for (const imp of ts.preProcessFile(readFileSync(f, 'utf8'), true, true).importedFiles) {
      if (!imp.fileName.startsWith('.')) continue // paket npm
      const ke = relatifAplikasi(resolve(dirname(f), imp.fileName))
      hasil.push({ dari, ke, spec: imp.fileName })
    }
  }
  return hasil
}

describe('lapisan frontend: inti <- modul <- aplikasi', () => {
  it('nol impor lintas lapisan di seluruh kode frontend', () => {
    const impor = imporRelatif()
    const perLapis = new Map<string, number>()
    for (const { dari } of impor) perLapis.set(lapisan(dari), (perLapis.get(lapisan(dari)) ?? 0) + 1)
    // ⛔ Penjaga yang membaca nol impor di satu lapis lulus atas apa pun.
    expect(perLapis.get('inti') ?? 0).toBeGreaterThan(30)
    // Setiap modul yang punya `frontend/` - dari folder, bukan daftar nama:
    // modul baru ikut terperiksa tanpa menyunting penjaga ini.
    const modul = akarSumberFrontend()
      .map((a) => /^modul\/([^/]+)\/frontend$/.exec(relatifAplikasi(a))?.[1])
      .filter((m): m is string => m !== undefined)
    expect(modul.length, 'modul berfrontend').toBeGreaterThanOrEqual(4)
    for (const m of modul) {
      expect(perLapis.get(`modul:${m}`) ?? 0, `impor modul ${m}`).toBeGreaterThan(0)
    }
    const langgar = impor
      .map((i) => ({ ...i, alasan: pelanggaranLapisan(i.dari, i.ke) }))
      .filter((i) => i.alasan !== null)
      .map((i) => `${i.dari} -> ${i.spec}: ${i.alasan}`)
    expect(langgar).toEqual([])
  })

  it('aturannya menggigit dua arah', () => {
    // ⛔ Nol nama modul sungguhan: `alfa`/`beta` tiruan.
    const kasus: [string, string, boolean][] = [
      ['modul/alfa/frontend/pages/X.tsx', 'modul/beta/frontend/api.ts', false],
      ['modul/beta/frontend/pages/X.test.ts', 'modul/alfa/frontend/labels.ts', false],
      ['modul/beta/frontend/rute.tsx', 'frontend/daftar.ts', false],
      ['modul/beta/frontend/rute.tsx', 'frontend/App.tsx', false],
      ['modul/beta/frontend/rute.tsx', 'modul/beta/frontend/pages/Inbox.tsx', true],
      ['modul/beta/frontend/api.ts', 'inti/frontend/klien.ts', true],
      ['inti/frontend/components/Shell.tsx', 'frontend/daftar.ts', false],
      ['inti/frontend/lib/daftarMenu.ts', 'modul/alfa/frontend/labels.ts', false],
      ['inti/frontend/components/Shell.tsx', 'inti/frontend/lib/daftarMenu.ts', true],
      ['frontend/daftar.ts', 'modul/alfa/frontend/rute.tsx', true],
      ['frontend/App.tsx', 'modul/alfa/frontend/api.ts', true],
      // Tabrakan awalan: `betax` diawali nama modul ini, tetapi modul LAIN.
      ['modul/betax/frontend/api.ts', 'modul/beta/frontend/api.ts', false],
      // Berkas `modul/<nama>/` di LUAR frontend/ bukan kode frontend modul itu.
      ['modul/alfa/frontend/api.ts', 'modul/alfa/docs/x.ts', false],
    ]
    for (const [dari, ke, boleh] of kasus) {
      expect(pelanggaranLapisan(dari, ke) === null, `${dari} -> ${ke}`).toBe(boleh)
    }
  })
})
