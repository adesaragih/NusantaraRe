// Uji identitas pelaku dari env — F0.6.
//
// ⛔ Form login dibuang `[perintah work owner 27-09-2026]`. Yang diuji kini:
// identitas dibentuk dari env, gerbangnya gagal TERTUTUP, dan peran yang
// dikarang di env tidak lolos.

import { afterEach, describe, expect, it, vi } from 'vitest'

import { PERAN } from '../labels'
import { headerIdentitas, pelakuStub, PERAN_TERSEDIA } from './sesi'

/** Menyetel env Vite untuk satu kasus uji. */
function pasangEnv(isi: Record<string, string | undefined>): void {
  vi.stubEnv('VITE_AUTH_STUB', isi.VITE_AUTH_STUB ?? '')
  vi.stubEnv('VITE_STUB_PELAKU', isi.VITE_STUB_PELAKU ?? '')
  vi.stubEnv('VITE_STUB_PERAN', isi.VITE_STUB_PERAN ?? '')
}

afterEach(() => {
  vi.unstubAllEnvs()
})

describe('gerbang stub gagal TERTUTUP', () => {
  it.each([undefined, '', 'false', 'TRUE', '1', 'ya'])(
    'VITE_AUTH_STUB=%s → nol identitas',
    (nilai) => {
      // ⛔ Hanya `'true'` persis. Apa pun yang lain - termasuk TIDAK DISETEL -
      // berarti tidak ada identitas, dan backend menolak setiap jalur
      // beridentitas. Itu keadaan yang benar sampai IAM ada.
      pasangEnv({ VITE_AUTH_STUB: nilai })
      expect(pelakuStub()).toBeNull()
      expect(headerIdentitas()).toEqual({})
    },
  )
})

describe('identitas dari env', () => {
  it('akun dan peran dibaca apa adanya', () => {
    pasangEnv({
      VITE_AUTH_STUB: 'true',
      VITE_STUB_PELAKU: 'UJI-SPV',
      VITE_STUB_PERAN: PERAN.spv,
    })
    expect(pelakuStub()).toEqual({ akunID: 'UJI-SPV', peran: [PERAN.spv] })
  })

  it('peran ganda dipisah koma, spasi dirapikan', () => {
    pasangEnv({
      VITE_AUTH_STUB: 'true',
      VITE_STUB_PERAN: ` ${PERAN.admin} , ${PERAN.spv} `,
    })
    expect(pelakuStub()?.peran).toEqual([PERAN.admin, PERAN.spv])
  })

  it('akun kosong jatuh ke bawaan UJI-ADMIN', () => {
    pasangEnv({ VITE_AUTH_STUB: 'true', VITE_STUB_PELAKU: '   ' })
    expect(pelakuStub()?.akunID).toBe('UJI-ADMIN')
  })

  it('peran kosong jatuh ke KETIGA peran, bukan satu', () => {
    // ⚠️ Pengembang yang menyalakan stub hampir selalu ingin melihat seluruh
    // antrian. Menyempitkannya ke satu peran membuat tiga dari empat tab
    // menghilang tanpa sebab yang terlihat.
    pasangEnv({ VITE_AUTH_STUB: 'true', VITE_STUB_PERAN: '' })
    expect(pelakuStub()?.peran).toEqual([...PERAN_TERSEDIA])
  })

  it('peran yang TIDAK DIKENAL dibuang', () => {
    // ⛔ Env dapat salah ketik, dan peran yang dikarang akan ditolak backend
    // dengan 403 yang tidak dapat dijelaskan pemakai.
    pasangEnv({
      VITE_AUTH_STUB: 'true',
      VITE_STUB_PERAN: `${PERAN.admin},SuperAdmin,ReasLifeAdminX`,
    })
    expect(pelakuStub()?.peran).toEqual([PERAN.admin])
  })

  it('seluruh peran asing → jatuh ke bawaan, bukan nol peran', () => {
    // Nol peran berarti nol tab, dan layar yang kosong tanpa sebab lebih
    // buruk daripada layar yang menampilkan bawaan yang dinyatakan.
    pasangEnv({ VITE_AUTH_STUB: 'true', VITE_STUB_PERAN: 'Bukan,Peran' })
    expect(pelakuStub()?.peran).toEqual([...PERAN_TERSEDIA])
  })

  it('ketiga peran ADR-U-0002 tersedia, tidak lebih', () => {
    expect(PERAN_TERSEDIA).toEqual([PERAN.admin, PERAN.medis, PERAN.spv])
  })
})

describe('header identitas', () => {
  it('X-Pelaku dan X-Peran terkirim, peran dipisah KOMA', () => {
    pasangEnv({
      VITE_AUTH_STUB: 'true',
      VITE_STUB_PELAKU: 'UJI-SEMUA',
      VITE_STUB_PERAN: `${PERAN.admin},${PERAN.spv}`,
    })
    expect(headerIdentitas()).toEqual({
      'X-Pelaku': 'UJI-SEMUA',
      'X-Peran': `${PERAN.admin},${PERAN.spv}`,
    })
  })

  it('stub mati → nol header, bukan header kosong', () => {
    // Header `X-Pelaku: ` yang kosong menyamarkan permintaan anonim sebagai
    // permintaan bernama di panel jaringan.
    pasangEnv({ VITE_AUTH_STUB: 'false' })
    expect(headerIdentitas()).toEqual({})
  })
})
