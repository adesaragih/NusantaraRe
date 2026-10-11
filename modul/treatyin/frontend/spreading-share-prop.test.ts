// Tab Share Prop — rincian Detail (`Section/DetailShare.xml`) disamakan
// dengan layar Pega pemakai.
//
// ⛔ RALAT 8 Oktober 2026 (sore) — koreksi pemilik proses atas keputusan pagi
// ("dropdown Spreading Type selalu tampil, grid manual dicabut"), dengan
// tangkapan layar Pega dan rancangan Section: *"spread nya bisa ditambah, cek
// lagi xml nya"*. `DetailShare.xml` punya DUA grid `.SpreadingList`:
//
//   L12  `.SpreadingTypeID != ''` — dropdown Spreading Type + grid BACA
//        (`Reins Type · Pct`, RD `PROPORTIONALARRG`, `FetchQSfromMaster`);
//   L16  `.SpreadingTypeID == ''` — grid MANUAL: `Reins Type` (dropdown RD
//        arrangement) · `Pct Share` (`pxNumber`), keduanya change →
//        `SetSpreadName`; `Add` di kepala / `Delete` per baris →
//        `AddDelSpreadingTreatyin` (tambah `Pct = 0` / hapus, lalu
//        `CountTotalPctSpead`); `pyEditAction = SpreadingTPDtl` — tiap baris
//        dibuka ke `.BreakDownSprdList` (Spread · Currency · Share (%) · Amount).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { createElement } from 'react'
import { describe, expect, it } from 'vitest'

import type { SimpulLimit } from './api'
import { selAngka } from './components/angka'
import TabShareProp, { jumlahPct } from './components/TabShareProp'
import { PenyediaHalaman } from './halaman'
import { DETAIL_SHARE } from './labelsShareProp'

const SRC = readFileSync(join(__dirname, 'components', 'TabShareProp.tsx'), 'utf8')
const i = SRC.indexOf('function RincianDetailShare')
const RINCIAN = SRC.slice(i, SRC.indexOf('export default function TabShareProp', i))

const nil = (v: string) => [{ Currency: 'IDR', Value: v }]

/** Render tab Share Prop dengan SATU Detail — semua rincian ada di DOM (`hidden`). */
function render(detail: Record<string, unknown>, mode: 'ubah' | 'lihat' = 'ubah'): string {
  const halaman = {
    RNMShareP: '25', BrokeragePercentP: '0', OptionLimit: '1',
    Limits: [{ TreatyType: 'QUOTA SHARE', Detail: [{ TreatyGroup: 'PROPERTY', TreatyGroupID: '10007', RNMShare: '25', RNMShareList: nil('225000000000'), ...detail }] }],
  }
  return renderToStaticMarkup(
    createElement(PenyediaHalaman, { penampung: { halaman, ubah: () => undefined } } as never, createElement(TabShareProp, { mode } as never)),
  )
}

