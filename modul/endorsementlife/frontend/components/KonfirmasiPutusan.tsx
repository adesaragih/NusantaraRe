// Konfirmasi keputusan endorsement (keputusan work owner 08-10-2026, pola Confirm PremiumList) - Modal inti:
// Cancel / X / Escape / backdrop menutup tanpa mengirim; Yes menjalankan `kirim` yang sama (`onLanjut`).

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { KONFIRMASI_PUTUSAN_EDM, PUTUSAN_EDM } from '../labels'

/** Panjang kutipan Comment di popup. */
export const PANJANG_KUTIPAN = 160

/** Judul + isi popup menurut Status: '1' Accept, selain itu Decline. */
export function teksKonfirmasiPutusan(status: string): { judul: string; isi: string } {
  return status === '1'
    ? { judul: KONFIRMASI_PUTUSAN_EDM.judulAccept, isi: KONFIRMASI_PUTUSAN_EDM.isiAccept }
    : { judul: KONFIRMASI_PUTUSAN_EDM.judulDecline, isi: KONFIRMASI_PUTUSAN_EDM.isiDecline }
}

/** Kutipan singkat Comment (kosong = tanpa kutipan). */
export function kutipanKomentar(komentar: string): string {
  const t = komentar.trim()
  return t.length > PANJANG_KUTIPAN ? `${t.slice(0, PANJANG_KUTIPAN)}…` : t
}

export default function KonfirmasiPutusan({
  status,
  komentar,
  onLanjut,
  onTutup,
}: {
  status: string
  komentar: string
  onLanjut: () => void
  onTutup: () => void
}) {
  const t = teksKonfirmasiPutusan(status)
  const kutipan = kutipanKomentar(komentar)
  return (
    <Modal
      judul={t.judul}
      onTutup={onTutup}
      labelBatal={KONFIRMASI_PUTUSAN_EDM.batal}
      aksi={
        <button type="button" className={status === '1' ? 'btn btn--primary' : 'btn btn--danger'} onClick={onLanjut}>
          {KONFIRMASI_PUTUSAN_EDM.ya}
        </button>
      }
    >
      <p className="edm-konfirmasi">{t.isi}</p>
      {kutipan !== '' && (
        <blockquote className="edm-konfirmasi__kutipan">
          <span>{PUTUSAN_EDM.comment}:</span> {kutipan}
        </blockquote>
      )}
    </Modal>
  )
}
