// Uji aksi tombol panel New — tambah baris dan rantai rumus.
//
// ⛔ Klasifikasi SETIAP tombol kerangka dijaga di sini: tombol yang belum
// punya rumus tercantum PERSIS. Menyambungkan rumus baru mengubah daftar
// ini, dan itu disengaja — perubahannya terlihat di uji, bukan diam.

import { afterEach, describe, expect, it, vi } from 'vitest'

import type { BarisBersarang, SisiPenyesuaian } from './api'
import type { AksiTombol, ButirKerangka, TombolKerangka } from './ekspor/jenis'
import { KERANGKA_INCLUDE, KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { hapusBarisGrid, labelTombol, rantaiSesudahHapus, tambahBarisGrid, tambahDari, terkunci, tulisanDari, unduhanDari } from './komponen/aksiTombol'
import { gabungPohon, nilaiJalur } from './komponen/baris'
import { tambah } from './komponen/desimal'
import { PREFIX_HITUNG_TREATYIN, rantaiRumus, ringkasanMDP } from './komponen/rumus'

/** Setiap tombol New (bukan Old) beserta tab asalnya. */
function tombolBaru(): { tab: string; t: TombolKerangka }[] {
  const out: { tab: string; t: TombolKerangka }[] = []
  const jalan = (tab: string, bs: readonly ButirKerangka[]) => {
    for (const b of bs) {
      if (b.t === 'tombol') out.push({ tab, t: b })
      if (b.t === 'blok') jalan(tab, b.anak)
      if (b.t === 'grid') for (const t of [...b.tombol, ...b.tombolKepala]) if (t !== null) out.push({ tab, t })
    }
  }
  for (const [k, v] of Object.entries(KERANGKA_TAB)) if (!k.includes('OldData')) jalan(k, v.isi)
  for (const [k, v] of Object.entries(KERANGKA_INCLUDE)) if (!k.includes('OldData')) jalan(k, v)
  // Section rincian baris panel New (Section Old: `…OldData`, `…_ReadOnly`, `CoBListReadOnly`).
  for (const [k, v] of Object.entries(KERANGKA_RINCIAN)) if (!/OldData|ReadOnly/.test(k)) jalan(k, v)
  return out
}

const golongan = (t: TombolKerangka) =>
  tambahDari(t) !== undefined || tambahBarisGrid(t) || (hapusBarisGrid(t) && rantaiSesudahHapus(t) !== undefined)
    ? 'baris'
    : unduhanDari(t) !== undefined
      ? 'unduh'
      : tulisanDari(t) !== undefined
        ? 'tulis'
        : rantaiRumus(t) !== undefined
          ? 'rumus'
          : 'mati'

const sisi = (medan: Record<string, string>, larik: Record<string, BarisBersarang[]> = {}): SisiPenyesuaian => ({ medan, larik })

describe('klasifikasi tombol panel New', () => {
  it('⭐ tombol berumus yang SUDAH tersambung', () => {
    const rumus = tombolBaru()
      .filter((x) => golongan(x.t) === 'rumus')
      .map((x) => `${x.tab.split('#')[1] ?? x.tab} · ${labelTombol(x.t)}`)
    expect([...new Set(rumus)].sort()).toEqual([
      // ⭐ 7 Oktober 2026 (lanjutan) — cabang Adjust Premium (rute `aktual`),
      // Achievement Refresh (`GetAchievement`), Value Difference (`selisih`).
      'Actual GNPI · Update Total',
      'DetailLimits · Refresh',
      'EGNPI · Update EGNPI Value',
      'EGNPI · Update Total',
      'Installment · Update Total',
      'Installment · Update Value',
      'Limits · Update Total',
      'Limits · Update Value in List',
      'Maximum Retention · Update Total',
      'Reporting Period · Apply',
      // ⭐ 7 Oktober 2026 — TreatyInPropshare (Share Prop) dan
      // TreatyInNonAddItem(share) + XOLAddSpreading + SetBrokerage (NP).
      'Share · Refresh',
      'Share · Update Summary',
      'Share · Update Total',
      'Share · Update Value in Share',
      'TreatyInActualLimits · Update Total',
      'TreatyInActualShare · Update Summary',
      'TreatyInActualSumary · Update Total',
      'TreatyInTabsNonProportionalValueDifference · Update Value',
    ])
  })

  it('⭐ Generate Excel (`GenerateCSVTreaty`) MENGUNDUH — nol perubahan keadaan', () => {
    const unduh = tombolBaru()
      .filter((x) => golongan(x.t) === 'unduh')
      .map((x) => `${x.tab.split('#')[1] ?? x.tab} · ${labelTombol(x.t)}`)
    expect([...new Set(unduh)]).toEqual(['DetailLimits · Generate Excel'])
  })

  it('⛔ tombol yang MASIH mati — daftar persis, supaya yang tersisa terlihat', () => {
    const mati = tombolBaru()
      .filter((x) => golongan(x.t) === 'mati')
      .map((x) => `${x.tab.split('#')[1] ?? x.tab} · ${labelTombol(x.t)}`)
    expect([...new Set(mati)].sort()).toEqual([
      // ⭐ 7 Oktober 2026 — Add*/Remove rincian Limits Prop, Delete Treaty
      // Group Layers, Remove deduksi Share, Achievement Refresh/Excel, Value
      // Difference, dan cabang Adjust Premium KINI hidup. Yang tersisa:
      // Submit Achievement (`InsertToLogAchievement` menulis log — jalur tulis),
      'DetailLimits · Submit',
      // Retro Share Prop (`AddFacRetroProp`, keputusan §17),
      'Share · Add',
      // halaman sesi `facsharedisp`,
      'Share · Hide Facultative Share (unused)', // pyLabel apa adanya
      // harness Facultative Share Actual (`showHarness` tanpa Activity),
      'TreatyInActualShare · Show Facultative Share',
      // dan varian NON-EDM tombol tulis (`TreatyInSubmit`,
      // `TreatyInDeclineConfirmation`, blok `!TreatyMasterInEDM`) — milik
      // layar Treaty In; 280 dari 280 penyesuaian ber-EDMState.
      'TreatyInfoSubmit · Decline offer',
      'TreatyInfoSubmit · Submit',
    ])
  })

  it('⭐ tombol TULIS EDM — dijalankan form lewat `/api/treaty-in/penyesuaian/*`', () => {
    const tulis = tombolBaru()
      .filter((x) => golongan(x.t) === 'tulis')
      .map((x) => `${x.tab.split('#')[1] ?? x.tab} · ${labelTombol(x.t)} · ${tulisanDari(x.t) ?? ''}`)
    expect([...new Set(tulis)].sort()).toEqual(['TreatyInfoSubmit · Decline offer · tolak', 'TreatyInfoSubmit · Submit · submit'])
  })
})

describe('tambah baris — disalin dari Activity', () => {
  const cari = (aktivitas: string, type?: string) =>
    tombolBaru().find((x) => x.t.aksi.some((a) => (a.aktivitas ?? a.transformasi) === aktivitas && (type === undefined || a.param?.Type === type)))?.t

  it('EGNPI: mata uang disalin dari baris Retention PERTAMA (TreatyInNonAddItem langkah 3)', () => {
    const t = cari('TreatyInNonAddItem', 'egnpi')
    expect(t).toBeDefined()
    const tb = t === undefined ? undefined : tambahDari(t)
    expect(tb?.larik).toBe('EGNPI')
    expect(tb?.baris(sisi({}, { Retention: [{ Currency: 'USD', CurrencyID: '10001' }, { Currency: 'IDR' }] }))).toEqual({
      ID: '',
      Currency: 'USD',
      CurrencyID: '10001',
    })
  })

  it('Limits NP: ID = jumlah layer sesudah ditambah (langkah 5, lalu keluar)', () => {
    const t = cari('TreatyInNonAddItem', 'limits')
    const tb = t === undefined ? undefined : tambahDari(t)
    expect(tb?.baris(sisi({}, { Limits: [{}, {}] }))).toEqual({ ID: '3' })
  })

  it('Accumulation: Period = awalan AccumulationPeriod + urutan (TreatyInAddAccumulation)', () => {
    const t = cari('TreatyInAddAccumulation')
    const tb = t === undefined ? undefined : tambahDari(t)
    expect(tb?.larik).toBe('AccumulationList')
    expect(tb?.baris(sisi({ AccumulationPeriod: 'quarter' }, { AccumulationList: [{}, {}] }))).toEqual({ Period: 'Q 3' })
    expect(tb?.baris(sisi({ AccumulationPeriod: 'none' }))).toEqual({ Period: 'T 1' })
    // Periode di luar keempatnya: `param.keyword` tidak pernah diisi.
    expect(tb?.baris(sisi({ AccumulationPeriod: 'other' }))).toEqual({ Period: '1' })
  })
})

describe('kunci baca-saja', () => {
  it("'selalu' terkunci; syarat terkunci bila salah satu benar; tanpa syarat bebas", () => {
    expect(terkunci('selalu', {})).toBe(true)
    expect(terkunci(['TreatyIn.ViewState = 1', 'TreatyIn.EDMMaterialType = 2'], { ViewState: '0', EDMMaterialType: '2' })).toBe(true)
    expect(terkunci(['TreatyIn.ViewState = 1', 'TreatyIn.EDMMaterialType = 2'], { ViewState: '0', EDMMaterialType: '1' })).toBe(false)
    expect(terkunci(null, { ViewState: '1' })).toBe(false)
  })
})

describe('rantai rumus — rute /hitung Treaty In', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('Update Total Retention: POST /hitung/retensi aksi total, hasil ditimpakan, kunci lain tidak hilang', async () => {
    const t = tombolBaru().find((x) => x.t.label === 'Update Total' && x.tab.endsWith('Maximum Retention'))?.t
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const panggil = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          retensi: [{ ID: '', Currency: 'IDR', Amount: '10' }],
          TotalRetentionAmountNP: [{ Currency: 'IDR', CurrencyID: '10026', Value: '10' }],
          pesan: [],
        }),
        { status: 200 },
      ),
    )
    const h = await r?.(sisi({}, { Retention: [{ Currency: 'IDR', Amount: '10', Note: 'catatan' }] }))
    const [url, init] = panggil.mock.calls[0] ?? []
    expect(String(url)).toContain(`${PREFIX_HITUNG_TREATYIN}/retensi`)
    expect(JSON.parse(String(init?.body))).toMatchObject({ aksi: 'total', indeks: 0 })
    expect(h?.larik.TotalRetentionAmountNP).toEqual([{ Currency: 'IDR', CurrencyID: '10026', Value: '10' }])
    expect(h?.larik.Retention?.[0]?.Note).toBe('catatan')
  })
})

