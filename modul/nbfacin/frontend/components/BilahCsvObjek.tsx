// Bilah "Please upload file with .csv format" General (tiket 66) - port `NB FacIn\Section\Periode.xml` sel 93-96
// (`ASM-FW-GISFW-DATA-OFFERFACIN!PERIODE`, disertakan InputInwardFacultative sel 107):
// - Download Template CSV (sel 93, `SetDownloadMasterCSV` -> Rule-File-Binary `TemplateCSVFire`): judul template 55 kolom;
// - Upload CSV (sel 94, `UploadCSV_DataPolis` / `UploadCSVDataPolis_Act`): berkas dibaca ke buffer layar (padanan
//   `TempWorkPage.ListB2BHostUpload`), belum menyimpan;
// - View Upload (sel 95, harness `ViewCSVResult_B2BHost`): isi buffer, baca-saja. Grid Pega (43 kolom, format B2B lama)
//   tidak sinkron dengan parsernya - di sini ditampilkan 55 kolom template `[keputusan agent]`;
// - Save Data (sel 96, `ChooseActivityMappingCSV_Act` -> `InsertUploadFire_act`): backend MENGGANTI semua objek dengan isi
//   CSV (seperti Pega); bila case sudah punya objek, konfirmasi dulu (keputusan work owner 08-10-2026 "Ganti semua +
//   konfirmasi"). Pesan hasil (Pega sel 98 `TempError.ERRMSG`) tampil di bawah bilah.

import { useRef, useState, type ChangeEvent, type ReactNode } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilObjek, simpanCsvObjek } from '../api'
import { PERIODE as P, TEKS_CSV_OBJEK as T } from '../labels'
import { bacaCsvObjek, KOLOM_CSV_OBJEK, TEMPLATE_CSV_OBJEK, type BarisCsvObjek } from './csvObjek'
import HalamanPopup from './HalamanPopup'

/** Baris per halaman popup View Upload. */
export const UKURAN_VIEW_UPLOAD = 10

function unduhTemplate() {
  const url = URL.createObjectURL(new Blob([TEMPLATE_CSV_OBJEK], { type: 'text/csv' }))
  const a = document.createElement('a')
  a.href = url
  a.download = 'templatecsvfire.csv'
  a.click()
  URL.revokeObjectURL(url)
}

export default function BilahCsvObjek({
  caseId,
  akhir,
  onTersimpan,
}: {
  caseId: string
  /** Tombol tambahan di ujung bilah (Insert Accumulation, akun IT). */
  akhir?: ReactNode
  /** Save Data berhasil: tab detail dimuat ulang. */
  onTersimpan: () => void
}) {
  const berkas = useRef<HTMLInputElement>(null)
  const [unggahan, setUnggahan] = useState<BarisCsvObjek[] | null>(null)
  const [pesan, setPesan] = useState<string[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [lihat, setLihat] = useState(false)
  const [hal, setHal] = useState(1)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)

  async function unggah(e: ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0]
    e.target.value = ''
    if (!f) return
    setGalat(null)
    const h = bacaCsvObjek(await f.text())
    if (h.galat !== '') {
      setUnggahan(null)
      setPesan([h.galat])
      return
    }
    setUnggahan(h.baris)
    setHal(1)
    setPesan([T.terbaca(h.baris.length)])
  }

  async function kirim() {
    if (unggahan === null) return
    setKonfirmasi(false)
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await simpanCsvObjek(caseId, unggahan)
      setPesan(h.pesan)
      setUnggahan(null)
      onTersimpan()
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  /** Save Data: tanpa unggahan -> pesan Pega; case sudah berobjek -> konfirmasi penggantian. */
  async function simpan() {
    if (unggahan === null) {
      setPesan([T.pleaseUpload])
      return
    }
    setGalat(null)
    try {
      const ada = (await ambilObjek(caseId)).baris.length > 0
      if (ada) setKonfirmasi(true)
      else await kirim()
    } catch (err) {
      setGalat(err)
    }
  }

  const awal = (hal - 1) * UKURAN_VIEW_UPLOAD
  return (
    <>
      <div className="nbf-opp__tombol">
        <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={unduhTemplate}>
          {P.downloadTemplateCsv.label}
        </button>
        <button type="button" className="btn nbf-tema__tombol btn--sm" disabled={menyimpan} onClick={() => berkas.current?.click()}>
          {P.uploadCsv.label}
        </button>
        <input ref={berkas} type="file" accept=".csv,text/csv" hidden aria-label={P.uploadCsv.label} onChange={(e) => void unggah(e)} />
        <button
          type="button"
          className="btn nbf-tema__tombol btn--sm"
          onClick={() => (unggahan === null ? setPesan([T.pleaseUpload]) : setLihat(true))}
        >
          {P.viewUpload.label}
        </button>
        <button type="button" className="btn nbf-tema__tombol btn--sm" disabled={menyimpan} onClick={() => void simpan()}>
          {menyimpan ? T.menyimpan : P.saveData.label}
        </button>
        {akhir}
      </div>
      <Gagal galat={galat} />
      {pesan.length > 0 && (
        <ul className="nbf-csvobj__pesan" aria-live="polite">
          {pesan.map((p, i) => (
            <li key={`${i}-${p}`}>{p}</li>
          ))}
        </ul>
      )}

      {lihat && unggahan !== null && (
        <Modal judul={P.viewUpload.label} onTutup={() => setLihat(false)} labelBatal={T.tutup} lebar>
          <HalamanPopup halaman={hal} ukuran={UKURAN_VIEW_UPLOAD} total={unggahan.length} onPindah={setHal} />
          <div className="table-wrap">
            <table className="nbf-tabel nbf-csvobj__tabel">
              <thead>
                <tr>
                  {KOLOM_CSV_OBJEK.map((k) => (
                    <th key={k} scope="col">
                      {k.replace(/_/g, ' ')}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {unggahan.slice(awal, awal + UKURAN_VIEW_UPLOAD).map((b, i) => (
                  <tr key={awal + i}>
                    {KOLOM_CSV_OBJEK.map((k) => (
                      <td key={k}>{b[k]}</td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Modal>
      )}

      {konfirmasi && (
        <Modal
          judul={P.saveData.label}
          onTutup={() => setKonfirmasi(false)}
          labelBatal={T.batal}
          aksi={
            <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void kirim()}>
              {T.ya}
            </button>
          }
        >
          <p>{T.konfirmasiGanti}</p>
        </Modal>
      )}
    </>
  )
}
