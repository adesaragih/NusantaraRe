// Uji tab Maximum Retention (non-prop) — bentuk layar lawan ekspor Pega.
//
// Dirender statis (`react-dom/server`). RUMUSNYA diuji di services
// (`hitung_retensi_test.go`); di sini hanya: apakah yang ekspor perlihatkan
// ADA di layar, dan apakah yang ekspor kunci TETAP terkunci.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisRetensi, OpsiLimits } from './api'
import TabRetensi from './components/TabRetensi'
import { TOTAL_RETENSI } from './labels'
import { DESIMAL_RETENSI, KOLOM_GRID_RETENSI, RETENSI, totalTerkunci } from './labelsRetensi'

const AKAR = __dirname
const SRC = readFileSync(join(AKAR, 'components', 'TabRetensi.tsx'), 'utf8')

const OPSI: OpsiLimits = {
  jenisTreaty: [],
  kelompokTreaty: [{ id: '10007', nama: 'PROPERTY', kembar: false }],
  mataUang: [{ id: '10026', nama: 'IDR', kembar: false }],
}

const BARIS: BarisRetensi[] = [
  {
    ID: '1',
    TreatyGroup: 'PROPERTY',
    TreatyGroupID: '10007',
    Currency: 'IDR',
    CurrencyID: '10026',
    Amount: '3500000000',
    ClassOfBusiness: '',
    ClassOfBusinessID: '',
    Note: 'catatan retensi',
  },
]

function render(mode: 'lihat' | 'ubah', edm = '') {
  return renderToStaticMarkup(
    <TabRetensi baris={BARIS} opsiAwal={OPSI} edmJenisMaterial={edm} mode={mode} />,
  )
}

describe('tab Maximum Retention — bentuk dari ekspor', () => {
  it('⭐ panel dan ketiga kolom grid ekspor ada', () => {
    const html = render('lihat')
    expect(html).toContain(RETENSI.panel)
    for (const k of KOLOM_GRID_RETENSI) {
      expect(html, k).toContain(k)
    }
  })

  it('⭐ panel total ada, berikut judul kolomnya', () => {
    const html = render('lihat')
    expect(html).toContain(TOTAL_RETENSI.judul)
    expect(html).toContain(TOTAL_RETENSI.kolomNilai)
  })

  it('⛔ Add dan Delete hanya di mode ubah', () => {
    expect(render('ubah')).toContain(RETENSI.tambah)
    expect(render('ubah')).toContain(RETENSI.hapus)
    expect(render('lihat')).not.toContain(RETENSI.tambah)
    expect(render('lihat')).not.toContain(RETENSI.hapus)
  })

  // ⚠️ SATU tombol, bukan dua. EGNPI punya `Update EGNPI Value`; retensi nol
  // punya kolom `Amount in IDR`, jadi nol konversi untuk dijalankan.
  it('⚠️ hanya SATU tombol — nol `Update Value`', () => {
    const html = render('ubah')
    expect(html).toContain(TOTAL_RETENSI.perbarui)
    expect(html).not.toContain('Update EGNPI Value')
    expect((SRC.match(/<button/g) ?? []).length).toBe(1)
  })

  // ⛔ `pyDisabledWhen` = `TreatyIn.EDMMaterialType = 2`, @202558.
  it('⛔ Update Total MATI bila EDMMaterialType = 2', () => {
    expect(totalTerkunci('2')).toBe(true)
    expect(totalTerkunci('1')).toBe(false)
    expect(render('ubah', '2')).toContain('disabled')
    const hidup = render('ubah', '1')
    const i = hidup.indexOf(TOTAL_RETENSI.perbarui)
    expect(hidup.slice(Math.max(0, i - 120), i)).not.toContain('disabled')
  })

  // ⛔ RALAT 7 Oktober 2026 — UJI INI DULU MENUNTUT HAL YANG SALAH.
  //
  // `pyReadOnly = true` DAN `pyReadOnlyCondition = TreatyIn.IsEditData = 1`
  // di sel yang sama; syaratnya yang berlaku, jadi `Note` aktif di mode
  // Edit. ⚠️ Syaratnya `IsEditData`, BUKAN `ViewState` seperti di tab
  // EGNPI — nol tambalan menyeluruh untuk keduanya.
  it('⭐ Note AKTIF di mode ubah, terkunci di mode lihat', () => {
    expect(render('ubah')).toContain('catatan retensi')
    expect(render('ubah')).not.toContain('readonly')
    expect(render('lihat')).toContain('readonly')
    // Tetap textarea sungguhan, bukan `Area` ber-onChange hampa.
    expect(SRC).toContain('<textarea')
    expect(SRC).toContain('readOnly={!bisaUbah}')
  })

  // ⚠️ Grid 0 desimal, panel total 2 — gambar 26. Nilai yang sama, dua
  // presisi, dua tempat.
  it('⚠️ grid dan panel total memakai desimal yang berbeda', () => {
    expect(DESIMAL_RETENSI.jumlah).toBe(0)
    expect(DESIMAL_RETENSI.nilaiTotal).toBe(2)
    expect(render('lihat')).toContain('3.500.000.000')
  })

  // ⭐ Total TERSIMPAN tampil sebelum tombol ditekan — panel kosong akan
  // terbaca sebagai "kontrak ini nol retensi".
  it('⭐ total awal dari dokumen tampil tanpa menekan tombol', () => {
    const html = renderToStaticMarkup(
      <TabRetensi
        baris={BARIS}
        totalAwal={[{ Currency: 'IDR', CurrencyID: '10026', Value: '3500000000' }]}
        opsiAwal={OPSI}
        mode="lihat"
      />,
    )
    expect(html).toContain('3.500.000.000,00')
    expect(html).not.toContain(TOTAL_RETENSI.tanpaBaris)
  })

  // ⭐ 8 Oktober 2026 — bentuk Pega: grid tetap tampil (kepala + `Add`),
  // baris kosongnya satu sel `No items`.
  it('⭐ nol baris → grid Pega berisi "No items" dan tombol Add', () => {
    const html = renderToStaticMarkup(
      <TabRetensi baris={[]} opsiAwal={OPSI} mode="ubah" />,
    )
    expect(html).toContain('<td colSpan="5" class="trin__kosong-pega">No items</td>')
    expect(html).toContain(`>${RETENSI.tambah}</button>`)
    expect(html).not.toContain(RETENSI.petunjukKosong)
  })

  it('⛔ nol rumus di layar — seluruhnya lewat /hitung/retensi', () => {
    expect(SRC).toContain('hitungRetensi(')
    const kode = SRC.split('\n')
      .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
      .join('\n')
    expect(kode).not.toContain('reduce(')
  })
})
