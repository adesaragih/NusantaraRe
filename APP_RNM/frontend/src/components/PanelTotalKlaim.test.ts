import { describe, expect, it } from 'vitest'

import { DETAIL } from '../assets/labels'
import { medanTotal, RULE_TOTAL_HILANG } from './PanelTotalKlaim'

// Uji panel total layar Detail.
//
// Yang dijaga bukan tampilannya melainkan SATU janji: layar ini tidak pernah
// menghitung uang yang aturannya tidak kita punyai.

describe('medanTotal', () => {
  it('lima total, tidak kurang dan tidak lebih', () => {
    // ⛔ Cacahnya dikunci. Total yang DIHILANGKAN dari layar sama merusaknya
    // dengan total yang dikarang, dan yang pertama tidak berbunyi: orang
    // hanya melihat layar yang tampak lengkap.
    expect(medanTotal()).toHaveLength(5)
  })

  it('labelnya VERBATIM dan berurutan seperti di section', () => {
    // Urutannya mengikuti b20914, b21201, b21488, b21775, b22063.
    expect(medanTotal().map((m) => m.label)).toEqual([
      DETAIL.totalShareNusantaraRe,
      DETAIL.totalSumInsured,
      DETAIL.totalSumReasured,
      DETAIL.totalShareRetro,
      DETAIL.totalClaimAmount,
    ])
  })

  it('TIDAK SATU PUN total membawa angka', () => {
    // ⛔ Inti berkas ini. `CheckTotalAdjustmentClaim` dirujuk sepuluh kali di
    // section itu tetapi NOL berkas rule-nya ada di korpus, jadi kita tidak
    // tahu baris mana yang ikut dihitung — seluruhnya, hanya yang `IsCheck`,
    // atau tanpa yang `STS_REJECT = 2`. Jalur tolak mencabut `IsCheck`, jadi
    // jawabannya berpengaruh, dan ini angka uang.
    //
    // Bila kelak seseorang menambahkan penjumlahan di sini, uji inilah yang
    // gagal lebih dulu dan menagih rule-nya — bukan neraca yang salah.
    for (const m of medanTotal()) {
      expect(m.belumBersumber, `${m.label} harus dinyatakan belum bersumber`).toBe(
        true,
      )
      expect(Object.keys(m).sort()).toEqual(['belumBersumber', 'label'])
    }
  })

  it('nama rule yang hilang disebut apa adanya', () => {
    // Supaya orang dapat mencari teks yang sama di ekspor Pega dan di OQ-H.
    expect(RULE_TOTAL_HILANG).toBe('CheckTotalAdjustmentClaim')
  })
})
