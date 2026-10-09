// Klik baris membuka rincian (9 Oktober 2026). Penyaringan unsur interaktif
// memerlukan DOM; di sini dipastikan baris biasa memanggil `bolak`.
import { describe, expect, it, vi } from 'vitest'

import { klikBaris } from './komponen/klikBaris'

describe('klikBaris', () => {
  it('klik di baris (bukan unsur interaktif) membuka/menutup rincian', () => {
    const bolak = vi.fn()
    klikBaris({ target: null }, bolak)
    expect(bolak).toHaveBeenCalledTimes(1)
  })
  it('unsur interaktif diabaikan', () => {
    const bolak = vi.fn()
    const tombol = { closest: () => ({}) }
    vi.stubGlobal('Element', function Element() {})
    Object.setPrototypeOf(tombol, (globalThis as unknown as { Element: { prototype: object } }).Element.prototype)
    klikBaris({ target: tombol as unknown as EventTarget }, bolak)
    expect(bolak).not.toHaveBeenCalled()
    vi.unstubAllGlobals()
  })
})
