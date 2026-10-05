// Aturan layar Company Detail - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend
// (`backend/services`); yang di sini hanya mencegah kiriman yang pasti ditolak dan menyusun tampilan.

import type { Opsi } from '../../../inti/frontend/components/ui/dasar'
import type {
  Akun,
  Alamat,
  Detail,
  HasilSalinLama,
  Isian,
  Negara,
  OrgLama,
  PIC,
  Pilihan,
  StatusSalinLama,
  Telfax,
} from './api'
import { CD } from './labels'

/** Isian Create: tanpa PIC dan alamat. */
export function isianKosong(): Isian {
  return { nama: '', title: '', npwp: '', country: '', businessField: '', parentId: '', note: '', pic: [], alamat: [] }
}

/** Isian ubah dari organisasi tersimpan. */
export function isianDari(d: Detail): Isian {
  return {
    nama: d.nama,
    title: d.title,
    npwp: d.npwp,
    country: d.country,
    businessField: d.businessField,
    parentId: d.parentId,
    note: d.note,
    pic: d.pic.map((p) => ({ ...p })),
    alamat: d.alamat.map((a) => ({ ...a, telfax: a.telfax.map((t) => ({ ...t })) })),
  }
}

export function picKosong(): PIC {
  return { userIdentifier: '', nama: '', position: '', gender: '', email: '', dateOfBirth: '', phone: '' }
}

export function alamatKosong(): Alamat {
  return { asal: '', type: '', address: '', telfax: [] }
}

export function telfaxKosong(): Telfax {
  return { type: '', code: '', no: '' }
}

/**
 * Galat pertama yang pasti ditolak backend; `null` = boleh dikirim. `title` = pilihan Title; `namaLama` = nama
 * tersimpan saat Edit (nama lama yang tidak diubah boleh tetap memuat title).
 */
export function periksa(isi: Isian, title: Pilihan[] = [], namaLama?: string): string | null {
  if (isi.nama.trim() === '') return CD.galatNama
  const t = galatTitle(isi.nama, title, namaLama)
  if (t !== null) return t
  if (isi.country.trim() === '') return CD.galatCountry
  if (isi.businessField.trim() === '') return CD.galatBusinessField
  for (const [i, p] of isi.pic.entries()) {
    if (p.nama.trim() === '') return CD.galatPICNama(i + 1)
    if (p.position.trim() === '') return CD.galatPICPosisi(i + 1)
  }
  const letak = new Map<string, number>()
  for (const [i, a] of isi.alamat.entries()) {
    if (a.type.trim() === '') return CD.galatAlamatType(i + 1)
    const teks = a.address.trim()
    if (teks === '') return CD.galatAlamat(i + 1)
    const sebelum = letak.get(teks)
    if (sebelum !== undefined) return CD.galatAlamatGanda(i + 1, sebelum)
    letak.set(teks, i + 1)
    for (const [j, t] of a.telfax.entries()) {
      if (t.type.trim() === '') return CD.galatJenisNomor(i + 1, j + 1)
      if (t.no.trim() === '') return CD.galatNomor(i + 1, j + 1)
    }
  }
  return null
}

/** Teks satu pilihan: label, atau kodenya bila label kosong (mis. kode area tanpa nama). */
export function teksPilihan(p: Pilihan, denganKode = false): string {
  if (p.label === '') return p.kode
  return denganKode ? `${p.kode} - ${p.label}` : p.label
}

/**
 * Opsi dropdown: pilihan AKTIF, ditambah nilai yang sedang terpasang bila tidak aktif lagi (nilai lama Pega) -
 * ditandai "(old value)" supaya tetap terbaca tanpa ditawarkan untuk baris baru.
 */
export function opsiDari(daftar: Pilihan[], terpasang: string, denganKode = false): Opsi[] {
  const opsi = daftar.filter((p) => p.aktif).map((p) => ({ value: p.kode, label: teksPilihan(p, denganKode) }))
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    const lama = daftar.find((p) => p.kode === terpasang)
    opsi.push({ value: terpasang, label: CD.nilaiLama(lama === undefined ? terpasang : teksPilihan(lama, denganKode)) })
  }
  return opsi
}

/** Opsi Title: nilainya LABEL (`CLIENT.TITLE` menyimpan label, seperti `RDBINSERTCLIENT`). */
export function opsiTitle(daftar: Pilihan[], terpasang: string): Opsi[] {
  const opsi = daftar.filter((p) => p.aktif).map((p) => ({ value: p.label, label: p.label }))
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    opsi.push({ value: terpasang, label: CD.nilaiLama(terpasang) })
  }
  return opsi
}

/** Kode `CLIENT.COUNTRY` sebuah negara: OLDID (bentuk data lama), atau ID bila tanpa OLDID. */
export function kodeNegara(n: Negara): string {
  return n.oldId === '' ? n.id : n.oldId
}

/** Opsi COUNTRY dari NATION; kode lama yang tidak ada di NATION tetap tampil sebagai nilai lama. */
export function opsiNegara(negara: Negara[], terpasang: string, namaTerpasang = ''): Opsi[] {
  const opsi = negara.map((n) => ({ value: kodeNegara(n), label: n.nama === '' ? kodeNegara(n) : n.nama }))
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    opsi.push({ value: terpasang, label: CD.nilaiLama(namaTerpasang === '' ? terpasang : namaTerpasang) })
  }
  return opsi
}

/**
 * Opsi PIC Name: nama akun login aktif (perintah work owner 05-10-2026: "dropdown dari tabel login"). Nilainya NAMA
 * (`CLIENT_PICLIST.NICKNAME` menyimpan nama); nama kembar tampil sekali. Nama PIC lama Pega yang bukan akun tetap
 * terbaca sebagai nilai lama.
 */
