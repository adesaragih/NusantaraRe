// Uji tab Limits — Prop (pemicu Activity, nilai awal Add, baca-saja,
// kunci materialitas) dan Non-Prop (grid layer, rincian Layers, Summary,
// Total All Layers).
//
// Dirender statis (`react-dom/server`) — BUKAN pengganti melihat layar di
// peramban. Rumusnya diuji di services (`hitung_limit_*_test.go`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { OpsiLimits, SimpulLimit } from './api'
import TabLimitsNonProp, { DITULIS_LIMIT_NP, RincianLayer } from './components/TabLimitsNonProp'
import TabLimitsProp from './components/TabLimitsProp'
import { GRID_DETAIL, TAB_KUNCI_MATERIAL, syaratQS, syaratSurplus } from './labelsLimitsProp'
import { kontrakRevisi } from './labelsLimitsNP'

const AKAR = __dirname
const LP = readFileSync(join(AKAR, 'components', 'TabLimitsProp.tsx'), 'utf8')
const NP = readFileSync(join(AKAR, 'components', 'TabLimitsNonProp.tsx'), 'utf8')

const OPSI: OpsiLimits = {
  jenisTreaty: [{ id: '10042', nama: 'SURPLUS', kembar: false }],
  kelompokTreaty: [{ id: '10007', nama: 'PROPERTY', kembar: false }],
  mataUang: [
    { id: '10026', nama: 'IDR', kembar: false },
    { id: '10001', nama: 'USD', kembar: false },
  ],
  jenisLayer: [
    { value: 'layer', label: 'layer' },
    { value: 'sublayer', label: 'sublayer' },
  ],
  cover: [{ value: 'risk', label: 'risk' }],
  relasiMataUang: [
    { value: 'OR', label: 'OR' },
    { value: 'AND', label: 'AND' },
  ],
  // ⛔ Label = PROMPT VALUE rule Property (`ReinstatementNote.xml`,
  // 8 Oktober 2026), bukan cermin kodenya. Data uji yang mencerminkan kode
  // membuat uji LULUS untuk kedua perilaku, dan itu yang menyembunyikan
  // cacatnya selama ini.
  catatanReinstatement: [
    { value: 'asamount', label: 'Additional Premium as to amount' },
    { value: 'astime', label: 'Additional Premium as to time' },
  ],
}

