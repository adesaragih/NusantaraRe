import { describe, expect, it } from 'vitest'

import { BERANDA, KETERANGAN_BELUM_DIMIGRASI } from '../inti/frontend/labels'
import { folderKorpusBelumDimigrasi } from '../inti/frontend/uji/sumber'
import { FOLDER_KORPUS } from './katalogKorpus'
import { kartuModul, ringkasanAntrean, type AntreanTahap } from './Beranda'

// Uji Beranda — butir bg.

describe('kartuModul', () => {
  it('dua puluh kartu, satu per kelompok MODUL', () => {
    // 18 sejak Treaty Contract Out ditambahkan (28-09-2026).
    // 20 sejak brief menu M_NAV_MENU (30-09-2026): Treaty In dan Treaty In
    // Adjustment - kedua folder korpus terakhir - masuk FOLDER_KORPUS (dulu MODUL).
    // + Marketing Officer, modul di luar korpus (migrasi inti 906, keputusan work owner 03-10-2026).
    // + Company Detail, modul di luar korpus (migrasi inti 907, keputusan work owner 04-10-2026).
    // + Accounts, modul di luar korpus (migrasi inti 908, keputusan work owner 04-10-2026).
    // + Master Data, modul di luar korpus (migrasi inti 904, keputusan work owner 04-10-2026).
    expect(kartuModul()).toHaveLength(24)
  })

  it('modul dimigrasi bertujuan, yang belum dimigrasi tanpa', () => {
    // ⛔ Kartu yang belum dimigrasi tetap BERDIRI. Menyembunyikannya
    // membuat aplikasi tampak lengkap padahal enam belas modul belum ada.
    //
    // Struktur tim satu folder per modul (30-09-2026): dulu "empat bertujuan,
    // enam belas tanpa" - angka yang harus disunting setiap modul yang
    // mendapat butir menu pertamanya. Kini: keempat modul lama tetap
    // bertujuan, dan kartu tanpa tujuan = folder korpus yang `MODUL.md`-nya
    // menyatakan `belum dimigrasi`.
    const kartu = kartuModul()
    const aktif = kartu.filter((k) => k.tujuan !== null).map((k) => k.nama)
    for (const lama of [FOLDER_KORPUS.claimLife, FOLDER_KORPUS.komiteClaimLife, FOLDER_KORPUS.premiumListLife, FOLDER_KORPUS.treatyContractOut]) {
      expect(aktif).toContain(lama)
    }
    const tanpaTujuan = kartu.filter((k) => k.tujuan === null).map((k) => k.nama)
    expect(tanpaTujuan.sort()).toEqual(folderKorpusBelumDimigrasi())
    expect(tanpaTujuan.length).toBeGreaterThan(0)
  })

  it('kartu Claim Life menunjuk HALAMAN AWAL-nya (menu datar 30-09-2026)', () => {
    const cl = kartuModul().find((k) => k.nama === FOLDER_KORPUS.claimLife)
    expect(cl?.tujuan).toBe('inbox')
    expect(cl?.label).toBe(FOLDER_KORPUS.claimLife)
  })
})

describe('ringkasanAntrean', () => {
  it('menjumlahkan keempat tahap', () => {
    const antrean: AntreanTahap[] = [
      { nomor: 1, nama: 'Input Register', total: 3 },
      { nomor: 2, nama: 'Outstanding Claim', total: 5 },
      { nomor: 3, nama: 'Medical Check', total: 0 },
      { nomor: 4, nama: 'Claim Analis', total: 2 },
    ]
    expect(ringkasanAntrean(antrean)).toBe(`10 ${BERANDA.antrean}`)
  })

  it('nol antrean tetap ANGKA, bukan keadaan', () => {
    // ⚠️ Nol adalah angka; "belum ada kotak masuk" adalah keadaan.
    // Keduanya tidak boleh tertukar.
    const antrean: AntreanTahap[] = [{ nomor: 1, nama: 'x', total: 0 }]
    expect(ringkasanAntrean(antrean)).toBe(`0 ${BERANDA.antrean}`)
    expect(ringkasanAntrean(antrean)).not.toBe(BERANDA.tanpaAntrean)
  })

  it('endpoint yang belum ada berkata begitu, bukan nol', () => {
    // ⛔ Menampilkan nol untuk modul yang endpointnya belum ada berarti
    // berbohong dengan angka yang terlihat benar.
    expect(ringkasanAntrean(null)).toBe(BERANDA.tanpaAntrean)
  })
})

describe('label Beranda', () => {
  it('kata "belum dimigrasi" satu tempat', () => {
    expect(KETERANGAN_BELUM_DIMIGRASI).toBe('belum dimigrasi')
  })
})
