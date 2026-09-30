import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { DETAIL_POLIS, JUDUL_KOLOM_PESERTA } from '../labels'
import { judulKolom, kalimatNomor, selPeserta } from './PremiumListDetail'

// Uji layar Premium List Detail — tiket 03.

const BERKAS = readFileSync(join(__dirname, 'PremiumListDetail.tsx'), 'utf8')

/**
 * Sumber TANPA komentar.
 *
 * ⛔ Prosa yang MENERANGKAN sebuah larangan tidak boleh dituduh melanggarnya.
 * Penjaga yang menuduh hal yang benar akan dilonggarkan orang, bukan dipatuhi.
 */
const SUMBER = BERKAS.split('\n')
  .filter((b) => {
    const t = b.trimStart()
    return !t.startsWith('//') && !t.startsWith('*')
  })
  .join('\n')

describe('kolom grid peserta', () => {
  it('daftar kolom datang dari server, tidak diketik ulang di layar', () => {
    // ⛔ Layar merender `hal.kolom`, bukan sebuah array literalnya sendiri.
    expect(SUMBER).toContain('hal.kolom.map')
    expect(SUMBER).not.toContain("'CERTIFICATE_NO'")
    expect(SUMBER).not.toContain("'SUM_INSURED'")
  })

  it('kolom tak dikenal tetap tampil dengan namanya', () => {
    // ⚠️ Kolom yang hilang dari layar karena judulnya belum terdaftar adalah
    // data yang hilang tanpa satu pun tanda.
    expect(judulKolom('KOLOM_BARU_YANG_BELUM_ADA')).toBe('KOLOM_BARU_YANG_BELUM_ADA')
  })

  it('judul yang dikenal diperindah dari labels.premiumlist', () => {
    expect(judulKolom('SUM_INSURED')).toBe(JUDUL_KOLOM_PESERTA.SUM_INSURED)
    expect(judulKolom('CERTIFICATE_NO')).toBe('Certificate No')
  })

  it('judul kolom TIDAK mengarang terjemahan Indonesia', () => {
    // ⛔ `PL_Detail_Sec.xml` nol `pyCaption`; judulnya nama kolom yang
    // dirapikan, bukan terjemahan yang kami karang sendiri.
    for (const judul of Object.values(JUDUL_KOLOM_PESERTA)) {
      expect(judul).not.toMatch(/Jumlah|Tanggal|Premi\b/)
    }
  })

  it('dua medan PL_Detail_Sec yang tanpa kolom tidak punya judul', () => {
    // ⛔ REINSTYPENAME dan RetrocadedShare tidak punya kolom di migrasi mana
    // pun; memberinya judul berarti menjanjikan kolom yang tidak pernah ada.
    expect(JUDUL_KOLOM_PESERTA.REINSTYPENAME).toBeUndefined()
    expect(JUDUL_KOLOM_PESERTA.RetrocadedShare).toBeUndefined()
  })

  it('NAME_OF_INSURED tidak punya judul — nol nama orang di grid', () => {
    expect(JUDUL_KOLOM_PESERTA.NAME_OF_INSURED).toBeUndefined()
    expect(SUMBER).not.toContain('NAME_OF_INSURED')
  })
})

describe('uang tetap teks', () => {
  it('nol Number(), parseFloat, atau toFixed di layar', () => {
    // ⛔ ADR-U-0003. Premi delapan angka desimal dibulatkan diam-diam oleh
    // ketiganya, dan pembulatan di jalan pulang tidak kalah salah dari
    // pembulatan saat menyimpan.
    expect(SUMBER).not.toMatch(/\bNumber\(/)
    expect(SUMBER).not.toContain('parseFloat')
    expect(SUMBER).not.toContain('toFixed')
  })

  it('sel kosong ditandai, bukan dibiarkan kosong', () => {
    expect(selPeserta('')).toBe('—')
    expect(selPeserta('   ')).toBe('—')
    expect(selPeserta(undefined)).toBe('—')
    // ⛔ Nol BUKAN kosong: premi nol dan premi belum diisi adalah dua
    // keadaan berbeda, dan hanya satu perlu dikerjakan.
    //
    // ⚠️ Rujukan ADR-U-0027 dicabut dari sini 28-09-2026 - ADR itu tentang
    // kolom nullable dan wajib-isi di kode, bukan tentang pemetaan ini.
    expect(selPeserta('0')).toBe('0')
    expect(selPeserta('1234.56789012')).toBe('1234.56789012')
  })
})

describe('nomor PL', () => {
  it('polis tanpa nomor dikatakan belum bernomor, bukan dikosongkan', () => {
    expect(kalimatNomor({ polisId: 'X', type: 'QR', businessCode: 'LF', plNumber: '', medanTanpaKolom: [] })).toBe(
      DETAIL_POLIS.belumBernomor,
    )
    expect(kalimatNomor({ polisId: 'X', type: 'QR', businessCode: 'LF', plNumber: '   ', medanTanpaKolom: [] })).toBe(
      DETAIL_POLIS.belumBernomor,
    )
  })

  it('polis bernomor menampilkan nomornya apa adanya', () => {
    expect(
      kalimatNomor({
        polisId: 'X',
        type: 'QR',
        businessCode: 'LF',
        plNumber: 'RNML-QRLF.09.26.00007',
        medanTanpaKolom: [],
      }),
    ).toBe('RNML-QRLF.09.26.00007')
  })

  it('kepala kosong tidak melempar', () => {
    expect(kalimatNomor(null)).toBe('')
  })

  it('tombol mati bila sudah bernomor atau belum punya peserta', () => {
    // ⛔ Tombol yang tetap hidup tetapi selalu menjawab hal yang sama
    // mengajari orang mengabaikan jawabannya.
    expect(SUMBER).toContain('disabled={sibuk || terbit || bernomor || tanpaPeserta}')
  })

  it('sebab tombolnya mati DIKATAKAN, bukan dibiarkan ditebak', () => {
    expect(SUMBER).toContain('DETAIL_POLIS.sudahBernomor')
    expect(SUMBER).toContain('DETAIL_POLIS.perluPeserta')
    // Kalimatnya menyebut apa yang harus dikerjakan lebih dahulu.
    expect(DETAIL_POLIS.perluPeserta).toContain('Unggah rincian peserta')
  })

  it('selisih kolom dijawab di layar, bukan hanya di komentar Go', () => {
    // ⛔ Siapa pun yang menghitung kolomnya akan bertanya kenapa empat puluh
    // medan layar lama menjadi tiga puluh delapan; jawaban yang hanya ada di
    // kode bukan jawaban bagi yang bertanya.
    expect(SUMBER).toContain('kepala.medanTanpaKolom')
    expect(SUMBER).toContain('medan layar lama tidak ditampilkan')
  })

  it('layar TIDAK merakit bentuk nomor sendiri', () => {
    // ⛔ AC tiket 03: nomornya dirakit di `models`, dari bahan yang seluruhnya
    // dibaca server. Layar yang menyusun bentuknya sendiri akan berselisih
    // dengan server tepat saat bentuknya berubah.
    expect(SUMBER).not.toContain('QR/QP/TP/TR')
    expect(SUMBER).not.toMatch(/padStart\(\s*5/)
  })
})
