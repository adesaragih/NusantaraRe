// Uji rantai rumus Section RINCIAN dan ketergantungan antartab panel New —
// 7 Oktober 2026.
//
// ⭐ Yang dijaga di sini:
//   - rantai berjalan atas AKAR + JALUR baris pemicu, dan hasilnya membawa
//     akar sebelum/sesudah (`awal`/`akhir`) untuk GABUNG TIGA ARAH;
//   - syarat aksi (`pyActionConditions`) dan parameter `.X` dibaca dari
//     halaman PEMICU (baris sel: `.Note`, `.Layer`);
//   - langkah bersyarat yang belum dibangun membatalkan rantai — nol
//     perubahan — hanya bila syaratnya terpenuhi;
//   - tambah baris Section rincian disalin dari Activity/DT korpus
//     `Treaty In Adjustment`.

import { afterEach, describe, expect, it, vi } from 'vitest'

import type { BarisBersarang, SisiPenyesuaian } from './api'
import type { AksiTombol, ButirKerangka, TombolKerangka } from './ekspor/jenis'
import { KERANGKA_INCLUDE, KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { PENYESUAIAN } from './labelsPenyesuaian'
import { rantaiSesudahHapus, rantaiSesudahTambah, tambahDari } from './komponen/aksiTombol'
import { gabungTigaArah } from './komponen/baris'
import { konteksBaris, type KonteksKerangka } from './komponen/KerangkaTab'
import { keSimpanTanggal, nilaiParam, PREFIX_HITUNG_TREATYIN, rantaiRumus, type HasilRantai } from './komponen/rumus'

type Baris = BarisBersarang

const sisi = (medan: Record<string, string>, larik: Record<string, Baris[]> = {}): SisiPenyesuaian => ({ medan, larik })

/** Aksi `change` sebuah medan (tanpa `larik`) atau sel grid `larik`. */
function aksiUbah(isi: readonly ButirKerangka[] | undefined, kunci: string, larik?: string): readonly AksiTombol[] {
  for (const b of isi ?? []) {
    if (b.t === 'blok') {
      const x = aksiUbah(b.anak, kunci, larik)
      if (x.length > 0) return x
    }
    if (larik === undefined && b.t === 'medan' && b.kunci === kunci && b.aksiUbah !== undefined) return b.aksiUbah
    if (larik !== undefined && b.t === 'grid' && b.larik === larik) {
      const a = b.aksiUbah[b.kunci.indexOf(kunci)]
      if (a !== null && a !== undefined) return a
    }
  }
  return []
}

/** Tombol grid `larik` (baris atau kepala) berlabel `label`. */
function tombolGrid(isi: readonly ButirKerangka[] | undefined, larik: string, label: string): TombolKerangka | undefined {
  for (const b of isi ?? []) {
    if (b.t === 'blok') {
      const x = tombolGrid(b.anak, larik, label)
      if (x !== undefined) return x
    }
    if (b.t === 'grid' && b.larik === larik) {
      const t = [...b.tombol, ...b.tombolKepala].find((x) => x !== null && x.label === label)
      if (t !== null && t !== undefined) return t
    }
  }
  return undefined
}

/** Tombol lepas (bukan grid) berlabel `label`. */
function tombol(isi: readonly ButirKerangka[] | undefined, label: string): TombolKerangka | undefined {
  for (const b of isi ?? []) {
    if (b.t === 'blok') {
      const x = tombol(b.anak, label)
      if (x !== undefined) return x
    }
    if (b.t === 'tombol' && b.label === label) return b
  }
  return undefined
}

const RIN = KERANGKA_RINCIAN
const TAB = (k: string) => KERANGKA_TAB[k]?.isi

/** `fetch` tiruan — satu jawaban per panggilan, BARU tiap kali (badan Response sekali baca). */
function jawab(...isi: unknown[]) {
  let i = 0
  return vi.spyOn(globalThis, 'fetch').mockImplementation(() => {
    const x = isi[Math.min(i, isi.length - 1)]
    i += 1
    return Promise.resolve(new Response(JSON.stringify(x), { status: 200 }))
  })
}

const badan = (sp: ReturnType<typeof jawab>, n = 0) => JSON.parse(String(sp.mock.calls[n]?.[1]?.body)) as Record<string, unknown>
const url = (sp: ReturnType<typeof jawab>, n = 0) => String(sp.mock.calls[n]?.[0])

afterEach(() => {
  vi.restoreAllMocks()
})

describe('gabung tiga arah — isian yang diketik selama rute menjawab TIDAK hilang', () => {
  it('yang rumus ubah menimpa; yang pemakai ketik (rumus tidak menyentuh) bertahan, juga di dalam baris', () => {
    const dasar = sisi({ A: '1', B: 'x' }, { L: [{ p: '1', q: 'a' }] })
    const kini = sisi({ A: '1', B: 'diketik' }, { L: [{ p: '1', q: 'diketik' }] })
    const hasil = sisi({ A: '2', B: 'x' }, { L: [{ p: '9', q: 'a' }] })
    expect(gabungTigaArah(dasar, kini, hasil)).toEqual(sisi({ A: '2', B: 'diketik' }, { L: [{ p: '9', q: 'diketik' }] }))
  })

  it('bentuk larik berubah (rumus menyusun ulang) → hasil rumus yang berlaku', () => {
    const dasar = sisi({}, { L: [{ p: '1' }] })
    const kini = sisi({}, { L: [{ p: '1', q: 'ketik' }] })
    const hasil = sisi({}, { L: [{ p: '1' }, { p: '2' }] })
    expect(gabungTigaArah(dasar, kini, hasil).larik.L).toEqual([{ p: '1' }, { p: '2' }])
  })

  it('larik yang rumus tidak ubah dibiarkan seperti SEKARANG', () => {
    const dasar = sisi({}, { L: [{ p: '1' }], M: [] })
    const kini = sisi({}, { L: [{ p: '1' }, { p: 'baru' }], M: [] })
    const hasil = sisi({}, { L: [{ p: '1' }], M: [{ x: '1' }] })
    expect(gabungTigaArah(dasar, kini, hasil).larik).toEqual({ L: [{ p: '1' }, { p: 'baru' }], M: [{ x: '1' }] })
  })
})

describe('konteks baris — hasil rantai berjalur diteruskan ke akar', () => {
  it('jalur diperpanjang per rincian; hasil ber-akar ke `terapkan` akar, bukan `ganti` baris', () => {
    const keAkar: unknown[] = []
    const diganti: Baris[] = []
    const akar: KonteksKerangka = {
      sisi: sisi({}),
      akar: sisi({}),
      halaman: {},
      jalur: [],
      terapkan: (h) => keAkar.push(h),
    }
    const k = konteksBaris(akar, { ID: '1' }, (b) => diganti.push(b), { larik: 'Limits', indeks: 0 })
    const dalam = konteksBaris(k, { ID: '2' }, () => undefined, { larik: 'Detail', indeks: 3 })
    expect(dalam.jalur).toEqual([
      { larik: 'Limits', indeks: 0 },
      { larik: 'Detail', indeks: 3 },
    ])
    const h: HasilRantai = { larik: {}, medan: {}, pesan: [], awal: sisi({}), akhir: sisi({ X: '1' }) }
    dalam.terapkan?.(h)
    expect(keAkar).toEqual([h])
    // Hasil tanpa akar (pemakaian lama) tetap mengganti barisnya.
    k.terapkan?.({ larik: {}, medan: { A: '1' }, pesan: [] })
    expect(diganti).toEqual([{ ID: '1', A: '1' }])
  })
})

describe('parameter dan tanggal', () => {
  it('`.X` dari halaman pemicu, `TreatyIn.X` dari akar, `.pxListSubscript` = sel/jalur + 1, kutip dibuang', () => {
    const s = sisi({ Layer: 'halaman' })
    const l = { panel: sisi({ RNMShare: '40' }), halaman: { '.Layer': 'true' }, sel: { larik: 'IOOLimitList', indeks: 2, kunci: 'Value' } }
    expect(nilaiParam('.Layer', s, l)).toBe('true')
    expect(nilaiParam('TreatyIn.RNMShare', s, l)).toBe('40')
    expect(nilaiParam('.pxListSubscript', s, l)).toBe('3')
    expect(nilaiParam('.pxListSubscript', s, { jalur: [{ larik: 'Share', indeks: 4 }] })).toBe('5')
    expect(nilaiParam('"reserve"', s, {})).toBe('reserve')
  })

  it('tanggal ke bentuk simpan: YYYY-MM-DD dan DD-MM-YYYY → YYYYMMDD; lainnya apa adanya', () => {
    expect(keSimpanTanggal('2024-01-18')).toBe('20240118')
    expect(keSimpanTanggal('18-01-2024')).toBe('20240118')
    expect(keSimpanTanggal('20240118')).toBe('20240118')
    expect(keSimpanTanggal(undefined)).toBe('')
  })
})

describe('Share Prop — rincian DetailShare/DetailLimits ke rute share-prop', () => {
  const akar = sisi(
    { RNMShareP: '10', Commencement: '01-01-2024' },
    { Limits: [{ ID: '1' }, { ID: '2', Detail: [{ RNMShare: '5', TreatyGroupID: '' }, { RNMShare: '1' }] }] },
  )

  it('% RNM Share Detail: halaman pemicu DITARUH di jalurnya; aksi `detail` ber-indeks jalur; hasil ke akar', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.DetailShare, 'RNMShare') })
    expect(r).toBeDefined()
    const Limits = [{ ID: '1' }, { ID: '2', Detail: [{ RNMShare: '7', TotalLimits: [{ Currency: 'IDR', Value: '1' }] }, { RNMShare: '1' }] }]
    const sp = jawab({ Limits, TotalShareRnmProp: [{ Currency: 'IDR', CurrencyID: '', Value: '9' }], TotalSpreadedRnmProp: [], TotalSpreadedRnmRIProp: [], pesan: [] })
    const halaman = sisi({ RNMShare: '7' }, {})
    const h = await r?.(halaman, { panel: akar, jalur: [{ larik: 'Limits', indeks: 1 }, { larik: 'Detail', indeks: 0 }], halaman: {} })
    expect(url(sp)).toContain(`${PREFIX_HITUNG_TREATYIN}/share-prop`)
    const b = badan(sp)
    expect(b).toMatchObject({ aksi: 'detail', indeksLimit: 1, indeksDetail: 0, RNMShareP: '10', Commencement: '20240101' })
    expect((b.Limits as Baris[])[1]).toMatchObject({ Detail: [{ RNMShare: '7' }, { RNMShare: '1' }] })
    expect(h?.awal.larik.Limits?.[1]).toMatchObject({ Detail: [{ RNMShare: '7' }, { RNMShare: '1' }] })
    expect(h?.akhir.larik.Limits).toEqual(Limits)
    expect(h?.akhir.larik.TotalShareRnmProp).toEqual([{ Currency: 'IDR', CurrencyID: '', Value: '9' }])
  })

  it('Spreading Type Detail: `spreading` membawa ParentReinsTypeID = .SpreadingTypeID halaman itu', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.DetailShare, 'SpreadingTypeID') })
    const sp = jawab({ Limits: akar.larik.Limits, TotalShareRnmProp: [], TotalSpreadedRnmProp: [], TotalSpreadedRnmRIProp: [], pesan: [] })
    await r?.(sisi({ SpreadingTypeID: '10241' }), {
      panel: akar,
      jalur: [{ larik: 'Limits', indeks: 1 }, { larik: 'Detail', indeks: 1 }],
      halaman: { '.SpreadingTypeID': '10241' },
    })
    expect(badan(sp)).toMatchObject({ aksi: 'spreading', ParentReinsTypeID: '10241', indeksLimit: 1, indeksDetail: 1 })
  })

  it('Treaty Group rincian Limits: nama dari master, ID = subscript (SetDetailsID), FetchQS TANPA induk bila QUOTA SHARE; LimitCalculation tidak (bukan SURPLUS)', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.DetailLimits, 'TreatyGroupID') })
    expect(r).toBeDefined()
    const sp = jawab({ Limits: akar.larik.Limits, TotalShareRnmProp: [], TotalSpreadedRnmProp: [], TotalSpreadedRnmRIProp: [], pesan: [] })
    await r?.(sisi({ TreatyGroupID: '7', TreatyType: 'QUOTA SHARE' }), {
      panel: akar,
      jalur: [{ larik: 'Limits', indeks: 1 }, { larik: 'Detail', indeks: 1 }],
      halaman: { '.TreatyType': 'QUOTA SHARE', '.TreatyGroupID': '7' },
      master: { jenisTreaty: [], mataUang: [], kelompokTreaty: [{ id: '7', nama: 'FIRE', kembar: false }] },
    })
    expect(sp).toHaveBeenCalledTimes(1)
    const b = badan(sp)
    expect(b).toMatchObject({ aksi: 'spreading', ParentReinsTypeID: '' })
    expect((b.Limits as Baris[])[1]).toMatchObject({ Detail: [{ RNMShare: '5' }, { TreatyGroup: 'FIRE', ID: '2' }] })
  })
})

