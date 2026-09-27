// Penjaga ADR-U-0003 di lapis terakhir — F0.1, brief lanjutan 6 §2 aturan 3.
//
// ⛔ Uang di proyek ini adalah TEKS DESIMAL EKSAK, dari Oracle sampai layar.
// `lib/desimal.ts` diadopsi dari REFERENSI_UI justru karena ia bekerja pada
// DIGIT lewat BigInt; nilainya hilang seketika bila seseorang "menyederhanakan"
// satu fungsi di dalamnya menjadi `Number(a) + Number(b)`.
//
// Kerusakannya TIDAK terlihat pada nilai kecil: `Number("1.10") + Number("2.20")`
// menghasilkan 3,3000000000000003 - dan dibulatkan tampilannya kembali menjadi
// "3,30". Ia baru terlihat pada nilai besar berdesimal panjang, yaitu justru
// nilai klaim reasuransi.
//
// Karena itu penjaganya TEKSTUAL: ia melarang jalannya, bukan menunggu satu
// nilai kebetulan membuktikan kerusakannya.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { jumlahDesimal, samaDengan, geserTitik, pangkasNolEkor } from './desimal'

const SUMBER = readFileSync(join(__dirname, 'desimal.ts'), 'utf8')

/** Buang komentar, supaya prosa yang MENYEBUT `Number(` tidak ikut tertuduh. */
function tanpaKomentar(teks: string): string {
  return teks
    .replace(/\/\*[\s\S]*?\*\//g, '')
    .split('\n')
    .filter((b) => !b.trimStart().startsWith('//'))
    .join('\n')
}

describe('desimal.ts tidak pernah melewatkan uang ke float', () => {
  const kode = tanpaKomentar(SUMBER)

  it.each(['Number(', 'parseFloat', 'parseInt', 'Math.round', 'toFixed('])(
    'nol pemakaian %s',
    (terlarang) => {
      expect(kode).not.toContain(terlarang)
    },
  )

  it('memakai BigInt — itulah yang membuat penjumlahannya eksak', () => {
    expect(kode).toContain('BigInt(')
  })
})

describe('perilaku eksak, bukan sekadar bentuk', () => {
  it('menjumlahkan desimal panjang tanpa kehilangan digit', () => {
    // ⛔ Kasus yang float64 rusakkan: 0,1 + 0,2 menjadi 0,30000000000000004.
    expect(jumlahDesimal(['0.1', '0.2']).total).toBe('0.3')
  })

  it('menjumlahkan nilai besar berdesimal panjang', () => {
    // Nilai sebesar ini dengan 9 desimal melampaui presisi float64.
    const hasil = jumlahDesimal(['9007199254740993.000000001', '0.000000001'])
    expect(hasil.total).toBe('9007199254740993.000000002')
  })

  it('isian rusak TIDAK dianggap nol — ia dilaporkan', () => {
    const hasil = jumlahDesimal(['1.00', 'bukan angka', '2.00'])
    expect(hasil.total).toBe('3.00')
    expect(hasil.takTerbaca).toEqual(['bukan angka'])
  })

  it('isian KOSONG dilewati — kosong berarti belum diisi, bukan rusak', () => {
    const hasil = jumlahDesimal(['1.00', '', '  ', '2.00'])
    expect(hasil.total).toBe('3.00')
    expect(hasil.takTerbaca).toEqual([])
  })

  it('membandingkan per digit: 100 = 100.000000000, tetapi 99.999999999 tidak', () => {
    expect(samaDengan('100', '100.000000000')).toBe(true)
    expect(samaDengan('100', '99.999999999')).toBe(false)
  })

  it('menggeser titik tanpa merusak nilai ekstrem', () => {
    // Ketiga mode gagal jalur float yang dicatat referensi:
    expect(geserTitik('0.00000000005', 2)).toBe('0.000000005')
    expect(geserTitik('1', -7)).toBe('0.0000001')
    // Teks yang bukan desimal dikembalikan APA ADANYA - layar tidak menebak.
    expect(geserTitik('0x10', 2)).toBe('0x10')
  })

  it('memangkas nol ekor tanpa membulatkan', () => {
    expect(pangkasNolEkor('100.000000000', 2)).toBe('100.00')
    // ⛔ BUKAN pembulatan: 99,999999999 tidak boleh menjadi 100,00 dan
    // menyembunyikan bahwa totalnya kurang.
    expect(pangkasNolEkor('99.999999999', 2)).toBe('99.999999999')
  })
})
