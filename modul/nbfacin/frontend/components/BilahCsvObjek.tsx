// Bilah "Please upload file with .csv format" General (tiket 66) - port `NB FacIn\Section\Periode.xml` sel 93-96
// (`ASM-FW-GISFW-DATA-OFFERFACIN!PERIODE`, disertakan InputInwardFacultative sel 107):
// - Download Template CSV (sel 93, `SetDownloadMasterCSV` -> Rule-File-Binary `TemplateCSVFire`): judul template 55 kolom;
// - Upload CSV (sel 94, `UploadCSV_DataPolis` / `UploadCSVDataPolis_Act`);
// - View Upload (sel 95, harness `ViewCSVResult_B2BHost`): isi unggahan, baca-saja. Grid Pega (43 kolom, format B2B lama)
//   tidak sinkron dengan parsernya - di sini ditampilkan kolom template `[keputusan agent]`;
// - Save Data (sel 96, `ChooseActivityMappingCSV_Act` -> `InsertUploadFire_act`): backend MENGGANTI semua objek dengan isi
//   CSV (seperti Pega); bila case sudah punya objek, konfirmasi dulu (keputusan work owner 08-10-2026 "Ganti semua +
//   konfirmasi"). Pesan hasil (Pega sel 98 `TempError.ERRMSG`) tampil di bawah bilah.
//
// Tiket 75 (work owner 10-10-2026 "UNGGAHAN CSV TIDAK USAH ADA BATASANNYA, JADI 1 JUTA JUGA HARUS BISA"; K75-1..2): tanpa batas
// baris. Berkas TIDAK dibaca utuh di peramban - hanya judulnya diperiksa dari potongan awal, lalu berkas di-stream ke server
// (`POST …/objek/csv/berkas`). Validasi dan Save Data berjalan di LATAR per kelompok; bilah membaca status berkala (kemajuan,
// galat per baris, Ulangi bila gagal). Status dipulihkan sesudah muat ulang (`…/objek/csv/aktif`).

import { useCallback, useEffect, useRef, useState, type ChangeEvent, type ReactNode } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilBarisCsvObjek,
  ambilTugasCsvAktif,
  ambilTugasCsvObjek,
  simpanTugasCsvObjek,
  ulangiTugasCsvObjek,
  unggahCsvObjek,
  type TugasCsvObjek,
} from '../api'
import { PERIODE as P, TEKS_CSV_OBJEK as T } from '../labels'
import { KOLOM_CSV_CARGO, periksaJudulBerkasCargo, TEMPLATE_CSV_CARGO } from './csvCargo'
import { KOLOM_CSV_MARINE_HULL, periksaJudulBerkasMarineHull, TEMPLATE_CSV_MARINE_HULL } from './csvMarineHull'
import { KOLOM_CSV_OBJEK, periksaJudulBerkas, TEMPLATE_CSV_OBJEK } from './csvObjek'
import HalamanPopup from './HalamanPopup'

/** Baris per halaman popup View Upload. */
export const UKURAN_VIEW_UPLOAD = 10

/** Jeda baca status selama pekerjaan berjalan. */
export const JEDA_STATUS_MS = 2000

/** Status yang masih berjalan di server (bilah terus membaca kemajuan). */
export const sedangBerjalan = (t: TugasCsvObjek | null) => t !== null && (t.status === 'memvalidasi' || t.status === 'menyimpan')

/** Pesan bilah untuk satu status pekerjaan (teks sistem baru, tiket 75). */
export function pesanTugas(t: TugasCsvObjek): string[] {
  switch (t.status) {
    case 'memvalidasi':
      return [T.memvalidasi(t.barisDiproses)]
    case 'valid':
      return [T.valid(t.totalBaris)]
    case 'gagal_validasi':
      return [T.gagalValidasi(t.jumlahGalat), ...t.galat, ...(t.jumlahGalat > t.galat.length ? [T.galatLain(t.jumlahGalat - t.galat.length)] : [])]
    case 'menyimpan':
      return [T.menyimpanKemajuan(t.barisTersimpan, t.totalBaris)]
    case 'selesai':
      return t.pesan
    case 'gagal':
      // galatProses = sebab berhentinya proses (mis. terputus / galat database) - bukan galat baris CSV.
      return [T.gagal(t.barisTersimpan, t.totalBaris), ...(t.galatProses ? [t.galatProses] : []), ...t.pesan]
  }
}

/** Jenis bilah CSV = importer Pega `ChooseActivityMappingCSV_Act` menurut COB: FIRE (L439) / MARINE CARGO (L539) / MARINE HULL
 *  (IsAneka L639, template `SetDownloadMasterCSV` IsMarineHull; tiket 73 P6). */
