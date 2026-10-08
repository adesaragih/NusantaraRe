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
  catatanReinstatement: [{ value: 'asamount', label: 'asamount' }],
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
    // 2026); Kind of Treaty tampil baca-saja.
    expect(html).toMatch(/<label class="field__label">Treaty Type<\/label><input(?=[^>]*role="combobox")(?=[^>]*value="2020 QS 89M FAC")[^>]*>/)
    expect(html).toMatch(/<label class="field__label">Kind of Treaty<\/label><input class="field__input field__input--readonly" type="text" readonly="" value="QUOTA SHARE 2020"/)
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
    expect(html).toMatch(/<button type="button" class="btn btn--sm tl-tambah" disabled="">(?:(?!<\/button>).)*Add<\/button>/)
    expect(html).toMatch(/<label class="field__label">Treaty Group<\/label><input class="field__input field__input--readonly"/)
    expect(TAB_KUNCI_MATERIAL).toContain('Event Limits')
    const bebas = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" opsi={OPSI} />)
    expect(bebas).toContain('<option value="10007" selected="">PROPERTY</option>')
  })

  it('⭐ Achievement: Refresh/Quarter Year → GetAchievement di services; Generate Excel .xlsx; Submit menunggu jalur tulis', () => {
    expect(LP).toContain('hitungAchievement({')
    expect(LP).toContain('jalankanAchievement(true, v)')
    expect(LP).toContain('jalankanAchievement(false, tahun)')
    expect(LP).toContain("setTahun('') // DT `Reset_DT`")
    expect(LP).toContain('unduhXlsx(')
    expect(LP).toMatch(/disabled title=\{ACHIEVEMENT\.kirimMenunggu\}/)
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
    Reinstatement_List: [{ ReinstatementValue: '1', ReinstatementPct: '100', AdditionalPct: '100', ReinstatementAmount1: '1000000', ReinstatementAmount2: '100' }],
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
    // Judul kartu layer = kolom Layers ekspor: `LayerType Layer of LayerPartType LayerPart`.
    expect(html).toContain('<span class="tl-kartu__judul">layer 1 of layer 1</span>')
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

  it('⛔ kolom Reinstatement Amount USD hanya bila Limit2 ≠ 0', () => {
    const tanpa = renderToStaticMarkup(
      <RincianLayer l={{ ...LAYER, Limit2: '0' }} opsi={OPSI} bisaUbah={false} modeUbah={false} onUbah={() => undefined} hitung={() => undefined} onGantiGrup={() => undefined} />,
    )
    expect(tanpa).not.toContain('Reinstatement Amount USD')
  })

  it('⛔ hanya medan yang Activity TULIS yang digabung; rumus di services', () => {
    expect(DITULIS_LIMIT_NP.adj).toEqual(['PremiumEarnedList', 'ROLPct'])
    expect(DITULIS_LIMIT_NP.mdp).toEqual(['MDPList', 'MDPMinList', 'ROLPct'])
    expect(DITULIS_LIMIT_NP.egnpi).toEqual(['EgnpiTotalList'])
    expect(NP).toContain('hitungLimitNP({ aksi, layers: dasar, egnpi, kurs, indeks, baris })')
    expect(NP).not.toMatch(/\* Number\(|parseFloat\(/)
  })

  it('⭐ kunci materialitas: No RIP dan ROL % tetap ikut mode saja', () => {
    expect(NP).toContain("bisaUbah={modeUbah} onUbah={set('NoRIPCalculation')}")
    expect(NP).toContain("bisaUbah={modeUbah} onUbah={set('ROLPct')}")
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
    // ⛔ Jangkarnya `label={LIMITS_NP.layer}` — `LIMITS_NP.layers` (judul
    // bagian) cocok lebih dulu bila yang dicari hanya namanya.
    for (const k of ['label={LIMITS_NP.layer}', 'label={LIMITS_NP.part}']) {
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
    expect(sebelum).toContain('onBlur={() => {')
    expect(sebelum).toContain('lapor({ jenis:')
  })

  // ⚠️ Dan kolom tanpa presisi ekspor TETAP di luar — pelajaran `1,00`
  // dari kartu Layers.
  it('⚠️ hanya kolom ber-golongan angka DAN berdesimal yang diformat', () => {
    expect(LP).toContain("k.golongan !== 'teks' && k.desimal !== null")
  })
})
