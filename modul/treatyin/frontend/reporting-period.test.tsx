// Uji tab Reporting Period — tombol Apply (`Activity/TreatyInSetReport.xml`).
//
// Rumusnya diuji di services (`periode_pelaporan_test.go`); di sini yang
// dijaga: layar MEMANGGIL rumus itu, tidak menulis rumus sendiri, dan kotak
// tanggal grid diisi nilai TERSIMPAN.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisPeriodeWarisan } from './api'
import TabReportingPeriod from './components/TabReportingPeriod'

const SUMBER = readFileSync(join(__dirname, 'components', 'TabReportingPeriod.tsx'), 'utf8')

const BARIS: BarisPeriodeWarisan[] = [
  {
    periode: 'Q 1',
    hitungOtomatis: '',
    tanggalAwal: '06/10/26',
    jatuhTempoKirim: '17/01/27',
    jatuhTempoKonfirmasi: '29/01/27',
    jatuhTempoBayar: '10/02/27',
    tanggalAwalAsli: '20261006',
    jatuhTempoKirimAsli: '20270117',
    jatuhTempoKonfirmasiAsli: '20270129',
    jatuhTempoBayarAsli: '20270210',
  },
]
const OPSI = [{ value: 'quarter', label: 'Quarter Year' }, { value: 'other', label: 'Others' }]

describe('tombol Apply', () => {
  it('⭐ Apply MEMANGGIL rumus services — nol rumus di layar', () => {
    // `awal` = `param.startdate` — kosong pada Apply, `.InitialDate` pada sel (7 Oktober 2026).
    expect(SUMBER).toContain('hitungPeriodePelaporan({ mulai, akhir, periode, interval, penyerahan, konfirmasi, pelunasan, awal })')
    expect(SUMBER).not.toMatch(/addMonths|setMonth|getDate\(\)|new Date\(/)
    // Bentuk lama: Apply mengembalikan baris dokumen tanpa menghitung.
    expect(SUMBER).not.toContain('setBaris(barisAwal)')
  })

  it('pilihan Period dari services; "Quarter Year" label nilai `quarter`', () => {
    const html = renderToStaticMarkup(<TabReportingPeriod baris={[]} opsiPeriode={OPSI} mode="ubah" />)
    expect(html).toContain('<option value="quarter">Quarter Year</option>')
  })

  it('⭐ label tanpa "(Days)" dan pesan statis di samping Apply — seperti layar Pega', () => {
    const html = renderToStaticMarkup(<TabReportingPeriod baris={[]} opsiPeriode={OPSI} mode="ubah" />)
    expect(html).not.toContain('(Days)')
    expect(html).toContain('>Apply</button>')
    expect(html).toContain('Start Date, Due. Must Not Be Empty')
  })

  it('mode lihat: Apply mati', () => {
    const html = renderToStaticMarkup(<TabReportingPeriod baris={[]} opsiPeriode={OPSI} mode="lihat" />)
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Apply<\/button>/)
  })
})

describe('grid hasil', () => {
  it('⛔ mode ubah: kotak tanggal diisi nilai TERSIMPAN (YYYY-MM-DD), bukan dd/mm/yy', () => {
    const html = renderToStaticMarkup(<TabReportingPeriod baris={BARIS} opsiPeriode={OPSI} mode="ubah" />)
    for (const v of ['2026-10-06', '2027-01-17', '2027-01-29', '2027-02-10']) expect(html).toContain(`value="${v}"`)
    expect(html).toContain('Auto Calculate')
  })

  it('mode lihat: teks tampil; Auto Calculate tersembunyi (`ViewState != 1`)', () => {
    const html = renderToStaticMarkup(<TabReportingPeriod baris={BARIS} opsiPeriode={OPSI} mode="lihat" />)
    expect(html).toContain('17/01/27')
    expect(html).not.toContain('Auto Calculate')
    expect(html).not.toContain('type="date" value="2027')
  })
})
