import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Penjaga statik: SATU jalan membaca galat backend, dan satu kunci saja.
//
// ⛔ Sebabnya dua cacat nyata yang keduanya tidak berbunyi di uji mana pun:
//
//  1. `services/api.ts` membaca `o.error` padahal backend menulis `galat`,
//     sehingga setiap pesan backend jatuh ke teks bawaan.
//  2. `pages/RegisterKlaim.tsx` punya pembaca galatnya SENDIRI berbentuk
//     axios (`e.response.data.galat`) - dan axios sudah dibuang di F0.2,
//     sehingga setiap penolakan backend tampil sebagai "Gagal menghubungi
//     server." Itu menuduh jaringan padahal backend menjawab.
//
// Keduanya sejenis: pembaca yang bentuknya tidak lagi cocok dengan yang
// dikirim, dan tidak ada satu pun uji yang memaksa keduanya bertemu.

// Dari letak berkas ini, bukan `process.cwd()`: sejak struktur tim satu folder
// per modul, `npm` dijalankan dari APP_RNM/, bukan dari frontend/.
const SRC = join(__dirname, '..')

/** Seluruh berkas .ts/.tsx di bawah src, kecuali berkas uji. */
function berkasSumber(dir = SRC, out: string[] = []): string[] {
  for (const nama of readdirSync(dir)) {
    const jalur = join(dir, nama)
    if (statSync(jalur).isDirectory()) {
      berkasSumber(jalur, out)
    } else if (
      (nama.endsWith('.ts') || nama.endsWith('.tsx')) &&
      !nama.includes('.test.')
    ) {
      out.push(jalur)
    }
  }
  return out
}

describe('envelope galat punya satu jalan', () => {
  const berkas = berkasSumber()

  it('ada berkas sumber yang terbaca', () => {
    // Tanpa ini, penelusur yang rusak akan membuat penjaga di bawah lulus
    // atas nol berkas - hijau yang tidak memeriksa apa pun.
    expect(berkas.length).toBeGreaterThan(10)
  })

  it('nol pembaca berbentuk transport lama (axios)', () => {
    // `e.response.data` adalah bentuk axios. Klien fetch kita melempar
    // `ApiFailure`, yang tidak punya `.response`.
    const tertuduh = berkas.filter((f) =>
      /\.response\s*\?\.\s*data|response\?\s*:\s*\{\s*data/.test(
        readFileSync(f, 'utf8'),
      ),
    )
    expect(tertuduh, 'bentuk axios tersisa; axios dibuang di F0.2').toEqual([])
  })

  it('kunci envelope hanya dibaca di inti/klien.ts, dan namanya "galat"', () => {
    const pembaca = berkas.filter((f) => /\bo\.(galat|error)\b/.test(readFileSync(f, 'utf8')))
    // Satu rumah saja. Pembaca kedua berarti dua aturan yang dapat bergeser
    // sendiri-sendiri, dan yang bergeser tidak akan terlihat.
    expect(pembaca.map((f) => f.replace(SRC, 'src'))).toEqual([
      join('src', 'inti', 'klien.ts'),
    ])
    expect(readFileSync(pembaca[0]!, 'utf8')).toContain('o.galat')
  })

  it('nol pembacaan kunci "error" di seluruh sumber', () => {
    // ⛔ Bukan sekadar "api.ts benar": kunci yang salah tidak boleh muncul
    // lagi di mana pun, sebab ia akan tampak masuk akal bagi pembaca
    // berikutnya justru karena bahasa Inggrisnya.
    const tertuduh = berkas.filter((f) =>
      /\b(data|isi|o|jawaban|badan)\s*\.\s*error\b/.test(readFileSync(f, 'utf8')),
    )
    expect(tertuduh, 'backend menulis "galat", bukan "error"').toEqual([])
  })
})
