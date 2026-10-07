// Popup pencari tabel AGENT - dipakai `Change SOB` (tiket 33) dan `Add / Select Ceding` (tiket 34).
//
// Kedua harness Pega memakai grid atas RD yang sama, `BrowseAgentNonLife_RD` (kelas Int-AGENT): `SOB`
// (lihat di bawah) dan `CedingCompany` (section `CedingCoHierarki`, sel grid Choose → `SetDataSobCeding_Act`;
// judul kolom `ID` · `Client ID` · `Name` tertulis di section itu). Yang berbeda hanya judul jendela.
//
// Port harness `SOB` (`D:\migrasi\RNM\NB FacIn\Harness\SOB.xml`; dibuka Periode sel 52 lewat showHarness
// `InputQuotation_PreAct`, WindowName "Change SOB", Target popup): kotak `Search` (sel 294,
// `SearchSOB.CARI1`) + grid `TempBusinessSource.pxResults` kelas `ASM-FW-GISFW-Data-Agent` (kolom `.ID`,
// `.ClientName`) + tombol `Choose` per baris. Judul kolom = tangkapan layar work owner 03-10-2026 (ID · Client
// ID · Name). Data = tabel AGENT dengan syarat work owner 03-10-2026 (StatusActive=1, AgentType2 <> 'LIFE
// INSURANCE', ClientID terisi; cari tidak peka huruf) - activity pencari Pega
// `SearchHierarkiSourceBizAgentTreatyIn_Act` TIDAK ada di korpus.
//
// Choose di Pega: activity tadi → `opener.location.reload` → `window.close`. Di sini (keputusan agent E-1)
// Choose mengisi Source of business di layar pembuka dan menutup popup; tersimpan lewat Save for later.
// Pencarian berjalan SAAT MENGETIK (jeda `JEDA_CARI_MS` sesudah ketikan terakhir) dan juga lewat Enter - SOB.xml
// dan tangkapan layar Pega tidak memuat tombol Search; daftar Pega menyaring dari kotak itu (laporan work owner
// 03-10-2026: "tombol search … belum berfungsi" saat hanya Enter yang mencari).

import { useEffect, useRef, useState } from 'react'

import { Field, Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariSOB, type BarisSOB, type HalamanSOB } from '../api'
import { POPUP_SOB, TEKS_INWARD } from '../labels'

/** Jeda sesudah ketikan terakhir sebelum mencari - satu permintaan per jeda ketik, bukan per huruf. */
export const JEDA_CARI_MS = 400

export default function PopupPilihAgent({
  judul,
  onTutup,
  onPilih,
}: {
  /** `POPUP_SOB.judul` atau `POPUP_CEDING.judulCari`. */
  judul: string
  onTutup: () => void
  onPilih: (b: BarisSOB) => void
}) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<HalamanSOB | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  // Nomor permintaan terakhir: jawaban lama dibuang (pola popup ChooseAccount).
  const nomorPermintaan = useRef(0)
  const kunciCari = useRef('')

  async function muat(cari: string, halaman: number) {
    const nomor = ++nomorPermintaan.current
    kunciCari.current = cari
    setMemuat(true)
    setGalat(null)
    try {
      const h = await cariSOB(cari, halaman)
      if (nomor === nomorPermintaan.current) setHasil(h)
    } catch (err) {
      if (nomor === nomorPermintaan.current) {
        setHasil(null)
        setGalat(err)
      }
    } finally {
      if (nomor === nomorPermintaan.current) setMemuat(false)
    }
  }

  // Dibuka: langsung memuat (kotak kosong = kata cari terakhir, jeda 0). Mengetik: memuat ulang dari halaman 1
  // sesudah jeda; ketikan baru membatalkan jadwal yang lama.
  useEffect(() => {
    const jadwal = window.setTimeout(() => void muat(kotak, 1), kotak === kunciCari.current ? 0 : JEDA_CARI_MS)
    return () => window.clearTimeout(jadwal)
  }, [kotak])

  return (
    <Modal judul={judul} onTutup={onTutup} onKirim={() => void muat(kotak, 1)} lebar>
      <div className="nbf-popup__kotak">
        <Field label={POPUP_SOB.cari.label} value={kotak} onChange={setKotak} />
      </div>
      {hasil && hasil.baris.length > 0 && (
        <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(kunciCari.current, h)} />
      )}
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {POPUP_SOB.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
          {hasil && hasil.baris.length > 0 && (
            <tbody>
              {hasil.baris.map((b) => (
                <tr key={b.id}>
                  <td>{b.id}</td>
                  <td>{b.clientId}</td>
                  <td>{b.name}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {POPUP_SOB.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {memuat && !hasil && <Memuat />}
      {hasil && hasil.baris.length === 0 && <Kosong pesan={TEKS_INWARD.tanpaSob} />}
      <Gagal galat={galat} />
    </Modal>
  )
}
