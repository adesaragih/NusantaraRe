// Paritas baca-saja form produk dengan XML (audit 02-10-2026).
//
// `InboxProductName.xml`: enam medan pemilih master ber-`pyReadOnly` true, `pyEditOptions` Read-only, dan
// `pyReadOnlyCondition` KOSONG = SELALU baca-saja (Ceding b4040, SOB b4428, R/I Risk Name b7362, Cause Of
// Loss b10626, Policy Holder b17062, Currency b28105); nilainya hanya dari tombol `Choose*`. Sel `PLAN LIST`
// `Bussines` (`.Name` b33504) dan `Benefit` b33658 juga selalu baca-saja (diisi autocomplete `Plan Name`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const KODE = readFileSync(join(__dirname, 'FormProduk.tsx'), 'utf8')

describe('form produk - medan yang di XML selalu baca-saja', () => {
  it('enam medan pemilih master tampil baca-saja, diisi hanya lewat tombol Choose*', () => {
    const medan = [
      ['UMUM_MPNL.ceding', 'UMUM_MPNL.chooseCeding'],
      ['UMUM_MPNL.sob', 'UMUM_MPNL.chooseSob'],
      ['UMUM_MPNL.riRisk', 'UMUM_MPNL.chooseRiRisk'],
      ['UMUM_MPNL.causeOfLoss', 'UMUM_MPNL.chooseCause'],
      ['INWARD_MPNL.policyHolder', 'INWARD_MPNL.choosePolicyHolder'],
      ['INWARD_MPNL.currency', 'INWARD_MPNL.chooseCurrency'],
    ] as const
    for (const [label, tombol] of medan) {
      const i = KODE.indexOf(`{medanMaster(\n          ${label},`)
      expect(i, label).toBeGreaterThan(-1)
      const potong = KODE.slice(i, KODE.indexOf('</button>', i))
      expect(potong, label).toContain(`bukaPemilih(${tombol},`)
    }
    // Pembungkusnya merender medan baca-saja TANPA syarat, tombolnya hanya di luar mode lihat.
    expect(KODE).toMatch(/function medanMaster\([^)]*\) \{\s*return \(\s*<div className="mpnl-medan-pilih">\s*<Field label=\{label\} value=\{nilai\} onChange=\{\(\) => undefined\} readOnly \/>\s*\{!lihat && tombol\}/)
    // Tidak ada lagi isian ketik (autocomplete master) atau dropdown Cause Of Loss yang dapat diubah.
    expect(KODE).not.toContain('<Saran<NilaiMaster>')
    expect(KODE).not.toContain("cariMaster('penyebab'")
  })

  it('Policy Holder lewat Choose tidak menjalankan SetTreatyName_Act (setPolicyHolder_DT b2416)', () => {
    const i = KODE.indexOf('bukaPemilih(INWARD_MPNL.choosePolicyHolder,')
    const potong = KODE.slice(i, KODE.indexOf('</button>', i))
    expect(potong).toContain('ubahInward({ policyHolderName: v.nama, policyHolder: v.id })')
    expect(potong).not.toContain('namaTreaty')
  })

  it('sel Bussines dan Benefit PLAN LIST tampil sebagai teks, tidak dapat diketik', () => {
    expect(KODE).toContain('<td>{b.name}</td>')
    expect(KODE).toContain('<td>{b.benefit}</td>')
    expect(KODE).not.toMatch(/label=\{PLAN_MPNL\.(bussines|benefit)\}/)
  })
})

describe('form produk - tanda wajib (keputusan work owner 02-10-2026)', () => {
  it('empat medan ber-pyRequired true di XML bertanda wajib *', () => {
    // Product Name b3585, Premium Factor (%) b25362, Annuity Interest (%) b26054, Premium Refund Factor (%) b26268.
    for (const label of ['UMUM_MPNL.productName', 'INWARD_MPNL.premiumFactor', 'INWARD_MPNL.annuityInterest', 'INWARD_MPNL.premiumRefundFactor']) {
      const i = KODE.indexOf(`label={${label}}`)
      expect(i, label).toBeGreaterThan(-1)
      const medan = KODE.slice(i, KODE.indexOf('/>', i))
      expect(medan, label).toMatch(/\brequired\b/)
    }
  })
})
