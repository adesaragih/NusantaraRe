import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { UNGGAH_CSV } from '../labels'
import { kunciPenolakan, ringkasanTinjau } from './UnggahCSVPeserta'

// Uji layar unggah CSV peserta — tiket 04.

const BERKAS = readFileSync(join(__dirname, 'UnggahCSVPeserta.tsx'), 'utf8')

/** Sumber TANPA komentar — prosa yang menerangkan larangan bukan pelanggaran. */
const SUMBER = BERKAS.split('\n')
  .filter((b) => {
    const t = b.trimStart()
    return !t.startsWith('//') && !t.startsWith('*')
  })
  .join('\n')

describe('dua langkah', () => {
  it('tinjau dan simpan adalah dua tombol, bukan satu', () => {
    expect(SUMBER).toContain('tinjauUnggahPolis')
    expect(SUMBER).toContain('simpanUnggahPolis')
    // ⛔ Labelnya DIRUJUK dari labels.premiumlist, tidak diketik ulang di
    // layar — dua salinan satu label berarti yang satu akan tertinggal.
    expect(SUMBER).toContain('UNGGAH_CSV.tinjau')
    expect(SUMBER).toContain('UNGGAH_CSV.simpan')
    expect(SUMBER).not.toContain(`'${UNGGAH_CSV.simpan}'`)
  })

  it('tombol simpan mati selama masih ada penolakan', () => {
    // ⛔ Tombol yang hidup lalu selalu dijawab 409 mengajari orang
    // mengabaikan galat.
    expect(SUMBER).toContain('tinjau.lolos && !sibuk')
    expect(SUMBER).toContain('disabled={!bolehSimpan}')
    // Batas produk kini penolakan Validate CSV (05-10-2026): tidak ada lagi
    // peringatan sesudah Calculate CSV; tombolnya mati bila tinjauan tidak lolos.
    expect(SUMBER).toContain('const bolehSimpan = tinjau !== null && tinjau.lolos && !sibuk')
    expect(SUMBER).not.toContain('peringatan')
  })

  it('sebab tombol simpan mati DIKATAKAN', () => {
    expect(SUMBER).toContain('UNGGAH_CSV.perbaikiDulu')
    expect(UNGGAH_CSV.perbaikiDulu).toContain('no row is saved')
  })

  it('tinjauan lama dibuang saat berkas berganti', () => {
    // ⛔ Membiarkannya membuat tombol simpan merujuk berkas yang berbeda dari
    // yang terlihat di pemilih berkas.
    expect(SUMBER).toContain('function pilihBerkas')
    expect(SUMBER).toContain('setTinjau(null)')
  })
})

describe('aturan pemisah dinyatakan di muka', () => {
  it('layar menyebut titik, dan menolak kedua bentuk pemisah ribuan', () => {
    // ⛔ Korpus menyebut aturan ini enam kali di nama langkahnya dan TIDAK
    // PERNAH di pesan yang dilihat pemakai — sehingga orang baru tahu
    // aturannya setelah berkasnya ditolak.
    expect(UNGGAH_CSV.aturanPemisah).toContain('DOT')
    expect(UNGGAH_CSV.aturanPemisah).toContain('1,234,567.89')
    expect(UNGGAH_CSV.aturanPemisah).toContain('1.234.567,89')
    expect(UNGGAH_CSV.aturanPemisah).toContain('dd/mm/yyyy')
    expect(SUMBER).toContain('UNGGAH_CSV.aturanPemisah')
  })
})

describe('tabel penolakan', () => {
  it('menampilkan nomor baris DAN nama kolom', () => {
    expect(SUMBER).toContain('UNGGAH_CSV.kolomBaris')
    expect(SUMBER).toContain('UNGGAH_CSV.kolomKolom')
    expect(SUMBER).toContain('p.baris')
    expect(SUMBER).toContain('p.kolom')
  })

  it('menampilkan pesan VERBATIM dan sebab yang tepat, keduanya', () => {
    // ⚠️ Pesan verbatim menjaga pengenalan; sebab menjaga kebenaran.
    expect(SUMBER).toContain('p.pesan')
    expect(SUMBER).toContain('p.sebab')
  })

  it('pesan penolakan TIDAK diketik ulang di layar', () => {
    // ⛔ Ia datang dari server, VERBATIM dari korpus. Dua sumber untuk satu
    // kalimat berarti yang satu akan tertinggal.
    expect(SUMBER).not.toContain('HARUS ADA')
    expect(SUMBER).not.toContain('TIDAK BOLEH DOUBLE')
  })

  it('kunci baris unik walau satu baris punya banyak alasan', () => {
    const a = { baris: 3, kolom: 'PLAN', pesan: 'x', sebab: 'y' }
    const b = { baris: 3, kolom: 'PLAN', pesan: 'x', sebab: 'z' }
    expect(kunciPenolakan(a, 0)).not.toBe(kunciPenolakan(b, 1))
  })

  it('baris 0 berarti penolakan atas BERKASNYA, bukan atas sebuah baris', () => {
    expect(SUMBER).toContain('p.baris === 0')
  })
})

describe('ringkasan tinjauan', () => {
  it('berkas bersih dikatakan lolos', () => {
    expect(
      ringkasanTinjau({ cacahBaris: 5, cacahDitolak: 0, ditolak: [], lolos: true }),
    ).toContain('all passed')
  })

  it('berkas rusak menyebut cacah BARIS dan cacah ALASAN, keduanya', () => {
    // ⚠️ Satu baris dapat punya banyak alasan; menyebut salah satunya saja
    // membuat orang salah menaksir seberapa besar perbaikannya.
    const s = ringkasanTinjau({
      cacahBaris: 10,
      cacahDitolak: 3,
      ditolak: [
        { baris: 2, kolom: 'PLAN', pesan: 'a', sebab: 'b' },
        { baris: 2, kolom: 'DOB', pesan: 'a', sebab: 'b' },
        { baris: 7, kolom: 'PLAN', pesan: 'a', sebab: 'b' },
      ],
      lolos: false,
    })
    expect(s).toContain('10 rows read')
    expect(s).toContain('2 rows rejected')
    expect(s).toContain('3 reasons')
    expect(s).toContain('Nothing was saved')
  })
})

describe('panel unggah ringkas (03-10-2026)', () => {
  it('satu baris; aturan format dilipat, bukan dibuang', () => {
    expect(SUMBER).toContain('<details className="pl-unggah__aturan">')
    expect(SUMBER).toMatch(/<details className="pl-unggah__aturan">[\s\S]{0,200}UNGGAH_CSV\.aturanPemisah/)
    expect(SUMBER).toContain('className="panel__title pl-unggah__judul"')
  })
})
