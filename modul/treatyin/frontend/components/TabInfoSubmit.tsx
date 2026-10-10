// Tab `Information & Submit` — FORM, bukan grid riwayat.
//
// ---------------------------------------------------------------------
// ⛔ CACAT YANG DIPERBAIKI 6 Oktober 2026
// ---------------------------------------------------------------------
// Tab ini menampilkan grid riwayat (Date · Operator · Approved · Suggest),
// padahal di Pega ia FORM: dua medan teks panjang dan dua tombol.
//
// ⚠️ Dan yang ditampilkan bukan sekadar salah — ia SALINAN. Riwayat sudah
// punya panelnya sendiri di kaki layar (`PanelHistory`), membaca tabel yang
// sama. Jadi layar menampilkan daftar yang sama dua kali, dan tab yang
// seharusnya tempat MENGIRIM justru tempat MEMBACA.
//
// ---------------------------------------------------------------------
// ⭐ BENTUKNYA DISALIN DARI EKSPOR, BUKAN DARI GAMBAR
// ---------------------------------------------------------------------
// `Section/TreatyInfoSubmit.xml`, sesudah `pyIncludedRuleXML` bersarang
// dibuang dengan hitung kedalaman:
//
//   TreatyIn.Information   `Additional Information`   pyControlDisplayTitle
//                                                     = Text area
//   TreatyIn.Comment       `Comment`                  = Text area
//
//   tombol `Submit`         pyActivity TreatyInSubmit      gaya Strong
//     pyCondition  TreatyIn.ViewState !='1'
//                  && TreatyIn.StatusAkseptasi != 'Resolve Complete'
//   tombol `Decline offer`                                 gaya Simple
//     pyCondition  TreatyIn.ViewState != '1'
//
// ⚠️ Varian KETIGA ada di ekspor — `Submit` ber-`pyActivity`
// `TreatyInSubmitEDM`, bersyarat sama. Ia milik cabang Adjustment dan TIDAK
// dibangun di sini; menaruh dua tombol `Submit` di satu layar membuat yang
// menekannya tidak tahu mana yang ia jalankan.
//
// ⭐ JALUR REVISI (7 Oktober 2026) — `RevisionState = '1'`:
//
//   cell 9   Comment         `ViewState != 1 || RevisionState=1`
//   cell 21  Submit (revisi) `RevisionState='1'` saja → `TreatyInSubmit`
//
// Kontrak revisi selalu ber-`ViewState=1` (`SetTreatyIn_Act` [6], dan
// `TreatyInSetEdit` menyetelnya kembali bila `RevisionState==1`), jadi
// Submit biasa dan Decline offer tersembunyi — tepat satu Submit tampil.

import { useState } from 'react'

