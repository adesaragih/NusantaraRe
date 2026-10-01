// Popup konfirmasi hapus Yes/Cancel - tiket 09, penyimpangan sadar 3.
//
// Pega menghapus TANPA konfirmasi (`DeleteTreatyLimit_Act`, `DeleteSecurityLife_Act`,
// `DeleteSecurityReinsurerLife_Act`, `DeleteRowBusiness`). Di sini popup lebih dulu: cacah anak yang
// IKUT terhapus dihitung server sesaat sebelum popup tampil (`GET …/dampak-hapus`); `Cancel` tidak
// mengirim apa pun; `Yes` mengirim cacah yang DILIHAT - server menolak (409) bila sudah lain.

import { HAPUS_MCRL } from '../labels'
import type { Dampak, JenisHapus } from '../api'
import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'

/** Baris rincian popup - kontrak menyebut ketiga anak, reinsurer menyebut security; daun nol baris. */
export function rincianDampak(jenis: JenisHapus, d: Dampak): string[] {
  const baris: string[] = []
  if (jenis === 'kontrak') {
    baris.push(`${d.reinsurer} ${HAPUS_MCRL.reinsurer}`, `${d.security} ${HAPUS_MCRL.security}`, `${d.business} ${HAPUS_MCRL.business}`)
  }
  if (jenis === 'reinsurer') baris.push(`${d.security} ${HAPUS_MCRL.security}`)
  return baris
}

/**
 * Penutup popup: selama `Yes` masih berjalan, `Cancel` (dan tutup popup) TIDAK berlaku - menutup
 * popup saat permintaan sudah di jalan membuat pengguna mengira nol baris tersentuh (AC 30/35).
 */
export function penutup(sibuk: boolean, tutup: () => void): () => void {
  return sibuk ? () => undefined : tutup
}

/** AC 36: baris tanpa anak tetap dikonfirmasi, dengan pesan tanpa anak. */
export function tanpaAnak(d: Dampak): boolean {
  return d.reinsurer + d.security + d.business === 0
}

export default function KonfirmasiHapus({
  judul,
  jenis,
  nama,
  dampak,
  galat,
  sibuk,
  onYa,
  onBatal,
  labelBatal,
}: {
  /** Teks tombol `Delete` panel asal. */
  judul: string
  jenis: JenisHapus
  nama: string
  dampak: Dampak | null
  galat: unknown
  sibuk: boolean
  onYa: () => void
  onBatal: () => void
  /** Teks tombol `Cancel` panel asal. */
  labelBatal: string
}) {
  return (
    <Modal
      judul={judul}
      onTutup={penutup(sibuk, onBatal)}
      labelBatal={labelBatal}
      aksi={
        <button type="button" className="btn btn--primary" disabled={sibuk || dampak === null} onClick={onYa}>
          {HAPUS_MCRL.ya}
        </button>
      }
    >
      <p>{HAPUS_MCRL.pertanyaan}</p>
      <p>
        <strong>{nama}</strong>
      </p>
      {dampak === null && galat === null && <p role="status">{HAPUS_MCRL.memuat}</p>}
      {galat !== null && <Gagal galat={galat} />}
      {dampak !== null &&
        (tanpaAnak(dampak) ? (
          <p>{HAPUS_MCRL.tanpaAnak}</p>
        ) : (
          <>
            <p>{HAPUS_MCRL.ikutTerhapus}</p>
            <ul className="mcrl-dampak">
              {rincianDampak(jenis, dampak).map((b) => (
                <li key={b}>{b}</li>
              ))}
            </ul>
          </>
        ))}
    </Modal>
  )
}
