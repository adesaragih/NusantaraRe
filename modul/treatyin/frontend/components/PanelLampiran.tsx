// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { useEffect, useState } from 'react'

import { Kosong, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilPanelLampiran,
  unggahLampiran,
  type BarisKategoriLampiran,
  type BarisLampiranWarisan,
  type HasilBerkasUnggah,
} from '../api'
import {
  KOLOM_LAMPIRAN,
  KOLOM_LIHAT_BERKAS,
  LAMPIRAN,
} from '../labels'

/**
 * Panel Attachment — Category · Count · Upload file · View File.
 *
 * ⛔ Disalin dari `Section/WorkAttachments.xml`: keempat kolomnya, kedua
 * tombolnya, dan spanduk birunya. Datanya dari `M_ATTACHMENTTREATY_2`,
 * tabel WARISAN — nol tabel baru dibuat untuk panel ini.
 *
 * ⛔ SELURUH kategori tampil, termasuk yang nol berkas: kolom `Count` tidak
 * akan pernah berbunyi `0` kalau barisnya disembunyikan saat kosong.
 *
 * ⭐ `Upload file` HIDUP sejak 8 Oktober 2026 — keputusan pemakai: unggahan
 * masuk `M_ATTACHMENTTREATY_2`. Rantainya `Section/WorkAttachments.xml`:
 * tombol per kategori (syarat `TreatyIn.ViewState !='1' ||
 * TreatyIn.RevisionState='1'`) → `SetkategoriDoc` → modal `ASM Attach
 * Content` (`TreatyAttachContent`) → Attach (`TreatySaveAttachment`) →
 * Refresh (`GetMasterTreatyCategory_Act`). Kode kategorinya dari
 * `M_KATEGORIMASTERTREATY` — kesebelasnya kini dipastikan.
 */
