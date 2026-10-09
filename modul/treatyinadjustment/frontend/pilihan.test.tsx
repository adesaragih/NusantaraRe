// Uji dropdown/autocomplete panel New — sumber daftar dari ekspor.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { ButirKerangka, SumberPilihan } from './ekspor/jenis'
import { GRID_KURS, KERANGKA_INCLUDE, KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { GridEkspor, type KonteksKerangka } from './komponen/KerangkaTab'
import { denganNilaiKini, opsiUntuk, type DataOpsi } from './komponen/pilihan'

const DATA: DataOpsi = {
  limits: {
    jenisTreaty: [{ id: '10042', nama: 'QUOTA SHARE', namaSoa: 'QS' }],
    kelompokTreaty: [{ id: '7', nama: 'PROPERTY' }],
    mataUang: [
      { id: '10026', nama: 'IDR' },
      { id: '10001', nama: 'USD' },
    ],
  },
  kepala: {
    bordereaux: [{ value: 'reporting', label: 'Reporting' }],
    caraPembukuan: [],
    caraPembukuanNonProp: [],
    periodePelaporan: [{ value: 'quarter', label: 'Quarterly' }],
  },
  agen: [{ id: '1', nama: 'PT ASURANSI A' }],
  halaman: { TempQuarter: [{ CARI6: '1' }], TempQuarterYear: [{ CARI6: '2024' }] },
}

/** Setiap sumber pilihan di kerangka, dengan kuncinya. */
function semuaSumber(): { kunci: string; sp: SumberPilihan }[] {
  const out: { kunci: string; sp: SumberPilihan }[] = []
  const jalan = (bs: readonly ButirKerangka[]) => {
    for (const b of bs) {
      if (b.t === 'blok') jalan(b.anak)
      if (b.t === 'medan' && b.pilihan !== undefined) out.push({ kunci: b.kunci, sp: b.pilihan })
      if (b.t === 'grid') b.pilihan.forEach((sp, i) => sp !== null && out.push({ kunci: b.kunci[i] ?? '', sp }))
    }
  }
  for (const k of Object.values(KERANGKA_TAB)) jalan(k.isi)
  for (const v of Object.values(KERANGKA_INCLUDE)) jalan(v)
  for (const v of Object.values(KERANGKA_RINCIAN)) jalan(v)
  jalan(Object.values(GRID_KURS))
  return out
}

describe('opsiUntuk — satu entri per sumber ekspor', () => {
  it('BrowseCurrency_RD bernilai `.ID` (grid kurs New): nilai pengenal, tampil kode', () => {
    // `id` ikut — pasangan kode yang `pySetValueOnSelect` isi.
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseCurrency_RD', nilai: 'ID', tampil: 'Currency' }, 'CurrencyID', DATA)).toEqual([
      { value: '10026', label: 'IDR', id: '10026' },
      { value: '10001', label: 'USD', id: '10001' },
    ])
  })

  it('RD bernilai nama: Treaty Group, Currency Event Limits, Treaty Type (`.Note`), agen', () => {
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseTreatyGroup_RD', nilai: 'TreatyGroupName' }, 'TreatyGroup', DATA)).toMatchObject([{ value: 'PROPERTY', label: 'PROPERTY' }])
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseCurrencyTreatyIn_RD', nilai: 'Currency' }, 'CurrencyRSMD', DATA)?.map((o) => o.value)).toEqual(['IDR', 'USD'])
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseReinsuranceType_RD', nilai: 'Note' }, 'TreatyType', DATA)).toMatchObject([{ value: 'QUOTA SHARE', label: 'QUOTA SHARE' }])
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseAgentNusaRe_RD', nilai: 'ClientName' }, 'ReinsName', DATA)).toMatchObject([{ value: 'PT ASURANSI A', label: 'PT ASURANSI A' }])
  })

  it('⛔ RD bernilai `.ID` (Treaty Group / Treaty Type rincian Limits Prop): nilai PENGENAL, tampil nama', () => {
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseTreatyGroup_RD', nilai: 'ID', tampil: 'TreatyGroupName' }, 'TreatyGroupID', DATA)).toEqual(
      DATA.limits?.kelompokTreaty.map((x) => ({ value: x.id, label: x.nama, id: x.id })),
    )
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseReinsuranceType_RD', nilai: 'ID', tampil: 'Note' }, 'TreatyTypeID', DATA)).toEqual(
      DATA.limits?.jenisTreaty.map((x) => ({ value: x.id, label: x.nama, id: x.id })),
    )
  })

  it('`associated` hanya yang berbukti; daftar belum dimuat → `undefined` (kotak teks)', () => {
    expect(opsiUntuk({ sumber: 'associated' }, 'ReportingPeriod', DATA)).toEqual([{ value: 'quarter', label: 'Quarterly' }])
    expect(opsiUntuk({ sumber: 'associated' }, 'AccumulationPeriod', DATA)?.map((o) => o.value)).toEqual(['quarter', 'half', 'month', 'none'])
    expect(opsiUntuk({ sumber: 'associated' }, 'OptionLimit', DATA)).toEqual([
      { value: '1', label: 'Of Cession to R/I' },
      // *(8 Okt, E)* rule `OptionLimit` (prompt list) — sama dengan Treaty In.
      { value: '2', label: 'Of 100% Limit' },
    ])
    expect(opsiUntuk({ sumber: 'reportdefinition', rd: 'BrowseTreatyGroup_RD' }, 'TreatyGroup', {})).toBeUndefined()
  })

  it('⭐ pageList halaman SESI Achievement: TempQuarter / TempQuarterYear, nilai `.CARI6`', () => {
    expect(opsiUntuk({ sumber: 'pageList', halaman: 'TempQuarterYear.pxResults' }, 'SearchData.CARI2', DATA)).toEqual([{ value: '2024', label: '2024' }])
    expect(opsiUntuk({ sumber: 'pageList', halaman: 'TempQuarter.pxResults' }, 'SearchData.CARI1', {})).toBeUndefined()
  })

  it('⛔ sumber kerangka yang BELUM punya daftar — daftar persis, supaya terlihat', () => {
    const tanpa = semuaSumber()
      .filter((x) => opsiUntuk(x.sp, x.kunci, DATA) === undefined)
      .map((x) => `${x.kunci} ← ${x.sp.rd ?? x.sp.sumber}`)
    // Semuanya dari Section RINCIAN Layers / Share / DetailLimits (7 Oktober
    // 2026) — medannya kotak teks sampai daftarnya disambungkan.
    expect([...new Set(tanpa)].sort()).toEqual([
      'ClassOfBusiness ← BrowseTreatyBusinessWOType_RD',
      'Cover ← associated',
      'CurrencyRelation ← associated',
      'LayerPartType ← associated',
      'LayerType ← associated',
      'ReinsTypeID ← BrowseTreatyArrangement_ParentReinsMasterTrt',
      'ReinstatementNote ← associated',
      'SpreadingTypeID ← BrowseTreatyArrangement_ParentReinsMasterTrt',
      'SpreadingTypeXOL ← BrowseTreatyArrangement_ParentReinsMasterTrt',
      'TreatyGroup ← BrowseBusinessGroup_RD',
      'TreatyGroup ← BrowseBusiness_RD',
    ])
  })

  it('nilai tersimpan di luar daftar DITAWARKAN (ditandai), bukan dijatuhkan', () => {
    expect(denganNilaiKini([{ value: 'IDR', label: 'IDR' }], 'SGD', '(x)')).toEqual([
      { value: 'SGD', label: 'SGD (x)' },
      { value: 'IDR', label: 'IDR' },
    ])
    expect(denganNilaiKini([{ value: 'IDR', label: 'IDR' }], 'IDR', '(x)')).toHaveLength(1)
  })
})

