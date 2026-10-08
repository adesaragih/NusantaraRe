// Aturan B pemilik proses, 6 Oktober 2026:
//
//   "apabila user melakukan inputan itu tidak langsung masuk ke table nya,
//    harus melakukan save atau submit dulu kemudian datanya masuk ke db"
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA PENJAGA INI DITULIS SEBELUM JALUR TULISNYA LAHIR
// ---------------------------------------------------------------------
// Hari ini aturan itu terpenuhi dengan sendirinya: nol rute tulis ada, dan
// setiap tab menyimpan suntingan di keadaan React setempat. Jadi penjaga ini
// tidak menangkap apa pun hari ini — dan justru itu gunanya.
//
// Aturan yang baru dijaga SESUDAH dilanggar menjaga hal yang sudah rusak.
// Kebocoran yang hendak dicegah bentuknya sepele: satu `onChange` yang
// memanggil penyimpan "supaya tidak hilang", ditulis dengan niat baik, dan
// sejak saat itu setiap ketukan papan ketik menyentuh basis data — tanpa
// seorang pun menekan Save.
//
// Keputusan lengkapnya di `docs/KEPUTUSAN-SASARAN-TULIS.md`.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Setiap berkas layar modul ini, berikut isinya. */
function berkasLayar(): { nama: string; isi: string }[] {
  const keluar: { nama: string; isi: string }[] = []
  for (const folder of ['', 'components', 'pages']) {
    const dir = folder === '' ? AKAR : join(AKAR, folder)
    for (const f of readdirSync(dir)) {
      if (!f.endsWith('.ts') && !f.endsWith('.tsx')) continue
      if (f.endsWith('.test.ts') || f.endsWith('.test.tsx')) continue
      keluar.push({ nama: folder === '' ? f : `${folder}/${f}`, isi: readFileSync(join(dir, f), 'utf8') })
    }
  }
  return keluar
}

/**
 * Nama fungsi `api.ts` yang MENULIS. Kosong hari ini.
 *
 * ⛔ Daftar ini WAJIB bertambah bersama fungsi tulis pertama. Fungsi tulis
 * yang tidak terdaftar membuat penjaga di bawah diam — dan penjaga yang diam
 * karena daftarnya basi lebih buruk daripada tidak ada penjaga.
 */
// ⭐ TERISI 7 Oktober 2026 — fungsi tulis PERTAMA lahir, dan penjaga ini
// menangkapnya persis seperti yang direncanakan ketika ia ditulis dalam
// keadaan kosong.
//
// ⛔ Sejak sekarang uji di bawah BENAR-BENAR memeriksa sesuatu: nol
// komponen tab boleh memanggil `simpanKontrak` dari `onChange`. Itulah
// kebocoran yang Aturan B cegah — satu `onChange` yang menyimpan "supaya
// tidak hilang", ditulis dengan niat baik, dan sejak itu tiap ketukan papan
// ketik menyentuh basis data tanpa seorang pun menekan Save.
// ⭐ DUA: `simpanKontrak` (Save) dan `kirimKontrak` (Submit) — tepat
// kedua tombol yang Aturan B sebut.
// ⭐ TIGA sejak 7 Oktober 2026: `mulaiRevisi` — tombol `Revision` daftar
// kontrak, yang di Pega pun menyimpan SEKETIKA (`SetTreatyIn_Act` [10–11]).
// ⭐ EMPAT sejak 8 Oktober 2026: `unggahLampiran` — tombol Attach panel
// Attachment (`TreatySaveAttachment`), yang di Pega pun menyimpan seketika.
const FUNGSI_TULIS: readonly string[] = ['simpanKontrak', 'kirimKontrak', 'mulaiRevisi', 'unggahLampiran']

