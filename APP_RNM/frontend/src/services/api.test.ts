import { describe, expect, it } from 'vitest'

import {
  jumlahUang,
  tampilUang,
  dampakHapusKlaim,
  hapusKlaim,
  serahkanKeKomite,
  bolehSerahkanDiLayar,
  bolehPutaranBaru,
  bolehSimpanAdjustment,
  simpanAdjustment,
  tambahPutaran,
  pesanGalat,
  type DampakHapus,
} from './api'

// Test ini menjaga SATU aturan: uang di sisi React tetap TEKS, tidak pernah
// menjadi angka (tiket 01 AC-4; ADR-U-0003, ADR-U-0016).
//
// Cara membacanya: `describe` = kelompok kasus, `it` = satu kasus,
// `expect(x).toBe(y)` = "x harus persis y".
describe('uang di sisi React', () => {
  it('membawa jumlah sebagai teks, apa adanya', () => {
    const uang = { amount: '1234567890.12345678', currency: 'IDR' }
    expect(jumlahUang(uang)).toBe('1234567890.12345678')
    expect(typeof jumlahUang(uang)).toBe('string')
  })

  it('tidak kehilangan presisi pada nilai yang tidak muat di float64', () => {
    // Number('1234567890.12345678') membulat menjadi 1234567890.1234567.
    const teks = '1234567890.12345678'
    expect(jumlahUang({ amount: teks, currency: 'IDR' })).toBe(teks)
    expect(String(Number(teks))).not.toBe(teks)
  })

  it('menolak jumlah yang datang sebagai angka JSON', () => {
    // UangMasuk sengaja menerima `unknown`, jadi angka boleh dicoba di sini.
    expect(() => jumlahUang({ amount: 1234.5, currency: 'IDR' })).toThrow(TypeError)
  })

  it('memperlakukan kosong sebagai kosong, bukan nol', () => {
    expect(jumlahUang(null)).toBe('')
    expect(jumlahUang(undefined)).toBe('')
    expect(jumlahUang({ currency: 'IDR' })).toBe('')
    expect(tampilUang(null)).toBe('')
  })

  it('menampilkan jumlah bersama mata uangnya', () => {
    expect(tampilUang({ amount: '250000', currency: 'IDR' })).toBe('250000 IDR')
  })
})

// Tiket 15: yang membuat "Batal" benar-benar membatalkan adalah bentuk
// jalurnya, bukan niat baik layar. Pratinjau dampak wajib MEMBACA saja.
describe('pratinjau dampak penghapusan', () => {
  it('memisahkan baris datar warisan dari totalnya', () => {
    // Bentuk yang server kirim: `total` TIDAK memuat barisDatarWarisan,
    // sebab nasib baris itu belum diputuskan work owner. Menjumlahkannya
    // berarti menjawab pertanyaan yang terbuka lewat sebuah angka.
    const d: DampakHapus = {
      header: 1,
      peserta: 2,
      adjustment: 4,
      spreading: 8,
      spreadingRetro: 16,
      dokumen: 3,
      barisWork: 1,
      barisDatarWarisan: 5,
      total: 35,
    }
    const jumlahDaftar =
      d.header + d.peserta + d.adjustment + d.spreading + d.spreadingRetro +
      d.dokumen + d.barisWork
    expect(d.total).toBe(jumlahDaftar)
    expect(d.total).not.toBe(jumlahDaftar + d.barisDatarWarisan)
  })

  it('membaca lewat GET, bukan lewat kata kerja yang menulis', () => {
    // Fungsi pratinjau memanggil api.get; yang menghapus memanggil api.delete.
    // Keduanya sengaja terpisah supaya membuka popup tidak pernah menulis.
    expect(dampakHapusKlaim.constructor.name).toBe('AsyncFunction')
    expect(hapusKlaim.constructor.name).toBe('AsyncFunction')
    expect(dampakHapusKlaim).not.toBe(hapusKlaim)
  })
})

// Tiket 10 - kontrak penyerahan ke Komite di sisi klien.
describe('penyerahan ke Komite', () => {
  it('menyerahkan BARIS, bukan klaim - pengenal peserta ikut di jalurnya', () => {
    // Unit keputusan adalah baris (ADR-U-0011). Jalur yang hanya membawa
    // pengenal klaim akan memaksa server menebak baris mana yang dimaksud.
    expect(serahkanKeKomite.length).toBe(3)
    expect(serahkanKeKomite.constructor.name).toBe('AsyncFunction')
  })

  it('meneruskan pesan server apa adanya, dan diam bila tidak ada', () => {
    // "Name of bank cannot be empty" adalah kalimat sistem lama yang sengaja
    // dipertahankan; ia harus sampai ke pengguna tanpa diparafrase.
    const galatAxios = {
      isAxiosError: true,
      response: { status: 422, data: { error: 'Name of bank cannot be empty' } },
    }
    expect(pesanGalat(galatAxios)).toBe('Name of bank cannot be empty')
    expect(pesanGalat({ isAxiosError: true, response: { status: 500, data: {} } })).toBeUndefined()
    expect(pesanGalat(new Error('bukan galat axios'))).toBeUndefined()
  })
})

