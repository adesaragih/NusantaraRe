// Penjaga TEKS LAYAR — nol catatan pengembang yang bocor ke pemakai.
//
// ---------------------------------------------------------------------
// ⛔ PERMINTAAN PEMILIK PROSES 8 Oktober 2026
// ---------------------------------------------------------------------
// *"hapus semua komen komen yg ada seperti ini di aplikasi untuk treaty in
// dan treaty adjustment"* — ditunjuk dengan tangkapan layar grid
// `Rate of Exchange` yang kosong, berbunyi:
//
//     Kurs per mata uang adalah tiket 20 (MATA_UANG_KONTRAK);
//     barisnya tidak dikarang.
//
// Empat puluh dua kalimat seperti itu hidup di kedua modul: nomor tiket,
// nama tabel Oracle, nama berkas ekspor Pega, cacah baris hasil sapuan, dan
// kabar "belum dibangun". Semuanya ditulis untuk PENGEMBANG, lalu
// dirender ke layar PEMAKAI.
//
// ---------------------------------------------------------------------
// ⚠️ MENGAPA PENJAGA, BUKAN SEKADAR DIHAPUS
// ---------------------------------------------------------------------
// Kalimat-kalimat itu lahir dengan niat baik — menerangkan mengapa sebuah
// grid kosong — dan niat yang sama akan melahirkannya lagi. Satu sapuan
// hari ini tidak menahan ronde berikutnya; uji ini menahan.
//
// ⭐ Yang DILARANG hanya istilah internal. Kalimat polos yang menolong
// pemakai TETAP BOLEH, dan beberapa memang sengaja dipertahankan:
//
//     'Belum ada baris EGNPI.'
//     'Tab ini belum dibangun.'
//     'Kosongkan penyaringnya untuk melihat seluruh daftar.'

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Istilah yang NOL urusannya dengan pemakai. */
const TERLARANG: ReadonlyArray<readonly [RegExp, string]> = [
  [/\btiket\s*\d/i, 'nomor tiket papan kerja'],
  [/POOLDATA|M_TREATY|M_ATTACHMENT|TREATY_IN\b|VERSI_KONTRAK|MATA_UANG_KONTRAK/, 'nama tabel Oracle'],
  [/\.xml\b|Section\/|Activity\s[A-Z]|flow action/i, 'nama berkas atau aturan Pega'],
  [/sistem lama|tabel warisan|dokumen warisan|tabel pendaratan/i, 'jeroan migrasi'],
  [/belum dibangun di layar|belum disalin dari ekspor|jalur simpan/i, 'kabar pekerjaan yang belum selesai'],
  [/dari 1\.85\d|dari 300 dokumen|yang disapu/i, 'cacah hasil sapuan korpus'],
]

/** Berkas teks layar modul ini — di sinilah seluruh kalimat pemakai hidup. */
function berkasLabel(): string[] {
  return readdirSync(AKAR)
    .filter((f) => f.startsWith('labels') && f.endsWith('.ts') && !f.includes('.test.'))
    .map((f) => join(AKAR, f))
}

/** Teks di dalam petik, dari baris yang BUKAN komentar. */
function teksLayar(berkas: string, minimal = 24): Array<{ baris: number; teks: string }> {
  const out: Array<{ baris: number; teks: string }> = []
  const petik = /'((?:[^'\\]|\\.)*)'|"((?:[^"\\]|\\.)*)"/g
  readFileSync(berkas, 'utf8')
    .split('\n')
    .forEach((l, i) => {
      const t = l.trim()
      // ⭐ Komentar BOLEH menyebut apa pun — di sanalah keterangannya
      // sekarang tinggal. Yang dijaga hanya yang sampai ke layar.
      if (t.startsWith('//') || t.startsWith('*') || t.startsWith('/*')) return
      if (t.startsWith('import ') || t.startsWith('} from ')) return
      for (let m = petik.exec(l); m !== null; m = petik.exec(l)) {
        const v = m[1] ?? m[2] ?? ''
        if (v.length > minimal) out.push({ baris: i + 1, teks: v })
      }
    })
  return out
}

describe('teks layar tidak membocorkan jeroan', () => {
  it('⛔ nol kalimat layar menyebut tiket, nama tabel, ekspor Pega, atau kabar migrasi', () => {
    const langgar: string[] = []
    for (const b of berkasLabel()) {
      for (const { baris, teks } of teksLayar(b)) {
        for (const [pola, sebab] of TERLARANG) {
          if (pola.test(teks)) {
            langgar.push(`${b.slice(AKAR.length + 1)}:${baris} — ${sebab}\n      ${teks.slice(0, 120)}`)
            break
          }
        }
      }
    }
    expect(langgar.join('\n')).toBe('')
  })

  // ⭐ CERMIN: kalimat polos TETAP ADA. Tanpa uji ini, cara termudah membuat
  // uji di atas hijau adalah mengosongkan SELURUH teks layar.
  it('⭐ kalimat polos yang menolong pemakai tidak ikut terbuang', () => {
    const semua = berkasLabel()
      .flatMap((b) => teksLayar(b, 0))
      .map((x) => x.teks)
    for (const tetap of [
      'Belum ada baris EGNPI.',
      'Belum ada baris Maximum Retention.',
      'Kosongkan penyaringnya untuk melihat seluruh daftar.',
    ]) {
      expect(semua).toContain(tetap)
    }
  })
})
