// Uji isian angka berpemisah ribuan — `formatKetik` / `keKabelAngka`.
//
// ⛔ JEBAKAN YANG DIJAGA BERKAS INI
// `modul/treatyinadjustment/.../KerangkaTab.tsx` dahulu membiarkan angka
// TANPA format, dan menuliskan sebabnya: *"memformat di tiap ketukan
// membuat koma desimal mustahil diketik"*. Itu benar untuk `formatNumber`,
// yang membulatkan dan membuang nol ekor. Uji di bawah memaku bahwa
// `formatKetik` TIDAK melakukan keduanya.

import { describe, expect, it } from 'vitest'

import { formatKetik, keKabelAngka, formatNumber, normalisasiKetikan } from './format'

describe('formatKetik — pemisah ribuan saat mengetik', () => {
  it('⭐ mengelompokkan ribuan', () => {
    expect(formatKetik('1000000', 2)).toBe('1.000.000')
    expect(formatKetik('1000', 2)).toBe('1.000')
    expect(formatKetik('999', 2)).toBe('999')
    expect(formatKetik('1000000,5', 2)).toBe('1.000.000,5')
  })

  // ⛔ INI JEBAKANNYA. Keempat keadaan di bawah muncul DI TENGAH
  // pengetikan, dan `formatNumber` menghapus ketiganya yang pertama.
  it('⛔ keadaan tengah-pengetikan BERTAHAN', () => {
    expect(formatKetik('12,', 2)).toBe('12,')
    expect(formatKetik('1,0', 2)).toBe('1,0')
    expect(formatKetik('1,00', 2)).toBe('1,00')
    expect(formatKetik('0', 2)).toBe('0')
  })

  // ⚠️ Dan pembandingnya: `formatNumber` memang menghapusnya. Dipaku
  // supaya tidak ada yang "menyederhanakan" `formatKetik` menjadi
  // panggilan `formatNumber`.
  it('⚠️ formatNumber MEMBUANGNYA — karena itu ia tidak dipakai saat mengetik', () => {
    expect(formatNumber('1.0', 2)).toBe('1')
    expect(formatNumber('1.00', 2)).toBe('1')
  })

  it('⚠️ ekor DIPOTONG, bukan dibulatkan', () => {
    expect(formatKetik('12,999', 2)).toBe('12,99')
    expect(formatKetik('12,991', 2)).toBe('12,99')
  })

  // ⭐ Papan tik angka banyak yang hanya punya titik — titik yang DIKETIK
  // diterjemahkan `normalisasiKetikan`, bukan ditebak `formatKetik`.
  //
  // ⛔ `formatKetik` SENDIRI membaca titik sebagai pemisah ribuan. Itu
  // disengaja: tanpa aturan tunggal, keluaran kotak yang diumpankan kembali
  // (`1.000` + `5` = `1.0005`) tidak dapat dibedakan dari desimal.
  it('⭐ titik yang DIKETIK menjadi koma lewat normalisasiKetikan', () => {
    expect(formatKetik(normalisasiKetikan('12.', '12'), 2)).toBe('12,')
    expect(formatKetik(normalisasiKetikan('1.000.', '1.000'), 2)).toBe('1.000,')
  })

  it('⛔ titik yang BUKAN baru diketik tetap pemisah ribuan', () => {
    // Panjangnya bertambah dua: bukan sisipan satu huruf, jadi apa adanya.
    expect(normalisasiKetikan('1.00', '1')).toBe('1.00')
    expect(formatKetik('12.99', 2)).toBe('1.299')
  })

  it('⭐ minus dipertahankan', () => {
    expect(formatKetik('-1000', 2)).toBe('-1.000')
    expect(formatKetik('-1000,25', 2)).toBe('-1.000,25')
  })

  it('⛔ huruf dan tanda asing dibuang', () => {
    expect(formatKetik('1a0b0c0', 2)).toBe('1.000')
    expect(formatKetik('', 2)).toBe('')
  })

  // ⚠️ Yang mengetik `0,5` MELEWATI keadaan `0`; membuang nol di depan
  // tanpa syarat akan memakan angkanya.
  it('⚠️ nol di depan dibuang, tetapi `0` sendirian bertahan', () => {
    expect(formatKetik('007', 2)).toBe('7')
    expect(formatKetik('0', 2)).toBe('0')
    expect(formatKetik('0,5', 2)).toBe('0,5')
    expect(formatKetik(',5', 2)).toBe('0,5')
  })

  // ⛔ CACAT YANG DITANGKAP UJI INI, 7 Oktober 2026: bentuk pertama
  // MENYAMBUNG digit sesudah pemisah, sehingga `1000,5` menjadi `10.005`
  // — keliru SEPULUH KALI LIPAT, tanpa satu pun tanda di layar. Kolom
  // bulat membuang ekornya.
  it('⛔ desimal 0 => ekor DIBUANG, bukan disambung', () => {
    expect(formatKetik('1000,5', 0)).toBe('1.000')
    expect(formatKetik('1000', 0)).toBe('1.000')
  })
})

