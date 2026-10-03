// Logika murni layar - aturan XML yang Pega jalankan di klien.

import { describe, expect, it } from 'vitest'

import {
  BATAS_DROPDOWN,
  PILIHAN_PEMBAYARAN,
  geserAktif,
  hitungMaxSumReasured,
  idBolehDisalin,
  ringkasSalin,
  saringLama,
  jepitHalaman,
  namaTreaty,
  potongHalaman,
  potongPilihan,
  produkBaru,
  salinBaris,
  salinProduk,
  tampilPremiumFactor,
  tampilViewOffice,
  tampilAngka,
  tampilTanggal,
  tampilViewRate,
  waktuPega,
} from './bentuk'
import { LAIN_MPNL } from './labels'
import type { ProdukLama } from './api'

describe('CountMaxSumReasured_Act - MAXSUMREASURED = MaxSumInsured - CedingLimit (eksak)', () => {
  it('mengurangi per digit, tanpa float', () => {
    expect(hitungMaxSumReasured('1000000000.30', '150000000.10')).toBe('850000000.20')
    expect(hitungMaxSumReasured('9007199254740993.000000002', '0.000000001')).toBe('9007199254740993.000000001')
    expect(hitungMaxSumReasured('0.3', '0.1')).toBe('0.2')
  })

  it('kosong = 0 (`local.* = 0` b239); hasil negatif apa adanya', () => {
    expect(hitungMaxSumReasured('', '')).toBe('0')
    expect(hitungMaxSumReasured('', '5')).toBe('-5')
    expect(hitungMaxSumReasured('5', '-2')).toBe('7')
  })

  it('koma desimal diterima seperti server; isian bukan angka = medan dibiarkan', () => {
    expect(hitungMaxSumReasured('12,5', '0,5')).toBe('12.0')
    expect(hitungMaxSumReasured('abc', '1')).toBeNull()
    expect(hitungMaxSumReasured('1e3', '1')).toBeNull()
  })
})

describe('SetTreatyName_Act - INWARDNAME = PRODUCTNAME + " " + POLICYHODERNAME', () => {
  it('tanpa pemangkasan, seperti Pega', () => {
    expect(namaTreaty('UJI PRODUK', 'UJI PEMEGANG')).toBe('UJI PRODUK UJI PEMEGANG')
    expect(namaTreaty('UJI PRODUK', '')).toBe('UJI PRODUK ')
  })
})

describe('CopyProduct - kedua ID dikosongkan, halaman tetap', () => {
  it('produk tersimpan → salinan baru bermarka salinanDari', () => {
    const asal = { ...produkBaru(), id: '100007' }
    asal.inward = { ...asal.inward, id: '100007', productId: '100007', insured: 'UJI' }
    asal.planList = [{ plan: 'UJI PLAN', planId: 'P1', name: '', benefit: '', riRate: '', riRateId: '' }]
    const s = salinProduk(asal)
    expect([s.id, s.inward.id, s.inward.productId, s.salinanDari]).toEqual(['', '', '', '100007'])
    expect(s.inward.insured).toBe('UJI')
    expect(s.planList).toEqual(asal.planList)
    // `Copy` dua kali: tetap menunjuk produk asal.
    expect(salinProduk(s).salinanDari).toBe('100007')
  })
})

describe('salin baris CopyFinancialWriting / CopyUnderWritingLimit', () => {
  it('menambah salinan di akhir, tanpa kunci asing (`asli`)', () => {
    const d = [{ minInsured: '1', maxInsured: '2', employee: 'a', nonEmployee: 'b', asli: '{"X":1}' }]
    expect(salinBaris(d, 0)).toEqual([d[0], { minInsured: '1', maxInsured: '2', employee: 'a', nonEmployee: 'b' }])
    expect(salinBaris(d, 5)).toEqual(d)
  })
})

describe('visibilitas XML', () => {
  it('Premium Factor hanya PAYMENT==3; View Rate bila RIRATE terisi; View Office Online untuk jenis kantor', () => {
    expect(tampilPremiumFactor('3')).toBe(true)
    expect(tampilPremiumFactor('1')).toBe(false)
    expect(tampilViewRate('')).toBe(false)
    expect(tampilViewRate('UJI RATE')).toBe(true)
    expect(['xls', 'xlsx', 'doc', 'docx', 'ppt', 'pptx'].every(tampilViewOffice)).toBe(true)
    expect(tampilViewOffice('pdf')).toBe(false)
  })

  it('pilihan pembayaran = kode 1-4 GenerateUpload_Act b1141', () => {
    expect(PILIHAN_PEMBAYARAN.map((o) => [o.value, o.label])).toEqual([
      ['1', 'Annual'],
      ['2', 'Semi Annual'],
      ['3', 'Quarterly'],
      ['4', 'Monthly'],
    ])
  })
})

describe('pyGridPaginator - 10 baris per halaman', () => {
  it('memotong dan menjepit halaman', () => {
    const semua = Array.from({ length: 23 }, (_, i) => i)
    expect(potongHalaman(semua, 3)).toEqual([20, 21, 22])
    expect(jepitHalaman(9, 23)).toBe(3)
    expect(jepitHalaman(0, 0)).toBe(1)
  })
})

describe('waktuPega - `.Date` History b62561 (pxDateTime)', () => {
  it('stempel Pega GMT tampil DD-MM-YYYY HH:mm waktu lokal (audit 02-10-2026)', () => {
    // WIB = GMT+7 (420 menit), dan GMT sendiri.
    expect(waktuPega('20261002T031500.000 GMT', 420)).toBe('02-10-2026 10:15')
    expect(waktuPega('20261002T031500.000 GMT', 0)).toBe('02-10-2026 03:15')
    // Lewat tengah malam: tanggalnya ikut maju.
    expect(waktuPega('20261001T200000.000 GMT', 420)).toBe('02-10-2026 03:00')
  })
  it('teks lain tampil apa adanya', () => {
    expect(waktuPega('', 420)).toBe('')
    expect(waktuPega('bukan stempel', 420)).toBe('bukan stempel')
  })
})

