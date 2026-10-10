// Tabel komite di bawah inbox Claim Life (perintah work owner 09-10-2026: menu Komite Claim Life dihapus, "anggap menu
// itu tidak pernah ada" - komite digabung ke menu Claim Life). Isinya Inbox Komite modul Komite Claim Life yang dipindah
// ke sini (disalin, bukan diimpor): kasus yang menunggu PELAKU (worklist `KomiteRouter`, penyaringnya di server) lewat
// rute pinjaman `GET /api/komite`, plus laporan harian "perlu intervensi" untuk admin (tiket 08). Klik baris = layar
// kasus komite DI TEMPAT (`onBukaModul`); Kembali ke inbox ini. `[tidak ada di korpus]`: bagian ini tampil hanya bila
// ada kasus yang menunggu, ada laporan admin, atau pembacaannya gagal - pelaku tanpa kasus komite tidak melihatnya.
//
// ⛔ UANG TETAP TEKS; status baris KATA.

import { useEffect, useState } from 'react'

import { PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import {
  daftarKomite,
  laporanHarianKomite,
  MODUL_KOMITE,
  type BarisKomite,
  type HalamanKomite,
  type LaporanHarianKomite,
} from '../api'
import { TABEL_KOMITE } from '../labels'

/** Sel kosong ditandai (ADR-U-0027). */
export function selKomite(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/** Tingkat berjalan dari seluruh tingkat — `KomiteCount`/`KomiteLoop`. */
export function tingkatKomite(b: BarisKomite): string {
  return `${String(b.tingkatBerjalan)} / ${String(b.komiteLoop)}`
}

export default function TabelKomite({
  peran,
  onBukaModul,
}: {
  /** Peran pelaku — laporan "perlu intervensi" hanya untuk admin (tiket 08). */
  peran: readonly KodePeran[]
  /** `PropsRute.onBukaModul` - baris tabel: layar kasus komite dibuka di tempat. */
  onBukaModul?: (modul: string, id: string) => boolean
}) {
  const [hal, setHal] = useState<HalamanKomite | null>(null)
  const [laporan, setLaporan] = useState<LaporanHarianKomite | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const admin = peran.includes(PERAN.admin)

  useEffect(() => {
    let aktif = true
    void (async () => {
      try {
        const h = await daftarKomite()
        // ⛔ Laporan kosong DINYATAKAN — ketiadaan laporan ≠ ketiadaan masalah.
        const l = admin ? await laporanHarianKomite() : null
        if (!aktif) return
        setHal(h)
        setLaporan(l)
      } catch (e) {
        // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
        if (aktif) setGalat(e)
      }
    })()
    return () => {
      aktif = false
    }
  }, [admin])

  const buka = (id: string) => {
    setPesan(null)
    if (onBukaModul?.(MODUL_KOMITE, id) !== true) setPesan(TABEL_KOMITE.takTerpasang)
  }

  if (galat === null && (hal === null || hal.baris.length === 0) && laporan === null) return null

  return (
    <section className="inbox claimlife__komite" aria-labelledby="claimlife-judul-komite">
      <h3 id="claimlife-judul-komite" className="inbox__judul">
        {TABEL_KOMITE.judul}
      </h3>
      {pesan && (
        <div className="alert alert--error" role="alert">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && hal.baris.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{TABEL_KOMITE.kasusId}</th>
              <th>{TABEL_KOMITE.tglUpdate}</th>
              <th>{TABEL_KOMITE.statusWork}</th>
              <th>{TABEL_KOMITE.nomorKlaim}</th>
              <th>{TABEL_KOMITE.tingkat}</th>
              <th>{TABEL_KOMITE.nilaiKlaim}</th>
              <th>{TABEL_KOMITE.mataUang}</th>
              <th>{TABEL_KOMITE.statusBaris}</th>
            </tr>
          </thead>
          <tbody>
            {hal.baris.map((b) => (
              <tr key={b.kasusId} className="inbox__baris" onClick={() => buka(b.kasusId)}>
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
      {hal !== null && hal.baris.length > 0 && (
        <p className="inbox__cacah" role="status">
          {hal.baris.length} dari {hal.total}
        </p>
      )}
      {laporan !== null && (
        <section className="inbox__intervensi">
          <h3>
            {TABEL_KOMITE.laporan} ({laporan.tanggal})
          </h3>
          {laporan.kosong ? (
            <p role="status">{TABEL_KOMITE.laporanKosong}</p>
          ) : (
            <ul>
              {laporan.baris.map((b, i) => (
                <li key={`${b.kasusId}-${b.jenis}-${String(i)}`}>
                  <button type="button" onClick={() => buka(b.kasusId)}>
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
