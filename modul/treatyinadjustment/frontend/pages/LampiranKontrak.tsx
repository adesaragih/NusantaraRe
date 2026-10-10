// Layar Attachment + History modul Treaty In Adjustment.
//
// ⛔ DISALIN dari ekspor Pega modul INI — `Section/WorkAttachments.xml`
// (Category · Count · Upload file · View File, Download All, Refresh,
// spanduk biru) dan `Section/ShowAttachmentTreaty.xml` (File Name · Type).
// Ketiga Section lampiran ada di KEDUA ekspor; yang dipakai di sini milik
// modul ini sendiri.
//
// ⛔ Datanya dari `POOLDATA.M_ATTACHMENTTREATY_2` — tabel WARISAN, 43 baris.
// NOL tabel baru, NOL migrasi.
//
// ⚠️ Ronde sebelumnya menambahkan repository dan models di modul ini tetapi
// BERHENTI di sana: nol services, nol handlers, nol layar. Itu persis yang
// komentar modul ini peringatkan — *"tabel terisi yang tidak dibaca siapa
// pun adalah pekerjaan yang terlihat selesai dan tidak sampai ke pemakai."*
// Berkas ini jalan menuju ke sana.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { formatDateTime } from '../../../../inti/frontend/lib/format'
import {
  ambilLampiran,
  ambilRiwayat,
  type BarisRiwayatWarisan,
  type LampiranKontrak as IsiLampiran,
} from '../api'
import { KOLOM_BERKAS_LAMPIRAN, KOLOM_HISTORY, KOLOM_LAMPIRAN, LAMPIRAN } from '../labels'

export interface LampiranKontrakProps {
  /**
   * Pengenal kontrak sistem lama.
   *
   * ⛔ TEKS. `M_ATTACHMENTTREATY_2.TREATYID` adalah `VARCHAR2(100)`, berbeda
   * dari `idKontrak` model baru yang angka. Mengubahnya menjadi angka
   * menolak pengenal warisan yang sah.
   */
  masterID: string
  /**
   * `TreatyIn.ProportionType` — `NonProportional` menamai kode `00007`
   * dengan nama Non-Prop (`GetMasterTreatyCategory_Act` [2.2]).
   */
  jenis?: string
}

/**
 * Panel Attachment SAJA — dipakai halaman ini dan layar Adjustment.
 *
 * ⭐ Dipisah dari History 5 Oktober 2026: layar Adjustment menaruh deret
 * tombol DI ANTARA Attachment dan History (`InputTreatyInAdjustment.xml`
 * @276129 → @336222 → @380295), dan History-nya bersumber lain
 * (`TreatyIn.CommentList`). Pembacaan lampirannya tidak berubah.
 *
 * ⭐ 8 Oktober 2026 — penamaan kategori diperbaiki (tangkapan layar Pega):
 * kesebelas NAMA dari `M_KATEGORIMASTERTREATY`, urut nama. Spanduk "empat
 * nama belum dipastikan" dan baris `00003 nama kategori belum dipastikan`
 * dibuang — pertanyaannya terjawab katalog itu.
 */
