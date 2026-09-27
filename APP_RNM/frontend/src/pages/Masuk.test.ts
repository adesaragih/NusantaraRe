// Penjaga layar masuk stub — F0.2.
//
// ⛔ Penjaga SIFAT, bukan render. Uji render menuntut DOM; yang harus dijaga
// di sini adalah tiga janji yang bila dilanggar tidak berbunyi sama sekali:
// nol sandi, gerbang env yang gagal tertutup, dan pita peringatannya ada.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const SUMBER = readFileSync(join(__dirname, 'Masuk.tsx'), 'utf8')

/** Buang komentar, supaya prosa yang MENYEBUT sebuah kata tidak tertuduh. */
function tanpaKomentar(teks: string): string {
  return teks
    .replace(/\/\*[\s\S]*?\*\//g, '')
    .split('\n')
    .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
    .join('\n')
}

const kode = tanpaKomentar(SUMBER)

describe('layar masuk stub', () => {
  // ⚠️ Yang dilarang MEKANISMEnya, bukan KATAnya. Ronde pertama penjaga
  // ini melarang kata "sandi" di mana pun - dan ia menuduh pita
  // peringatan yang justru berbunyi "tidak ada sandi". Penjaga yang
  // menuduh kalimat yang benar akan dilonggarkan orang, bukan dipatuhi.
  it.each([
    ['isian sandi', /type=["']password["']/],
    ['komponen FieldSandi', /FieldSandi/],
    ['keadaan bernama sandi', /\b(set)?[sS]andi\s*[,)=:]/],
    ['keadaan bernama password', /\b(set)?[pP]assword\s*[,)=:]/],
    ['autoComplete sandi', /autoComplete=["'][^"']*password/],
  ])('nol %s', (_nama, pola) => {
    expect(kode).not.toMatch(pola)
  })

  it('pita peringatannya justru MENYEBUT bahwa tidak ada sandi', () => {
    // Kata itu harus ADA di layar - yang dilarang mekanismenya.
    //
    // ⚠️ Spasi dinormalkan lebih dulu: JSX memenggal kalimat pada batas
    // baris, sehingga 'tidak ada sandi' tersimpan terpisah newline.
    // Mencarinya mentah-mentah gagal karena PEMFORMATAN, bukan karena
    // kalimatnya hilang.
    const satuBaris = kode.replace(/\s+/g, ' ')
    expect(satuBaris).toContain('tidak ada sandi')
  })

  it('bergerbang VITE_AUTH_STUB dan GAGAL TERTUTUP', () => {
    // Gerbangnya di store/sesi.ts; layar hanya memanggilnya.
    expect(kode).toContain('bolehMasukStub()')
    // Cabang "tidak boleh" mengembalikan layar TANPA tombol masuk.
    const tolak = kode.slice(kode.indexOf('if (!stubMenyala)'))
    expect(tolak).toContain('VITE_AUTH_STUB')
    expect(tolak.slice(0, tolak.indexOf('return (\n    <main className="masuk">\n      <form')))
      .not.toContain('type="submit"')
  })

  it('memasang pita peringatan "mode stub"', () => {
    expect(kode).toContain('Mode stub.')
    expect(kode).toContain('role="status"')
  })

  it('peran dipilih dengan KOTAK CENTANG, bukan pilihan tunggal', () => {
    // `pelakuDari` memecah `X-Peran` pada koma; memaksa satu peran akan
    // menutup jalur pelaku yang memegang lebih dari satu.
    expect(kode).toContain('type="checkbox"')
    expect(kode).not.toContain('type="radio"')
  })

  it('menolak masuk tanpa peran, di layar — bukan menunggu 403', () => {
    expect(kode).toContain('peran.length === 0')
  })

  it('menolak masuk bila penyimpanan menolak menyimpan sesi', () => {
    expect(kode).toContain('!sesi.simpan(s)')
  })
})