describe('dropdown master - pengganti tombol Choose* (keputusan work owner 02-10-2026)', () => {
  it('server diminta BATAS_DROPDOWN + 1 baris: baris lebih = daftar terpotong, dinyatakan', () => {
    const semua = Array.from({ length: BATAS_DROPDOWN + 1 }, (_, i) => i)
    expect(potongPilihan(semua)).toEqual({ tampil: semua.slice(0, BATAS_DROPDOWN), lebih: true })
    expect(potongPilihan(semua.slice(0, BATAS_DROPDOWN))).toEqual({ tampil: semua.slice(0, BATAS_DROPDOWN), lebih: false })
    expect(potongPilihan([])).toEqual({ tampil: [], lebih: false })
    // Kalimat potongan menyebut batas yang sama.
    expect(LAIN_MPNL.dropdownTerpotong).toContain(String(BATAS_DROPDOWN))
  })

  it('panah atas/bawah dan PageUp/PageDown menggeser pilihan aktif tanpa keluar daftar', () => {
    expect(geserAktif(-1, 1, 3)).toBe(0)
    expect(geserAktif(0, 1, 3)).toBe(1)
    expect(geserAktif(2, 1, 3)).toBe(2)
    expect(geserAktif(0, -1, 3)).toBe(0)
    expect(geserAktif(1, -Infinity, 3)).toBe(0)
    expect(geserAktif(1, Infinity, 3)).toBe(2)
    expect(geserAktif(5, 0, 3)).toBe(2)
    expect(geserAktif(0, 1, 0)).toBe(-1)
  })
})

describe('mode lihat - nilai tampil seperti layar Pega (foto layar work owner 02-10-2026)', () => {
  it('angka: pemisah ribuan titik, desimal koma - eksak dari teks, tanpa float', () => {
    expect(tampilAngka('250000000')).toBe('250.000.000')
    expect(tampilAngka('1500000.30')).toBe('1.500.000,30')
    expect(tampilAngka('-1234.5')).toBe('-1.234,5')
    // Lebih dari 2^53 tetap eksak (bukan float): setiap digit utuh, berkelompok tiga. Hasil tidak ditulis literal -
    // empat kelompok angka bertitik dibaca penjaga TestMPNLNolAlamatLayanan sebagai alamat IP.
    const besar = tampilAngka('9007199254740993')
    expect(besar.split('.').join('')).toBe('9007199254740993')
    expect(besar.split('.').map((k) => k.length)).toEqual([1, 3, 3, 3, 3, 3])
    expect(tampilAngka('180')).toBe('180')
    expect(tampilAngka('0')).toBe('0')
    // Teks yang bukan angka kanonik (data lama) tampil apa adanya - tidak ditebak.
    expect(tampilAngka('12,5')).toBe('12,5')
    expect(tampilAngka('')).toBe('')
  })

  it('tanggal YYYY-MM-DD tampil DD/MM/YYYY (`Begin Date` 01/08/2023 di Pega)', () => {
    expect(tampilTanggal('2023-08-01')).toBe('01/08/2023')
    expect(tampilTanggal('')).toBe('')
    expect(tampilTanggal('bukan tanggal')).toBe('bukan tanggal')
  })
})

describe('Copy Old - popup produk lama (permintaan work owner 03-10-2026)', () => {
  const lama = (id: string, ubah: Partial<ProdukLama> = {}): ProdukLama => ({
    id, productName: `UJI PRODUK ${id}`, ceding: 'UJI CEDING', treatyNumber: `UJI-${id}`, inwardName: 'UJI TREATY',
    createOp: 'UJI-OP', updateOp: 'UJI-OP', bolehDisalin: true, alasan: [], catatan: [], ...ubah,
  })
  const daftar = [lama('100901'), lama('100902', { ceding: 'UJI LAIN', bolehDisalin: false, alasan: ['UJI ALASAN'] }), lama('100903')]

  it('Search menyaring ID, Product Name, Ceding, Treaty Number, Treaty Name - tanpa membedakan huruf', () => {
    expect(saringLama(daftar, '').map((d) => d.id)).toEqual(['100901', '100902', '100903'])
    expect(saringLama(daftar, 'lain').map((d) => d.id)).toEqual(['100902'])
    expect(saringLama(daftar, ' uji-100903 ').map((d) => d.id)).toEqual(['100903'])
    expect(saringLama(daftar, '10090').length).toBe(3)
  })

  it('yang dikirim Process Copy: hanya ID terpilih yang boleh disalin, urutan daftar', () => {
    expect(idBolehDisalin(daftar, new Set(['100903', '100902', '100901', 'UJI-HILANG']))).toEqual(['100901', '100903'])
    expect(idBolehDisalin(daftar, new Set())).toEqual([])
  })

  it('ringkasan hasil per status', () => {
    expect(
      ringkasSalin([
        { id: '1', status: 'disalin', pesan: [] },
        { id: '2', status: 'disalin', pesan: [] },
        { id: '3', status: 'ditolak', pesan: ['x'] },
        { id: '4', status: 'gagal', pesan: [] },
        { id: '5', status: 'sudahAda', pesan: [] },
      ]),
    ).toEqual({ disalin: 2, sudahAda: 1, ditolak: 1, gagal: 1 })
  })
})

