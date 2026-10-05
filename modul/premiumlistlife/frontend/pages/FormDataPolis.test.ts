import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { isiDariDataPolis, type DataPolis } from '../api'
import { LABEL_DATA_POLIS } from '../labels'
import { kolomWajibDataPolis, typeRetro } from './FormDataPolis'

// Uji data polis layar Input Premium Detail — tiket 03 bagian 2.

const BERKAS = readFileSync(join(__dirname, 'FormDataPolis.tsx'), 'utf8')
const DETAIL = readFileSync(join(__dirname, 'PremiumListDetail.tsx'), 'utf8')

const contoh: DataPolis = {
  caseId: 'NBLF-1',
  tahap: 'Input Premium Detail',
  bolehDisimpan: true,
  type: 'QR',
  productNameId: 'UJI-P1',
  productName: 'UJI-PRODUK',
  sourceOfBusiness: 'UJI-S1',
  sobName: 'UJI-SOB',
  cedingCo: 'UJI-L01',
  cedingCoName: 'UJI-CEDING',
  policyHolder: 'UJI-C1',
  policyHolderName: 'UJI-PEMEGANG',
  riSlipRnm: '',
  proRateType: '1',
  moId: 'UJI-M1',
  marketingCode: 'UJI-MC',
  marketingName: 'UJI-MARKETING',
  annuityInterest: '0.05',
  premiumRefundFactor: '1',
  retroId: '',
  retroName: '',
  securityReinsurerId: '',
  securityReinsurer: '',
  dateReceived: '2026-10-01T00:00:00Z',
  wpc: null,
  rekapDihapus: false,
  pilihan: { type: [], proRateType: [] },
}

describe('kolom wajib data polis', () => {
  it('QR lengkap — tidak ada yang kurang', () => {
    expect(kolomWajibDataPolis(isiDariDataPolis(contoh))).toEqual([])
  })

  it('TP/TR menuntut R/I SLIP, Billing Name, dan Retrocessionaire', () => {
    const isi = { ...isiDariDataPolis(contoh), type: 'TP' }
    expect(kolomWajibDataPolis(isi)).toEqual([
      LABEL_DATA_POLIS.riSlip,
      LABEL_DATA_POLIS.billing,
      LABEL_DATA_POLIS.retro,
    ])
    expect(typeRetro('TR')).toBe(true)
    expect(typeRetro('QP')).toBe(false)
  })

  it('Email Received Date wajib, dikirim YYYY-MM-DD (02-10-2026)', () => {
    expect(isiDariDataPolis(contoh).dateReceived).toBe('2026-10-01')
    const kosong = isiDariDataPolis({ ...contoh, dateReceived: null })
    expect(kosong.dateReceived).toBe('')
    expect(kolomWajibDataPolis(kosong)).toEqual([LABEL_DATA_POLIS.dateReceived])
    expect(BERKAS).toContain("onChange={ubah('dateReceived')}")
  })

  it('badan simpan tanpa medan milik server', () => {
    const isi = isiDariDataPolis(contoh)
    for (const k of ['caseId', 'tahap', 'bolehDisimpan', 'wpc', 'pilihan']) {
      expect(Object.keys(isi)).not.toContain(k)
    }
  })
})

describe('struktur layar', () => {
  it('Save Data terkunci selama ada yang kosong, dan dinyatakan tanpa perhitungan', () => {
    expect(BERKAS).toContain('disabled={sibuk || kurang.length > 0}')
    // Save Data tidak lagi membawa peringatan batas produk (kini penolakan Validate CSV);
    // Type / Product Name berganti → pesan rekap dihapus (keputusan work owner 05-10-2026).
    expect(BERKAS).not.toContain('data.peringatan')
    expect(BERKAS).toContain('setRekapDihapus(hasil.rekapDihapus === true)')
    expect(BERKAS).toContain('{LABEL_DATA_POLIS.rekapDihapus}')
  })

  it('Choose Product Name hanya sebelum bernomor', () => {
    expect(BERKAS).toContain("!bernomor && tombol(LABEL_DATA_POLIS.pilihProduk, 'produk')")
  })

  it('terpasang di halaman Premium List Detail', () => {
    expect(DETAIL).toContain('<FormDataPolis')
  })

  it('label VERBATIM ShowLifePremiumDetail.xml', () => {
    expect(LABEL_DATA_POLIS.proRateType).toBe('Premium Payment Method')
    expect(LABEL_DATA_POLIS.billing).toBe('Billing Name')
    expect(LABEL_DATA_POLIS.simpan).toBe('Save Data')
  })
})