export type JenisCsv = 'fire' | 'cargo' | 'marinehull'

/** Per jenis: kolom cadangan View Upload, isi + nama berkas Download Template CSV (nama berkas Pega), cek judul berkas. */
export const CSV_PER_JENIS: Record<JenisCsv, { kolom: readonly string[]; template: string; berkas: string; periksa: (b: Blob) => Promise<string> }> = {
  fire: { kolom: KOLOM_CSV_OBJEK, template: TEMPLATE_CSV_OBJEK, berkas: 'templatecsvfire.csv', periksa: periksaJudulBerkas },
  cargo: { kolom: KOLOM_CSV_CARGO, template: TEMPLATE_CSV_CARGO, berkas: 'templatecsvmarinecargo.csv', periksa: periksaJudulBerkasCargo },
  marinehull: { kolom: KOLOM_CSV_MARINE_HULL, template: TEMPLATE_CSV_MARINE_HULL, berkas: 'templatecsvmarinehull.csv', periksa: periksaJudulBerkasMarineHull },
}

/** Download Template CSV - FIRE `TemplateCSVFire`; Marine Cargo / Marine Hull = judul template Pega (K74-14 / K73-60). */
function unduhTemplate(jenis: JenisCsv) {
  const url = URL.createObjectURL(new Blob([CSV_PER_JENIS[jenis].template], { type: 'text/csv' }))
  const a = document.createElement('a')
  a.href = url
  a.download = CSV_PER_JENIS[jenis].berkas
  a.click()
  URL.revokeObjectURL(url)
}