describe('render — grid kurs panel New di mode Edit', () => {
  const sisi = { medan: { ViewState: '0' }, larik: { CurrencyList: [{ CurrencyID: '10001', Conversion: '15000' }] } }
  const konteks = (opsi?: KonteksKerangka['opsi']): KonteksKerangka => ({
    sisi,
    akar: sisi,
    halaman: { ViewState: '0' },
    ubah: true,
    opsi,
  })

  // ⛔ DIBALIK 8 Oktober 2026 — dahulu `<select>`. Permintaan pemilik
  // proses untuk SELURUH dropdown Treaty In dan Adjustment: *"kondisi saat
  // melakukan pengetikannya seharusnya terlihat layaknya melakukan mencari,
  // kemudian data yang keluar adalah yang 100% mirip dengan yang diketik"*.
  //
  // Kontrolnya kini combobox (`PilihSaring`): satu kotak yang dapat diketik,
  // daftarnya tersaring saat mengetik, dan ketikan tanpa memilih
  // dikembalikan ke pilihan terakhir.
  it('Currency menjadi combobox berisi daftar BrowseCurrency_RD, nilai tersimpan terpilih', () => {
    const html = renderToStaticMarkup(<GridEkspor g={GRID_KURS.baru} k={konteks((sp, kunci) => opsiUntuk(sp, kunci, DATA))} />)
    expect(html).toMatch(/role="combobox"[^>]*aria-label="Currency"|aria-label="Currency"[^>]*role="combobox"/)
    // ⚠️ Teks pilihan terbaca di KOTAKNYA. Daftar `<li>` baru dirender
    // saat combobox dibuka, jadi render statis tidak memuatnya.
    expect(html).toContain('value="USD"')
    expect(html).toContain('aria-expanded="false"')
    // ⛔ Sel tabel: `<th>` sudah menamainya, jadi NOL label di dalam sel.
    expect(html).not.toContain('<label class="field__label">Currency</label>')
  })

  it('daftar belum dimuat → kotak teks, bukan combobox kosong', () => {
    const html = renderToStaticMarkup(<GridEkspor g={GRID_KURS.baru} k={konteks()} />)
    expect(html).not.toContain('role="combobox"')
    expect(html).not.toContain('<select')
  })
})
