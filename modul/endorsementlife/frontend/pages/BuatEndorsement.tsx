// Buat endorsement - `Harness/EndorsmentLife_harnes.xml` → `Section/EndorsmentLife_Section.xml`
// (dibuka `Create Addendum` `InboxEndorsementLife.xml` b6620 → `showHarness` b6984).
//
// `SetErrorBatalEndorsement_Act` berjalan saat `Policy No` (b1249), `EDM Type` (b1545), dan `EDM Date` (b3347)
// berubah - di sini `cekKelayakan` sesudah isian itu berubah. `Submit` b4226 tampil hanya bila nol pesan
// gerbang dan `EDM Type` `1`/`3`; ia menjalankan `MappingEDMLife` lalu `openAssignment` (buka layar kasus).
//
// ⏸️ `EDM Type Perubahan Data` (b1645) dan `EDM Type Batal` (b1937) tidak dibangun: opsinya tidak ada di korpus
// dan tidak punya kolom (OQ-EDM-005). ➖ `Process Policy No` b2414 `VIS=never`.

import { useEffect, useState } from 'react'

import { Area, Field, FieldTanggal, Gagal, IkonTutup, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import { buatKasus, cekKelayakan, type KelayakanEDM } from '../api'
import { BUAT_EDM, INBOX_EDM, OPSI_EDM_TYPE, UMUM_EDM } from '../labels'
import '../endorsementlife.css'

/** Opsi `EDM Type` - `1` Perubahan Data, `3` Batal (spec §5). */
const OPSI = Object.entries(OPSI_EDM_TYPE).map(([value, label]) => ({ value, label }))

export default function BuatEndorsement({ onTutup, onDibuat }: { onTutup: () => void; onDibuat: (id: string) => void }) {
  const [policyNo, setPolicyNo] = useState('')
  const [edmType, setEdmType] = useState('')
  const [edmDate, setEdmDate] = useState('')
  const [deskripsi, setDeskripsi] = useState('')
  const [kelayakan, setKelayakan] = useState<KelayakanEDM | null>(null)
  const [memeriksa, setMemeriksa] = useState(false)
  const [mengirim, setMengirim] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)

  // Gerbang dijalankan ulang setiap kali salah satu dari ketiga isian pemicu berubah (b1249, b1545, b3347).
  useEffect(() => {
    if (policyNo === '' && edmType === '') {
      setKelayakan(null)
      return
    }
    let batal = false
    setMemeriksa(true)
    const t = setTimeout(() => {
      cekKelayakan(policyNo.trim(), edmType)
        .then((k) => {
          if (!batal) setKelayakan(k)
        })
        .catch((e: unknown) => {
          if (!batal) setGalat(e)
        })
        .finally(() => {
          if (!batal) setMemeriksa(false)
        })
    }, 300)
    return () => {
      batal = true
      clearTimeout(t)
    }
  }, [policyNo, edmType, edmDate])

  const bolehSubmit = kelayakan !== null && kelayakan.boleh && !memeriksa && !mengirim

  async function kirim() {
    setGalat(null)
    setMengirim(true)
    try {
      const h = await buatKasus({ policyNo: policyNo.trim(), edmType, edmDate: keInputTanggal(edmDate), description: deskripsi })
      onDibuat(h.id)
    } catch (e) {
      setGalat(e)
    } finally {
      setMengirim(false)
    }
  }

  return (
    <section className="panel edm-buat">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{INBOX_EDM.createAddendum}</h2>
        <button type="button" className="btn btn--ghost btn--sm" aria-label={UMUM_EDM.tutup} title={UMUM_EDM.tutup} onClick={onTutup}>
          <IkonTutup />
        </button>
      </header>
      <div className="edm-buat__isian">
        <Field label={BUAT_EDM.policyNo} value={policyNo} onChange={setPolicyNo} autoFocus />
        <Pilih label={BUAT_EDM.edmType} value={edmType} onChange={setEdmType} opsi={OPSI} />
        <Area label={BUAT_EDM.description} value={deskripsi} onChange={setDeskripsi} />
        <FieldTanggal label={BUAT_EDM.edmDate} value={edmDate} onChange={setEdmDate} />
      </div>
      {memeriksa && <p className="edm-catatan" role="status">{UMUM_EDM.memeriksa}</p>}
      {kelayakan !== null && kelayakan.pesan.length > 0 && (
        <ul className="edm-pesan" role="alert">
          {kelayakan.pesan.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {bolehSubmit && (
        <div className="edm-aksi">
          <button type="button" className="btn btn--primary" onClick={() => void kirim()}>
            {BUAT_EDM.submit}
          </button>
        </div>
      )}
    </section>
  )
}