describe('Limits Prop — LimitCalculation dari sel grid: syarat `.Note`, autocalculate `.Layer`', () => {
  const aksi = aksiUbah(RIN.DetailLimits, 'Value', 'IOOLimitList')
  const akar = sisi({}, { Limits: [{ ID: '1', Detail: [{ ID: '1', QSPct: '30', IOOLimitList: [{ Currency: 'IDR', Value: '100', Note: 'QUOTA SHARE', Layer: 'true' }] }] }] })
  const jalur = [{ larik: 'Limits', indeks: 0 }, { larik: 'Detail', indeks: 0 }]
  const halamanDetail = sisi({ ID: '1', QSPct: '30' }, { IOOLimitList: [{ Currency: 'IDR', Value: '100', Note: 'QUOTA SHARE', Layer: 'true' }] })
  const sel = { larik: 'IOOLimitList', indeks: 0, kunci: 'Value' }

  it('baris bukan QUOTA SHARE → aksi DILEWATI (nol panggilan); Layer tak dicentang → Activity KELUAR', async () => {
    const r = rantaiRumus({ aksi })
    const sp = jawab({})
    const a = await r?.(halamanDetail, { panel: akar, jalur, sel, halaman: { '.Note': 'SURPLUS', '.Layer': 'true' } })
    const b = await r?.(halamanDetail, { panel: akar, jalur, sel, halaman: { '.Note': 'QUOTA SHARE', '.Layer': 'false' } })
    expect(sp).not.toHaveBeenCalled()
    expect(a?.akhir).toEqual(a?.awal)
    expect(b?.akhir).toEqual(b?.awal)
  })

  it('QUOTA SHARE + Layer dicentang → POST /hitung/limit `qs`; HANYA medan yang Activity tulis digabung', async () => {
    const r = rantaiRumus({ aksi })
    const sp = jawab({ ID: '1', RetentionPct: '70', CessionPct: '30', RetentionList: [{ Currency: 'IDR', Value: '70' }], CessionList: [{ Currency: 'IDR', Value: '30' }], IOOPct: 'X' })
    const h = await r?.(halamanDetail, { panel: akar, jalur, sel, halaman: { '.Note': 'QUOTA SHARE', '.Layer': 'true' } })
    expect(url(sp)).toContain(`${PREFIX_HITUNG_TREATYIN}/limit`)
    expect(badan(sp)).toMatchObject({ jenis: 'qs', tambah: '', otomatis: true, detail: { QSPct: '30' } })
    expect((badan(sp).pohon as Baris[]).length).toBe(1)
    const d = (h?.akhir.larik.Limits?.[0]?.Detail as Baris[] | undefined)?.[0]
    expect(d).toMatchObject({ RetentionPct: '70', CessionPct: '30', QSPct: '30' })
    expect(d?.IOOPct).toBeUndefined()
  })

  it('Remove 100% Limit: syarat dibaca dari baris yang DIHAPUS', async () => {
    const t = tombolGrid(RIN.DetailLimits, 'IOOLimitList', 'Remove')
    expect(t).toBeDefined()
    const r = t === undefined ? undefined : rantaiSesudahHapus(t)
    expect(r).toBeTypeOf('function')
    const sp = jawab({ RetentionList: [], CessionList: [] })
    if (typeof r === 'function') {
      await r(sisi({ QSPct: '30' }, { IOOLimitList: [] }), { panel: akar, jalur, sel: { ...sel, kunci: '' }, halaman: { '.Note': 'QUOTA SHARE', '.Layer': 'true' } })
    }
    expect(sp).toHaveBeenCalledTimes(1)
  })

  it('mata uang sel: SetCurrName_Act mengisi PASANGAN pengenal dari master; LimitCalculation dilewati (bukan QUOTA SHARE)', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.DetailLimits, 'Currency', 'IOOLimitList') })
    expect(r).toBeDefined()
    const sp = jawab({})
    const h = await r?.(sisi({}, { IOOLimitList: [{ Currency: 'USD', CurrencyID: '10026', Note: 'SURPLUS' }] }), {
      panel: akar,
      jalur,
      sel: { larik: 'IOOLimitList', indeks: 0, kunci: 'Currency' },
      halaman: { '.Note': 'SURPLUS' },
      master: { jenisTreaty: [], kelompokTreaty: [], mataUang: [{ id: '10001', nama: 'USD', kembar: false }] },
    })
    expect(sp).not.toHaveBeenCalled()
    expect(h?.larik.IOOLimitList).toEqual([{ Currency: 'USD', CurrencyID: '10001', Note: 'SURPLUS' }])
  })

  it('Treaty Type rincian Limits: Kind of Treaty = SOA Name (Name bila kosong), ID = subscript, jenis turun ke Detail', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.LimitProportional, 'TreatyTypeID') })
    expect(r).toBeDefined()
    const h = await r?.(sisi({ TreatyTypeID: '10035' }, { Detail: [{ ID: '1' }] }), {
      panel: sisi({}, { Limits: [{}, { TreatyTypeID: '10035', Detail: [{ ID: '1' }] }] }),
      jalur: [{ larik: 'Limits', indeks: 1 }],
      halaman: {},
      master: { mataUang: [], kelompokTreaty: [], jenisTreaty: [{ id: '10035', nama: 'QUOTA SHARE', namaSoa: '', kembar: false }] },
    })
    expect(h?.medan).toMatchObject({ TreatyType: 'QUOTA SHARE', ID: '2' })
    expect(h?.larik.Detail).toEqual([{ ID: '1', TreatyType: 'QUOTA SHARE' }])
  })
})

