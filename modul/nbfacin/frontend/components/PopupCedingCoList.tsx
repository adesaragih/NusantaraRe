// Popup tombol `Change Ceding Co` di layar Inward Facultative (tiket 34).
//
// Port harness/section `ShowCedingCoList` (`D:\migrasi\RNM\NB FacIn\Section\ShowCedingCoList.xml`; dibuka
// Periode sel 53, Target popup): grid atas `.Quotation.CedingCoList` - kepala `Ceding Co` (wajib) + tombol
// `Add / Select Ceding`; tiap baris nama ceding (baca-saja) + `Delete`; tombol `Submit` di kaki. Tampilan =
// tangkapan layar work owner 03-10-2026 (daftar kosong: "No items").
//
// Pega: Add / Select Ceding → `AddCedingList_act` menambah baris KOSONG lalu membuka popup `CedingCompany`;
// Choose (`SetDataSobCeding_Act`) mengisi baris terakhir dengan `.ID` / `.ClientName` dari AGENT. Delete →
// `DeleteCeding_Act` (buang baris + buang kode/nama dari gabungan `;`). Submit → `SetCedingCo_Act` (gabung
// kode dan nama berpemisah `;` ke `.CedingCo` / `.CedingCoName`) → muat ulang layar pembuka → tutup.
//
// Keputusan agent (tiket 34):
// - E-5 baris ditambah HANYA saat Choose - tidak ada baris kosong tertinggal bila popup pencari ditutup tanpa
//   memilih (di Pega baris kosong itu tertahan oleh tanda wajib sel 17);
// - E-6 popup daftar dan popup pencari tampil BERGANTIAN, bukan bertumpuk: Modal inti memasang Escape di
//   document dan tidak memakai portal, sehingga dua modal bertumpuk tertutup bersamaan oleh satu Escape;
// - E-7 Submit menyerahkan daftar ke layar Inward (tersimpan lewat Save for later, sejalan E-1); Cancel / X
//   membuang perubahan. Ceding yang sama dapat dipilih dua kali - Pega tidak memeriksanya.

import { useState } from 'react'

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisCeding } from '../api'
import { POPUP_CEDING, TEKS_INWARD } from '../labels'
import PopupPilihAgent from './PopupPilihAgent'

export default function PopupCedingCoList({
  awal,
  onTutup,
  onSubmit,
}: {
  awal: BarisCeding[]
  onTutup: () => void
  onSubmit: (daftar: BarisCeding[]) => void
}) {
  const [daftar, setDaftar] = useState<BarisCeding[]>(awal)
  const [mencari, setMencari] = useState(false)

  if (mencari) {
    return (
      <PopupPilihAgent
        judul={POPUP_CEDING.judulCari}
        onTutup={() => setMencari(false)}
        onPilih={(b) => {
          setDaftar((d) => [...d, { id: b.id, name: b.name }])
          setMencari(false)
        }}
      />
    )
  }

  return (
    <Modal
      judul={POPUP_CEDING.judul}
      onTutup={onTutup}
      lebar
      aksi={
        <button type="button" className="btn btn--primary" onClick={() => onSubmit(daftar)}>
          {POPUP_CEDING.submit.label}
        </button>
      }
    >
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col">
                {POPUP_CEDING.kolom.label}
                <span className="field__req">*</span>
              </th>
              <th scope="col" className="table__actions">
                <button type="button" className="btn btn--ghost btn--sm" onClick={() => setMencari(true)}>
                  {POPUP_CEDING.tambah.label}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {daftar.length === 0 ? (
              <tr>
                <td colSpan={2}>{TEKS_INWARD.kosong}</td>
              </tr>
            ) : (
              daftar.map((c, i) => (
                <tr key={`${i}-${c.id}`}>
                  <td>{c.name}</td>
                  <td className="table__actions">
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => setDaftar((d) => d.filter((_, j) => j !== i))}
                    >
                      {POPUP_CEDING.hapus.label}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </Modal>
  )
}
