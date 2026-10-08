// Uji mode Edit/View, tab Limits proporsional, dan dropdown kepala.
//
// Dirender statis (`react-dom/server`) — BUKAN pengganti melihat layar di
// peramban. `<details>` dirender tertutup, tetapi isinya tetap di HTML.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { SimpulLimit } from './api'
import TabGridWarisan from './components/TabGridWarisan'
import TabLimitsProp, { formatLimit } from './components/TabLimitsProp'
import { GRID_DETAIL, TAB_DETAIL_LIMIT } from './labelsLimitsProp'
import { modeUntuk } from './pages/DaftarKontrakTreatyIn'

const AKAR = __dirname
const FORM = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')

const POHON: SimpulLimit[] = [
  {
    TreatyType: 'QUOTA SHARE',
    TreatyTypeID: '10042',
    Detail: [
      {
        TreatyGroup: 'PROPERTY',
        TreatyType: 'QUOTA SHARE',
        QSPct: '25',
        RetentionPct: '75',
        COBList: [{ ClassOfBusiness: 'FIRE' }],
        IOOLimitList: [
          // `.Layer` — kotak centang bernilai teks (terukur: `true`/`false`).
          { Currency: 'IDR', Value: '1500000000', Note: 'QUOTA SHARE', Layer: 'true' },
          { Currency: 'USD', Value: '100000', Note: 'LAIN', Layer: 'true' },
        ],
        RetentionList: [],
        CessionList: [{ Currency: 'IDR', Value: '375000000' }],
      },
    ],
  },
]

const kontrol = (html: string) => [...html.matchAll(/<(input|textarea|select)\b[^>]*>/g)].map((m) => m[0])
const tombol = (html: string, label: string) => (html.match(new RegExp(`>${label}</button>`, 'g')) ?? []).length

