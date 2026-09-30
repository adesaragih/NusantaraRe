import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'

import ts from 'typescript'
import { describe, expect, it } from 'vitest'

import { berkasTS, relatifAplikasi } from './uji/sumber'

// Penjaga LAPISAN frontend - refactor bentuk B paket 8 (30-09-2026).
//
// Padanan `inti/backend/penjaga/impor_lintas_modul_test.go` di backend.
// Aturannya, atas SETIAP berkas .ts/.tsx kode frontend (termasuk berkas uji):
//
//   1. `inti/**` hanya mengimpor `inti/**` - inti tidak mengenal modul.
//   2. `modul/X/**` hanya mengimpor `inti/**` dan `modul/X/**`. Yang mengenal
//      semua modul hanya lapisan aplikasi: `App.tsx`, `Beranda.tsx`,
//      `main.tsx`, dan daftar modul `modul/daftar.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): jalur dibaca relatif
// APP_RNM - `inti/frontend/**` adalah inti, `modul/<nama>/frontend/**` adalah
// modul, dan selebihnya (`frontend/**`, perakit) adalah aplikasi. Selama
// modul masih di `frontend/src/modul/<nama>/`, letak itu dihitung modul juga.
//
// ⚠️ Impor dibaca pengurai TypeScript (`preProcessFile`), bukan pola teks:
// teks yang MENYEBUT sebuah impor di dalam string atau komentar tidak
// dihitung, dan impor dinamis `import('...')` ikut dihitung.

/** Lapisan sebuah berkas (jalur relatif APP_RNM, bergaris-miring). */
function lapisan(rel: string): string {
  if (rel.startsWith('inti/frontend/')) return 'inti'
  const m = /^modul\/([^/]+)\/frontend\//.exec(rel) ?? /^frontend\/src\/modul\/([^/]+)\//.exec(rel)
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
  return 'modul hanya mengimpor inti/ dan dirinya sendiri; yang merakit modul hanya modul/daftar.ts dan App.tsx'
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
    for (const m of ['claim-life', 'premiumlist-life', 'komite-claim-life', 'treaty-contract-out']) {
      expect(perLapis.get(`modul:${m}`) ?? 0, `impor modul ${m}`).toBeGreaterThan(10)
    }
    const langgar = impor
      .map((i) => ({ ...i, alasan: pelanggaranLapisan(i.dari, i.ke) }))
      .filter((i) => i.alasan !== null)
      .map((i) => `${i.dari} -> ${i.spec}: ${i.alasan}`)
    expect(langgar).toEqual([])
  })

  it('aturannya menggigit dua arah', () => {
    const kasus: [string, string, boolean][] = [
      ['modul/claimlife/frontend/pages/X.tsx', 'modul/premiumlistlife/frontend/api.ts', false],
      ['modul/komiteclaimlife/frontend/pages/X.test.ts', 'modul/claimlife/frontend/labels.ts', false],
      ['modul/komiteclaimlife/frontend/rute.tsx', 'frontend/src/modul/daftar.ts', false],
      ['modul/komiteclaimlife/frontend/rute.tsx', 'frontend/src/App.tsx', false],
      ['modul/komiteclaimlife/frontend/rute.tsx', 'modul/komiteclaimlife/frontend/pages/InboxKomite.tsx', true],
      ['modul/komiteclaimlife/frontend/api.ts', 'inti/frontend/klien.ts', true],
      ['inti/frontend/components/Shell.tsx', 'frontend/src/modul/daftar.ts', false],
      ['inti/frontend/lib/daftarMenu.ts', 'modul/treatycontractout/frontend/labels.ts', false],
      ['inti/frontend/components/Shell.tsx', 'inti/frontend/lib/daftarMenu.ts', true],
      ['frontend/src/modul/daftar.ts', 'modul/treatycontractout/frontend/rute.tsx', true],
      ['frontend/src/App.tsx', 'modul/claimlife/frontend/api.ts', true],
      // Tabrakan awalan: `komiteclaimlifex` diawali nama modul ini, tetapi modul LAIN.
      ['modul/komiteclaimlifex/frontend/api.ts', 'modul/komiteclaimlife/frontend/api.ts', false],
      // Berkas `modul/<nama>/` di LUAR frontend/ bukan kode frontend modul itu.
      ['modul/claimlife/frontend/api.ts', 'modul/claimlife/docs/x.ts', false],
      // Letak lama selama paket 2: `frontend/src/modul/<nama>/` tetap modul.
      ['frontend/src/modul/claim-life/pages/X.tsx', 'frontend/src/modul/premiumlist-life/api.ts', false],
      ['frontend/src/modul/claim-life/api.ts', 'inti/frontend/klien.ts', true],
    ]
    for (const [dari, ke, boleh] of kasus) {
      expect(pelanggaranLapisan(dari, ke) === null, `${dari} -> ${ke}`).toBe(boleh)
    }
  })
})
