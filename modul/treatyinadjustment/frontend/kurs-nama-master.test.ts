// Uji `SetCurrNameMasterTreaty_Act` — sel Currency grid Rate of Exchange
// panel New. *(8 Okt, E)*
//
// Bukti: `Treaty In Adjustment/Activity/SetCurrNameMasterTreaty_Act.xml`
// [2] Obj-Browse CURRENCY `.ID = .CurrencyID` (@4723/@4763), [3]
// Property-Set `.Currency = CURRENCY.pxResults(1).Currency` (@18421/@18463);
// aksi sel dari `GRID_KURS.baru.aksiUbah[0]` (kerangka bangkitan).

import { describe, expect, it } from 'vitest'

import type { SisiPenyesuaian } from './api'
import { GRID_KURS } from './ekspor/kerangka.gen'
import { rantaiRumus } from './komponen/rumus'

const MASTER = {
  jenisTreaty: [],
  kelompokTreaty: [],
  mataUang: [
    { id: '10026', nama: 'IDR' },
    { id: '10001', nama: 'USD' },
  ],
}

const sisi = (CurrencyList: Record<string, string>[]): SisiPenyesuaian => ({ medan: {}, larik: { CurrencyList } })

describe('SetCurrNameMasterTreaty_Act — grid kurs panel New', () => {
  const aksi = GRID_KURS.baru.aksiUbah[GRID_KURS.baru.kunci.indexOf('CurrencyID')] ?? []

  it('aksi sel CurrencyID = postValue → SetCurrNameMasterTreaty_Act → refresh (ekspor)', () => {
    expect(aksi.map((a) => a.aktivitas ?? a.aksi)).toEqual(['postValue', 'SetCurrNameMasterTreaty_Act', 'refresh'])
  })

  it('rantainya kini BERUMUS (dahulu mati: Activity tanpa rumus)', () => {
    expect(rantaiRumus({ aksi })).toBeTypeOf('function')
  })

  it('mengisi `.Currency` dari `.CurrencyID` baris PEMICU saja', async () => {
    const r = rantaiRumus({ aksi })
    const awal = sisi([
      { CurrencyID: '10001', Currency: 'IDR', Conversion: '15000' },
      { CurrencyID: '10026', Currency: 'lama' },
    ])
    const h = await r?.(awal, { sel: { larik: 'CurrencyList', indeks: 0, kunci: 'CurrencyID' }, halaman: {}, master: MASTER })
    expect(h?.larik.CurrencyList).toEqual([
      { CurrencyID: '10001', Currency: 'USD', Conversion: '15000' },
      { CurrencyID: '10026', Currency: 'lama' },
    ])
  })

  it('pengenal tak dikenal → `.Currency` kosong (`pxResults(1)` tidak ada)', async () => {
    const r = rantaiRumus({ aksi })
    const h = await r?.(sisi([{ CurrencyID: '99999', Currency: 'USD' }]), {
      sel: { larik: 'CurrencyList', indeks: 0, kunci: 'CurrencyID' },
      halaman: {},
      master: MASTER,
    })
    expect(h?.larik.CurrencyList).toEqual([{ CurrencyID: '99999', Currency: '' }])
  })

  it('daftar master belum dimuat → nol perubahan, ada pesan', async () => {
    const r = rantaiRumus({ aksi })
    const h = await r?.(sisi([{ CurrencyID: '10001', Currency: '' }]), {
      sel: { larik: 'CurrencyList', indeks: 0, kunci: 'CurrencyID' },
      halaman: {},
    })
    expect(h?.larik.CurrencyList).toBeUndefined()
    expect(h?.pesan.length).toBeGreaterThan(0)
  })
})