describe('daftar → mode form', () => {
  it('⭐ Edit membuka mode ubah; View, Copy, Revision membuka mode lihat', () => {
    expect(modeUntuk('Edit')).toBe('ubah')
    for (const a of ['View', 'Copy', 'Revision']) expect(modeUntuk(a)).toBe('lihat')
  })

  it('rute meneruskan mode; Add (kontrak baru) membuka mode ubah', () => {
    expect(RUTE).toContain('mode={mode}')
    expect(RUTE).toMatch(/onTambah=\{\(\) => \{\s*setMode\('ubah'\)/)
  })

  it('⛔ mode lihat mematikan isian lewat fieldset — strip tab dan Close di luarnya', () => {
    expect((FORM.match(/<fieldset className="trin__mode" disabled=\{!bisaUbah\}/g) ?? []).length).toBe(1)
    // ⭐ Fieldset isi tab: mati di mode lihat KECUALI tab Information & Submit
    // kontrak revisi (`TreatyInfoSubmit` cell 9 dan 21) — 7 Oktober 2026.
    expect(FORM).toContain('disabled={!bisaUbah && !(revisi && tabTampil === TAB_REVISI)}')
    const strip = FORM.indexOf('<StripTab tab={tab}')
    const tutup = FORM.lastIndexOf('</fieldset>', strip)
    expect(tutup).toBeGreaterThan(0)
    expect(FORM.lastIndexOf('onClick={onKembali}')).toBeGreaterThan(FORM.lastIndexOf('</fieldset>'))
  })
})

describe('Accounting Mode / Bordereaux — nilai tersimpan, label dari services', () => {
  it('⛔ dropdown diisi nilai TERSIMPAN, bukan label hasil terjemahan', () => {
    expect(FORM).toContain('setPembukuan(k.caraPembukuanAsli)')
    expect(FORM).toContain('setPembukuanNonProp(k.caraPembukuanNonPropAsli)')
    expect(FORM).toContain('setBordereaux(k.bordereauxAsli)')
  })

  it('⛔ pilihannya pasangan nilai↔label dari services, bukan nilai mentah sebagai label', () => {
    expect(FORM).toContain('opsi={opsi?.caraPembukuan ?? []}')
    expect(FORM).toContain('opsi={opsi?.caraPembukuanNonProp ?? []}')
    expect(FORM).toContain('opsi={opsi?.bordereaux ?? []}')
    expect(FORM).not.toMatch(/caraPembukuanNilai\.map|bordereauxNilai\.map/)
  })
})

describe('tab Limits proporsional — tiga tingkat dari ekspor', () => {
  it('⛔ nol medan layer non-prop (Cover, MDP, ROL, padanan 70,6%)', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="lihat" />)
    for (const salah of ['Cover', 'Currency Relation', 'MDP %', 'ROL %', 'padanan 70,6%']) expect(html).not.toContain(salah)
  })

  it('tingkat 1–3 hadir: Kind of Treaty → Treaty Type + Treaty Group → DetailLimits', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="lihat" />)
    for (const t of ['Kind of Treaty', 'QUOTA SHARE', 'Treaty Type', 'PROPERTY', 'Class of Business', 'FIRE', '100% Limit', 'Retention', 'Cession to R/I']) {
      expect(html).toContain(t)
    }
    // Sebelas tab DetailLimits, urutan ekspor.
    expect(TAB_DETAIL_LIMIT.map((x) => x.judul)).toEqual([
      'Event Limits', 'Deduction In A', 'Deduction', 'Reserve', 'Experience Premium Refund', 'PLA',
      'Cash Loss Limit', 'Claim Cooperation', 'LPC', 'EPI', 'Achievement',
    ])
    for (const x of TAB_DETAIL_LIMIT) expect(html).toContain(`>${x.judul}<`)
  })

  it('⛔ QS % hanya untuk QUOTA SHARE; Lines hanya untuk jenis surplus', () => {
    const qs = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="lihat" />)
    expect(qs).toContain('QS %')
    expect(qs).not.toContain('>Lines<')
    const sp: SimpulLimit[] = [{ TreatyType: 'SURPLUS', Detail: [{ TreatyGroup: 'X', TreatyType: 'SURPLUS', Surplus: '3' }] }]
    const html = renderToStaticMarkup(<TabLimitsProp pohon={sp} mode="lihat" />)
    expect(html).toContain('Lines')
    expect(html).not.toContain('QS %')
  })

  it('⛔ desimal per kolom dari ekspor, nol di ekor dipertahankan', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="lihat" />)
    expect(html).toContain('1.500.000.000,00')
    expect(html).toContain('375.000.000,00')
    expect(formatLimit('uang', 2, '1')).toBe('1,00')
    expect(formatLimit('uang', null, '43994125000.0')).toBe('43.994.125.000')
    expect(formatLimit('uang', -1, '1.123456789')).toBe('1,123456789')
  })

  it('⛔ syarat SEL: Layer hanya pada baris ber-Note QUOTA SHARE', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="lihat" />)
    // Kotak centang `Layer` HANYA pada baris ber-Note QUOTA SHARE.
    expect((html.match(/type="checkbox"/g) ?? []).length).toBe(1)
    // Keterangannya `pyCheckboxCaption` `Auto calculate` (layar Pega).
    expect(html).toMatch(/type="checkbox" aria-label="Auto calculate"[^>]*checked=""/)
    expect(html).toContain('Auto calculate</label>')
  })

  it('⭐ mode ubah: Add/Delete di tiap tingkat dan grid ber-Add; isian dapat diketik', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" />)
    // 9 grid ber-Add di ekspor (IOO, Retention, Cession, Deduction, Reserve,
    // PLA, CashLoss, ClaimCoop, EPI). Yang terlihat di render pertama:
    // Kind of Treaty + Treaty Group + tiga grid di atas strip tab — grid di
    // dalam tab dirender ketika tabnya dipilih (tab aktif: Event Limits).
    const gridAdd = Object.values(GRID_DETAIL).filter((g) => g.tambah).length
    expect(gridAdd).toBe(9)
    expect(tombol(html, 'Add')).toBe(2 + 3)
    expect(tombol(html, 'Delete')).toBeGreaterThanOrEqual(2)
    const bebas = kontrol(html).filter((t) => !/readOnly=""|readonly=""|disabled=""/i.test(t))
    expect(bebas.length).toBeGreaterThan(5)
    // ⛔ RALAT 7 Oktober 2026 — pernyataan ini DULU menuntut nilai MENTAH
    // di kotak isian (`value="1500000000"`), dengan alasan "isian memegang
    // nilai TERSIMPAN, bukan terjemahan tampil".
    //
    // Alasannya benar untuk TEKS, dan keliru untuk ANGKA. Pemilik proses
    // melaporkan grid `100% Limit` / `Retention` / `Cession to R/I`
    // menampilkan `1500000000` apa adanya. `FieldAngka` memisahkan kedua
    // bentuk itu: yang TAMPIL berpemisah ribuan, yang DIKIRIM tetap bentuk
    // kabel — jadi tuntutan lamanya tidak hilang, ia pindah ke
    // `inti/frontend/lib/angkaKetik.test.ts` (`keKabelAngka`).
    expect(html).toContain('1.500.000.000')
    expect(html).not.toContain('value="1500000000"')
  })

  it('⛔ mode lihat: nol tombol Add/Delete, nol isian dapat diketik', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="lihat" />)
    expect(tombol(html, 'Add')).toBe(0)
    expect(tombol(html, 'Delete')).toBe(0)
    expect(kontrol(html).filter((t) => !/readOnly=""|readonly=""|disabled=""/i.test(t))).toEqual([])
  })

  it('⛔ Class of Business tanpa Add — tombolnya `… && 1=2` di ekspor', () => {
    expect(GRID_DETAIL.COBList?.tambah).toBe(false)
  })
})

