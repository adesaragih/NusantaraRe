// Tab **EGNPI** cabang NON-PROPORSIONAL — bentuk dan rumus dari ekspor.
//
// ---------------------------------------------------------------------
// ⛔ YANG DIGANTI, DAN MENGAPA
// ---------------------------------------------------------------------
// Sebelumnya tab ini `TabGridWarisan`: delapan kolom, hanya dibaca, nol
// tombol. Ekspor Pega memperlihatkan hal yang berbeda — grid yang DAPAT
// disunting dengan `Add`/`Delete`, dua panel total, dan DUA tombol berumus:
//
//   Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-3 `EGNPI`
//     panel  `Estimate Gross Net Premium Income`  grid + Add + Delete
//     panel  `Total EGNPI Amount` | `Value`       per mata uang
//     panel  `Total Amount in IDR` · `Total Proportion %`
//     tombol `Update Total`        → TreatyInNPSetTotal(type=egnpi)
//     tombol `Update EGNPI Value`  → TreatyInEGNPIListValue
//
// ⚠️ Jadi yang hilang bukan hiasan: tanpa `Add` nol baris EGNPI dapat lahir,
// dan tanpa kedua tombol itu `Amount in IDR` serta `Proportion %` nol pernah
// terisi — padahal tab Limits MEMBACA keduanya lewat `TotalEgnpi`.
//
// ---------------------------------------------------------------------
// ⭐ RUMUSNYA DI SERVICES, BUKAN DI SINI
// ---------------------------------------------------------------------
// `backend/services/hitung_egnpi.go` — disalin dari Activity dan diuji.
// Layar ini hanya mengirim isian dan menampilkan jawabannya; nol aritmetika
// di berkas ini, sama seperti tab Limits.
//
// ---------------------------------------------------------------------
// ⭐ BENTUK TAMPILAN meniru tab Limits (permintaan pemilik proses)
// ---------------------------------------------------------------------
// Pega memberi tab ini SATU grid yang barisnya dapat dibuka menjadi rincian
// (`Section/DetailEGNPI.xml`). Itu persis bentuk `KartuLipat` tab Limits:
// kepala kartu memuat ringkasan yang di Pega menjadi kolom grid, dan isinya
// medan rincian. Jadi "meniru Limits" dan "sama seperti Pega" di sini satu
// hal, bukan dua.
//
// ⛔ Aturan B pemilik proses tetap berlaku: mengetik TIDAK menyentuh basis
// data. Suntingan hidup di keadaan layar sampai Save/Submit — dan rute yang
// dipanggil di sini ber-awalan `/hitung/`, yang menyatakan dirinya nol tulis.

import { useCallback, useEffect, useState } from 'react'

import { PemicuUbah } from './pemicuUbah'

import { Field, FieldAngka, Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilOpsiLimits,
  hitungEgnpi,
  type AksiEgnpi,
  type BarisEgnpi,
  type KursLimitNP,
  type OpsiLimits,
  type PilihanWarisan,
} from '../api'
import { DESIMAL_EGNPI, EGNPI } from '../labelsEgnpi'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'
import TanggalRedup from './TanggalRedup'
import DropdownWarisan from './DropdownWarisan'
import { Bagian, Chip, KartuLipat, KepalaBagian, TombolHapus, TombolTambah } from './limitsUI'
import { formatLimit } from './TabLimitsProp'

/** Angka tampil — pemformat modul, bukan pemformat kedua. */
const tampil = (nilai: string, desimal: number) => formatLimit('uang', desimal, nilai)

/** Satu baris total per mata uang. */
interface BarisTotal {
  Currency: string
  CurrencyID: string
  Value: string
}

