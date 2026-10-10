// Hitung langsung saat mengetik (work owner 10-10-2026: "udah di ubah premi OGP, tapi deduction A tidak berubah
// otomatis, harus di triger dulu").

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { JEDA_HITUNG_MS, hitungSaatKetik, pertahankanKetikan } from './hitungLangsung'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })

describe('kapan ketikan dihitung tanpa menunggu blur', () => {
  it('medan ber-aksi dengan angka bukan nol; kosong / nol menunggu blur (CountOGPONP mengenolkan persen OGP/ONP)', () => {
    expect(hitungSaatKetik(true, '101340187.9')).toBe(true)
    expect(hitungSaatKetik(true, '0.5')).toBe(true)
    expect(['', '0', '0.00', ' '].map((v) => hitungSaatKetik(true, v))).toEqual([false, false, false, false])
    expect(hitungSaatKetik(false, '1000')).toBe(false)
    expect(JEDA_HITUNG_MS).toBeGreaterThanOrEqual(300)
    expect(JEDA_HITUNG_MS).toBeLessThanOrEqual(1000)
  })
})

describe('jawaban server tidak menimpa ketikan yang lebih baru', () => {
  it('medan yang berubah sesudah permintaan dikirim dipertahankan; hasil hitung server lainnya dipakai', () => {
    const dikirim = hal({ 'PolicyTreatyIn.PremiOgp': '123', 'PolicyTreatyIn.ResultOgp1': '1' })
    const kini = hal({ 'PolicyTreatyIn.PremiOgp': '1234', 'PolicyTreatyIn.ResultOgp1': '1' })
    const server = hal({ 'PolicyTreatyIn.PremiOgp': '123', 'PolicyTreatyIn.ResultOgp1': '39.975' })
    expect(pertahankanKetikan(server, dikirim, kini).nilai).toEqual({
      'PolicyTreatyIn.PremiOgp': '1234',
      'PolicyTreatyIn.ResultOgp1': '39.975',
    })
    // tidak ada ketikan baru: halaman server apa adanya
    expect(pertahankanKetikan(server, dikirim, dikirim)).toBe(server)
  })
})

describe('pemasangan', () => {
  it('isian angka: jeda sesudah ketikan, onSelesai terkini, blur tidak mengirim ulang nilai yang sudah terkirim', () => {
    const kotak = readFileSync(join(__dirname, 'components', 'KotakMedan.tsx'), 'utf8')
    expect(kotak).toContain('if (!hitungSaatKetik(!!medan.aksi?.length, x)) return')
    expect(kotak).toContain('selesaiKini.current(medan, x)')
    expect(kotak).toContain('}, JEDA_HITUNG_MS)')
    expect(kotak).toContain('if (!sudah) onSelesai(medan, kini)')
  })

  it('layar: jawaban yang sudah disusul diabaikan; ketikan sesudah kirim dipertahankan', () => {
    const layar = readFileSync(join(__dirname, 'pages', 'LayarKasus.tsx'), 'utf8')
    expect(layar).toContain('if (ke !== hitungTerakhir.current) return')
    expect(layar).toContain('setH((kini) => (kini ? pertahankanKetikan(ly.halaman, dikirim, kini) : ly.halaman))')
  })
})
