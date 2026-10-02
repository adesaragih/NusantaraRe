import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { DETAIL_POLIS, JUDUL_KOLOM_PESERTA } from '../labels'
import { judulKolom, kalimatNomor, kolomBerisi, selPeserta, susunKolom } from './PremiumListDetail'

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
    expect(SUMBER).toContain('susunKolom(kolomBerisi(hal.kolom, hal.baris))')
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

  it('NAME_OF_INSURED, DOB, dan dua Gross Valuation tampil (keputusan work owner 02-10-2026)', () => {
    // Kolomnya tetap DATANG DARI SERVER (`models.KolomGridTambahan`); layar
    // hanya memberi judul dan urutan, dan tidak mengetik namanya di berkas layar.
    expect(JUDUL_KOLOM_PESERTA.NAME_OF_INSURED).toBe('Name Of Insured')
    expect(JUDUL_KOLOM_PESERTA.DOB).toBe('DOB')
    expect(JUDUL_KOLOM_PESERTA.GROSS_VALUATION_BEGIN_DATE).toBe('Gross Valuation Begin Date')
    expect(JUDUL_KOLOM_PESERTA.GROSS_VALUATION_EXPIRED_DATE).toBe('Gross Valuation Expired Date')
    expect(SUMBER).not.toContain('NAME_OF_INSURED')
    expect(susunKolom(['NAME_OF_INSURED', 'PLAN', 'CERTIFICATE_NO'])).toEqual(['CERTIFICATE_NO', 'NAME_OF_INSURED', 'PLAN'])
  })
})

describe('hanya kolom yang berisi (02-10-2026)', () => {
  const baris = [
    { nilai: { PLAN: 'UJI-PLAN', STNC: '', CLAIM: '0' } },
    { nilai: { PLAN: '', STNC: '   ', CLAIM: '0' } },
  ]

  it('kolom kosong di semua baris disembunyikan, urutan server tetap', () => {
    expect(kolomBerisi(['STNC', 'PLAN', 'LAPSE_DATE', 'CLAIM'], baris)).toEqual(['PLAN'])
  })

  it('kolom yang seluruhnya nol disembunyikan (keputusan work owner 02-10-2026)', () => {
    expect(kolomBerisi(['CLAIM'], baris)).toEqual([])
    const ragam = [{ nilai: { A: '0.00', B: '-0', C: '.0', D: '0.5', E: '10', F: '100.00' } }]
    expect(kolomBerisi(['A', 'B', 'C', 'D', 'E', 'F'], ragam)).toEqual(['D', 'E', 'F'])
  })

  it('kolom tampil bila SATU baris saja bukan nol; sel nolnya tetap ditulis 0', () => {
    expect(kolomBerisi(['CLAIM'], [...baris, { nilai: { CLAIM: '12.5' } }])).toEqual(['CLAIM'])
    expect(selPeserta('0')).toBe('0')
  })

  it('jumlah kolom tersembunyi dinyatakan di layar', () => {
    expect(SUMBER).toContain('empty columns hidden')
  })
})

describe('STNC dan WPC paling kanan (02-10-2026)', () => {
  it('urutan kolom CSV ceding, kolom lain sesudahnya, STNC lalu WPC di ujung kanan', () => {
    expect(susunKolom(['BEGIN_DATE', 'WPC', 'PLAN', 'STNC', 'RATE', 'POLICY_HOLDER', 'GROSS_PREMIUM'])).toEqual([
      'POLICY_HOLDER',
      'PLAN',
      'BEGIN_DATE',
      'GROSS_PREMIUM',
      'RATE',
      'STNC',
      'WPC',
    ])
  })

  it('kolom yang tersembunyi tidak dimunculkan kembali', () => {
    expect(susunKolom(['PLAN', 'WPC'])).toEqual(['PLAN', 'WPC'])
    expect(susunKolom(['PLAN'])).toEqual(['PLAN'])
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

  it('kepala HANYA nomor kasus, Period, dan PL Number (02-10-2026)', () => {
    expect(SUMBER).toContain('kepala.polisId')
    expect(SUMBER).toContain('periodeTampil(periode)')
    expect(SUMBER).toContain('DETAIL_POLIS.nomor')
    expect(SUMBER).not.toContain('kepala.type')
    expect(SUMBER).not.toContain('kepala.businessCode')
  })

  it('TANPA tombol Generate PL Number (keputusan work owner 01-10-2026)', () => {
    // Nomor PL terbit saat polis disimpan, bukan lewat tombol; di
    // ShowLifePremiumDetail sel PL_NUMBER pun `pyVisible never`.
    expect(SUMBER).not.toContain('DETAIL_POLIS.terbitkan')
    expect(SUMBER).not.toContain('terbitkanNomorPL')
  })

  it('selisih kolom vs PL_Detail_Sec TIDAK ditampilkan (keputusan work owner 02-10-2026)', () => {
    // Catatan pengembang, bukan informasi pemakai; jawabannya tetap di
    // `models.MedanGridTanpaKolom`.
    expect(SUMBER).not.toContain('legacy screen fields not shown')
    expect(SUMBER).not.toContain('pl-detail__absen')
  })

  it('layar TIDAK merakit bentuk nomor sendiri', () => {
    // ⛔ AC tiket 03: nomornya dirakit di `models`, dari bahan yang seluruhnya
    // dibaca server. Layar yang menyusun bentuknya sendiri akan berselisih
    // dengan server tepat saat bentuknya berubah.
    expect(SUMBER).not.toContain('QR/QP/TP/TR')
    expect(SUMBER).not.toMatch(/padStart\(\s*5/)
  })
})
