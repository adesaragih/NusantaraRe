// Uji render tab Co-Ins Scale — mode lihat dan mode ubah.
//
// Dirender statis (`react-dom/server`) — BUKAN pengganti melihat layar di
// peramban, dan klik tidak dapat diuji di sini. Yang dijaga: apa yang
// tampil di tiap mode, sesuai gambar 18 (lihat) dan tangkapan layar
// pemakai 6 Oktober 2026 (ubah).

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisSkalaKoasuransiWarisan } from './api'
import TabCoInsScale from './components/TabCoInsScale'
import { CO_INS_SCALE, KOLOM_COIN_SCALE } from './labels'

const BARIS: BarisSkalaKoasuransiWarisan[] = [
  { bagianKoasuransi: 'Panel A', persenLimit: '12.5', penyusun: '', disusunPada: '' },
]

const kepala = (html: string): string => html.slice(html.indexOf('<thead>'), html.indexOf('</thead>'))
const badan = (html: string): string => html.slice(html.indexOf('<tbody>'), html.indexOf('</tbody>'))

describe('Co-Ins Scale — mode lihat (gambar 18)', () => {
  const html = renderToStaticMarkup(<TabCoInsScale baris={[]} mode="lihat" />)

  it('kepala dua kolom tampil walau kosong, dan kosongnya "No items"', () => {
    for (const k of KOLOM_COIN_SCALE) expect(kepala(html)).toContain(k)
    expect(badan(html)).toContain(CO_INS_SCALE.kosong)
  })

  it('⛔ tanpa tombol tambah/hapus — `TreatyIn.IsEditData!=\'1\'` tidak terpenuhi', () => {
    expect(html).not.toContain(CO_INS_SCALE.tambah)
    expect(html).not.toContain(CO_INS_SCALE.hapus)
  })

  it('kedua medan Max Co tampil sebagai kotak baca-saja', () => {
    expect(html).toContain(CO_INS_SCALE.nonGroup)
    expect(html).toContain(CO_INS_SCALE.group)
    expect(html.match(/<input[^>]*readonly=""/g)?.length).toBe(2)
  })

  it('baris tersimpan dibaca sebagai teks, persen dengan `%`', () => {
    const isi = badan(renderToStaticMarkup(<TabCoInsScale baris={BARIS} mode="lihat" />))
    expect(isi).toContain('Panel A')
    expect(isi).toContain('12,5')
    expect(isi).toContain('%')
    expect(isi).not.toContain('<input')
    expect(isi).not.toContain(CO_INS_SCALE.kosong)
  })
})

describe('Co-Ins Scale — mode ubah (tangkapan layar pemakai)', () => {
  it('tombol `Tambah` DI DALAM kepala grid', () => {
    const html = renderToStaticMarkup(<TabCoInsScale baris={[]} mode="ubah" />)
    expect(kepala(html)).toContain(CO_INS_SCALE.tambah)
    expect(badan(html)).toContain(CO_INS_SCALE.kosong)
  })

  it('tiap baris dapat disunting di tempatnya dan punya tombol `Hapus`', () => {
    const isi = badan(renderToStaticMarkup(<TabCoInsScale baris={BARIS} mode="ubah" />))
    expect(isi).toContain('value="Panel A"')
    expect(isi).toContain('value="12.5"')
    expect(isi).toContain(CO_INS_SCALE.simbolPersen)
    expect(isi).toContain(CO_INS_SCALE.hapus)
  })

  it('kedua medan Max Co dapat diisi', () => {
    const html = renderToStaticMarkup(<TabCoInsScale baris={[]} mode="ubah" />)
    expect(html).not.toContain('readonly=""')
  })
})
