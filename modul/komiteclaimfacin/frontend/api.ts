// Klien Komite Claim Fac In - `/api/komite-claim-fac-in` (`modul/komiteclaimfacin/backend/handlers/rute.go`). Layar
// datang dari server sudah tersusun (`models.Layar`): ubin, bagian, label VERBATIM Section ShowTransfer (+ rincian
// `SpreadingDetail` / `DetailAdjustmentFac`, pop-up `ShowRetro`), nilai klaim induk; layar hanya merender dan mengirim
// isian keputusan. Pola `modul/komiteclaimnonprop/frontend/api.ts` (disalin, bukan impor). Daftar kerja penyetuju
// (`GET /kasus`) dibaca tabel komite inbox Claim Fac In (`modul/claimfacin/frontend/api.ts`), bukan modul ini.

import { BATAS_WAKTU_MS, kegagalanDari, minta, rakitURL } from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

export const PREFIX_KCFI = '/api/komite-claim-fac-in'

/** `tautan` = pxLink: nilai sel = kunci `Layar.modal` (kosong = tanpa tautan). */
export type JenisNilai = 'teks' | 'teksPanjang' | 'angka' | 'tanggal' | 'tanggalJam' | 'centang' | 'tautan'

/** Satu sel hanya-baca berlabel - `models.Medan`. */
export interface Medan {
  label: string
  nilai: string
  jenis: JenisNilai | ''
}

export interface KolomGrid {
  label: string
  properti: string
  jenis: JenisNilai | ''
}

/** Grid hanya-baca - `models.Grid`. `rincian[i]` = expand pane baris ke-i (masterDetail); `null` = tidak dapat dibuka. */
export interface Grid {
  judul: string
  kolom: KolomGrid[]
  baris: Record<string, string>[] | null
  rincian?: (Bagian[] | null)[]
  kaki?: Medan[]
}

export interface Bagian {
  kunci: string
  judul?: string
  medan?: Medan[]
  grid?: Grid[]
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
  usulTutup: boolean
  usulCadang: boolean
}

export interface IsianLayar {
  nilai: Keputusan
  /** Dua Propose tampil (`VIS .TransferType = 2`). */
  tampilUsul: boolean
  /** Dua Propose terbuka (`NA .KomiteCount != '1'` salah). */
  terbuka: boolean
  pilihanTerima: Pilihan[]
  label: Partial<Record<keyof Keputusan, string>>
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
  /** 2 ADJUSTMENT, 3 REJECT, 4 CLOSE. */
  transferType: string
  komiteLoop: number
  komiteCount: number
  acceptStatus: string
  usulTutup: string
  usulCadang: string
  tahap: string
  statusWork: string
  pembuatId: string
  pembuatNama: string
  tglCreate: string
  tglUpdate: string
  tangga: Anggota[] | null
}

export interface Layar {
  kasus: Kasus
  judul: string[]
  ubin: Medan[]
  bagian: Bagian[]
  isian: IsianLayar
  tombol: Tombol[]
  /** Isi pop-up local action ShowRetro per kunci sel `tautan`. */
  modal?: Record<string, Bagian[]>
  bolehKerja: boolean
  pesan?: string[]
}

export interface HasilKeputusan {
  kasusId: string
  selesai: boolean
  komiteCount: number
  acceptedNo?: string
  /** Pesan sesudah Submit (mis. OQ dokumen PDF) - tampil dulu sebelum kembali ke inbox Claim Fac In. */
  info?: string
}

export function bukaKasus(id: string): Promise<Layar> {
  return minta<Layar>(`${PREFIX_KCFI}/kasus/${encodeURIComponent(id)}`)
}

/** Validasi layar gagal di server (422 `{galat, pesan}`, `services.GalatValidasi`): Submit dibatalkan. */
export class GalatValidasiKomite extends Error {
  constructor(readonly pesan: string[]) {
    super(pesan.join('; '))
    this.name = 'GalatValidasiKomite'
  }
}

/** Jawaban gagal Submit: 422 ber-`pesan` = GalatValidasiKomite; selainnya (403 / 409 / ...) amplop galat inti. */
export function galatPutuskan(status: number, teks: string): Error {
  if (status === 422) {
    try {
      const o = JSON.parse(teks) as { pesan?: unknown }
      if (Array.isArray(o.pesan)) {
        return new GalatValidasiKomite(o.pesan.filter((p): p is string => typeof p === 'string'))
      }
    } catch {
      // bukan JSON: amplop galat inti yang menamainya
    }
  }
  return kegagalanDari(status, teks)
}

/**
 * Submit. Dibaca langsung (bukan `minta`) karena jawaban 422 membawa `pesan` yang dibuang amplop galat inti (pola
 * `aksiKasus` Claim Fac In).
 */
export async function putuskan(id: string, kep: Keputusan): Promise<HasilKeputusan> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  let jawab: Response
  try {
    jawab = await fetch(rakitURL(`${PREFIX_KCFI}/kasus/${encodeURIComponent(id)}/putuskan`), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...headerIdentitas() },
      body: JSON.stringify(kep),
      signal: kendali.signal,
    })
  } finally {
    clearTimeout(jam)
  }
  const teks = await jawab.text()
  if (!jawab.ok) throw galatPutuskan(jawab.status, teks)
  try {
    return JSON.parse(teks) as HasilKeputusan
  } catch {
    throw kegagalanDari(jawab.status, teks)
  }
}