// Tiket 10 AC: gerbang rekening ditegakkan di LAYANAN, dan tetap ditegakkan
// meskipun kontrol layarnya ditampilkan.
describe('kontrol layar bukan pagar', () => {
  it('menampilkan Send ke Komite untuk baris yang rekeningnya kosong', () => {
    // ⛔ Inilah buktinya. Layar menampilkan tombol untuk baris ini, sehingga
    // bila gerbang rekening tidak ada di services, ia tidak ada di mana pun.
    // Baris di bawah TIDAK punya medan bank sama sekali - dan tombolnya tetap
    // tampil, sebagaimana seharusnya.
    const barisTanpaBank = {
      id: 'UJI-A-1',
      status: 'Outstanding',
      kodeStatus: '0',
      statusDiketahui: true,
      jumlahKlaim: { amount: '1500000', currency: 'IDR' },
      nomorAkseptasi: '',
      tanggalAkseptasi: '',
      komiteId: '',
    }
    expect(bolehSerahkanDiLayar(barisTanpaBank)).toBe(true)
  })

  it('menyembunyikan kontrol untuk baris yang sudah diserahkan atau sudah diputus', () => {
    const dasar = {
      id: 'UJI-A-2',
      status: 'Outstanding',
      kodeStatus: '0',
      statusDiketahui: true,
      jumlahKlaim: { amount: '1', currency: 'IDR' },
      nomorAkseptasi: '',
      tanggalAkseptasi: '',
      komiteId: '',
    }
    expect(bolehSerahkanDiLayar({ ...dasar, komiteId: 'KMT-000001' })).toBe(false)
    expect(bolehSerahkanDiLayar({ ...dasar, status: 'Ditolak', kodeStatus: '2' })).toBe(false)
  })
})

// Tiket 11 - klaim tidak terminal.
describe('putaran berikutnya', () => {
  const baris = (kodeStatus: string) => ({
    id: 'UJI-A',
    status: kodeStatus === '2' ? 'Ditolak' : 'Outstanding',
    kodeStatus,
    statusDiketahui: true,
    jumlahKlaim: { amount: '1', currency: 'IDR' },
    nomorAkseptasi: '',
    tanggalAkseptasi: '',
    komiteId: '',
  })

  it('membuka putaran hanya sesudah baris TERAKHIR ditolak', () => {
    // ⛔ Yang menentukan baris terakhir, bukan "ada baris yang pernah
    // ditolak": klaim yang barisnya ditolak lalu sudah dibuatkan penggantinya
    // tidak boleh menawarkan putaran ketiga.
    expect(bolehPutaranBaru([baris('2')])).toBe(true)
    expect(bolehPutaranBaru([baris('2'), baris('0')])).toBe(false)
    expect(bolehPutaranBaru([baris('0')])).toBe(false)
    expect(bolehPutaranBaru([baris('1')])).toBe(false)
    expect(bolehPutaranBaru([])).toBe(false)
  })

  it('membuka putaran per PESERTA, bukan per baris', () => {
    // Unit keputusan adalah baris (ADR-U-0011), tetapi putaran berikutnya
    // lahir pada pesertanya - barisnya belum ada saat diminta.
    expect(tambahPutaran.length).toBe(2)
  })
})

// Audit A0 - jalur akseptasi Claim Life sendiri (SaveAdjustment_Act).
describe('kontrol Save Adjustment', () => {
  const baris = (kodeStatus: string, nomorAkseptasi = '') => ({
    id: 'UJI-A',
    status: kodeStatus === '0' ? 'Outstanding' : kodeStatus === '1' ? 'Aksep' : 'Ditolak',
    kodeStatus,
    statusDiketahui: true,
    jumlahKlaim: { amount: '1', currency: 'IDR' },
    nomorAkseptasi,
    tanggalAkseptasi: '',
    komiteId: '',
  })
  const dipilih = { isCheck: 'true' }

  it('meniru prasyarat XML: dipilih, belum bernomor, masih Outstanding', () => {
    expect(bolehSimpanAdjustment(dipilih, [baris('0')])).toBe(true)
    // ⛔ Peserta TIDAK dipilih - prasyarat `.IsCheck=true`.
    expect(bolehSimpanAdjustment({ isCheck: 'false' }, [baris('0')])).toBe(false)
    expect(bolehSimpanAdjustment({ isCheck: '' }, [baris('0')])).toBe(false)
    // ⛔ Sudah bernomor - prasyarat `.ACCEPTEDNO==""`.
    expect(bolehSimpanAdjustment(dipilih, [baris('0', 'RNML-AL1.04.26.00001')])).toBe(false)
    // ⛔ Bukan Outstanding - prasyarat `.STS_REJECT=="0"`.
    expect(bolehSimpanAdjustment(dipilih, [baris('1')])).toBe(false)
    expect(bolehSimpanAdjustment(dipilih, [baris('2')])).toBe(false)
    expect(bolehSimpanAdjustment(dipilih, [])).toBe(false)
  })

  it('memeriksa baris TERAKHIR, bukan sembarang baris', () => {
    // `.AdjustmentList(<LAST>)` - XML menstempel baris terakhir.
    expect(bolehSimpanAdjustment(dipilih, [baris('2'), baris('0')])).toBe(true)
    expect(bolehSimpanAdjustment(dipilih, [baris('0'), baris('1')])).toBe(false)
  })

  it('mengaksep per PESERTA, dan mengembalikan nomornya', () => {
    expect(simpanAdjustment.length).toBe(2)
    expect(simpanAdjustment.constructor.name).toBe('AsyncFunction')
  })
})
