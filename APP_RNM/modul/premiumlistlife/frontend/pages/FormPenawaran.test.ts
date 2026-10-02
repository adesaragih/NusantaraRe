import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { isiDariPenawaran, tanggalMasukan, type PenawaranPolis } from '../api'
import { LABEL_PENAWARAN } from '../labels'
import { hitungMaxTBC, jenisAsuransiDari, kolomWajibKosong, opsiDari, saringAngkaDesimal, tampilTBC, tanggalRiwayat, terapkanPilihan } from './FormPenawaran'

// Uji form penawaran — tiket 01 bagian 3.

const BERKAS = readFileSync(join(__dirname, 'FormPenawaran.tsx'), 'utf8')
const INPUT_OFFER = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')

const contoh: PenawaranPolis = {
  caseId: 'NBLF-1',
  flag: '0',
  tahap: 'Input Offer Life',
  bolehDisimpan: true,
  noOffer: '',
  cedingCo: 'UJI-L01',
  cedingCoName: 'UJI-CEDING',
  policyHolder: '',
  policyHolderName: '',
  typeCeding: '4',
  typeCedingName: 'XOL',
  jenisAsuransi: 'Non Proportional',
  businessCode: 'L2',
  businessName: 'GROUP TERM LIFE',
  dateReceived: '2026-09-30T00:00:00Z',
  description: 'UJI',
  batasUsiaPeserta: 62,
  periodePertanggungan: '1 YEAR',
  sumInsured: '6000000000',
  tanggalPenawaran: '2026-07-30T00:00:00Z',
  tanggalRespon: null,
  tanggalKonfirmasi: '2026-09-29T00:00:00Z',
  tbc: 25,
  tanggalTbc: '2026-10-24T00:00:00Z',
  statusUpdate: 'UJI-STATUS',
  keteranganMarketing: 'UJI-CATATAN',
  qqName: '',
  jenisUsaha: '',
  ketentuanUnderwriting: '',
  tanggalKonfirmasiBalik: null,
  tanggalRealisasi: null,
  tanggalBind: null,
  statusFinal: '',
  status: 'Pending',
  riwayat: [],
  pilihan: {
    typeCeding: [{ kode: '1', nama: 'QS' }],
    classOfBusiness: [],
    status: [],
  },
}

describe('isian dari server', () => {
  it('badan simpan membawa KODE, bukan nama turunan', () => {
    const isi = isiDariPenawaran(contoh)
    expect(isi.typeCeding).toBe('4')
    expect(isi.businessCode).toBe('L2')
    expect(Object.keys(isi)).not.toContain('typeCedingName')
    expect(Object.keys(isi)).not.toContain('businessName')
    expect(isi.dateReceived).toBe('2026-09-30')
  })

  it('tanggal kosong tetap kosong, bukan tanggal nol', () => {
    expect(tanggalMasukan(null)).toBe('')
    expect(tanggalMasukan('')).toBe('')
  })

  it('pilihan server menjadi opsi apa adanya', () => {
    expect(opsiDari([{ kode: '1', nama: 'QS' }])).toEqual([{ value: '1', label: 'QS' }])
  })
})

describe('popup pilihan (setCeding_act / setPolicyHolder_act)', () => {
  it('kode dan nama berpindah bersama', () => {
    const isi = isiDariPenawaran(contoh)
    const c = terapkanPilihan(isi, 'ceding', { id: 'UJI-L02', nama: 'UJI-CEDING-2' })
    expect([c.cedingCo, c.cedingCoName]).toEqual(['UJI-L02', 'UJI-CEDING-2'])
    const p = terapkanPilihan(isi, 'pemegang', { id: 'UJI-C1', nama: 'UJI-PEMEGANG' })
    expect([p.policyHolder, p.policyHolderName]).toEqual(['UJI-C1', 'UJI-PEMEGANG'])
    expect(p.cedingCo).toBe('UJI-L01')
  })
})

describe('riwayat', () => {
  it('tanggal nol Go tampil sebagai tanda kosong', () => {
    expect(tanggalRiwayat('0001-01-01T00:00:00Z')).toBe('—')
    expect(tanggalRiwayat('2026-09-30T10:15:00Z')).toBe('30/09/2026 10:15')
  })
})

describe('struktur layar', () => {
  it('nol daftar pilihan ditulis di layar — datang dari server', () => {
    expect(BERKAS).not.toContain("'QS'")
    expect(BERKAS).not.toContain('INDIVIDUAL TERM LIFE')
    expect(BERKAS).toContain('data.pilihan.typeCeding')
  })

  it('form hanya di tahap penawaran', () => {
    expect(INPUT_OFFER).toContain('tahap === TAHAP_POLIS.penawaran && <FormPenawaran')
  })

  it('label VERBATIM InputOfferLife.xml', () => {
    expect(LABEL_PENAWARAN.typeCeding).toBe('System Reinsurance')
    expect(LABEL_PENAWARAN.businessCode).toBe('Class of Business')
    expect(LABEL_PENAWARAN.simpan).toBe('Save Offer')
  })
})