describe('baris bersarang — jalur berindeks Pega', () => {
  it('`Larik(n).Medan` menelusuri larik anak; indeks mulai 1', () => {
    const b = { Layer: '1', RnmLimitList: [{ Currency: 'IDR', Value: '100' }, { Currency: 'USD', Value: '7' }] }
    expect(nilaiJalur(b, 'Layer')).toBe('1')
    expect(nilaiJalur(b, 'RnmLimitList(2).Value')).toBe('7')
    expect(nilaiJalur(b, 'RnmLimitList(3).Value')).toBe('')
  })

  it('larik TAMPIL memakai sumbernya (TreatyInNonAddItem 14.3): Display = RnmLimitList / GrossPremiumList', () => {
    const b = { RnmLimitList: [{ Currency: 'IDR', Value: '100' }], GrossPremiumList: [{ Currency: 'IDR', Value: '5' }] }
    expect(nilaiJalur(b, 'RnmLimitListDisplay(1).Value')).toBe('100')
    expect(nilaiJalur(b, 'RnmGrossPremiDisplay(1).Value')).toBe('5')
  })

  it('pohon digabung ke larik: baris membawa anaknya, larik lain utuh', () => {
    const s = gabungPohon({
      medan: {},
      larik: { Limits: [{ Layer: '1' }], EGNPI: [{ Amount: '1' }] },
      pohon: { Limits: [{ Layer: '1', MDPList: [{ Currency: 'IDR', Value: '9' }] }] },
    })
    expect(s.larik.Limits?.[0]?.MDPList).toEqual([{ Currency: 'IDR', Value: '9' }])
    expect(s.larik.EGNPI).toEqual([{ Amount: '1' }])
  })
})

