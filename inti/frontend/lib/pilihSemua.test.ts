import { describe, expect, it } from 'vitest'

import { jungkitPilihSemua, terpilihSemua, type KeadaanPilihSemua } from './pilihSemua'

// `SelectAllClaimLife_act.xml:247` — penjungkit TIGA keadaan.

describe('jungkitPilihSemua', () => {
  it('kosong → true: tekanan pertama MENCENTANG', () => {
    // ⛔ Inilah sebab keadaannya teks dan bukan boolean. Bila kosong
    // diperlakukan sama dengan 'false', tekanan pertama akan melepas centang
    // pada daftar yang memang belum tercentang - tombol yang tampak rusak.
    expect(jungkitPilihSemua('')).toBe('true')
  })

  it('true → false', () => {
    expect(jungkitPilihSemua('true')).toBe('false')
  })

  it('false → true', () => {
    expect(jungkitPilihSemua('false')).toBe('true')
  })

  it('nilai asing diperlakukan sebagai belum tercentang', () => {
    // Cabang terakhir @if bersarang Pega menampung apa pun selain '' dan
    // 'true'. Ditiru apa adanya: nilai yang tidak kita kenal berarti daftar
    // belum tercentang, sehingga tombolnya mencentang.
    for (const asing of ['TRUE', '1', 'ya', 'null']) {
      expect(jungkitPilihSemua(asing)).toBe('true')
    }
  })

  it('dua tekanan kembali ke true, bukan ke kosong', () => {
    // Keadaan '' tidak pernah kembali: ia hanya keadaan awal. Penjungkit yang
    // dapat kembali ke '' akan membuat siklusnya tiga langkah, dan pemakai
    // melihat satu tekanan yang seolah tidak melakukan apa-apa.
    const satu = jungkitPilihSemua('')
    const dua = jungkitPilihSemua(satu)
    const tiga = jungkitPilihSemua(dua)
    expect([satu, dua, tiga]).toEqual(['true', 'false', 'true'])
  })
})

describe('terpilihSemua', () => {
  it.each<[KeadaanPilihSemua, boolean]>([
    ['', false],
    ['true', true],
    ['false', false],
  ])('keadaan %s → %s', (keadaan, mau) => {
    expect(terpilihSemua(keadaan)).toBe(mau)
  })

  it('hanya teks "true" PERSIS yang berarti tercentang', () => {
    // b429 menyalin Select.CARI1 ke .IsAccept apa adanya, jadi pembacaannya
    // harus sama persis di kedua sisi. 'TRUE' bukan 'true'.
    expect(terpilihSemua('TRUE')).toBe(false)
    expect(terpilihSemua('1')).toBe(false)
  })
})
