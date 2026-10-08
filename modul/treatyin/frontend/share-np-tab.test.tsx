// Uji tab Share NON-PROP — bentuk, syarat tampil, dan tombol dari ekspor
// (`TreatyInTabsNonProportional.xml` @1695720, `Share.xml`).
//
// Dirender statis (`react-dom/server`) — BUKAN pengganti melihat layar di
// peramban. Rumusnya diuji di services (`hitung_share_np*_test.go`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisShareNP, ShareNP } from './api'
import TabShareNonProp, { RincianShare, SHARE_NP_KOSONG, judulBarisShare } from './components/TabShareNonProp'
import { GRID_RINCIAN_SPREADING, GRID_TOTAL_SHARE, SHARE_NP } from './labelsShareNP'

const AKAR = __dirname
const KOMP = readFileSync(join(AKAR, 'components', 'TabShareNonProp.tsx'), 'utf8')
const FORM = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')

const nilai = (Currency: string, Value: string) => ({ Currency, CurrencyID: '', Value })

function baris(isi: Partial<BarisShareNP> = {}): BarisShareNP {
  return {
    LayerType: 'Layer',
    Layer: '1',
    LayerPartType: 'Layer',
    LayerPart: '1',
    Cover: 'risk',
    RNMShare: '40',
    TreatyGroupList: [{ TreatyGroup: 'PROPERTY', TreatyGroupID: '10002' }],
    Limit: '',
    Limit2: '',
    SpreadingTypeXOL: '2022 QS 145M TRT',
    SpreadingTypeIDXOL: '10227',
    SpreadingTotalPctXOL: '100',
    SpreadingListXOL: [
      {
        ReinsTypeName: 'QS (OR)',
        ReinsTypeID: '10028',
        ParentReinsTypeID: '10227',
        Pct: '40',
        Rp: '',
        Usd: '',
        RnmLimitList: [],
        GrossPremiumList: [],
        GrossPremiumMinList: [],
        DeductionTotalList: [],
        NetPremiumList: [],
      },
      {
        ReinsTypeName: 'QS (R/I)',
        ReinsTypeID: '10004',
        ParentReinsTypeID: '10227',
        Pct: '60',
        Rp: '',
        Usd: '',
        RnmLimitList: [],
        GrossPremiumList: [],
        GrossPremiumMinList: [],
        DeductionTotalList: [],
        NetPremiumList: [],
      },
    ],
    DeductionList: [
      {
        Comment: 'Brokerage fee',
        Currency: 'IDR',
        CurrencyID: '',
        Deduction: '2000000',
        DeductionPct: '10',
        DeductionPctCalculate: 'true',
      },
    ],
    RnmLimitList: [nilai('IDR', '400000000'), nilai('USD', '40000')],
    GrossPremiumList: [nilai('IDR', '20000000')],
    GrossPremiumMinList: [],
    DeductionTotalList: [nilai('IDR', '2000000')],
    NetPremiumList: [nilai('IDR', '18000000')],
    RNMSpreadedListXOL: [],
    RNMSpreadedListRIXOL: [],
    RNMSpreadedListGrossXOL: [],
    RNMSpreadedListGrossRIXOL: [],
    RNMSpreadedListGrossMinXOL: [],
    RNMSpreadedListGrossRIMinXOL: [],
    RNMSpreadedListDeductXOL: [],
    RNMSpreadedListDeductRIXOL: [],
    RNMSpreadedListNetXOL: [],
    RNMSpreadedListNetRIXOL: [],
    ...isi,
  }
}

function share(isi: Partial<ShareNP> = {}): ShareNP {
  return {
    ...SHARE_NP_KOSONG,
    RNMShare: '40',
    BrokeragePercent: '10',
    ShareReins: [{ ID: '1', ReinsID: '100', ReinsName: 'SINGAPORE REINSURANCE CORPORATION LIMITED', Layer: '', SharePct: '5' }],
    Share: [baris()],
    LimitShareSummaryList: [
      {
        LayerType: 'Layer',
        Layer: '1',
        LayerPartType: 'Layer',
        LayerPart: '1',
        Note: 'Layer1 of Layer1',
        Limit: '400000000',
        Limit2: '40000',
        MDP: '20000000',
        MDP2: '',
        Deductible: '2000000',
        Deductible2: '',
        NetPremi: '18000000',
        NetPremi2: '',
      },
    ],
    Total: { TotalSpreadedRnmProp: [nilai('IDR', '9476000000')], TotalSpreadedRnmRIProp: [nilai('IDR', '14224000000')] },
    ...isi,
  }
}

