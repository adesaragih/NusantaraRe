// Form "Opportunity" - dibuka dari `Create opportunity` di portal (tiket 26).
//
// ⛔ Sumber: TANGKAPAN LAYAR Pega kiriman work owner 02-10-2026 (`docs/02-layar/tangkapan/`) dan gambar
// keadaan awal `D:\migrasi\RNM\DDL\HALAMAN DEPAN NB.JPG`. Form ini tidak ada di korpus: di Pega ia lahir
// dari `createWork` atas `D_crmAppExtPage.WorkClass_Opportunity`, data page yang tidak diekspor
// (`SFAPortalOpportunitiesHeader.xml` L2486-L2504).
//
// Urutan dan kelompok = gambar: baris atas Estimated Closing Date + Owner; kolom kiri (Business Prospect
// Name, Group Business, Class Of Business, Type Of Inward [+ Type Of Facultative]), kolom kanan (Phase,
// Stage, Opportunity Source, Business Status); Description selebar form.
//
// Keputusan agent (tiket 26): dropdown hanya berisi nilai yang TERLIHAT di gambar (C-1) - Opportunity
// Source lengkap dari tangkapan layar dropdown terbuka; Type Of Inward awalnya kosong dan Type Of
// Facultative baru tampil bila Facultative dipilih (C-7, `[dugaan]` dari dua gambar); Search Group
// Business membuka popup ChooseAccount (C-8) yang daftarnya belum dapat dimuat (DDL `T_M_ACCOUNT` belum
// ada); dua tombol Group Business lain nonaktif (C-2); Stage read-only (C-3); Class Of Business kotak teks
// (C-4); tanpa tombol simpan dan tanpa endpoint (C-5); nol catatan pengembang (C-6).

import { useState } from 'react'

import { Area, BelumTersedia, Field, FieldTanggal, Modal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  FORM_OPPORTUNITY as F,
  NILAI_AWAL_OPPORTUNITY as AWAL,
  OPSI_OPPORTUNITY_SOURCE,
  POPUP_CHOOSE_ACCOUNT as POPUP,
  TEKS_FORM_OPPORTUNITY as TEKS,
  TOMBOL_FORM_OPPORTUNITY as TOMBOL,
} from '../labels'

const tanpaUbah = () => {}

/** Satu nilai yang terlihat di gambar - bukan daftar pilihan lengkap (C-1). */
const satu = (nilai: string): Opsi[] => [{ value: nilai, label: nilai }]

/** Daftar lengkap dari tangkapan layar dropdown terbuka (Opportunity Source, 02-10-2026). */
const daftar = (nilai: readonly string[]): Opsi[] => nilai.map((n) => ({ value: n, label: n }))

/**
 * Popup tombol `Search Group Business` (C-8): kotak Search + tombol Search, lalu grid Insured ID ·
 * Insured Name · Group Business dengan tombol Choose per baris. Sumber datanya (`T_M_ACCOUNT`) belum
 * tersambung, jadi grid tampil sebagai `BelumTersedia` - bukan tabel kosong - dan Search nonaktif.
 */
export function PopupChooseAccount({ onTutup }: { onTutup: () => void }) {
  return (
    <Modal judul={POPUP.judul} onTutup={onTutup} lebar>
      <div className="nbf-popup__cari">
        <div className="nbf-popup__kotak">
          <Field label={POPUP.cari} value="" onChange={tanpaUbah} readOnly />
        </div>
        <button type="button" className="btn btn--ghost btn--sm" disabled>
          {POPUP.tombolCari}
        </button>
      </div>
      <div className="table-wrap">
        <table className="table">
          <thead>
            <tr>
              <th scope="col" />
              {POPUP.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
        </table>
      </div>
      <BelumTersedia apa={TEKS.daftarAccount} />
    </Modal>
  )
}

export default function FormOpportunity({ pemilik }: { pemilik: string }) {
  const [tanggalTutup, setTanggalTutup] = useState('')
  const [namaProspek, setNamaProspek] = useState('')
  const [classOfBusiness, setClassOfBusiness] = useState('')
  const [typeOfInward, setTypeOfInward] = useState('')
  const [typeOfFacultative, setTypeOfFacultative] = useState<string>(AWAL.typeOfFacultative)
  const [phase, setPhase] = useState<string>(AWAL.phase)
  const [sumber, setSumber] = useState('')
  const [statusBisnis, setStatusBisnis] = useState<string>(AWAL.statusBisnis)
  const [deskripsi, setDeskripsi] = useState('')
  const [cariGrup, setCariGrup] = useState(false)

  return (
    <div className="nbfacin">
      <section className="panel">
        <h4 className="panel__title">{F.judul}</h4>
        <div className="nbf-opp__atas">
          <FieldTanggal label={F.tanggalTutup} value={tanggalTutup} onChange={setTanggalTutup} required />
          <div className="field">
            <span className="field__label">{F.owner}</span>
            <div className="nbf-opp__owner">{pemilik}</div>
          </div>
        </div>
        <div className="nbf-opp__kolom">
          <div className="nbf-opp__tumpuk">
            <Field label={F.namaProspek} value={namaProspek} onChange={setNamaProspek} required />
            <div className="field">
              <span className="field__label">{F.grupBisnis}</span>
              <div className="nbf-opp__tombol">
                <button type="button" className="btn btn--ghost btn--sm" onClick={() => setCariGrup(true)}>
                  {TOMBOL.cariGrup}
                </button>
                <button type="button" className="btn btn--sm" disabled>
                  {TOMBOL.perusahaanBaru}
                </button>
                <button type="button" className="btn btn--sm" disabled>
                  {TOMBOL.grupBaru}
                </button>
              </div>
            </div>
            <Field label={F.classOfBusiness} value={classOfBusiness} onChange={setClassOfBusiness} required />
            <div className="nbf-opp__pasangan">
              <Pilih
                label={F.typeOfInward}
                value={typeOfInward}
                onChange={setTypeOfInward}
                opsi={satu(AWAL.typeOfInward)}
                kosong={AWAL.inwardKosong}
                required
              />
              {typeOfInward === AWAL.typeOfInward && (
                <Pilih
                  label={F.typeOfFacultative}
                  value={typeOfFacultative}
                  onChange={setTypeOfFacultative}
                  opsi={satu(AWAL.typeOfFacultative)}
                  required
                />
              )}
            </div>
          </div>
          <div className="nbf-opp__tumpuk">
            <Pilih label={F.phase} value={phase} onChange={setPhase} opsi={satu(AWAL.phase)} required />
            <Field label={F.stage} value={AWAL.stage} onChange={tanpaUbah} readOnly />
            <Pilih label={F.sumber} value={sumber} onChange={setSumber} opsi={daftar(OPSI_OPPORTUNITY_SOURCE)} kosong={AWAL.sumberKosong} />
            <Pilih label={F.statusBisnis} value={statusBisnis} onChange={setStatusBisnis} opsi={satu(AWAL.statusBisnis)} required />
          </div>
        </div>
        <div className="nbf-opp__bawah">
          <Area label={F.deskripsi} value={deskripsi} onChange={setDeskripsi} baris={5} />
        </div>
      </section>
      {cariGrup && <PopupChooseAccount onTutup={() => setCariGrup(false)} />}
    </div>
  )
}
