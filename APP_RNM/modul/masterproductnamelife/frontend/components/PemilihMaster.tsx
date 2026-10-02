// Tujuh pemilih master `Choose*` (PARITAS §4) - FlowAction `ChooseCeding`, `ChooseSOB`, `ChoosePolicyHolder`,
// `ChooseCurrency`, `ChooseRIRisk`, `ChooseRIRate`, `ChooseCauseOfLoss` → section `*_Section`.
//
// Isi section (sama di ketujuhnya): medan `Search` (`SearchPolicyHolder.CARI1`), Enter →
// `SearchPolicyHolder_act` 1 b236 `CARI1 = @toUpperCase(CARI1)` (server); grid RD berparam `CARI1` dengan kolom
// `ID` / `Name` (`RIRate Name` untuk R/I Rate) dan tombol baris `Choose` → `set*_DT` + `closeContainer`.
// Kaki FlowAction: `Submit` / `Cancel` (b32 / b31 `ChooseCeding.xml`) - keduanya menutup tanpa memilih.
// Grid terbuka dengan `CARI1` kosong = seluruh baris (`Contains ""`).

import { useCallback, useEffect, useRef, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariMaster, type JenisMaster, type NilaiMaster } from '../api'
import { UKURAN_HALAMAN_MPNL, potongHalaman } from '../bentuk'
import { LAIN_MPNL, PEMILIH_MPNL } from '../labels'

export default function PemilihMaster({
  judul,
  jenis,
  kolomNama = PEMILIH_MPNL.kolomName,
  onPilih,
  onTutup,
}: {
  /** Teks tombol pembukanya VERBATIM (`Choose Ceding Name`, ...). */
  judul: string
  jenis: JenisMaster
  kolomNama?: string
  onPilih: (v: NilaiMaster) => void
  onTutup: () => void
}) {
  const [kata, setKata] = useState('')
  const [daftar, setDaftar] = useState<NilaiMaster[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  // Grid RD berhalaman (`pyGridPaginator`): hasil RD bisa puluhan ribu baris (`CLIENT`) - yang dirender satu halaman.
  const [halaman, setHalaman] = useState(1)

  // Nomor permintaan terakhir: jawaban yang tiba SESUDAH permintaan yang lebih baru diabaikan (audit
  // 02-10-2026 - daftar penuh `CLIENT` yang lambat dulu dapat menimpa hasil pencarian yang lebih cepat).
  const terakhir = useRef(0)
  const muat = useCallback(
    async (cari: string) => {
      const nomor = ++terakhir.current
      setDaftar(null)
      setGalat(null)
      setHalaman(1)
      try {
        const d = await cariMaster(jenis, cari)
        if (nomor === terakhir.current) setDaftar(d.daftar)
      } catch (e) {
        if (nomor === terakhir.current) setGalat(e)
      }
    },
    [jenis],
  )

  useEffect(() => {
    void muat('')
  }, [muat])

  return (
    <Modal
      judul={judul}
      onTutup={onTutup}
      labelBatal={PEMILIH_MPNL.cancel}
      aksi={
        <button type="button" className="btn btn--primary" onClick={onTutup}>
          {PEMILIH_MPNL.submit}
        </button>
      }
    >
      <div className="field">
        <label className="field__label">{PEMILIH_MPNL.search}</label>
        <input
          className="field__input"
          value={kata}
          onChange={(e) => {
            setKata(e.target.value)
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault()
              void muat(kata)
            }
          }}
        />
      </div>
      {daftar === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={LAIN_MPNL.kosong} />}
      {daftar !== null && daftar.length > UKURAN_HALAMAN_MPNL && (
        <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_MPNL} total={daftar.length} onPindah={setHalaman} />
      )}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{PEMILIH_MPNL.kolomId}</th>
              <th>{kolomNama}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {potongHalaman(daftar, halaman).map((v) => (
              <tr key={v.id} className="inbox__baris">
                <td>{v.id}</td>
                <td>{v.nama}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      onPilih(v)
                      onTutup()
                    }}
                  >
                    {PEMILIH_MPNL.choose}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Modal>
  )
}
