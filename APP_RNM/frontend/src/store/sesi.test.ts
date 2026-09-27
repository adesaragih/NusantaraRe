// Uji sesi dan header identitas — F0.2.
//
// Bukti selesai yang brief tuntut: masuk → `X-Pelaku`/`X-Peran` terkirim;
// keluar → hilang; galat jaringan → `BACKEND_TIDAK_TERJANGKAU`.

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { PERAN } from '../assets/labels'
import { headerIdentitas, sesi, PERAN_TERSEDIA } from './sesi'

/** sessionStorage tiruan — lingkungan uji `node` tidak punya DOM. */
function pasangPenyimpanan(): Map<string, string> {
  const isi = new Map<string, string>()
  vi.stubGlobal('sessionStorage', {
    getItem: (k: string) => isi.get(k) ?? null,
    setItem: (k: string, v: string) => {
      isi.set(k, v)
    },
    removeItem: (k: string) => {
      isi.delete(k)
    },
  })
  return isi
}

beforeEach(() => {
  pasangPenyimpanan()
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('sesi disimpan dan dibaca kembali', () => {
  it('belum masuk berarti null, bukan sesi kosong', () => {
    expect(sesi.baca()).toBeNull()
  })

  it('masuk lalu baca mengembalikan akun dan perannya', () => {
    expect(sesi.simpan({ akunID: 'UJI-ADMIN', peran: [PERAN.admin] })).toBe(true)
    expect(sesi.baca()).toEqual({ akunID: 'UJI-ADMIN', peran: [PERAN.admin] })
  })

  it('keluar membuang identitasnya', () => {
    sesi.simpan({ akunID: 'UJI-ADMIN', peran: [PERAN.admin] })
    sesi.hapus()
    expect(sesi.baca()).toBeNull()
  })

  it('satu pelaku boleh memegang LEBIH DARI SATU peran', () => {
    // `pelakuDari` memecah `X-Peran` pada koma, jadi layar tidak boleh
    // memaksa memilih satu.
    sesi.simpan({ akunID: 'UJI-SEMUA', peran: [PERAN.admin, PERAN.medis, PERAN.spv] })
    expect(sesi.baca()?.peran).toHaveLength(3)
  })

  it('peran yang tidak dikenal DIBUANG saat dibaca kembali', () => {
    // ⛔ Isi sessionStorage dapat disunting siapa pun yang membuka devtools.
    sessionStorage.setItem(
      'rnm.sesi',
      JSON.stringify({ akunID: 'UJI-X', peran: [PERAN.admin, 'SuperAdmin'] }),
    )
    expect(sesi.baca()?.peran).toEqual([PERAN.admin])
  })

  it('sesi tanpa peran BUKAN sesi', () => {
    sessionStorage.setItem('rnm.sesi', JSON.stringify({ akunID: 'UJI-X', peran: [] }))
    expect(sesi.baca()).toBeNull()
  })

  it('sesi tanpa akun BUKAN sesi', () => {
    sessionStorage.setItem('rnm.sesi', JSON.stringify({ akunID: '  ', peran: [PERAN.admin] }))
    expect(sesi.baca()).toBeNull()
  })

  it('isi rusak tidak memecahkan layar', () => {
    sessionStorage.setItem('rnm.sesi', 'bukan json')
    expect(sesi.baca()).toBeNull()
  })

  it('penyimpanan yang MELEMPAR diperlakukan sebagai belum masuk', () => {
    // Mode privat / site-data diblokir.
    vi.stubGlobal('sessionStorage', {
      getItem: () => {
        throw new Error('ditolak')
      },
      setItem: () => {
        throw new Error('ditolak')
      },
      removeItem: () => {
        throw new Error('ditolak')
      },
    })
    expect(sesi.baca()).toBeNull()
    expect(sesi.simpan({ akunID: 'UJI', peran: [PERAN.admin] })).toBe(false)
    expect(() => {
      sesi.hapus()
    }).not.toThrow()
  })

  it('ketiga peran ADR-U-0002 tersedia, tidak lebih', () => {
    expect(PERAN_TERSEDIA).toEqual([PERAN.admin, PERAN.medis, PERAN.spv])
  })
})

describe('header identitas', () => {
  it('belum masuk → nol header, bukan header kosong', () => {
    // Header `X-Pelaku: ` yang kosong menyamarkan permintaan anonim sebagai
    // permintaan bernama di panel jaringan.
    expect(headerIdentitas()).toEqual({})
  })

  it('sesudah masuk → X-Pelaku dan X-Peran terkirim', () => {
    sesi.simpan({ akunID: 'UJI-ADMIN', peran: [PERAN.admin] })
    expect(headerIdentitas()).toEqual({
      'X-Pelaku': 'UJI-ADMIN',
      'X-Peran': PERAN.admin,
    })
  })

  it('peran ganda dipisah KOMA — persis yang pelakuDari urai', () => {
    sesi.simpan({ akunID: 'UJI-SEMUA', peran: [PERAN.admin, PERAN.spv] })
    expect(headerIdentitas()['X-Peran']).toBe(`${PERAN.admin},${PERAN.spv}`)
  })

  it('sesudah keluar → header hilang lagi', () => {
    sesi.simpan({ akunID: 'UJI-ADMIN', peran: [PERAN.admin] })
    expect(headerIdentitas()['X-Pelaku']).toBe('UJI-ADMIN')
    sesi.hapus()
    expect(headerIdentitas()).toEqual({})
  })
})
