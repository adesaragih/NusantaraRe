// Penjaga tampilan form produk yang TIDAK bergantung korpus XML (keputusan work owner 05-10-2026).
// `FormProduk.test.ts` membaca korpus saat dikumpulkan dan tidak berjalan di mesin tanpa korpus;
// tes di sini tetap berjalan.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const KODE = readFileSync(join(__dirname, 'FormProduk.tsx'), 'utf8')

describe('form produk - tampilan 05-10-2026', () => {
  it('checkbox On Retention tidak dirender; isOrs tidak ditulis layar dan hitungOutward tidak dikirim', () => {
    expect(KODE).not.toContain('type="checkbox"')
    expect(KODE).not.toContain('UMUM_MPNL.onRetention')
    expect(KODE).not.toMatch(/hitungOutward:\s*true/)
    expect(KODE).not.toMatch(/isOrs:/)
  })

  it('PLAN LIST: pilihan Plan Name disaring Bussines yang sudah dipakai baris lain', () => {
    expect(KODE).toContain('tampil: saringPlanBisnis(h.tampil, p.planList, i)')
  })

  it('riwayat komentar tampil terbaru di atas tanpa mengurut ulang state', () => {
    expect(KODE).toContain('urutKomentarTampil(p.commentList).map(({ k, i }) => (')
    expect(KODE).not.toMatch(/commentList:\s*[^,}]*(sort|reverse)/)
    expect(KODE).not.toMatch(/p\.commentList\.(sort|reverse)\(/)
  })
})
