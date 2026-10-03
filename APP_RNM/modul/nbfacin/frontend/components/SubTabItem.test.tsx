// Paritas sub-tab Object Item (tiket 39) - section Pega `PropertyItemList` / `PropertyItemFacIn_Section`. Data uji
// sintetis. Uang = teks desimal; penjumlahan diuji EKSAK (tanpa float).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { FORM_ITEM as F, GRID_ITEM, TEKS_INWARD, TEKS_ITEM } from '../labels'
import SubTabItem, { galatItem, itemBaru, pilihJenis, setelAdjustable, tandaDesimal, totalPerMataUang } from './SubTabItem'

const SUMBER = readFileSync(join(__dirname, 'SubTabItem.tsx'), 'utf8').replace(/\r\n/g, '\n')
const it1 = (ubah: Partial<ReturnType<typeof itemBaru>>) => ({ ...itemBaru(), ...ubah })

describe('SubTabItem - tampilan', () => {
  it('grid kosong: kolom berurutan + Tambah, "No items"; grid Total ikut kosong', () => {
    const html = renderToStaticMarkup(<SubTabItem items={[]} ubah={() => {}} />)
    const kolom = [...html.matchAll(/<th[^>]*>([\s\S]*?)<\/th>|<th[^>]*\/>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(kolom).toEqual(['', ...GRID_ITEM.kolom.map((k) => k.label), 'Tambah', ...GRID_ITEM.total.map((k) => k.label)])
    expect(html.split(TEKS_INWARD.kosong)).toHaveLength(3)
  })

  it('baris: TSI 4 desimal gaya Indonesia; Total per mata uang', () => {
    const html = renderToStaticMarkup(
      <SubTabItem
        items={[it1({ itemType: 'UJI A', currency: 'USD', tsi: '1234567.5' }), it1({ itemType: 'UJI B', currency: 'USD', tsi: '0.25' })]}
        ubah={() => {}}
      />,
    )
    expect(html).toContain('<td class="nbf-angka">1.234.567,5</td>')
    expect(html).toContain('<td>USD</td><td class="nbf-angka">1.234.567,75</td>')
  })
})

describe('SubTabItem - aturan', () => {
  it('total per mata uang EKSAK (bukan float), hanya total > 0, urutan muncul', () => {
    const t = totalPerMataUang([
      it1({ currency: 'IDR', tsi: '0.1' }),
      it1({ currency: 'USD', tsi: '99999999999999999999.12345678' }),
      it1({ currency: 'IDR', tsi: '0.2' }),
      it1({ currency: 'USD', tsi: '0.00000002' }),
      it1({ currency: 'EUR', tsi: '0' }),
      it1({ currency: '', tsi: '5' }),
    ])
    expect(t).toEqual([
      { currency: 'IDR', total: '0.3' },
      { currency: 'USD', total: '99999999999999999999.12345680' },
    ])
  })

  it('memilih Object Item Type mengisi nama dan Note (SetObjItemType_Act)', () => {
    const i = pilihJenis(itemBaru(), { kode: 'K1', nama: 'UJI JENIS', keterangan: 'UJI KET' }, 'K1')
    expect([i.itemTypeId, i.itemType, i.note]).toEqual(['K1', 'UJI JENIS', 'UJI KET'])
    expect(pilihJenis(i, undefined, '').itemType).toBe('')
  })

  it('Adjustable dimatikan -> PctAdjustOther 100 (ResetPct_Adjustment); dinyalakan -> nilai tetap', () => {
    expect(setelAdjustable(it1({ isAdjustable: true, pctAdjustOther: '70' }), false).pctAdjustOther).toBe('100')
    expect(setelAdjustable(it1({ pctAdjustOther: '70' }), true).pctAdjustOther).toBe('70')
  })

  it('Currency wajib (A133): kosong = galat (menahan Save); terisi = sah', () => {
    expect(galatItem(itemBaru()).currency).toBe(TEKS_ITEM.currencyWajib)
    expect(galatItem(it1({ currency: 'USD' })).currency).toBeUndefined()
  })

  it('galat: Unit <= 0, TSI bukan angka / minus, % adjustment di luar 60..100 (hanya bila Adjustable)', () => {
    expect(galatItem(it1({ unit: '0' })).unit).toBe(TEKS_ITEM.unit)
    expect(galatItem(it1({ unit: '3' })).unit).toBeUndefined()
    expect(galatItem(it1({ tsi: '1,5' })).tsi).toBe(TEKS_ITEM.tsiBukanAngka)
    expect(galatItem(it1({ tsi: '-0.01' })).tsi).toBe(TEKS_ITEM.tsiMinus)
    expect(galatItem(it1({ tsi: '-0.00' })).tsi).toBeUndefined()
    expect(galatItem(it1({ isAdjustable: true, pctAdjustOther: '59.99' })).pct).toBe(TEKS_ITEM.pctAdjust)
    expect(galatItem(it1({ isAdjustable: true, pctAdjustOther: '100.01' })).pct).toBe(TEKS_ITEM.pctAdjust)
    expect(galatItem(it1({ isAdjustable: true, pctAdjustOther: '60' })).pct).toBeUndefined()
    expect(galatItem(it1({ isAdjustable: true, pctAdjustOther: '100' })).pct).toBeUndefined()
    expect(galatItem(it1({ isAdjustable: false, pctAdjustOther: '5' })).pct).toBeUndefined()
    expect(itemBaru().pctAdjust2).toBe('100')
  })

  it('tandaDesimal eksak', () => {
    expect([tandaDesimal('0.000'), tandaDesimal('-0.0001'), tandaDesimal('0.0001')]).toEqual([0, -1, 1])
  })

  it('form: label "Remark of <jenis>", Adjustment Pct. dropdown bila tidak Adjustable; tanpa Number()/parseFloat', () => {
    expect(SUMBER).toContain('`${F.remark.label} ${i.itemType}`.trim()')
    expect(SUMBER).toMatch(/i\.isAdjustable \? \(\s*<Field label=\{F\.pctAdjust\.label\}[\s\S]*?: \(\s*<Pilih label=\{F\.pctAdjust\.label\}/)
    // `Number(` telanjang (bukan `formatNumber(`) = jalan float yang dilarang ADR-0003.
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
    expect(F.tsi.label).toBe('TSI Object Item (All Unit)')
  })
})

describe('SubTabItem - Condition (Condition.xml)', () => {
  it('kolom grid Condition menampilkan label, bukan kode', () => {
    const html = renderToStaticMarkup(<SubTabItem items={[it1({ itemType: 'UJI', condition: '2' })]} ubah={() => {}} />)
    expect(html).toContain('<td>UJI</td><td>Fair</td>')
  })
})