describe('desimal eksak', () => {
  it('menjumlah tanpa float; skala hasil = skala terbesar', () => {
    expect(tambah('1500000000', '250000000.5')).toBe('1750000000.5')
    expect(tambah('0.1', '0.2')).toBe('0.3')
    expect(tambah('-1.25', '1')).toBe('-0.25')
    expect(tambah('', '5')).toBe('5')
  })
})

describe('TreatyInSummaryMDP — disalin langkah demi langkah', () => {
  it('jumlah per layer per IDR/USD; mata uang lain diabaikan; Note dirangkai', () => {
    const layers = [
      { LayerType: 'Layer ', Layer: '1', LayerPartType: 'Part ', LayerPart: '1', MDPList: [
        { Currency: 'IDR', Value: '100.5' }, { Currency: 'USD', Value: '7' }, { Currency: 'EUR', Value: '9' }, { Currency: 'IDR', Value: '0.5' },
      ] },
      { LayerType: 'Layer ', Layer: '2', LayerPartType: 'Part ', LayerPart: '1', MDPList: [{ Currency: 'USD', Value: '3' }] },
      // Kunci SAMA dengan layer pertama → dijumlah ke barisnya (langkah 3.2).
      { LayerType: 'Layer ', Layer: '1', LayerPartType: 'Part ', LayerPart: '1', MDPList: [{ Currency: 'IDR', Value: '1' }] },
      // Tanpa MDP → tanpa baris.
      { LayerType: 'Layer ', Layer: '3', LayerPartType: 'Part ', LayerPart: '1', MDPList: [] },
    ]
    expect(ringkasanMDP(layers)).toEqual([
      { LayerType: 'Layer ', Layer: '1', LayerPartType: 'Part ', LayerPart: '1', IDR: '102.0', USD: '7', Note: 'Layer 1 of Part 1' },
      { LayerType: 'Layer ', Layer: '2', LayerPartType: 'Part ', LayerPart: '1', IDR: '0', USD: '3', Note: 'Layer 2 of Part 1' },
    ])
  })
})

