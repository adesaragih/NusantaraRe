import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { BERANDA, KETERANGAN_BELUM_DIMIGRASI } from '../inti/frontend/labels'
import type { DaftarBeranda } from '../inti/frontend/modul'
import { batangKotakMasuk, jenisKotakMasuk, modulBerkotakMasuk, saringDaftar, type KotakModul } from './Beranda'
import { MODUL_FRONTEND } from './daftar'

// Uji Beranda — butir bg.

describe('label Beranda', () => {
  it('kata "belum dimigrasi" satu tempat', () => {
    expect(KETERANGAN_BELUM_DIMIGRASI).toBe('belum dimigrasi')
  })
})

// Kotak masuk Beranda mengikuti workbasket (keputusan work owner 06-10-2026): modul yang mendaftarkan
// `antreanBeranda` di PENDAFTARAN_MENU-nya - kini NB Treaty In dan EDM Treaty In (keputusan work owner 07-10-2026)
// - dan aktif.
describe('kotak masuk Beranda', () => {
  it('modul berpenghitung: NB Treaty In dan EDM Treaty In, dan hilang bila modulnya nonaktif', () => {
    expect(modulBerkotakMasuk(MODUL_FRONTEND).map((m) => m.nama)).toEqual(['edmtreatyin', 'nbtreatyin'])
    expect(modulBerkotakMasuk(MODUL_FRONTEND, ['claimlife'])).toEqual([])
  })

})

// Panel kotak masuk (keputusan work owner 06-10-2026: seperti dasbor Pega "Note | Count" + grafik per workbasket,
// "yang lebih modern, ga harus pie chart"): batang per workbasket urut terbanyak, lalu cacah per jenis (modul).
describe('panel kotak masuk', () => {
  const m = (nama: string, kelompok: string) => ({ nama, kelompok }) as unknown as KotakModul['modul']
  const kotak: KotakModul[] = [
    { modul: m('alfa', 'UJI Alfa'), baris: [{ workbasket: 'A1', nama: 'UJI A1', jumlah: 1 }, { workbasket: 'A2', nama: 'UJI A2', jumlah: 0 }] },
    { modul: m('beta', 'UJI Beta'), baris: [{ workbasket: 'B1', nama: 'UJI B1', jumlah: 4 }] },
  ]

  it('batang: urut terbanyak, lebar relatif ke yang terbanyak, porsi dari total, warna berurutan; NOL tidak tampil', () => {
    // "kalo jumlahnya 0 ga usah muncul" (work owner 06-10-2026)
    expect(batangKotakMasuk(kotak).map((b) => [b.nama, b.jumlah, b.persen, b.porsi, b.warna])).toEqual([
      ['UJI B1', 4, 100, 80, 0],
      ['UJI A1', 1, 25, 20, 1],
    ])
    expect(batangKotakMasuk([{ modul: m('alfa', 'UJI Alfa'), baris: [{ workbasket: 'A', nama: 'UJI A', jumlah: 0 }] }])).toEqual([])
  })

  it('warna berputar setelah palet habis', () => {
    const banyak: KotakModul[] = [
      { modul: m('alfa', 'UJI Alfa'), baris: Array.from({ length: 8 }, (_, i) => ({ workbasket: `W${i}`, nama: `UJI W${i}`, jumlah: 8 - i })) },
    ]
    expect(batangKotakMasuk(banyak).map((b) => b.warna)).toEqual([0, 1, 2, 3, 4, 5, 0, 1])
  })

  // Klik workbasket -> per jenis -> daftar berkas (permintaan work owner 06-10-2026).
  const bersama: KotakModul[] = [
    { modul: m('alfa', 'UJI Alfa'), baris: [{ workbasket: 'W', nama: 'UJI W', jumlah: 2 }, { workbasket: 'X', nama: 'UJI X', jumlah: 1 }] },
    { modul: m('beta', 'UJI Beta'), baris: [{ workbasket: 'W', nama: 'UJI W', jumlah: 3 }] },
    { modul: m('gama', 'UJI Gama'), baris: [{ workbasket: 'W', nama: 'UJI W', jumlah: 0 }] },
  ]

  it('workbasket yang sama di beberapa modul = SATU batang berjumlah gabungan', () => {
    expect(batangKotakMasuk(bersama).map((b) => [b.workbasket, b.jumlah, b.persen, b.porsi])).toEqual([
      ['W', 5, 100, 83],
      ['X', 1, 20, 17],
    ])
  })

  it('per jenis: hanya modul di workbasket yang dibuka, urut terbanyak; modul bernilai nol tidak tampil', () => {
    expect(jenisKotakMasuk(bersama, 'W').map((j) => [j.modul.kelompok, j.jumlah])).toEqual([
      ['UJI Beta', 3],
      ['UJI Alfa', 2],
    ])
    expect(jenisKotakMasuk(bersama, 'X').map((j) => [j.modul.kelompok, j.jumlah])).toEqual([['UJI Alfa', 1]])
    expect(jenisKotakMasuk(bersama, 'TIDAK-ADA')).toEqual([])
  })

  it('per jenis tertutup sampai batang workbasket diklik; klik jenis membuka daftar berkas workbasket itu', () => {
    const sumber = readFileSync(join(__dirname, 'Beranda.tsx'), 'utf8')
    expect(sumber).toContain('{wbBuka !== undefined && (')
    expect(sumber).toContain('aria-expanded={wbBuka?.workbasket === b.workbasket}')
    expect(sumber).toContain('{ modul: j.modul, workbasket: wbBuka.workbasket, judul: wbBuka.nama }')
    // klik batang tidak lagi langsung membuka daftar berkas
    expect(sumber).not.toContain('pilih({ modul: b.modul')
  })
})

