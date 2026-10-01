import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Layar login tanpa kedip (laporan work owner 01-10-2026: "halaman login suka
// kedip kedip"). Penjaga statik atas App.tsx - dua sebab yang ditemukan:
//
//   1. Pemeriksaan sesi dulu merender paragraf polos di halaman putih, lalu
//      berganti ke kartu login bergradien: kedip di SETIAP muat ulang.
//   2. `GET /api/menu` dikirim saat App menyala, SEBELUM login. Tanpa sesi ia
//      dijawab 401, dan 401 memicu `PERISTIWA_SESI_BERAKHIR` di tengah
//      pemeriksaan sesi.
const APP = readFileSync(join(__dirname, 'App.tsx'), 'utf8')

describe('App: pemeriksaan sesi tanpa kedip', () => {
  it('layar periksa sesi memakai latar layar login', () => {
    expect(APP).toContain('<main className="halaman-masuk" aria-busy="true">')
    expect(APP).not.toContain('className="polis__catatan"')
  })

  it('menu dibaca sesudah identitas diketahui, ulang bila akunnya berganti', () => {
    expect(APP).toContain("const kunciMenu = stub ? 'stub' : (profil?.akunId ?? null)")
    expect(APP).toContain('if (kunciMenu === null) return')
    expect(APP).toContain('}, [kunciMenu])')
  })
})