const render = (s: ShareNP, mode: 'lihat' | 'ubah' = 'ubah', lain: { edmState?: string; edmJenisMaterial?: string } = {}) =>
  renderToStaticMarkup(<TabShareNonProp share={s} layers={[]} mode={mode} reasuradur={[]} induk={{}} {...lain} />)

describe('tab Share Non-Prop — bentuk ekspor', () => {
  it('⭐ disambung ke cabang NON-PROP form, membaca `warisan.shareNP`', () => {
    expect(FORM).toContain('<TabShareNonProp')
    expect(FORM).toContain('share={shareNP ?? warisan?.shareNP}')
    // ⭐ Clipboard bersama: Update Summary membaca layer yang SEDANG disunting
    // di tab Limits, dan kedua tab hidup lintas pindah tab.
    expect(FORM).toContain('layers={limitsNP?.layers ?? warisan?.limitsPohon ?? []}')
    expect(FORM).toContain('onUbah={setShareNP}')
    expect(FORM).toContain('onUbah={(layers, akar) => setLimitsNP({ layers, akar })}')
    expect(FORM).toContain("commencement={keSimpan(mulai) || (warisan?.tanggalMulaiAsli ?? '')}")
    expect(FORM).not.toContain('<SubTabShare')
  })

  it('⭐ panel Share: isian akar berurut ekspor, Update Summary di mode Edit', () => {
    const html = render(share())
    const urut = [
      SHARE_NP.persenRnm,
      SHARE_NP.persenBrokerage,
      SHARE_NP.acrossTheBoard,
      SHARE_NP.shareKeRetro,
      'Reinsurer Name',
      'RNM Share',
      SHARE_NP.ringkasan,
      SHARE_NP.totalSemua,
    ]
    let posisi = -1
    for (const t of urut) {
      const p = html.indexOf(t, posisi + 1)
      expect(p, t).toBeGreaterThan(posisi)
      posisi = p
    }
    expect(html).toContain(SHARE_NP.perbaruiRingkasan)
    expect(html).toContain(SHARE_NP.perbaruiTotal)
    // `pyDefaultValue = true` — dicentang.
    expect(html).toMatch(/<input type="checkbox" checked=""/)
  })

  it('⛔ Brokerage From Other Retro, Facultative Reinsurers, Share to RNM HANYA bila Share to Other Retro > 0', () => {
    const tanpa = render(share({ FacultativeShare: '0' }))
    for (const t of [SHARE_NP.brokerageRetro, 'Facultative Reinsurers', SHARE_NP.shareKeRnm]) {
      expect(tanpa).not.toContain(t)
    }
    const dengan = render(share({ FacultativeShare: '10', RNMShare: '0', RnmShareDeducted: '-10' }))
    for (const t of [SHARE_NP.brokerageRetro, 'Facultative Reinsurers']) {
      expect(dengan).toContain(t)
    }
    // Tangkapan layar Pega: "Share to RNM : -10 %".
    expect(dengan).toContain(`${SHARE_NP.shareKeRnm} <strong>-10</strong> %`)
  })

  it('⭐ grid RNM Share: layer · Part of · 100% Limit ×2 · MDP ×2 · % Share', () => {
    const html = render(share())
    expect(html).toContain(
      '<th scope="col" colSpan="2">100% Limit</th><th scope="col" colSpan="2">100% Limit</th><th scope="col" colSpan="2">MDP</th><th scope="col" colSpan="2">MDP</th><th scope="col">% Share</th>',
    )
    expect(html).toContain('<td class="tl-share-of">Part of</td>')
    expect(html).toContain('400.000.000,00')
    expect(html).toContain('40.000,00')
    expect(judulBarisShare(baris())).toBe('Layer 1 Part of Layer 1')
  })

  it('⭐ Summarry of RNM Share (ejaan Pega) dan sembilan grid total berurut ekspor', () => {
    const html = render(share())
    expect(html).toContain('Summarry of RNM Share')
    expect(html).toContain('Layer1 of Layer1')
    const judul = GRID_TOTAL_SHARE.flat()
      .filter((g) => g !== null)
      .map((g) => g.judul)
    expect(judul).toEqual([
      'Total RNM Limit (RNM Share)',
      'Total OR Limit',
      'Total R/I Limit',
      'Total Gross Min Premium',
      'Total Gross Premium (MDP)',
      'Total Deduction',
      'Total Net Premium',
      'Total OR Net Premium',
      'Total R/I Net Premium',
    ])
    let posisi = -1
    for (const j of judul) {
      const p = html.indexOf(`<th scope="col">${j}</th>`)
      expect(p, j).toBeGreaterThan(posisi)
      posisi = p
    }
    expect(html).toContain('9.476.000.000,00')
    expect(html).toContain('14.224.000.000,00')
  })

  it('⛔ mode lihat: nol tombol; EDM material 2: Update Summary/Total mati', () => {
    const lihat = render(share(), 'lihat')
    expect(lihat).not.toContain(SHARE_NP.perbaruiRingkasan)
    expect(lihat).not.toContain(SHARE_NP.perbaruiTotal)
    expect(lihat).not.toContain('tl-tambah')
    const edm = render(share(), 'ubah', { edmJenisMaterial: '2' })
    expect(edm).toMatch(/<button type="button" class="btn btn--primary btn--sm" disabled="">Update Summary<\/button>/)
  })

  it('⛔ Update Value in Share HANYA untuk kontrak revisi EDM (`TreatyMasterInEDM`)', () => {
    expect(render(share())).not.toContain(SHARE_NP.perbaruiNilai)
    expect(render(share(), 'ubah', { edmState: '1' })).toContain(SHARE_NP.perbaruiNilai)
  })

  it('⭐ kosong → "No items" dengan petunjuk Update Summary', () => {
    const html = render(share({ Share: [] }))
    expect(html).toContain(SHARE_NP.petunjukKosong)
  })

  it('⛔ nol rumus di layar — hanya memanggil services', () => {
    expect(KOMP).toContain('hitungShareNP(')
    expect(KOMP).not.toMatch(/@divide|\* 100|\/ 100/)
  })
})

