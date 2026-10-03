import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { TEKS_UI } from '../../lib/teksUI'
import { kataSaatBuka, sorotBerikut } from './pilihSaring'

// Kode saja - baris komentar dibuang (komentarnya MENYEBUT `<datalist>`).
const KODE = readFileSync(join(__dirname, 'pilihSaring.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*') && !b.trimStart().startsWith('/**'))
  .join('\n')

describe('PilihSaring — dropdown yang dapat difilter', () => {
  it('sorotan ArrowUp/ArrowDown dijepit di dalam daftar', () => {
    expect(sorotBerikut(0, 3, 1)).toBe(1)
    expect(sorotBerikut(2, 3, 1)).toBe(2)
    expect(sorotBerikut(0, 3, -1)).toBe(0)
    expect(sorotBerikut(5, 0, 1)).toBe(0)
  })
  it('membuka daftar dengan teks pilihan terakhir = tampilkan SEMUA', () => {
    expect(kataSaatBuka('UJI REAS SATU', 'UJI REAS SATU')).toBe('')
    expect(kataSaatBuka('uji', 'UJI REAS SATU')).toBe('uji')
  })
  it('combobox sungguhan: panah, daftar ber-role, keyboard, bukan datalist', () => {
    expect(KODE).toContain('role="combobox"')
    expect(KODE).toContain('role="listbox"')
    expect(KODE).toContain('className="pilih-saring__panah"')
    for (const k of ['ArrowDown', 'ArrowUp', 'Enter', 'Escape']) expect(KODE).toContain(`"${k}"`)
    expect(KODE).not.toMatch(/<datalist/)
  })
  it('pilihan berubah HANYA saat butir dipilih; ketikan dikembalikan saat fokus pergi', () => {
    const tutup = KODE.slice(KODE.indexOf('function tutup'), KODE.indexOf('function pilih'))
    expect(tutup).toContain('setKetik(teksTerpilih)')
    const ketik = KODE.slice(KODE.indexOf('function ubahKetik'))
    expect(ketik.slice(0, ketik.indexOf('\n  }\n'))).not.toContain('onPilih')
  })
  it('Escape menutup daftar saja, bukan popup tempatnya berada', () => {
    const esc = KODE.slice(KODE.indexOf('e.key === "Escape"'))
    expect(esc.slice(0, 200)).toContain('e.stopPropagation()')
  })
  it('teks bawaan dari TEKS_UI — id tetap untuk modul lain, en lengkap', () => {
    expect(TEKS_UI.id.tidakDiDaftar).toBe('(tidak ada di daftar referensi)')
    expect(TEKS_UI.id.memuat).toBe('Memuat...')
    expect(Object.keys(TEKS_UI.en).sort()).toEqual(Object.keys(TEKS_UI.id).sort())
    expect(TEKS_UI.en.tidakCocok).toBe('No matches')
  })
})
