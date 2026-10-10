// Uji tab Share NON-PROP — bentuk, syarat tampil, dan tombol dari ekspor
// (`TreatyInTabsNonProportional.xml` @1695720, `Share.xml`).
//
// Dirender statis (`react-dom/server`) — BUKAN pengganti melihat layar di
// peramban. Rumusnya diuji di services (`hitung_share_np*_test.go`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisShareNP, RingkasanShareNP, ShareNP } from './api'
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

  // ⭐ 8 Oktober 2026 — bentuk Pega: grid tetap tampil, satu sel "No items",
  // tanpa petunjuk tambahan ("jangan ada design tambahan").
  it('⭐ kosong → grid Pega dengan satu sel "No items"', () => {
    const html = render(share({ Share: [] }))
    expect(html).toContain(`<td colSpan="15" class="trin__kosong-pega">${SHARE_NP.tanpaBaris}</td>`)
    expect(html).not.toContain(SHARE_NP.petunjukKosong)
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

  it('⭐ Spreading Type bernama: dropdown dari RD, anak susunan, Spreading Total Pct', () => {
    const html = panel(baris())
    // ⛔ Combobox sejak 8 Oktober 2026 (*"mengetik harus terasa seperti
    // mencari"*): teks pilihan ada di KOTAKNYA, dan daftar butirnya baru
    // dirender ketika dibuka.
    expect(html).toContain('value="2022 QS 145M TRT"')
    // ⛔ DIBALIK 8 Oktober 2026 (sore) — grid `.SpreadingListXOL` Pega memang
    // ber-`readOnly` (@408855, kerangka `baca: "selalu"`), tetapi pemilik
    // proses meminta spread tetap dapat DITAMBAH walau Spreading Type sudah
    // dipilih: *"konsepnya harus sama seperti yang di prop juga yang di mana
    // spread nya bisa ditambah"*. Grid manual kini dipakai di kedua cabang.
    expect(html).toContain('Add')
    // ⛔ Kolom `Reins Type` TIDAK dapat diubah di cabang ini — permintaan
    // pemilik proses 8 Oktober 2026. Isinya datang dari RD, jadi ia teks,
    // bukan dropdown; yang dapat disunting hanya `Pct` dan jumlah barisnya.
    expect(html).toContain('>QS (OR)</td>')
    expect(html).not.toMatch(/<select[^>]*aria-label="Reins Type"/)
    // Label ekspor `Spreading Total Pct` (Share.xml @473985) + teks `%`.
    expect(html).toContain(SHARE_NP.spreadingTotalPct)
    expect(html).toContain('Deduction Details')
    expect(html).toContain('Brokerage fee')
    // ⭐ 8 Oktober 2026 — grid Treaty Group baca-saja (`CoBListReadOnly`).
    expect(html).toContain('>PROPERTY</td>')
    expect(html).not.toMatch(/<input class="field__input" list="[^"]*"[^>]*value="PROPERTY"/)
  })

  it('⭐ `pyReadOnlyCondition` MENIMPA `pyReadOnly`: Cover aktif di mode Edit; Layer SELALU baca-saja', () => {
    // `hanya_baca()` atas `Section/Share.xml`: sel ber-`pyReadOnly=true` DAN
    // `pyReadOnlyCondition = TreatyIn.ViewState = 1` — terkunci HANYA di mode
    // lihat. ⭐ 8 Oktober 2026: KECUALI sel layer ber-`pyDisabledNew =
    // always` (kerangka Adjustment `baca: "selalu"`) — teks tanpa label.
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
    expect(ubah).toContain('value="Risk"')
    // `Layer` SELALU baca-saja: teks, bukan kontrol berdaftar.
    expect(ubah).toContain('>Layer<')
    expect((ubah.match(/<span class="trin__teks-sel">Layer<\/span>/g) ?? []).length).toBe(2)
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
    expect(lihat).not.toContain('role="combobox"')
    expect(lihat).toMatch(
      /<label class="field__label">Cover<\/label><input class="field__input field__input--readonly" type="text" readonly="" value="Risk"/,
    )
  })

  it('⭐ Spreading Type kosong: grid manual Reins Type · Pct Share, Total Share Pct', () => {
    const html = panel(baris({ SpreadingTypeXOL: '', SpreadingListXOL: [] }))
    expect(html).toContain('Pct Share')
    expect(html).toContain(SHARE_NP.totalSharePct)
    expect(html).not.toContain(SHARE_NP.spreadingTotalPct)
  })

  // ⭐ 9 Oktober 2026 — DUA spreading persis Pega (keputusan pemakai):
  // terisi = spreading LAMA (grid `readOnly` @408855, nol Add/Delete);
  // kosong = spreading BARU (grid manual @54364, Add/Delete, TANPA dropdown).
  it('⭐ spreading LAMA (Spreading Type terisi): dropdown + grid baca, nol Add/Delete', () => {
    const html = panel(baris())
    expect(html).toContain(SHARE_NP.spreadingType)
    expect(html).not.toContain('Pct Share')
    expect(html).not.toContain('tl-tambah')
    expect(html).not.toContain(`${SHARE_NP.hapus} ${SHARE_NP.spreading} 1`)
    expect(html).not.toMatch(/<select[^>]*aria-label="Reins Type"/)
  })

  it('⭐ spreading BARU (Spreading Type kosong) mode Edit: Add/Delete, TANPA dropdown Spreading Type', () => {
    const html = panel(baris({ SpreadingTypeXOL: '', SpreadingListXOL: [{ ReinsTypeName: '', ReinsTypeID: '', Pct: '0' }] as never }))
    expect(html).not.toContain(SHARE_NP.spreadingType)
    expect(html).toContain('tl-tambah')
    expect(html).toContain(`${SHARE_NP.hapus} ${SHARE_NP.spreading} 1`)
    expect(html).toMatch(/<select[^>]*aria-label="Reins Type"/)
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
      // ⭐ 8 Oktober 2026 — kepala ekspor dua sel: judul · '' (grid Pega).
      const p = html.indexOf(`<th scope="col">${j}</th><th scope="col" class="trin__angka"></th>`, posisi + 1)
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
    expect(html).toContain('role="combobox"')
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
    // ⛔ BUKAN `<datalist>` lagi: ketikan bebas tidak boleh lolos jadi nilai.
    expect(html).toContain('value="IDR"')
    expect(html).not.toMatch(/<input[^>]*\slist=/)
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

// ⛔ ANGKA `Summarry of RNM Share` — laporan pemilik proses 8 Oktober 2026:
// *"maksimal di belakang koma berapa, dan kalau kepanjangan jangan sampai ke
// bawah, angkanya harus tetap menyamping"*.
//
// Sebabnya sel ini diformat APA ADANYA (`formatLimit('uang', -1, …)`), jadi
// hasil pembagian membawa ekornya utuh: `11.652.832,7749999503125`. Ekor itu
// membungkus ke baris berikutnya dan menaikkan tinggi seluruh barisnya.
describe('Summarry of RNM Share — angka dibulatkan dan tidak membungkus', () => {
  const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')

  it('⛔ maksimal DUA desimal, dan nol di ekor dibuang', () => {
    const dasar = share().LimitShareSummaryList[0] as RingkasanShareNP
    const html = render(
      share({ LimitShareSummaryList: [{ ...dasar, MDP: '11652832.7749999503125', Deductible: '20119177.6749' }] }),
      'lihat',
    )
    expect(html).toContain('11.652.832,77')
    expect(html).toContain('20.119.177,67')
    // Ekor panjangnya TIDAK boleh sampai ke layar.
    expect(html).not.toContain('7749999503125')
    expect(html).not.toContain('6749<')
  })

  it('⭐ bilangan bulat tetap TANPA `,00` — catatan ekspor `pxNumber`', () => {
    const dasar = share().LimitShareSummaryList[0] as RingkasanShareNP
    const html = render(share({ LimitShareSummaryList: [{ ...dasar, Limit: '175000000', Limit2: '0' }] }), 'lihat')
    expect(html).toContain('>175.000.000<')
    expect(html).not.toContain('175.000.000,00')
  })

  // ⚠️ `white-space` HARUS di selnya: `<col>` hanya menghormati width,
  // background, border dan visibility — di sana ia diam-diam diabaikan.
  it('⛔ sel angka tidak membungkus, dan aturannya di SEL bukan di `<col>`', () => {
    const i = CSS.indexOf('.treatyin .trin__tabel--pega td.trin__angka')
    expect(i).toBeGreaterThan(0)
    expect(CSS.slice(i, CSS.indexOf('}', i))).toContain('white-space: nowrap')
    // Boundary-nya penting: tanpa `` pola ini ikut mencocoki
    // `.trin__buka-kolom` dan gagal atas aturan yang benar.
    expect(CSS).not.toMatch(/[\s,>]col\s*\{[^}]*white-space/)
  })
})