describe('tambah baris Section rincian — disalin dari Activity/DT Adjustment', () => {
  it('AddDeduction: Currency = .Currency halaman (korpus Adjustment, BUKAN CurrencyIOOLimit)', () => {
    const t = tombolGrid(RIN.DetailLimits, 'DeductionList', 'Add')
    const tb = t === undefined ? undefined : tambahDari(t)
    expect(tb?.larik).toBe('DeductionList')
    expect(tb?.baris(sisi({ Currency: 'IDR', CurrencyIOOLimit: 'USD' }))).toEqual({ Currency: 'IDR' })
  })

  it('AddValue(type) dan AddLimitRetentionCession(type) — larik menurut parameternya', () => {
    const cek = (larik: string) => {
      const t = tombolGrid(RIN.DetailLimits, larik, 'Add')
      return t === undefined ? undefined : tambahDari(t)
    }
    expect(cek('ReserveList')?.baris(sisi({}))).toEqual({ ID: '' })
    expect(cek('EPIList')?.larik).toBe('EPIList')
    expect(cek('IOOLimitList')?.baris(sisi({ TreatyType: 'QUOTA SHARE' }))).toEqual({ Note: 'QUOTA SHARE' })
    expect(cek('RetentionList')?.larik).toBe('RetentionList')
    expect(cek('CessionList')?.larik).toBe('CessionList')
  })

  it('AddClassofBusiness: ParentID = subscript Kind of Treaty, TreatyType = jenis induk', () => {
    const t = tombolGrid(RIN.LimitProportional, 'Detail', 'Add')
    const tb = t === undefined ? undefined : tambahDari(t)
    expect(tb?.baris(sisi({ TreatyType: 'SURPLUS' }), 3)).toEqual({ ID: '', ParentID: '3', TreatyType: 'SURPLUS' })
  })

  it('Add Treaty Group layer: addRow lalu DT TreatyTypeSetIndex — ID layer = subscript', async () => {
    const t = tombolGrid(RIN.Layers, 'TreatyGroupList', 'Add Treaty Group')
    expect(t).toBeDefined()
    const r = t === undefined ? undefined : rantaiSesudahTambah(t)
    expect(r).toBeTypeOf('function')
    if (typeof r !== 'function') return
    const h = await r(sisi({ ID: '' }, { TreatyGroupList: [{}] }), { panel: sisi({}, { Limits: [{}, {}, { ID: '' }] }), jalur: [{ larik: 'Limits', indeks: 2 }] })
    expect(h.medan.ID).toBe('3')
  })
})

