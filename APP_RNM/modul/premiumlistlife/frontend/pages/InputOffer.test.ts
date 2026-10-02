import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KEPUTUSAN_POLIS } from '../labels'
import { bolehRejectDiTahap, TAHAP_POLIS } from '../api'
import { ringkasanAkibat } from './InputOffer'

// Uji layar keputusan penawaran — tiket 01 bagian 2; GILIRAN-14 butir bq.

const BERKAS = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
/** Sumber tanpa komentar — prosa yang menerangkan tidak dituduh kode. */
const SUMBER = BERKAS.replace(/\{\/\*[\s\S]*?\*\/\}/g, '')
  .split('\n')
  .filter((b) => {
    const t = b.trimStart()
    return !t.startsWith('//') && !t.startsWith('*')
  })
  .join('\n')

describe('Reject hanya di tahap yang punya konektornya', () => {
  it('tahap penawaran TIDAK menawarkan Reject', () => {
    // ⛔ `Reject` muncul TEPAT SEKALI di seluruh flow — Transition9 b2306
    // pada Decision2, sesudah Input Premium Detail. Tombol yang pasti
    // dijawab 409 adalah tombol yang mengajari orang mengabaikan galat.
    expect(bolehRejectDiTahap(TAHAP_POLIS.penawaran)).toBe(false)
  })

  it('tahap detail menawarkannya', () => {
    expect(bolehRejectDiTahap(TAHAP_POLIS.detail)).toBe(true)
  })

  it('tahap summary tidak — ia nol konektor keputusan', () => {
    expect(bolehRejectDiTahap(TAHAP_POLIS.summary)).toBe(false)
  })

  it('layar memagarinya, bukan hanya backend', () => {
    expect(SUMBER).toContain('bolehRejectDiTahap(tahap)')
  })
})

describe('Decision3 tidak ditanyakan (butir bq)', () => {
  it('nol tombol Offer/Premium dan nol rute penggolong di layar', () => {
    // ⛔ `Offer`/`Premium` hasil decision table `IsFlagOnGoingPolicy` atas
    // bendera kasus — backend menerapkannya di dalam `Confirm`.
    expect(SUMBER).not.toContain('golongkanPenawaran')
    expect(SUMBER).not.toContain('PENGGOLONG_POLIS')
    expect(SUMBER).not.toContain('menungguPenggolong')
  })

  it('setiap keputusan yang berhasil melepas kasus dari layar ini', () => {
    // Sesudah bq, `Confirm` selalu menutup (Offer) atau memindahkan
    // (Premium) — tidak ada keadaan "menunggu" yang menahan layar.
    expect(SUMBER).toMatch(/setAkibat\(hasil\)\s*\n\s*onSelesai\(\)/)
  })
})

describe('ringkasanAkibat', () => {
  it('Confirm + bendera "0" (Offer) — tertutup Resolved-Completed', () => {
    const s = ringkasanAkibat({ tahapTujuan: '', statusWork: 'Resolved-Completed' })
    expect(s).toContain('Resolved-Completed')
  })

  it('Confirm + bendera "1" (Premium) — berpindah ke Input Premium Detail', () => {
    const s = ringkasanAkibat({ tahapTujuan: TAHAP_POLIS.detail, statusWork: '' })
    expect(s).toContain(TAHAP_POLIS.detail)
  })

  it('berpindah menyebut tahap tujuannya', () => {
    const s = ringkasanAkibat({ tahapTujuan: TAHAP_POLIS.penawaran, statusWork: '' })
    expect(s).toContain(TAHAP_POLIS.penawaran)
  })
})

describe('label VERBATIM', () => {
  it('ketiga keputusan', () => {
    expect(KEPUTUSAN_POLIS.confirm).toBe('Confirm')
    expect(KEPUTUSAN_POLIS.reject).toBe('Reject')
    expect(KEPUTUSAN_POLIS.decline).toBe('Decline')
  })
})

describe('periode produksi mm/yyyy', () => {
  it('YYYY-MM dari server tampil sebagai MM/YYYY', async () => {
    const { periodeTampil } = await import('./InputOffer')
    expect(periodeTampil('2026-10')).toBe('10/2026')
    expect(periodeTampil('2026-3')).toBe('03/2026')
  })

  it('bentuk asing ditampilkan apa adanya', async () => {
    const { periodeTampil } = await import('./InputOffer')
    expect(periodeTampil('Q4 2026')).toBe('Q4 2026')
    expect(periodeTampil('')).toBe('')
  })
})

describe('judul halaman', () => {
  it('tahap penawaran berjudul "Input Offer Life", tahap detail "Input Premium Detail"', async () => {
    const { judulKeputusan } = await import('./InputOffer')
    expect(judulKeputusan(TAHAP_POLIS.penawaran)).toBe('Input Offer Life')
    expect(judulKeputusan(TAHAP_POLIS.detail)).toBe('Input Premium Detail')
  })
})