describe('grid tab lain — Add/Delete hidup di mode ubah', () => {
  const props = { judul: 'EGNPI', kolom: ['A', 'B'], baris: [['1', '2']], petunjukKosong: '-' }

  it('mode ubah: Add, Delete, sel dapat diketik', () => {
    const html = renderToStaticMarkup(<TabGridWarisan {...props} mode="ubah" />)
    expect(tombol(html, 'Add')).toBe(1)
    expect(tombol(html, 'Delete')).toBe(1)
    expect(html).toMatch(/<input[^>]*value="1"/)
  })

  it('mode lihat: tanpa tombol, sel teks', () => {
    const html = renderToStaticMarkup(<TabGridWarisan {...props} mode="lihat" />)
    expect(tombol(html, 'Add')).toBe(0)
    expect(html).not.toContain('<input')
  })

  it('⛔ riwayat tanpa Add — barisnya lahir dari Submit', () => {
    const html = renderToStaticMarkup(<TabGridWarisan {...props} mode="ubah" bisaTambah={false} />)
    expect(tombol(html, 'Add')).toBe(0)
    // ⛔ ALAMATNYA DIGANTI 6 Oktober 2026, maksudnya TETAP.
    //
    // Bunyi sebelumnya menuntut `bisaTambah={false}` di cabang tab
    // `Information & Submit`, sebab tab itu dahulu merender grid riwayat.
    // Pukul 16:26 cabangnya diganti `<TabInfoSubmit …/>` — form Information
    // + Comment dari `Section/TreatyInfoSubmit.xml` — dan grid riwayat
    // keluar dari tab itu. Riwayat kini HANYA di `PanelHistory` di kaki
    // layar. Jadi yang dijaga: tab itu tidak merender grid riwayat lagi,
    // dan `PanelHistory` tidak punya tombol apa pun.
    const i = FORM.indexOf("tabTampil === 'Information & Submit'")
    expect(i).toBeGreaterThan(-1)
    const j = FORM.indexOf('tabTampil ===', i + 1)
    expect(FORM.slice(i, j === -1 ? undefined : j)).not.toContain('KOLOM_CATATAN')
    const RIWAYAT = readFileSync(join(AKAR, 'components', 'PanelHistory.tsx'), 'utf8')
    expect(RIWAYAT).not.toContain('<button')
  })
})

