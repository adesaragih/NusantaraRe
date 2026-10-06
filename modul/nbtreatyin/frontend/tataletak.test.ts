// Tata letak dua kolom layar kasus (perintah work owner 05-10-2026, screenshot layar Pega: "coba ikuti layoutnya begini,
// buat rapih ya"; lanjutan: "RNM Share 1 baris", "choose bisnis paling atas", "banyak sekali tulisan IDR"): kolom
// kiri / kanan section General, Remark melebar di bawah; bagian uang Gross di atas, kolom OGP dan kolom ONP. Setiap
// medan definisi tetap muncul TEPAT sekali, dengan urutan aslinya di dalam kolomnya.

import { describe, expect, it } from 'vitest'

import { POLIS } from './api'
import {
  MEDAN_ADMIN_UANG,
  MEDAN_ADMIN_UMUM,
  MEDAN_ATASAN_TOTAL,
  MEDAN_ATASAN_UANG,
  MEDAN_ATASAN_UMUM,
  type Medan,
} from './medan'
import { deretQ, TATA_UANG_ADMIN, TATA_UANG_ATASAN, tataUmum, TOTAL_ATASAN, type TataUang } from './tataletak'

const label = (ms: Medan[]) => ms.map((m) => m.label)
const jalur = (ms: Medan[]) => ms.map((m) => m.jalur)
const semuaUang = (t: TataUang) => [
  ...t.atasKiri,
  ...t.atasKanan,
  ...t.kiri.flatMap((k) => k.medan),
  ...t.kanan.flatMap((k) => k.medan),
]
const tanpaKodeShare = (ms: Medan[]) => jalur(ms).filter((j) => j !== POLIS + 'ShareCurrency')

describe('tata letak section General', () => {
  it('admin: kolom kiri Master ID .. RNM Share, kanan Statement Date .. Payment Type, Remark di bawah', () => {
    const t = tataUmum(MEDAN_ADMIN_UMUM)
    expect(label(t.kiri).slice(0, 4)).toEqual(['Master ID', 'Commencement', 'Statement Period', 'Source Of Business'])
    expect(t.kiri.at(-1)?.jalur).toBe(POLIS + 'ShareValue')
    expect(label(t.kanan).slice(0, 4)).toEqual(['Statement Date', 'Termination', 'To', 'Ceding Company'])
    expect(label(t.kanan)).toContain('Payment Type')
    expect(label(t.bawah)).toEqual(['Remark'])
    expect(jalur([...t.kiri, ...t.kanan, ...t.bawah])).toEqual(tanpaKodeShare(MEDAN_ADMIN_UMUM))
  })

  it('atasan: pembagian yang sama, seluruh medan tetap ada berurutan', () => {
    const t = tataUmum(MEDAN_ATASAN_UMUM)
    expect(t.kanan[0]?.label).toBe('Statement Date')
    expect(label(t.bawah)).toEqual(['Remark'])
    expect(jalur([...t.kiri, ...t.kanan, ...t.bawah])).toEqual(tanpaKodeShare(MEDAN_ATASAN_UMUM))
  })

  it('RNM Share satu baris: "RNM Share  IDR 163.125.000" - baris ShareCurrency dilebur ke ShareValue', () => {
    for (const ms of [MEDAN_ADMIN_UMUM, MEDAN_ATASAN_UMUM]) {
      const share = tataUmum(ms).kiri.filter((m) => m.jalur.startsWith(POLIS + 'Share'))
      expect(share).toHaveLength(1)
      expect(share[0]).toMatchObject({
        jalur: POLIS + 'ShareValue',
        label: 'RNM Share',
        mataUang: POLIS + 'ShareCurrency',
      })
    }
    // ShareCurrency tidak tampil: ShareValue tetap apa adanya.
    const tanpa = MEDAN_ADMIN_UMUM.filter((m) => m.jalur !== POLIS + 'ShareCurrency')
    expect(tataUmum(tanpa).kiri.find((m) => m.jalur === POLIS + 'ShareValue')?.label).toBe('ShareValue')
  })

  it('deret Q / U/Y menjadi satu baris; medan lain tetap sendiri-sendiri', () => {
    const kanan = tataUmum(MEDAN_ADMIN_UMUM).kanan
    const d = deretQ(kanan)
    const deret = d.find((x) => Array.isArray(x)) as Medan[]
    expect(label(deret)).toEqual(['Q', '/', 'U/Y'])
    expect(d.flat()).toEqual(kanan)
    expect(deretQ(kanan.filter((m) => m.label !== '/')).some((x) => Array.isArray(x))).toBe(false)
  })
})

describe('tata letak bagian uang', () => {
  it('admin: Gross Premium 100% | Claim 100%; kolom OGP sampai Balance Due To Us; kolom ONP sampai PPN 2.2%', () => {
    const t = TATA_UANG_ADMIN
    expect(label(t.atasKiri)).toEqual(['Gross Premium 100%'])
    expect(label(t.atasKanan)).toEqual(['Claim 100%'])
    expect(t.kiri[0]?.judul).toBe('OGP')
    expect(label(t.kiri[0]!.medan)[0]).toBe('Premi Ogp')
    expect(label(t.kiri.flatMap((k) => k.medan)).at(-1)).toBe('Balance Due To Us')
    expect(t.kanan[0]?.judul).toBe('ONP')
    expect(label(t.kanan.flatMap((k) => k.medan))).toEqual([
      'Premi Onp',
      '(%) Deduction In A (ONP)',
      'Deduction In A (ONP)',
      '(%) Deduction In B (ONP)',
      'Deduction In B (ONP)',
      'Deduction1',
      'Deduction2',
      'PPH 2%',
      'PPN 2.2%',
    ])
    expect(jalur(semuaUang(t)).sort()).toEqual(jalur(MEDAN_ADMIN_UANG).sort())
  })

  it('atasan: Gross di atas, kolom OGP dan ONP; seluruh medan tetap ada', () => {
    const t = TATA_UANG_ATASAN
    expect(label(t.atasKiri)).toEqual(['Gross Premium 100%'])
    expect(t.kiri[0]?.judul).toBe('OGP')
    expect(t.kanan[0]?.judul).toBe('ONP')
    expect(jalur(semuaUang(t)).sort()).toEqual(jalur(MEDAN_ATASAN_UANG).sort())
  })

  it('kode mata uang polis tidak diulang di setiap medan uang (sudah tampil di medan Currency)', () => {
    expect(MEDAN_ADMIN_UANG.some((m) => m.mataUang === POLIS + 'Currency')).toBe(true)
    for (const t of [TATA_UANG_ADMIN, TATA_UANG_ATASAN]) {
      expect(semuaUang(t).filter((m) => m.mataUang !== undefined)).toEqual([])
    }
    expect(TOTAL_ATASAN.flat().filter((m) => m.mataUang !== undefined)).toEqual([])
  })

  it('total atasan dua kolom: premium | klaim', () => {
    expect(TOTAL_ATASAN.map(label)).toEqual([
      ['Total %Share', 'Total Premium'],
      ['Total %Share Claim', 'Total Claim'],
    ])
    expect(jalur(TOTAL_ATASAN.flat())).toEqual(jalur(MEDAN_ATASAN_TOTAL))
  })
})