describe('sel penawaran migrasi 059', () => {
  it('angka dan uang dikirim sebagai TEKS, tanggal sebagai YYYY-MM-DD', () => {
    const isi = isiDariPenawaran(contoh)
    expect(isi.batasUsiaPeserta).toBe('62')
    expect(isi.tbc).toBe('25')
    expect(isi.sumInsured).toBe('6000000000')
    expect(isi.tanggalKonfirmasi).toBe('2026-09-29')
    expect(isi.tanggalRespon).toBe('')
  })

  it('Input TBC hanya tampil bila Confirmation Date terisi', () => {
    const isi = isiDariPenawaran(contoh)
    expect(tampilTBC(isi)).toBe(true)
    expect(tampilTBC({ ...isi, tanggalKonfirmasi: '' })).toBe(false)
  })

  it('Sum Insured bukan input number — uang tetap teks', () => {
    expect(BERKAS).not.toMatch(/sumInsured[\s\S]{0,120}type="number"/)
  })

  it('kasus Input Premium TETAP dapat mengubah sel penawaran (keputusan work owner 02-10-2026)', () => {
    // Yang mengunci hanya tahapnya; FlagOnGoingPolicy tidak lagi mengunci.
    expect(BERKAS).toContain('const kunci = !bisa\n')
    expect(BERKAS).not.toContain("kunci = !bisa || data.flag === '1'")
  })
})

describe('sel penawaran migrasi 060', () => {
  it('sel layar Pega tampil, kecuali yang disembunyikan work owner', () => {
    for (const medan of ['ketentuanUnderwriting', 'tanggalKonfirmasiBalik',
      'tanggalRealisasi', 'tanggalBind', 'statusFinal'] as const) {
      expect(BERKAS).toContain(`LABEL_PENAWARAN.${medan}`)
      expect(isiDariPenawaran(contoh)[medan]).toBe('')
    }
  })

  it('Insured Name dan Occupation disembunyikan, tetapi nilainya tetap terkirim', () => {
    expect(BERKAS).not.toContain('LABEL_PENAWARAN.qqName')
    expect(BERKAS).not.toContain('LABEL_PENAWARAN.jenisUsaha')
    const isi = isiDariPenawaran({ ...contoh, qqName: 'UJI-TERTANGGUNG', jenisUsaha: 'UJI-USAHA' })
    expect(isi.qqName).toBe('UJI-TERTANGGUNG')
    expect(isi.jenisUsaha).toBe('UJI-USAHA')
  })
})

describe('Max TBC = Confirmation Date + Input TBC hari (SetMaxTBCLife_Act)', () => {
  it('contoh layar Pega: 29/09/2026 + 25 = 24/10/2026', () => {
    expect(hitungMaxTBC('2026-09-29', '25')).toBe('24/10/2026')
  })

  it('melewati akhir tahun dan tahun kabisat', () => {
    expect(hitungMaxTBC('2026-12-20', '15')).toBe('04/01/2027')
    expect(hitungMaxTBC('2028-02-28', '1')).toBe('29/02/2028')
    expect(hitungMaxTBC('2026-09-29', '0')).toBe('29/09/2026')
  })

  it('kosong bila salah satunya kosong atau bukan bilangan bulat', () => {
    expect(hitungMaxTBC('', '25')).toBe('')
    expect(hitungMaxTBC('2026-09-29', '')).toBe('')
    expect(hitungMaxTBC('2026-09-29', '2,5')).toBe('')
    expect(hitungMaxTBC('2026-09-29', '-3')).toBe('')
  })

  it('layar menghitung dari isian yang diketik, read-only', () => {
    expect(BERKAS).toContain('hitungMaxTBC(isi.tanggalKonfirmasi, isi.tbc)')
  })
})

describe('Reinsurance Type turunan System Reinsurance (SetReinsuranceType)', () => {
  it('XOL -> Non Proportional, selain itu Proportional', () => {
    expect(jenisAsuransiDari('4')).toBe('Non Proportional')
    expect(jenisAsuransiDari('1')).toBe('Proportional')
    expect(jenisAsuransiDari('')).toBe('Proportional')
  })

  it('read-only, dihitung dari pilihan yang sedang dipilih, tidak dikirim', () => {
    expect(BERKAS).toContain('jenisAsuransiDari(isi.typeCeding)')
    expect(Object.keys(isiDariPenawaran(contoh))).not.toContain('jenisAsuransi')
  })
})