describe('Limits Non-Prop — rincian Layers / CoBList / DetailEGNPI', () => {
  const akar = sisi(
    {},
    {
      Limits: [
        { ID: '1', Layer: '1', AdjRate: '2', TreatyGroupList: [{ TreatyGroup: 'FIRE' }], Reinstatement_List: [{ AdditionalPct: '10' }] },
        { ID: '2', Layer: '2' },
      ],
      EGNPI: [{ TreatyGroup: 'FIRE', Currency: 'IDR', Amount: '100' }],
      CurrencyList: [{ Currency: 'USD', Conversion: '15000' }],
    },
  )

  it('Delete Treaty Group: deleteRow lalu TotalEgnpi — HANYA EgnpiTotalList seluruh layer yang diambil', async () => {
    const t = tombolGrid(RIN.Layers, 'TreatyGroupList', 'Delete')
    const r = t === undefined ? undefined : rantaiSesudahHapus(t)
    expect(r).toBeTypeOf('function')
    const sp = jawab({
      layers: [
        { ID: '1', AdjRate: 'diubah?', EgnpiTotalList: [{ Currency: 'IDR', Value: '0' }] },
        { ID: '2', EgnpiTotalList: [] },
      ],
      pesan: [],
    })
    if (typeof r !== 'function') return
    // Halaman pemicu = layer UTUH (baris yang dibuka), Treaty Group-nya sudah dihapus.
    const h = await r(sisi({ ID: '1', Layer: '1', AdjRate: '2' }, { TreatyGroupList: [], Reinstatement_List: [{ AdditionalPct: '10' }] }), { panel: akar, jalur: [{ larik: 'Limits', indeks: 0 }], sel: { larik: 'TreatyGroupList', indeks: 0, kunci: '' } })
    expect(badan(sp)).toMatchObject({ aksi: 'egnpi', indeks: 0 })
    expect(h.akhir.larik.Limits?.[0]).toMatchObject({ AdjRate: '2', EgnpiTotalList: [{ Currency: 'IDR', Value: '0' }], TreatyGroupList: [] })
  })

  it('Treaty Group CoBList: TotalEgnpi DUA kali di ekspor → SATU panggilan rute', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.CoBList, 'TreatyGroup') })
    expect(r).toBeDefined()
    const sp = jawab({ layers: [{}, {}], pesan: [] })
    await r?.(sisi({ TreatyGroup: 'MOTOR' }), { panel: akar, jalur: [{ larik: 'Limits', indeks: 0 }, { larik: 'TreatyGroupList', indeks: 0 }] })
    expect(sp).toHaveBeenCalledTimes(1)
  })

  it('Reinstatement % Additional: DT CalculateReinstatement → `reinst-jumlah` ber-indeks layer dan baris sel', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.Layers, 'AdditionalPct', 'Reinstatement_List') })
    expect(r).toBeDefined()
    const sp = jawab({ layers: [{ Reinstatement_List: [{ AdditionalPct: '10', ReinstatementAmount1: '5' }] }, {}], pesan: [] })
    const h = await r?.(sisi({ ID: '1' }, { Reinstatement_List: [{ AdditionalPct: '10' }] }), {
      panel: akar,
      jalur: [{ larik: 'Limits', indeks: 0 }],
      sel: { larik: 'Reinstatement_List', indeks: 0, kunci: 'AdditionalPct' },
    })
    expect(badan(sp)).toMatchObject({ aksi: 'reinst-jumlah', indeks: 0, baris: 0 })
    expect(h?.akhir.larik.Limits?.[0]?.Reinstatement_List).toEqual([{ AdditionalPct: '10', ReinstatementAmount1: '5' }])
  })

  it('EGNPI rincian Amount: SetAmountConversion → `konversi` baris itu; kunci kosong jawaban Go tidak dilahirkan', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.DetailEGNPI, 'Amount') })
    expect(r).toBeDefined()
    const sp = jawab({ egnpi: [{ TreatyGroup: 'FIRE', Currency: 'IDR', Amount: '200', AmountIDR: '200', ClassOfBusinessID: '' }], pesan: [] })
    const h = await r?.(sisi({ TreatyGroup: 'FIRE', Currency: 'IDR', Amount: '200' }), { panel: akar, jalur: [{ larik: 'EGNPI', indeks: 0 }] })
    expect(badan(sp)).toMatchObject({ aksi: 'konversi', indeks: 0 })
    expect(h?.akhir.larik.EGNPI).toEqual([{ TreatyGroup: 'FIRE', Currency: 'IDR', Amount: '200', AmountIDR: '200' }])
  })
})