describe('Treaty Type — dropdown `BrowseReinsuranceType_RD`', () => {
  const OPSI_JENIS = [
    { id: '10042', nama: 'SURPLUS', kembar: false },
    { id: '10035', nama: 'QUOTA SHARE', kembar: false },
  ]

  it('⭐ dropdown: kode tersimpan terpilih, label `.Note` — dapat diketik seperti Ceding', () => {
    const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON} mode="ubah" opsi={{ jenisTreaty: OPSI_JENIS, kelompokTreaty: [], mataUang: [] }} />)
    // Kotak ketik-pilih (`DropdownWarisan`) memperlihatkan Name pilihan tersimpan.
    expect(html).toMatch(/<label class="field__label">Treaty Type<\/label><input(?=[^>]*role="combobox")(?=[^>]*value="SURPLUS")[^>]*>/)
    expect(html).not.toContain('teks pilihannya tidak ada di ekspor')
    // ⛔ Bukan lagi `<select>` untuk Treaty Type.
    expect(html).not.toMatch(/<option value="10042" selected="">SURPLUS<\/option>/)
  })

  it('⛔ pilihannya diminta dari server, bukan ditulis di layar', () => {
    const sumber = readFileSync(join(AKAR, 'components', 'TabLimitsProp.tsx'), 'utf8')
    expect(sumber).toContain('ambilOpsiLimits()')
    expect(sumber).not.toMatch(/value: 'SURPLUS'|value: 'QUOTA SHARE'/)
  })
})

it('Treaty Type mode lihat: label `.Note` baca-saja, bukan kode', () => {
  const html = renderToStaticMarkup(
    <TabLimitsProp pohon={POHON} mode="lihat" opsi={{ jenisTreaty: [{ id: '10042', nama: 'SURPLUS', kembar: false }], kelompokTreaty: [], mataUang: [] }} />,
  )
  expect(html).toMatch(/value="SURPLUS"[^>]*readOnly=""|readOnly=""[^>]*value="SURPLUS"/i)
})

describe('tab Limits seperti layar Pega (lampiran 6 Oktober 2026)', () => {
  const OPSI = {
    jenisTreaty: [{ id: '10042', nama: 'SURPLUS', kembar: false }],
    kelompokTreaty: [{ id: '10007', nama: 'PROPERTY', kembar: false }],
    mataUang: [
      { id: '10026', nama: 'IDR', kembar: false },
      { id: '10001', nama: 'USD', kembar: false },
    ],
  }
  const POHON_G: SimpulLimit[] = [
    { TreatyType: 'SURPLUS', TreatyTypeID: '10042', Detail: [{ TreatyGroup: 'PROPERTY', TreatyGroupID: '10007', IOOLimitList: [{ Currency: 'USD', Value: '0' }] }] },
  ]
  const html = renderToStaticMarkup(<TabLimitsProp pohon={POHON_G} mode="ubah" opsi={OPSI} />)

  it('Treaty Group dropdown `BrowseTreatyGroup_RD` — kode terpilih', () => {
    expect(html).toMatch(/<option value="10007" selected="">PROPERTY<\/option>/)
  })
  it('sel mata uang grid = dropdown CURRENCY bernilai `.Currency`', () => {
    expect(html).toMatch(/<option value="USD" selected="">USD<\/option>/)
    expect(html).toContain('<option value="IDR">IDR</option>')
  })
  it('baris grid: Remove; baris Kind of Treaty/Treaty Group: Delete', () => {
    expect(html).toContain('>Remove</button>')
    expect(html).toContain('>Delete</button>')
  })
  it('⭐ Kind of Treaty punya isian sendiri (autocomplete @586908); Treaty Group bernomor', () => {
    expect(html).toMatch(/<label[^>]*>Kind of Treaty<\/label>/)
    expect(html).toMatch(/<span class="tl-nomor">1<\/span><span class="tl-kartu__judul">PROPERTY<\/span>/)
  })
})
