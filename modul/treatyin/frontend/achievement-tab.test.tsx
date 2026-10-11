// Uji tab Achievement In IDR — bentuk layar lawan ekspor Pega.
//
// RUMUSNYA diuji di services (`hitung_achievement_ringkas_test.go`); di sini
// hanya bentuknya: keenam kolom, tiga sel kaki, dan sifat hanya-bacanya.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { ACHIEVEMENT, DESIMAL_ACHIEVEMENT, KAKI_ACHIEVEMENT, KOLOM_ACHIEVEMENT } from './labelsAchievement'

const AKAR = __dirname
const SRC = readFileSync(join(AKAR, 'components', 'TabAchievement.tsx'), 'utf8')
const HAL = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')

describe('tab Achievement In IDR — bentuk dari ekspor', () => {
  // ⭐ Keenam kolom `AchievementCombine.xml`, URUT.
  it('⭐ keenam kolom ekspor, urut dan berejaan Pega', () => {
    expect(KOLOM_ACHIEVEMENT.map((k) => k.label)).toEqual([
      'Treaty Group',
      'Reins Type',
      'Gross Premium Before Claim in IDR',
      'Net Premium Before Claim in IDR',
      'Incured Claim in IDR',
      'Net Loss Ratio',
    ])
  })

  // ⛔ `Incured`, bukan `Incurred` — ejaan ekspor dipertahankan.
  it('⛔ ejaan Pega `Incured` dipertahankan', () => {
    const l = KOLOM_ACHIEVEMENT.map((k) => k.label).join('|')
    expect(l).toContain('Incured Claim in IDR')
    expect(l).not.toContain('Incurred')
  })

  // ⛔ KAKI HANYA TIGA SEL, dan `Gross Premium` BUKAN salah satunya.
  it('⛔ kaki hanya tiga sel; Gross Premium nol punya total', () => {
    expect(KAKI_ACHIEVEMENT).toHaveLength(KOLOM_ACHIEVEMENT.length)
    const terisi = KAKI_ACHIEVEMENT.filter((k) => k !== null)
    expect(terisi).toEqual(['SumTotalAchievNetPremium', 'SumTotalAchievIncured', 'SumLossRatio'])
    // Kolom 3 = Gross Premium; kakinya harus KOSONG.
    expect(KAKI_ACHIEVEMENT[2]).toBeNull()
  })

  // ⭐ Label kaki TAMPIL walau `pyCondition = 1=2`, sebab
  // `pyVisible = ALWAYS` — syarat hanya berlaku saat `pyVisible = OTHER`.
  it('⭐ label `Total in IDR` dirender', () => {
    expect(ACHIEVEMENT.totalIDR).toBe('Total in IDR')
    expect(SRC).toContain('{ACHIEVEMENT.totalIDR}')
  })

  // ⛔ SELURUH sel hanya-baca — nol `Field`, nol `onChange`, nol tombol.
  it('⛔ nol sel yang dapat diisi, nol tombol', () => {
    expect(SRC).not.toContain('<Field')
    expect(SRC).not.toContain('<button')
    expect(SRC).not.toContain('<input')
  })

  // ⛔ Nol rumus di layar — seluruhnya dari services.
  it('⛔ nol aritmetika di layar', () => {
    const kode = SRC.split('\n')
      .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
      .join('\n')
    expect(kode).toContain('hitungAchievement(')
    expect(kode).not.toContain('reduce(')
    expect(kode).not.toContain('* 100')
  })

  // ⭐ `Net Loss Ratio` PERSEN, lima lainnya bukan — supaya rasio tidak
  // diformat sebagai rupiah.
  it('⭐ rasio digolongkan persen, bukan uang', () => {
    const lr = KOLOM_ACHIEVEMENT.find((k) => k.kunci === 'LossRatio')
    expect(lr?.jenis).toBe('persen')
    expect(KOLOM_ACHIEVEMENT.filter((k) => k.jenis === 'uang')).toHaveLength(3)
    expect(DESIMAL_ACHIEVEMENT.persen).toBe(2)
  })

  // ⛔ Tabnya DIRENDER — sebelum ini namanya ada tetapi nol cabang
  // merendernya, sehingga membukanya menampilkan isi tab pertama.
  it('⛔ form punya cabang untuk tab ini', () => {
    expect(HAL).toContain("tabTampil === 'Achievement In IDR'")
    expect(HAL).toContain('<TabAchievement')
  })

  it('⭐ baris total dibedakan dari data di CSS', () => {
    expect(SRC).toContain('trin__ach-total')
    expect(CSS).toContain('.trin__ach-total')
  })
})
