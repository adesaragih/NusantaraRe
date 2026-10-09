// Penjaga layar modul Treaty In Adjustment.
//
// ⛔ Membaca BERKAS SUMBER, bukan merender — pola yang sama dengan
// `modul/treatyin/frontend/layar.test.ts`. Yang dijaga teks dan susunannya,
// dan merender menambah ketergantungan tanpa menambah yang dijaga.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KOLOM_BERKAS_LAMPIRAN, KOLOM_HISTORY, KOLOM_LAMPIRAN, LAMPIRAN } from './labels'
import { LEBAR_DAFTAR, LEBAR_TOMBOL_BARIS } from './labelsPenyesuaian'
import { HALAMAN_TREATYINADJUSTMENT } from './menu'

const AKAR = __dirname
const LAYAR = readFileSync(join(AKAR, 'pages', 'LampiranKontrak.tsx'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const CSS = readFileSync(join(AKAR, 'treatyinadjustment.css'), 'utf8')
const REPO_RIWAYAT = readFileSync(
  join(AKAR, '..', 'backend', 'repository', 'warisan_riwayat.go'),
  'utf8',
)

describe('panel Attachment — jalan menuju ke sana ADA', () => {
  it('⛔ halamannya terdaftar, dirutekan, dan punya berkasnya', () => {
    // Ronde sebelumnya berhenti di repository: nol services, nol handlers,
    // nol layar — "tabel terisi yang tidak dibaca siapa pun adalah pekerjaan
    // yang terlihat selesai dan tidak sampai ke pemakai."
    expect([...HALAMAN_TREATYINADJUSTMENT]).toContain('treatyinadjustment-lampiran')
    expect(RUTE).toContain("halaman === 'treatyinadjustment-lampiran'")
    expect(RUTE).toContain('LampiranKontrak')
    expect(LAYAR.length).toBeGreaterThan(500)
  })

  it('keempat kolomnya disalin dari WorkAttachments.xml apa adanya', () => {
    expect([...KOLOM_LAMPIRAN]).toEqual(['Category', 'Count', 'Upload file', 'View File'])
    expect([...KOLOM_BERKAS_LAMPIRAN]).toEqual(['File Name', 'Type', 'Uploaded', 'By'])
    expect([...KOLOM_HISTORY]).toEqual(['Date', 'PIC', 'Approval', 'Comment'])
  })

  it('spanduk biru disalin apa adanya — ia ATURAN, bukan hiasan', () => {
    expect(LAMPIRAN.spanduk).toBe('Recommended safe substitute should be . or _')
    expect(LAYAR).toContain('tria__spanduk')
    expect(CSS).toMatch(/\.tria__spanduk \{/)
  })

  // ⭐ RALAT 8 Oktober 2026 — penamaan kategori diperbaiki (tangkapan layar
  // Pega): nama dari `M_KATEGORIMASTERTREATY`, disusun backend. Layar TIDAK
  // menghafal nama apa pun dan tidak lagi menandai "belum dipastikan".
  it('⭐ kategori tampil dengan NAMA dari backend; kode hanya bila tanpa nama', () => {
    expect(LAYAR).toContain("<td>{k.nama !== '' ? k.nama : k.kode}</td>")
    expect(LAYAR).not.toContain('kategoriBelumPasti')
    expect(LAYAR).not.toContain('namaBelumBerumah')
    expect(Object.keys(LAMPIRAN)).not.toContain('namaBelumBerumah')
    for (const n of ['Binding, signed share Email', 'Claim Data', 'Info Pack']) {
      expect(LAYAR).not.toContain(`>${n}<`)
    }
  })

  it('⭐ jenis kontrak diteruskan — nama Non-Prop untuk kode 00007', () => {
    expect(LAYAR).toContain('ambilLampiran(masterID, jenis)')
    const ADJ = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
    expect(ADJ).toContain('jenis={cabang}')
  })

  it('⛔ Count TIDAK diformat — ia cacah butir, bukan uang', () => {
    expect(LAYAR).toContain('String(k.cacah)')
    expect(LAYAR).not.toMatch(/selAngka|formatNumber/)
  })

  it('Download All MATI, Refresh HIDUP — dan bedanya beralasan', () => {
    // `Download All` menarik berkas dari `T_STORAGE_IMAGE`, jalur yang belum
    // dibangun. `Refresh` hanya membaca ulang, dan itu boleh.
    const i = LAYAR.indexOf('tria__aksi')
    expect(i).toBeGreaterThan(0)
    const blok = LAYAR.slice(i, i + 900)
    expect(blok).toContain('unduhSemua')
    expect(blok).toContain('segarkan')
    expect((blok.match(/disabled/g) ?? []).length).toBe(1)
    expect(blok).toContain('setMuatUlang')
  })

  it('kosong memakai Kosong (belum ada DATA), teksnya No items', () => {
    expect(LAMPIRAN.tanpaIsi).toBe('No items')
    expect(LAYAR).toContain('<Kosong')
    expect(LAYAR).not.toContain('tria__belum')
  })

  it('⛔ nol tabel baru disebut untuk lampiran', () => {
    for (const salah of ['M_TREATYIN_LAMPIRAN', 'M_TREATYINADJ_LAMPIRAN']) {
      expect(LAYAR).not.toContain(salah)
    }
  })

  it('⭐ History KINI TERISI — dibaca modul ini sendiri, nol impor silang', () => {
    // Ronde sebelumnya meninggalkannya kosong dengan alasan "sumbernya
    // milik modul sebelah". Yang ditolak penjaga adalah impor paket Go,
    // bukan pembacaan tabel.
    expect(LAYAR).toContain('ambilRiwayat')
    expect(LAYAR).toContain('riwayat.map')
    expect(LAYAR).toContain('T_VIEW_COMMENT')
    // ⛔ Dan ia benar-benar membaca sendiri: nol impor dari modul sebelah.
    expect(LAYAR).not.toMatch(/from '.*modul\/treatyin\//)
    expect(REPO_RIWAYAT).not.toMatch(/nusantarare\/modul\/treatyin\//)
  })

  it('⛔ pembacaan History TERPISAH dari lampiran', () => {
    // Keduanya dari tabel berbeda; kegagalan satu tidak boleh
    // mengosongkan yang lain.
    expect((LAYAR.match(/useEffect\(/g) ?? []).length).toBeGreaterThanOrEqual(2)
  })

  it('riwayat dibaca ORDER BY URUTAN — larik Pega berurut', () => {
    expect(REPO_RIWAYAT).toContain('ORDER BY URUTAN')
  })
})

describe('§2.2 — satu pemilih kontrak, bukan dua daftar', () => {
  it('RantaiVersi melaporkan NOMOR WARISAN, bukan id model baru', () => {
    const RV = readFileSync(join(AKAR, 'pages', 'RantaiVersi.tsx'), 'utf8')
    expect(RV).toContain('onPilihWarisan')
    // ⛔ `TREATYID` dan `MASTERID` keduanya pengenal sistem lama bertipe
    // teks; mengirim `id` model baru mengembalikan nol baris tanpa galat.
    expect(RV).toContain('nomorKontrakWarisan')
  })

  it('pengenalnya hidup di RUTE — dua layar memakainya, satu memilihnya', () => {
    expect(RUTE).toContain('const [masterID, setMasterID]')
    expect(RUTE).toContain('masterID={masterID}')
    expect(RUTE).toContain('onPilihWarisan')
  })

  it('⛔ NOL layar daftar kontrak warisan kedua di modul ini', () => {
    const { readdirSync } = require('node:fs') as typeof import('node:fs')
    const halaman = readdirSync(join(AKAR, 'pages'))
    // ⭐ `PenyesuaianKontrak.tsx` (5 Oktober 2026) BUKAN daftar kontrak
    // kedua: ia mendaftar PENYESUAIAN (`TREATY_IN_EDM`), layar sistem lama
    // `InputTreatyInAdjustment` sendiri. Daftar kontrak tetap satu.
    expect(halaman.sort()).toEqual(['LampiranKontrak.tsx', 'PenyesuaianKontrak.tsx', 'RantaiVersi.tsx'])
    const P = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
    expect(P).not.toMatch(/ambilKontrak|TREATY_IN(?!_EDM)/)
  })
})


// ⛔ KEPALA KOLOM TETAP TERLIHAT SAAT TABEL DIGULIR — permintaan pemilik
// proses 7 Oktober 2026, sepasang dengan penjaga di modul Treaty In.
describe('kepala kolom tidak ikut tenggelam', () => {
  it('sticky DIBATASI `.table-wrap`, berlatar penuh, bergaris `box-shadow`', () => {
    const i = CSS.indexOf('.treatyinadjustment .table-wrap .tria__tabel thead th')
    expect(i).toBeGreaterThan(0)
    const aturan = CSS.slice(i, i + 260)
    expect(aturan).toMatch(/position: sticky/)
    expect(aturan).toMatch(/background: var\(--bg\)/)
    // ⚠️ WAJIB di sini: `.tria__tabel` ber-`border-collapse: collapse`,
    // dan `border-bottom` pada mode itu milik kisi tabel — ia tergulir pergi.
    expect(aturan).toMatch(/box-shadow: inset 0 -1px 0 var\(--border\)/)
    expect(CSS).toMatch(/\.tria__tabel \{[^}]*border-collapse: collapse/)
  })
})

// ⛔ GRID DAFTAR TIDAK BOLEH MEMEPET — keluhan pemakai 8 Oktober 2026:
// *"perbaiki design tablenya biar gak mepet"*.
//
// Yang terlihat di layar: `SAHABAT INSURANCE` terpecah menjadi
// `SAHAB AT INSUR ANCE`, dan `08-10-2026` pecah dua baris.
//
// ⚠️ Sebabnya BUKAN padding, melainkan LEBAR: kolomnya berukuran persen dan
// persen menyusut mengikuti layar tanpa batas bawah.
describe('grid daftar penyesuaian — lebar dan pemenggalan', () => {
  const CSS2 = readFileSync(join(AKAR, 'treatyinadjustment.css'), 'utf8')
  const HAL = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')

  it('⛔ sel tabel memutus di SELA KATA, bukan di tengahnya', () => {
    // ⚠️ `anywhere` TIDAK dilarang menyeluruh: panel teks panjang
    // (`tria__teks`) memang membutuhkannya untuk kalimat tanpa spasi.
    // Yang dijaga hanya aturan SEL TABEL.
    const i = CSS2.indexOf('.treatyinadjustment .tria__tabel th,')
    expect(i).toBeGreaterThan(0)
    const aturan = CSS2.slice(i, CSS2.indexOf('}', i))
    expect(aturan).toContain('overflow-wrap: break-word')
    expect(aturan).not.toContain('overflow-wrap: anywhere')
  })

  // ⛔ RALAT 8 Oktober 2026 — *"perbaiki tampilan adjustment treaty"*: lebar
  // minimum 1.551px + `table-layout: fixed` tetap memenggal `Complete`,
  // `ADESAMUEL`, `NonProportional` di kolom sempit ekspor dan mendorong kolom
  // pertama keluar layar. Lebar ekspor kini PERBANDINGAN acuan di `<col>`;
  // tata letak `auto` menjamin tiap kolom minimal selebar kata terpanjangnya.
  it('⭐ lebar ekspor tetap acuan `<col>`; tabel `auto` selebar layar — kata nol terpenggal', () => {
    // 136+134+93+93+141+106+137+135+160+149+83+72 = 1.439, + 2 x 56 = 1.551.
    const jumlah = LEBAR_DAFTAR.reduce((a, b) => a + b, 0) + 2 * LEBAR_TOMBOL_BARIS
    expect(jumlah).toBe(1551)
    const i = CSS2.indexOf('.treatyinadjustment .tria__tabel--daftar {')
    expect(i).toBeGreaterThan(0)
    const aturan = CSS2.slice(i, CSS2.indexOf('}', i))
    expect(aturan).toContain('table-layout: auto')
    expect(aturan).not.toContain('table-layout: fixed')
    expect(aturan).not.toContain('min-width')
    expect(HAL).toContain('tria__tabel tria__tabel--daftar')
  })

  it('⛔ nowrap tanggal ada di SEL, bukan di `<col>` — `<col>` mengabaikannya', () => {
    const j = CSS2.indexOf('.treatyinadjustment .tria__sel-tanggal {')
    expect(j).toBeGreaterThan(0)
    expect(CSS2.slice(j, CSS2.indexOf('}', j))).toContain('white-space: nowrap')
    expect(CSS2).not.toContain('.tria__kol-tanggal')
    expect(HAL).toContain("JENIS_DAFTAR[i] === 'tanggal' ? 'tria__sel-tanggal' : undefined")
    // Dan `<col>` hanya membawa lebarnya.
    expect(HAL).toContain('<col key={i} style={{ width: persenLebar(lebar, i) }} />')
  })
})