describe('tab Limits Prop — sesuai Activity', () => {
  const POHON: SimpulLimit[] = [
    {
      TreatyType: 'QUOTA SHARE',
      Detail: [{ TreatyGroup: 'PROPERTY', TreatyGroupID: '10007', TreatyType: 'QUOTA SHARE', QSPct: '40', RetentionPct: '60', CessionPct: '40' }],
    },
  ]

  it('⛔ Retention % dan Cession % SELALU baca-saja (ditulis LimitCalculation) — teks di samping label, bukan isian', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" opsi={OPSI} />)
    // Bentuk Pega: nilai, spasi, `%` — `100 %`, `0 %`, `100 %` (pemakai, 7 Oktober 2026).
    expect(html).toContain('Retention<span class="tl-persen">60 %</span>')
    expect(html).toContain('Cession to R/I<span class="tl-persen">40 %</span>')
    expect(html).toContain('100% Limit<span class="tl-persen">100 %</span>')
    expect(html).not.toMatch(/<label class="field__label">%<\/label>/)
  })

  it('⭐ Retention % / Cession % langsung dari QS % (`LimitCalculation` [3]) — Treaty Group baru tidak kosong', async () => {
    const { persenQuotaShare } = await import('./components/TabLimitsProp')
    // QS 100 (tangkapan layar Pega): Retention 0 %, Cession 100 %.
    expect(persenQuotaShare({ QSPct: '100' })).toEqual({ retensi: '0', cession: '100' })
    expect(persenQuotaShare({ QSPct: '40' })).toEqual({ retensi: '60', cession: '40' })
    // Tanpa galat float.
    expect(persenQuotaShare({ QSPct: '33.33' })).toEqual({ retensi: '66.67', cession: '33.33' })
    // QS % kosong: nilai tersimpan apa adanya — Pega pun belum menghitungnya.
    expect(persenQuotaShare({ RetentionPct: '', CessionPct: '' })).toEqual({ retensi: '', cession: '' })
    const html = renderToStaticMarkup(
      <TabLimitsProp pohon={[{ TreatyType: 'QUOTA SHARE', Detail: [{ TreatyGroup: 'PROPERTY', TreatyType: 'QUOTA SHARE', QSPct: '100' }] }]} mode="ubah" opsi={OPSI} />,
    )
    expect(html).toContain('Retention<span class="tl-persen">0 %</span>')
    expect(html).toContain('Cession to R/I<span class="tl-persen">100 %</span>')
  })

  it('persen di samping label: angka apa adanya, `0 %` tetap tampil', async () => {
    const { persenSamping } = await import('./components/TabLimitsProp')
    expect(persenSamping('0')).toBe('0 %')
    expect(persenSamping('100')).toBe('100 %')
    expect(persenSamping('')).toBe('%')
  })

  it('⭐ QS % tampil di Treaty Group — jenis induk DISEBAR saat Treaty Type berubah (`TreatyTypeSetIndex`)', async () => {
    const { terapkanJenisTreaty } = await import('./components/TabLimitsProp')
    // Treaty Group ditambahkan SEBELUM Treaty Type dipilih: jenisnya kosong.
    const induk: SimpulLimit = { TreatyType: 'QUOTA SHARE', TreatyTypeID: '10035', Detail: [{ TreatyGroup: 'PROPERTY', TreatyType: '' }] }
    const hasil = terapkanJenisTreaty(induk, 2)
    expect(hasil.ID).toBe('3') // `.ID = .pxListSubscript`
    expect((hasil.Detail as SimpulLimit[])[0]?.TreatyType).toBe('QUOTA SHARE')
    // Sesudahnya QS % dan ketiga persen tampil.
    const html = renderToStaticMarkup(<TabLimitsProp pohon={[hasil]} mode="ubah" opsi={OPSI} />)
    expect(html).toContain('QS %')
    expect(html).toContain('100% Limit<span class="tl-persen">100 %</span>')
    // Tanpa penyebaran, jenis kosong menyembunyikan semuanya.
    const lama = renderToStaticMarkup(<TabLimitsProp pohon={[induk]} mode="ubah" opsi={OPSI} />)
    expect(lama).not.toContain('class="tl-persen"')
    expect(LP).toContain('onUbah(terapkanJenisTreaty(x, indeks))')
  })

  it('⭐ Treaty Type = Name menu Reinsurance Type; Kind of Treaty = SOA Name (Name bila SOA kosong)', async () => {
    const { kindOfTreatyDari } = await import('./components/TabLimitsProp')
    expect(kindOfTreatyDari({ id: '10210', nama: '2020 QS 89M FAC', namaSoa: 'QUOTA SHARE 2020', kembar: false })).toBe('QUOTA SHARE 2020')
    expect(kindOfTreatyDari({ id: '10035', nama: 'QUOTA SHARE', namaSoa: '', kembar: false })).toBe('QUOTA SHARE')
    expect(LP).toContain('namaDari={kindOfTreatyDari}')
    const opsiSoa: OpsiLimits = { ...OPSI, jenisTreaty: [{ id: '10210', nama: '2020 QS 89M FAC', namaSoa: 'QUOTA SHARE 2020', kembar: false }] }
    const pohon: SimpulLimit[] = [{ TreatyType: 'QUOTA SHARE 2020', TreatyTypeID: '10210', Detail: [] }]
    const html = renderToStaticMarkup(<TabLimitsProp pohon={pohon} mode="ubah" opsi={opsiSoa} />)
    // Dropdown berlabel Name — kotak ketik-pilih seperti Ceding (7 Oktober
    // 2026); Kind of Treaty tampil di SEL GRID induknya (bentuk Pega,
    // 8 Oktober 2026), bukan kotak baca-saja tambahan.
    expect(html).toMatch(/<label class="field__label">Treaty Type<\/label><input(?=[^>]*role="combobox")(?=[^>]*value="2020 QS 89M FAC")[^>]*>/)
    expect(html).toContain('<td>QUOTA SHARE 2020</td>')
    expect(html).not.toMatch(/<label class="field__label">Kind of Treaty<\/label>/)
  })

  // ⭐ RALAT 7 Oktober 2026 (perbandingan layar Pega): kolom kotak centang
  // `.Layer` berkepala KOSONG dan kotaknya berketerangan `Auto calculate`
  // (`pyCheckboxCaption`) — bukan kepala `Layer` yang dikarang di sini.
  it('⭐ grid nilai berkepala Currency · Value · (kosong); kotak centang berketerangan Auto calculate', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" opsi={OPSI} />)
    expect(html).toContain('<th scope="col">Currency</th><th scope="col">Value</th><th scope="col"></th>')
    expect(html).not.toContain('<th scope="col">Layer</th>')
    expect(LP).toContain('<label className="trin__cek">')
  })

  it('⭐ nilai awal Add — Note = jenis treaty, ParentID, ID kosong', () => {
    expect(LP).toContain('IOOLimitList: { Note: jenis }')
    expect(LP).toContain('RetentionList: { Note: jenis }')
    expect(LP).toContain('CessionList: { Note: jenis }')
    expect(LP).toContain("DeductionList: { Currency: teksDari(d, 'CurrencyIOOLimit') }")
    expect(LP).toContain("{ ID: '', ParentID: String(indeks + 1), TreatyType: teksDari(l, 'TreatyType') }")
    expect(LP).toContain("setLimits([...limits, { ID: '', Detail: [] }])")
  })

  it('⭐ pemicu grid: 100% Limit (qs, .Layer), Retention (surplus, man, .Layer), Deduction, Reserve', () => {
    expect(LP).toContain("jalankanLC({ ...d, IOOLimitList: p.baris }, 'qs', '', teksDari(p.lama, 'Layer') === 'true')")
    expect(LP).toContain("jalankanLC({ ...d, RetentionList: p.baris }, 'surplus', 'man', teksDari(p.lama, 'Layer') === 'true')")
    expect(LP).toContain('hitungDeduksi({ sts, indeks: p.r')
    expect(LP).toContain('hitungCadangan({ PremiumReservePct')
    // Treaty Group berubah → LimitCalculation(surplus) HANYA untuk SURPLUS.
    expect(LP).toContain("if (jenis === 'SURPLUS') jalankanLC(x, 'surplus', '', true)")
  })

  it('⛔ syarat QS/Surplus PERSIS seperti Pega — keputusan pemakai, tidak dilonggarkan', () => {
    expect(syaratQS('QUOTA SHARE')).toBe(true)
    expect(syaratQS('QUOTA SHARE 2020')).toBe(false)
    expect(syaratSurplus('SURPLUS')).toBe(true)
    expect(syaratSurplus('SPECIAL SURPLUS')).toBe(true)
    expect(syaratSurplus('SURPLUS 2019')).toBe(false)
  })

  it('⛔ grid hasil hitungan baca-saja di mode mana pun', () => {
    expect(GRID_DETAIL.DeductionTotalList?.bacaSaja).toBe(true)
    expect(GRID_DETAIL.AchievementLists?.bacaSaja).toBe(true)
  })

  it('⭐ kunci materialitas (EDMMaterialType = 2): Treaty Group & Add terkunci walau Edit', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" opsi={OPSI} edmJenisMaterial="2" />)
    // ⭐ 8 Oktober 2026 — `Add` di sel kepala kolom tombol grid Pega.
    expect(html).toContain('<button type="button" class="btn btn--sm" disabled="">Add</button>')
    expect(html).toMatch(/<label class="field__label">Treaty Group<\/label><input class="field__input field__input--readonly"/)
    expect(TAB_KUNCI_MATERIAL).toContain('Event Limits')
    const bebas = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" opsi={OPSI} />)
    // ⛔ Combobox sejak 8 Oktober 2026 (*"mengetik harus terasa seperti
    // mencari"*): teks pilihan ada di KOTAKNYA, dan daftar butirnya baru
    // dirender ketika dibuka.
    expect(bebas).toContain('value="PROPERTY"')
  })

  it('⭐ Achievement: Refresh/Quarter Year → GetAchievement di services; Generate Excel .xlsx; Submit → LOG_ACHIEVEMENT', () => {
    expect(LP).toContain('hitungAchievement({')
    expect(LP).toContain('jalankanAchievement(true, v)')
    expect(LP).toContain('jalankanAchievement(false, tahun)')
    expect(LP).toContain("setTahun('') // DT `Reset_DT`")
    expect(LP).toContain('unduhXlsx(')
    // Submit HIDUP sejak 8 Oktober 2026 (keputusan pemakai): satu baris log
    // per baris AchievementLists, pengenal `TreatyIn.ID`.
    expect(LP).not.toMatch(/disabled title=\{ACHIEVEMENT\.kirimMenunggu\}/)
    expect(LP).toContain('catatLogAchievement({')
    expect(LP).toContain('onClick={kirim}')
    expect(LP).toContain('KOLOM_ACH_PARAMETER')
  })

  it('⭐ Class of Business = autocomplete BrowseTreatyBusinessWOType_RD', () => {
    expect(LP).toContain('ambilKelasBisnis(idGrup)')
    expect(GRID_DETAIL.COBList?.kolom[0]?.auto).toBe('kelasBisnis')
  })
})

