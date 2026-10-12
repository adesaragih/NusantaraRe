import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

// ⛔ Laporan pemakai 10 Oktober 2026: tab Exclusions dan Special Conditions
// menyimpan teks yang SAMA. Kedua tab merender `TabTeksPanjang` di posisi
// pohon yang sama; tanpa `key` React memakai ulang komponennya dan teks tab
// pertama tersemai ke properti tab kedua.
describe('Exclusions dan Special Conditions — satu komponen, dua keadaan', () => {
  const form = readFileSync(join(__dirname, 'pages/FormKontrakTreatyIn.tsx'), 'utf8')

  it('setiap `TabTeksPanjang` ber-`key` menurut tabnya', () => {
    const pakai = form.match(/<TabTeksPanjang/g) ?? []
    const ber = form.match(/<TabTeksPanjang[^>]*?key=\{tabTampil\}/gs) ?? []
    expect(pakai.length).toBe(2)
    expect(ber.length).toBe(2)
  })
})
