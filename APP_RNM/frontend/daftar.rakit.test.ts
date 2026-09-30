import { describe, expect, it } from 'vitest'

import type { MenuModul, RuteModul } from '../inti/frontend/modul'
import { rakitModulFrontend, type BerkasMenu, type BerkasRute, type Halaman } from './daftar'

// Perakit daftar modul frontend - struktur tim satu folder per modul
// (30-09-2026). Daftarnya terbentuk dari FOLDER (`import.meta.glob`), jadi yang
// dijaga di sini adalah pasangan `menu.ts` + `rute.tsx` per folder - dengan
// modul tiruan, tanpa bergantung pada modul yang kebetulan ada.

const ruteKosong: RuteModul<Halaman> = () => null

function menu(nama: string): { PENDAFTARAN_MENU: MenuModul } {
  return { PENDAFTARAN_MENU: { nama, kelompok: `Kelompok ${nama}`, halaman: [], menu: [] } }
}

describe('rakitModulFrontend', () => {
  it('memasangkan menu.ts dan rute.tsx per folder, berurutan menurut nama folder', () => {
    const m: BerkasMenu = { '../modul/beta/frontend/menu.ts': menu('beta'), '../modul/alfa/frontend/menu.ts': menu('alfa') }
    const r: BerkasRute = {
      '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong },
      '../modul/beta/frontend/rute.tsx': { RUTE_MODUL: ruteKosong },
    }
    const hasil = rakitModulFrontend(m, r)
    expect(hasil.map((x) => x.nama)).toEqual(['alfa', 'beta'])
    expect(hasil[0]?.kelompok).toBe('Kelompok alfa')
    expect(hasil[0]?.Rute).toBe(ruteKosong)
  })

  it('folder yang hanya punya salah satu berkasnya DITOLAK dengan menyebut foldernya', () => {
    expect(() => rakitModulFrontend({ '../modul/alfa/frontend/menu.ts': menu('alfa') }, {})).toThrow(
      'modul/alfa/frontend punya menu.ts tanpa rute.tsx',
    )
    expect(() => rakitModulFrontend({}, { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } })).toThrow(
      'modul/alfa/frontend punya rute.tsx tanpa menu.ts',
    )
  })

  it('nama modul yang berbeda dari nama foldernya DITOLAK', () => {
    expect(() =>
      rakitModulFrontend(
        { '../modul/alfa/frontend/menu.ts': menu('lain') },
        { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } },
      ),
    ).toThrow('menyebut nama "lain"')
  })
})
