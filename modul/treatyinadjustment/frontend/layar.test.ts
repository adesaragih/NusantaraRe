// Penjaga layar modul Treaty In Adjustment.
//
// ⛔ Membaca BERKAS SUMBER, bukan merender — pola yang sama dengan
// `modul/treatyin/frontend/layar.test.ts`. Yang dijaga teks dan susunannya,
// dan merender menambah ketergantungan tanpa menambah yang dijaga.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KOLOM_BERKAS_LAMPIRAN, KOLOM_HISTORY, KOLOM_LAMPIRAN, LAMPIRAN } from './labels'
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

  it('⛔ kategori belum dipastikan menampilkan KODE, bukan nama tebakan', () => {
    expect(LAYAR).toContain('k.dipastikan')
    expect(LAYAR).toContain('kategoriBelumPasti')
    for (const n of ['Binding, signed share Email', 'Claim Data', 'Info Pack']) {
      expect(LAYAR).not.toContain(`>${n}<`)
    }
    expect(LAMPIRAN.namaBelumBerumah).toContain('Letter of Acknowledgment / LOA')
    expect(LAMPIRAN.namaBelumBerumah).toContain('PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md')
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
