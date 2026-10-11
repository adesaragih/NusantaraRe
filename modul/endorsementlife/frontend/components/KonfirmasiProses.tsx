// Konfirmasi sebelum Process Endorsement (keputusan work owner 08-10-2026) - Modal inti: Cancel / X / Escape /
// backdrop menutup tanpa memproses; Continue menjalankan aksi Save yang sama (`onLanjut`).

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { PROSES_EDM } from '../labels'

export default function KonfirmasiProses({
  kalimat,
  onLanjut,
  onTutup,
}: {
  kalimat: readonly string[]
  onLanjut: () => void
  onTutup: () => void
}) {
  return (
    <Modal
      judul={PROSES_EDM.judulKonfirmasi}
      onTutup={onTutup}
      labelBatal={PROSES_EDM.batal}
      aksi={
        <button type="button" className="btn btn--primary" onClick={onLanjut}>
          {PROSES_EDM.lanjut}
        </button>
      }
    >
      {kalimat.map((k) => (
        <p key={k} className="edm-konfirmasi">
          {k}
        </p>
      ))}
    </Modal>
  )
}
