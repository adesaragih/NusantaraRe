// Paritas popup Ceding Co List dengan section Pega `ShowCedingCoList` + tangkapan layar work owner 03-10-2026
// (tiket 34). Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { POPUP_CEDING, TEKS_INWARD } from '../labels'
import PopupCedingCoList from './PopupCedingCoList'

const kosong = renderToStaticMarkup(<PopupCedingCoList awal={[]} onTutup={() => {}} onSubmit={() => {}} />)
const berisi = renderToStaticMarkup(
  <PopupCedingCoList
    awal={[
      { id: 'UJI-C1', name: 'UJI CEDING SATU' },
      { id: 'UJI-C2', name: 'UJI CEDING DUA' },
    ]}
    onTutup={() => {}}
    onSubmit={() => {}}
  />,
)
const SUMBER = readFileSync(join(__dirname, 'PopupCedingCoList.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('PopupCedingCoList', () => {
  it('judul Ceding Co List; kepala Ceding Co bertanda wajib + tombol Add / Select Ceding; Submit di kaki', () => {
    expect(kosong).toContain(`>${POPUP_CEDING.judul}</h3>`)
    expect(kosong).toContain(`${POPUP_CEDING.kolom.label}<span class="field__req">*</span></th>`)
    expect(kosong).toMatch(new RegExp(`<button type="button" class="btn btn--ghost btn--sm">${POPUP_CEDING.tambah.label}</button>`))
    expect(kosong).toMatch(new RegExp(`<button type="button" class="btn btn--primary">${POPUP_CEDING.submit.label}</button>`))
  })

  it('daftar kosong = No items (tangkapan layar)', () => {
    expect(kosong).toContain(`<td colSpan="2">${TEKS_INWARD.kosong}</td>`)
    expect(kosong).not.toContain(`>${POPUP_CEDING.hapus.label}</button>`)
  })

  it('daftar berisi: nama per baris berurutan, masing-masing tombol Delete; kode tidak ditampilkan', () => {
    const nama = [...berisi.matchAll(/<tr><td>([^<]*)<\/td>/g)].map((m) => m[1])
    expect(nama).toEqual(['UJI CEDING SATU', 'UJI CEDING DUA'])
    expect(berisi.split(`>${POPUP_CEDING.hapus.label}</button>`)).toHaveLength(3)
    expect(berisi).not.toContain('UJI-C1')
    expect(berisi).not.toContain(TEKS_INWARD.kosong)
  })

  it('Add / Select Ceding membuka pencari Ceding Company menggantikan daftar (E-6); Choose menambah baris di akhir (E-5)', () => {
    expect(SUMBER).toContain('onClick={() => setMencari(true)}')
    expect(SUMBER).toMatch(/if \(mencari\) \{\s*return \(\s*<PopupPilihAgent\s*judul=\{POPUP_CEDING\.judulCari\}/)
    expect(SUMBER).toContain('setDaftar((d) => [...d, { id: b.id, name: b.name }])')
    expect(SUMBER).toContain('onTutup={() => setMencari(false)}')
  })

  it('Delete membuang baris itu saja; Submit menyerahkan daftar; Cancel tanpa menyerahkan', () => {
    expect(SUMBER).toContain('setDaftar((d) => d.filter((_, j) => j !== i))')
    expect(SUMBER).toContain('onClick={() => onSubmit(daftar)}')
    expect(SUMBER).toContain('onTutup={onTutup}')
  })
})
