// Layar buat endorsemen - `Harness/TreatyCreateEdm` + `Section/TreatyCreateEdm` (tombol portal "Create New Addendum
// Treaty": showHarness target newDocument, judul "Create EDM"). Di React: popup `Modal` (pola popup NB Treaty In).
//
//   S2 "Errors" (`TrtERR.CARI1 != ''`)  teks galat, hanya-baca
//   No Polis Treaty   `DisplayData.CARI1` pxTextInput - change -> TrtEdmCheckPolicyError -> CheckNopolisAvailability
//                     (`GET /periksa-polis`)
//   No Master Treaty  `DisplayData.CARI2` disabled selalu - diisi jawaban periksa
//   Source of Change  `DisplayData.CARI3` pxDropdown (DT `TreatyEDMListType` -> acuan `jenisEdm`), tidak wajib
//   Create            tampil `TrtERR.CARI1 = '' && DisplayData.CARI1 != '' && DisplayData.CARI2 != ''` ->
//                     CreateEDMT(edmtype) -> openWorkByHandle (`onBuka`)
//
// ⛔ Sel tanggal tanpa label (`1=2`) dan tombol tanpa label (`NEVER`) tidak dirender. Source of Change TIDAK
// disyaratkan tombol Create (XML) - tidak diwajibkan di sini.

import { useRef, useState } from 'react'

import { Field, Gagal, Modal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { buatKasus, periksaPolis } from '../api'
import { BUAT, JUDUL, TOMBOL } from '../labels'

/** Syarat tampil tombol Create (`TreatyCreateEdm` S8 pyVisible). */
export function tampilTombolBuat(galat: string, noPolis: string, noMaster: string): boolean {
  return galat === '' && noPolis !== '' && noMaster !== ''
}

export default function BuatEDM({
  opsiJenis,
  onBuka,
  onTutup,
}: {
  /** Pilihan "Source of Change" (`Acuan.jenisEdm`). */
  opsiJenis: Opsi[]
  onBuka: (id: string) => void
  onTutup: () => void
}) {
  const [noPolis, setNoPolis] = useState('')
  const [diperiksa, setDiperiksa] = useState('')
  const [noMaster, setNoMaster] = useState('')
  const [galatTrt, setGalatTrt] = useState('')
  const [jenis, setJenis] = useState('')
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  // Nomor yang pemeriksaannya TERAKHIR berangkat - jawaban untuk nomor lain diabaikan (tinjauan kode 06-10-2026:
  // jawaban polis A yang tiba sesudah polis B tidak boleh mengisi No Master B).
  const terakhir = useRef('')

  /** Event change sel No Polis Treaty: hanya bila nilainya berubah sejak pemeriksaan terakhir. */
  const periksa = async () => {
    const n = noPolis.trim()
    if (n === diperiksa) return
    setDiperiksa(n)
    terakhir.current = n
    setGalat(null)
    setNoMaster('')
    setGalatTrt('')
    setSibuk(true)
    try {
      const r = await periksaPolis(n)
      if (terakhir.current !== n) return
      setNoMaster(r.noMaster)
      setGalatTrt(r.galat)
    } catch (e: unknown) {
      if (terakhir.current !== n) return
      setGalat(e)
    } finally {
      if (terakhir.current === n) setSibuk(false)
    }
  }

  const buat = async () => {
    setGalat(null)
    setSibuk(true)
    try {
      const k = await buatKasus({ noPolis: noPolis.trim(), edmType: jenis })
      onBuka(k.id)
    } catch (e: unknown) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  // nilai yang diketik sesudah pemeriksaan terakhir belum diperiksa: tombol Create menunggu event change
  const sesuai = noPolis.trim() === diperiksa
  const tampil = sesuai && tampilTombolBuat(galatTrt, noPolis.trim(), noMaster)

  return (
    <Modal
      judul={JUDUL.buat}
      onTutup={onTutup}
      aksi={
        tampil ? (
          <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void buat()}>
            {TOMBOL.buat}
          </button>
        ) : undefined
      }
    >
      {/* S2 "Errors" (TrtERR.CARI1 != '') */}
      {galatTrt !== '' && <div className="alert alert--warn">{galatTrt}</div>}
      {galat !== null && <Gagal galat={galat} />}
      <div className="edmt__kolom edmt__buat">
        <div onBlur={() => void periksa()}>
          <Field label={BUAT.noPolis} value={noPolis} onChange={setNoPolis} />
        </div>
        <Field label={BUAT.noMaster} value={noMaster} onChange={() => undefined} readOnly />
        <Pilih label={BUAT.jenis} value={jenis} opsi={opsiJenis} onChange={setJenis} />
      </div>
    </Modal>
  )
}