describe('Share Non-Prop — syarat aksi dan langkah bersyarat yang belum dibangun', () => {
  const share = TAB('TreatyInTabsNonProportional#Share')
  const jawabShare = () => jawab({ share: { RNMShare: '40', Share: [], Total: {} }, pesan: [], pesanBaris: [] })

  it('Share Across The Board: XOLAddSpreading hanya bila dicentang, SetBrokerage hanya bila Brokerage > 0', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(share, 'RNMShareAcrossTheBoard') })
    expect(r).toBeDefined()
    let sp = jawabShare()
    await r?.(sisi({ RNMShareAcrossTheBoard: 'true', BrokeragePercent: '0' }), { jalur: [], halaman: { RNMShareAcrossTheBoard: 'true', BrokeragePercent: '0' } })
    expect(sp.mock.calls.map((_, i) => badan(sp, i).aksi)).toEqual(['rnm'])
    vi.restoreAllMocks()
    sp = jawabShare()
    await r?.(sisi({ RNMShareAcrossTheBoard: 'false', BrokeragePercent: '5' }), { jalur: [], halaman: { RNMShareAcrossTheBoard: 'false', BrokeragePercent: '5' } })
    expect(sp.mock.calls.map((_, i) => badan(sp, i).aksi)).toEqual(['set-brokerage'])
  })

  it('% Brokerage: SetBrokerage LALU XOLAddSpreading — dua langkah, urutan ekspor', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(share, 'BrokeragePercent') })
    const sp = jawabShare()
    await r?.(sisi({ BrokeragePercent: '5' }), { jalur: [], halaman: {} })
    expect(sp.mock.calls.map((_, i) => badan(sp, i).aksi)).toEqual(['set-brokerage', 'rnm'])
  })

  it('Update Summary: FacultativeShare kosong → 0 dulu; SATU POST `summary`; langkah EDM (7897987) dilewati', async () => {
    const t = tombol(share, 'Update Summary')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const sp = jawabShare()
    const h = await r?.(sisi({ FacultativeShare: '', RNMShare: '40' }), { jalur: [], halaman: { FacultativeShare: '', EDMMaterialType: '1' } })
    expect(sp).toHaveBeenCalledTimes(1)
    expect(badan(sp)).toMatchObject({ aksi: 'summary', share: { FacultativeShare: '0' } })
    expect(h?.pesan).toEqual([])
  })

  it('% RNM Share baris dari ActualValue.Share (EDMState = 3): XOLAddSpreadingDetail atas Share(idx) LALU …Actual atas ActualValue.Share(idx)', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.Share, 'RNMShare') })
    expect(r).toBeDefined()
    const sp = jawab(
      { share: { Share: [], Total: {} }, pesan: [], pesanBaris: [] },
      { actual: { medan: {}, larik: { Share: [{}, { RNMShare: '20', GrossPremiumList: [{ Currency: 'IDR', Value: '5' }] }] } }, selisih: null, akar: { medan: {}, larik: {} }, pesan: [] },
    )
    const akar = sisi({ EDMState: '3' }, { Share: [{}, {}], 'ActualValue.Share': [{}, { RNMShare: '20' }] })
    const h = await r?.(sisi({ RNMShare: '20' }), { panel: akar, jalur: [{ larik: 'ActualValue.Share', indeks: 1 }], halaman: { EDMState: '3' } })
    expect(url(sp, 0)).toContain('/share-np')
    expect(badan(sp, 0)).toMatchObject({ aksi: 'rnm-baris', indeks: 1 })
    expect(url(sp, 1)).toContain('/aktual')
    expect(badan(sp, 1)).toMatchObject({ aksi: 'rnm-baris', indeks: 1 })
    expect(h?.akhir.larik['ActualValue.Share']?.[1]?.GrossPremiumList).toEqual([{ Currency: 'IDR', Value: '5' }])
    // ⛔ 8 Oktober 2026 — catatan "langkah belum dibangun" dicabut dari
    // layar atas permintaan pemilik proses; yang belum dibangun tetap belum,
    // hanya tidak lagi diumumkan kepada pemakai.
    expect(PENYESUAIAN.langkahBelum).toBe('')
  })

  it('% RNM Share baris, EDMState lain: `rnm-baris` ber-indeks baris Share; hasil ke akar', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.Share, 'RNMShare') })
    const sp = jawab({ share: { Share: [{}, {}, { RNMShare: '20', GrossPremiumList: [{ Currency: 'IDR', Value: '1' }] }], Total: {} }, pesan: [], pesanBaris: [{ indeks: 2, pesan: 'baris 3' }, { indeks: 0, pesan: 'baris 1' }] })
    const akar = sisi({ EDMState: '2' }, { Share: [{}, {}, { RNMShare: '10' }] })
    const h = await r?.(sisi({ RNMShare: '20' }), { panel: akar, jalur: [{ larik: 'Share', indeks: 2 }], halaman: { EDMState: '2' } })
    expect(badan(sp)).toMatchObject({ aksi: 'rnm-baris', indeks: 2 })
    expect((badan(sp).share as { Share: Baris[] }).Share[2]).toMatchObject({ RNMShare: '20' })
    expect(h?.akhir.larik.Share?.[2]?.GrossPremiumList).toEqual([{ Currency: 'IDR', Value: '1' }])
    // Hanya pesan milik baris ini.
    expect(h?.pesan).toEqual(['baris 3'])
  })

  it('Deduction rincian Share: CalculateDeduction → share-np `deduksi`, baris = sel, sts dari parameter', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.Share, 'DeductionPct', 'DeductionList') })
    expect(r).toBeDefined()
    const sp = jawabShare()
    await r?.(sisi({}, { DeductionList: [{}, { DeductionPct: '5' }] }), {
      panel: sisi({}, { Share: [{ DeductionList: [{}, { DeductionPct: '5' }] }] }),
      jalur: [{ larik: 'Share', indeks: 0 }],
      sel: { larik: 'DeductionList', indeks: 1, kunci: 'DeductionPct' },
    })
    expect(url(sp)).toContain('/share-np')
    expect(badan(sp)).toMatchObject({ aksi: 'deduksi', indeks: 0, baris: 1, sts: 'pct' })
  })
})