/** Grid dua kolom hanya-baca — panel `Total EGNPI Amount`. */
function GridTotal({ baris }: { baris: readonly BarisTotal[] }) {
  return (
    <div className="table-wrap trin__share-total">
      <table className="trin__tabel">
        <thead>
          <tr>
            <th scope="col">{EGNPI.totalPerMataUang}</th>
            <th scope="col">{EGNPI.nilai}</th>
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={2}>{EGNPI.tanpaBaris}</td>
            </tr>
          )}
          {baris.map((b, i) => (
            <tr key={String(i)}>
              <td>{b.Currency}</td>
              <td>{tampil(b.Value, DESIMAL_EGNPI.nilaiPerMataUang)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/** Rincian satu baris — `Section/DetailEGNPI.xml`, urut dan hak ubahnya. */
/**
 * Rincian satu baris — `Section/DetailEGNPI.xml`, urut dan hak ubahnya.
 *
 * ⛔ BENTUK DIPERBAIKI 7 Oktober 2026 atas tangkapan layar Pega pemilik
 * proses. Tiga hal yang sebelumnya salah:
 *
 *  1. `As At` kotak teks biasa → ia `pxDateTime` di ekspor, dan tangkapan
 *     layar memperlihatkan ikon kalender. Kini `FieldTanggal`.
 *  2. Autocomplete ketik-bebas → pemilik proses: samakan dengan dropdown
 *     `Ceding` dan `Source of Business`. Kini `DropdownWarisan`, komponen
 *     yang SAMA dengan kedua pemilih itu — bukan tiruannya.
 *  3. `Amount` dan `Amount in IDR` dirender sebagai medan tunggal → di
 *     ekspor keduanya SEPASANG: label menempel pada medan mata uang, dan
 *     medan nilai di sebelahnya TANPA label sendiri.
 *       `Amount`         .Currency (pilihan) + .Amount
 *       `Amount in IDR`  "IDR" (tetap)       + .AmountIDR
 *     Tangkapan layar memperlihatkan persis itu, berikut kotak `IDR`
 *     kelabu di kiri baris kedua.
 */
function RincianBaris({
  b,
  bisaUbah,
  ambilMataUang,
  ambilKelompokTreaty,
  onUbah,
  onKonversi,
}: {
  b: BarisEgnpi
  bisaUbah: boolean
  ambilMataUang: () => Promise<PilihanWarisan[]>
  ambilKelompokTreaty: () => Promise<PilihanWarisan[]>
  onUbah: (baru: BarisEgnpi) => void
  /**
   * `SetAmountConversion`. `baru` = baris yang BARU diubah di event yang sama
   * — ⛔ tanpanya konversi berjalan atas baris LAMA (keadaan belum dirender)
   * dan jawabannya menimpa mata uang yang baru dipilih (laporan pemakai
   * 8 Oktober 2026: "set idr atau yg lain tidak bisa malah hilang").
   */
  onKonversi: (baru?: BarisEgnpi) => void
}) {
  return (
    <Bagian judul={EGNPI.panel}>
      <div className="trin__egnpi-medan">
        {/* ⚠️ EKSPOR BERSELISIH DENGAN DIRINYA SENDIRI di medan ini: grid
            membuatnya DAPAT diubah, rinciannya `pyReadOnly=true` — tetapi
            sel yang sama memuat `pyReadOnlyCondition = ViewState = 1`, jadi
            ia aktif di mode Edit. Lihat `alat-baca-ekspor` Aturan 3. */}
        {bisaUbah ? (
          <DropdownWarisan
            label={EGNPI.treatyGroup}
            ambil={ambilKelompokTreaty}
            nilai={b.TreatyGroup}
            onPilih={(o) => {
              onUbah({ ...b, TreatyGroup: o.nama, TreatyGroupID: o.id })
            }}
          />
        ) : (
          <Field label={EGNPI.treatyGroup} value={b.TreatyGroup} readOnly onChange={() => undefined} />
        )}

        {/* ⭐ `pxDateTime` — tanggal yang DAPAT DIKETIK atau dipilih dari
            kalender (`TanggalRedup`), bentuk kabel DD-MM-YYYY, sama dengan
            seluruh tanggal modul ini. */}
        {bisaUbah ? (
          <TanggalRedup
            label={EGNPI.asAt}
            value={b.AsDate}
            onChange={(v) => {
              onUbah({ ...b, AsDate: v })
            }}
          />
        ) : (
          <Field label={EGNPI.asAt} value={b.AsDate} readOnly onChange={() => undefined} />
        )}

        {/* Baris `Amount` — pilihan mata uang lalu kotak nilai.
            ⛔ `Field` inti NOL punya `onBlur`, dan ia tidak diubah dari sini:
            ia dipakai seluruh aplikasi. Pembungkus ini yang menangkapnya
            (`onBlur` React menggelembung), jadi keluar dari kotak Amount
            memicu `SetAmountConversion`.
            ⚠️ Ekspor memicunya pada SETIAP perubahan; di sini pemicunya
            kehilangan fokus — satu panggilan jaringan per huruf yang diketik
            bukan kesetiaan, melainkan layar yang tersendat. */}
        {/* ⭐ `change` Pega: konversi berjalan hanya bila Amount BERUBAH
            (mata uang memicunya sendiri saat dipilih, `onPilih`). */}
        <PemicuUbah className="trin__egnpi-pasangan" nilai={b.Amount} aktif={bisaUbah} aksi={onKonversi}>
          {bisaUbah ? (
            <DropdownWarisan
              label={EGNPI.jumlah}
              ambil={ambilMataUang}
              nilai={b.Currency}
              onPilih={(o) => {
                // Mata uang dipilih → konversi SEGERA (Pega: `change` dropdown),
                // atas baris yang SUDAH memuat pilihannya.
                const baru = { ...b, Currency: o.nama, CurrencyID: o.id }
                onUbah(baru)
                onKonversi(baru)
              }}
            />
          ) : (
            <Field label={EGNPI.jumlah} value={b.Currency} readOnly onChange={() => undefined} />
          )}
          <FieldAngka
            label=""
            value={b.Amount}
            desimal={DESIMAL_EGNPI.jumlah}
            readOnly={!bisaUbah}
            onChange={(v) => {
              onUbah({ ...b, Amount: v })
            }}
          />
        </PemicuUbah>

        {/* Baris `Amount in IDR` — mata uangnya TETAP `IDR`.
            ⭐ Di ekspor sel kiri `.pyTemplateRichTextEditor`: pemegang tempat
            yang hanya memperlihatkan satuannya. Karena itu ia hanya-baca di
            sini, bukan pilihan. */}
        <div className="trin__egnpi-pasangan">
          <Field label={EGNPI.jumlahIDR} value={EGNPI.satuanIDR} readOnly onChange={() => undefined} />
          {/* ⭐ Nilainya DAPAT diubah di ekspor sekalipun
              `SetAmountConversion` menimpanya — dibiarkan begitu: menguncinya
              menghapus jalan keluar pemakai ketika mata uangnya nol di grid
              Rate of Exchange. */}
          <FieldAngka
            label=""
            value={b.AmountIDR}
            desimal={DESIMAL_EGNPI.jumlahIDR}
            readOnly={!bisaUbah}
            onChange={(v) => {
              onUbah({ ...b, AmountIDR: v })
            }}
          />
        </div>

        {/* ⚠️ `Proportion %` dapat diubah di ekspor, tetapi `Update Total`
            MENIMPANYA setiap kali ditekan. Itu bunyi ekspornya. */}
        {/* ⭐ PERSEN — `%` ditempel saat tidak sedang diketik. */}
        <FieldAngka
          label={EGNPI.proporsi}
          value={b.Proportion}
          desimal={DESIMAL_EGNPI.proporsi}
          persen
          readOnly={!bisaUbah}
          onChange={(v) => {
            onUbah({ ...b, Proportion: v })
          }}
        />

        {/* ⛔ RALAT 7 Oktober 2026 — medan ini sempat dikunci mati atas dasar
            `pyReadOnly` telanjang; `pyReadOnlyCondition = ViewState = 1` yang
            berlaku, jadi ia aktif di mode Edit. */}
        <Field
          label={EGNPI.keterangan}
          value={b.Note}
          readOnly={!bisaUbah}
          onChange={(v) => {
            onUbah({ ...b, Note: v })
          }}
        />
      </div>
    </Bagian>
  )
}

export default function TabEgnpi({
  baris,
  kurs,
  retensi,
  opsiAwal,
  mode = 'lihat',
}: {
  baris: readonly BarisEgnpi[]
  kurs: readonly KursLimitNP[]
  /** `TreatyIn.Retention` — baris pertamanya memberi mata uang baris baru. */
  retensi: readonly { Currency: string; CurrencyID: string }[]
  /** Isi dropdown; bila nol diberikan, tab ini mengambilnya sendiri. */
  opsiAwal?: OpsiLimits
  mode?: ModeForm
}) {
  // ⭐ PENAMPUNG HALAMAN (`../halaman.tsx`) — 7 Oktober 2026. Isian bertahan
  // saat pindah tab, dan tab lain membaca properti yang SAMA.
  //
  // ⛔ PERBAIKAN "input tiba-tiba hilang": efek `setRows([...baris])` atas
  // perubahan `baris` DICABUT. `baris` dibuat ulang (`.map`) tiap form
  // dirender, jadi efek itu mengembalikan baris ke data kontrak setiap kali
  // apa pun di form berubah. Kontrak lain = penampung dikosongkan + tab
  // dilahirkan ulang (`key` fieldset form).
  const [rows, setRows] = useProperti<BarisEgnpi[]>('EGNPI', () => [...baris])
  const [totalIDR, setTotalIDR] = useProperti('TotalEgnpiAmount', '')
  const [totalProporsi, setTotalProporsi] = useProperti('TotalEgnpiProportion', '')
  const [perMataUang, setPerMataUang] = useProperti<BarisTotal[]>('TotalEgnpiAmountNP', [])
  const [pesan, setPesan] = useState<string[]>([])
  const [gagal, setGagal] = useState('')

  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? { jenisTreaty: [], kelompokTreaty: [], mataUang: [] })

  const bisaUbah = mode === 'ubah'

  useEffect(() => {
    if (opsiAwal !== undefined) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setOpsi(o)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [opsiAwal])


  // ⛔ `DropdownWarisan` menjadikan `ambil` TANGGUNGAN efeknya, jadi fungsi
  // baru tiap render akan mengambil ulang daftarnya tanpa henti. Keduanya
  // dibekukan `useCallback`, dan daftarnya sudah ada di `opsi` — satu
  // pengambilan untuk seluruh baris, bukan dua per baris.
  const ambilKelompokTreaty = useCallback(() => Promise.resolve(opsi.kelompokTreaty), [opsi])
  const ambilMataUang = useCallback(() => Promise.resolve(opsi.mataUang), [opsi])

  function jalankan(aksi: AksiEgnpi, indeks = 0, dasar: readonly BarisEgnpi[] = rows) {
    setGagal('')
    hitungEgnpi({ aksi, egnpi: dasar, kurs, retensi, indeks })
      .then((h) => {
        setRows(h.egnpi)
        setTotalIDR(h.TotalEgnpiAmount)
        setTotalProporsi(h.TotalEgnpiProportion)
        setPerMataUang(h.TotalEgnpiAmountNP ?? [])
        setPesan(h.pesan ?? [])
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }

  return (
    <Panel judul={EGNPI.judul}>
      <div className="tl-rincian">
        <KepalaBagian
          judul={EGNPI.panel}
          jumlah={rows.length}
          aksi={
            bisaUbah && (
              <TombolTambah
                label={EGNPI.tambah}
                onClick={() => {
                  // `TreatyInNonAddItem(egnpi)` — mata uangnya mewarisi
                  // `TreatyIn.Retention(1)`; services yang menentukannya.
                  jalankan('tambah')
                }}
              />
            )
          }
        />

        {rows.length === 0 ? (
          <Kosong pesan={EGNPI.petunjukKosong} />
        ) : (
          <div className="tl-daftar">
            {rows.map((b, i) => (
              <KartuLipat
                key={i}
                nomor={i + 1}
                judul={b.TreatyGroup}
                judulKosong={EGNPI.barisBaru}
                bukaAwal={b.TreatyGroup === ''}
                meta={
                  <>
                    <Chip label={EGNPI.asAt} nilai={b.AsDate} />
                    <Chip label={EGNPI.mataUang} nilai={b.Currency} />
                    <Chip label={EGNPI.jumlah} nilai={tampil(b.Amount, DESIMAL_EGNPI.jumlah)} />
                    <Chip label={EGNPI.jumlahIDR} nilai={tampil(b.AmountIDR, DESIMAL_EGNPI.jumlahIDR)} />
                    <Chip label={EGNPI.proporsi} nilai={tampil(b.Proportion, DESIMAL_EGNPI.proporsi)} />
                  </>
                }
                aksi={
                  bisaUbah && (
                    <TombolHapus
                      label={EGNPI.hapus}
                      labelAkses={`${EGNPI.hapus} ${EGNPI.treatyGroup} ${i + 1}`}
                      onClick={() => {
                        jalankan('hapus', i)
                      }}
                    />
                  )
                }
              >
                <RincianBaris
                  b={b}
                  bisaUbah={bisaUbah}
                  ambilMataUang={ambilMataUang}
                  ambilKelompokTreaty={ambilKelompokTreaty}
                  onUbah={(baru) => {
                    setRows(rows.map((x, j) => (j === i ? baru : x)))
                  }}
                  onKonversi={(baru) => {
                    jalankan('konversi', i, baru === undefined ? rows : rows.map((x, j) => (j === i ? baru : x)))
                  }}
                />
              </KartuLipat>
            ))}
          </div>
        )}

        {pesan.length > 0 && (
          <ul className="tl-pesan" role="alert">
            {pesan.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        )}
        {gagal !== '' && (
          <p className="tl-pesan" role="alert">
            {gagal}
          </p>
        )}

        <Bagian judul={EGNPI.totalPerMataUang}>
          <GridTotal baris={perMataUang} />
        </Bagian>

        <Bagian
          judul={EGNPI.totalIDR}
          aksi={
            bisaUbah && (
              <div className="trin__aksi">
                <button
                  type="button"
                  className="btn btn--primary btn--sm"
                  onClick={() => {
                    jalankan('total')
                  }}
                >
                  {EGNPI.perbaruiTotal}
                </button>
                <button
                  type="button"
                  className="btn btn--sm"
                  onClick={() => {
                    jalankan('nilai')
                  }}
                >
                  {EGNPI.perbaruiNilai}
                </button>
              </div>
            )
          }
        >
          <div className="trin__egnpi-total">
            <Field
              label={`${EGNPI.totalIDR} (${EGNPI.satuanIDR})`}
              value={tampil(totalIDR, DESIMAL_EGNPI.totalIDR)}
              readOnly
              onChange={() => undefined}
            />
            <Field
              label={EGNPI.totalProporsi}
              value={tampil(totalProporsi, DESIMAL_EGNPI.totalProporsi)}
              readOnly
              onChange={() => undefined}
            />
          </div>
        </Bagian>
      </div>
    </Panel>
  )
}