describe('Aturan B — isian tidak masuk DB sampai Save atau Submit', () => {
  it('⛔ nol fungsi tulis dipanggil dari `onChange`', () => {
    const pelanggar: string[] = []
    for (const { nama, isi } of berkasLayar()) {
      for (const fn of FUNGSI_TULIS) {
        // Cari pemanggilan fungsi tulis di dalam badan sebuah `onChange`.
        const pola = new RegExp(`onChange=\\{[^}]*\\b${fn}\\s*\\(`, 's')
        if (pola.test(isi)) pelanggar.push(`${nama} → ${fn}`)
      }
    }
    expect(pelanggar).toEqual([])
  })

  // ⛔ Dan daftar di atas harus MENGIKUTI `api.ts`. Fungsi tulis yang lahir
  // tanpa didaftarkan membuat uji di atas lulus tanpa memeriksa apa pun.
  //
  // ⭐ `POST` SENDIRI BUKAN TANDA PENULISAN, dan penjaga bentuk pertama
  // salah di titik ini: `hitungPeriodePelaporan` ber-`POST` semata-mata
  // karena ia membawa badan permintaan. Rutenya menyatakan dirinya
  // *"MURNI menghitung dari isian layar: nol baca, nol tulis basis data"*.
  //
  // Pembedanya AWALAN JALUR `/hitung/` — konvensi yang sudah hidup di modul
  // ini. Rute hitung menerima isian layar dan mengembalikan angka; ia tidak
  // pernah menyentuh tabel, jadi ia tidak melanggar Aturan B sekalipun
  // dipanggil dari `onChange`.
  it('⛔ setiap fungsi TULIS `api.ts` terdaftar di penjaga ini', () => {
    const api = readFileSync(join(AKAR, 'api.ts'), 'utf8')
    const menulis: string[] = []
    for (const m of api.matchAll(
      /export async function (\w+)([\s\S]{0,700}?)metode:\s*'(POST|PUT|PATCH|DELETE)'/g,
    )) {
      const nama = m[1]
      const badan = m[2] ?? ''
      if (nama === undefined) continue
      if (badan.includes('/hitung/')) continue // menghitung, bukan menyimpan
      menulis.push(nama)
    }
    for (const n of menulis) {
      expect(FUNGSI_TULIS).toContain(n)
    }
  })

  // ⚠️ Dan rute `/hitung/` WAJIB tetap tidak menyentuh basis data. Begitu
  // salah satunya menyimpan, pengecualian di atas berubah menjadi lubang.
  it('⚠️ rute `/hitung/` menyatakan dirinya nol tulis', () => {
    const rute = readFileSync(
      join(AKAR, '..', 'backend', 'handlers', 'rute_warisan.go'),
      'utf8',
    )
    const i = rute.indexOf('/hitung/')
    expect(i).toBeGreaterThan(0)
    expect(rute.slice(Math.max(0, i - 400), i)).toContain('nol tulis basis data')
  })

  // ⭐ 7 Oktober 2026 — jalur tulisnya ADA. Submit dan Decline offer kini
  // hidup, dan HANYA lewat callback tombol yang form berikan: tab ini sendiri
  // tidak memanggil fungsi tulis mana pun. Decline offer melewati konfirmasi
  // `TreatyInDeclineConfirmation` lebih dulu.
  it('⭐ Revision menulis HANYA dari tombolnya — `doubleclick`, persis cell 994', () => {
    const src = readFileSync(join(AKAR, 'pages', 'DaftarKontrakTreatyIn.tsx'), 'utf8')
    expect(src).toMatch(/onDoubleClick=\{\(\) => \{\s*tekanRevisi\(b\.id\)/)
    expect(src).toContain('mulaiRevisi(id)')
    expect(src.match(/mulaiRevisi\(/g)?.length).toBe(1)
  })

  it('⭐ Submit dan Decline offer hidup lewat callback tombol, bukan fungsi tulis di tab', () => {
    const src = readFileSync(join(AKAR, 'components', 'TabInfoSubmit.tsx'), 'utf8')
    for (const fn of FUNGSI_TULIS) expect(src).not.toContain(`${fn}(`)
    expect(src).toContain('onClick={onKirim}')
    expect(src).toContain('TOMBOL_TULIS.tanyaTolak')
    expect(src).toMatch(/onKirim=\{\(\) => \{\s*setKonfirmasi\(false\)\s*onTolak\?\.\(\)/)
    const form = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
    expect(form).toContain("tekanKirim('submit')")
    expect(form).toContain("tekanKirim('decline')")
  })
})