export function PanelLampiranKontrak({ masterID, jenis = '' }: LampiranKontrakProps) {
  const [isi, setIsi] = useState<IsiLampiran | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(masterID !== '')
  // ⭐ Tombol `Refresh` layar lama: menaikkan penanda ini memicu pembacaan
  // ulang. Ia TIDAK menyimpan apa pun — panel ini baca-saja.
  const [muatUlang, setMuatUlang] = useState(0)

  useEffect(() => {
    if (masterID === '') {
      setMemuat(false)
      return
    }
    let dibuang = false
    setMemuat(true)
    setGalat(null)
    ambilLampiran(masterID, jenis)
      .then((l) => {
        if (!dibuang) setIsi(l)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
      .finally(() => {
        if (!dibuang) setMemuat(false)
      })
    return () => {
      dibuang = true
    }
  }, [masterID, jenis, muatUlang])

  const kategori = isi?.kategori ?? []
  const berkas = isi?.berkas ?? []

  return (
    <>
      {galat !== null && <Gagal galat={galat} />}
      {memuat && <Memuat />}

      <Panel judul={LAMPIRAN.judul}>
        {/* Spanduk biru — ATURAN nama berkas, bukan hiasan. Ditegakkan
            `services.NamaBerkasAman` di pintu masuk; KEPUTUSAN §18
            menetapkan ia TIDAK berlaku surut atas 43 baris yang sudah ada. */}
        <span className="tria__spanduk" role="note">
          {LAMPIRAN.spanduk}
        </span>

        <div className="table-wrap">
          <table className="tria__tabel">
            <thead>
              <tr>
                {KOLOM_LAMPIRAN.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {kategori.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_LAMPIRAN.length}>
                    <Kosong pesan={LAMPIRAN.tanpaIsi} petunjuk={LAMPIRAN.petunjukLampiran} />
                  </td>
                </tr>
              )}
              {kategori.map((k) => (
                <tr key={k.kode}>
                  {/* NAMA katalog; kode hanya untuk kode di data yang
                      tidak ada di katalog mana pun. */}
                  <td>{k.nama !== '' ? k.nama : k.kode}</td>
                  {/* ⛔ Cacah TIDAK diformat — ia butir, bukan uang. */}
                  <td>{String(k.cacah)}</td>
                  <td />
                  <td />
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="tria__aksi" role="group" aria-label={LAMPIRAN.judul}>
          {/* ⛔ `Download All` MATI: ia menarik berkas dari `T_STORAGE_IMAGE`,
              dan jalur itu belum dibangun. Tombol hidup yang tidak mengunduh
              apa pun lebih buruk daripada tombol mati yang mengatakan
              sebabnya. */}
          <button type="button" className="btn" disabled>
            {LAMPIRAN.unduhSemua}
          </button>
          {/* ⭐ `Refresh` HIDUP — ia hanya membaca ulang, dan membaca ulang
              adalah hal yang panel baca-saja memang boleh lakukan. */}
          <button
            type="button"
            className="btn"
            onClick={() => {
              setMuatUlang((n) => n + 1)
            }}
          >
            {LAMPIRAN.segarkan}
          </button>
        </div>

        <div className="table-wrap">
          <table className="tria__tabel">
            <thead>
              <tr>
                {KOLOM_BERKAS_LAMPIRAN.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {berkas.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_BERKAS_LAMPIRAN.length}>
                    <Kosong pesan={LAMPIRAN.tanpaIsi} petunjuk={LAMPIRAN.petunjukLampiran} />
                  </td>
                </tr>
              )}
              {berkas.map((b) => (
                <tr key={b.id}>
                  <td>{b.namaBerkas}</td>
                  <td>{b.jenisMime}</td>
                  <td>{b.diunggah}</td>
                  <td>{b.pengunggah}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>

    </>
  )
}

/**
 * Panel History — Date · PIC · Approval · Comment.
 *
 * ⛔ Hanya MENAMPILKAN baris yang diberikan; sumbernya urusan pemanggil.
 * Halaman ini memberinya `T_VIEW_COMMENT`; layar Adjustment memberinya
 * `TreatyIn.CommentList` dokumen penyesuaiannya sendiri — diukur, nol baris
 * `T_VIEW_COMMENT` berpengenal penyesuaian.
 */
export function PanelRiwayat({
  riwayat,
  petunjuk = LAMPIRAN.petunjukHistory,
}: {
  riwayat: readonly BarisRiwayatWarisan[]
  petunjuk?: string
}) {
  return (
    <>
      {/* ⭐ History KINI TERISI. Ronde sebelumnya meninggalkannya kosong
          dengan alasan "sumbernya milik modul sebelah" — dan alasan itu
          setengah keliru: yang ditolak `TestModulTidakMengimporModulLain`
          adalah impor paket Go, bukan pembacaan tabel. Modul ini membaca
          `T_VIEW_COMMENT` sendiri lewat repository-nya, dan nol baris
          mengimpor `treatyin`. */}
      <Panel judul={LAMPIRAN.judulHistory}>
        <div className="table-wrap">
          <table className="tria__tabel">
            <thead>
              <tr>
                {KOLOM_HISTORY.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {riwayat.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_HISTORY.length}>
                    <Kosong pesan={LAMPIRAN.tanpaIsi} petunjuk={petunjuk} />
                  </td>
                </tr>
              )}
              {riwayat.map((b, i) => (
                <tr key={i}>
                  {/* ⛔ Stempel Pega (`20261009T050155.604 GMT`) tampil apa
                      adanya sampai 8 Oktober 2026 — deretan angka mesin yang
                      harus dipecah sendiri oleh pembacanya. Jam ikut: dua
                      catatan sehari hanya dibedakan olehnya. */}
                  <td>{formatDateTime(b.tanggal)}</td>
                  <td>{b.operator}</td>
                  <td>{b.disetujui}</td>
                  <td>{b.catatan}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>
    </>
  )
}

export default function LampiranKontrak({ masterID }: LampiranKontrakProps) {
  const [riwayat, setRiwayat] = useState<BarisRiwayatWarisan[]>([])
  const [galat, setGalat] = useState<unknown>(null)

  // ⛔ Pembacaan KEDUA, terpisah dari lampiran — dan itu disengaja:
  // keduanya dari tabel berbeda, dan kegagalan salah satunya tidak boleh
  // mengosongkan yang lain.
  useEffect(() => {
    if (masterID === '') {
      setRiwayat([])
      return
    }
    let dibuang = false
    setGalat(null)
    ambilRiwayat(masterID)
      .then((r) => {
        if (!dibuang) setRiwayat(r)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [masterID])

  return (
    <>
      <PanelLampiranKontrak masterID={masterID} />
      {galat !== null && <Gagal galat={galat} />}
      <PanelRiwayat riwayat={riwayat} />
    </>
  )
}
