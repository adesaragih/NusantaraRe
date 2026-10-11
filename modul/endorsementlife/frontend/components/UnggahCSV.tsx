// Modal `Upload CSV` b8973 → FlowAction `UploadCSV_LifeEndorsement` (judul b191): pilih berkas, lalu
// server mengurai dan memvalidasinya TANPA menyimpan (`POST /unggah`) - tinjauan sebelum `Add CSV Data`.

import { useState, type ChangeEvent } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { kodeStatusGalat } from '../../../../inti/frontend/klien'
import { periksaCSV, unduhTemplatCSV, type PeriksaCSVEDM, type PesanCSVEDM } from '../api'
import { UMUM_EDM, UNGGAH_EDM } from '../labels'

/** Kalimat ringkasan tinjauan (disederhanakan 10-10-2026): ada masalah → merah, semua lolos → hijau. */
export function ringkasanCSV(h: PeriksaCSVEDM): { kalimat: string; masalah: boolean } {
  if (h.ditolak > 0) {
    return {
      kalimat: UNGGAH_EDM.ringkasMasalah.split('{n}').join(String(h.ditolak)).split('{total}').join(String(h.total)),
      masalah: true,
    }
  }
  return { kalimat: UNGGAH_EDM.ringkasLolos.split('{total}').join(String(h.total)), masalah: false }
}

/**
 * Tinjauan unggahan (AC 37) untuk user awam (10-10-2026): satu kalimat ringkasan berwarna, kolom asing (abu, bila ada),
 * lalu tabel dua kolom Row / Problem - kalimat masalah dari server sudah menyebut nama kolomnya. Row = nomor baris
 * spreadsheet (judul = 1).
 */
export function TabelPesanCSV({ hasil }: { hasil: PeriksaCSVEDM }) {
  const r = ringkasanCSV(hasil)
  return (
    <div className="edm-csv">
      <p
        className={r.masalah ? 'edm-csv__ringkasan edm-csv__ringkasan--masalah' : 'edm-csv__ringkasan edm-csv__ringkasan--lolos'}
        role={r.masalah ? 'alert' : 'status'}
      >
        {r.kalimat}
      </p>
      {hasil.diabaikan.length > 0 && (
        <p className="edm-catatan">{UNGGAH_EDM.kolomAsing.split('{kolom}').join(hasil.diabaikan.join(', '))}</p>
      )}
      {hasil.pesan.length > 0 && (
        <>
          <p className="edm-catatan">{UNGGAH_EDM.catatanBaris}</p>
          <div className="edm-gulir">
            <table className="inbox__tabel">
              <thead>
                <tr>
                  <th>{UMUM_EDM.kolomBaris}</th>
                  <th>{UNGGAH_EDM.kolomMasalah}</th>
                </tr>
              </thead>
              <tbody>
                {hasil.pesan.map((p: PesanCSVEDM, i) => (
                  <tr key={`${p.baris}-${p.kolom}-${i}`}>
                    <td className="edm-angka">{p.baris}</td>
                    <td>{p.pesan}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
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
      {/* Berkas/judul ditolak (422, kalimat server tanpa awalan teknis): arahkan ke templat (10-10-2026). */}
      {galat !== null && kodeStatusGalat(galat) === 422 && (
        <p className="edm-catatan edm-csv__petunjuk">
          {UNGGAH_EDM.petunjukTemplat}{' '}
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              unduhTemplatCSV(kasusId).catch((e: unknown) => setGalat(e))
            }}
          >
            {UNGGAH_EDM.downloadTemplate}
          </button>
        </p>
      )}
      {hasil !== null && <TabelPesanCSV hasil={hasil} />}
    </Modal>
  )
}
