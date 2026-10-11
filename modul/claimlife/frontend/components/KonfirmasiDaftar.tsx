// Konfirmasi sebelum Register Claim (keputusan work owner 10-10-2026: "saat register claim di klik, ada konfirmasi
// dahulu") - Modal inti, pola KonfirmasiProses endorsementlife. Cancel / X / Escape / backdrop menutup TANPA mendaftar
// (pilihan tetap); "Register Claim" menjalankan pendaftaran (`onLanjut`). Selama proses tombolnya nonaktif berlabel
// "Registering…" dan modal tidak dapat ditutup - satu klik, satu permintaan.

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { REGISTER_BANYAK } from '../labels'
import type { KelompokPilihan } from '../pilihanPeserta'

export default function KonfirmasiDaftar({
  kelompok,
  sibuk,
  onLanjut,
  onTutup,
}: {
  kelompok: readonly KelompokPilihan[]
  sibuk: boolean
  onLanjut: () => void
  onTutup: () => void
}) {
  return (
    <Modal
      judul={REGISTER_BANYAK.judulKonfirmasi}
      onTutup={onTutup}
      tanpaTutup={sibuk}
      labelBatal={REGISTER_BANYAK.batal}
      aksi={
        <button type="button" className="btn btn--primary" onClick={onLanjut} disabled={sibuk}>
          {sibuk ? REGISTER_BANYAK.mendaftarkan : REGISTER_BANYAK.daftarkan}
        </button>
      }
    >
      <p>{REGISTER_BANYAK.kalimatKonfirmasi.split('{n}').join(String(kelompok.length))}</p>
      <div className="register__gulir">
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{REGISTER_BANYAK.kolomPL}</th>
              <th>{REGISTER_BANYAK.kolomUang}</th>
              <th>{REGISTER_BANYAK.kolomJumlah}</th>
            </tr>
          </thead>
          <tbody>
            {kelompok.map((k) => (
              <tr key={`${k.nomorPremiList}|${k.mataUang}`}>
                <td>{k.nomorPremiList}</td>
                <td>{k.mataUang}</td>
                <td>{k.peserta.length}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p>{REGISTER_BANYAK.lanjutkan}</p>
    </Modal>
  )
}