describe('Billing Name dan Retrocessionaire', () => {
  it('hanya tampil untuk Type TP/TR', () => {
    const awal = BERKAS.search(/\{typeRetro\(isi\.type\) && \(\s*<div className="pl-dp-bagian">/)
    expect(awal).toBeGreaterThan(-1)
    const blok = BERKAS.slice(awal)
    expect(blok).toContain('LABEL_DATA_POLIS.billing')
    expect(blok).toContain('LABEL_DATA_POLIS.retro')
  })
})

describe('Retrocessionaire wajib (TP/TR) dan lebar kolom seragam', () => {
  it('Retrocessionaire masuk daftar wajib untuk TP/TR, tidak untuk QR', () => {
    const tp = { ...isiDariDataPolis(contoh), type: 'TP', riSlipRnm: 'UJI-RNML', retroName: 'UJI-B' }
    expect(kolomWajibDataPolis(tp)).toEqual([LABEL_DATA_POLIS.retro])
    expect(kolomWajibDataPolis(isiDariDataPolis(contoh))).not.toContain(LABEL_DATA_POLIS.retro)
  })

  it('Billing dan Retrocessionaire berbagi satu grid', () => {
    expect(BERKAS).toContain('pl-dp-grid')
  })

  it('R/I SLIP RNM No. di bawah catatan Billing dan di atas Billing Name', () => {
    const catatan = BERKAS.indexOf('LABEL_DATA_POLIS.catatanBilling')
    const riSlip = BERKAS.indexOf('label={LABEL_DATA_POLIS.riSlip}')
    const billing = BERKAS.indexOf('label={LABEL_DATA_POLIS.billing}')
    expect(catatan).toBeGreaterThan(-1)
    expect(riSlip).toBeGreaterThan(catatan)
    expect(billing).toBeGreaterThan(riSlip)
  })
})

describe('tombol View isi Product Name (05-10-2026)', () => {
  it('View tampil di samping Product Name, juga sesudah bernomor, dan membuka popup rincian', () => {
    expect(BERKAS).toContain('{RINCIAN_PRODUK.tombol}')
    expect(BERKAS).toContain("disabled={isi.productNameId.trim() === ''}")
    expect(BERKAS).toContain('<ModalRincianProduk produkID={isi.productNameId}')
    const modal = readFileSync(join(__dirname, '..', 'components', 'ModalRincianProduk.tsx'), 'utf8')
    expect(modal).toContain('ambilRincianProduk(produkID)')
    // View Rate baris PLAN LIST: kunci RIRATEID dari server, tabel rate di bawah grid (05-10-2026).
    expect(modal).toContain('ambilRateProduk(riRateId)')
    // View R/I Risk: kunci riRiskId dari server, view RIRISK_LIFE (05-10-2026).
    expect(modal).toContain('ambilRiskProduk(riRiskId)')
    // Tombol View menempel di kotak R/I Risk Name (05-10-2026).
    expect(modal).toContain("m.label === RINCIAN_PRODUK.medanRisk && (isi.riRiskId ?? '') !== '' ? (")
    expect(modal).toContain('{riskBuka ? RINCIAN_PRODUK.tutupSingkat : RINCIAN_PRODUK.tombol}')
    expect(modal).toContain('{terbuka ? RINCIAN_PRODUK.tutupRate : RINCIAN_PRODUK.viewRate}')
  })
})

describe('tata letak seragam Premium List Detail (02-10-2026)', () => {
  it('Status penawaran tidak ditampilkan (keputusan work owner 02-10-2026)', () => {
    expect(BERKAS).not.toContain('tampil(LABEL_DATA_POLIS.status,')
    expect(BERKAS).not.toContain('offer.status)')
  })

  it('setiap bagian memakai grid seragam; tiga kolom lama tidak dipakai lagi', () => {
    // Empat bagian tetap + satu bagian Retrocession (TP/TR).
    expect(BERKAS.match(/className="pl-dp-grid"/g)?.length).toBe(4)
    expect(BERKAS).not.toContain('pl-datapolis__kolom-tiga')
    expect(BERKAS).not.toContain('pl-offer__pasangan')
  })

  it('tombol pilih menempel di kotaknya dan tetap bernama lengkap untuk pembaca layar', () => {
    // Empat isian bertombol pilih + Product Name (Choose + View, kelas bersyarat - 05-10-2026).
    expect(BERKAS.match(/className="pl-dp-pilih pl-dp-lebar"/g)?.length).toBe(4)
    expect(BERKAS).toContain("'pl-dp-pilih pl-dp-pilih--dua pl-dp-lebar' : 'pl-dp-pilih pl-dp-lebar'")
    expect(BERKAS).toContain('aria-label={label}')
    expect(BERKAS).toContain('{TEKS_TOMBOL_PILIH}')
  })

  it('System Reinsurance di kiri Premium Payment Method (02-10-2026)', () => {
    const sistem = BERKAS.indexOf('tampil(LABEL_DATA_POLIS.typeCeding,')
    const bayar = BERKAS.indexOf('label={LABEL_DATA_POLIS.proRateType}')
    expect(sistem).toBeGreaterThan(-1)
    expect(bayar).toBeGreaterThan(sistem)
  })

  it('tiga kolom berdampingan; di dalamnya grid dua sel, isian panjang selebar kolom (02-10-2026)', () => {
    const css = readFileSync(join(__dirname, '..', 'premiumlistlife.css'), 'utf8')
    expect(css).toMatch(/\.pl-dp-kolom-tiga \{[^}]*repeat\(3, minmax\(0, 1fr\)\)/)
    expect(css).toMatch(/\.pl-dp-grid \{[^}]*repeat\(2, minmax\(0, 1fr\)\)/)
    expect(css).toMatch(/\.pl-dp-lebar \{[^}]*grid-column: 1 \/ -1/)
    expect(BERKAS.match(/className="pl-dp-kolom-tiga"/g)?.length).toBe(1)
  })
})

describe('Underwriting Policy dan Marketing Note sebagai text area (03-10-2026)', () => {
  it('kedua layar memakai AreaTeks untuk keduanya, terkunci pun tetap text area', () => {
    const penawaran = readFileSync(join(__dirname, 'FormPenawaran.tsx'), 'utf8')
    for (const sumber of [BERKAS, penawaran]) {
      expect(sumber).toMatch(/<AreaTeks[\s\S]{0,80}ketentuanUnderwriting/)
      expect(sumber).toMatch(/<AreaTeks[\s\S]{0,80}keteranganMarketing/)
    }
    expect(BERKAS).not.toContain('tampil(LABEL_DATA_POLIS.ketentuanUnderwriting')
    expect(BERKAS).not.toContain('tampil(LABEL_DATA_POLIS.keteranganMarketing')
    const area = readFileSync(join(__dirname, '..', 'components', 'AreaTeks.tsx'), 'utf8')
    expect(area).toContain('<textarea')
    expect(area).toContain('readOnly={readOnly}')
  })
})
