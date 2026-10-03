// Popup konfirmasi hapus Ya/Batal — tiket 10 Treaty Contract Out.
//
// Penyimpangan sadar 4: Pega menghapus kontrak/reinsurer TANPA konfirmasi. Di
// sini popup menyebut jumlah baris TIAP jenis yang ikut terhapus (dihitung
// server dengan saringan yang sama dengan penghapusannya) (AC 43). Batal
// tidak mengirim apa pun; Ya mengirim jumlah yang dilihat — server menolak
// bila sudah lain.
//
// Catatan "klausul tidak terhapus" DIBUANG dari popup [keputusan work owner
// 01-10-2026]: klausul tetap tidak dihapus server (AC 44), hanya tidak lagi
// disebut di sini.

import { HAPUS_TCO } from '../labels'
import type { DampakHapusTCO } from '../api'
import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'

/** Baris rincian popup — kontrak menyebut tiga anak; reinsurer menyebut security. */
export function rincianDampak(d: DampakHapusTCO, jenis: 'kontrak' | 'reinsurer'): string[] {
  if (jenis === 'reinsurer') return [`${d.security} ${HAPUS_TCO.security}`]
  return [`${d.reinsurer} ${HAPUS_TCO.reinsurer}`, `${d.security} ${HAPUS_TCO.security}`, `${d.business} ${HAPUS_TCO.business}`]
}

/** Peringatan kontrak lain (OQ-TCO-21), tanpa angkanya — tunggal/jamak mengikuti cacah `n`. */
export function teksBersama(n: number): string {
  return `${n === 1 ? HAPUS_TCO.bersamaSatu : HAPUS_TCO.bersamaBanyak} ${HAPUS_TCO.bersamaBusiness}`
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
          {jenis === 'kontrak' && dampak.bersama > 0 && (
            <p className="alert alert--warn" role="alert">
              <strong>{dampak.bersama}</strong> {teksBersama(dampak.bersama)}
            </p>
          )}
        </>
      )}
    </Modal>
  )
}