describe('panel rincian baris — flow action `Share`', () => {
  const panel = (b: BarisShareNP, modeUbah = true) =>
    renderToStaticMarkup(
      <RincianShare
        b={b}
        induk={[
          {
            reinsTypeId: '10227',
            reinsTypeName: '2022 QS 145M TRT',
            parentReinsTypeId: '00',
            treatyYearId: '1000610',
            pct: '',
            rp: '',
            usd: '',
          },
        ]}
        bisaUbah={modeUbah}
        modeUbah={modeUbah}
        onUbah={() => undefined}
        hitung={() => undefined}
      />,
    )

  it('⭐ Spreading Type bernama: dropdown dari RD, anak susunan, Total Spreading Pct', () => {
    const html = panel(baris())
    expect(html).toContain('<option value="2022 QS 145M TRT" selected="">2022 QS 145M TRT</option>')
    // Sel grid spreading bernama ber-`Auto` — isian di mode Edit.
    expect(html).toMatch(/<input class="field__input"[^>]*value="QS \(OR\)"/)
    expect(html).toContain(SHARE_NP.totalSpreadingPct)
    expect(html).toContain('Deduction Details')
    expect(html).toContain('Brokerage fee')
    expect(html).toMatch(/<input class="field__input" list="[^"]*"[^>]*value="PROPERTY"/)
  })

  it('⭐ `pyReadOnlyCondition` MENIMPA `pyReadOnly`: Layer · Cover · Treaty Group aktif di mode Edit', () => {
    // `hanya_baca()` atas `Section/Share.xml`: kelima sel ber-`pyReadOnly=true`
    // DAN `pyReadOnlyCondition = TreatyIn.ViewState = 1` — terkunci HANYA di
    // mode lihat.
    const opsi = {
      jenisTreaty: [],
      kelompokTreaty: [{ id: '10002', nama: 'PROPERTY', kembar: false }],
      mataUang: [],
      jenisLayer: [{ value: 'Layer', label: 'Layer' }],
      cover: [{ value: 'risk', label: 'Risk' }],
    }
    const ubah = renderToStaticMarkup(
      <RincianShare b={baris()} induk={[]} opsi={opsi} bisaUbah modeUbah onUbah={() => undefined} hitung={() => undefined} />,
    )
    expect(ubah).toContain('<option value="risk" selected="">Risk</option>')
    expect((ubah.match(/<option value="Layer" selected="">Layer<\/option>/g) ?? []).length).toBe(2)
    const lihat = renderToStaticMarkup(
      <RincianShare
        b={baris()}
        induk={[]}
        opsi={opsi}
        bisaUbah={false}
        modeUbah={false}
        onUbah={() => undefined}
        hitung={() => undefined}
      />,
    )
    expect(lihat).not.toContain('<select')
    expect(lihat).toMatch(
      /<label class="field__label">Cover<\/label><input class="field__input field__input--readonly" type="text" readonly="" value="Risk"/,
    )
  })

  it('⭐ Spreading Type kosong: grid manual Reins Type · Pct Share, Total Share Pct', () => {
    const html = panel(baris({ SpreadingTypeXOL: '', SpreadingListXOL: [] }))
    expect(html).toContain('Pct Share')
    expect(html).toContain(SHARE_NP.totalSharePct)
    expect(html).not.toContain(SHARE_NP.totalSpreadingPct)
  })

  it('⛔ mode lihat: Spreading Type kosong tidak menampilkan dropdown', () => {
    const html = panel(baris({ SpreadingTypeXOL: '', SpreadingListXOL: [] }), false)
    expect(html).not.toContain(SHARE_NP.spreadingType)
    expect(html).not.toContain('tl-tambah')
  })

  it('⭐ blok `hidden, reference` TAMPIL bersama Spreading bernama — 15 grid dasar · OR · R/I', () => {
    const html = panel(baris({ RNMSpreadedListXOL: [nilai('IDR', '160000000')] }))
    const judul = GRID_RINCIAN_SPREADING.flat().map((g) => g.judul)
    expect(judul).toHaveLength(15)
    let posisi = -1
    for (const j of judul) {
      const p = html.indexOf(`<th scope="colgroup" colSpan="2">${j}</th>`, posisi + 1)
      expect(p, j).toBeGreaterThan(posisi)
      posisi = p
    }
    expect(html).toContain('160.000.000,00')
    // Spreading manual: blok itu tidak tampil (wadahnya `.SpreadingTypeXOL != ''`).
    expect(panel(baris({ SpreadingTypeXOL: '', SpreadingListXOL: [] }))).not.toContain('Spreading OR Limit')
  })

  it('⭐ Reins Type spreading manual: induk SEMUA grup (`TempSprd.TreatyGroupID` kosong → filter dilewati)', () => {
    const lain = {
      reinsTypeId: '10236',
      reinsTypeName: '2023 QS 150M TRT',
      parentReinsTypeId: '00',
      treatyYearId: '1000644',
      pct: '',
      rp: '',
      usd: '',
    }
    const b = baris({
      SpreadingTypeXOL: '',
      SpreadingListXOL: [
        {
          ReinsTypeName: '',
          ReinsTypeID: '',
          ParentReinsTypeID: '',
          Pct: '40',
          Rp: '',
          Usd: '',
          RnmLimitList: [],
          GrossPremiumList: [],
          GrossPremiumMinList: [],
          DeductionTotalList: [],
          NetPremiumList: [],
        },
      ],
    })
    const html = renderToStaticMarkup(
      <RincianShare b={b} induk={[]} indukManual={[lain]} bisaUbah modeUbah onUbah={() => undefined} hitung={() => undefined} />,
    )
    expect(html).toContain('<option value="10236">2023 QS 150M TRT</option>')
  })

  it('⭐ Currency Deduction Details = autocomplete daftar mata uang (`BrowseCurrency_RD`)', () => {
    const opsi = {
      jenisTreaty: [],
      kelompokTreaty: [],
      mataUang: [
        { id: '10026', nama: 'IDR', kembar: false },
        { id: '10001', nama: 'USD', kembar: false },
      ],
    }
    const html = renderToStaticMarkup(
      <RincianShare b={baris()} induk={[]} opsi={opsi} bisaUbah modeUbah onUbah={() => undefined} hitung={() => undefined} />,
    )
    expect(html).toMatch(/<input class="field__input" list="[^"]*"[^>]*value="IDR"/)
    expect(html).toContain('<option value="USD"></option>')
  })

  it('⭐ pesan Activity milik baris tampil di panel barisnya (Pega: pesan medan, bukan kepala tab)', () => {
    const html = renderToStaticMarkup(
      <RincianShare
        b={baris({ SpreadingTypeXOL: '', SpreadingListXOL: [] })}
        induk={[]}
        pesan={['Total share must equal RNM share.!!']}
        bisaUbah
        modeUbah
        onUbah={() => undefined}
        hitung={() => undefined}
      />,
    )
    expect(html).toContain('<ul class="tl-pesan" role="alert"><li>Total share must equal RNM share.!!</li></ul>')
    // Aksi satu baris hanya mengganti pesan baris itu.
    expect(KOMP).toMatch(/const AKSI_BARIS: readonly AksiShareNP\[\] = \[[^\]]*'deduksi'[^\]]*'spreading-pct'/)
  })
})
