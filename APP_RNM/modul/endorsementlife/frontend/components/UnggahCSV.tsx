// Modal `Upload CSV` b8973 → FlowAction `UploadCSV_LifeEndorsement` (judul b191): pilih berkas, lalu
// server mengurai dan memvalidasinya TANPA menyimpan (`POST /unggah`) - tinjauan sebelum `Add CSV Data`.

import { useState, type ChangeEvent } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { periksaCSV, type PeriksaCSVEDM, type PesanCSVEDM } from '../api'
import { UMUM_EDM, UNGGAH_EDM } from '../labels'

/** Daftar penolakan berbaris dan berkolom (AC 37). */
export function TabelPesanCSV({ hasil }: { hasil: PeriksaCSVEDM }) {
  return (
    <div className="edm-csv">
      <dl className="edm-csv__ringkas">
        <dt>{UMUM_EDM.barisDibaca}</dt>
        <dd>{hasil.total}</dd>
        <dt>{UMUM_EDM.barisDitolak}</dt>
        <dd>{hasil.ditolak}</dd>
        {hasil.diabaikan.length > 0 && (
          <>
            <dt>{UMUM_EDM.kolomDiabaikan}</dt>
            <dd>{hasil.diabaikan.join(', ')}</dd>
          </>
        )}
      </dl>
      {hasil.pesan.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel" role="alert">
            <thead>
              <tr>
                <th>{UMUM_EDM.kolomBaris}</th>
                <th>{UMUM_EDM.kolomKolom}</th>
                <th>{UMUM_EDM.kolomPesan}</th>
              </tr>
            </thead>
            <tbody>
              {hasil.pesan.map((p: PesanCSVEDM, i) => (
                <tr key={`${p.baris}-${p.kolom}-${i}`}>
                  <td className="edm-angka">{p.baris}</td>
                  <td>{p.kolom}</td>
                  <td>{p.pesan}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {hasil.terpotong && <p className="edm-catatan">{UMUM_EDM.pesanTerpotong}</p>}
    </div>
  )
}

export default function UnggahCSV({
  kasusId,
  onBerkas,
  onTutup,
}: {
  kasusId: string
  onBerkas: (berkas: File, hasil: PeriksaCSVEDM) => void
  onTutup: () => void
}) {
  const [hasil, setHasil] = useState<PeriksaCSVEDM | null>(null)
  const [memeriksa, setMemeriksa] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)

  const pilih = async (e: ChangeEvent<HTMLInputElement>) => {
    const berkas = e.target.files?.[0]
    if (berkas === undefined) return
    setMemeriksa(true)
    setGalat(null)
    setHasil(null)
    try {
      const h = await periksaCSV(kasusId, berkas)
      setHasil(h)
      onBerkas(berkas, h)
    } catch (err) {
      setGalat(err)
    } finally {
      setMemeriksa(false)
    }
  }

  return (
    <Modal judul={UNGGAH_EDM.judulModal} onTutup={onTutup} labelBatal={UMUM_EDM.tutup} lebar>
      <label className="edm-csv__berkas">
        <span>{UMUM_EDM.pilihBerkas}</span>
        <input type="file" accept=".csv,text/csv" disabled={memeriksa} onChange={(e) => void pilih(e)} />
      </label>
      {memeriksa && <p className="edm-catatan">{UMUM_EDM.mengunggah}</p>}
      {galat !== null && <Gagal galat={galat} />}
      {hasil !== null && <TabelPesanCSV hasil={hasil} />}
    </Modal>
  )
}
