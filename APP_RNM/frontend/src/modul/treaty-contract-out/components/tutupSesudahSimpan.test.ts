// Form Treaty Contract Out TERTUTUP sesudah simpan berhasil — keputusan work
// owner 30-09-2026: "kalo berhasil di save, form edit nya langsung ke close".

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

/** Potongan sesudah panggilan simpan hingga akhir blok `try` / fungsi. */
function sesudahSimpan(berkas: string, panggilan: string): string {
  const kode = readFileSync(join(__dirname, berkas), 'utf8')
  const i = kode.indexOf(panggilan)
  expect(i, `${berkas}: ${panggilan}`).toBeGreaterThan(0)
  return kode.slice(i, kode.indexOf('\n', kode.indexOf('setForm(', i)))
}

describe('form tertutup sesudah simpan berhasil', () => {
  it.each([
    ['PanelKontrakTahun.tsx', 'await simpanKontrakTahun('],
    ['PanelReinsurerKombinasi.tsx', 'await simpanReinsurerKombinasi('],
    ['PanelSecurityReinsurer.tsx', 'await simpanSecurity('],
    ['PanelBusinessKombinasi.tsx', 'await simpanBusinessKombinasi('],
    ['PanelJenisKlausul.tsx', 'await simpanKlausul('],
    ['../pages/InboxTreatyContract.tsx', 'await simpanTahunTreaty('],
  ])('%s', (berkas, panggilan) => {
    const potongan = sesudahSimpan(berkas, panggilan)
    // Setter form PERTAMA sesudah simpan adalah `setForm(null)`, bukan membuka ulang isian.
    expect(potongan).toMatch(/setForm\(null\)$/)
    expect(potongan).not.toMatch(/buka\(form|setForm\(form/)
  })
})
