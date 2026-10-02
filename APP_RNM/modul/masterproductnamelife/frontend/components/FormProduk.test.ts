// Paritas baca-saja form produk dengan XML (audit 02-10-2026).
//
// `InboxProductName.xml`: enam medan pemilih master ber-`pyReadOnly` true, `pyEditOptions` Read-only, dan
// `pyReadOnlyCondition` KOSONG = SELALU baca-saja (Ceding b4040, SOB b4428, R/I Risk Name b7362, Cause Of
// Loss b10626, Policy Holder b17062, Currency b28105): nilainya hanya dari daftar master, tidak diketik. Sel
// `PLAN LIST` `Bussines` (`.Name` b33504) dan `Benefit` b33658 juga selalu baca-saja (diisi autocomplete `Plan Name`).
//
// Keputusan work owner 02-10-2026 ("perubahan pada tampilan untuk semua Choose ubah jadi dropdown saja"): tombol
// `Choose*` + popup FlowAction diganti dropdown master `DropdownMaster` - ketujuh master, termasuk R/I Rate baris
// `PLAN LIST`. Nilai tetap HANYA dari daftar master.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const KODE = readFileSync(join(__dirname, 'FormProduk.tsx'), 'utf8')

/** Potongan satu elemen `<DropdownMaster ... />` yang memuat `penanda`. */
function dropdown(penanda: string): string {
  const i = KODE.indexOf(penanda)
  if (i < 0) return ''
  const awal = KODE.lastIndexOf('<DropdownMaster', i)
  return KODE.slice(awal, KODE.indexOf('/>', i) + 2)
}

describe('form produk - ketujuh pemilih master berupa dropdown (keputusan work owner 02-10-2026)', () => {
  it('setiap master: label verbatim, jenis RD, mode lihat, dan penerima set*_DT yang menyalin ID + nama', () => {
    const medan = [
      ['label={UMUM_MPNL.ceding}', "jenis=\"ceding\"", 'ubahUmum({ ceding: v.nama, cedingId: v.id })'],
      ['label={UMUM_MPNL.sob}', "jenis=\"sob\"", 'ubahUmum({ sobName: v.nama, sobId: v.id })'],
      ['label={UMUM_MPNL.riRisk}', "jenis=\"ri-risk\"", 'ubahUmum({ riRisk: v.nama, riRiskId: v.id })'],
      ['label={UMUM_MPNL.causeOfLoss}', "jenis=\"penyebab\"", 'ubahUmum({ cause: v.nama, causeId: v.id })'],
      ['label={INWARD_MPNL.policyHolder}', "jenis=\"pemegang-polis\"", 'ubahInward({ policyHolderName: v.nama, policyHolder: v.id })'],
      ['label={INWARD_MPNL.currency}', "jenis=\"mata-uang\"", 'ubahInward({ currency: v.nama, currencyId: v.id })'],
      ['labelAria={PLAN_MPNL.riRate}', "jenis=\"ri-rate\"", 'ganti<BarisPlan>(x.planList, i, { riRate: v.nama, riRateId: v.id })'],
    ] as const
    for (const [label, jenis, pilih] of medan) {
      const d = dropdown(label)
      expect(d, label).not.toBe('')
      expect(d, label).toContain(jenis)
      expect(d, label).toContain('lihat={lihat}')
      expect(d, label).toContain(pilih)
    }
    expect(KODE.split('<DropdownMaster').length - 1).toBe(medan.length)
    // `SetRIRate` 1 b249: kepala kolom nama pemilih R/I Rate `RIRate Name`.
    expect(dropdown('labelAria={PLAN_MPNL.riRate}')).toContain('kolomNama={PEMILIH_MPNL.kolomRiRateName}')
  })

  it('tombol Choose* dan popup pemilih tidak dirender lagi; nilai tetap tidak dapat diketik', () => {
    expect(KODE).not.toMatch(/\{(?:UMUM|INWARD|PLAN)_MPNL\.choose\w*\}/)
    expect(KODE).not.toContain('PemilihMaster')
    expect(KODE).not.toContain('bukaPemilih')
    // Tidak ada isian ketik (autocomplete master) atau dropdown Cause Of Loss yang menerima teks bebas.
    expect(KODE).not.toContain('<Saran<NilaiMaster>')
    expect(KODE).not.toContain("cariMaster('penyebab'")
  })

  it('Policy Holder lewat dropdown tidak menjalankan SetTreatyName_Act (setPolicyHolder_DT b2416)', () => {
    expect(dropdown('label={INWARD_MPNL.policyHolder}')).not.toContain('namaTreaty')
  })

  it('mode lihat: medan master tampil baca-saja seperti sebelumnya, sel R/I Rate PLAN LIST tetap teks', () => {
    const DROPDOWN = readFileSync(join(__dirname, 'DropdownMaster.tsx'), 'utf8')
    expect(DROPDOWN).toContain(
      'if (lihat) {\n    return label ? <Field label={label} value={nilai} onChange={() => undefined} readOnly /> : <span>{nilai}</span>\n  }',
    )
    // Daftar dibaca dengan saringan `Search` dan batas + 1 (potongan dinyatakan).
    expect(DROPDOWN).toContain('cariMaster(jenis, kata, BATAS_DROPDOWN + 1)')
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