describe('⭐ DUA cabang Spreading — persis DetailShare.xml', () => {
  // ⛔ DIBALIK 8 Oktober 2026 (sore) — dahulu uji ini menuntut grid BACA.
  //
  // `DetailShare.xml` memang memisahkan: `.SpreadingTypeID != ''` → L12 tanpa
  // tombol. Pemilik proses meminta spread tetap dapat DITAMBAH walau
  // Spreading Type sudah dipilih (*"konsepnya harus sama seperti yang di prop
  // juga yang di mana spread nya bisa ditambah"*), jadi grid manual dipakai
  // di KEDUA cabang — sama dengan layar Non-Prop.
  //
  // ⭐ Yang TIDAK berubah: dropdown Spreading Type dan `Spreading Total Pct`
  // tetap milik cabang ini, dan `Total Share Pct` tetap milik cabang manual.
  // ⭐ 9 Oktober 2026 — DIGANTI keputusan pemakai: *"yang tidak bisa ditambah
  // itu khusus untuk spreading lama"*. Cabang bernama kembali PERSIS Pega
  // (L12): grid BACA, nol Add/Delete, Pct tidak dapat diubah.
  it('⛔ Spreading Type TERISI (spreading LAMA) → dropdown + grid BACA, nol Add/Delete', () => {
    const html = render({
      SpreadingTypeID: '10227',
      SpreadingType: '2023 QS 150M TRT',
      // ⚠️ Baris hasil `FetchQSfromMaster` SELALU membawa `ReinsTypeID`
      // (`fetchQS`: `"ReinsTypeID": a.ReinsTypeID`). Data uji tanpa ID
      // membuat selnya jatuh ke `Choose` dan menyembunyikan bahwa nama
      // baris nyata tetap terbaca.
      SpreadingList: [{ ReinsTypeID: '10004', ReinsTypeName: 'QS (OR)', Pct: '40' }],
    })
    expect(html).toContain(`>${DETAIL_SHARE.jenisSpreading}</label>`)
    expect(html).toContain(`>${DETAIL_SHARE.totalSpreadingPct}`)
    expect(html).not.toContain(`>${DETAIL_SHARE.tambahSpread}</button>`)
    expect(html).not.toContain(`>${DETAIL_SHARE.hapusSpread}</button>`)
    for (const k of DETAIL_SHARE.kolomSpreading) expect(html).toContain(`>${k}</th>`)
    expect(html).toContain('>QS (OR)</td>')
    expect(html).not.toMatch(/<select[^>]*aria-label="Reins Type 1"/)
    expect(html).not.toMatch(/<input[^>]*aria-label="Pct Share 1"/)
    expect(html).not.toContain(DETAIL_SHARE.totalSharePct)
  })

  it('Spreading Type KOSONG → grid MANUAL (L16): Reins Type · Pct Share, Add, Delete, Total Share Pct', () => {
    const html = render({
      SpreadingTypeID: '',
      SpreadingList: [{ ReinsTypeID: '10500', ReinsTypeName: '2025 QS 181M TRT', Pct: '25', BreakDownSprdList: [] }],
      SpreadingTotalPct: '25',
    })
    expect(html).not.toContain(`>${DETAIL_SHARE.jenisSpreading}</label>`)
    for (const k of DETAIL_SHARE.kolomSpreadingManual) expect(html).toContain(`>${k}</th>`)
    expect(html).toContain(`>${DETAIL_SHARE.tambahSpread}</button>`)
    expect(html).toContain(`>${DETAIL_SHARE.hapusSpread}</button>`)
    expect(html).toMatch(/<select class="field__input" aria-label="Reins Type 1">/)
    expect(html).toMatch(/aria-label="Pct Share 1" value="25"/)
    expect(html).toContain(`${DETAIL_SHARE.totalSharePct} 25,00%`)
  })

  it('baris manual dibuka ke `SpreadingTPDtl` — Spread · Currency · Share (%) · Amount', () => {
    const html = render({
      SpreadingTypeID: '',
      SpreadingList: [{
        ReinsTypeID: '10500', ReinsTypeName: '2025 QS 181M TRT', Pct: '25',
        BreakDownSprdList: [
          { ReinsName: 'QS (OR)', Currency: 'IDR', SharePct: '40', Amount: '22500000000' },
          { ReinsName: 'QS (R/I)', Currency: 'IDR', SharePct: '60', Amount: '33750000000' },
        ],
      }],
    })
    for (const k of DETAIL_SHARE.kolomPecahan) expect(html).toContain(`>${k}</th>`)
    expect(html).toContain('>QS (OR)</td>')
    expect(html).toContain('>22.500.000.000,00</td>')
    expect(html).toContain('>40,00</td>')
  })

  it('mode lihat: grid manual BACA — nol select, nol Add/Delete', () => {
    const html = render({ SpreadingTypeID: '', SpreadingList: [{ ReinsTypeID: '1', ReinsTypeName: '2025 QS 181M TRT', Pct: '25' }] }, 'lihat')
    expect(html).toContain('>2025 QS 181M TRT</td>')
    expect(html).not.toContain(`>${DETAIL_SHARE.tambahSpread}</button>`)
    expect(html).not.toContain('aria-label="Reins Type 1"')
  })

  it('Reins Type / Pct Share → `SetSpreadName` (aksi `sebar-nama`)', () => {
    expect((RINCIAN.match(/onRumus\('sebar-nama'/g) ?? []).length).toBe(2)
    // ⛔ KEDUA dropdown memanggil RD TANPA Treaty Group — PENYIMPANGAN yang
    // diputuskan pemilik proses 8 Oktober 2026 (*"gimana pun caranya asal itu
    // ada isinya"*).
    //
    // `Section/DetailShare.xml` MENGIRIM `TreatyGroupID`; mengikutinya
    // mengosongkan dropdown untuk kontrak yang susunannya ADA di
    // `PROPORTIONALARRG`, hanya terdaftar di Treaty Group lain.
    //
    // ⚠️ Dibuang di DROPDOWN dan di PENCARIAN sekaligus (`fetchQS`,
    // `hitung_share_prop.go`). Membuangnya hanya di dropdown — yang sempat
    // terjadi — membuat Spreading Type terpilih sementara grid spreadingnya
    // tetap kosong: setengah penyimpangan terbaca seperti berhasil.
    expect(RINCIAN).toContain("useIndukSpreading(commencement, !terkunci && manual)")
    expect(RINCIAN).toContain("useIndukSpreading(commencement, !terkunci && !manual)")
    expect(RINCIAN).not.toContain("teksDari(d, 'TreatyGroupID'), commencement")
    // `TreatyDescID` TETAP dikirim — nol pemanggil yang menghilangkannya.
    expect(SRC).toContain("ambilIndukSpreading('', DESC_SPREADING_PROP, commencement)")
    expect(SRC).toContain("const DESC_SPREADING_PROP = '10001'")
  })

  it('pilihan Spreading Type → `FetchQSfromMaster` (aksi `spreading`); dikosongkan → nama ikut kosong', () => {
    expect(RINCIAN).toContain("onRumus('spreading', { ...d, SpreadingTypeID: v, ...(v === '' ? { SpreadingType: '' } : {}) })")
  })

  it('`CountTotalPctSpead` — Σ Pct, ketikan berkoma ikut terbaca', () => {
    const r = (Pct: string) => ({ Pct }) as SimpulLimit
    expect(jumlahPct([r('25'), r('10,5'), r('0')])).toBe('35.5')
    expect(jumlahPct([])).toBe('0')
  })

  it('`Total Spreading Pct :` lalu Value Spreading OR / R/I (cabang L12)', () => {
    expect(DETAIL_SHARE.totalSpreadingPct).toBe('Total Spreading Pct :')
    const t = RINCIAN.indexOf('DETAIL_SHARE.totalSpreadingPct')
    expect(RINCIAN.indexOf('DETAIL_SHARE.sebaranOR', t)).toBeGreaterThan(t)
    expect(RINCIAN.indexOf('DETAIL_SHARE.sebaranRI', t)).toBeGreaterThan(t)
  })
})

// Perbandingan layar Pega 7 Oktober 2026 — format dari ekspor.
describe('⭐ tampilan tab Share Prop = layar Pega', () => {
  it('`% RNM Share` Treaty Group: `pxNumber` 2 desimal TANPA %, lalu ShareNote — `25,00 of 80% of 100%`', () => {
    expect(selAngka(['uang', 2], '25') + ' of 80% of 100%').toBe('25,00 of 80% of 100%')
    expect(SRC).toContain("{selAngka(['uang', 2], teksDari(d, 'RNMShare'))}")
    expect(SRC).not.toContain("selAngka(['persen', 2], teksDari(d, 'RNMShare'))")
  })

  it('grid rincian Detail (RNM Share, Value Spreading OR/R/I) berkepala kolom kedua KOSONG', () => {
    // RNM Share + Value Spreading OR/R/I di cabang L12, dan OR/R/I lagi di cabang L16.
    expect((SRC.match(/judulNilai=""/g) ?? []).length).toBe(5)
  })

  it('`% RNM Share` / `% Brokerage` kepala: tanpa simbol, placeholder `%`', () => {
    expect((SRC.match(/desimal=\{2\} placeholder="%"/g) ?? []).length).toBe(2)
    expect(SRC).not.toMatch(/persenRnmShare\} value=\{rnmShare\} desimal=\{2\} persen/)
  })
})