describe('rantai Limits dan Share — satu klik, satu panggilan rute', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })
  const tombol = (tab: string, label: string) => tombolBaru().find((x) => x.tab.endsWith(tab) && x.t.label === label)?.t

  it('Limits Update Total: SATU POST limit-np `total` (SummaryLimit tercakup) + SummaryMDP di layar', async () => {
    const t = tombol('#Limits', 'Update Total')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const panggil = vi.spyOn(globalThis, 'fetch').mockImplementation(() =>
      Promise.resolve(
        new Response(
          JSON.stringify({
            layers: [{ ID: '1', Layer: '1', MDPList: [{ Currency: 'IDR', Value: '5' }] }],
            Total: { TotalLimitMDPNP: [{ Currency: 'IDR', CurrencyID: '10026', Value: '5' }] },
            TotalLimitsROL: '2.5',
            LimitSummaryList: [{ Note: 'Layer 1', MDP: '5' }],
            pesan: [],
          }),
          { status: 200 },
        ),
      ),
    )
    const h = await r?.({ medan: {}, larik: { Limits: [{ ID: '1', Layer: '1', LayerType: 'Layer ', MDPList: [{ Currency: 'IDR', Value: '5' }] }] } })
    expect(panggil).toHaveBeenCalledTimes(1)
    const [url, init] = panggil.mock.calls[0] ?? []
    expect(String(url)).toContain(`${PREFIX_HITUNG_TREATYIN}/limit-np`)
    expect(JSON.parse(String(init?.body))).toMatchObject({ aksi: 'total' })
    expect(h?.medan.TotalLimitsROL).toBe('2.5')
    expect(h?.larik.LimitSummaryList).toEqual([{ Note: 'Layer 1', MDP: '5' }])
    expect(h?.larik.TotalLimitMDPNP).toEqual([{ Currency: 'IDR', CurrencyID: '10026', Value: '5' }])
    // Layer DITIMPAKAN — kunci yang rute tidak kembalikan (LayerType) tetap ada.
    expect(h?.larik.Limits?.[0]?.LayerType).toBe('Layer ')
    expect(h?.larik.MDPSummaryList?.[0]?.IDR).toBe('5')
  })

  it('Share Update Total: SATU POST share-np `total`; skalar, larik, dan sembilan total dipetakan kembali', async () => {
    const t = tombol('#Share', 'Update Total')
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const panggil = vi.spyOn(globalThis, 'fetch').mockImplementation(() =>
      Promise.resolve(
        new Response(
          JSON.stringify({
            share: {
              RNMShare: '40',
              Share: [{ Layer: '1', RNMShare: '40', GrossPremiumList: [{ Currency: 'IDR', Value: '2' }] }],
              LimitShareSummaryList: [{ Note: 'Layer 1' }],
              Total: { TotalShareNetNP: [{ Currency: 'IDR', CurrencyID: '10026', Value: '2' }] },
            },
            pesan: [],
          }),
          { status: 200 },
        ),
      ),
    )
    const h = await r?.({ medan: { RNMShare: '40', ID: '1000080/R01' }, larik: { Limits: [{ Layer: '1' }] } })
    expect(panggil).toHaveBeenCalledTimes(1)
    const badan = JSON.parse(String(panggil.mock.calls[0]?.[1]?.body)) as { aksi: string; share: { RNMShare: string; Total: Record<string, unknown> }; idKontrak: string }
    expect(badan.aksi).toBe('total')
    expect(badan.share.RNMShare).toBe('40')
    expect(Object.keys(badan.share.Total)).toHaveLength(9)
    expect(badan.idKontrak).toBe('1000080/R01')
    expect(h?.medan.RNMShare).toBe('40')
    expect(h?.larik.Share?.[0]?.GrossPremiumList).toEqual([{ Currency: 'IDR', Value: '2' }])
    expect(h?.larik.TotalShareNetNP).toEqual([{ Currency: 'IDR', CurrencyID: '10026', Value: '2' }])
  })
})

