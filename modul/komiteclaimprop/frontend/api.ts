// Klien Komite Claim Prop - `/api/komite-claim-prop` (`modul/komiteclaimprop/backend/handlers/rute.go`). Layar datang
// dari server sudah tersusun (`models.Layar`): bagian, label VERBATIM Section ShowTransfer, nilai klaim induk; layar
// hanya merender dan mengirim isian keputusan.

import { BATAS_WAKTU_MS, kegagalanDari, minta, rakitURL, unduhBerkasBeridentitas } from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

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
  /** `.TransferType = 2`: Subjectivity dan dua Propose tampil; false = kasus Close Without Payment (TT 4). */
  adjustment: boolean
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

export interface HasilKeputusan {
  kasusId: string
  selesai: boolean
  komiteCount: number
  acceptedNo?: string
  /** Keputusan tersimpan, tetapi PDF akseptasi gagal disimpan (`services.PesanDokumenGagal`). */
  galatDokumen?: string
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

/** Grid kategori lampiran klaim induk - HANYA kategori yang berisi (`models.SusunLampiranKomite`). */
export interface KategoriLampiran {
  id: string
  label: string
  countAttach: number
}

/** Satu dokumen klaim induk. */
export interface Lampiran {
  id: string
  namaFile: string
  kategori: string
  mime: string
  /** Create Date `DD/MM/YYYY HH:mm`. */
  tanggal: string
  /** PXCREATEOPERATOR - username pengunggah. */
  operator: string
  /** Objek penyimpanan tercatat (syarat View / View Office Online). */
  adaObjek: boolean
}

export interface LampiranKomite {
  kategori: KategoriLampiran[]
  lampiran: Lampiran[]
}

/** Lampiran klaim induk kasus komite (perintah work owner 11-10-2026), hanya-baca. */
export function ambilLampiran(id: string): Promise<LampiranKomite> {
  return minta<LampiranKomite>(`${PREFIX_KCP}/kasus/${encodeURIComponent(id)}/lampiran`)
}

function jalurLampiran(id: string, a: Lampiran): string {
  return `${PREFIX_KCP}/kasus/${encodeURIComponent(id)}/lampiran/${encodeURIComponent(a.id)}`
}

/** File: unduh satu dokumen klaim induk - fetch beridentitas. */
export function unduhLampiran(id: string, a: Lampiran): Promise<void> {
  return unduhBerkasBeridentitas(`${jalurLampiran(id, a)}/isi`, a.namaFile)
}

/** Isi satu dokumen klaim induk sebagai Blob - `View` pdf / gambar di popup penampil. */
export async function ambilIsiLampiran(id: string, a: Lampiran): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`${jalurLampiran(id, a)}/isi`), {
      method: 'GET',
      headers: { ...headerIdentitas() },
      signal: kendali.signal,
    })
    if (!jawab.ok) throw kegagalanDari(jawab.status, await jawab.text())
    return await jawab.blob()
  } finally {
    clearTimeout(jam)
  }
}

/** View Office Online - URL bertanda tangan untuk penampil kantor. */
export function tautanOfficeLampiran(id: string, a: Lampiran): Promise<{ url: string }> {
  return minta<{ url: string }>(`${jalurLampiran(id, a)}/office`)
}
