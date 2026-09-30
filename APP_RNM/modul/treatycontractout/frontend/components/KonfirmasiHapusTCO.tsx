// Popup konfirmasi hapus Ya/Batal — tiket 10 Treaty Contract Out.
//
// Penyimpangan sadar 4: Pega menghapus kontrak/reinsurer TANPA konfirmasi. Di
// sini popup menyebut jumlah baris TIAP jenis yang ikut terhapus (dihitung
// server dengan saringan yang sama dengan penghapusannya) dan menyatakan
// EKSPLISIT bahwa klausul tidak terhapus (AC 43, 44). Batal tidak mengirim
// apa pun; Ya mengirim jumlah yang dilihat — server menolak bila sudah lain.

import { HAPUS_TCO } from '../labels'
import type { DampakHapusTCO } from '../api'
import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'

/** Baris rincian popup — kontrak menyebut tiga anak + klausul; reinsurer menyebut security. */
export function rincianDampak(d: DampakHapusTCO, jenis: 'kontrak' | 'reinsurer'): string[] {
  if (jenis === 'reinsurer') return [`${d.security} ${HAPUS_TCO.security}`]
  return [`${d.reinsurer} ${HAPUS_TCO.reinsurer}`, `${d.security} ${HAPUS_TCO.security}`, `${d.business} ${HAPUS_TCO.business}`]
}

export default function KonfirmasiHapusTCO({
  jenis,
  nama,
  dampak,
  galat,
  sibuk,
  onYa,
  onBatal,
}: {
  jenis: 'kontrak' | 'reinsurer'
  nama: string
  dampak: DampakHapusTCO | null
  galat: unknown
  sibuk: boolean
  onYa: () => void
  onBatal: () => void
}) {
  return (
    <Modal
      judul={jenis === 'kontrak' ? HAPUS_TCO.judulKontrak : HAPUS_TCO.judulReinsurer}
      onTutup={onBatal}
      labelBatal={HAPUS_TCO.batal}
      aksi={
        <button type="button" className="btn btn--primary" disabled={sibuk || dampak === null} onClick={onYa}>
          {HAPUS_TCO.ya}
        </button>
      }
    >
      <p>
        <strong>{nama}</strong>
      </p>
      {dampak === null && galat === null && <p role="status">{HAPUS_TCO.memuatDampak}</p>}
      {galat !== null && <Gagal galat={galat} />}
      {dampak !== null && (
        <>
          <p>{HAPUS_TCO.ikutTerhapus}</p>
          <ul>
            {rincianDampak(dampak, jenis).map((b) => (
              <li key={b}>{b}</li>
            ))}
          </ul>
          {jenis === 'kontrak' && (
            <p className="polis__catatan" role="note">
              {dampak.klausulTetap} {HAPUS_TCO.klausulTetap}
            </p>
          )}
          {jenis === 'kontrak' && dampak.bersama > 0 && (
            <p className="alert alert--warn" role="alert">
              <strong>{dampak.bersama}</strong> {HAPUS_TCO.bersama}
            </p>
          )}
        </>
      )}
    </Modal>
  )
}
