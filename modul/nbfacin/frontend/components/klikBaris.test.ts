// Klik baris = buka / tutup detail (work owner 10-10-2026). Elemen tiruan - tanpa DOM.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { tanganiKlikBaris, tombolBukaBaris } from './klikBaris'

type El = { tagName: string; children: El[]; classList?: { contains: (c: string) => boolean }; closest?: (s: string) => El | null; disabled?: boolean; click?: () => void }

/** Baris tiruan: td(tombol .nbf-buka) + td(teks). `closest` meniru pencarian ke atas. */
function baris(opsi: { disabled?: boolean } = {}) {
  let diklik = 0
  const tombol: El = { tagName: 'BUTTON', children: [], classList: { contains: (c) => c === 'nbf-buka' }, disabled: opsi.disabled, click: () => diklik++ }
  const tr: El = { tagName: 'TR', children: [] }
  const tdTombol: El = { tagName: 'TD', children: [tombol] }
  const teks: El = { tagName: 'TD', children: [] }
  tr.children.push(tdTombol, teks)
  teks.closest = (s) => (s === 'tr' ? tr : null)
  tombol.closest = (s) => (s.includes('button') ? tombol : s === 'tr' ? tr : null)
  return { tr, teks, tombol, jumlah: () => diklik }
}
const klik = (target: unknown, ubah: Partial<{ button: number; defaultPrevented: boolean }> = {}) => ({ target, button: 0, defaultPrevented: false, ...ubah })

describe('klikBaris', () => {
  it('klik teks di baris (judul / nama) meneruskan ke tombol buka', () => {
    const b = baris()
    expect(tanganiKlikBaris(klik(b.teks), () => '')).toBe(true)
    expect(b.jumlah()).toBe(1)
  })

  it('klik pada tombol / isian sendiri tidak diteruskan (aksinya jalan sendiri)', () => {
    const b = baris()
    expect(tanganiKlikBaris(klik(b.tombol), () => '')).toBe(false)
    expect(b.jumlah()).toBe(0)
  })

  it('tidak diteruskan: tombol nonaktif, sedang menyeleksi teks, tombol mouse lain, baris tanpa .nbf-buka', () => {
    const nonaktif = baris({ disabled: true })
    expect(tanganiKlikBaris(klik(nonaktif.teks), () => '')).toBe(false)
    const b = baris()
    expect(tanganiKlikBaris(klik(b.teks), () => 'teks diseleksi')).toBe(false)
    expect(tanganiKlikBaris(klik(b.teks, { button: 2 }), () => '')).toBe(false)
    const tanpa: El = { tagName: 'TR', children: [{ tagName: 'TD', children: [] }] }
    expect(tombolBukaBaris(tanpa)).toBeNull()
    expect(b.jumlah()).toBe(0)
  })

  it('dipasang sekali di rute modul (semua tabel NB FacIn); baris berkursor tangan', () => {
    const RUTE = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8')
    expect(RUTE).toContain('useEffect(() => pasangKlikBaris(), [])')
    const CSS = readFileSync(join(__dirname, '..', 'nbfacin.css'), 'utf8')
    expect(CSS).toContain('.nbfacin tr:has(> td > .nbf-buka) {\n  cursor: pointer;')
  })
})
