import { describe, expect, it } from 'vitest'

import { DETAIL } from '../../assets/labels.claimlife'
import type { TotalPeserta } from '../../services/api'
import { medanTotal, RULE_TOTAL_HILANG } from './PanelTotalPeserta'

// Uji panel total layar Detail.
//
// ⛔ Berkas ini DITULIS ULANG 27-09-2026. Versinya yang lama menjaga janji
// yang ternyata salah: "layar ini tidak pernah menghitung uang yang aturannya
// tidak kita punyai" — dan ia menegakkan janji itu dengan menuntut kelima
// total TIDAK membawa angka. Aturannya ternyata kami punyai
// (`SavePesertaClaim.xml` 8.1 b4221 / 8.2 b4592), dan totalnya enam.
//
// Yang dijaga sekarang: cacah, urutan, dan bahwa angkanya datang dari
// backend apa adanya — bukan dijumlah ulang di layar.

/** Total uji dengan enam nilai BERBEDA, supaya dua medan tertukar berbunyi. */
const TOTAL: TotalPeserta = {
  cedingRetention: { amount: '11', currency: 'IDR' },
  shareNusantaraRe: { amount: '22', currency: 'IDR' },
  sumInsured: { amount: '33', currency: 'IDR' },
  sumReasured: { amount: '44', currency: 'IDR' },
  shareRetro: { amount: '55', currency: 'IDR' },
  jumlahKlaim: { amount: '66', currency: 'IDR' },
}

describe('medanTotal', () => {
  it('enam total, tidak kurang dan tidak lebih', () => {
    // ⛔ Cacahnya dikunci, dan angkanya SENGAJA enam bukan lima. Ronde
    // sebelumnya mengunci LIMA - `Total Ceding Retention` b20629 luput
    // karena pencacahannya memakai rujukan `CheckTotalAdjustmentClaim` dan
    // hanya total itu yang tidak punya aksi refresh. Uji yang mengunci
    // cacah yang salah lebih berbahaya daripada tidak ada uji: ia membuat
    // kekurangannya tampak disengaja.
    expect(medanTotal(TOTAL)).toHaveLength(6)
  })

  it('labelnya VERBATIM dan berurutan seperti di section', () => {
    // Urutannya b20629, b20914, b21201, b21488, b21775, b22063.
    expect(medanTotal(TOTAL).map((m) => m.label)).toEqual([
      DETAIL.totalCedingRetention,
      DETAIL.totalShareNusantaraRe,
      DETAIL.totalSumInsured,
      DETAIL.totalSumReasured,
      DETAIL.totalShareRetro,
      DETAIL.totalClaimAmount,
    ])
  })

  it('tiap label membawa ANGKANYA SENDIRI, tidak tertukar', () => {
    // Keenam nilainya berbeda di TOTAL, jadi dua medan yang tertukar di
    // medanTotal akan berbunyi di sini - dan hanya di sini. Dua angka uang
    // yang bertukar tempat tidak mengubah satu pun tipe.
    expect(medanTotal(TOTAL).map((m) => m.jumlah)).toEqual([
      '11',
      '22',
      '33',
      '44',
      '55',
      '66',
    ])
  })

  it('angkanya TEKS apa adanya, tidak dijumlah ulang di layar', () => {
    // Uang dijumlah sekali, di backend yang punya apd.Decimal. Bila kelak
    // seseorang menjumlah di layar, ia akan memakai number - dan uji ini
    // yang gagal lebih dulu (ADR-U-0003, ADR-U-0016).
    for (const m of medanTotal(TOTAL)) {
      expect(typeof m.jumlah, `${m.label} harus teks`).toBe('string')
      expect(Object.keys(m).sort()).toEqual(['jumlah', 'label', 'mataUang'])
    }
  })

  it('nol tetap "0", bukan kosong', () => {
    // Peserta tanpa baris adjustment bertotal NOL di Pega (penampungnya
    // mulai dari literal 0, b4027..b4159). Menampilkannya sebagai kosong
    // akan berbunyi "belum ada datanya" untuk peserta yang datanya lengkap.
    const nol: TotalPeserta = {
      cedingRetention: { amount: '0', currency: 'IDR' },
      shareNusantaraRe: { amount: '0', currency: 'IDR' },
      sumInsured: { amount: '0', currency: 'IDR' },
      sumReasured: { amount: '0', currency: 'IDR' },
      shareRetro: { amount: '0', currency: 'IDR' },
      jumlahKlaim: { amount: '0', currency: 'IDR' },
    }
    expect(medanTotal(nol).map((m) => m.jumlah)).toEqual([
      '0',
      '0',
      '0',
      '0',
      '0',
      '0',
    ])
  })

  it('nama rule yang hilang TETAP disebut apa adanya', () => {
    // ⚠️ Walau angkanya kini ada. Rujukan menggantungnya nyata - sepuluh
    // pemanggilan, nol berkas - dan OQ-H masih menanyakan kenapa rule
    // refresh-nya tidak ikut diekspor. Menghapus nama ini memutus jejak
    // antara kode dan pertanyaan yang belum dijawab.
    expect(RULE_TOTAL_HILANG).toBe('CheckTotalAdjustmentClaim')
  })
})

// ⛔ KUNCI KONTRAK DUA SISI. Tiga cacat dengan bentuk yang sama sudah terjadi
// di modul ini: envelope `galat` vs `error`, rute tanpa pemanggil, dan
// `IsCheck` "1" vs "true". Ketiganya lolos karena tiap sisi hijau sendirian.
//
// Keenam nama di bawah HARUS sama persis dengan tag JSON di
// `internal/models/totalpeserta.go`, yang dikunci
// `TestNamaJSONTotalPesertaDikunci`. Bila salah satu sisi berganti nama,
// React akan membaca `undefined` dan menampilkan medan uang KOSONG - tanpa
// satu pun galat, di layar yang tampak baik-baik saja.
describe('kontrak JSON TotalPeserta', () => {
  it('keenam nama medan dikunci', () => {
    expect(Object.keys(TOTAL).sort()).toEqual([
      'cedingRetention',
      'jumlahKlaim',
      'shareNusantaraRe',
      'shareRetro',
      'sumInsured',
      'sumReasured',
    ])
  })
})