// ---------------------------------------------------------------------
// ⛔ CACAT YANG DILAPORKAN PEMILIK PROSES 7 Oktober 2026
// ---------------------------------------------------------------------
// *"kenapa malah mentok begini tidak bisa meng input lebih"* — kotak
// `Amount` mentok di `1,00` dan `Amount in IDR` di `1.000`.
//
// Sebabnya: kotak terkendali mengumpankan KELUARANNYA SENDIRI kembali.
// Sesudah mengetik empat digit layar berbunyi `1.000`; digit kelima membuat
// peramban mengirim `1.0005`. Bentuk pertama `formatKetik` membaca titik
// PEMISAH RIBUAN itu sebagai pemisah DESIMAL, lalu memotong ekornya ke dua
// digit — kembali ke `1,00`. Setiap ketukan berikutnya mengulang hal yang
// sama, jadi angkanya tidak pernah tumbuh.
//
// ⭐ ATURANNYA KINI TUNGGAL DAN TIDAK AMBIGU: titik SELALU pemisah ribuan,
// koma SELALU pemisah desimal. Titik yang DIKETIK diterjemahkan menjadi koma
// di batas masukan (`normalisasiKetikan`), bukan ditebak di sini.
describe('formatKetik — keluaran sendiri diumpankan kembali', () => {
  it('⛔ digit kelima sesudah `1.000` TIDAK mengecilkan angkanya', () => {
    expect(formatKetik('1.0005', 2)).toBe('10.005')
    expect(formatKetik('1.0005', 0)).toBe('10.005')
  })

  it('⛔ mengetik terus dari nol sampai jutaan — nol langkah mundur', () => {
    let tampil = ''
    for (const digit of '1234567') {
      tampil = formatKetik(tampil + digit, 2)
    }
    expect(tampil).toBe('1.234.567')
  })

  it('⛔ dan dengan desimal: ekor bertahan sambil bagian bulat tumbuh', () => {
    let tampil = formatKetik('1234,5', 2)
    expect(tampil).toBe('1.234,5')
    tampil = formatKetik(tampil + '6', 2)
    expect(tampil).toBe('1.234,56')
    // Digit ketujuh di ekor DIPOTONG, bukan memindahkan koma.
    tampil = formatKetik(tampil + '7', 2)
    expect(tampil).toBe('1.234,56')
  })

  it('⭐ titik SELALU ribuan — nol tebakan', () => {
    expect(formatKetik('1.000', 2)).toBe('1.000')
    expect(formatKetik('12.345.678', 2)).toBe('12.345.678')
  })
})

