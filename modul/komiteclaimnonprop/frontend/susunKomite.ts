// Penyusun tata letak layar komite - fungsi murni (pola Komite Claim Prop, perintah work owner 09-10-2026 "ikuti tampilan
// klaim prop"). Bagian dari server (`models.Layar.Bagian`, urutan Section ShowTransfer) dikelompokkan ke kartu berjudul;
// isi, label medan / grid, dan urutan di dalam kelompok TIDAK diubah. Tangga dirender sebagai langkah; ringkasan kepala
// dibaca dari medan yang sudah dikirim server.

import type { Anggota, Bagian, Layar } from './api'
import { TATA_KCNP } from './labels'

/** Satu kartu layar: judul + bagian server di dalamnya (urutan server). */
export interface Kelompok {
  kunci: string
  judul: string
  bagian: Bagian[]
}

/** Kelompok kartu menurut `Bagian.kunci` server (`modul/komiteclaimnonprop/backend/models/layar.go`), berurutan. */
const KELOMPOK: readonly { kunci: string; judul: string; isi: readonly string[] }[] = [
  { kunci: 'klaim', judul: TATA_KCNP.kelompokKlaim, isi: ['klaim', 'klaimKanan'] },
  { kunci: 'kerugian', judul: TATA_KCNP.kelompokKerugian, isi: ['kejadian'] },
  { kunci: 'akseptasi', judul: TATA_KCNP.kelompokAkseptasi, isi: ['klaimAkseptasi', 'totalKlaim'] },
  { kunci: 'alokasi', judul: TATA_KCNP.kelompokAlokasi, isi: ['alokasi', 'dibayar'] },
  { kunci: 'spreading', judul: TATA_KCNP.kelompokSpreading, isi: ['spreading'] },
  { kunci: 'bayar', judul: TATA_KCNP.kelompokBayar, isi: ['bayar', 'bank', 'bank2'] },
  { kunci: 'catatan', judul: TATA_KCNP.kelompokCatatan, isi: ['teksKomite'] },
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
  /** Tingkat yang sedang berjalan; `null` = tidak ada yang menunggu. */
  tingkat: { ke: number; dari: number; jabatan: string } | null
}

/** Label sel inline S9 ShowTransfer (`models.LabelNoKlaim`). */
const LABEL_NO_KLAIM = 'Claim No / Claim ID'

function nilaiMedan(l: Layar, kunci: string, label: string): string {
  const m = l.bagian.find((b) => b.kunci === kunci)?.medan?.find((x) => x.label === label)
  return m?.nilai.trim() ?? ''
}

export function ringkasan(l: Layar): Ringkasan {
  const tangga = l.kasus.tangga ?? []
  const i = tangga.findIndex((a) => a.keputusan === '0')
  const berjalan = i >= 0 ? tangga[i] : undefined
  return {
    // S9 "Claim No / Claim ID": `NoClaim "/" CLMNO`; NoClaim kosong = ID klaim saja
    noKlaim: nilaiMedan(l, 'klaimKanan', LABEL_NO_KLAIM).replace(/^\/\s*/, '') || l.kasus.klaimId,
    polis: nilaiMedan(l, 'klaim', 'Policy No'),
    tertanggung: nilaiMedan(l, 'klaim', 'Insured Name'),
    tingkat: berjalan ? { ke: berjalan.urut, dari: tangga.length, jabatan: berjalan.jabatan } : null,
  }
}