describe('Beranda berbahasa Inggris (06-10-2026)', () => {
  it('label kotak masuk dan nama menu', () => {
    expect(BERANDA.judul).toBe('Home')
    expect(BERANDA.kotakMasuk).toBe('Your inbox')
    expect(BERANDA.perWorkbasket).toBe('By workbasket')
    expect(BERANDA.perJenis).toBe('By case type')
  })

  it('tanpa teks Indonesia tertanam di Beranda.tsx', () => {
    const sumber = readFileSync(join(__dirname, 'Beranda.tsx'), 'utf8')
    for (const teks of ['gagal dimuat', 'KETERANGAN_BELUM_DIMIGRASI']) {
      expect(sumber, teks).not.toContain(teks)
    }
  })
})

// Beranda = salam + kotak masuk saja (perintah work owner 06-10-2026: kartu tahap Claim Life, tabel Modul dan
// Peran Anda "buang aja, ga perlu").
describe('Beranda tanpa kartu Claim Life, tabel Modul, Peran Anda', () => {
  it('ketiga bagian itu tidak dirender lagi', () => {
    const sumber = readFileSync(join(__dirname, 'Beranda.tsx'), 'utf8')
    for (const teks of ['beranda__antrean', 'beranda__kisi', 'beranda-modul', 'beranda-peran', 'claimlife', 'kartuModul']) {
      expect(sumber, teks).not.toContain(teks)
    }
  })
})

// Filter daftar berkas kotak masuk (perintah work owner 06-10-2026: "tambahkan filter disini").
describe('saringDaftar', () => {
  const daftar: DaftarBeranda = {
    kolom: [
      { kunci: 'id', label: 'Offer No' },
      { kunci: 'nama', label: 'Insured Name' },
    ],
    baris: [
      { id: 'UJI-1', sel: { id: 'UJI-1', nama: 'UJI ALFA', tersembunyi: 'ZETA' } },
      { id: 'UJI-2', sel: { id: 'UJI-2', nama: 'UJI BETA' } },
    ],
  }
  const id = (q: string) => saringDaftar(daftar, q).map((r) => r.id)

  it('kueri kosong atau spasi = semua baris', () => {
    expect(id('')).toEqual(['UJI-1', 'UJI-2'])
    expect(id('   ')).toEqual(['UJI-1', 'UJI-2'])
  })

  it('cocok di kolom mana pun, huruf besar-kecil diabaikan, spasi tepi dibuang', () => {
    expect(id('beta')).toEqual(['UJI-2'])
    expect(id(' uji-1 ')).toEqual(['UJI-1'])
    expect(id('uji')).toEqual(['UJI-1', 'UJI-2'])
  })

  it('hanya kolom yang tampil yang dicari; tak cocok = kosong', () => {
    expect(id('zeta')).toEqual([])
    expect(id('gama')).toEqual([])
  })

  it('kapsul filter di kepala daftar, dikosongkan saat daftar lain dibuka', () => {
    const sumber = readFileSync(join(__dirname, 'Beranda.tsx'), 'utf8')
    expect(sumber).toContain('role="search"')
    expect(sumber).toContain('saringDaftar(daftar, kueri)')
    expect(sumber).toContain("setKueri('')")
    expect(BERANDA.saring).toBe('Filter cases')
    expect(BERANDA.tidakCocok).toBe('No cases match the filter.')
  })
})

// Back dari layar kasus: workbasket, jenis, daftar berkas, dan filter tetap seperti ditinggalkan (perintah work owner
// 06-10-2026: "saat di back, ini jangan ilang"). Beranda dibongkar selama layar kasus tampil, jadi App yang mengingat.
describe('pilihan panel diingat App', () => {
  const beranda = readFileSync(join(__dirname, 'Beranda.tsx'), 'utf8')
  const app = readFileSync(join(__dirname, 'App.tsx'), 'utf8')

  it('Beranda mulai dari ingatan dan melaporkan setiap perubahan', () => {
    expect(beranda).toContain('useState<PilihanKotak | null>(ingatan?.pilihan ?? null)')
    expect(beranda).toContain('useState<string | null>(ingatan?.bukaWb ?? null)')
    expect(beranda).toContain("useState(ingatan?.kueri ?? '')")
    expect(beranda).toContain('onIngat?.({ bukaWb, pilihan, kueri })')
  })

  it('memuat daftar sesudah Back tidak mengosongkan filter; filter dikosongkan hanya oleh klik pemakai', () => {
    const efek = beranda.slice(beranda.indexOf('const muat = pilihan?.modul.daftarBeranda') - 200, beranda.indexOf('}, [pilihan])'))
    expect(efek).not.toContain('setKueri')
  })

  it('App menyimpan ingatan per akun', () => {
    expect(app).toContain('ingatan={ingatanBeranda?.akun === masuk.akunID ? ingatanBeranda.isi : undefined}')
    expect(app).toContain('onIngat={ingatBeranda}')
  })
})