export function opsiAkun(akun: Akun[], terpasang: string): Opsi[] {
  const opsi: Opsi[] = []
  for (const a of akun) {
    if (!opsi.some((o) => o.value === a.nama)) opsi.push({ value: a.nama, label: a.nama })
  }
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    opsi.push({ value: terpasang, label: CD.nilaiLama(terpasang) })
  }
  return opsi
}

/** JOB_POSITION akun bernama `nama` (PIC Position, perintah work owner 05-10-2026); `undefined` = bukan akun login. */
export function jabatanAkun(akun: Akun[], nama: string): string | undefined {
  return akun.find((a) => a.nama === nama)?.jabatan
}

/** Nilai opsi "Others" kode area - bukan kode sungguhan, tidak pernah dikirim. */
export const KODE_LAIN = '__lain__'

/**
 * Opsi kode area: kodehp aktif (dan nilai lama yang terpasang) + "Others" untuk kode yang diisi sendiri (perintah
 * work owner 05-10-2026: "tambahin others, bisa isi sendiri").
 */
export function opsiKodeArea(kode: Pilihan[], terpasang: string): Opsi[] {
  return [...opsiDari(kode, terpasang, true), { value: KODE_LAIN, label: CD.lainnya }]
}

/** Kode area terpasang yang tidak ada di daftar kodehp = "Others" yang diisi sendiri. */
export function kodeAreaLain(kode: Pilihan[], terpasang: string): boolean {
  return terpasang !== '' && !kode.some((p) => p.kode === terpasang)
}

/** Label sebuah kode di daftar pilihan; kode tak dikenal dikembalikan apa adanya. */
export function labelKode(daftar: Pilihan[], kode: string): string {
  if (kode === '') return ''
  const p = daftar.find((x) => x.kode === kode)
  return p === undefined ? kode : teksPilihan(p)
}

/** Kata nama huruf besar: titik dibuang (`P.T.` -> `PT`), tanda baca lain pemisah - sama dengan backend `kataNama`. */
export function kataNama(nama: string): string[] {
  return nama
    .replace(/\./g, '')
    .toUpperCase()
    .split(/[^\p{L}\p{N}]+/u)
    .filter((k) => k !== '')
}

/**
 * LABEL title AKTIF yang muncul sebagai KATA di nama (`PT ABC`, `P.T. ABC`); kosong = tidak ada. Title nonaktif
 * (TN., NY., NN. - sapaan orang) tidak dihitung, sama dengan backend.
 */
export function titleDiNama(nama: string, title: Pilihan[]): string {
  const peta = new Map(
    title.filter((p) => p.aktif).map((p) => [p.label.replace(/\./g, '').toUpperCase().trim(), p.label]),
  )
  for (const k of kataNama(nama)) {
    const label = peta.get(k)
    if (label !== undefined && k !== '') return label
  }
  return ''
}

/** Galat title di nama (work owner 04-10-2026: "ga boleh disimpan"); Edit: hanya bila nama diubah. */
export function galatTitle(nama: string, title: Pilihan[], namaLama?: string): string | null {
  const t = titleDiNama(nama, title)
  if (t === '') return null
  if (namaLama !== undefined && nama.trim().toUpperCase() === namaLama.trim().toUpperCase()) return null
  return CD.galatTitleDiNama(t)
}

/** Jumlah halaman daftar (sekurangnya 1). */
export function jumlahHalaman(total: number, ukuran: number): number {
  return Math.max(1, Math.ceil(total / Math.max(1, ukuran)))
}

/** Ganti satu unsur larik tanpa mengubah larik asalnya. */
export function gantiDi<T>(larik: T[], i: number, unsur: T): T[] {
  return larik.map((x, j) => (j === i ? unsur : x))
}

/** Buang satu unsur larik tanpa mengubah larik asalnya. */
export function buangDi<T>(larik: T[], i: number): T[] {
  return larik.filter((_, j) => j !== i)
}

/** Baris popup Copy Old yang cocok dengan kata cari (ORG ID atau nama, tanpa beda huruf). */
export function saringLama(daftar: OrgLama[], kata: string): OrgLama[] {
  const k = kata.trim().toUpperCase()
  if (k === '') return daftar
  return daftar.filter((o) => o.idView.toUpperCase().includes(k) || o.nama.toUpperCase().includes(k))
}

/** Catatan satu baris Copy Old: apa yang akan disalin. */
export function catatanLama(o: OrgLama): string {
  const bagian: string[] = []
  if (o.baru) bagian.push(CD.catatanBaru)
  if (o.pic > 0) bagian.push(CD.catatanPIC(o.pic))
  if (o.nomor > 0) bagian.push(CD.catatanNomor(o.nomor))
  if (o.isi.length > 0) bagian.push(CD.catatanIsi(o.isi))
  return bagian.join('; ')
}

/** Cacah hasil Process Copy per status. */
export function ringkasSalin(hasil: HasilSalinLama[]): Record<StatusSalinLama, number> {
  const r: Record<StatusSalinLama, number> = { disalin: 0, sudahAda: 0, ditolak: 0, gagal: 0 }
  for (const h of hasil) r[h.status]++
  return r
}

/** Ukuran satu permintaan Process Copy: setiap permintaan jauh di bawah batas waktu layar (30 detik). */
export const UKURAN_SALIN_LAMA = 20

/** Pecah ID menjadi kelompok berurutan berukuran `ukuran`. */
export function potong<T>(daftar: T[], ukuran: number): T[][] {
  const out: T[][] = []
  for (let i = 0; i < daftar.length; i += Math.max(1, ukuran)) out.push(daftar.slice(i, i + Math.max(1, ukuran)))
  return out
}
