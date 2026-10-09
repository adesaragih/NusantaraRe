// Ringkasan salinan lampiran master di pesan tombol Save/Submit draf.
//
// Pega (`TreatyRevisionCopyAttachment`) menelan berkas yang gagal; server
// aplikasi ini melaporkan nasib tiap berkas (`HasilSimpan.SalinanLampiran`),
// dan layar menuliskannya di bawah pesan prosedur. Tanpa salinan (penyesuaian
// tersimpan, atau master tanpa lampiran) pesan prosedur tidak berubah.

import type { SalinanLampiran } from '../api'
import { SALINAN_LAMPIRAN } from '../labelsSalinanLampiran'

/** Kalimat ringkasan salinan; '' bila tidak ada yang disalin. */
export function pesanSalinanLampiran(s: SalinanLampiran | null | undefined): string {
  if (s === undefined || s === null) return ''
  if (s.pesan !== undefined && s.pesan !== '') return s.pesan
  const berkas = s.berkas ?? []
  let t = SALINAN_LAMPIRAN.ringkas
    .replace('{sumber}', s.sumber)
    .replace('{tersalin}', String(s.tersalin))
    .replace('{jumlah}', String(berkas.length))
  const gagal = berkas.filter((b) => !b.berhasil)
  if (gagal.length > 0) {
    t += ' ' + SALINAN_LAMPIRAN.gagal + ' ' + gagal.map((b) => `${b.nama} (${b.pesan})`).join('; ')
  }
  return t
}

/** Pesan prosedur ditambah ringkasan salinan (bila ada). */
export function gabungPesanSalinan(pesan: string, s: SalinanLampiran | null | undefined): string {
  const tambahan = pesanSalinanLampiran(s)
  return tambahan === '' ? pesan : `${pesan} ${tambahan}`
}
