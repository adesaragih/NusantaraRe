// Penyusun tata letak layar komite - fungsi murni (permintaan work owner 09-10-2026 "layout komite diperbaiki, lebih enak
// dilihat dan userfriendly"). Bagian dari server (`models.Layar.Bagian`, urutan Section ShowTransfer) dikelompokkan ke
// kartu berjudul; isi, label medan / grid, dan urutan di dalam kelompok TIDAK diubah. Tangga dirender sebagai langkah
// di kolom samping; ringkasan kepala dibaca dari medan / grid yang sudah dikirim server.

import type { Anggota, Bagian, Layar } from './api'
import { TATA_KCP } from './labels'

/** Satu kartu layar: judul + bagian server di dalamnya (urutan server). */
export interface Kelompok {
  kunci: string
  judul: string
  bagian: Bagian[]
}

/** Kelompok kartu menurut `Bagian.kunci` server (`modul/komiteclaimprop/backend/models/layar.go`), berurutan. */
const KELOMPOK: readonly { kunci: string; judul: string; isi: readonly string[] }[] = [
  { kunci: 'klaim', judul: TATA_KCP.kelompokKlaim, isi: ['klaim', 'klaimKanan'] },
  { kunci: 'kerugian', judul: TATA_KCP.kelompokKerugian, isi: ['kerugian'] },
  { kunci: 'estimasi', judul: TATA_KCP.kelompokEstimasi, isi: ['estimasi', 'totalEstimasi'] },
  { kunci: 'adjustment', judul: TATA_KCP.kelompokAdjustment, isi: ['riwayatAdjustment', 'deductible'] },
  { kunci: 'spreading', judul: TATA_KCP.kelompokSpreading, isi: ['spreading'] },
  { kunci: 'bayar', judul: TATA_KCP.kelompokBayar, isi: ['bayar', 'bank'] },
  { kunci: 'catatan', judul: TATA_KCP.kelompokCatatan, isi: ['teksKomite'] },
]

/** Kunci bagian tangga ("Committe Accept Status") - dirender sebagai langkah, bukan kartu. */
export const KUNCI_TANGGA = 'tangga'

/** Kelompokkan bagian server ke kartu; bagian yang tidak dikenal tetap tampil (kartu sendiri, di akhir). */
export function kelompokkanBagian(bagian: readonly Bagian[]): { kelompok: Kelompok[]; tangga?: Bagian } {
  const dikenal = new Set(KELOMPOK.flatMap((k) => k.isi))
  const kelompok: Kelompok[] = []
  for (const k of KELOMPOK) {
    const isi = bagian.filter((b) => k.isi.includes(b.kunci))
    if (isi.length > 0) kelompok.push({ kunci: k.kunci, judul: k.judul, bagian: isi })
  }
  for (const b of bagian) {
    if (!dikenal.has(b.kunci) && b.kunci !== KUNCI_TANGGA) {
      kelompok.push({ kunci: b.kunci, judul: b.judul ?? '', bagian: [b] })
    }
  }
  return { kelompok, tangga: bagian.find((b) => b.kunci === KUNCI_TANGGA) }
}

/** Keadaan satu tingkat tangga: disetujui, ditolak, sedang berjalan (menunggu PERTAMA, KomiteRouter S6.1), menunggu. */
export type KeadaanLangkah = 'setuju' | 'tolak' | 'berjalan' | 'menunggu'

export interface Langkah {
  urut: number
  jabatan: string
  keadaan: KeadaanLangkah
  tanggal: string
  komentar: string
}

/** Tangga kasus -> langkah (kode `KOMITE_APPROVAL`: 0 menunggu, 1 setuju, 2 tolak). */
export function langkahTangga(tangga: readonly Anggota[]): Langkah[] {
  let berjalan = false
  return tangga.map((a) => {
    let keadaan: KeadaanLangkah = 'menunggu'
    if (a.keputusan === '1') keadaan = 'setuju'
    else if (a.keputusan === '2') keadaan = 'tolak'
    else if (!berjalan) {
      keadaan = 'berjalan'
      berjalan = true
    }
    return { urut: a.urut, jabatan: a.jabatan, keadaan, tanggal: a.tanggal, komentar: a.komentar }
  })
}

/** Ringkasan kepala layar komite. */
export interface Ringkasan {
  noKlaim: string
  polis: string
  tertanggung: string
  /** Baris History Adjustment kasus komite ini (KomiteNo = ID kasus). */
  adjustment: { mataUang: string; nilai: string } | null
  /** Tingkat yang sedang berjalan; `null` = tidak ada yang menunggu. */
  tingkat: { ke: number; dari: number; jabatan: string } | null
}

function nilaiMedan(l: Layar, kunci: string, label: string): string {
  const m = l.bagian.find((b) => b.kunci === kunci)?.medan?.find((x) => x.label === label)
  return m?.nilai.trim() ?? ''
}

export function ringkasan(l: Layar): Ringkasan {
  const riwayat = l.bagian.find((b) => b.kunci === 'riwayatAdjustment')?.grid?.[0]?.baris ?? []
  const baris = riwayat.find((b) => b.KomiteNo === l.kasus.id)
  const tangga = l.kasus.tangga ?? []
  const i = tangga.findIndex((a) => a.keputusan === '0')
  const berjalan = i >= 0 ? tangga[i] : undefined
  return {
    noKlaim: nilaiMedan(l, 'klaimKanan', 'NoClaim') || l.kasus.klaimId,
    polis: nilaiMedan(l, 'klaim', 'Policy No'),
    tertanggung: nilaiMedan(l, 'klaim', 'Insured Name'),
    adjustment: baris ? { mataUang: baris.Currency ?? '', nilai: baris.AdjustmentValue ?? '' } : null,
    tingkat: berjalan ? { ke: berjalan.urut, dari: tangga.length, jabatan: berjalan.jabatan } : null,
  }
}