describe('tab Limits Non-Prop — grid layer → Layers → Summary → Total', () => {
  const LAYER: SimpulLimit = {
    LayerType: 'layer', Layer: '1', LayerPartType: 'layer', LayerPart: '1',
    Currency: 'IDR', Currency2: 'USD', Limit: '1000000', Limit2: '100', Deductible: '500', Deductible2: '0',
    AgregateLimit: '2000000', CurrencyRelation: 'OR', Cover: 'risk', AdjRate: '10', MDPPct: '80', MDPMinPct: '50', ROLPct: '5',
    ReinstatementValue: '1', NoRIPCalculation: 'false', IsCombineMDP: 'true',
    TreatyGroupList: [{ TreatyGroup: 'PROPERTY', TreatyGroupID: '10007', IsROLProfile: 'true', ClassOfBusinessList: [{ ClassOfBusiness: 'FIRE' }] }],
    EgnpiTotalList: [{ Currency: 'IDR', Value: '500000' }],
    PremiumEarnedList: [{ Currency: 'IDR', Value: '50000' }],
    MDPList: [{ Currency: 'IDR', Value: '40000' }],
    MDPMinList: [],
    Reinstatement_List: [{ ReinstatementValue: '1', ReinstatementNote: 'asamount', ReinstatementPct: '100', AdditionalPct: '100', ReinstatementAmount1: '1000000', ReinstatementAmount2: '100' }],
  }
  const AKAR_ISI = {
    LimitSummaryList: [{ Note: 'layer1 of layer1', Limit: '1000000' }],
    Total: { TotalLimitIOONP: [{ Currency: 'IDR', Value: '1000000' }] },
    TotalLimitsROL: '5',
  }

  it('⭐ grid layer: kolom ekspor, baris "layer 1 of layer 1", Summary, Total All Layers', () => {
    const html = renderToStaticMarkup(<TabLimitsNonProp pohon={[LAYER]} akar={AKAR_ISI} egnpi={[]} kurs={[]} mode="lihat" opsi={OPSI} />)
    for (const t of ['100% Limits ( IDR )', 'Deductible ( IDR )', '100% Limits ( USD )', 'Deductible ( USD )',
      'Summary of Limit', 'layer1 of layer1', 'Total All Layers', 'Total 100% Limit', 'Total Deductible',
      'Total Premium Earned', 'Total MDP', 'Total ROL', '1.000.000,00']) {
      expect(html).toContain(t)
    }
    // ⭐ 8 Oktober 2026 — baris GRID Pega (bukan kartu): lima sel kolom
    // `Layers` ekspor = `LayerType · Layer · of · LayerPartType · LayerPart`.
    expect(html).toContain('<td>layer</td><td>1</td><td>of</td><td>layer</td><td>1</td>')
  })

  // ⛔ KELUHAN PEMILIK PROSES 8 Oktober 2026: *"di nonprop bagian
  // ReinstatementNote itu seharusnya mengambil Prompt value bukan standard
  // value"*.
  //
  // Rule Property `ReinstatementNote.xml` (kelas
  // `ASM-FW-GISFW-Data-TreatyInLimits`) memasangkan di `pyPromptTableList`:
  //   asamount → Additional Premium as to amount
  //   astime   → Additional Premium as to time
  //
  // ⚠️ Sampai berkas itu tiba, kami MENYIMPULKAN `as amount` / `as time`
  // dengan dasar "kode gandeng dipecah jadi kata". Melesetnya bukan sedikit:
  // labelnya kalimat penuh yang nol hubungannya dengan ejaan kodenya.
  it('⛔ sel Note menampilkan PROMPT VALUE, bukan kode tersimpannya', () => {
    // Mode LIHAT dan mode UBAH — keduanya, sebab cacatnya justru hanya di
    // mode lihat dan uji yang memeriksa satu mode akan melewatkannya lagi.
    for (const bisaUbah of [false, true]) {
      const html = renderToStaticMarkup(
        <RincianLayer l={LAYER} opsi={OPSI} bisaUbah={bisaUbah} modeUbah={bisaUbah} onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
      )
      expect(html, `bisaUbah=${String(bisaUbah)}`).toContain('Additional Premium as to amount')
      // ⚠️ Label yang benar saja tidak cukup: `<select>` di mode LIHAT juga
      // memuat labelnya sebagai `<option>`, jadi memeriksa teksnya saja
      // meloloskan sel yang diam-diam dapat diubah.
      expect(html.includes('role="combobox"'), `combobox saat bisaUbah=${String(bisaUbah)}`).toBe(bisaUbah)
    }
    const html = renderToStaticMarkup(
      <RincianLayer l={LAYER} opsi={OPSI} bisaUbah={false} modeUbah={false} onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
    )
    // Kode mentahnya tidak boleh bocor sebagai teks sel.
    expect(html).not.toContain('>asamount<')
    // Dan bukan pula kesimpulan lama yang terbukti salah.
    expect(html).not.toContain('as amount<')
    expect(html).not.toContain('tl-kartu')
    // Mode lihat: nol tombol.
    expect(html).not.toContain('add Layer')
    expect(html).not.toContain('Update Total')
  })

  it('mode ubah: add Layer + Update Total; Update Value in List HANYA kontrak revisi', () => {
    const biasa = renderToStaticMarkup(<TabLimitsNonProp pohon={[LAYER]} egnpi={[]} kurs={[]} edmState="0" mode="ubah" opsi={OPSI} />)
    expect(biasa).toContain('>add Layer</button>')
    expect(biasa).toContain('>Update Total</button>')
    expect(biasa).not.toContain('Update Value in List')
    const revisi = renderToStaticMarkup(<TabLimitsNonProp pohon={[LAYER]} egnpi={[]} kurs={[]} edmState="1" mode="ubah" opsi={OPSI} />)
    expect(revisi).toContain('>Update Value in List</button>')
    expect(kontrakRevisi('3')).toBe(true)
    expect(kontrakRevisi('0')).toBe(false)
  })

  it('⭐ rincian Layers: seluruh medan ekspor, Reinstatement, PE, MDP, ROL', () => {
    const html = renderToStaticMarkup(
      <RincianLayer l={LAYER} opsi={OPSI} bisaUbah={false} modeUbah={false} onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
    )
    for (const t of ['Part of', 'Treaty Group', 'ROL Profile', 'Egnpi this layer', 'Cover', 'Currency Relation',
      '100 % Limit', 'Agregate Year Limit', 'Deductible', 'No Reinstatement Premium Calculation', 'Reinstatement',
      'Reinstatement No.', '% Additional Premium', 'Reinstatement Premium Amount IDR (MDP x % Add Premium)',
      'Reinstatement %', 'Reinstatement Amount IDR', 'Reinstatement Amount USD',
      'Adjustment Rate %', 'Premium Earned', 'Min Premium %', 'Min Premium Amt', 'MDP %', 'MDP', '(*) Combine MDP', 'ROL %']) {
      expect(html).toContain(t)
    }
    expect(html).toContain('PROPERTY')
  })

  // ⭐ 9 Oktober 2026 — SELALU tampil walau Limit2 kosong/0 (layar Pega
  // produksi; syarat `.Limit2 != 0` @801378 tidak diikuti, keputusan pemakai).
  it('⭐ kolom Reinstatement Amount USD tetap tampil walau Limit2 kosong', () => {
    const tanpa = renderToStaticMarkup(
      <RincianLayer l={{ ...LAYER, Limit2: '0' }} opsi={OPSI} bisaUbah={false} modeUbah={false} onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
    )
    expect(tanpa).toContain('Reinstatement Amount USD')
    const kosong = renderToStaticMarkup(
      <RincianLayer l={{ ...LAYER, Limit2: '' }} opsi={OPSI} bisaUbah={false} modeUbah={false} onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
    )
    expect(kosong).toContain('Reinstatement Amount USD')
  })

  it('⛔ hanya medan yang Activity TULIS yang digabung; rumus di services', () => {
    expect(DITULIS_LIMIT_NP.adj).toEqual(['PremiumEarnedList', 'ROLPct'])
    // ⭐ 8 Oktober 2026 — `Reinstatement_List` ikut: sesudah MDP berubah,
    // Reinstatement Premium Amount disegarkan dengan rumus yang SAMA
    // (`segarkanReinstatement`, keputusan pemakai).
    expect(DITULIS_LIMIT_NP.mdp).toEqual(['MDPList', 'MDPMinList', 'ROLPct', 'Reinstatement_List'])
    expect(DITULIS_LIMIT_NP.egnpi).toEqual(['EgnpiTotalList'])
    expect(NP).toContain('hitungLimitNP({ aksi, layers: dasar, egnpi, kurs, indeks, baris })')
    expect(NP).not.toMatch(/\* Number\(|parseFloat\(/)
  })

  it('⭐ kunci materialitas: No RIP tetap ikut mode saja', () => {
    expect(NP).toContain("bisaUbah={modeUbah} onUbah={set('NoRIPCalculation')}")
  })

  // ⛔ 8 Oktober 2026: `.ROLPct` di Section `Layers`/`LayersEDM` ber-
  // `pyDisabledNew=always` — dulu terbuka di mode Edit, "1,659" berkoma
  // terbaca nol dan `Total ROL` jadi 0,00.
  it('⛔ ROL % selalu read-only — hanya rumus DetailCalculationROL yang mengisinya', () => {
    expect(NP).toContain("bisaUbah={false} onUbah={set('ROLPct')}")
    expect(NP).not.toContain("bisaUbah={modeUbah} onUbah={set('ROLPct')}")
    const html = renderToStaticMarkup(
      <RincianLayer l={LAYER} opsi={OPSI} bisaUbah modeUbah onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
    )
    const sel = /<label class="field__label">ROL %<\/label><input([^>]*)>/.exec(html)
    expect(sel).not.toBeNull()
    expect(sel?.[1]).toMatch(/readonly/i)
  })
})

// ---------------------------------------------------------------------
// ⛔ ANGKA BIASA TIDAK DIFORMAT — laporan pemilik proses 7 Oktober 2026
// ---------------------------------------------------------------------
// *"pada tab limits di layer nya tidak perlu format nya krn itu memang
// angka biasa"* — kotak `Layer` dan `Part` tampil `1,00`.
//
// Sebabnya satu kata: `desimal ?? 2`. `desimal === null` BUKAN "pakai
// bawaan" — ia berarti presisi kolomnya TIDAK dinyatakan ekspor
// (`pyDecimalPlaces` kosong), dan `Layer`/`Part` bukan uang sama sekali
// melainkan NOMOR URUT.
describe('⛔ kolom tanpa presisi ekspor TIDAK diformat', () => {
  const NP = readFileSync(join(AKAR, 'components', 'TabLimitsNonProp.tsx'), 'utf8')
  const LP = readFileSync(join(AKAR, 'components', 'TabLimitsProp.tsx'), 'utf8')

  /**
   * ⛔ KOMENTAR DISARING LEBIH DULU. Catatan yang MENJELASKAN mengapa
   * `desimal ?? 2` keliru justru harus boleh menyebutnya — kekeliruan
   * yang sama dengan penjaga kode regu di `models/tangga_akseptasi_test.go`.
   */
  const kode = (src: string) =>
    // Blok `/* … */` (termasuk komentar JSX `{/* … */}`) dibuang UTUH:
    // barisnya di tengah blok tidak diawali tanda komentar apa pun, dan
    // penyaring per-baris melewatkannya.
    src
      .replace(/\/\*[\s\S]*?\*\//g, '')
      .split(String.fromCharCode(10))
      .filter((b) => !b.trimStart().startsWith('//'))
      .join(String.fromCharCode(10))

  it('⛔ nol `desimal ?? 2` — itu yang membuat `Layer` tampil `1,00`', () => {
    for (const [nama, src] of [['NonProp', NP], ['Prop', LP]] as const) {
      expect(kode(src), nama).not.toContain('desimal ?? 2')
    }
  })

  // ⚠️ Dan penyaringnya harus TETAP MENGGIGIT kode sungguhan.
  it('⚠️ penyaring komentar tidak membutakan penjaga', () => {
    expect(kode('// desimal ?? 2 dilarang')).not.toContain('desimal ?? 2')
    expect(kode('{/* desimal ?? 2 dilarang */}')).not.toContain('desimal ?? 2')
    expect(kode('  desimal={desimal ?? 2}')).toContain('desimal ?? 2')
  })

  it('⭐ `desimal === null` bercabang ke medan TANPA format', () => {
    expect(NP).toContain('desimal === null ?')
    expect(LP).toContain("m.golongan === 'teks' || m.desimal === null ?")
  })

  // ⚠️ Dan `Layer`/`Part` memang dikirim tanpa presisi — kalau suatu hari
  // keduanya diberi angka, uji di atas tidak lagi melindunginya.
  it('⚠️ Layer dan Part dikirim `desimal={null}`', () => {
    // ⛔ Jangkarnya nilai medannya — sejak 8 Oktober 2026 kedua medan TANPA
    // label (`Layers.xml` @11943, tangkapan layar 31), jadi `label={…}` tak
    // lagi dapat dijadikan jangkar.
    for (const k of ["nilai={teksDari(l, 'Layer')}", "nilai={teksDari(l, 'LayerPart')}"]) {
      const i = NP.indexOf(k)
      expect(i, k).toBeGreaterThan(-1)
      expect(NP.slice(i, NP.indexOf('/>', i)), k).toContain('desimal={null}')
    }
  })

  // ⛔ Golongan `teks` juga dilewati: `>=30% up to < 50%` kebetulan
  // berangka, dan menguraikannya sebagai bilangan merusaknya.
  it('⛔ golongan `teks` nol diurai sebagai bilangan', () => {
    expect(LP).toContain("m.golongan === 'teks'")
  })
})

// ---------------------------------------------------------------------
// ⭐ PEMILIH = DROPDOWN, bukan autocomplete — 7 Oktober 2026
// ---------------------------------------------------------------------
// *"perbaiki di limits non prop agar menjadi dropdown seperti dropdown
// ceding"*. Komponennya SAMA dengan pemilih `Ceding` dan `Source of
// Business` (`DropdownWarisan`), dibungkus `DropdownDaftar` supaya daftar
// yang sudah di tangan tidak diambil ulang tiap render.
describe('⭐ pemilih memakai dropdown yang sama dengan Ceding', () => {
  const berkas = ['TabLimitsNonProp.tsx', 'TabLimitsProp.tsx', 'TabRetensi.tsx', 'TabEgnpi.tsx'] as const

  it('⛔ nol `IsianAuto` tersisa di keempat tab', () => {
    for (const f of berkas) {
      const src = readFileSync(join(AKAR, 'components', f), 'utf8')
      expect(src, f).not.toContain('<IsianAuto')
    }
  })

  it('⭐ keempatnya memakai DropdownWarisan — langsung atau lewat DropdownDaftar', () => {
    for (const f of berkas) {
      const src = readFileSync(join(AKAR, 'components', f), 'utf8')
      expect(src.includes('DropdownDaftar') || src.includes('DropdownWarisan'), f).toBe(true)
    }
  })

  // ⚠️ `ambil` adalah TANGGUNGAN efek `DropdownWarisan`; fungsi baru tiap
  // render mengambil ulang daftarnya tanpa henti. Pembungkusnya membekukan
  // sekali, supaya nol pemanggil perlu mengingatnya.
  it('⚠️ pembungkusnya membekukan `ambil` dengan useCallback', () => {
    const src = readFileSync(join(AKAR, 'components', 'IsianAuto.tsx'), 'utf8')
    expect(src).toContain('export function DropdownDaftar')
    expect(src).toContain('useCallback')
  })
})

// ---------------------------------------------------------------------
// ⭐ SEL ANGKA GRID BERPEMISAH RIBUAN — laporan pemilik proses 7 Okt 2026
// ---------------------------------------------------------------------
// Grid `100% Limit`, `Retention`, dan `Cession to R/I` cabang Prop
// menampilkan `1500000000` apa adanya di kotak isian.
describe('⭐ sel angka grid memakai FieldAngka', () => {
  const NP = readFileSync(join(AKAR, 'components', 'TabLimitsNonProp.tsx'), 'utf8')
  const LP = readFileSync(join(AKAR, 'components', 'TabLimitsProp.tsx'), 'utf8')

  it('⭐ kedua cabang memakai FieldAngka di sel nilai grid', () => {
    for (const [nama, src] of [['NonProp', NP], ['Prop', LP]] as const) {
      expect(src, nama).toContain('<FieldAngka')
    }
  })

  // ⛔ Pembungkus `onBlur` WAJIB tetap ada di cabang Prop: `FieldAngka`
  // memakai `onBlur`-nya sendiri untuk keluar dari mode ketik, sementara
  // peristiwa `change` Pega (`lapor`) harus tetap terkirim. Tanpa
  // pembungkus, rumus yang bergantung pada nilai itu berhenti berjalan.
  it('⛔ peristiwa `change` Pega tetap dilaporkan dari sel angka', () => {
    // ⛔ Jangkarnya `desimal={k.desimal}` — SEL GRID. `FieldAngka` pertama
    // di berkas ini milik pembungkus medan `MedanLimit`, bukan grid, dan
    // menjangkar ke sana memeriksa tempat yang salah.
    const i = LP.indexOf('desimal={k.desimal}')
    expect(i).toBeGreaterThan(-1)
    const sebelum = LP.slice(Math.max(0, i - 700), i)
    // ⭐ 8 Oktober 2026 — `PemicuUbah` (lepas fokus DAN berubah), bukan `onBlur` mentah.
    expect(sebelum).toContain('<PemicuUbah')
    expect(sebelum).toContain('lapor({ jenis:')
  })

  // ⚠️ Dan kolom tanpa presisi ekspor TETAP di luar — pelajaran `1,00`
  // dari kartu Layers.
  it('⚠️ hanya kolom ber-golongan angka DAN berdesimal yang diformat', () => {
    expect(LP).toContain("k.golongan !== 'teks' && k.desimal !== null")
  })
})