export default function PanelLampiran({
  kategori: kategoriAwal,
  berkas: berkasAwal,
  idKontrak = '',
  bisaUnggah = false,
}: {
  kategori: readonly BarisKategoriLampiran[]
  berkas: readonly BarisLampiranWarisan[]
  /** `TreatyIn.ID` — kosong untuk kontrak yang belum tersimpan. */
  idKontrak?: string
  /** `TreatyIn.ViewState !='1' || TreatyIn.RevisionState='1'`. */
  bisaUnggah?: boolean
}) {
  // Isi panel — dari kontrak yang dimuat, lalu dari Refresh/Attach sendiri
  // (panel saja; isian form yang belum di-Save tidak ikut dibaca ulang).
  const [kategori, setKategori] = useState<readonly BarisKategoriLampiran[]>(kategoriAwal)
  const [berkas, setBerkas] = useState<readonly BarisLampiranWarisan[]>(berkasAwal)
  useEffect(() => {
    setKategori(kategoriAwal)
    setBerkas(berkasAwal)
  }, [kategoriAwal, berkasAwal])
  // Modal `ASM Attach Content` — kategori yang dipilih (`SetkategoriDoc`).
  const [unggahKe, setUnggahKe] = useState<BarisKategoriLampiran | null>(null)
  const [dipilih, setDipilih] = useState<File[]>([])
  const [sibuk, setSibuk] = useState(false)
  const [hasil, setHasil] = useState<{ galat: string; berkas: HasilBerkasUnggah[] } | null>(null)
  const segarkan = () => {
    if (idKontrak === '') return
    setSibuk(true)
    ambilPanelLampiran(idKontrak)
      .then((p) => {
        setKategori(p.kategoriLampiran)
        setBerkas(p.lampiran)
      })
      .catch((e: unknown) => {
        setHasil({ galat: e instanceof Error ? e.message : String(e), berkas: [] })
      })
      .finally(() => {
        setSibuk(false)
      })
  }
  const lampirkan = () => {
    if (unggahKe === null || dipilih.length === 0) return
    setSibuk(true)
    setHasil(null)
    unggahLampiran(idKontrak, unggahKe.kode, dipilih)
      .then((h) => {
        setKategori(h.kategoriLampiran)
        setBerkas(h.lampiran)
        setHasil({ galat: '', berkas: h.berkas })
        setUnggahKe(null)
        setDipilih([])
      })
      .catch((e: unknown) => {
        setHasil({ galat: e instanceof Error ? e.message : String(e), berkas: [] })
      })
      .finally(() => {
        setSibuk(false)
      })
  }
  /**
   * ⭐ PERMINTAAN PERUBAHAN, satu-satunya tempat ronde ini menyimpang dari
   * Pega — dan pemilik proses yang memintanya, pada keterangan gambar `25`:
   *
   *   "dan view upload saran dibuatkan pop up dan ini diubah menjadi
   *    design nya bagus"
   *
   * ⛔ Yang gambar 25 perlihatkan adalah JENDELA CHROME TERPISAH
   * (`ShowAttachmentTreaty`), lengkap dengan bilah alamat `appdev...`. Di
   * sini ia menjadi modal di dalam halaman.
   *
   * ⚠️ `Modal` DIPANGGIL dari `inti`, tidak ditulis ulang — modal kedua di
   * repositori ini berarti dua perilaku tutup, dua perangkap fokus, dan dua
   * tempat untuk salah.
   */
  const [berkasDilihat, setBerkasDilihat] = useState<string | null>(null)
  const kategoriDilihat = kategori.find((k) => k.kode === berkasDilihat) ?? null
  return (
    <Panel judul={LAMPIRAN.judul}>
      {/* Spanduk biru — ATURAN nama berkas, bukan hiasan. Ditegakkan di
          services (`NamaBerkasAman`) dan dinyatakan di sini. */}
      <span className="trin__spanduk" role="note">
        {LAMPIRAN.spanduk}
      </span>
      {/* `Refresh` — `GetMasterTreatyCategory_Act` (`!pyIsMobile`).
          `Download All` ber-`pyCondition never` di ekspor: tidak dirender. */}
      <div className="trin__aksi">
        <button type="button" className="btn btn--sm" disabled={sibuk || idKontrak === ''} onClick={segarkan}>
          {LAMPIRAN.segarkan}
        </button>
      </div>
      {hasil !== null && (
        <div className={hasil.galat !== '' ? 'alert alert--error' : 'alert alert--info'} role={hasil.galat !== '' ? 'alert' : 'status'}>
          {hasil.galat}
          {hasil.berkas.map((b) => (
            <div key={b.nama}>
              {b.nama}: {b.berhasil ? LAMPIRAN.terunggah : b.pesan}
            </div>
          ))}
        </div>
      )}
      {bisaUnggah && idKontrak === '' && <p className="trin__redup">{LAMPIRAN.simpanDulu}</p>}

      <div className="table-wrap">
        <table className="trin__tabel">
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
                  <Kosong pesan={LAMPIRAN.tanpaLampiran} petunjuk={LAMPIRAN.petunjukLampiran} />
                </td>
              </tr>
            )}
            {kategori.map((k) => (
              <tr key={k.kode}>
                {/* ⭐ NAMA apa adanya — kesebelasnya, persis tangkapan layar
                    pemilik proses 6 Oktober 2026. Sebelumnya empat baris
                    berbunyi KODE-nya (`00003 nama kategori belum
                    dipastikan`), sebab panel dirender dari katalog basis
                    data alih-alih dari daftar layar.

                    ⛔ Pertanyaan kode↔nama TIDAK tertutup oleh ini, dan
                    tidak boleh terlihat tertutup: `dipastikan` tetap
                    dibawa, dan jalur UNGGAH — yang memerlukan kodenya —
                    masih mati. Lihat
                    `docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`. */}
                <td>{k.nama}</td>
                {/* ⛔ Cacah TIDAK diformat — ia butir, bukan uang. */}
                <td>{String(k.cacah)}</td>
                <td>
                  {/* ⭐ `Upload file` — tampil bila `TreatyIn.ViewState !='1'
                      || TreatyIn.RevisionState='1'`; butuh kode kategori dan
                      kontrak yang sudah ber-ID. */}
                  {bisaUnggah && (
                    <button
                      type="button"
                      className="btn btn--sm"
                      disabled={sibuk || idKontrak === '' || k.kode === ''}
                      onClick={() => {
                        setUnggahKe(k)
                        setDipilih([])
                        setHasil(null)
                      }}
                    >
                      {LAMPIRAN.unggah}
                    </button>
                  )}
                </td>
                <td>
                  {/* ⭐ `View File` HIDUP — ia hanya menampilkan apa yang
                      sudah dibaca, dan menampilkan itu yang panel baca-saja
                      memang boleh lakukan. */}
                  <button
                    type="button"
                    className="btn btn--sm"
                    onClick={() => {
                      setBerkasDilihat(k.kode)
                    }}
                  >
                    {LAMPIRAN.lihatBerkas}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {unggahKe !== null && (
        <Modal
          judul={LAMPIRAN.judulUnggah}
          onTutup={() => {
            setUnggahKe(null)
          }}
          labelBatal={LAMPIRAN.batal}
          onKirim={lampirkan}
          aksi={
            <button type="submit" className="btn btn--primary" disabled={sibuk || dipilih.length === 0}>
              {sibuk ? LAMPIRAN.mengunggah : LAMPIRAN.lampirkan}
            </button>
          }
        >
          <p>
            <strong>{unggahKe.nama}</strong>
          </p>
          <input
            type="file"
            multiple
            aria-label={LAMPIRAN.pilihBerkas}
            onChange={(e) => {
              setDipilih(Array.from(e.target.files ?? []))
            }}
          />
          {dipilih.length > 0 && (
            <ul>
              {dipilih.map((f) => (
                <li key={f.name}>{f.name}</li>
              ))}
            </ul>
          )}
        </Modal>
      )}

      {berkasDilihat !== null && (
        <Modal
          judul={LAMPIRAN.judulLihatBerkas}
          onTutup={() => {
            setBerkasDilihat(null)
          }}
          labelBatal={LAMPIRAN.tutup}
          lebar
        >
          <div className="table-wrap">
            <table className="trin__tabel">
              <thead>
                <tr>
                  {/* ⛔ DUA kolom, persis gambar 25 — `File Name` dan
                      `Type`. Tidak tiga, tidak empat. */}
                  {KOLOM_LIHAT_BERKAS.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {berkasKategori(berkas, kategoriDilihat).length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_LIHAT_BERKAS.length}>{LAMPIRAN.tanpaBaris}</td>
                  </tr>
                )}
                {berkasKategori(berkas, kategoriDilihat).map((b) => (
                  <tr key={b.id}>
                    <td>{b.namaBerkas}</td>
                    <td>{b.jenisMime}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Modal>
      )}
    </Panel>
  )
}

/**
 * Berkas milik SATU kategori.
 *
 * ⚠️ Disaring menurut `kode`, bukan menurut nama: nama kategori yang
 * pasangannya belum dipastikan ditampilkan sebagai kodenya, dan menyaring
 * menurut yang tampil akan menyaring menurut dua hal yang berbeda.
 */
function berkasKategori(
  berkas: readonly BarisLampiranWarisan[],
  kategori: BarisKategoriLampiran | null,
): readonly BarisLampiranWarisan[] {
  if (kategori === null) return []
  // ⛔ Kode bila ia diketahui, NAMA bila tidak. Keempat kategori yang
  // kodenya belum dipastikan akan selalu berbunyi `No items` kalau
  // penyaringnya hanya kode — dan kosong yang salah terbaca persis seperti
  // kosong yang benar.
  if (kategori.kode !== '') return berkas.filter((b) => b.kodeKategori === kategori.kode)
  return berkas.filter((b) => b.namaKategori === kategori.nama)
}