describe('rumus Installment — rute /hitung/angsuran Treaty In', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  const jawab = (isi: unknown) => vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify(isi), { status: 200 }))

  it('Update Value: status=update, Installment dari isian, OLDDATA ikut; larik Installment DIGANTI utuh', async () => {
    const t = tombolBaru().find((x) => x.t.label === 'Update Value' && x.tab.endsWith('#Installment'))?.t
    const r = t === undefined ? undefined : rantaiRumus(t)
    expect(r).toBeDefined()
    const panggil = jawab({
      angsuran: [{ Currency: 'IDR', AmountTotal: '100', PctTotal: '100', InstallmentList: [{ Installment: '1', InstallmentPct: '100', Amount: '100' }] }],
      TotalInstallmentNP: [{ Currency: 'IDR', CurrencyID: '', Value: '100' }],
      InstallmentNo: '1',
      pesan: [],
    })
    const kini = sisi({ InstallmentNo: '1', EDMState: '2' }, { Installment: [{ Currency: 'IDR', pxListSubscript: '1' }], TotalShareNetNP: [{ Currency: 'IDR', Value: '100' }] })
    const lama = { medan: {}, larik: { Installment: [{ Currency: 'IDR', InstallmentList: [{ InstallmentPct: '60', DueDate: '20250101' }] }] } }
    const h = await r?.(kini, { lama })
    const [url, init] = panggil.mock.calls[0] ?? []
    expect(String(url)).toContain(`${PREFIX_HITUNG_TREATYIN}/angsuran`)
    const badan = JSON.parse(String(init?.body)) as Record<string, unknown>
    expect(badan).toMatchObject({ aksi: 'nilai', status: 'update', installmentNo: '1', edmState: '2', netPremium: [{ Currency: 'IDR', CurrencyID: '', Value: '100' }] })
    expect(badan.angsuranLama).toMatchObject([{ Currency: 'IDR', InstallmentList: [{ InstallmentPct: '60', DueDate: '20250101' }] }])
    // Langkah 2 membuang larik lama — `pxListSubscript` tidak dibawa pulang.
    expect(h?.larik.Installment).toEqual([{ Currency: 'IDR', AmountTotal: '100', PctTotal: '100', InstallmentList: [{ Installment: '1', InstallmentPct: '100', Amount: '100' }] }])
    expect(h?.larik.TotalInstallmentNP).toEqual([{ Currency: 'IDR', CurrencyID: '', Value: '100' }])
    expect(h?.medan.InstallmentNo).toBe('1')
  })

  it('isian Installment (perilaku change): postValue + TreatyInSetValueInstallment TANPA status', async () => {
    const m = KERANGKA_TAB['TreatyInTabsNonProportional#Installment']?.isi.find((b) => b.t === 'medan' && b.kunci === 'InstallmentNo')
    const aksi = m?.t === 'medan' ? m.aksiUbah : undefined
    expect(aksi?.map((a) => a.aksi)).toEqual(['postValue', 'refresh'])
    const r = aksi === undefined ? undefined : rantaiRumus({ aksi })
    const panggil = jawab({ angsuran: [], TotalInstallmentNP: [], InstallmentNo: '0', pesan: ['Installment Value cannot be less than 1'] })
    const h = await r?.(sisi({ InstallmentNo: '0' }))
    expect(JSON.parse(String(panggil.mock.calls[0]?.[1]?.body))).toMatchObject({ aksi: 'nilai', status: '', installmentNo: '0' })
    expect(h?.pesan).toEqual(['Installment Value cannot be less than 1'])
    expect(h?.larik.Installment).toEqual([])
  })

  it('sel % Installment (rincian): SetTotalInstallment editpercentage — hanya Amount, AmountTotal, PctTotal yang berubah', async () => {
    const g = KERANGKA_RINCIAN.Installments?.find((b) => b.t === 'grid')
    const i = g?.t === 'grid' ? g.kunci.indexOf('InstallmentPct') : -1
    const aksi = g?.t === 'grid' ? g.aksiUbah[i] : null
    expect(aksi).toEqual([{ aksi: 'refresh', aktivitas: 'SetTotalInstallment', param: { status: 'editpercentage' } }])
    const r = aksi === null || aksi === undefined ? undefined : rantaiRumus({ aksi })
    const panggil = jawab({
      angsuran: [{ Currency: 'IDR', AmountTotal: '30', PctTotal: '30', InstallmentList: [{ InstallmentPct: '30', Amount: '30' }] }],
      TotalInstallmentNP: null,
      InstallmentNo: '',
      pesan: [],
    })
    const halaman = sisi({ Currency: 'IDR', AmountTotal: '9' }, { InstallmentList: [{ InstallmentPct: '30', Amount: '9', WPC: '45', pxListSubscript: '1' }] })
    const panel = sisi({}, { TotalShareNetNP: [{ Currency: 'IDR', Value: '100' }] })
    const h = await r?.(halaman, { panel })
    const badan = JSON.parse(String(panggil.mock.calls[0]?.[1]?.body)) as Record<string, unknown>
    expect(badan).toMatchObject({ aksi: 'total-baris', status: 'editpercentage', indeks: 0, netPremium: [{ Currency: 'IDR', Value: '100' }] })
    expect(h?.larik.InstallmentList).toEqual([{ InstallmentPct: '30', Amount: '30', WPC: '45', pxListSubscript: '1' }])
    expect(h?.medan).toEqual({ AmountTotal: '30', PctTotal: '30' })
  })

  it('Update Total Installment: TotalInstallmentNP dari AmountTotal tiap halaman', async () => {
    const t = tombolBaru().find((x) => x.t.label === 'Update Total' && x.tab.endsWith('#Installment'))?.t
    const r = t === undefined ? undefined : rantaiRumus(t)
    const panggil = jawab({ angsuran: [], TotalInstallmentNP: [{ Currency: 'IDR', CurrencyID: '', Value: '5' }], InstallmentNo: '', pesan: [] })
    const h = await r?.(sisi({}, { Installment: [{ Currency: 'IDR', AmountTotal: '5', InstallmentList: [] }] }))
    expect(JSON.parse(String(panggil.mock.calls[0]?.[1]?.body))).toMatchObject({ aksi: 'total', angsuran: [{ Currency: 'IDR', AmountTotal: '5' }] })
    expect(h?.larik).toEqual({ TotalInstallmentNP: [{ Currency: 'IDR', CurrencyID: '', Value: '5' }] })
  })

  it('⛔ hapus baris yang disusul Activity tanpa rumus DIMATIKAN — rantai separuh jalan tidak dijalankan', () => {
    const hapus = (aksi: AksiTombol[]): TombolKerangka => ({ t: 'tombol', at: 0, label: 'Remove', syarat: [], aksi })
    expect(rantaiSesudahHapus(hapus([{ aksi: 'deleteRow' }]))).toBeNull()
    expect(rantaiSesudahHapus(hapus([{ aksi: 'deleteRow' }, { aksi: 'refresh', aktivitas: 'InsertToLogAchievement' }]))).toBeUndefined()
    expect(rantaiSesudahHapus(hapus([{ aksi: 'deleteRow' }, { aksi: 'refresh', aktivitas: 'TreatyInNPSetTotal', param: { type: 'installment' } }]))).toBeTypeOf('function')
  })
})