describe('Prop akar — Accumulation, Reporting Period, Exclusions', () => {
  it('Period Accumulation: DT pra-refresh MENGOSONGKAN daftar lebih dulu → `periode` dengan daftar kosong; tanggal bentuk simpan', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(TAB('TreatyInTabsProportional#Accumulation'), 'AccumulationPeriod') })
    expect(r).toBeDefined()
    const sp = jawab({ AccumulationList: [{ Period: 'Q 1', ReportDate: '20240101', SubDays: '', SubDueDate: '' }] })
    const h = await r?.(sisi({ AccumulationPeriod: 'quarter', ReportingStart: '01-01-2024', ReportingEnd: '2024-12-31' }, { AccumulationList: [{ Period: 'X 1' }, { Period: 'X 2' }] }), { jalur: [], halaman: {} })
    expect(badan(sp)).toMatchObject({ aksi: 'periode', AccumulationList: [], ReportingStart: '20240101', ReportingEnd: '20241231' })
    expect(h?.larik.AccumulationList).toEqual([{ Period: 'Q 1', ReportDate: '20240101', SubDays: '', SubDueDate: '' }])
  })

  it('Sel Initial Date: hanya bila Auto Calculate baris dicentang; startdate = tanggal sel itu', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(TAB('TreatyInTabsProportional#Reporting Period'), 'InitialDate', 'ReportingPeriodList') })
    expect(r).toBeDefined()
    const sp = jawab({ baris: [], galat: {} })
    const akar = sisi({ ReportingStart: '20240101', ReportingEnd: '20241231', ReportingPeriod: 'quarter' }, { ReportingPeriodList: [{ InitialDate: '2024-04-01', AutoCalculate: 'false' }] })
    await r?.(akar, { jalur: [], sel: { larik: 'ReportingPeriodList', indeks: 0, kunci: 'InitialDate' }, halaman: { '.InitialDate': '2024-04-01', '.AutoCalculate': 'false' } })
    expect(sp).not.toHaveBeenCalled()
    await r?.(akar, { jalur: [], sel: { larik: 'ReportingPeriodList', indeks: 0, kunci: 'InitialDate' }, halaman: { '.InitialDate': '2024-04-01', '.AutoCalculate': 'true' } })
    expect(badan(sp)).toMatchObject({ awal: '20240401', mulai: '20240101' })
  })

  it('Exclusions: DT TreatyInCopyConditions — Exclusions/SpecialConditions disalin dari ejaan P, nol rute', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(TAB('TreatyInTabsProportional#Exclusions'), 'ExclusionsP') })
    expect(r).toBeDefined()
    const sp = jawab({})
    const h = await r?.(sisi({ ExclusionsP: 'perang', SpecialConditionsP: 'khusus' }), { jalur: [], halaman: {} })
    expect(sp).not.toHaveBeenCalled()
    expect(h?.medan).toEqual({ Exclusions: 'perang', SpecialConditions: 'khusus' })
  })

  it('Refresh Share Prop: `share` atas akar', async () => {
    const t = tombol(KERANGKA_INCLUDE.TreatyInShareProp, 'Refresh') ?? tombol(TAB('TreatyInTabsProportional#Share'), 'Refresh')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const sp = jawab({ Limits: [], TotalShareRnmProp: [], TotalSpreadedRnmProp: [], TotalSpreadedRnmRIProp: [], pesan: [] })
    await r?.(sisi({ RNMShareP: '10' }, { Limits: [] }), { jalur: [], halaman: {} })
    expect(badan(sp)).toMatchObject({ aksi: 'share', RNMShareP: '10' })
  })
})