/** Popup View Upload - baris unggahan dibaca dari server per halaman; kolom = judul berkas dari server, cadangan per jenis. */
function LihatUnggahan({ caseId, idUnggah, jenis, onTutup }: { caseId: string; idUnggah: string; jenis: JenisCsv; onTutup: () => void }) {
  const [hal, setHal] = useState(1)
  const [data, setData] = useState<{ total: number; baris: Record<string, string>[]; kolom?: string[] } | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  useEffect(() => {
    let batal = false
    setData(null)
    ambilBarisCsvObjek(caseId, idUnggah, hal, UKURAN_VIEW_UPLOAD).then(
      (h) => {
        if (!batal) setData(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId, idUnggah, hal])
  const kolom: readonly string[] = data?.kolom ?? CSV_PER_JENIS[jenis].kolom
  return (
    <Modal judul={P.viewUpload.label} onTutup={onTutup} labelBatal={T.tutup} lebar>
      {data !== null && <HalamanPopup halaman={hal} ukuran={UKURAN_VIEW_UPLOAD} total={data.total} onPindah={setHal} />}
      {data === null && galat === null && <Memuat />}
      <Gagal galat={galat} />
      {data !== null && (
        <div className="table-wrap">
          <table className="nbf-tabel nbf-csvobj__tabel">
            <thead>
              <tr>
                {kolom.map((k) => (
                  <th key={k} scope="col">
                    {k.replace(/_/g, ' ')}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data.baris.map((b, i) => (
                <tr key={(hal - 1) * UKURAN_VIEW_UPLOAD + i}>
                  {kolom.map((k) => (
                    <td key={k}>{b[k] ?? ''}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Modal>
  )
}

export default function BilahCsvObjek({
  caseId,
  jenis = 'fire',
  akhir,
  onTersimpan,
}: {
  caseId: string
  /** Tiket 74 P6 / 73 P6: 'cargo' / 'marinehull' (cek judul + template + View Upload per jenis; server memilih importer menurut COB). */
  jenis?: JenisCsv
  /** Tombol tambahan di ujung bilah (Insert Accumulation, akun IT). */
  akhir?: ReactNode
  /** Save Data selesai: tab detail dimuat ulang. */
  onTersimpan: () => void
}) {
  const berkas = useRef<HTMLInputElement>(null)
  const [tugas, setTugas] = useState<TugasCsvObjek | null>(null)
  const [pesan, setPesan] = useState<string[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [mengunggah, setMengunggah] = useState(false)
  const [lihat, setLihat] = useState(false)
  const [konfirmasi, setKonfirmasi] = useState(false)
  // Status sebelumnya - `onTersimpan` dipanggil sekali saat berpindah dari 'menyimpan' ke 'selesai'.
  const statusLama = useRef<string | null>(null)

  const terima = useCallback(
    (t: TugasCsvObjek | null) => {
      setTugas(t)
      setPesan(t === null ? [] : pesanTugas(t))
      if (t !== null && t.status === 'selesai' && statusLama.current === 'menyimpan') onTersimpan()
      statusLama.current = t?.status ?? null
    },
    [onTersimpan],
  )

  // Pulihkan pekerjaan terakhir case (mis. sesudah muat ulang saat penyimpanan latar masih berjalan). Status case sebelumnya
  // dikosongkan dulu (bug DEV 10-10-2026: unggahan case FIRE tampil di case MARINE CARGO karena `aktif` case baru = null).
  useEffect(() => {
    let batal = false
    statusLama.current = null
    setTugas(null)
    setPesan([])
    setGalat(null)
    ambilTugasCsvAktif(caseId).then(
      (t) => {
        if (batal || t === null) return
        statusLama.current = t.status
        setTugas(t)
        // Pekerjaan lama yang sudah selesai tidak perlu dipamerkan lagi; yang berjalan / gagal / siap simpan ditampilkan.
        if (t.status !== 'selesai') setPesan(pesanTugas(t))
      },
      () => undefined,
    )
    return () => {
      batal = true
    }
  }, [caseId])

  // Selama berjalan di server: baca status berkala.
  const idBerjalan = sedangBerjalan(tugas) ? tugas!.idUnggah : null
  useEffect(() => {
    if (idBerjalan === null) return
    let batal = false
    const jam = window.setInterval(() => {
      ambilTugasCsvObjek(caseId, idBerjalan).then(
        (t) => {
          if (!batal) terima(t)
        },
        (err: unknown) => {
          if (!batal) setGalat(err)
        },
      )
    }, JEDA_STATUS_MS)
    return () => {
      batal = true
      window.clearInterval(jam)
    }
  }, [caseId, idBerjalan, terima])

  async function unggah(e: ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0]
    e.target.value = ''
    if (!f) return
    setGalat(null)
    const salahJudul = await CSV_PER_JENIS[jenis].periksa(f)
    if (salahJudul !== '') {
      setPesan([salahJudul])
      return
    }
    setMengunggah(true)
    setPesan([T.mengunggah(f.name, (f.size / (1024 * 1024)).toLocaleString('id-ID', { maximumFractionDigits: 1 }))])
    try {
      terima(await unggahCsvObjek(caseId, f))
    } catch (err) {
      setPesan([])
      setGalat(err)
    } finally {
      setMengunggah(false)
    }
  }

  async function mulaiSimpan() {
    if (tugas === null) return
    setKonfirmasi(false)
    setGalat(null)
    try {
      terima(await simpanTugasCsvObjek(caseId, tugas.idUnggah))
    } catch (err) {
      setGalat(err)
    }
  }

  /** Save Data: hanya dari unggahan yang valid; case sudah berobjek -> konfirmasi penggantian. */
  function simpan() {
    if (tugas === null || tugas.status !== 'valid') {
      setPesan(tugas !== null && tugas.status !== 'selesai' ? pesanTugas(tugas) : [T.pleaseUpload])
      return
    }
    if (tugas.adaObjek) setKonfirmasi(true)
    else void mulaiSimpan()
  }

  async function ulangi() {
    if (tugas === null) return
    setGalat(null)
    try {
      terima(await ulangiTugasCsvObjek(caseId, tugas.idUnggah))
    } catch (err) {
      setGalat(err)
    }
  }

  const sibuk = mengunggah || sedangBerjalan(tugas)
  const bisaLihat = tugas !== null && tugas.status !== 'gagal_validasi' && tugas.status !== 'selesai'
  return (
    <>
      <div className="nbf-opp__tombol">
        <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => unduhTemplate(jenis)}>
          {P.downloadTemplateCsv.label}
        </button>
        <button type="button" className="btn nbf-tema__tombol btn--sm" disabled={sibuk} onClick={() => berkas.current?.click()}>
          {P.uploadCsv.label}
        </button>
        <input ref={berkas} type="file" accept=".csv,text/csv" hidden aria-label={P.uploadCsv.label} onChange={(e) => void unggah(e)} />
        <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => (bisaLihat ? setLihat(true) : setPesan([T.pleaseUpload]))}>
          {P.viewUpload.label}
        </button>
        <button type="button" className="btn nbf-tema__tombol btn--sm" disabled={sibuk} onClick={simpan}>
          {tugas?.status === 'menyimpan' ? T.menyimpan : P.saveData.label}
        </button>
        {tugas?.status === 'gagal' && (
          <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--sm" onClick={() => void ulangi()}>
            {T.ulangi}
          </button>
        )}
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

      {lihat && tugas !== null && <LihatUnggahan caseId={caseId} idUnggah={tugas.idUnggah} jenis={jenis} onTutup={() => setLihat(false)} />}

      {konfirmasi && (
        <Modal
          judul={P.saveData.label}
          onTutup={() => setKonfirmasi(false)}
          labelBatal={T.batal}
          aksi={
            <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void mulaiSimpan()}>
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
