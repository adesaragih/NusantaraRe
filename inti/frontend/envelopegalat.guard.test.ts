import { readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import { akarSumberFrontend, berkasSumberTS, relatifAplikasi } from './uji/sumber'

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

// Seluruh sumber frontend non-uji: perakit, `inti/frontend`, dan setiap
// `modul/<nama>/frontend` (`uji/sumber.ts`, struktur tim satu folder per
// modul 30-09-2026) - bukan satu folder yang diam-diam menyempit.

describe('envelope galat punya satu jalan', () => {
  const berkas = berkasSumberTS()

  it('ada berkas sumber yang terbaca', () => {
    // Tanpa ini, penelusur yang rusak akan membuat penjaga di bawah lulus
    // atas nol berkas - hijau yang tidak memeriksa apa pun.
    expect(berkas.length).toBeGreaterThan(10)
    // Dan dari SETIAP akar kode: akar yang tak terbaca menyempitkan penjaga.
    for (const akar of akarSumberFrontend()) {
      expect(berkas.some((f) => f.startsWith(akar)), relatifAplikasi(akar)).toBe(true)
    }
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
    expect(pembaca.map(relatifAplikasi)).toEqual(['inti/frontend/klien.ts'])
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
