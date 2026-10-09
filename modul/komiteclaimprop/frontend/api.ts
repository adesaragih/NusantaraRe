// Klien Komite Claim Prop - `/api/komite-claim-prop` (`modul/komiteclaimprop/backend/handlers/rute.go`). Layar datang
// dari server sudah tersusun (`models.Layar`): bagian, label VERBATIM Section ShowTransfer, nilai klaim induk; layar
// hanya merender dan mengirim isian keputusan.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_KCP = '/api/komite-claim-prop'

export type JenisNilai = 'teks' | 'teksPanjang' | 'angka' | 'tanggal' | 'tanggalJam'

/** Satu sel hanya-baca berlabel - `models.Medan`. */
export interface Medan {
  label: string
  nilai: string
  jenis: JenisNilai | ''
}

export interface KolomGrid {
  label: string
  properti: string
  jenis: JenisNilai
}

export interface Grid {
  judul: string
  kolom: KolomGrid[]
  baris: Record<string, string>[]
}

export interface Bagian {
  kunci: string
  judul?: string
  medan?: Medan[]
  grid?: Grid[]
  sel?: Medan[][]
}

export interface Pilihan {
  nilai: string
  label: string
}

export interface Tombol {
  label: string
  aksi: 'lihat' | 'batal' | 'putuskan'
  aktif: boolean
  alasan?: string
}

/** Isian Submit - `models.Keputusan`. */
export interface Keputusan {
  acceptStatus: string
  comment: string
  isSubjectivity: boolean
  subjectivityNote: string
  usulTutup: boolean
  usulCadang: boolean
}

export interface IsianLayar {
  nilai: Keputusan
  /** `pyDisabledWhen .KomiteCount!='1'` salah. */
  terbuka: boolean
  pilihanTerima: Pilihan[]
  /** Dropdown "Subjectivity Note" (SubjectivityNote.xml). */
  pilihanSubjectivityNote: Pilihan[]
  label: Record<keyof Keputusan, string>
}

export interface Anggota {
  id: string
  urut: number
  operatorId: string
  jabatan: string
  keputusan: string
  komentar: string
  tanggal: string
}

export interface Kasus {
  id: string
  klaimId: string
  adjustmentId: string
  komiteLoop: number
  komiteCount: number
  acceptStatus: string
  statusWork: string
  pembuatNama: string
  tangga: Anggota[]
}

export interface Layar {
  kasus: Kasus
  judul: string[]
  bagian: Bagian[]
  isian: IsianLayar
  tombol: Tombol[]
  bolehKerja: boolean
  pesan?: string[]
}

/** Satu baris daftar kerja - `models.BarisKerja`. */
export interface BarisKerja {
  kasusId: string
  klaimId: string
  noKlaim: string
  tingkat: number
  komiteCount: number
  komiteLoop: number
  jabatan: string
  nilai: string
  mataUang: string
  statusBaris: string
  statusWork: string
  tglUpdate: string
}

export interface HasilKeputusan {
  kasusId: string
  selesai: boolean
  komiteCount: number
  acceptedNo?: string
  /** Keputusan tersimpan, tetapi PDF akseptasi gagal disimpan (`services.PesanDokumenGagal`). */
  galatDokumen?: string
}

export function daftarKerja(): Promise<BarisKerja[]> {
  return minta<BarisKerja[]>(`${PREFIX_KCP}/kasus`)
}

export function bukaKasus(id: string): Promise<Layar> {
  return minta<Layar>(`${PREFIX_KCP}/kasus/${encodeURIComponent(id)}`)
}

export function putuskan(id: string, kep: Keputusan): Promise<HasilKeputusan> {
  return minta<HasilKeputusan>(`${PREFIX_KCP}/kasus/${encodeURIComponent(id)}/putuskan`, {
    metode: 'POST',
    badan: kep,
  })
}