describe('Value Difference, Achievement, Payment Date — rute Treaty In', () => {
  it('Update Value Value Difference: POST selisih (akar, OLDDATA, ActualValue); ValueDifference DIGANTI utuh berkunci bertitik', async () => {
    const t = tombol(KERANGKA_INCLUDE.TreatyInTabsNonProportionalValueDifference, 'Update Value')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const sp = jawab({
      selisih: { medan: { RNMShare: '10' }, larik: { Share: [{ Layer: '1' }] } },
      sebelumProrata: { medan: { RNMShare: '20' }, larik: {} },
      actual: { medan: {}, larik: { LimitShareSummaryList: [{ Note: 'x' }] } },
      FacultativeShareList: null,
      pesan: [],
    })
    const akar = sisi({ RNMShare: '50', 'ValueDifference.Lama': 'buang', 'ActualValue.RNMShare': '60' }, { 'ValueDifference.EGNPI': [{ Amount: '1' }] })
    const lama = sisi({ RNMShare: '40' })
    const h = await r?.(akar, { jalur: [], halaman: {}, lama })
    expect(url(sp)).toContain(`${PREFIX_HITUNG_TREATYIN}/selisih`)
    expect(badan(sp)).toMatchObject({ aksi: 'edm', akar: { medan: { RNMShare: '50' } }, lama: { medan: { RNMShare: '40' } }, actual: { medan: { RNMShare: '60' } } })
    expect(h?.medan).toMatchObject({ 'ValueDifference.RNMShare': '10', 'ValueDifference.Lama': '', 'ValueBeforeProrate.RNMShare': '20' })
    expect(h?.larik['ValueDifference.EGNPI']).toEqual([])
    expect(h?.larik['ValueDifference.Share']).toEqual([{ Layer: '1' }])
    expect(h?.larik['ActualValue.LimitShareSummaryList']).toEqual([{ Note: 'x' }])
  })

  it('Achievement: Refresh → GetAchievement atas seluruh Limits; daftar Quarter dan FlagExcel ke halaman sesi akar', async () => {
    const t = tombol(RIN.DetailLimits, 'Refresh')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const sp = jawab({ limits: [[{ AchievementLists: [{ Quarter: '1' }] }]], kuartal: ['1', '2'], tahunKuartal: ['2024'], flagExcel: true })
    const akar = sisi({ ID: '1000080', 'SearchData.CARI1': '1' }, { Limits: [{ TreatyType: 'QS', Detail: [{ TreatyGroup: 'FIRE' }] }] })
    const h = await r?.(sisi({ TreatyGroup: 'FIRE' }), { panel: akar, jalur: [{ larik: 'Limits', indeks: 0 }, { larik: 'Detail', indeks: 0 }], halaman: {} })
    expect(badan(sp)).toMatchObject({ idKontrak: '1000080', cari: false, asAt: '1', limits: [{ TreatyType: 'QS', Detail: [{ TreatyGroup: 'FIRE' }] }] })
    expect(h?.akhir.larik.TempQuarter).toEqual([{ CARI6: '1' }, { CARI6: '2' }])
    expect(h?.akhir.medan['FlagExcel.CARI1']).toBe('1')
    expect((h?.akhir.larik.Limits?.[0]?.Detail as BarisBersarang[] | undefined)?.[0]?.AchievementLists).toEqual([{ Quarter: '1' }])
  })

  it('Achievement: Quarter Year → GetAchievement(search); As At → DT Reset_DT (Quarter Year kosong)', async () => {
    const rTahun = rantaiRumus({ aksi: aksiUbah(RIN.DetailLimits, 'SearchData.CARI2') })
    const sp = jawab({ limits: [], kuartal: [], tahunKuartal: [], flagExcel: false })
    const jalurDetail = [{ larik: 'Limits', indeks: 0 }, { larik: 'Detail', indeks: 0 }]
    await rTahun?.(sisi({}), { panel: sisi({ 'SearchData.CARI2': '2024' }), jalur: jalurDetail, halaman: {} })
    expect(badan(sp)).toMatchObject({ cari: true, tahun: '2024' })
    const rAsAt = rantaiRumus({ aksi: aksiUbah(RIN.DetailLimits, 'SearchData.CARI1') })
    const h = await rAsAt?.(sisi({}), { panel: sisi({ 'SearchData.CARI2': '2024' }), jalur: jalurDetail, halaman: {} })
    expect(h?.akhir.medan['SearchData.CARI2']).toBe('')
  })

  it('Payment Date: sel Due Date → tanggal-bayar halaman itu, tanggal bentuk simpan; HANYA PaymentDate yang diambil', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.Installments, 'DueDate', 'InstallmentList') })
    expect(r).toBeDefined()
    const sp = jawab({ angsuran: [{ Currency: 'IDR', AmountTotal: '', PctTotal: '', InstallmentList: [{ PaymentDate: '20240319', Amount: 'X' }] }], TotalInstallmentNP: null, InstallmentNo: '', pesan: [] })
    const h = await r?.(sisi({ Currency: 'IDR' }, { InstallmentList: [{ DueDate: '2024-01-18', WPC: '60', Amount: '5' }] }), {
      panel: sisi({}, { Installment: [{ Currency: 'IDR' }] }),
      jalur: [{ larik: 'Installment', indeks: 0 }],
      sel: { larik: 'InstallmentList', indeks: 0, kunci: 'DueDate' },
      halaman: {},
    })
    expect(badan(sp)).toMatchObject({ aksi: 'tanggal-bayar', angsuran: [{ InstallmentList: [{ DueDate: '20240118', WPC: '60' }] }] })
    expect(h?.larik.InstallmentList).toEqual([{ DueDate: '2024-01-18', WPC: '60', Amount: '5', PaymentDate: '20240319' }])
  })

  it('Payment Date: sel WPC → tanggal-bayar-semua atas SEMUA halaman Installment akar', async () => {
    const r = rantaiRumus({ aksi: aksiUbah(RIN.Installments, 'WPC', 'InstallmentList') })
    const sp = jawab({ angsuran: [{ Currency: 'IDR', InstallmentList: [{ PaymentDate: 'a' }] }, { Currency: 'USD', InstallmentList: [{ PaymentDate: 'b' }] }], TotalInstallmentNP: null, InstallmentNo: '', pesan: [] })
    const akar = sisi({}, { Installment: [{ Currency: 'IDR', InstallmentList: [{ DueDate: '20240101' }] }, { Currency: 'USD', InstallmentList: [{ DueDate: '20240101' }] }] })
    const h = await r?.(sisi({ Currency: 'USD' }, { InstallmentList: [{ DueDate: '20240101', WPC: '3' }] }), { panel: akar, jalur: [{ larik: 'Installment', indeks: 1 }], halaman: {} })
    expect(badan(sp)).toMatchObject({ aksi: 'tanggal-bayar-semua' })
    expect((h?.akhir.larik.Installment?.[0]?.InstallmentList as BarisBersarang[] | undefined)?.[0]?.PaymentDate).toBe('a')
    expect((h?.akhir.larik.Installment?.[1]?.InstallmentList as BarisBersarang[] | undefined)?.[0]?.PaymentDate).toBe('b')
  })
})