describe('kolom wajib dan tombol Save Offer', () => {
  const lengkap = {
    ...isiDariPenawaran(contoh),
    policyHolder: 'UJI-C1',
    policyHolderName: 'UJI-PEMEGANG',
    status: 'Pending',
  }

  it('isian lengkap — tidak ada yang kurang', () => {
    expect(kolomWajibKosong(lengkap)).toEqual([])
  })

  it('melaporkan seluruh yang kosong, urut seperti layar', () => {
    const kosong = {
      ...lengkap,
      cedingCoName: '',
      policyHolderName: ' ',
      typeCeding: '',
      businessCode: '',
      status: '',
      description: '',
    }
    expect(kolomWajibKosong(kosong)).toEqual([
      'Ceding Name',
      'Policy Holder',
      'System Reinsurance',
      'Class of Business',
      'Status',
      'Comment',
    ])
  })

  it('Email Received Date wajib (keputusan work owner 02-10-2026)', () => {
    expect(kolomWajibKosong({ ...lengkap, dateReceived: '' })).toEqual([LABEL_PENAWARAN.dateReceived])
    expect(kolomWajibKosong({ ...lengkap, dateReceived: '', status: '' })).toEqual([
      LABEL_PENAWARAN.dateReceived,
      'Status',
    ])
    expect(BERKAS).toContain("required={medan === 'dateReceived'}")
  })

  it('Status wajib — radio belum dipilih menahan simpan', () => {
    expect(kolomWajibKosong({ ...lengkap, status: '' })).toEqual(['Status'])
  })

  it('tombol terkunci selama ada yang kosong, dan simpan() pun menolak', () => {
    expect(BERKAS).toContain('disabled={sibuk || kurang.length > 0}')
    expect(BERKAS).toContain('kolomWajibKosong(isi).length > 0) return')
  })
})

describe('status tersimpan (migrasi 062)', () => {
  it('radio Status terisi dari status terakhir yang disimpan', () => {
    expect(isiDariPenawaran(contoh).status).toBe('Pending')
  })
})

describe('gaya layar (premiumlistlife.css)', () => {
  const CSS = readFileSync(join(__dirname, '..', 'premiumlistlife.css'), 'utf8')
  const DATA_POLIS = readFileSync(join(__dirname, 'FormDataPolis.tsx'), 'utf8')

  it('setiap kelas pl- yang dipakai layar terdefinisi di CSS modul', () => {
    const dipakai = new Set(BERKAS.match(/pl-offer[\w-]*/g) ?? [])
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) {
      expect(CSS).toContain('.' + k)
    }
  })

  it('setiap tombol memakai kelas dasar btn — tanpa itu tampil sebagai tombol bawaan peramban', () => {
    for (const sumber of [BERKAS, DATA_POLIS]) {
      expect(sumber).not.toMatch(/className="btn--/)
    }
  })

  it('CSS memakai token inti, nol warna heksa', () => {
    expect(CSS).not.toMatch(/#[0-9a-fA-F]{3,6}\b/)
  })
})

describe('Sum Insured hanya angka', () => {
  it('huruf, simbol, dan koma dibuang; satu titik desimal dipertahankan', () => {
    expect(saringAngkaDesimal('6.000.000abc')).toBe('6.000000')
    expect(saringAngkaDesimal('Rp 6,000,000,000')).toBe('6000000000')
    expect(saringAngkaDesimal('1234.56')).toBe('1234.56')
    expect(saringAngkaDesimal('-1e5')).toBe('15')
    expect(saringAngkaDesimal('')).toBe('')
  })

  it('kolom Sum Insured memakai penyaringnya', () => {
    expect(BERKAS).toContain("ubah('sumInsured')(saringAngkaDesimal(v))")
  })
})

describe('dropdown berbahasa Inggris', () => {
  it('setiap Pilih di layar memakai teks kosong "-- choose --"', () => {
    const DATA_POLIS = readFileSync(join(__dirname, 'FormDataPolis.tsx'), 'utf8')
    for (const sumber of [BERKAS, DATA_POLIS]) {
      const pilih = sumber.match(/<Pilih\b/g) ?? []
      const kosong = sumber.match(/kosong=\{TEKS_PILIH\}/g) ?? []
      expect(pilih.length).toBeGreaterThan(0)
      expect(kosong.length).toBe(pilih.length)
    }
  })
})

describe('gaya seluruh layar modul', () => {
  const CSS = readFileSync(join(__dirname, '..', 'premiumlistlife.css'), 'utf8')
  const layar = ['FormDataPolis.tsx', 'PremiumListDetail.tsx', 'UnggahCSVPeserta.tsx', 'InputOffer.tsx']

  it('setiap kelas pl- yang dipakai terdefinisi di premiumlistlife.css', () => {
    for (const berkas of layar) {
      const sumber = readFileSync(join(__dirname, berkas), 'utf8')
      for (const k of new Set(sumber.match(/\bpl-[a-z]+(?:__[a-z-]+)?/g) ?? [])) {
        if (k === 'pl-detail' || k === 'pl-datapolis' || k === 'pl-unggah') continue // pembungkus tanpa gaya
        expect(CSS, `${berkas}: .${k}`).toContain('.' + k)
      }
    }
  })

  it('tidak ada tombol tanpa kelas dasar btn', () => {
    for (const berkas of layar) {
      const sumber = readFileSync(join(__dirname, berkas), 'utf8')
      const tombol = sumber.match(/<button\b[^>]*>/g) ?? []
      for (const t of tombol) {
        expect(t, berkas).toMatch(/className="btn\b/)
      }
    }
  })
})
