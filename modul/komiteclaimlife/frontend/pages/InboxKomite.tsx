// Inbox Komite — tiket 01 Komite Claim Life.
//
// Worklist `KomiteRouter` (`Flow/KomiteLife_Flow.xml` b619 `pyRouteTo
// Custom`): kasus yang menunggu PELAKU — anggota pertama di tangga yang belum
// memutuskan. Penyaringnya di server; layar tidak menyaring ulang.
//
// ⛔ UANG TETAP TEKS; status baris KATA (AC 17, AC 29 spec).

import { useCallback, useEffect, useState } from 'react'

import { PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import { EFEK_KOMITE, INBOX_KOMITE, KOLOM_INBOX_KOMITE } from '../labels'
import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilInboxKomite,
  ambilLaporanHarianKomite,
  type BarisInboxKomite,
  type HalamanInboxKomite,
  type LaporanHarianKomite,
} from '../api'

/** Sel kosong ditandai (ADR-U-0027). */
export function selKomite(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/** Tingkat berjalan dari seluruh tingkat — `KomiteCount`/`KomiteLoop`. */
export function tingkatKomite(b: BarisInboxKomite): string {
  return `${String(b.tingkatBerjalan)} / ${String(b.komiteLoop)}`
}

export default function InboxKomite({
  onBuka,
  peran = [],
}: {
  onBuka: (kasusID: string) => void
  /** Peran pelaku — laporan "perlu intervensi" hanya untuk admin (tiket 08). */
  peran?: readonly KodePeran[]
}) {
  const [hal, setHal] = useState<HalamanInboxKomite | null>(null)
  const [laporan, setLaporan] = useState<LaporanHarianKomite | null>(null)
  const admin = peran.includes(PERAN.admin)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      setHal(await ambilInboxKomite())
      // ⛔ Laporan kosong DINYATAKAN — ketiadaan laporan ≠ ketiadaan masalah.
      if (admin) setLaporan(await ambilLaporanHarianKomite())
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
      setGalat(e)
      setHal(null)
    } finally {
      setSibuk(false)
    }
  }, [admin])

  useEffect(() => {
    void muat()
  }, [muat])

  return (
    <section className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{INBOX_KOMITE.judul}</h2>
        <p className="inbox__aturan">{INBOX_KOMITE.aturan}</p>
      </header>
      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && hal.baris.length === 0 && <Kosong pesan={INBOX_KOMITE.kosong} />}
      {hal !== null && hal.baris.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{KOLOM_INBOX_KOMITE.kasusId}</th>
              <th>{KOLOM_INBOX_KOMITE.tglUpdate}</th>
              <th>{KOLOM_INBOX_KOMITE.statusWork}</th>
              <th>{KOLOM_INBOX_KOMITE.nomorKlaim}</th>
              <th>{KOLOM_INBOX_KOMITE.tingkat}</th>
              <th>{KOLOM_INBOX_KOMITE.nilaiKlaim}</th>
              <th>{KOLOM_INBOX_KOMITE.mataUang}</th>
              <th>{KOLOM_INBOX_KOMITE.statusBaris}</th>
            </tr>
          </thead>
          <tbody>
            {hal.baris.map((b) => (
              <tr
                key={b.kasusId}
                className="inbox__baris"
                onClick={() => {
                  onBuka(b.kasusId)
                }}
              >
                <td>{b.kasusId}</td>
                <td>{selKomite(b.tglUpdate)}</td>
                <td>{selKomite(b.statusWork)}</td>
                <td>{selKomite(b.nomorKlaim)}</td>
                <td>{tingkatKomite(b)}</td>
                <td>{selKomite(b.nilaiKlaim)}</td>
                <td>{selKomite(b.mataUang)}</td>
                <td>{b.statusBaris}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {hal !== null && (
        <p className="inbox__cacah" role="status">
          {hal.baris.length} dari {hal.total}
        </p>
      )}
      {laporan !== null && (
        <section className="inbox__intervensi">
          <h3>
            {EFEK_KOMITE.laporan} ({laporan.tanggal})
          </h3>
          {laporan.kosong ? (
            <p role="status">{EFEK_KOMITE.laporanKosong}</p>
          ) : (
            <ul>
              {laporan.baris.map((b, i) => (
                <li key={`${b.kasusId}-${b.jenis}-${String(i)}`}>
                  <button
                    type="button"
                    onClick={() => {
                      onBuka(b.kasusId)
                    }}
                  >
                    {b.kasusId}
                  </button>{' '}
                  {b.jenis} — {b.keadaan} sejak {b.sejak}
                </li>
              ))}
            </ul>
          )}
        </section>
      )}
    </section>
  )
}