describe('cabang Adjust Premium (EDMState = 3)', () => {
  const ap = (t: string) => TAB(`TreatyInTabsNonProportionalAdjustPremi#${t}`)

  it('Add Actual GNPI: TreatyInNonAddItem(actualpremium) — ID kosong, mata uang dari Retention(1)', () => {
    let tb: ReturnType<typeof tambahDari>
    const cari = (bs: readonly ButirKerangka[] | undefined) => {
      for (const b of bs ?? []) {
        if (b.t === 'blok') cari(b.anak)
        if (b.t === 'grid' && b.larik === 'ActualValue.EGNPI') for (const t of [...b.tombol, ...b.tombolKepala]) if (t !== null && tambahDari(t) !== undefined) tb = tambahDari(t)
      }
    }
    cari(ap('Actual GNPI'))
    expect(tb?.larik).toBe('ActualValue.EGNPI')
    expect(tb?.baris(sisi({}, { Retention: [{ Currency: 'USD', CurrencyID: '10001' }] }))).toEqual({ ID: '', Currency: 'USD', CurrencyID: '10001' })
  })

  it('Actual GNPI · Update Total: POST aktual nilai; ActualValue DIGANTI utuh, akar hanya kunci yang ditulis', async () => {
    const t = tombol(ap('Actual GNPI'), 'Update Total')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const sp = jawab({
      actual: { medan: { TotalEgnpiAmount: '500' }, larik: { EGNPI: [{ Amount: '500' }] } },
      selisih: { medan: { RNMShare: '0' }, larik: {} },
      akar: { medan: { RnmShareDeducted: '30' }, larik: {} },
      pesan: [],
    })
    const akar = sisi({ RNMShare: '40', 'ActualValue.Lama': 'x' }, { Limits: [{ Layer: '1' }], 'ActualValue.EGNPI': [{ Amount: '500' }] })
    const h = await r?.(akar, { jalur: [], halaman: {} })
    expect(badan(sp)).toMatchObject({ aksi: 'nilai', akar: { medan: { RNMShare: '40' }, larik: { Limits: [{ Layer: '1' }] } }, actual: { larik: { EGNPI: [{ Amount: '500' }] } } })
    expect(h?.medan).toMatchObject({ 'ActualValue.TotalEgnpiAmount': '500', 'ActualValue.Lama': '', RnmShareDeducted: '30', 'ValueDifference.RNMShare': '0' })
  })

  it('Actual Limits · Update Total: SATU POST limits (kedua Summary Actual tercakup)', async () => {
    const t = tombol(KERANGKA_INCLUDE.TreatyInActualLimits, 'Update Total')
    const r = t === undefined ? undefined : rantaiRumus(t)
    const sp = jawab({ actual: { medan: {}, larik: {} }, selisih: null, akar: { medan: {}, larik: {} }, pesan: [] })
    await r?.(sisi({}), { jalur: [], halaman: {} })
    expect(sp).toHaveBeenCalledTimes(1)
    expect(badan(sp)).toMatchObject({ aksi: 'limits' })
  })
})
