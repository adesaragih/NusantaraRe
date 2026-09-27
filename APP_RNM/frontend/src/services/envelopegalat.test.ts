import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, cariPesertaLife, pesanGalat } from './api'

// Kontrak envelope galat, SISI KLIEN.
//
// ⛔ Sebab uji ini ada, dan itu sungguh terjadi: backend menulis
// `{"galat": "<kalimat>"}` sejak tiket 01 (`handlers.galat`), tetapi klien ini
// membaca `o.error`. Akibatnya SETIAP pesan backend — 401, 403, 409, 503 —
// jatuh ke teks bawaan "Permintaan ditolak backend". Pemakai melihat pita
// merah yang tidak menyebutkan apa pun, padahal backend mengirim kalimat yang
// tepat, dan tidak ada satu pun uji yang gagal.
//
// Cacatnya tidak berbunyi di satu sisi mana pun: backend benar, klien benar
// menurut komentarnya sendiri, dan hanya PERTEMUANNYA yang salah. Pasangan
// uji ini ada di `internal/handlers/envelopegalat_test.go`.

/** Menjawab satu permintaan dengan badan dan status tertentu. */
function jawabDengan(badan: string, status = 409): void {
  vi.stubGlobal(
    'fetch',
    vi.fn(() =>
      Promise.resolve(
        new Response(badan, {
          status,
          headers: { 'Content-Type': 'application/json' },
        }),
      ),
    ),
  )
}

beforeEach(() => {
  vi.unstubAllGlobals()
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('envelope galat backend', () => {
  it('kalimat dari kunci "galat" sampai ke pemakai', async () => {
    const kalimat = 'perpindahan itu tidak ada di tangga kerja klaim'
    jawabDengan(JSON.stringify({ galat: kalimat }))

    const e = await cariPesertaLife('UJI-PL-1').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(ApiFailure)
    expect((e as ApiFailure).detail.message).toBe(kalimat)
    // Pesan yang dibaca layar, bukan hanya medan di dalamnya.
    expect((e as ApiFailure).message).toBe(kalimat)
    expect(pesanGalat(e)).toBe(kalimat)
  })

  it('kunci "error" TIDAK diterima diam-diam', async () => {
    // ⛔ Menerima kedua kunci akan menambal gejalanya dan MENYEMBUNYIKAN
    // penyebabnya: sejak itu kedua sisi tidak pernah dipaksa bertemu lagi,
    // dan perbedaan berikutnya akan lolos dengan cara yang sama persis.
    // Backend kita tidak pernah menulis `error`; klien karena itu tidak
    // boleh mengenalinya.
    jawabDengan(JSON.stringify({ error: 'kalimat dari envelope yang salah' }))

    const e = await cariPesertaLife('UJI-PL-1').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(ApiFailure)
    expect((e as ApiFailure).detail.message).toBeUndefined()
    expect((e as ApiFailure).message).toBe('Permintaan ditolak backend')
  })

  it('status tetap diteruskan apa adanya', async () => {
    jawabDengan(JSON.stringify({ galat: 'database belum dikonfigurasi' }), 503)
    const e = await cariPesertaLife('UJI-PL-1').catch((x: unknown) => x)
    expect((e as ApiFailure).status).toBe(503)
  })

  it('badan kosong tetap menjadi ApiFailure, bukan pecah', async () => {
    jawabDengan('', 500)
    const e = await cariPesertaLife('UJI-PL-1').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(ApiFailure)
    expect((e as ApiFailure).message).toBe('Permintaan ditolak backend')
  })

  it('galat kosong tidak menggantikan teks bawaan dengan kekosongan', async () => {
    // Pita merah yang isinya string kosong lebih buruk daripada teks bawaan:
    // ia terbaca seperti kerusakan layar, bukan penolakan.
    jawabDengan(JSON.stringify({ galat: '' }))
    const e = await cariPesertaLife('UJI-PL-1').catch((x: unknown) => x)
    expect((e as ApiFailure).message).toBe('Permintaan ditolak backend')
  })

  it('galat bukan teks diabaikan, bukan ditampilkan sebagai [object Object]', async () => {
    jawabDengan(JSON.stringify({ galat: { pesan: 'bersarang' } }))
    const e = await cariPesertaLife('UJI-PL-1').catch((x: unknown) => x)
    expect((e as ApiFailure).detail.message).toBeUndefined()
  })
})
