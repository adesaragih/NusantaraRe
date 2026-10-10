// Kontrol Number hanya menerima angka — permintaan pemakai 8 Oktober 2026:
// "inputan yg khusus angka tidak boleh input karakter lain". Medan `pxNumber`
// ekspor yang dapat disunting memakai `FieldAngka` (inti), yang membuang
// huruf/simbol di tiap ketukan (`formatKetik`, dipaku `inti/frontend/lib/
// angkaKetik.test.ts`).

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { formatKetik, keKabelAngka } from '../../../inti/frontend/lib/format'
import type { MedanKerangka } from './ekspor/jenis'
import { MedanEkspor, type KonteksKerangka } from './komponen/KerangkaTab'

const sisi = { medan: { QSPct: '80', AdjRate: '0.125', Nama: 'abc', Aneh: 'IDR' }, larik: {} }
const k: KonteksKerangka = { sisi, akar: sisi, halaman: { ViewState: '0' }, ubah: true, ubahMedan: () => undefined }
const medan = (kunci: string, format: string, desimal: number | null = 2): MedanKerangka => ({
  t: 'medan', at: 1, label: kunci, dari: 'sisi', kunci, format, desimal, syarat: [],
})

describe('pxNumber = kotak angka saja', () => {
  it('medan pxNumber yang dapat disunting memakai kotak angka', () => {
    const html = renderToStaticMarkup(<MedanEkspor m={medan('QSPct', 'pxNumber')} k={k} />)
    expect(html).toContain('inputMode="decimal"')
    expect(html).toContain('value="80,00"')
  })

  it('desimal tak dinyatakan ekspor = tidak dibatasi (tarif tidak terpotong)', () => {
    const html = renderToStaticMarkup(<MedanEkspor m={medan('AdjRate', 'pxNumber', null)} k={k} />)
    expect(html).toContain('value="0,125"')
  })

  it('medan teks tetap teks bebas; nilai pxNumber yang bukan angka tampil apa adanya', () => {
    expect(renderToStaticMarkup(<MedanEkspor m={medan('Nama', 'pxTextInput')} k={k} />)).not.toContain('inputMode')
    const aneh = renderToStaticMarkup(<MedanEkspor m={medan('Aneh', 'pxNumber')} k={k} />)
    expect(aneh).not.toContain('inputMode')
    expect(aneh).toContain('value="IDR"')
  })

  it('huruf dan simbol yang diketik dibuang', () => {
    expect(formatKetik('12a3b,5x', 2)).toBe('123,5')
    expect(keKabelAngka(formatKetik('Rp 1.000abc', 2))).toBe('1000')
    expect(formatKetik('abc', 2)).toBe('')
  })
})