describe('keKabelAngka — kembali ke bentuk backend', () => {
  it('⭐ titik desimal, nol pemisah ribuan', () => {
    expect(keKabelAngka('1.000.000,5')).toBe('1000000.5')
    expect(keKabelAngka('1.000')).toBe('1000')
    expect(keKabelAngka('12,99')).toBe('12.99')
  })

  // ⚠️ Koma tanpa digit BUKAN angka sah di kabel — layar tetap
  // memperlihatkan komanya, yang dikirim `12`.
  it('⚠️ koma menggantung dikirim tanpa koma', () => {
    expect(keKabelAngka('12,')).toBe('12')
  })

  it('⭐ minus dan kosong', () => {
    expect(keKabelAngka('-1.000,25')).toBe('-1000.25')
    expect(keKabelAngka('')).toBe('')
  })

  // ⛔ KEDUANYA WAJIB TETAP SEPASANG.
  it('⛔ bolak-balik tidak mengubah angka', () => {
    // ⛔ Bentuk KABEL memakai TITIK desimal; ia diterjemahkan ke koma
    // sebelum masuk — persis yang `FieldAngka` lakukan saat kotak dipegang.
    for (const kabel of ['1000000.5', '12.99', '0.5', '-1000.25', '7', '1000']) {
      expect(keKabelAngka(formatKetik(kabel.replace('.', ','), 2))).toBe(kabel)
    }
  })
})

// ---------------------------------------------------------------------
// GELUNG KOTAK TERKENDALI — tiruan persis `FieldAngka`
// ---------------------------------------------------------------------
// ⛔ Uji fungsi satuan di atas TIDAK cukup, dan cacat 7 Oktober 2026
// membuktikannya: tiap fungsinya benar sendiri-sendiri, yang salah adalah
// KELUARANNYA DIUMPANKAN KEMBALI. Yang ditiru di bawah adalah gelung
// lengkapnya — persis urutan yang `FieldAngka` jalankan tiap ketukan.
describe('gelung kotak terkendali', () => {
  /** Satu ketukan: sisipkan `huruf` di ujung, lalu olah seperti komponen. */
  function ketuk(tampil: string, huruf: string, desimal: number) {
    const mentah = tampil + huruf
    const t = formatKetik(normalisasiKetikan(mentah, tampil), desimal)
    return { tampil: t, kabel: keKabelAngka(t) }
  }

  it('⛔ mengetik 1234567 — angkanya TUMBUH tiap ketukan', () => {
    let s = { tampil: '', kabel: '' }
    const jejak: string[] = []
    for (const h of '1234567') {
      s = ketuk(s.tampil, h, 2)
      jejak.push(s.tampil)
    }
    expect(jejak).toEqual(['1', '12', '123', '1.234', '12.345', '123.456', '1.234.567'])
    expect(s.kabel).toBe('1234567')
  })

  it('⛔ kolom BULAT juga tumbuh — `Amount in IDR` mentok di `1.000`', () => {
    let s = { tampil: '', kabel: '' }
    for (const h of '1000000') s = ketuk(s.tampil, h, 0)
    expect(s.tampil).toBe('1.000.000')
    expect(s.kabel).toBe('1000000')
  })

  it('⭐ koma lalu dua desimal', () => {
    let s = { tampil: '', kabel: '' }
    for (const h of '1234') s = ketuk(s.tampil, h, 2)
    s = ketuk(s.tampil, ',', 2)
    expect(s.tampil).toBe('1.234,')
    s = ketuk(s.tampil, '5', 2)
    s = ketuk(s.tampil, '6', 2)
    expect(s.tampil).toBe('1.234,56')
    expect(s.kabel).toBe('1234.56')
  })

  it('⭐ titik papan tik angka menjadi koma', () => {
    let s = { tampil: '', kabel: '' }
    for (const h of '12') s = ketuk(s.tampil, h, 2)
    s = ketuk(s.tampil, '.', 2)
    expect(s.tampil).toBe('12,')
    s = ketuk(s.tampil, '9', 2)
    s = ketuk(s.tampil, '9', 2)
    expect(s.tampil).toBe('12,99')
    expect(s.kabel).toBe('12.99')
  })

  // ⚠️ Nilai yang DIMUAT dari backend juga harus dapat dilanjutkan.
  it('⚠️ melanjutkan nilai muatan backend', () => {
    const kabel = '1000000.5'
    let tampil = formatKetik(kabel.replace('.', ','), 2) // onFocus
    expect(tampil).toBe('1.000.000,5')
    const s = ketuk(tampil, '2', 2)
    expect(s.tampil).toBe('1.000.000,52')
    expect(s.kabel).toBe('1000000.52')
  })
})