import { Area, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { INFO_SUBMIT } from '../labels'
import { TOMBOL_TULIS } from '../labelsTulis'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'
import { SelKosongPega, TataPegaBlok } from './tataPega'

export default function TabInfoSubmit({
  mode = 'lihat',
  statusAkseptasi,
  onKirim,
  onTolak,
  sibuk = false,
  revisi = false,
}: {
  mode?: ModeForm
  /** `TreatyIn.StatusAkseptasi` — menentukan tombol `Submit` tampil. */
  statusAkseptasi: string
  /**
   * ⭐ `Submit` → `TreatyInSubmit` (form menyusun dokumen dan menjalankannya).
   * Tidak diberikan (uji satu tab) = tombolnya mati.
   */
  onKirim?: () => void
  /** `Decline offer` → sesudah konfirmasi, `TreatyInDeclineConfirmation_postact`. */
  onTolak?: () => void
  /** Tombol tulis sedang berjalan. */
  sibuk?: boolean
  /** `TreatyIn.RevisionState = '1'` — kontrak sedang direvisi. */
  revisi?: boolean
}) {
  // Modal `TreatyInDeclineConfirmation`.
  const [konfirmasi, setKonfirmasi] = useState(false)
  // ⭐ `TreatyIn.Information` / `TreatyIn.Comment` di penampung halaman.
  const [info, setInfo] = useProperti('Information', '')
  const [komentar, setKomentar] = useProperti('Comment', '')

  // ⛔ `ViewState != '1'` di ekspor berarti: BUKAN mode lihat. Kedua tombol
  // memakai syarat itu, dan `Submit` menambahkan satu lagi.
  const bisaUbah = mode === 'ubah'
  // `ViewState = '1'`: mode lihat, ATAU kontrak revisi di mode mana pun.
  const viewState1 = !bisaUbah || revisi
  const bolehKirim = !viewState1 && statusAkseptasi !== 'Resolve Complete'

  return (
    <Panel judul={INFO_SUBMIT.judul}>
      {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/TreatyInfoSubmit.xml`
          `Inline grid double` L816, urutan sel: Additional Information
          (`Stacked with labels left` L1114) | tombol `1=2` (L1606) | Comment
          (`Stacked with labels left` L2061) | tombol `1=2` (L2669) | …
          Sel `1=2` tidak tampil tetapi TETAP memesan slotnya, jadi kedua
          Text area BERTUMPUK di separuh KIRI — persis gambar Pega 22. Tiap
          sel `.trin__kolom` = bentuk label-kiri kepala form. */}
      <TataPegaBlok tata="g2">
      <div className="trin__kolom trin__infosubmit">
        {/* Additional Information TIDAK ikut dibuka jalur revisi — hanya
            Comment (cell 9). Form yang membuka tab ini untuk revisi tidak
            mematikan isinya, jadi medan ini mematikan dirinya sendiri. */}
        <fieldset className="trin__mode" disabled={!bisaUbah}>
          <Area
            label={INFO_SUBMIT.infoTambahan}
            value={info}
            onChange={setInfo}
            baris={5}
          />
        </fieldset>
      </div>
      <SelKosongPega />
      <div className="trin__kolom">
        <Area label={INFO_SUBMIT.komentar} value={komentar} onChange={setKomentar} baris={5} />
      </div>
      <SelKosongPega />
      </TataPegaBlok>

      {/* ⛔ Kedua tombol TIDAK dirender di mode lihat — `pyCondition`-nya
          `TreatyIn.ViewState != '1'`, dan tombol yang Pega sembunyikan tidak
          boleh muncul di sini. */}
      {(!viewState1 || revisi) && (
        // ⭐ `Inline grid double` L5197 (cabang EDM L7668 sama): sel KIRI =
        // layout Submit/Submit revisi (L5218, selalu memesan slotnya), sel
        // KANAN = Decline offer — Decline duduk di awal separuh kanan (gambar
        // Pega 22 dan 41).
        // ⛔ Kedua tombol DIBUNGKUS sel `<div>`: tombol yang langsung menjadi
        // butir grid ikut melar selebar separuh layar. Di gambar 22/41
        // `Decline offer` berukuran tombol biasa di awal selnya.
        <div className="trin__aksi trin__tata trin__tata--g2" role="group" aria-label={INFO_SUBMIT.judul}>
          {/* ⭐ HIDUP sejak 7 Oktober 2026 — keputusan pemilik proses:
              Submit menyimpan ke tabel masing-masing DAN menjalankan tangga
              akseptasi (`TreatyInCheckError` → `Akseptasi_DT` →
              `AddCommentList_Act` → simpan). Bukan lewat Pega. */}
          <div>
          {(bolehKirim || revisi) && (
            <button type="button" className="btn btn--primary" disabled={onKirim === undefined || sibuk} onClick={onKirim}>
              {INFO_SUBMIT.kirim}
            </button>
          )}
          </div>
          <div>
          {!viewState1 && (
            <button
              type="button"
              className="btn"
              disabled={onTolak === undefined || sibuk}
              onClick={() => {
                setKonfirmasi(true)
              }}
            >
              {INFO_SUBMIT.tolak}
            </button>
          )}
          </div>
        </div>
      )}
      {konfirmasi && (
        <Modal
          judul={TOMBOL_TULIS.judulTolak}
          onTutup={() => {
            setKonfirmasi(false)
          }}
          labelBatal={TOMBOL_TULIS.batal}
          onKirim={() => {
            setKonfirmasi(false)
            onTolak?.()
          }}
          aksi={
            <button type="submit" className="btn btn--primary">
              {TOMBOL_TULIS.tolak}
            </button>
          }
        >
          <p>{TOMBOL_TULIS.tanyaTolak}</p>
        </Modal>
      )}
    </Panel>
  )
}
